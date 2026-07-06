import "@testing-library/jest-dom";
import {screen} from "@testing-library/react";
import {render} from "@/app/ui/dashboard/utils/test-render";
import useWebSocket, {ReadyState} from "react-use-websocket";

jest.mock("react-use-websocket");
jest.mock("@/app/lib/data");

import {getInfraCapabilities} from "@/app/lib/data";
import MachineTabs from "./machine-tabs";

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

describe("MachineTabs", () => {
  it("shows only the Docker tab in a docker-only environment", async () => {
    mockedGetInfraCapabilities.mockResolvedValue({
      docker: {installed: true, version: "v1.10.10"},
      vsphere: {installed: false, version: ""},
    });
    render(<MachineTabs/>);
    expect(await screen.findByText("Docker Machines")).toBeInTheDocument();
    expect(screen.queryByText("vSphere Machines")).not.toBeInTheDocument();
  });

  it("shows both tabs in a mixed environment", async () => {
    mockedGetInfraCapabilities.mockResolvedValue({
      docker: {installed: true, version: "v1.10.10"},
      vsphere: {installed: true, version: "v1.12.0"},
    });
    render(<MachineTabs/>);
    expect(await screen.findByText("Docker Machines")).toBeInTheDocument();
    expect(screen.getByText("vSphere Machines")).toBeInTheDocument();
  });

  it("shows a clear message and no provider tab when neither is installed", async () => {
    mockedGetInfraCapabilities.mockResolvedValue({
      docker: {installed: false, version: ""},
      vsphere: {installed: false, version: ""},
    });
    render(<MachineTabs/>);
    expect(await screen.findByText(/no supported infrastructure provider detected/i)).toBeInTheDocument();
    expect(screen.queryByText("Docker Machines")).not.toBeInTheDocument();
    expect(screen.queryByText("vSphere Machines")).not.toBeInTheDocument();
  });
});
