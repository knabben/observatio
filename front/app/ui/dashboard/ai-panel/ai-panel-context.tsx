'use client';

import React, {createContext, useCallback, useContext, useEffect, useMemo, useState} from 'react';

export interface WSRequest {
  id: string;
  type: string;
  agent_id: string;
  content: string;
  timestamp: string;
  actor: string;
  /** Streaming marker: absent/empty for a one-shot message, "delta" to append to the in-progress
   *  message with the same id, or "done" once that streamed reply has finished. */
  event?: string;
}

/** Whatever object a detail screen is currently showing, registered for AI auto-context. */
export interface ObjectContext {
  kind: string;
  name: string;
  namespace: string;
  status: string;
  keySpecFields: Record<string, string>;
}

/** Turns an ObjectContext into the pre-filled query text (FR-006). */
export function formatObjectContext(ctx: ObjectContext): string {
  const specLine = Object.entries(ctx.keySpecFields)
    .map(([k, v]) => `${k}=${v}`)
    .join(', ');
  const specSuffix = specLine ? ` Key spec fields: ${specLine}.` : '';
  return `On ${ctx.kind} "${ctx.name}" (namespace ${ctx.namespace}): ${ctx.status}.${specSuffix}`;
}

interface AIPanelContextValue {
  isOpen: boolean;
  open: (prefill?: string) => void;
  close: () => void;
  messages: WSRequest[];
  setMessages: React.Dispatch<React.SetStateAction<WSRequest[]>>;
  currentObjectContext: ObjectContext | null;
  setCurrentObjectContext: (ctx: ObjectContext | null) => void;
  queryField: string;
  setQueryField: (value: string) => void;
  queryFieldTouched: boolean;
}

const AIPanelContext = createContext<AIPanelContextValue | null>(null);

/**
 * App-wide AI troubleshooting panel state: open/closed, the conversation, whichever object is
 * currently in view, and the editable query field. Mounted once in the dashboard layout so it
 * survives client-side navigation; never persisted beyond the browser session (Constitution IV).
 */
export function AIPanelProvider({children}: {children: React.ReactNode}) {
  const [isOpen, setIsOpen] = useState(false);
  const [messages, setMessages] = useState<WSRequest[]>([]);
  const [currentObjectContext, setCurrentObjectContext] = useState<ObjectContext | null>(null);
  const [queryField, setQueryFieldState] = useState('');
  const [queryFieldTouched, setQueryFieldTouched] = useState(false);

  // Auto-refresh the query field from whatever object is currently in view, but only while
  // the operator hasn't started editing it (FR-006, FR-007, FR-009).
  useEffect(() => {
    if (queryFieldTouched) return;
    setQueryFieldState(currentObjectContext ? formatObjectContext(currentObjectContext) : '');
  }, [currentObjectContext, queryFieldTouched]);

  const setQueryField = useCallback((value: string) => {
    setQueryFieldState(value);
    setQueryFieldTouched(true);
  }, []);

  const open = useCallback((prefill?: string) => {
    if (prefill !== undefined) {
      setQueryFieldState(prefill);
      setQueryFieldTouched(false);
    }
    setIsOpen(true);
  }, []);

  const close = useCallback(() => setIsOpen(false), []);

  const value = useMemo<AIPanelContextValue>(() => ({
    isOpen,
    open,
    close,
    messages,
    setMessages,
    currentObjectContext,
    setCurrentObjectContext,
    queryField,
    setQueryField,
    queryFieldTouched,
  }), [isOpen, open, close, messages, currentObjectContext, queryField, setQueryField, queryFieldTouched]);

  return <AIPanelContext.Provider value={value}>{children}</AIPanelContext.Provider>;
}

export function useAIPanel(): AIPanelContextValue {
  const ctx = useContext(AIPanelContext);
  if (!ctx) {
    throw new Error('useAIPanel must be used within an AIPanelProvider');
  }
  return ctx;
}
