# Contract: Data Channels (WebSocket + REST)

The backend contract is **unchanged**; this documents the shapes the frontend consumes and the hardened
handling applied to them. No backend modifications are in scope.

---

## WebSocket — live resource stream

**Endpoint**: `WS_URL` (from `AppConfig`, was hardcoded). One connection per resource object type.

**Inbound frame** (existing shape):

```ts
interface WSResponse {
  type: string;          // operation: create/update/delete (raw string)
  data?: unknown[];      // resource items; MAY be absent on keepalive/empty frames
}
```

**Frontend handling contract**:
- A frame with no `data` is a **no-op** (keepalive) — MUST NOT clear the current list. (FR-005)
- `create`/`update` operations merge/replace by stable `metadata.name`; `delete` removes by name.
- Connection lifecycle drives `ChannelState` (connecting/ready/empty/error). (FR-003)
- Reconnect is bounded: 8 attempts, exponential backoff to 30s, then terminal `error`. (FR-007)
- No item-level access assumes field presence; sort tolerates missing `metadata.name`. (FR-001, FR-008)

---

## WebSocket — AI troubleshooting chat

**Endpoint**: `WS_URL_CHATBOT` (from `AppConfig`).

**Outbound**: request carries a client message id (generated via `uuid` v4 — secure-context-safe) plus the
affected resource's condition context (AI stays grounded — Constitution IV).

**Frontend handling contract**:
- Send only when `readyState === OPEN`; otherwise surface a failure and reset in-progress state. (FR-035)
- Message content rendered as **safe plain text** (no `dangerouslySetInnerHTML`, valid element nesting). (FR-032)
- Conversation scoped to the single session (Constitution IV); no cross-session persistence.

---

## REST — dashboard reads

**Base**: `API_URL` (from `AppConfig`). Endpoints (existing): cluster hierarchy, summary, component
versions, cluster classes, etc.

**Frontend handling contract**:
- Every fetch checks `res.ok` before `res.json()`; non-ok → `error` state with message, never parsed as
  data. (FR-004)
- Network/parse failures are caught and surfaced, not thrown uncaught. (FR-004)
- Empty results render `EmptyState` (FR-002); topology handles `{nodes: [], edges: []}` gracefully (FR-012).
