import "@testing-library/jest-dom";
import {screen} from "@testing-library/react";
import {render} from "@/app/ui/dashboard/utils/test-render";
import useWebSocket, {ReadyState} from "react-use-websocket";

jest.mock("react-use-websocket");
jest.mock("@/app/lib/data");

import {getInfraCapabilities} from "@/app/lib/data";
import ClusterTabs from "./cluster-tabs";

const mockedGetInfraCapabilities = getInfraCapabilities as jest.MockedFunction<typeof getInfraCapabilities>;
const mockedUseWebSocket = useWebSocket as jest.MockedFunction<typeof useWebSocket>;

beforeEach(() => {
  mockedUseWebSocket.mockReturnValue({
    sendJsonMessage: jest.fn(),
    lastJsonMessage: null,
    readyState: ReadyState.CONNECTING,
  } as unknown as ReturnType<typeof useWebSocket>);
});

afterEach(() => {
  jest.resetAllMocks();
});

describe("ClusterTabs", () => {
  it("shows only the Docker tab in a docker-only environment", async () => {
    mockedGetInfraCapabilities.mockResolvedValue({
      docker: {installed: true, version: "v1.10.10"},
      vsphere: {installed: false, version: ""},
    });
    render(<ClusterTabs/>);
    expect(await screen.findByText("Docker Clusters")).toBeInTheDocument();
    expect(screen.queryByText("vSphere Clusters")).not.toBeInTheDocument();
  });

  it("shows only the vSphere tab in a vsphere-only environment", async () => {
    mockedGetInfraCapabilities.mockResolvedValue({
      docker: {installed: false, version: ""},
      vsphere: {installed: true, version: "v1.12.0"},
    });
    render(<ClusterTabs/>);
    expect(await screen.findByText("vSphere Clusters")).toBeInTheDocument();
    expect(screen.queryByText("Docker Clusters")).not.toBeInTheDocument();
  });

  it("shows both tabs in a mixed environment", async () => {
    mockedGetInfraCapabilities.mockResolvedValue({
      docker: {installed: true, version: "v1.10.10"},
      vsphere: {installed: true, version: "v1.12.0"},
    });
    render(<ClusterTabs/>);
    expect(await screen.findByText("Docker Clusters")).toBeInTheDocument();
    expect(screen.getByText("vSphere Clusters")).toBeInTheDocument();
  });

  it("shows a clear message and no provider tab when neither is installed", async () => {
    mockedGetInfraCapabilities.mockResolvedValue({
      docker: {installed: false, version: ""},
      vsphere: {installed: false, version: ""},
    });
    render(<ClusterTabs/>);
    expect(await screen.findByText(/no supported infrastructure provider detected/i)).toBeInTheDocument();
    expect(screen.queryByText("Docker Clusters")).not.toBeInTheDocument();
    expect(screen.queryByText("vSphere Clusters")).not.toBeInTheDocument();
  });
});
