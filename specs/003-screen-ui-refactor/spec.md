# Feature Specification: Screen Refactoring & UI Tech-Debt Remediation

**Feature Branch**: `003-screen-ui-refactor`
**Created**: 2026-07-05
**Status**: Draft
**Input**: User description: "Review ALL the screens one by one and suggest a refactoring, fixing layouts and broken issues that can appear, take as a tech debt enumeration of open issues, review careful each screen and change what can be wrong"

## Overview

The Observātiō dashboard presents cluster health across five screens: the **Dashboard overview**, **Clusters**, **Machine Deployments**, **Machines**, and the shared **navigation/AI-troubleshooting** shell. A screen-by-screen audit found recurring defects that undermine the product's core promise of reliable, real-time visibility: screens can crash on partial data, hang forever on a silent WebSocket, render half-empty or overflowing layouts, and offer no feedback when data is empty or a connection fails. This feature is a remediation pass that makes every screen render correctly across data and viewport conditions, behave predictably under failure, and consolidate the heavily duplicated screen components so the same bug cannot reappear in four places.

This is a **non-functional refactor**: no new product capabilities are introduced. Existing behavior is preserved except where the current behavior is itself the defect (e.g., an infinite spinner, a crash, a mislabeled field, a control that does not do what its label claims). The one build/packaging concern in scope is validating that the refactored frontend still embeds into and ships as the existing single Go binary.

## Clarifications

### Session 2026-07-05

- Q: When the SPA is embedded and served by the single Go binary (same-origin), how should the frontend address the backend API/WebSocket? → A: Same-origin/relative — the browser derives the API and WebSocket URLs from the page origin it was served from; `NEXT_PUBLIC_*` overrides remain **only** for the split development mode (frontend dev server on :3000 → backend on :8080).
- Q: How should the "embedding + build is seamless" validation be realized? → A: A dedicated automated make target (e.g. `make verify-binary`) that runs the full build, launches `output/observatio`, and asserts the embedded UI root returns HTTP 200 and a live API/WebSocket endpoint responds on the same origin, with no external assets required — failing on any embed/serve regression.
- Q: Must the previously resolved items (same-origin addressing, single-binary build+verify target) and the approved dependency-update policy be reflected in the plan artifacts? → A: Yes — the plan is authoritative and MUST incorporate every previously resolved item; `plan.md`, `research.md`, and `contracts/` are updated so no design artifact contradicts this spec, and the safe within-major dependency updates (frontend + backend) are recorded as planned work.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Screens never crash or hang on real cluster data (Priority: P1)

An operator opens each screen against a live cluster whose resources may be mid-provisioning, partially populated, or missing optional fields. Every screen renders a meaningful result — data, an empty-state message, or an error message — and never shows a blank white area, an unhandled crash, or a spinner that never resolves.

**Why this priority**: Crashes and infinite spinners directly break the product's reason to exist (visibility into cluster health). An operator troubleshooting an outage cannot tolerate a dashboard that itself fails when a resource is in a degraded or transitional state — which is exactly when they need it most.

**Independent Test**: Point each screen at cluster resources with missing/partial `status`, `metadata.name`, `conditions`, and `paused` fields, and at a backend that connects but sends no data; confirm each screen shows data, a "no items" message, or an error message rather than crashing or spinning indefinitely.

**Acceptance Scenarios**:

1. **Given** a cluster resource missing its `status` block, **When** the operator opens the Clusters, Machines, or Machine Deployments list, **Then** the row and detail panel render without throwing and show a neutral/unknown status indicator instead of a blank screen.
2. **Given** a WebSocket that connects but never delivers a message within a reasonable time, **When** any list screen loads, **Then** the screen resolves to an empty-state or timeout message instead of showing the loading spinner forever.
3. **Given** an empty result set (no clusters, machines, deployments, versions, or cluster classes), **When** the corresponding screen loads, **Then** a clear "No … found" message is shown instead of a header-only table or an empty chart.
4. **Given** the backend returns an HTTP error or the WebSocket connection drops, **When** any screen requests data, **Then** the operator sees an actionable error message and a way to retry, not a silent failure or a perpetual loader.
5. **Given** a resource whose numeric field is `0` (e.g. CPU cores per socket), **When** the machine specification renders, **Then** the value is displayed correctly and no stray `0` or malformed markup leaks into the table.

---

### User Story 2 - Layouts render correctly across screen sizes and data volumes (Priority: P1)

An operator uses the dashboard on a laptop, an external monitor, and occasionally a tablet. Every screen fills its available space without overflow, horizontal scrollbars, half-empty columns, or content clipped off-screen, and adapts as the viewport narrows.

**Why this priority**: Broken and lopsided layouts are the most visible quality signal and directly impede reading cluster state. Fixed-pixel canvases and single-child two-column grids currently force horizontal scrolling and waste half the screen, degrading usability on every viewport.

**Independent Test**: Load each screen at desktop (≥1440px), laptop (~1280px), and tablet (~768px) widths with both small and large data sets; verify no horizontal page scroll, no clipped content, no permanently empty half-panels, and that multi-column sections stack gracefully as width decreases.

**Acceptance Scenarios**:

1. **Given** the Dashboard overview at laptop width, **When** the cluster topology renders, **Then** the topology fits within its panel and recenters to show all nodes rather than forcing horizontal scroll or rendering nodes off-screen.
2. **Given** any specification/detail panel, **When** it contains a single content block, **Then** it uses the full panel width instead of rendering at half width with a permanently blank adjacent column.
3. **Given** a table with long names, namespaces, or server/thumbprint values, **When** the viewport is narrow, **Then** the table scrolls within its own container instead of overflowing the page.
4. **Given** any screen viewed at tablet width, **When** it contains a two-column layout, **Then** the columns stack vertically rather than colliding or cramping.
5. **Given** the primary navigation on a narrow viewport, **When** the operator needs to switch screens, **Then** the navigation remains usable (labels or accessible names are available) and does not consume a disproportionate share of the viewport.

---

### User Story 3 - Consistent, accessible navigation and status feedback (Priority: P2)

An operator relies on the sidebar to know which screen they are on and on status indicators to read health at a glance. Navigation highlights the active section (including nested routes), every interactive element is keyboard-reachable and screen-reader labeled, and status indicators visually distinguish "healthy", "not ready", and "unknown".

**Why this priority**: Navigation and status legibility are used on every screen and every session. Incorrect active-state, unlabeled icon-only controls, and status indicators that all pulse identically make the dashboard harder to trust and inaccessible to keyboard/AT users, but they do not block core rendering, so they rank below crashes and layout breakage.

**Independent Test**: Navigate to a nested route and confirm the parent nav item is highlighted; traverse every clickable name/link/icon-button by keyboard alone; inspect status indicators for healthy vs. not-ready vs. unknown resources.

**Acceptance Scenarios**:

1. **Given** the operator is on a nested route (e.g. a specific cluster's detail view), **When** the sidebar renders, **Then** the corresponding top-level navigation item is shown as active.
2. **Given** a clickable resource name in any table, **When** the operator uses keyboard navigation, **Then** the element is focusable and activatable via keyboard, and screen readers announce it as an interactive control.
3. **Given** a resource in a "not ready" or error state, **When** its status indicator renders, **Then** the indicator is visually distinct from a "healthy" indicator and does not use an animation that implies work-in-progress for a failed state.
4. **Given** an icon-only navigation item on a narrow viewport, **When** a screen reader encounters it, **Then** it announces a meaningful label rather than an empty or generic name.

---

### User Story 4 - Consolidated screen components and consistent theming (Priority: P2)

A developer maintaining the dashboard changes shared behavior (a table column, a status rule, a detail header) in one place and has it apply consistently to Clusters, Machines, Machine Deployments, and their infrastructure variants — instead of editing four near-identical copies that have already drifted apart.

**Why this priority**: The lister/table/details/specification components are copy-pasted across four resource areas, and defects (index keys, dead CSS classes, mislabeled fields, missing null guards) already exist in some copies but not others. Consolidation prevents the same bug from silently reappearing and is what makes the P1 fixes durable — but it is internal quality, invisible to operators, so it ranks below user-facing defects.

**Independent Test**: Confirm the list, table, detail-header, and specification presentations for each resource area derive from a shared, configurable implementation, and that a single change to shared presentation logic is reflected across all resource areas.

**Acceptance Scenarios**:

1. **Given** the shared list/table/detail presentation, **When** a developer changes a shared rule (e.g. how the row key or status color is derived), **Then** the change applies uniformly to all resource areas without per-area edits.
2. **Given** the theme's accent colors and typography, **When** any screen renders, **Then** colors and fonts resolve from a single defined source rather than scattered hardcoded values and dead/undefined style tokens.
3. **Given** the fonts and color-scheme configuration, **When** the app loads, **Then** the intended fonts are actually applied and the light/dark presentation is internally consistent (no clash between forced scheme and system-preference styling).

---

### User Story 5 - Controls behave as labeled (Priority: P3)

An operator interacting with the search control and the AI-troubleshooting panel gets behavior consistent with each control's label and appearance: "Search" filters results, informational status chips are not mistaken for toggles, and the AI chat panel is safe and stays within its region.

**Why this priority**: These are narrower, screen-specific correctness issues affecting fewer flows than the systemic P1/P2 items, but they represent controls that actively mislead the operator or introduce risk (unsanitized content rendering, a "Search" that cannot be typed into, status chips that appear clickable).

**Independent Test**: Exercise the search control to confirm it filters visible results (or is corrected to match its actual behavior); confirm condition/status chips cannot be toggled by the operator; confirm AI-panel content renders safely and the panel stays within its container and can be collapsed after expanding.

**Acceptance Scenarios**:

1. **Given** the control labeled "Search", **When** the operator interacts with it, **Then** it either filters the visible list as its label implies or is relabeled/restyled to match what it actually does.
2. **Given** read-only condition/status indicators, **When** the operator clicks one, **Then** it does not appear to change state as if it were an editable toggle.
3. **Given** the AI-troubleshooting panel receives message content, **When** it renders that content, **Then** the content is displayed safely without executing embedded markup, and the panel stays within its card region rather than forcing full-viewport height.
4. **Given** the operator expands the AI chat panel, **When** they want to return to the previous size, **Then** a collapse action is available.

---

### User Story 6 - Runs as a single self-contained binary (Priority: P2)

An operator or CI builds the project with one make target and receives a single binary that, when run, serves the full dashboard UI and its API/WebSocket on one origin — no separate frontend server, no external asset directory, and no per-environment rebuild.

**Why this priority**: The single-binary embed is how the product is actually shipped and run; a broken embed pipeline means the refactored UI never reaches users. It ranks P2 because it protects existing delivery rather than changing an operator-facing screen — but it gates release.

**Independent Test**: On a clean checkout, run the build+verify make target; execute the produced binary; confirm the UI root returns HTTP 200, static assets load from the binary itself, and the API/WebSocket respond on the same origin.

**Acceptance Scenarios**:

1. **Given** a clean checkout, **When** the build+verify make target runs, **Then** it produces one self-contained binary and the automated smoke check reports the embedded UI root returns HTTP 200.
2. **Given** the running binary, **When** the browser loads the UI, **Then** the SPA calls the API and WebSocket on the same origin it was served from (no `localhost:8080` or external host baked in).
3. **Given** a broken or missing embed (assets absent/stale), **When** the verify make target runs, **Then** it fails with a clear error instead of producing a binary that serves a blank or 404 UI.
4. **Given** the binary is relocated to another host or port, **When** it is run there, **Then** the UI and API/WebSocket continue to work without any rebuild (same-origin addressing).

---

### Edge Cases

- **Partial resources**: `metadata`, `metadata.name`, `status`, `status.phase`, `conditions`, `paused`, `ownerReferences`, `configRef`, and infrastructure fields may each be absent; every screen MUST tolerate any subset being missing.
- **Empty collections**: zero clusters / machines / deployments / versions / cluster classes / conditions / cluster classes MUST each produce a labeled empty state, not a header-only table or empty chart.
- **Zero-valued numbers**: numeric fields equal to `0` (replicas, cores, unavailable replicas) MUST render as data, and MUST NOT be treated as "absent" or leak stray characters into markup.
- **Unknown availability**: a resource whose availability/readiness is unknown (field absent) MUST be shown as "unknown", distinct from both "healthy" and "failed".
- **Connection failure**: WebSocket disconnect, connect-but-no-data, and REST HTTP error responses MUST each surface a user-visible state and MUST NOT loop reconnection attempts unbounded or hang on a loader.
- **Large data volumes**: long field values and long lists MUST scroll within their container without breaking page layout.
- **Duplicate identifiers**: list items sharing a display name (e.g. version components across kinds) MUST still render with stable, unique identity.
- **Non-secure context**: features that depend on secure-context browser APIs MUST degrade gracefully when unavailable.
- **Missing/stale embed**: if the frontend assets are absent or not refreshed before the binary is compiled, the build+verify make target MUST fail loudly rather than produce a binary that serves a blank or 404 UI.
- **Relocated binary**: the binary run on a different host/port MUST continue serving UI + API on that origin with no rebuild (same-origin addressing), and MUST NOT reference a baked-in `localhost` backend.

## Requirements *(mandatory)*

### Functional Requirements — State & Robustness (US1)

- **FR-001**: Every screen MUST render without throwing when any optional resource field (`metadata`, `metadata.name`, `status`, `status.phase`, `conditions`, `paused`, `ownerReferences`, `configRef`, infrastructure fields) is absent.
- **FR-002**: Every list and table MUST display a labeled empty-state message when its collection is empty.
- **FR-003**: Every data-loading region MUST resolve out of its loading state within a bounded time, presenting an empty-state or error message if no data arrives, so no screen can spin indefinitely.
- **FR-004**: Every data source (WebSocket and REST) MUST surface connection and HTTP error conditions to the operator as an actionable message, and MUST NOT treat error responses as valid data.
- **FR-005**: An empty, malformed, or keepalive data frame MUST NOT silently clear an already-populated list.
- **FR-006**: Numeric fields equal to `0` MUST render as data and MUST NOT be misinterpreted as absent; conditional rendering MUST NOT emit stray values into markup.
- **FR-007**: Reconnection behavior MUST be bounded (limited attempts and/or backoff) rather than an unbounded tight loop, and MUST report a terminal failure state.
- **FR-008**: Sorting and filtering of resource lists MUST tolerate items with missing sort keys without failing the entire list.

### Functional Requirements — Layout & Responsiveness (US2)

- **FR-009**: No screen MUST introduce a horizontal page scrollbar at supported viewport widths (see Assumptions for the supported range).
- **FR-010**: Panels containing a single content block MUST occupy the full available width and MUST NOT reserve a permanently empty adjacent column.
- **FR-011**: Multi-column layouts MUST stack vertically as viewport width decreases so content does not collide or cramp.
- **FR-012**: The cluster topology visualization MUST fit within its container, adapt to available width, and present all nodes within view (with a means to recenter/zoom).
- **FR-013**: Tables with potentially wide content MUST scroll within their own container rather than overflowing the page.
- **FR-014**: The dashboard content area MUST remain scrollable at all supported viewport sizes so tall content is never clipped without a scroll affordance.
- **FR-015**: Loading indicators MUST be centered within the region they occupy.

### Functional Requirements — Navigation, Status & Accessibility (US3)

- **FR-016**: Primary navigation MUST indicate the active section for both exact and nested routes.
- **FR-017**: The active navigation item MUST be programmatically identifiable to assistive technology, not conveyed by color alone.
- **FR-018**: All interactive elements (clickable names, links, icon buttons, back/expand controls) MUST be keyboard-focusable, keyboard-activatable, and have accessible names.
- **FR-019**: Icon-only controls MUST expose a descriptive accessible label.
- **FR-020**: Status indicators MUST visually and semantically distinguish "healthy", "not ready/failed", and "unknown", and MUST NOT apply an in-progress animation to a failed/static state.
- **FR-021**: Text and status colors MUST meet accessibility contrast expectations against their background.
- **FR-022**: The primary navigation MUST remain usable on narrow viewports without consuming a disproportionate share of the screen.

### Functional Requirements — Consolidation & Consistency (US4)

- **FR-023**: The list, table, detail-header, and specification presentations for Clusters, Machines, Machine Deployments, and their infrastructure variants MUST derive from shared, configurable components such that shared presentation logic exists in one place.
- **FR-024**: A change to shared presentation logic (row identity, status derivation, header layout) MUST apply consistently across all resource areas.
- **FR-025**: List item identity MUST be derived from a stable unique resource identifier, not array position.
- **FR-026**: Accent colors, status colors, and typography MUST resolve from a single defined source; hardcoded and dead/undefined style tokens MUST be removed.
- **FR-027**: The intended fonts MUST be applied application-wide, and light/dark presentation MUST be internally consistent with no conflict between forced and system-preference styling.
- **FR-028**: Component documentation/comments MUST accurately describe each component's actual responsibility (no stale claims of behavior that lives elsewhere).
- **FR-029**: Field labels MUST accurately describe the value shown (e.g. a duration MUST NOT be labeled as a creation timestamp).

### Functional Requirements — Control Correctness (US5)

- **FR-030**: The "Search" control MUST either filter the visible results as its label implies, or be corrected in label/appearance to match its actual behavior.
- **FR-031**: Read-only condition/status indicators MUST NOT present as editable toggles.
- **FR-032**: Message/dynamic content rendered in the AI-troubleshooting panel MUST be displayed safely, without executing embedded markup, and MUST use valid document structure.
- **FR-033**: The AI-troubleshooting panel MUST remain within its container region and MUST NOT force full-viewport height or reserve empty header space inside a card.
- **FR-034**: Any panel that can be expanded MUST also be collapsible back to its prior size.
- **FR-035**: Message-send flows MUST verify the connection is open before sending and MUST reset transient in-progress state when a send cannot complete.

### Functional Requirements — Configuration

- **FR-036**: The frontend MUST address the backend API and WebSocket on the **same origin** it was served from (relative/derived URLs), so the embedded single-binary deployment works on any host/port with no per-origin rebuild. Absolute endpoints MAY be overridden via build-time `NEXT_PUBLIC_*` variables **only** for the split development mode (frontend dev server → separate backend); no hardcoded `localhost` address may remain in the embedded production path.

### Functional Requirements — Build & Packaging (Single Binary)

- **FR-037**: The entire stack (static frontend + API/WebSocket server) MUST be runnable from a single self-contained binary with the frontend assets embedded; the running binary MUST serve the UI and the API/WebSocket from the same origin with no dependency on external asset files.
- **FR-038**: A single make target MUST build the frontend, embed its output into the binary, and automatically validate that the resulting binary serves the embedded UI root (HTTP 200) and a live API/WebSocket endpoint on the same origin — failing the target if embedding or serving is broken.
- **FR-039**: The build+verify make target MUST be runnable from a clean checkout and MUST NOT depend on a separately running frontend dev server or externally hosted assets.

### Key Entities

- **Screen**: A top-level operator view — Dashboard overview, Clusters, Machine Deployments, Machines — each composed of a list/table, a detail/specification panel, and shared navigation.
- **Resource**: A ClusterAPI object (Cluster, MachineDeployment, Machine) and its infrastructure variant, with a `metadata`, `status`/`conditions`, and specification fields, any of which may be partially populated.
- **Status Indicator**: A visual element conveying resource health in one of three states — healthy, not-ready/failed, unknown.
- **Data Channel**: The live (WebSocket) and request/response (REST) sources feeding screens, each with connected, empty, and error conditions; addressed on the same origin as the served UI.
- **Shared Presentation Component**: The consolidated list/table/detail/specification building blocks reused across all resource areas.
- **Deployable Binary**: The single self-contained artifact produced by the build, embedding the exported frontend and serving UI + API/WebSocket on one origin; validated by the build+verify make target.

## Tech-Debt Inventory (per screen) *(reference for planning)*

<!--
  Concrete findings from the screen-by-screen audit, retained to scope planning
  and verify coverage. Each maps to one or more requirements above.
-->

### Dashboard overview

- Fixed 860×500 topology canvas overflows narrow viewports; bare flow with no fit-to-view, empty-state, background, or controls (FR-012, FR-002).
- Suppressed error masking possibly-undefined topology data; stale node state on refetch (FR-001).
- Non-responsive 5/7 and 6/6 column splits with fixed chart height (FR-011).
- `conditions.map` / `status.toLowerCase()` unguarded → crash on partial cluster-class data (FR-001).
- Read-only condition chips rendered as interactive toggles; BigInt rendered as raw child (FR-031, FR-001).
- Missing empty states for summary chart, versions table, and cluster-class table (FR-002).
- Duplicated fetch-hook boilerplate with copy-pasted wrong error strings ("cluster summary" in hierarchy/versions/class) (FR-028).
- Title links to its own page; low-contrast header color (FR-021).

### Clusters (and infra variant)

- `cluster.metadata.name` / `status.phase` / `paused.toString()` unguarded → crashes; inconsistent null-safety vs infra copy (FR-001).
- Single-child two-column specification grid renders at half width with blank column (FR-010).
- Array-index row keys; no empty states; no table scroll containers (FR-025, FR-002, FR-013).
- Non-accessible `<a>` without href; always-pulsing red indicator; magic sizes/colors (FR-018, FR-020, FR-026).
- Dead Tailwind classes (`text-bold`, `text-medium`); types declare fields non-null that code treats as optional (FR-026, FR-001).
- Infra failure fields (`failureReason`/`failureMessage`) and `modules` typed but never surfaced (FR-004).

### Machines & Machine Deployments (and machine infra variant)

- `0 && <JSX>` leaks a literal `0` into the machine spec table (FR-006).
- Mixed guarded/unguarded access on the same `status` object → partial crash (FR-001).
- Loose `== 0` treats unknown availability as failed (FR-020).
- Array-index keys vs name-keys drift between copies; missing `'use client'` in one table; two different centering mechanisms (FR-025, FR-024).
- No empty states, no scroll containers, `Created` label shows an age/duration (FR-002, FR-013, FR-029).
- Unused `v1beta2` status sub-type and dead failure fields (FR-004).

### Shared shell (navigation, base, AI, data)

- Infinite spinner when socket connects but sends nothing; empty frame clears the list; unbounded reconnect loop (FR-003, FR-005, FR-007).
- Sidebar has no mobile collapse; 250px logo dominates small viewports; strict-equality active state misses nested routes (FR-022, FR-016).
- Icon-only nav links lack accessible names; no `aria-current` (FR-018, FR-017).
- "Search" is a non-typeable select mislabeled as search (FR-030).
- AI panel: unsanitized content rendering, invalid nesting, `AppShell`-in-grid forcing 100vh, one-way expand, unchecked send, secure-context-only UUID (FR-032, FR-033, FR-034, FR-035).
- Loader not vertically centered; hardcoded `localhost` backend URL; no `res.ok` checks (FR-015, FR-036, FR-004).
- Dead/undefined font tokens; forced dark scheme clashing with system-preference CSS; scattered hardcoded greens (FR-026, FR-027).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% of screens render a data, empty-state, or error result — and zero unhandled crashes — when exercised against resources with any single optional field missing.
- **SC-002**: Zero screens remain in a loading state longer than a bounded threshold when the backend connects but sends no data; each resolves to an empty-state or error message.
- **SC-003**: Zero horizontal page scrollbars and zero permanently-empty half-panels across all screens at desktop, laptop, and tablet widths.
- **SC-004**: Every empty collection across all screens shows a labeled empty-state message (0 header-only tables, 0 empty charts without a message).
- **SC-005**: Every interactive element on every screen is reachable and activatable by keyboard alone and exposes an accessible name (100% keyboard traversal).
- **SC-006**: Active navigation is correctly indicated for 100% of exact and nested routes.
- **SC-007**: Status indicators correctly distinguish healthy, not-ready/failed, and unknown states for 100% of sampled resources, including zero-valued and absent fields.
- **SC-008**: The number of independent copies of the list/table/detail/specification presentation is reduced to a single shared implementation per concern (from four resource-area copies to one configurable source).
- **SC-009**: The embedded single-binary build serves the UI and reaches the API/WebSocket on the same origin with zero hardcoded backend addresses; relocating the binary to another host/port requires no rebuild.
- **SC-010**: The full frontend test suite passes and covers the previously-crashing partial-data and empty-collection cases for each screen.
- **SC-011**: A single make target builds the frontend, embeds it into the binary, and automatically verifies the produced binary serves the embedded UI root (HTTP 200) and a same-origin API/WebSocket endpoint — failing on any embed/serve regression — runnable from a clean checkout with no separate frontend server.

## Assumptions

- **Scope is remediation, not redesign**: existing screens, routes, and data model stay; only defective layout, state handling, accessibility, consistency, and mislabeled controls are corrected. No new product capabilities are added beyond making labeled controls behave correctly.
- **Supported viewport range**: primary targets are desktop and laptop widths, with usable (stacked, scroll-contained) layouts down to tablet width (~768px). Full phone-width optimization and a full mobile navigation drawer are desirable but treated as best-effort within US2/US3, not a launch blocker.
- **"Search" resolution**: the currently non-functional "Search" control is assumed to be corrected to filter the visible list; if a full typeahead is disproportionate, relabeling/restyling to match actual behavior is an acceptable alternative (FR-030).
- **Bounded loading threshold**: a specific timeout value for FR-003 will be chosen during planning using standard web expectations; the requirement is that some bounded resolution exists, not a particular number.
- **Backend contract unchanged**: the backend API and WebSocket message shapes are stable; this work adapts the frontend to tolerate partial/empty payloads the backend may already send. The only intentional touch of the Go/build layer is (a) serving the SPA same-origin and (b) the build+verify make target — no API/message-shape changes.
- **Single-binary delivery (existing, validated)**: the project already embeds the exported frontend into the Go binary via `//go:embed` (`webserver/internal/web/handlers/build`) and the `make build` copy step; this feature makes the frontend address the backend same-origin (FR-036) and adds an automated make target (FR-037–FR-039) that verifies the embed/build end-to-end. Same-origin addressing (Clarification 2026-07-05) supersedes the earlier build-time absolute-URL approach; `NEXT_PUBLIC_*` is retained solely for split dev mode.
- **Accessibility target**: WCAG AA-level contrast and keyboard operability are the reference standard for FR-018–FR-021.
- **Testing**: per the project constitution, changes are accompanied by frontend tests covering the partial-data, empty-collection, and error-state cases enumerated here.
