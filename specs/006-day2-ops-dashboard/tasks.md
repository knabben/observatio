---

description: "Task list for feature implementation"
---

# Tasks: Day-2 Operations Dashboard

**Input**: Design documents from `/specs/006-day2-ops-dashboard/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md

**Tests**: Constitution Principle V (Test-Driven Quality) mandates test coverage for all backend
and frontend changes, so test tasks are included and MUST be written first and MUST fail before
their corresponding implementation task.

**Organization**: Tasks are grouped by user story (US1–US5, priorities from spec.md) to enable
independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1, US2, US3, US4, US5)

## Path Conventions

Web app: `webserver/` (Go backend), `front/` (Next.js frontend) — per plan.md Project Structure.

---

## Phase 1: Setup

**Purpose**: Project/dependency initialization per plan.md

- [X] T001 Promote `k8s.io/apiextensions-apiserver` from indirect to direct dependency in `webserver/go.mod` (`go get k8s.io/apiextensions-apiserver@v0.31.3 && go mod tidy`), needed for CRD version introspection (research.md R6)
- [X] T002 [P] Create backend package skeleton `webserver/internal/infra/clusterapi/day2ops/` (package doc comment only, per plan.md Project Structure)
- [X] T003 [P] Create frontend module skeleton directory `front/app/ui/dashboard/components/ops/` (per plan.md Project Structure)

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Shared WS envelope, types, and drill-in scaffolding that every user story's detectors plug into

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T004 Define `Day2OpsEvent`, `HealthRollup`, `DebugLayer`, `DebugPath`, `RiskWarning`, `FailureSeverity` Go types in `webserver/internal/infra/clusterapi/day2ops/types.go` per data-model.md
- [X] T005 [P] Define matching TypeScript types (`Day2OpsEvent`, `HealthRollup`, `DebugLayer`, `DebugPath`, `RiskWarning`, `FailureSeverity`) in `front/app/ui/dashboard/shared/use-day2-ops.ts` per data-model.md
- [X] T006 [P] Add `'degraded'` to `StatusState` in `front/app/ui/dashboard/shared/status.ts`, render it in `front/app/ui/dashboard/shared/status-indicator.tsx`, and add its amber color to `STATUS_COLORS` in `front/app/styles/theme.ts`
- [X] T007 Implement `webserver/internal/web/watchers/day2ops.go`'s `WatchDay2Ops(ctx, conn, objType)`: open Cluster/Machine/MachineDeployment K8s watches directly (via `clusterapi.NewDynamicClient`, same primitive every existing watcher uses), fan their events into one loop, recompute a `Day2OpsEvent` (initially empty `Rollups`/`DebugPaths`/`Risks`/`Severities`, `SourceUnavailable` reflecting live API-server reachability) on each event, and write it to `conn`; register `day2ops` in `websocket.go`'s `watchHandlers` map per contracts/day2ops-ws-event.md (research.md R9 — corrected from an earlier "shared pool" design; depends on T004)
- [X] T008 [P] Implement `use-day2-ops.ts` WS consumer hook treating each `day2ops` event as a full-state replace, per the contracts/day2ops-ws-event.md consumer contract (depends on T005)
- [X] T009 Implement `GET /api/day2ops/detail` handler skeleton in `webserver/internal/web/handlers/kubernetes/day2ops.go` (400 on invalid GVR/name, 404 on not-found, 500 otherwise, per contracts/day2ops-detail-api.md), registered in `webserver/internal/web/handlers/handlers.go` (depends on T004)
- [X] T010 [P] Go test for `day2ops.go` request parsing/validation (GVR + namespace + name) in `webserver/internal/web/handlers/kubernetes/day2ops_test.go`, mirroring the existing `raw_test.go` table-test pattern (depends on T009)

**Checkpoint**: WS envelope and REST drill-in scaffolding work end-to-end with placeholder data — user stories now add real detection logic on top.

---

## Phase 3: User Story 1 - Centralized landing view (Priority: P1) 🎯 MVP

**Goal**: A redesigned `front/app/dashboard/page.tsx` shows live, per-category (Cluster, MachineDeployment, Machine) healthy/degraded/failed rollups, narrows in place per category, and shows explicit "all clear" and "data unavailable" states.

**Independent Test**: Seed a mix of healthy/unhealthy objects across categories, open the dashboard, and confirm every category's rollup is visible and accurate without navigating to any other page (per spec.md US1 Acceptance Scenarios).

### Tests for User Story 1

- [X] T011 [P] [US1] Go test for rollup computation (healthy/degraded/failed counts per category, `Unavailable` flag) in `webserver/internal/infra/clusterapi/day2ops/rollup_test.go`
- [X] T012 [P] [US1] Jest test for `health-rollup-card.tsx` rendering healthy/degraded/failed/unavailable states in `front/app/ui/dashboard/components/ops/health-rollup-card.test.tsx`
- [X] T013 [P] [US1] Jest test for `ops-dashboard.tsx` in-place category narrowing and "all clear"/"data unavailable" states in `front/app/ui/dashboard/components/ops/ops-dashboard.test.tsx`

### Implementation for User Story 1

- [X] T014 [US1] Implement `webserver/internal/infra/clusterapi/day2ops/rollup.go`: compute `HealthRollup[]` per category, reusing existing counts from `webserver/internal/infra/clusterapi/dashboard.go`'s `GenerateClusterSummary` plus new `Degraded`/`Unavailable` logic (depends on T004; T011 must fail first)
- [X] T015 [US1] Wire `rollup.go` into `day2ops.go`'s `Day2OpsEvent.Data.Rollups` and `SourceUnavailable` (depends on T007, T014)
- [X] T016 [US1] Implement `front/app/ui/dashboard/components/ops/health-rollup-card.tsx` (depends on T012)
- [X] T017 [US1] Implement `front/app/ui/dashboard/components/ops/ops-dashboard.tsx`: orchestrator consuming `use-day2-ops.ts`, rendering category rollup cards, in-place category narrowing, and "all clear"/"data unavailable" states (depends on T008, T016, T013)
- [X] T018 [US1] Redesign `front/app/dashboard/page.tsx` to render `<OpsDashboard/>`, re-homing the existing `ClusterSummary`/`ClusterClassLister`/`ClusterHierarchy`/`ClusterVersions` widgets as constituent rollup building blocks (depends on T017)

**Checkpoint**: User Story 1 is fully functional and independently testable/demoable.

---

## Phase 4: User Story 2 - Layered root-cause debugging path (Priority: P1)

**Goal**: For any unhealthy object, the dashboard shows which layer(s) — conditions, phase, provider resource, controller activity — explain the failure, using only dashboard data.

**Independent Test**: Seed a known failure at a specific layer (e.g., a `DockerMachine` provider error) and confirm the dashboard highlights that exact layer with supporting evidence, using only the dashboard UI (per spec.md US2 Acceptance Scenarios).

### Tests for User Story 2

- [X] T019 [P] [US2] Go test for `debugpath.go` layer synthesis — conditions/phase/provider_resource "implicated" cases, and `controller_activity` populated only when earlier layers are inconclusive — in `webserver/internal/infra/clusterapi/day2ops/debugpath_test.go`
- [X] T020 [P] [US2] Go test for `GET /api/day2ops/detail` returning a full `DebugPath` payload, extending `day2ops_test.go`
- [X] T021 [P] [US2] Jest test for `debugging-path.tsx` rendering ordered, layer-labeled evidence in `front/app/ui/dashboard/components/ops/debugging-path.test.tsx`

### Implementation for User Story 2

- [X] T022 [US2] Implement `webserver/internal/infra/clusterapi/day2ops/debugpath.go`: synthesize `DebugLayer[]` from Machine `Conditions`+`Phase` and the provider-infra object's `Conditions` (already-watched data, research.md R1) (depends on T004; T019 must fail first)
- [X] T023 [US2] Add `webserver/internal/infra/clusterapi/fetchers/events.go`: read `corev1.Event`s scoped to `involvedObject`, feeding `debugpath.go`'s `controller_activity` layer only when conditions/phase/provider_resource are inconclusive (research.md R2, FR-007) (depends on T022)
- [X] T024 [US2] Wire `debugpath.go` into (a) `day2ops.go`'s `Day2OpsEvent.Data.DebugPaths` with evidence capped to one line per layer, so layers/status render directly on the landing screen per FR-004 (data-model.md, contracts/day2ops-ws-event.md), and (b) the `GET /api/day2ops/detail` handler, returning the full, uncapped `DebugPath` JSON per contracts/day2ops-detail-api.md for on-demand expansion (depends on T007, T009, T022, T023; T020 must fail first)
- [X] T025 [US2] Implement `front/app/ui/dashboard/components/ops/debugging-path.tsx`: renders the capped `DebugPath` inline (from `use-day2-ops.ts`) for every unhealthy object on the landing screen, with an expand action that fetches the full evidence list from the detail endpoint on demand (depends on T021)
- [X] T026 [US2] Wire `debugging-path.tsx` inline into `ops-dashboard.tsx`'s unhealthy rollup entries (depends on T017, T025)

**Checkpoint**: User Story 2 is fully functional independently of US3/US4.

---

## Phase 5: User Story 3 - Proactive risk detection (Priority: P2)

**Goal**: Certificate expiry, stalled rollouts, provider/CRD version skew, and infrastructure drift are each detected and flagged, with an explicit "check could not be performed" state when a check can't run.

**Independent Test**: Seed each of the four risk conditions independently and confirm each is surfaced as a distinct, correctly categorized warning (per spec.md US3 Acceptance Scenarios).

### Tests for User Story 3

- [X] T027 [P] [US3] Go test for `risk_certexpiry.go` (warning-window boundary, not-evaluable case) in `webserver/internal/infra/clusterapi/day2ops/risk_certexpiry_test.go`
- [X] T028 [P] [US3] Go test for `risk_stalledrollout.go` (grace-period boundary, PDB/finalizer likely-cause detection) in `webserver/internal/infra/clusterapi/day2ops/risk_stalledrollout_test.go`
- [X] T029 [P] [US3] Go test for `risk_versionskew.go` in `webserver/internal/infra/clusterapi/day2ops/risk_versionskew_test.go`
- [X] T030 [P] [US3] Go test for `risk_drift.go` (`generation`/`observedGeneration` mismatch) in `webserver/internal/infra/clusterapi/day2ops/risk_drift_test.go`
- [X] T031 [P] [US3] Jest test for `risk-warnings.tsx` grouping/rendering, including the `not_evaluable` state, in `front/app/ui/dashboard/components/ops/risk-warnings.test.tsx`

### Implementation for User Story 3

- [X] T032 [P] [US3] Add `webserver/internal/web/watchers/machineset.go`, following the existing `machinedeployment.go` pattern (research.md R5)
- [X] T033 [P] [US3] Add `webserver/internal/infra/clusterapi/fetchers/secrets.go`: read-only fetch scoped to CAPI-managed CA/etcd/proxy cert Secrets (research.md R4)
- [X] T034 [US3] Implement `webserver/internal/infra/clusterapi/day2ops/risk_certexpiry.go` using T033 (depends on T033; T027 must fail first)
- [X] T035 [US3] Implement `webserver/internal/infra/clusterapi/day2ops/risk_stalledrollout.go` using the T032 MachineSet watcher plus a PodDisruptionBudget (`policy/v1`) lookup for the likely-cause heuristic (depends on T032; T028 must fail first)
- [X] T036 [US3] Implement `webserver/internal/infra/clusterapi/day2ops/risk_versionskew.go`, extending the existing `clusterctlv1.ProviderList` read with a CRD-version check via the `apiextensions-apiserver` clientset (depends on T001; T029 must fail first)
- [X] T037 [US3] Implement `webserver/internal/infra/clusterapi/day2ops/risk_drift.go` using `generation`/`observedGeneration` already present on watched objects (research.md R3) (T030 must fail first)
- [X] T038 [US3] Wire all four risk detectors into `day2ops.go`'s `Day2OpsEvent.Data.Risks` (depends on T007, T034, T035, T036, T037)
- [X] T039 [US3] Implement `front/app/ui/dashboard/components/ops/risk-warnings.tsx`, grouped by category, including the `not_evaluable` state (depends on T031)
- [X] T040 [US3] Wire `risk-warnings.tsx` into `ops-dashboard.tsx` (depends on T017, T039)

**Checkpoint**: User Story 3 is fully functional independently of US2/US4.

---

## Phase 6: User Story 4 - Failure-severity awareness (Priority: P3)

**Goal**: Detected issues are classified as self-healing, needs-investigation, provider-degraded, or management-critical, with escalating and visually distinct urgency.

**Independent Test**: Simulate one condition from each severity level and confirm the dashboard labels each with a distinct, correctly escalating severity and guidance (per spec.md US4 Acceptance Scenarios).

### Tests for User Story 4

- [X] T041 [P] [US4] Go test for `severity.go` classification (self_healing / needs_investigation / provider_degraded / management_critical, including CA-secret-missing → management_critical) in `webserver/internal/infra/clusterapi/day2ops/severity_test.go`
- [X] T042 [P] [US4] Jest test for `severity-banner.tsx` rendering escalating urgency and never showing self-healing activity with alert-level styling in `front/app/ui/dashboard/components/ops/severity-banner.test.tsx`

### Implementation for User Story 4

- [X] T043 [P] [US4] Add `webserver/internal/web/watchers/machinehealthcheck.go`, watching the first-class CAPI `MachineHealthCheck` type (research.md R7)
- [X] T044 [US4] Implement `webserver/internal/infra/clusterapi/day2ops/severity.go`: self-healing/needs-investigation from `MachineHealthCheck` status + `maxUnhealthy` breach, provider-degraded from Pod/Deployment status in known controller namespaces (`capi-system`, provider namespaces), management-critical from live API-server reachability, and CA-secret-missing (reusing T033's Secret fetcher) (depends on T043, T033; T041 must fail first)
- [X] T045 [US4] Wire `severity.go` into `day2ops.go`'s `Day2OpsEvent.Data.Severities` (depends on T007, T044)
- [X] T046 [US4] Implement `front/app/ui/dashboard/components/ops/severity-banner.tsx` (depends on T042)
- [X] T047 [US4] Wire `severity-banner.tsx` into `ops-dashboard.tsx` as the top-level, hard-to-miss banner (depends on T017, T046)

**Checkpoint**: User Stories 1–4 are independently functional.

---

## Phase 7: User Story 5 - Deep-dive into controller logs (Priority: P2)

**Goal**: A new "Logs" destination (lateral nav item + debugging-path deep-dive) streams the relevant controller's Pod log output via the standard Kubernetes Pod-log subresource, and shows static SSH connection instructions for a Machine's node access — no Docker daemon access, no stored credentials, Machine/node-level log streaming explicitly deferred as a TODO.

**Independent Test**: From an object whose debugging path implicates the controller-activity layer, open the deep-dive action and confirm the correct controller's log output streams in; separately confirm the Logs view is directly reachable from the lateral navigation, and that a VM-based-provider Machine's node-access deep-dive shows only SSH connection instructions (per spec.md US5 Acceptance Scenarios).

### Tests for User Story 5

- [X] T048 [P] [US5] Go test for `controllerlogs.go` (resolves Deployment → current Pod, returns/streams logs, 404 when no Pod backs the Deployment, 503 when logs can't be retrieved) in `webserver/internal/infra/clusterapi/fetchers/controllerlogs_test.go`
- [X] T049 [P] [US5] Go test for `logs.go` handler (`GET /api/logs/controller` query validation and status codes; `GET /api/logs/node-access` response shape) in `webserver/internal/web/handlers/kubernetes/logs_test.go`
- [X] T050 [P] [US5] Jest test for `logs-view.tsx` (controller selection, streamed log pane, "logs unavailable" state) in `front/app/ui/dashboard/logs/logs-view.test.tsx`
- [X] T051 [P] [US5] Jest test for `node-access-panel.tsx` (renders only the SSH command/address/disclaimer, no credential input of any kind) in `front/app/ui/dashboard/logs/node-access-panel.test.tsx`

### Implementation for User Story 5

- [X] T052 [US5] Implement `webserver/internal/infra/clusterapi/fetchers/controllerlogs.go`: resolve a controller Deployment's current Pod and return/stream its logs via the standard Kubernetes Pod-log subresource (research.md R10) (depends on T048 must fail first)
- [X] T053 [US5] Implement `webserver/internal/web/handlers/kubernetes/logs.go`: `GET /api/logs/controller` (depends on T052) and `GET /api/logs/node-access` (reading `Machine.status.addresses`), per contracts/logs-api.md, registered in `webserver/internal/web/handlers/handlers.go` (depends on T052; T049 must fail first)
- [X] T054 [US5] Implement `front/app/ui/dashboard/logs/logs-view.tsx`: controller selector + streamed log pane + "logs unavailable" state (depends on T050)
- [X] T055 [US5] Implement `front/app/ui/dashboard/logs/node-access-panel.tsx` (depends on T051)
- [X] T056 [US5] Add `front/app/dashboard/logs/layout.tsx` + `page.tsx` rendering `<LogsView/>`, and add the new "Logs" entry to `front/app/ui/dashboard/nav-links.tsx` (depends on T054)
- [X] T057 [US5] Wire the "deep dive" action from `debugging-path.tsx` into the Logs view when `controller_activity` is implicated (linking to the correct controller), and into `node-access-panel.tsx` for a Machine's node-access deep-dive (depends on T025, T056, T055)

**Checkpoint**: All five user stories are independently functional.

---

## Phase 8: Polish & Cross-Cutting Concerns

- [X] T058 [P] Run quickstart.md validation scenarios end-to-end against a `kind-capi-mgmt`-style test cluster
- [X] T059 [P] Verify `make build`, `make run-tests-backend`, and `make run-tests-frontend` all pass
- [X] T060 Annotate this tasks.md with any deviations discovered mid-implementation (project convention from features 004/005)

### Discovered mid-implementation

- **T007 architecture correction (research.md R9)**: The plan originally described a shared
  broadcast "pool" for the aggregator. Inspecting the actual watcher architecture
  (`webserver/internal/web/handlers/system/websocket.go`) showed every existing resource watcher
  is a per-connection 1:1 relay — there is no shared pub/sub bus; `pool.go` is used only by the AI
  chatbot. `WatchDay2Ops` was implemented as a new entry in the same `watchHandlers` dispatch table
  (per-connection fan-in of several GVRs into one stream), not a shared aggregator.
- **T032/T043 — no standalone `machineset.go`/`machinehealthcheck.go` watcher files**: these kinds
  are only ever needed internally by the Day2Ops aggregator (not as their own list pages), so their
  GVRs and store-tracking were added directly inside `webserver/internal/web/watchers/day2ops.go`
  rather than creating unused standalone watcher functions matching the 1:1-relay convention.
- **T035 — PDB-based likely-cause not implemented**: only the finalizer-based likely-cause for a
  stalled rollout was built; checking for a blocking PodDisruptionBudget would require workload
  cluster access, out of scope per spec.md's Assumptions.
- **T036 — version-skew heuristic changed from the original plan**: comparing a provider's release
  version against its CRD API version conflates two unrelated versioning schemes and can't be done
  correctly in general. Implemented instead as: a CRD's `status.storedVersions` containing a
  version no longer in `spec.versions[].served` — a directly observable, well-defined upgrade
  hazard, achievable with the same `apiextensions-apiserver` client (research.md R6 updated).
- **Live bug — a single missing/optional CRD aborted the entire Day2Ops connection**: the first
  live test against a real Docker-only `kind-capi-mgmt` cluster (no vSphere provider installed)
  showed `WatchDay2Ops` failing outright with "the server could not find the requested resource"
  for `vspheremachines`, because every watched GVR was opened unconditionally and any single
  failure aborted the whole connection — the dashboard showed nothing at all, with the browser
  reporting a raw WebSocket error event. Fixed by (a) making per-GVR watch failures non-fatal
  (skip and log, only fail if *zero* watches could be opened at all), and (b) proactively detecting
  installed infrastructure providers via the existing `clusterapi.GenerateInfrastructureCapability`
  mechanism (same one used by `/api/infra/capabilities`) so provider-specific GVRs are only
  attempted for providers actually installed, rather than reactively discovering the absence of
  every possible provider's CRDs by trial and error.
- **Live bug — `GET /api/clusters/topology` (pre-existing, unrelated to this feature) also
  hardcoded a VSphereMachine-only seed**: found while investigating the above; the same class of
  bug caused the existing Cluster Topology widget to 500 on any non-vSphere cluster.
  `clusterapi.GenerateClusterTopology`/`processOwnerHierarchy` were made provider-agnostic (seeded
  from core `Machine` objects instead of `VSphereMachine`), and node coloring was extended to
  reflect actual object health (red for unhealthy, via a generic `status.ready`/`Ready`-condition
  check) rather than only the previous position-based layer palette. This was reported live by the
  user and fixed as part of this feature's Polish phase, though it is not itself a Day-2 Ops
  dashboard requirement.

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — can start immediately
- **Foundational (Phase 2)**: Depends on Setup — BLOCKS all user stories
- **User Stories (Phase 3–7)**: All depend on Foundational; independently implementable/testable in any order thereafter (P1 stories US1/US2 recommended first, matching spec.md priorities)
- **Polish (Phase 8)**: Depends on all desired user stories being complete

### User Story Dependencies

- **US1 (P1)**: No dependency on US2/US3/US4/US5. US2/US3/US4's dashboard wiring tasks (T026, T040, T047) depend on `ops-dashboard.tsx` (T017) existing, but US2/US3/US4's detection logic itself has no dependency on US1.
- **US2 (P1)**: Independent detection logic; only its dashboard wiring (T026) depends on T017.
- **US3 (P2)**: Independent detection logic; only its dashboard wiring (T040) depends on T017.
- **US4 (P3)**: Independent detection logic; only its dashboard wiring (T047) depends on T017.
- **US5 (P2)**: Independent Logs/node-access implementation (T048–T056 have no dependency on any other story); only its final wiring task (T057) depends on US2's `debugging-path.tsx` (T025) existing, since that's where the deep-dive action lives.

### Parallel Opportunities

- T002, T003 (Setup) in parallel.
- T005, T006 (Foundational, frontend) in parallel with T004/T007/T009/T010 (Foundational, backend).
- Once Foundational completes, US1/US2/US3/US4/US5's detection/implementation tasks (everything except each story's final dashboard-wiring task) can proceed fully in parallel across stories.
- Within each story, all `[P]`-marked test tasks run in parallel; within US3, T032/T033 (new watcher + new fetcher) run in parallel; within US4, T043 runs in parallel with nothing else in that story (T044 depends on it); within US5, T048–T051 (all four tests) run in parallel, and T054/T055 (the two frontend components) run in parallel.

---

## Parallel Example: User Story 3

```bash
# Tests (parallel):
Task: "Go test for risk_certexpiry.go in webserver/internal/infra/clusterapi/day2ops/risk_certexpiry_test.go"
Task: "Go test for risk_stalledrollout.go in webserver/internal/infra/clusterapi/day2ops/risk_stalledrollout_test.go"
Task: "Go test for risk_versionskew.go in webserver/internal/infra/clusterapi/day2ops/risk_versionskew_test.go"
Task: "Go test for risk_drift.go in webserver/internal/infra/clusterapi/day2ops/risk_drift_test.go"
Task: "Jest test for risk-warnings.tsx in front/app/ui/dashboard/components/ops/risk-warnings.test.tsx"

# New watcher + fetcher (parallel):
Task: "Add MachineSet watcher in webserver/internal/web/watchers/machineset.go"
Task: "Add Secret fetcher in webserver/internal/infra/clusterapi/fetchers/secrets.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (blocks all stories)
3. Complete Phase 3: User Story 1
4. **STOP and VALIDATE**: confirm the consolidated landing view against quickstart.md's US1 scenarios
5. Deploy/demo if ready — this alone replaces "visit 3 separate list pages" with one screen

### Incremental Delivery

1. Setup + Foundational → WS/REST scaffolding ready
2. US1 → consolidated landing view (MVP)
3. US2 → layered debugging path (the feature's core differentiator)
4. US3 → proactive risk detection
5. US4 → severity classification/prioritization
6. US5 → controller-logs deep-dive (depends on US2's debugging-path UI existing for its wiring step, otherwise independent)
7. Each story adds value without breaking previously delivered stories, since detection logic is additive to the shared `Day2OpsEvent` and dashboard wiring is a single append per story

### Parallel Team Strategy

With multiple developers, after Foundational completes: one developer per user story (US1, US2, US3, US4, US5), since detection logic is independent per story and only converges at each story's single dashboard-wiring task (US5's final task additionally waits on US2's debugging-path component).
