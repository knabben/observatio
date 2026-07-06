/* eslint-disable @typescript-eslint/no-explicit-any */

import useWebSocket, {ReadyState} from "react-use-websocket";
import {WS_URL_WATCHER, WS_URL_CHATBOT} from "@/app/lib/config";

export {WS_URL_WATCHER, WS_URL_CHATBOT};

export type WSResponse = {
  type: string;
  data: any;
}

type WebSocketExtraOptions = {
  onReconnectStop?: (attempts: number) => void;
};

/**
 * Establishes and manages a WebSocket connection via `react-use-websocket`.
 * Reconnection is bounded: at most 8 attempts with exponential backoff (capped at 30s),
 * after which `onReconnectStop` fires so consumers can surface a terminal error state
 * instead of looping forever.
 */
export function WebSocket(URL: string = WS_URL_WATCHER, options: WebSocketExtraOptions = {}) {
  return useWebSocket(
    URL, {
      share: false,
      shouldReconnect: () => true,
      reconnectAttempts: 8,
      reconnectInterval: (attempt: number) => Math.min(1000 * 2 ** attempt, 30000),
      onReconnectStop: options.onReconnectStop,
    },
  )
}

/**
 * Sends an initial WebSocket request when the connection is open.
 */
export function sendInitialRequest(readyState: number, type: string, sendJsonMessage: any) {
  if (readyState === ReadyState.OPEN) {
    sendJsonMessage({type});
  }
}

export enum WSOperationType {
  ADDED = "ADDED",
  MODIFIED = "MODIFIED",
  DELETED = "DELETED"
}

/**
 * Processes a WebSocket response and returns the updated items list.
 *
 * IMPORTANT: an empty/malformed frame (no `.data`) is treated as a no-op and the
 * current list is returned UNCHANGED — a keepalive or partial frame must never wipe
 * an already-populated list.
 */
export function receiveAndPopulate(
  response: any,
  items: any[],
): any[] {
  if (!response?.data) {
    return items;
  }
  if (isItemUpdateOperation(response.type)) {
    return updateItemsList(items, response.data)
  }
  return items.filter(item => item.metadata?.name !== response.data.metadata?.name);
}

function isItemUpdateOperation(type: string): boolean {
  return type === WSOperationType.ADDED || type === WSOperationType.MODIFIED;
}

/**
 * Adds or replaces an item with a matching `metadata.name`.
 */
function updateItemsList<T extends { metadata?: {name?: string} }>(items: T[], newItem: T): T[] {
  const existingItemIndex = items.findIndex(item => item.metadata?.name === newItem.metadata?.name);

  if (existingItemIndex !== -1) {
    return [
      ...items.slice(0, existingItemIndex),
      newItem,
      ...items.slice(existingItemIndex + 1)
    ];
  }

  return [...items, newItem];
}

/** Stable name-based comparator that tolerates missing `metadata.name`. */
export function byMetadataName<T extends { metadata?: {name?: string} }>(a: T, b: T): number {
  return (a?.metadata?.name ?? '').localeCompare(b?.metadata?.name ?? '');
}
