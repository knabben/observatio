# Contract: Shared UI Building Blocks

The consolidated components under `front/app/ui/dashboard/shared/`. These are the internal UI contracts the
per-resource screens (clusters, machines, mds, and infra variants) compose. Behavior below is what tests
assert (Constitution V).

---

## `ObjectTable<T>`

```ts
interface ColumnDef<T> {
  header: string;
  render: (item: T) => React.ReactNode;   // MUST be null-safe
  align?: 'left' | 'center' | 'right';
  width?: number;
}

interface ObjectTableProps<T> {
  items: T[] | undefined;
  columns: ColumnDef<T>[];
  getRowKey: (item: T) => string;          // stable unique id — NEVER array index
  onSelect?: (item: T) => void;
  emptyLabel: string;                       // shown when items is empty/undefined
}
```

**Contract**:
- `items` empty or `undefined` → renders `EmptyState` with `emptyLabel` (never a header-only table). (FR-002)
- Rows keyed by `getRowKey`, not index. (FR-025)
- Wraps the table in a horizontal scroll container. (FR-013)
- If `onSelect` is set, the row's primary cell is a keyboard-focusable, `aria`-labeled button. (FR-018)
- Never throws on items with missing fields (each `render` is null-safe). (FR-001)

---

## `StatusIndicator`

```ts
interface StatusIndicatorProps {
  state: 'healthy' | 'notready' | 'unknown';
  label?: string;   // accessible text; defaults per state
}
```

**Contract**:
- `healthy` → solid positive color, **no animation**; `notready` → solid negative color, **no `processing`
  pulse**; `unknown` → neutral. (FR-020)
- Exposes an accessible name; status is not conveyed by color alone. (FR-017)

---

## `EmptyState` / `ErrorState`

```ts
interface EmptyStateProps { label: string; }
interface ErrorStateProps { message: string; onRetry?: () => void; }
```

**Contract**:
- `EmptyState` renders a labeled message. (FR-002)
- `ErrorState` renders the message and, when `onRetry` is provided, a keyboard-accessible retry control that
  re-enters the `connecting` state. (FR-004)

---

## `useResourceStream<T>` (hook)

```ts
function useResourceStream<T>(objectType: string): {
  state: 'connecting' | 'ready' | 'empty' | 'error';
  items: T[];
  retry: () => void;
};
```

**Contract**:
- Encapsulates the `ChannelState` machine (data-model §3). Resolves out of `connecting` within 10s. (FR-003)
- Bounded reconnect (8 attempts, exponential backoff), terminal `error` on `onReconnectStop`. (FR-007)
- Ignores empty/malformed frames — never clears a populated list. (FR-005)
- Replaces the four copy-pasted fetch hooks. (FR-023, FR-024)

---

## `BaseLister` / `ObjectDetails` (refactored)

**Contract**:
- `BaseLister` renders `connecting → CenteredLoader`, `empty → EmptyState`, `error → ErrorState(retry)`,
  `ready → ObjectTable`. Loader is vertically centered. (FR-003, FR-015)
- `ObjectDetails` renders a shared header (name/namespace/age via `DetailFieldDef`) + tabs; guards empty
  `tabs`; labels are accurate (`Age`, not `Created`). (FR-023, FR-029)
- Both are responsive (stack on narrow viewports) and use theme tokens, not hardcoded hex. (FR-011, FR-026)
