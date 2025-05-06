import useWebSocket, {ReadyState} from "react-use-websocket";
import {WS_URL} from "@/app/lib/consts";

export type WSResponse = {
  type: string;
  data: any;
}

/**
 * A function that establishes and manages a WebSocket connection using the specified `useWebSocket` integration.
 * The WebSocket URL is defined by the `WS_URL` constant. The connection is not shared among multiple instances
 * and is configured to always attempt reconnection if the connection is lost.
 *
 * @return {Object} The WebSocket connection object provided by the `useWebSocket` integration.
 */
export function WebSocket() {
  return useWebSocket(
    WS_URL, {
      share: false,
      shouldReconnect: () => true
    },
  )
}

/**
 * Sends an initial WebSocket request when the connection is open.
 * This function is typically used to establish initial subscription or request specific data types
 * from the WebSocket server.
 *
 * @param readyState - The current state of the WebSocket connection (from ReadyState enum)
 * @param type - The type of data to request from the server
 * @param sendJsonMessage - The function to send JSON messages through WebSocket
 *
 * @example
 * ```typescript
 * sendInitialRequest(readyState, "cluster-infra", sendJsonMessage);
 * ```
 */
export function sendInitialRequest(readyState: number, type: string, sendJsonMessage: any) {
  if (readyState === ReadyState.OPEN) {
    sendJsonMessage({types: [type]});
  }
}


/**
 * Processes a WebSocket response and updates the provided items list based on the response type.
 *
 * @param {WSResponse | null} response - The WebSocket response object, or null if no response was received.
 * @param {T[]} items - The current list of items to be updated.
 * @param {(items: T[]) => void} setItems - A callback function used to update the list of items with the new state.
 * @return {void} This function does not return a value.
 */
export enum WSOperationType {
  ADDED = "ADDED",
  MODIFIED = "MODIFIED",
  DELETED = "DELETED"
}

export interface NamedItem {
  name: string;
}

export function receiveAndPopulate<T extends NamedItem>(
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
    return items.filter(item => item.name !== response.data.name);
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
 *
 * @template T The type of the items in the list, which must include a `name` property.
 * @param {T[]} items The original list of items.
 * @param {T} newItem The new item to be added or used to replace an existing item.
 * @return {T[]} A new array with the updated list of items.
 */
function updateItemsList<T extends { name: string }>(items: T[], newItem: T): T[] {
  const existingItemIndex = items.findIndex(item => item.name === newItem.name);

  if (existingItemIndex !== -1) {
    return [
      ...items.slice(0, existingItemIndex),
      newItem,
      ...items.slice(existingItemIndex + 1)
    ];
  }

  return [...items, newItem];
}