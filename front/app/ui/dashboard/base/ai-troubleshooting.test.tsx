import "@testing-library/jest-dom";
import {fireEvent, screen} from "@testing-library/react";
import useWebSocket, {ReadyState} from "react-use-websocket";
import {render} from "@/app/ui/dashboard/utils/test-render";
import {FilterItemsByName} from "@/app/dashboard/utils";
import Search from "@/app/ui/dashboard/search";
import ClusterClassLister from "@/app/ui/dashboard/components/dashboard/clusterclass";
import AITroubleshooting from "./ai-troubleshooting";

jest.mock("react-use-websocket");
jest.mock("@/app/lib/data");

const mockedUseWebSocket = useWebSocket as jest.MockedFunction<typeof useWebSocket>;

// jsdom does not implement scrollTo/ResizeObserver; ScrollArea uses both on mount.
beforeAll(() => {
  window.HTMLElement.prototype.scrollTo = jest.fn();
  window.Element.prototype.scrollTo = jest.fn();
  (global as unknown as {ResizeObserver: unknown}).ResizeObserver = class {
    observe() {}
    unobserve() {}
    disconnect() {}
  };
});

function mockSocket(lastJsonMessage: unknown, readyState: ReadyState = ReadyState.OPEN) {
  mockedUseWebSocket.mockReturnValue({
    sendJsonMessage: jest.fn(),
    lastJsonMessage,
    readyState,
  } as unknown as ReturnType<typeof useWebSocket>);
}

describe("Search", () => {
  it("filters the visible list as the user types (client-side, by name)", () => {
    const items = [{metadata: {name: "alpha"}}, {metadata: {name: "beta"}}];
    expect(FilterItemsByName("al", items)).toEqual([{metadata: {name: "alpha"}}]);
    expect(FilterItemsByName("", items)).toEqual(items);
    expect(FilterItemsByName("zzz", items)).toEqual([]);
  });

  it("calls onChange with the typed value, not a static pick-one selection", () => {
    const onChange = jest.fn();
    render(<Search value="" onChange={onChange}/>);
    const input = screen.getByRole("textbox", {name: /search by name/i});
    fireEvent.change(input, {target: {value: "clu"}});
    expect(onChange).toHaveBeenCalledWith("clu");
  });
});

describe("ClusterClass status display", () => {
  it("renders status as a read-only badge, not an interactive toggle", async () => {
    const {getClusterClasses} = jest.requireMock("@/app/lib/data");
    getClusterClasses.mockResolvedValue([
      {name: "cc1", namespace: "default", generation: BigInt(1), conditions: [{type: "Ready", status: "True"}]},
    ]);
    const {container} = render(<ClusterClassLister/>);
    expect(await screen.findByText("Ready")).toBeInTheDocument();
    // A Chip renders a hidden checkbox input; a read-only Badge renders none.
    expect(container.querySelector("input")).toBeNull();
  });
});

describe("AITroubleshooting panel", () => {
  beforeEach(() => {
    mockSocket(null);
  });

  it("renders AI/user message content as safe plain text, never parsed as HTML", () => {
    mockSocket({
      id: "1",
      type: "chatbot",
      agent_id: "cluster-agent",
      actor: "agent",
      timestamp: "now",
      content: '<img src=x onerror="window.__pwned=true">',
    });
    const {container} = render(
      <AITroubleshooting objectType="cluster" objectName="c1" objectNamespace="default" conditions={[]}/>,
    );
    expect(screen.getByText('<img src=x onerror="window.__pwned=true">')).toBeInTheDocument();
    expect(container.querySelector("img")).toBeNull();
  });

  it("expands and collapses the panel via the same reversible control", () => {
    render(<AITroubleshooting objectType="cluster" objectName="c1" objectNamespace="default" conditions={[]}/>);
    const expandButton = screen.getByRole("button", {name: /expand ai troubleshooting panel/i});
    fireEvent.click(expandButton);
    const collapseButton = screen.getByRole("button", {name: /collapse ai troubleshooting panel/i});
    fireEvent.click(collapseButton);
    expect(screen.getByRole("button", {name: /expand ai troubleshooting panel/i})).toBeInTheDocument();
  });
});
