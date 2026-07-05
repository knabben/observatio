# Phase 0 Research: Screen Refactoring & UI Tech-Debt Remediation

All Technical Context unknowns and technology decisions are resolved below. Each decision is scoped to
the existing stack (Next.js 15 static export, React 19, Mantine 7.17, Tailwind 4, XYFlow 12,
react-use-websocket 4.13) so no new runtime dependency is introduced (Constitution: Technology Stack).

---

## R1 — Same-origin backend addressing (FR-036) — *revised per Clarification 2026-07-05*

**Decision**: Introduce `front/app/lib/config.ts` exporting `API_URL`, `WS_URL`, and `WS_URL_CHATBOT`
that **default to the page's own origin** (relative REST paths; WebSocket URL derived from
`window.location` — `ws:`/`wss:` matching `http:`/`https:`, same host/port). All modules import from
`config.ts` instead of the hardcoded `URL` const in `data.tsx`. A build-time `NEXT_PUBLIC_API_URL` /
`NEXT_PUBLIC_WS_URL` override is honored **only** to support the split development mode (frontend dev
server on :3000 → backend on :8080); it is unset in the embedded production build.

**Rationale**: The exported SPA is embedded via `//go:embed` and served by the **same Go binary** that
serves the API/WebSocket (see R11), so the browser is already on the correct origin. Deriving endpoints
from `window.location` makes one binary work on any host/port with **no per-origin rebuild** — the
previously-chosen build-time absolute-URL approach would bake a wrong `localhost:8080` into the shipped
binary. Origin derivation is a client runtime computation, so `output: "export"` (no server runtime) is
not a constraint here. No new dependency.

**Alternatives considered**:
- Build-time absolute `NEXT_PUBLIC_*` URLs (the prior R1 decision) — **rejected**: requires a rebuild per
  target origin and bakes a wrong default into the single binary; contradicts the single-binary goal.
- Runtime `/config.json` fetch at boot — adds a request + loading gate every session; unnecessary when the
  origin is already known from `window.location`.
- Server-side env via API routes — impossible under `output: "export"`.

---

## R2 — Bounded loading & WebSocket lifecycle (FR-003, FR-004, FR-005, FR-007)

**Decision**: Model each live view as an explicit state machine — `connecting → ready(data) | empty |
error` — driven by `react-use-websocket`'s `readyState` plus a **10s data-arrival timeout**. Configure the
socket with bounded reconnection: `reconnectAttempts: 8`, `reconnectInterval` using exponential backoff
(`Math.min(1000 * 2**attempt, 30000)`), and `onReconnectStop` → terminal `error` state. An empty/malformed
frame (no `.data`) is treated as a **no-op**, never as "clear the list".

**Rationale**: `react-use-websocket` 4.13 already exposes `readyState`, `shouldReconnect`,
`reconnectAttempts`, `reconnectInterval` (accepts a function of attempt count), and `onReconnectStop` — all
needed primitives exist. The 10s threshold matches standard web "something is wrong" expectations (well
above the <2s happy-path budget from Constitution II) and satisfies FR-003's "bounded resolution" without
guessing a backend-specific number. Distinguishing empty-frame from empty-result fixes the silent
list-wipe (FR-005).

**Alternatives considered**:
- Infinite reconnect (current `shouldReconnect: () => true`) — the audited bug; hammers a down server.
- No timeout, rely only on `readyState` — a socket that connects but never emits data still hangs (the
  audited infinite-spinner case).

---

## R3 — Responsive layout strategy (US2, FR-009–FR-015)

**Decision**: Replace fixed 12-grid spans and pixel canvases with Mantine responsive props and container
sizing:
- `Grid.Col`/`SimpleGrid` use responsive objects, e.g. `span={{ base: 12, md: 7 }}` / `cols={{ base: 1, sm: 2 }}`.
- Single-child panels drop the two-column wrapper (full-width) — fixes half-empty grids (FR-010).
- Tables wrap in `Table.ScrollContainer minWidth={…}` (FR-013).
- Topology container uses `width: '100%'` + a fixed/aspect height and `<ReactFlow fitView>` with
  `<Background/>` + `<Controls/>` (FR-012).
- Dashboard content area uses `overflow-y-auto` at all breakpoints (not `md:`-only) so tall content scrolls
  on small screens (FR-014).

**Rationale**: Mantine 7 grid/`SimpleGrid` accept breakpoint objects natively; `Table.ScrollContainer` is
the built-in overflow primitive; XYFlow's `fitView` recenters API-positioned nodes. All in-stack, no new deps.

**Alternatives considered**:
- Pure Tailwind grid rewrite — would fight Mantine's layout components already in use; inconsistent.
- CSS `container-queries` — Tailwind 4 supports them, but breakpoint props are simpler and already the
  project idiom.

---

## R4 — Consolidating the four duplicated resource areas (US4, FR-023–FR-025)

**Decision**: Introduce config-driven shared building blocks under `app/ui/dashboard/shared/`:
- `ObjectTable<T>` — takes a `ColumnDef<T>[]` (`{ header, render, align?, width? }`) + `getRowKey(item)`
  (stable id, never index) + `onSelect` + `emptyLabel`. Renders header, `Table.ScrollContainer`, rows, and
  an empty state.
- `StatusIndicator` — maps a tri-state (`healthy | notready | unknown`) to color + label + a11y text.
- `EmptyState` / `ErrorState` — labeled messages; `ErrorState` includes a retry action.
- `useResourceStream` — one hook encapsulating the WS state machine (R2), replacing the four copy-pasted
  fetch hooks and the per-file `handleFetchError`.
Per-resource files (`clusters`, `machines`, `mds`, and infra variants) collapse to a `ColumnDef` array +
detail field config that compose these blocks. `BaseLister`/`ObjectDetails` are refactored to consume them.

**Rationale**: The audit found the same bugs (index keys, missing null-guards, `text-bold`, `Created`
mislabel, missing empty states) duplicated 3–4× and already drifting. A single generic implementation makes
FR-025 (stable keys) and FR-002 (empty states) structural, so they cannot regress per-copy. Generics over
the existing CAPI-mapped types preserve Constitution III.

**Alternatives considered**:
- A data-grid library (e.g. mantine-react-table, TanStack Table) — new dependency requiring Constitution
  Technology-Stack justification; overkill for read-only tables of a handful of columns.
- Leave duplication, fix each copy — the audited status quo; guarantees future drift.

---

## R5 — Centralized theming, colors, and fonts (US4, FR-026, FR-027)

**Decision**: Create `app/styles/theme.ts` with a Mantine `createTheme` defining the brand accent
(consolidating the scattered `#aaf16a`, `#8feb83`, `#48654a`, `#a1f54d` greens into one `primaryColor` scale)
and semantic status colors (healthy/notready/unknown), applied via `MantineProvider theme={…}`. Fix
`fonts.ts` identifiers/weights and either define the `--font-geist-*` CSS vars from the actual fonts or remove
the dead `@theme` references in `globals.css`; drop the hardcoded `body { font-family: Arial }`. Reconcile the
forced `defaultColorScheme` with the `prefers-color-scheme` CSS so light/dark is internally consistent.

**Rationale**: One theme source removes the audit's "scattered greens + dead font tokens + scheme clash"
class of findings and makes contrast (FR-021) tunable in one place. Uses Mantine's native theming — no new dep.

**Alternatives considered**:
- Tailwind theme tokens only — Mantine components ignore Tailwind's theme; would leave two sources.

---

## R6 — Safe AI-panel content rendering & containment (US5, FR-032–FR-035)

**Decision**: Stop using `dangerouslySetInnerHTML` for bot/message content. Render message text as plain
text (with whitespace/line-break preservation) inside a valid element (not a `<div>` in a `<p>`); if
rich formatting is required later, sanitize first — but the default is safe plain-text rendering. Remove the
page-level `AppShell` from inside the grid column and replace with a self-contained flex/Card layout that
respects the parent region (no forced 100vh, no phantom header). Make expand reversible (toggle
`span` 6↔12). Guard `requestIA` on `readyState === OPEN`, reset `isLoading` on failure, and use functional
state updates. Replace `crypto.randomUUID()` with the already-present `uuid` v4 (secure-context-safe).

**Rationale**: Eliminates the XSS finding and invalid nesting, keeps the panel inside its card (Constitution
IV — AI grounded and contained), and fixes the one-way-expand and stuck-loading bugs. `uuid` is already a
dependency, so no new install.

**Alternatives considered**:
- Add a markdown renderer + sanitizer (e.g. `react-markdown` + `dompurify`) — new dependencies; deferred
  unless rich formatting is an explicit requirement (it is not in this spec).

---

## R7 — Status tri-state semantics (US3, FR-020)

**Decision**: Derive status via a pure helper `toStatusState(item)` returning `healthy | notready | unknown`:
readiness field present & true → healthy; present & false → notready; **absent/undefined → unknown**. Use
strict comparisons (no `== 0`), so unknown availability renders as unknown, not failed. `StatusIndicator`
renders healthy = solid green (no animation), notready = solid red (no `processing` pulse), unknown = neutral
gray. Not conveyed by color alone — each carries an accessible label (FR-017/FR-020).

**Rationale**: Fixes the audit's loose `== 0` (unknown shown as failed), the always-pulsing red indicator,
and the color-only signaling. Centralizing in one helper covers all screens uniformly.

**Alternatives considered**:
- Keep per-table inline logic — the audited drift source.

---

## R8 — Accessibility for interactive elements & navigation (US3, FR-016–FR-019, FR-022)

**Decision**: Clickable resource names become real `<button>`/Mantine `UnstyledButton` (or anchors with
`href`) — keyboard-focusable, activatable, `aria`-labeled. Icon-only nav links get `aria-label`/`title` and
`aria-current="page"` on the active item. Active-state detection uses `pathname.startsWith(link.href)` (with a
root exact-match guard) so nested routes highlight the parent (FR-016). Sidenav gains a responsive
collapse (Mantine `Burger` + `Drawer`/`AppShell` navbar toggle) so navigation is usable on narrow viewports
without a 250px logo dominating (FR-022); logo uses responsive `sizes`.

**Rationale**: Directly resolves the keyboard/AT findings and nested-route highlight bug using Mantine's
built-in `Burger`/`Drawer`/`AppShell` primitives (no new deps).

**Alternatives considered**:
- CSS-only `:hover` menu — not keyboard/AT accessible.

---

## R9 — "Search" control resolution (US5, FR-030)

**Decision**: Convert the mislabeled `Combobox`-over-button into a real text filter: a Mantine `TextInput`
(or `Autocomplete`) whose value filters the visible list by name/namespace client-side, wired through the
existing `FilterItems`/lister state. If a resource-scoped typeahead proves disproportionate during
implementation, the accepted fallback (per spec Assumptions) is relabeling the control to match its actual
select behavior — but the default target is a working filter.

**Rationale**: Matches the control to its label (FR-030) using in-stack Mantine inputs; leverages the
existing (null-hardened) `FilterItems` util.

**Alternatives considered**:
- Server-side search — requires backend changes, out of scope (backend contract unchanged).

---

## R10 — Single-binary embedding & build+verify make target (US6, FR-037–FR-039) — *Clarification 2026-07-05*

**Decision**: Keep the existing embed pipeline (`front` static export → copied to
`webserver/internal/web/handlers/build` → `//go:embed build/*` in `handlers.go`, served by
`system/spa.go`) and add a **`make verify-binary`** target that: (1) runs `make build`; (2) launches
`output/observatio` on an ephemeral port in the background; (3) polls the embedded UI root (`GET /` →
HTTP 200 with SPA HTML) and a live API/WebSocket endpoint on the **same origin**; (4) tears the process
down; (5) exits non-zero on any failure. The SPA fallback handler MUST return `index.html` for unknown
non-API routes so client-side routes resolve. `make build` remains the composition of
`build-frontend` + embed-copy + `build-backend`; `verify-binary` depends on `build`.

**Rationale**: The embedding already exists; the gap is an automated, clean-checkout-runnable proof that
the frontend refactor still embeds and self-serves (FR-038/FR-039). A background-launch + curl smoke test
uses only shell + the binary — no new dependency — and catches embed regressions (missing/stale assets,
broken SPA fallback, wrong-origin calls) that a plain compile would not. Same-origin addressing (R1) is
what lets the smoke test assert the UI and API share one origin.

**Alternatives considered**:
- Compile-only check (does `go build` succeed) — **rejected**: a binary can compile yet serve a blank/404
  UI if assets are missing or the SPA fallback is wrong.
- Go integration test using `httptest` against the embed FS — viable and complementary, but does not
  exercise the actual built binary end-to-end; the make target is the user-requested "seamless from make
  target" signal. May be added as a backend unit test in addition.

---

## R11 — Dependency currency: safe within-major updates (frontend + backend) — *user directive*

**Decision**: Update all existing dependencies to the **latest version within their current major**
(minor/patch only) for both `front` (pnpm) and `webserver` (Go modules). Explicitly **defer** the breaking
major jumps: frontend `@mantine/* 7→9`, `@jest/globals`/`@types/jest` `29→30`; backend
`anthropics/anthropic-sdk-go v0.2-beta→v1`, `k8s.io/* 0.32→0.36`, `sigs.k8s.io/cluster-api 1.9→1.13`,
`controller-runtime 0.19→0.24` (and the coupled `cluster-api-provider-vsphere`). Frontend via
`pnpm update` (respects `^` ranges) plus explicit within-major bumps for pinned `next`/`eslint-config-next`;
backend via targeted `go get` on the safe modules (`gin`, `gorilla/websocket`, `cobra`, `testify`) +
`go mod tidy`. Gate on `make run-tests-frontend`, `make run-tests-backend`, and `make build` staying green.

**Rationale**: Keeps the build/tests green and avoids a large, breaking migration that would collide with
this Mantine-7-based refactor and the CAPI/k8s-coupled backend (user selected "safe updates only"). The
k8s/CAPI/controller-runtime set is version-locked together and must move as one coordinated effort — out of
scope here.

**Alternatives considered**:
- Full latest incl. majors — **rejected by user**: red build + multi-day migration, conflicts with the
  refactor.
- No updates — leaves known minor/patch fixes unapplied; the user explicitly requested updating all.

---

## Resolved unknowns summary

| Unknown | Resolution |
|---------|------------|
| Backend endpoint addressing | **Same-origin** derived from `window.location`; `NEXT_PUBLIC_*` dev-only override (R1) |
| Bounded loading threshold value | 10s data-arrival timeout + backoff reconnect (R2) |
| Reconnect policy | 8 attempts, exponential backoff to 30s, terminal error (R2) |
| Responsive mechanism | Mantine breakpoint props + `Table.ScrollContainer` + XYFlow `fitView` (R3) |
| Consolidation pattern | Config-driven `ObjectTable`/shared blocks + `useResourceStream` (R4) |
| Theming source | Single `theme.ts` + fixed fonts/scheme (R5) |
| AI content safety | Plain-text render, no `AppShell`-in-grid, `uuid` for ids (R6) |
| Status semantics | Pure tri-state helper, strict comparisons (R7) |
| Accessibility & nav | Real buttons, `aria-current`, `startsWith` active state, responsive nav (R8) |
| Search behavior | Real client-side text filter, relabel fallback (R9) |
| Single-binary embed validation | `make verify-binary` smoke-tests the built binary end-to-end (R10) |
| Dependency currency | Safe within-major updates, frontend + backend; majors deferred (R11) |

**No NEEDS CLARIFICATION markers remain.**
