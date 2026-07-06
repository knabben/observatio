---

description: "Task list for feature implementation"
---

# Tasks: Object YAML View & Global AI Troubleshooting Panel

**Input**: Design documents from `/specs/005-object-viz-ai-panel/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/raw-object-api.md, quickstart.md

**Tests**: Included per Constitution Principle V ("Test-Driven Quality" — all backend/frontend
changes MUST be accompanied by tests); not full upfront TDD, but every story ships tests alongside
its implementation.

**Organization**: Tasks are grouped by user story (from spec.md) to enable independent
implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1, US2, US3)

## Path Conventions

Existing web app: `webserver/` (Go backend), `front/` (Next.js frontend) — see plan.md Project
Structure for exact new/changed files. The seven object detail components touched throughout this
feature: `clusters/details.tsx`, `clusters/infra/infra-details.tsx`, `clusters/infra/docker-details.tsx`,
`machines/details.tsx`, `machines/infra/infra-details.tsx`, `machines/infra/docker-details.tsx`,
`mds/details.tsx` (all under `front/app/ui/dashboard/components/`).

---

## Phase 1: Setup

- [X] T001 Confirm `make build`, `make run-tests-backend`, and `make run-tests-frontend` all pass
  on a clean checkout of `005-object-viz-ai-panel` before starting implementation.

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Stand up the app-wide AI panel shell (empty-context version) that both US1 and US2
extend, and that removes the old per-object embed's home.

**⚠️ CRITICAL**: No user story work can begin until this phase is complete.

- [X] T002 Create `front/app/ui/dashboard/ai-panel/ai-panel-context.tsx`: `AIPanelProvider` holding
  `isOpen`, `open()`/`close()`, `messages`, `currentObjectContext` (initially always `null`),
  `queryField`, `queryFieldTouched`, and a `useAIPanel()` hook — mirrors the
  `InfraCapabilityContext` pattern from feature 004 (data-model.md `AIPanelState`).
- [X] T003 [P] Create `front/app/ui/dashboard/ai-panel/ai-panel.tsx`: the `Drawer`-based panel UI
  (reusing the `Drawer` pattern already used by `front/app/ui/dashboard/sidenav.tsx`), built from
  the existing `ChatBot` markup in `front/app/ui/dashboard/base/ai-troubleshooting.tsx` but with
  every hardcoded color (`#0f0f23`, `#1a1a3e`, `#00d4aa`, `#4a4a6a`, the `rgba(...)` message
  bubbles) replaced with the dashboard's existing theme tokens (`var(--mantine-color-brand-*)`)
  per FR-010.
- [X] T004 [P] Create `front/app/ui/dashboard/ai-panel/ai-panel-trigger.tsx`: the persistent global
  open control (FR-001).
- [X] T005 Mount `<AIPanelProvider>` around `{children}` and render `<AIPanelTrigger>` in
  `front/app/dashboard/layout.tsx` (depends on T002-T004).
- [X] T006 [P] Add `front/app/ui/dashboard/ai-panel/ai-panel-context.test.tsx`: open/close toggling,
  and a message sent before closing is still present after reopening (FR-003).

**Checkpoint**: The global, collapsible, theme-consistent AI panel exists and is reachable — ready
for US1/US2 to build on.

---

## Phase 3: User Story 1 - AI troubleshooting is available from anywhere (Priority: P1) 🎯 MVP

**Goal**: The panel from Phase 2 is reachable from every screen, the embedded "AI Troubleshooting"
tab is gone from every object type, and the Object Conditions table lives on in "Specification".

**Independent Test**: From the Clusters list, the Dashboard overview, and a Machine detail screen,
confirm the same panel opens each time and can be collapsed/reopened without losing its state;
confirm no object detail screen has an "AI Troubleshooting" tab, and each still shows its conditions
table under "Specification".

### Tests for User Story 1

- [X] T007 [P] [US1] Jest test: `AIPanelTrigger`/`AIPanelProvider` open the same panel instance
  regardless of which screen renders them, in `front/app/ui/dashboard/ai-panel/ai-panel-trigger.test.tsx`.
- [X] T008 [P] [US1] Jest test: extend each of the seven detail components' existing test files (or
  add one) asserting the tabs list no longer contains "AI Troubleshooting" and does contain
  "Specification" with the conditions table still rendered.

### Implementation for User Story 1

- [X] T009 [US1] Delete `front/app/ui/dashboard/base/ai-troubleshooting.tsx` now that its markup has
  been migrated into `ai-panel.tsx` (T003) (depends on T003).
- [X] T010 [P] [US1] Update `front/app/ui/dashboard/components/clusters/details.tsx`
  (`ClusterDetails`): remove the "AI Troubleshooting" tab; ensure the conditions table remains
  visible within "Specification".
- [X] T011 [P] [US1] Same edit for `front/app/ui/dashboard/components/clusters/infra/infra-details.tsx`
  (`ClusterInfraDetails`, vSphere).
- [X] T012 [P] [US1] Same edit for `front/app/ui/dashboard/components/clusters/infra/docker-details.tsx`
  (`ClusterInfraDockerDetails`). **Discovered mid-implementation**: this component (and T015's
  Docker machine equivalent) had *only* an "AI Troubleshooting" tab — no "Specification" tab
  existed at all. Removing the tab without adding one would have left zero tabs. Added a new
  minimal `docker-specification.tsx` (Cluster/LoadBalancerIP/Ready) alongside it, matching the
  vSphere Specification pattern, so the screen isn't regressed to empty.
- [X] T013 [P] [US1] Same edit for `front/app/ui/dashboard/components/machines/details.tsx`
  (`MachineDetails`).
- [X] T014 [P] [US1] Same edit for `front/app/ui/dashboard/components/machines/infra/infra-details.tsx`
  (`MachineInfraDetails`, vSphere).
- [X] T015 [P] [US1] Same edit for `front/app/ui/dashboard/components/machines/infra/docker-details.tsx`
  (`MachineInfraDockerDetails`) — same missing-Specification-tab gap as T012, same fix
  (`docker-specification.tsx` for machines).
- [X] T016 [P] [US1] Same edit for `front/app/ui/dashboard/components/mds/details.tsx`
  (`MachineDeploymentDetails`).

**Also discovered**: none of the 7 `Specification` components rendered a conditions table before
this feature (conditions were ONLY ever shown inside the now-removed AI Troubleshooting tab).
Extracted the Chip-based table markup from the old `ai-troubleshooting.tsx` into a new shared
`front/app/ui/dashboard/shared/conditions-table.tsx`, and added it to each of the 5
condition-bearing Specification components (the two Docker variants have no conditions data to
show). The old `ai-troubleshooting.test.tsx` also bundled unrelated `Search` and `ClusterClassLister`
tests (pre-existing organizational debt) — relocated to `search.test.tsx` and
`dashboard.test.tsx` respectively rather than deleted, and its XSS-safety test was migrated to
`ai-panel.test.tsx`.

**Checkpoint**: US1 is independently testable — the panel is global, the old embed is gone
everywhere, and nothing that was visible before (status/conditions) has disappeared.

---

## Phase 4: User Story 2 - The AI panel starts from rich, automatic context (Priority: P1)

**Goal**: Opening the panel while viewing an object pre-fills a rich description of it; a
per-object "Ask AI about this" quick-action opens the panel pre-filled in one click.

**Independent Test**: Open the panel from a specific object's detail screen and confirm the
pre-fill covers identity + status/conditions + key spec fields; open it from a list/overview screen
and confirm it's empty/general instead; edit the pre-fill and confirm the edit (not the original)
is what sends; click a detail screen's quick-action and confirm one click opens the panel
pre-filled.

### Tests for User Story 2

- [ ] T017 [P] [US2] Jest test for `use-current-object-context.ts`: registers on mount, unregisters
  on unmount, and — critically — the panel's `queryField` only auto-refreshes while
  `queryFieldTouched` is `false`, in `front/app/ui/dashboard/ai-panel/use-current-object-context.test.ts`.
- [ ] T018 [P] [US2] Jest test for `ask-ai-button.tsx`: clicking it opens the panel with `queryField`
  already set from that screen's `ObjectContext`.

### Implementation for User Story 2

- [ ] T019 [US2] Create `front/app/ui/dashboard/ai-panel/use-current-object-context.ts`: hook each
  detail screen calls with its `ObjectContext` (kind/name/namespace/status/keySpecFields per
  data-model.md); updates `AIPanelState.currentObjectContext` and refreshes `queryField` only when
  `queryFieldTouched` is `false` (FR-006, FR-009); sets `queryFieldTouched` on manual edit (FR-008);
  clears registration on unmount so a list/overview screen's absence of a call means no context
  (FR-007) (depends on T002).
- [ ] T020 [P] [US2] Create `front/app/ui/dashboard/ai-panel/ask-ai-button.tsx`: the per-screen
  quick-action calling `useAIPanel().open()` with `queryField` seeded from that screen's
  `ObjectContext` (FR-016) (depends on T002, T019).
- [ ] T021 [P] [US2] Wire `use-current-object-context` + `AskAIButton` into
  `clusters/details.tsx` (identity/status/key spec fields for a Cluster).
- [ ] T022 [P] [US2] Wire into `clusters/infra/infra-details.tsx` (vSphere ClusterInfra fields:
  server, thumbprint, control plane endpoint).
- [ ] T023 [P] [US2] Wire into `clusters/infra/docker-details.tsx` (Docker ClusterInfra fields:
  load balancer IP, ready).
- [ ] T024 [P] [US2] Wire into `machines/details.tsx` (Machine fields: node name, provider ID,
  bootstrap, version).
- [ ] T025 [P] [US2] Wire into `machines/infra/infra-details.tsx` (vSphere MachineInfra fields:
  template, CPU/memory/disk).
- [ ] T026 [P] [US2] Wire into `machines/infra/docker-details.tsx` (Docker MachineInfra fields:
  provider ID, ready).
- [ ] T027 [P] [US2] Wire into `mds/details.tsx` (MachineDeployment fields: bootstrap/infra
  template refs, version).

**Checkpoint**: US1 + US2 together are independently testable — the panel is global and starts
from rich, correct, editable context everywhere it's opened.

---

## Phase 5: User Story 3 - Inspect the complete raw object on any detail screen (Priority: P2)

**Goal**: Every object detail screen gains a "YAML" tab rendering the complete underlying object
(not the curated subset) as an expandable/collapsible tree.

**Independent Test**: Open the detail view for a Cluster, a Machine, and a Machine Deployment; on
each, open the new tab and confirm the complete object (every field the backend returns, not a
curated subset) renders as a readable, expandable/collapsible tree, remains responsive for large
objects, and reflects live updates.

### Tests for User Story 3

- [ ] T028 [P] [US3] Go test for `HandleRawObject` in
  `webserver/internal/web/handlers/kubernetes/raw_test.go`: missing required query param → 400;
  unresolvable group/version/resource → 400/404 (per contracts/raw-object-api.md).
- [ ] T029 [P] [US3] Jest test for `to-tree-data.ts`: nested objects, arrays, scalar leaves, and an
  empty object all convert to valid `TreeNodeData[]` with unique `value` paths, in
  `front/app/ui/dashboard/shared/to-tree-data.test.ts`.
- [ ] T030 [P] [US3] Jest test for `object-tree.tsx`: renders a tree for a fetched object, shows a
  loading state before the fetch resolves, and re-fetches when the passed `resourceVersion` prop
  changes, in `front/app/ui/dashboard/shared/object-tree.test.tsx`.

### Implementation for User Story 3

- [ ] T031 [US3] Add `HandleRawObject` in `webserver/internal/web/handlers/kubernetes/raw.go`:
  parses `group`/`version`/`resource`/`namespace`/`name` query params, calls
  `clusterapi.NewDynamicClient(ctx).Resource(gvr).Namespace(ns).Get(ctx, name, metav1.GetOptions{})`,
  writes `obj.Object` as JSON; register `GET /api/raw` in
  `webserver/internal/web/handlers/handlers.go` (contracts/raw-object-api.md).
- [ ] T032 [P] [US3] Create `front/app/ui/dashboard/shared/to-tree-data.ts`: converts arbitrary
  JSON into Mantine `TreeNodeData[]` (data-model.md).
- [ ] T033 [US3] Create `front/app/ui/dashboard/shared/object-tree.tsx`: fetches `/api/raw` with the
  GVR/namespace/name for the given object on first mount, renders via Mantine's `Tree`, and
  re-fetches when the object's `resourceVersion` changes (research.md R2) (depends on T031, T032).
- [ ] T034 [P] [US3] Add a small frontend GVR constant table (alongside `object-tree.tsx` or in
  `front/app/lib/`) mapping each screen to its `{group, version, resource}` (Cluster, DockerCluster,
  VSphereCluster, Machine, DockerMachine, VSphereMachine, MachineDeployment) per
  contracts/raw-object-api.md.
- [ ] T035 [P] [US3] Add the "YAML" tab (rendering `<ObjectTree .../>`) to `clusters/details.tsx`
  (depends on T033, T034).
- [ ] T036 [P] [US3] Same for `clusters/infra/infra-details.tsx`.
- [ ] T037 [P] [US3] Same for `clusters/infra/docker-details.tsx`.
- [ ] T038 [P] [US3] Same for `machines/details.tsx`.
- [ ] T039 [P] [US3] Same for `machines/infra/infra-details.tsx`.
- [ ] T040 [P] [US3] Same for `machines/infra/docker-details.tsx`.
- [ ] T041 [P] [US3] Same for `mds/details.tsx`.

**Checkpoint**: All three user stories are independently functional.

---

## Phase 6: Polish & Cross-Cutting Concerns

- [ ] T042 [P] Run through `quickstart.md` manually: all six detail-screen variants' YAML tab, and
  the global panel's open/collapse/reopen, auto-context, quick-action, and mobile-viewport behavior.
- [ ] T043 [P] Visually verify light/dark theme contrast of the restyled AI panel against the rest
  of the dashboard (FR-010, SC-004).
- [ ] T044 Run `make build`, `make run-tests-backend`, and `make run-tests-frontend` end-to-end
  before opening a PR (Constitution Principle V gate).

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — start immediately.
- **Foundational (Phase 2)**: Depends on Setup — BLOCKS all user stories (the panel shell must exist
  before US1 can remove the old embed or US2 can wire context into it).
- **User Stories (Phase 3-5)**: All depend on Foundational completion.
  - **US1 → US2**: US2's context-wiring tasks (T021-T027) touch the same seven detail files as
    US1's tab-removal tasks (T010-T016) — do US1 first per file, or coordinate carefully if
    parallelizing across developers.
  - **US3** has no dependency on US1/US2 (different files: backend `raw.go`, new `shared/`
    components) and can proceed in parallel with either.
- **Polish (Phase 6)**: Depends on all desired user stories being complete.

### Parallel Opportunities

- T003, T004 (Foundational, marked [P]) run in parallel once T002 exists.
- Within US1, T010-T016 (seven different detail files) all run in parallel.
- Within US2, T021-T027 (seven different detail files) all run in parallel once T019/T020 exist.
- Within US3, T035-T041 (seven different detail files) all run in parallel once T033/T034 exist;
  T028-T030 (tests) and T031-T032 (backend handler + tree-data util) can proceed in parallel with
  each other.
- US3's entire phase can run in parallel with US1+US2 (disjoint files) if staffed separately.

---

## Parallel Example: User Story 1

```bash
# Seven independent file edits, same change shape, run together:
Task: "Remove AI Troubleshooting tab from clusters/details.tsx"
Task: "Remove AI Troubleshooting tab from clusters/infra/infra-details.tsx"
Task: "Remove AI Troubleshooting tab from machines/details.tsx"
Task: "Remove AI Troubleshooting tab from mds/details.tsx"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup.
2. Complete Phase 2: Foundational (blocks everything).
3. Complete Phase 3: User Story 1.
4. **STOP and VALIDATE**: run quickstart.md's global-panel scenarios (1-6); confirm no detail screen
   still has the old embedded tab.
5. Deploy/demo if ready — fixes the core "AI help isn't reachable everywhere" problem.

### Incremental Delivery

1. Setup + Foundational → global panel shell exists.
2. US1 → validate independently → deploy (MVP: panel is global, old embed gone).
3. US2 → validate independently → deploy (adds rich auto-context + quick-action).
4. US3 → validate independently → deploy (adds the YAML tree tab) — can be built in parallel with
   US1/US2 by a second developer since it touches disjoint files.

---

## Notes

- [P] tasks touch different files with no unmet dependencies.
- [Story] labels map every user-story-phase task back to spec.md for traceability.
- No task in this list requires a new dependency or a Mantine version bump (research.md R3/R4).
- Commit after each task or logical group; stop at any checkpoint to validate a story independently.
