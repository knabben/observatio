# Phase 1 Data Model: Screen Refactoring & UI Tech-Debt Remediation

This is a frontend refactor with **no persistent storage** and **no backend model changes**. The "data
model" here is the set of **frontend view-models and configuration types** introduced to consolidate the
duplicated screens and make state handling explicit. All resource types continue to map to ClusterAPI CRDs
(Constitution III); this document only adds presentation-layer types over them.

---

## 1. Resource (existing, hardened)

Represents a ClusterAPI object (`Cluster`, `MachineDeployment`, `Machine`) or an infrastructure variant.

**Change**: The existing `types.tsx` declarations mark `metadata`, `status`, `conditions`, etc. as required
while the code treats them as optional. **Nullability is corrected** so types model reality: optional fields
become `field?: T`, forcing null-safe access at compile time (resolves the FR-001 crash class).

| Field | Type | Notes |
|-------|------|-------|
| `metadata?` | `Meta` | `name?`, `namespace?`, `ownerReferences?` all optional |
| `status?` | `Status` | `phase?`, `conditions?`, readiness flags optional |
| `spec`-level fields | resource-specific | `paused?`, `topology?`, `replicas: number` (was `string`) |

**Validation rules**: no field access may assume presence; every optional field is guarded or defaulted at
the presentation boundary.

---

## 2. StatusState (new)

Normalized tri-state health used by every status indicator (FR-020).

| Value | Meaning | Derivation |
|-------|---------|------------|
| `healthy` | ready flag present and true | `ready === true` |
| `notready` | ready flag present and false | `ready === false` |
| `unknown` | ready flag absent/undefined | `ready == null` |

- Derived by a pure function `toStatusState(item): StatusState` (strict comparisons only; no `== 0`).
- **State transitions**: `unknown → healthy | notready` as data populates; a resource never renders as
  `notready` solely because a field is missing.

---

## 3. ChannelState (new)

Explicit state machine for a live data view (FR-003–FR-007), replacing the implicit boolean `loading`.

| State | Meaning | Enter condition | Exit |
|-------|---------|-----------------|------|
| `connecting` | socket opening / awaiting first frame | mount | first data frame → `ready`; 10s timeout → `error`; open+empty result → `empty` |
| `ready` | data present | valid frame with items | new frames update in place |
| `empty` | connected, zero items | data frame with empty collection | new non-empty frame → `ready` |
| `error` | connection failed / reconnects exhausted | `onReconnectStop`, HTTP not-ok, or connect timeout | manual retry → `connecting` |

**Rules**:
- An empty/malformed frame (no `.data`) is a no-op and MUST NOT transition `ready → empty` (FR-005).
- `error` exposes a retry action that re-enters `connecting`.

---

## 4. ColumnDef<T> (new) — table configuration

Drives the shared `ObjectTable` so each resource area is data, not duplicated markup (FR-023).

| Field | Type | Notes |
|-------|------|-------|
| `header` | `string` | column label |
| `render` | `(item: T) => ReactNode` | cell content; null-safe |
| `align?` | `'left' \| 'center' \| 'right'` | single alignment mechanism (no Tailwind/Mantine mix) |
| `width?` | `number` | optional fixed width from theme scale |

Table config also carries `getRowKey: (item: T) => string` (stable unique id — never index, FR-025),
`onSelect: (item: T) => void`, and `emptyLabel: string`.

---

## 5. DetailFieldDef<T> (new) — detail/specification configuration

Drives the shared detail header and specification panels (FR-023, FR-029).

| Field | Type | Notes |
|-------|------|-------|
| `label` | `string` | accurate label (e.g. `Age`, not mislabeled `Created`) |
| `value` | `(item: T) => ReactNode` | null-safe; renders `—` for absent values |
| `span?` | responsive object | for responsive detail grids |

---

## 6. AppConfig (new) — same-origin endpoint configuration

Single source for backend endpoints (FR-036). **Default is same-origin** (derived from `window.location`);
`NEXT_PUBLIC_*` is a dev-only override. See `contracts/environment-config.md`.

| Field | Type | Source |
|-------|------|--------|
| `API_URL` | `string` | `NEXT_PUBLIC_API_URL` (dev) \| `''` ⇒ relative to origin (prod) |
| `WS_URL` | `string` | `NEXT_PUBLIC_WS_URL` (dev) \| `ws(s)://<origin>/ws` (prod) |
| `WS_URL_CHATBOT` | `string` | `NEXT_PUBLIC_WS_URL_CHATBOT` (dev) \| `ws(s)://<origin>/chatbot` (prod) |

## 6b. DeployableBinary (new) — packaging artifact

The single self-contained binary (`output/observatio`) embedding the exported frontend (US6, FR-037–039).

| Aspect | Value |
|--------|-------|
| Embed source | `front/output` → `webserver/internal/web/handlers/build` (`//go:embed build/*`) |
| Serving | UI + API/WebSocket on one origin; SPA fallback → `index.html` |
| Validation | `make verify-binary` (build → launch → smoke-test → exit code) |

---

## 7. ThemeTokens (new) — centralized styling

Single Mantine theme + semantic tokens (FR-026, FR-021).

| Token | Purpose |
|-------|---------|
| `primaryColor` (accent scale) | consolidates scattered greens |
| `status.healthy / notready / unknown` | status indicator colors (contrast-checked) |
| font family vars | actual configured fonts (dead `--font-geist-*` removed) |

---

## Entity relationships

```text
Screen ──has many──> Resource ──derives──> StatusState
  │                     │
  │                     └──rendered via──> ColumnDef[] / DetailFieldDef[]
  │
  ├──feeds from──> ChannelState (WebSocket/REST)
  └──styled by──> ThemeTokens ; configured by AppConfig
```

No database, migrations, or persisted state are involved.
