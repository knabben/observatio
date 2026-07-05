'use client';

import {useCallback, useEffect, useRef, useState} from 'react';
import {ReadyState} from 'react-use-websocket';
import {byMetadataName, receiveAndPopulate, sendInitialRequest, WebSocket} from '@/app/lib/websocket';

/**
 * Explicit lifecycle for a live resource view.
 *
 * - `connecting` : socket opening / awaiting first frame
 * - `ready`      : data present
 * - `empty`      : connected, zero items
 * - `error`      : connection failed or reconnects exhausted
 */
export type ChannelState = 'connecting' | 'ready' | 'empty' | 'error';

/** Bounded time to await the first data frame before resolving out of `connecting`. */
export const DATA_TIMEOUT_MS = 10_000;

interface ResourceStream<T> {
  state: ChannelState;
  items: T[];
  retry: () => void;
}

/**
 * Subscribes to the live resource stream for `objectType` and exposes an explicit
 * state machine so screens never hang on a silent socket or wipe a populated list on
 * an empty frame.
 */
export function useResourceStream<T extends {metadata?: {name?: string}}>(
  objectType: string,
): ResourceStream<T> {
  const [items, setItems] = useState<T[]>([]);
  const [state, setState] = useState<ChannelState>('connecting');
  const [nonce, setNonce] = useState(0);
  const receivedRef = useRef(false);

  const onReconnectStop = useCallback(() => {
    if (!receivedRef.current) setState('error');
  }, []);

  const {sendJsonMessage, lastJsonMessage, readyState} = WebSocket(undefined, {onReconnectStop});

  const retry = useCallback(() => {
    receivedRef.current = false;
    setItems([]);
    setState('connecting');
    setNonce((n) => n + 1);
  }, []);

  // (Re)send the initial subscription whenever the socket opens or we retry.
  useEffect(() => {
    sendInitialRequest(readyState, objectType, sendJsonMessage);
  }, [readyState, sendJsonMessage, objectType, nonce]);

  // Bounded resolution: if no data arrives in time, leave `connecting` for `empty`
  // (socket open) or `error` (never connected) rather than spinning forever.
  useEffect(() => {
    const timer = setTimeout(() => {
      if (!receivedRef.current) {
        setState((prev) =>
          prev !== 'connecting' ? prev : readyState === ReadyState.OPEN ? 'empty' : 'error',
        );
      }
    }, DATA_TIMEOUT_MS);
    return () => clearTimeout(timer);
  }, [nonce, readyState]);

  // Incoming frames. Null/empty frames are ignored (never clear the list).
  useEffect(() => {
    if (lastJsonMessage == null) return;
    setItems((prev) => {
      const next = (receiveAndPopulate(lastJsonMessage, prev as unknown[]) as T[]).sort(
        byMetadataName,
      );
      receivedRef.current = true;
      setState(next.length === 0 ? 'empty' : 'ready');
      return next;
    });
  }, [lastJsonMessage]);

  return {state, items, retry};
}
