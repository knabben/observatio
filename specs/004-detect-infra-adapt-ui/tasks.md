---

description: "Task list for feature implementation"
---

# Tasks: Infrastructure Provider Detection & Adaptive Listing Screens

**Input**: Design documents from `/specs/004-detect-infra-adapt-ui/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/infra-detection-api.md, quickstart.md

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
Structure for exact new/changed files.

---

## Phase 1: Setup

**Purpose**: De-risk the one open technical question from research.md (R3) and confirm a clean
starting baseline before touching any code.

- [X] T001 [P] Verify `sigs.k8s.io/cluster-api/test/infrastructure/docker/api/v1beta1` (Docker
  `DockerCluster`/`DockerMachine` types) resolves and compiles under the pinned `v1.9.6` module.
  **Result: does NOT resolve as a subpackage** — confirmed via `go get`, which revealed
  `test/infrastructure/docker` is a separate nested Go module and pulled `cluster-api` v1.9.6→v1.13.3
  plus a cascade of `k8s.io/*`/`controller-runtime` upgrades. Reverted `go.mod`/`go.sum` immediately.
  Falling back to the dynamic/unstructured client approach (research.md R3, revised) for T015/T016/T018 —
  no `dockerv1` scheme registration (drop from T005).
- [ ] T002 [P] Confirm `make build`, `make run-tests-backend`, and `make run-tests-frontend` all pass
  on a clean checkout of `004-detect-infra-adapt-ui` before starting implementation.

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Shared provider-detection primitives every user story depends on.

**⚠️ CRITICAL**: No user story work can begin until this phase is complete.

- [X] T003 Create `webserver/internal/infra/providerkind/providerkind.go` with
  `FromKind(kind string) string`, mapping `DockerCluster`/`DockerMachine` → `docker`,
  `VSphereCluster`/`VSphereMachine` → `vsphere`, anything else (including empty) → `unknown`.
- [X] T004 [P] Add unit tests in `webserver/internal/infra/providerkind/providerkind_test.go`
  covering docker/vsphere/unknown/empty-string `Kind` inputs.
- [X] T005 ~~Register the Docker infra scheme~~ — SKIPPED per T001 finding: Docker infra objects are
  read via the dynamic/unstructured client (no typed `dockerv1` package), so no `runtime.Scheme`
  registration is needed or possible for them.
- [X] T006 Add `InfrastructureCapability` and `ProviderStatus` structs to
  `webserver/internal/infra/models/capability.go` per data-model.md.
- [X] T007 Add `GenerateInfrastructureCapability(ctx, c) (models.InfrastructureCapability, error)` in
  `webserver/internal/infra/clusterapi/dashboard.go`, reusing the existing
  `clusterctlv1.ProviderList` lookup already performed by `GenerateComponentVersions`
  (filter `Kind == "InfrastructureProvider"` and `ProviderName` in `{docker, vsphere}`; research.md R1).
- [X] T008 [P] Add tests for `GenerateInfrastructureCapability` in
  `webserver/internal/infra/clusterapi/dashboard_test.go` using a CAPI fake client seeded with
  `clusterctlv1.Provider` fixtures: docker-only, vsphere-only, both, neither installed.
- [X] T009 Add `HandleInfraCapabilities` in `webserver/internal/web/handlers/kubernetes/dashboard.go`
  and register `GET /api/infra/capabilities` in `webserver/internal/web/handlers/handlers.go`
  (depends on T007).
- [X] T010 Add a `Provider string` field to `models.Cluster`
  (`webserver/internal/infra/models/cluster.go`) and populate it via
  `providerkind.FromKind(InfrastructureRef.Kind)` wherever `models.Cluster` is built (the processor
  that already sets `InfrastructureRef`) (depends on T003).
- [X] T011 Add a `Provider string` field to `models.Machine`
  (`webserver/internal/infra/models/machine.go`) and populate it the same way from the raw
  `clusterv1.Machine.Spec.InfrastructureRef.Kind` (depends on T003).
- [X] T012 [P] Add `getInfraCapabilities()` + `InfrastructureCapability`/`ProviderStatus` types to
  `front/app/lib/data.tsx` (co-located with every other REST fetch helper in this codebase, rather
  than a new `capabilities.ts` file, per existing convention) for `GET /api/infra/capabilities`.

**Checkpoint**: Foundation ready — all user stories can now proceed.

---

## Phase 3: User Story 1 - Listing screens reflect the actual infrastructure provider in use (Priority: P1) 🎯 MVP

**Goal**: Docker-only, vSphere-only, and mixed environments each render the correct provider-specific
infra tab/view — never the wrong one, never a static hardcoded one.

**Independent Test**: Point the dashboard at a Docker-backed environment and confirm a
Docker-appropriate infra view appears (no vSphere-labeled tab); point it at vSphere and confirm the
existing view is unchanged; point it at a mixed environment and confirm both appear, correctly
scoped.

### Tests for User Story 1

- [X] T013 [P] [US1] Backend test for the dispatch/selection logic (auto-selects the installed
  provider, honors `?provider=` override, rejects an unrecognized value) in
  `webserver/internal/web/handlers/kubernetes/cluster_test.go`. Scoped to `resolveInfraProvider`
  directly rather than the full HTTP handler: `HandleClusterInfraList`/`HandleMachineInfra` always
  construct a real `client.Client` via `clusterapi.NewClientWithScheme` (no fake-client seam),
  matching this codebase's existing pattern where handlers themselves are untested and their pure
  logic is extracted and tested instead (e.g. `ProcessCluster`).
- [X] T014 [P] [US1] Frontend test: `ClusterTabs` renders the correct tab set (docker-only,
  vsphere-only, both, neither) from a mocked `getInfraCapabilities()` response, in
  `front/app/ui/dashboard/components/clusters/cluster-tabs.test.tsx` (co-located with the
  component it tests, per this codebase's convention, rather than a `page.test.tsx`).

### Implementation for User Story 1

- [X] T015 [P] [US1] Add `models.ClusterInfraDocker` struct to
  `webserver/internal/infra/models/cluster.go` per data-model.md (revised: `Cluster`,
  `LoadBalancerIP`, `Ready`, `Conditions` — no `dockerv1` type reference).
- [X] T016 [P] [US1] Add `fetchers.ListClusterInfraDocker` in
  `webserver/internal/infra/clusterapi/fetchers/cluster_infra_docker.go`, listing
  `DockerCluster` objects (`infrastructure.cluster.x-k8s.io/v1beta1`) via
  `clusterapi.NewDynamicClient` + `unstructured.UnstructuredList`, reading
  `spec.loadBalancerIP`/`status.ready`/`status.conditions` field-by-field (research.md R3 revised —
  no typed `dockerv1` package).
- [X] T017 [US1] Extend `HandleClusterInfraList` in
  `webserver/internal/web/handlers/kubernetes/cluster.go` to accept an optional `?provider=` query
  param, default to the first `installed` provider from `GenerateInfrastructureCapability`, dispatch
  to the Docker or vSphere fetcher accordingly, and return `404` when the resolved/requested provider
  isn't installed (depends on T007, T015, T016).
- [X] T018 [US1] Mirror the same dispatch for Machines: extend `HandleMachineInfra` in
  `webserver/internal/web/handlers/kubernetes/machine.go`, adding
  `fetchers.ListMachineInfraDocker` (dynamic client + unstructured decode of `DockerMachine`, same
  approach as T016) and `models.MachineInfraDocker` (FR-008) (depends on T007).
- [X] T019/T020 [P] [US1] Revised approach: rather than generalizing the vSphere-specific
  `infra-lister.tsx`/`infra-table.tsx` in place (risking FR-010's "existing vSphere behavior
  unchanged" guarantee, and awkward given Docker's flat `ready` vs vSphere's nested
  `status.ready`), added parallel Docker components mirroring the vSphere trio exactly:
  `clusters/infra/docker-{lister,table,details}.tsx` and
  `machines/infra/docker-{lister,table,details}.tsx`, plus `ClusterInfraDockerType`/
  `MachineInfraDockerType` in each `types.tsx`.
  **Also discovered mid-implementation (research.md R6)**: the existing infra listers are
  WebSocket-driven (`BaseLister` → `useResourceStream` → `ws://.../ws/watcher`), not REST — the
  `/api/clusters/infra/list`/`/api/machines/infra/list` REST endpoints (R4) are unused by the UI.
  Added `watchers.WatchDockerClusters`/`WatchDockerMachines` (dynamic-client watch on
  `dockerclusters`/`dockermachines`, reusing the existing `WatchResourceViaWebSocket` helper — so
  the earlier `ResourceVersion:"0"` fix applies automatically) and registered new WS object types
  `"cluster-infra-docker"`/`"machine-infra-docker"` in `system/websocket.go`. Exported
  `fetchers.ProcessDockerCluster`/`ProcessDockerMachine` so both the REST fetcher and the WS
  watcher share one decode path.
- [X] T021 [US1] Replaced the static 2-tab `Tabs` in `front/app/dashboard/clusters/page.tsx` with a
  new `ClusterTabs` client component driven by `getInfraCapabilities()` — a provider's tab renders
  only when `installed` is true; shows a "no supported infrastructure provider detected" message
  when neither is (covers US3's empty-state ahead of schedule). `page.tsx` is now a thin shell.
- [X] T022 [US1] Same pattern applied to Machines via a new `MachineTabs` component;
  `machines/page.tsx` is now a thin shell.

**Checkpoint**: US1 is independently testable — Docker/vSphere/mixed environments each show exactly
the right infra view(s).

---

## Phase 4: User Story 2 - Operators can identify a cluster's infrastructure provider at a glance (Priority: P2)

**Goal**: Every row in the main Clusters (and Machines) list shows a provider (+ version) indicator,
or "Unknown" when undetermined.

**Independent Test**: Load the main Clusters list against a mixed environment and confirm every row
shows a correct provider(+version) badge matching that resource's actual backing infrastructure.

### Tests for User Story 2

- [X] T023 [P] [US2] Frontend test for the provider badge component: correct label/version for
  docker/vsphere, "Unknown" for an unrecognized provider, in
  `front/app/ui/dashboard/shared/provider-badge.test.tsx`.

### Implementation for User Story 2

- [X] T024 [P] [US2] Create `front/app/ui/dashboard/shared/provider-badge.tsx`: takes a `provider`
  string and the capabilities response, renders "Docker vX.Y.Z" / "vSphere vX.Y.Z" / "Unknown".
- [X] T025 [US2] Add a `provider` field to `ClusterType` in
  `front/app/ui/dashboard/components/clusters/types.tsx` (mirrors backend T010).
- [X] T026 [US2] Wire `provider-badge.tsx` into
  `front/app/ui/dashboard/components/clusters/table.tsx` row rendering, resolving the version by
  looking up `provider` in the capabilities response (depends on T012, T024, T025).
- [X] T027 [US2] Wire the same badge into the Machines list table (depends on T011, T024).

**Checkpoint**: US2 is independently testable — every Clusters/Machines row shows a correct
provider(+version) or Unknown badge.

---

## Phase 5: User Story 3 - Unsupported or undetectable providers degrade gracefully (Priority: P3)

**Goal**: An unrecognized provider, or an environment with no supported provider installed, never
crashes or blanks the screen — it degrades to a clear, generic state.

**Independent Test**: Inject a cluster with an infrastructure reference outside the supported set,
and separately simulate an environment with neither provider installed; confirm both cases render a
clear, non-crashing result.

### Tests for User Story 3

- [X] T028 [P] [US3] Backend test: `TestProcessCluster_Provider`/`TestProcessMachine_Provider` in
  `processor/cluster_test.go`/`machine_test.go` confirm an unrecognized (or, for Cluster, absent)
  `infrastructureRef.kind` never errors/panics and resolves to `"unknown"`. The
  `GET /api/infra/capabilities` "neither installed" case was already covered by
  `Test_GenerateInfrastructureCapability` (T008) and `Test_resolveInfraProvider` (T013).
- [X] T029 [P] [US3] Frontend test: `cluster-tabs.test.tsx`/`machine-tabs.test.tsx` (T014) already
  cover the "no supported infrastructure provider detected" message for both screens; added an
  assertion to the existing `clusters/table.test.tsx` partial-cluster case confirming an
  unrecognized/absent-provider cluster still renders in the main list with an "Unknown" badge
  rather than being dropped.

### Implementation for User Story 3

- [X] T030 [US3] Already delivered as part of T021/T022 (`ClusterTabs`/`MachineTabs` render the
  "no supported infrastructure provider detected" `EmptyState` when both `docker.installed` and
  `vsphere.installed` are false).
- [X] T031 [US3] Verified by code review (consistent with T013's handler-testability note —
  `HandleClusterInfraList`/`HandleMachineInfra` always construct a real client, so this isn't
  fake-client-testable): both handlers check `capability.<Provider>.Installed` before dispatching
  and call `http.Error(w, ..., http.StatusNotFound)` otherwise — never a panic or a silent `200`
  (depends on T017, T018).

**Checkpoint**: All three user stories are independently functional.

---

## Phase 6: Polish & Cross-Cutting Concerns

- [X] T032 [P] Ran `quickstart.md` scenario 1 (Docker-only) live end-to-end against the real
  `kind-capi-mgmt` cluster — capabilities, cluster/machine `provider` field, infra REST
  auto-select/404, and both new WebSocket types all confirmed against real cluster data. See
  quickstart.md's "Live validation" note. Scenarios 2-5 (vSphere-only, mixed, unknown, none) are
  covered by the fake-client/component test suites — no additional live cluster of those shapes
  was available in this environment.
- [X] T033 [P] Updated route documentation/comments in
  `webserver/internal/web/handlers/handlers.go` for the new `/api/infra/capabilities` endpoint,
  the provider-dispatch behavior of the `/infra/list` routes, and the new WebSocket object types.
- [X] T034 Ran `make run-tests-backend`, `make run-tests-frontend`, and `make build` end-to-end —
  all green. `make build`'s production lint pass caught one unused import
  (`infra-capability-context.tsx`) that `tsc`/`jest` had missed; fixed and re-verified.

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — start immediately.
- **Foundational (Phase 2)**: Depends on Setup — BLOCKS all user stories.
- **User Stories (Phase 3-5)**: All depend on Foundational completion; independently testable and
  deliverable in priority order (P1 → P2 → P3), or in parallel if staffed.
- **Polish (Phase 6)**: Depends on all desired user stories being complete.

### User Story Dependencies

- **US1 (P1)**: No dependency on US2/US3. Delivers the core "right view for the right provider"
  behavior — the MVP.
- **US2 (P2)**: Independent of US1's infra-tab work; only needs the `Provider` field from
  Foundational (T010/T011) and the capabilities fetch (T012). Can be built in parallel with US1.
- **US3 (P3)**: Builds on the same Foundational primitives; its frontend empty-state task (T030)
  slots into the same `clusters/page.tsx` touched by US1 (T021) — sequence after US1 if one
  developer, otherwise coordinate on that file.

### Parallel Opportunities

- T001/T002 (Setup) run in parallel.
- T004, T008, T012 (Foundational, marked [P]) run in parallel with their non-parallel siblings once
  their own dependencies are met.
- US1 and US2 can be implemented in parallel by different developers once Foundational is done.
- Within each story, all [P]-marked tasks (typically model/fetcher/test files that don't collide)
  run in parallel.

---

## Parallel Example: User Story 1

```bash
# Backend model + fetcher, frontend component generalization — different files, run together:
Task: "Add models.ClusterInfraDocker struct to webserver/internal/infra/models/cluster.go"
Task: "Add fetchers.ListClusterInfraDocker in webserver/internal/infra/clusterapi/fetchers/cluster_infra_docker.go"
Task: "Generalize infra-lister.tsx and infra-table.tsx to accept a provider config"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup.
2. Complete Phase 2: Foundational (blocks everything).
3. Complete Phase 3: User Story 1.
4. **STOP and VALIDATE**: run quickstart.md scenarios 1-3 (docker-only, vsphere-only, mixed).
5. Deploy/demo if ready — this alone fixes the core "wrong/no provider view" problem.

### Incremental Delivery

1. Setup + Foundational → foundation ready.
2. US1 → validate independently → deploy (MVP).
3. US2 → validate independently → deploy (adds at-a-glance provider/version badges).
4. US3 → validate independently → deploy (adds graceful degradation for edge environments).

---

## Notes

- [P] tasks touch different files with no unmet dependencies.
- [Story] labels map every user-story-phase task back to spec.md for traceability.
- Commit after each task or logical group; stop at any checkpoint to validate a story independently.
- No task in this list requires a new `go.mod`/`package.json` dependency (per research.md R3).
