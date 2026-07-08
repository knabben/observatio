import {renderHook, act} from "@testing-library/react";
import useWebSocket, {ReadyState} from "react-use-websocket";
import {useDay2Ops, Day2OpsData} from "./use-day2-ops";

jest.mock("react-use-websocket");

const mockedUseWebSocket = useWebSocket as jest.MockedFunction<typeof useWebSocket>;

type MockState = {
  lastJsonMessage: unknown;
  readyState: ReadyState;
  onReconnectStop?: (attempts: number) => void;
};

function mockSocket(state: MockState) {
  mockedUseWebSocket.mockImplementation((_url, options) => {
    state.onReconnectStop = options?.onReconnectStop as (attempts: number) => void;
    return {
      sendJsonMessage: jest.fn(),
      lastJsonMessage: state.lastJsonMessage,
      readyState: state.readyState,
    } as unknown as ReturnType<typeof useWebSocket>;
  });
}

const sampleData: Day2OpsData = {
  rollups: [{category: "cluster", healthy: 2, degraded: 0, failed: 1, unavailable: false}],
  debugPaths: [],
  risks: [],
  severities: [],
  sourceUnavailable: false,
};

describe("useDay2Ops", () => {
  afterEach(() => jest.resetAllMocks());

  it("starts unloaded with empty data", () => {
    const state: MockState = {lastJsonMessage: null, readyState: ReadyState.CONNECTING};
    mockSocket(state);

    const {result} = renderHook(() => useDay2Ops());

    expect(result.current.loaded).toBe(false);
    expect(result.current.data.rollups).toEqual([]);
  });

  it("replaces the full data snapshot on a day2ops frame", () => {
    const state: MockState = {lastJsonMessage: null, readyState: ReadyState.OPEN};
    mockSocket(state);

    const {result, rerender} = renderHook(() => useDay2Ops());

    state.lastJsonMessage = {type: "MODIFIED", event: "day2ops", data: sampleData};
    rerender();

    expect(result.current.loaded).toBe(true);
    expect(result.current.data).toEqual(sampleData);
  });

  it("ignores frames for other event types", () => {
    const state: MockState = {lastJsonMessage: null, readyState: ReadyState.OPEN};
    mockSocket(state);

    const {result, rerender} = renderHook(() => useDay2Ops());

    state.lastJsonMessage = {type: "ADDED", event: "cluster", data: {metadata: {name: "c1"}}};
    rerender();

    expect(result.current.loaded).toBe(false);
  });

  it("marks the data source unavailable when reconnect attempts are exhausted", () => {
    const state: MockState = {lastJsonMessage: null, readyState: ReadyState.CONNECTING};
    mockSocket(state);

    const {result} = renderHook(() => useDay2Ops());

    act(() => {
      state.onReconnectStop?.(8);
    });

    expect(result.current.data.sourceUnavailable).toBe(true);
  });
});
