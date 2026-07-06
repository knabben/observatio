import "@testing-library/jest-dom";
import {screen, waitFor} from "@testing-library/react";
import {render} from "@/app/ui/dashboard/utils/test-render";

jest.mock("@/app/lib/data");

import {
  getClusterHierarchy,
  getClusterSummary,
  getComponentsVersion,
  getClusterClasses,
} from "@/app/lib/data";

import ClusterHierarchy from "./clusterhierarchy";
import ClusterSummary from "./clustersummary";
import ClusterVersions from "./clusterversions";
import ClusterClassLister from "./clusterclass";

const mockedGetClusterHierarchy = getClusterHierarchy as jest.MockedFunction<typeof getClusterHierarchy>;
const mockedGetClusterSummary = getClusterSummary as jest.MockedFunction<typeof getClusterSummary>;
const mockedGetComponentsVersion = getComponentsVersion as jest.MockedFunction<typeof getComponentsVersion>;
const mockedGetClusterClasses = getClusterClasses as jest.MockedFunction<typeof getClusterClasses>;

// xyflow measures its container via ResizeObserver, which jsdom does not implement.
beforeAll(() => {
  (global as unknown as {ResizeObserver: unknown}).ResizeObserver = class {
    observe() {}
    unobserve() {}
    disconnect() {}
  };
});

afterEach(() => {
  jest.resetAllMocks();
});

describe("ClusterHierarchy", () => {
  it("renders an empty state when there is no topology data", async () => {
    mockedGetClusterHierarchy.mockResolvedValue({nodes: [], edges: []});
    render(<ClusterHierarchy/>);
    expect(await screen.findByText(/no cluster topology found/i)).toBeInTheDocument();
  });

  it("renders the topology without crashing when nodes are present", async () => {
    mockedGetClusterHierarchy.mockResolvedValue({
      nodes: [{id: "1", type: "default", data: {label: "cluster-a"}, position: {x: 0, y: 0}, style: {background: "#fff", color: "#000", border: "1px"}}],
      edges: [],
    });
    render(<ClusterHierarchy/>);
    await waitFor(() => expect(screen.queryByText(/no cluster topology found/i)).not.toBeInTheDocument());
  });

  it("renders an error message when the fetch fails", async () => {
    mockedGetClusterHierarchy.mockRejectedValue(new Error("network down"));
    render(<ClusterHierarchy/>);
    expect(await screen.findByText(/failed to load cluster topology/i)).toBeInTheDocument();
  });
});

describe("ClusterSummary", () => {
  it("renders zero-value counts as legible zeros, not blanks", async () => {
    mockedGetClusterSummary.mockResolvedValue({
      clusterProvisioned: 0,
      clusterFailed: 0,
      machineProvisioned: 0,
      machineFailed: 0,
      machineDeploymentProvisioned: 0,
      machineDeploymentFailed: 0,
    });
    render(<ClusterSummary/>);
    expect(await screen.findByText("Cluster running")).toBeInTheDocument();
    expect(screen.getAllByText("0").length).toBeGreaterThan(0);
  });

  it("tolerates a partial response with missing fields", async () => {
    mockedGetClusterSummary.mockResolvedValue({clusterProvisioned: 3});
    render(<ClusterSummary/>);
    expect(await screen.findByText("Cluster running")).toBeInTheDocument();
  });
});

describe("ClusterVersions", () => {
  it("renders an empty state for zero components", async () => {
    mockedGetComponentsVersion.mockResolvedValue([]);
    render(<ClusterVersions/>);
    expect(await screen.findByText(/no components found/i)).toBeInTheDocument();
  });

  it("renders a partial component list without crashing", async () => {
    mockedGetComponentsVersion.mockResolvedValue([{name: "capi", kind: "Deployment", version: "v1.9.0"}]);
    render(<ClusterVersions/>);
    expect(await screen.findByText("capi")).toBeInTheDocument();
  });
});

describe("ClusterClassLister", () => {
  it("renders an empty state for zero cluster classes", async () => {
    mockedGetClusterClasses.mockResolvedValue([]);
    render(<ClusterClassLister/>);
    expect(await screen.findByText(/no cluster classes found/i)).toBeInTheDocument();
  });

  it("renders a partial cluster class (missing name/conditions, bigint generation) without crashing", async () => {
    mockedGetClusterClasses.mockResolvedValue([{namespace: "default", generation: BigInt(2)}]);
    render(<ClusterClassLister/>);
    expect(await screen.findByText("default")).toBeInTheDocument();
    expect(screen.getByText("2")).toBeInTheDocument();
  });

  it("renders status as a read-only badge, not an interactive toggle", async () => {
    mockedGetClusterClasses.mockResolvedValue([
      {name: "cc1", namespace: "default", generation: BigInt(1), conditions: [{type: "Ready", status: "True"}]},
    ]);
    const {container} = render(<ClusterClassLister/>);
    expect(await screen.findByText("Ready")).toBeInTheDocument();
    // A Chip renders a hidden checkbox input; a read-only Badge renders none.
    expect(container.querySelector("input")).toBeNull();
  });
});
