/* eslint-disable @typescript-eslint/no-explicit-any */

import useWebSocket, {ReadyState} from "react-use-websocket";

export const WS_URL_WATCHER = typeof window !== 'undefined' ? `ws://${window.location.hostname}:8080/ws/watcher` : 'ws://localhost:8080/ws/watcher';
export const WS_URL_CHATBOT = typeof window !== 'undefined' ? `ws://${window.location.hostname}:8080/ws/analysis` : 'ws://localhost:8080/ws/analysis';

export type WSResponse = {
  type: string;
  data: any;
}

/**
 * A function that establishes and manages a WebSocket connection using the specified `useWebSocket` integration.
 * The WebSocket URL is defined by the `URL` parameter. The connection is not shared among multiple instances
 * and is configured to always attempt reconnection if the connection is lost.
 */
export function WebSocket(URL: string = WS_URL_WATCHER) {
  return useWebSocket(
    URL, {
      share: false,
      shouldReconnect: () => true
    },
  )
}

/**
 * Sends an initial WebSocket request when the connection is open.
 * This function is typically used to establish initial subscription or request specific data types
 * from the WebSocket server.
 */
export function sendInitialRequest(readyState: number, type: string, sendJsonMessage: any) {
  if (readyState === ReadyState.OPEN) {
    sendJsonMessage({type: type});
  }
}


/**
 * Processes a WebSocket response and updates the provided items list based on the response type.
 */
export enum WSOperationType {
  ADDED = "ADDED",
  MODIFIED = "MODIFIED",
  DELETED = "DELETED"
}

export function receiveAndPopulate(
  response: any,
  items: any[],
): any {
  if (!response?.data) {
    return [];
  }
  const isUpdateOperation = isItemUpdateOperation(response.type);
  if (isUpdateOperation) {
    return updateItemsList(items, response.data)
  } else {
    return items.filter(item => item.metadata?.name !== response.data.metadata?.name);
  }
}

// Extract operation type checking to a separate function
function isItemUpdateOperation(type: WSOperationType): boolean {
  return type === WSOperationType.ADDED || type === WSOperationType.MODIFIED;
}

/**
 * Updates a list of items by adding or replacing an item with a matching name.
 * If an item with the same name already exists in the list, it will be replaced
 * with the new item. Otherwise, the new item will be added to the list.
 */
function updateItemsList<T extends { metadata: {name: string} }>(items: T[], newItem: T): T[] {
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