'use strict'

import {sendInitialRequest} from "./websocket";
import {ReadyState} from "react-use-websocket";
import {describe, it, expect, jest} from "@jest/globals";

describe("sendInitialRequest", () => {

  it("should send the initial request when WebSocket connection is open", () => {
    const mockSendJsonMessage = jest.fn();
    const readyState = ReadyState.OPEN;
    const type = "cluster-infra";
    sendInitialRequest(readyState, type, mockSendJsonMessage);
    expect(mockSendJsonMessage).toHaveBeenCalledWith({types: [type]});
  });

  it("should not send the initial request when WebSocket connection is not open", () => {
    const mockSendJsonMessage = jest.fn();
    const readyState = ReadyState.CLOSED;
    const type = "cluster-infra";
    sendInitialRequest(readyState, type, mockSendJsonMessage);
    expect(mockSendJsonMessage).not.toBeCalled();
  });

  it("should handle different data types correctly", () => {
    const mockSendJsonMessage = jest.fn();
    const readyState = ReadyState.OPEN;
    const type = "user-data";
    sendInitialRequest(readyState, type, mockSendJsonMessage);
    expect(mockSendJsonMessage).toHaveBeenCalledWith({types: [type]});
  });
});