import {renderHook, act} from "@testing-library/react";
import useWebSocket, {ReadyState} from "react-use-websocket";
import {useResourceStream, DATA_TIMEOUT_MS} from "./resource-hooks";

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

describe("useResourceStream", () => {
  beforeEach(() => {
    jest.useFakeTimers();
  });

  afterEach(() => {
    jest.useRealTimers();
    jest.resetAllMocks();
  });

  it("starts in the connecting state", () => {
    const state: MockState = {lastJsonMessage: null, readyState: ReadyState.CONNECTING};
    mockSocket(state);

    const {result} = renderHook(() => useResourceStream("cluster"));

    expect(result.current.state).toBe("connecting");
    expect(result.current.items).toEqual([]);
  });

  it("resolves to ready when a populated frame arrives", () => {
    const state: MockState = {lastJsonMessage: null, readyState: ReadyState.OPEN};
    mockSocket(state);

    const {result, rerender} = renderHook(() => useResourceStream("cluster"));

    state.lastJsonMessage = {type: "ADDED", data: {metadata: {name: "c1"}}};
    rerender();

    expect(result.current.state).toBe("ready");
    expect(result.current.items).toEqual([{metadata: {name: "c1"}}]);
  });

  it("resolves to empty after the data timeout when the socket is open with no data", () => {
    const state: MockState = {lastJsonMessage: null, readyState: ReadyState.OPEN};
    mockSocket(state);

    const {result} = renderHook(() => useResourceStream("cluster"));

    act(() => {
      jest.advanceTimersByTime(DATA_TIMEOUT_MS);
    });

    expect(result.current.state).toBe("empty");
  });

  it("resolves to error after the data timeout when the socket never connects", () => {
    const state: MockState = {lastJsonMessage: null, readyState: ReadyState.CONNECTING};
    mockSocket(state);

    const {result} = renderHook(() => useResourceStream("cluster"));

    act(() => {
      jest.advanceTimersByTime(DATA_TIMEOUT_MS);
    });

    expect(result.current.state).toBe("error");
  });

  it("moves to error when reconnect attempts are exhausted", () => {
    const state: MockState = {lastJsonMessage: null, readyState: ReadyState.CONNECTING};
    mockSocket(state);

    const {result} = renderHook(() => useResourceStream("cluster"));

    act(() => {
      state.onReconnectStop?.(8);
    });

    expect(result.current.state).toBe("error");
  });

  it("ignores an empty/malformed frame instead of clearing a populated list", () => {
    const state: MockState = {lastJsonMessage: null, readyState: ReadyState.OPEN};
    mockSocket(state);

    const {result, rerender} = renderHook(() => useResourceStream("cluster"));

    state.lastJsonMessage = {type: "ADDED", data: {metadata: {name: "c1"}}};
    rerender();
    expect(result.current.state).toBe("ready");

    state.lastJsonMessage = {type: "MODIFIED"};
    rerender();

    expect(result.current.state).toBe("ready");
    expect(result.current.items).toEqual([{metadata: {name: "c1"}}]);
  });

  it("retry resets to connecting and clears items", () => {
    const state: MockState = {lastJsonMessage: null, readyState: ReadyState.OPEN};
    mockSocket(state);

    const {result, rerender} = renderHook(() => useResourceStream("cluster"));

    state.lastJsonMessage = {type: "ADDED", data: {metadata: {name: "c1"}}};
    rerender();
    expect(result.current.state).toBe("ready");

    act(() => {
      result.current.retry();
    });

    expect(result.current.state).toBe("connecting");
    expect(result.current.items).toEqual([]);
  });
});
