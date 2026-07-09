---

description: "Task list for feature implementation"
---

# Tasks: First-Class Pages for MachineHealthCheck, KubeadmControlPlane, MachineSet, and ClusterClass

**Input**: Design documents from `/specs/007-capi-object-pages/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md

**Tests**: Constitution Principle V (Test-Driven Quality) mandates test coverage for all backend
and frontend changes. Per research.md R1, the shared shell (`BaseLister`, `ObjectDetails`) is
already generically tested — new tests only need to cover each kind's thin per-kind wiring
(processor conversion, and the Specification tab's kind-specific fields), not re-test the shell.

**Organization**: Tasks are grouped by user story (US1–US4, priorities from spec.md). All four
stories are structurally independent — none blocks another — so there is no shared Foundational
phase; each story adds its own line to the few shared dispatch files (`websocket.go`,
`resource-gvr.ts`, `nav-links.tsx`), the same way every existing kind already does.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1, US2, US3, US4)

## Path Conventions

Web app: `webserver/` (Go backend), `front/` (Next.js frontend) — per plan.md Project Structure.
Every task mirrors the existing Machines page file-for-file (research.md R1).

---

## Phase 1: Setup

**Purpose**: Confirm the one new import resolves; no new dependency, no new generic infrastructure
needed (research.md R1, R3) — this phase is intentionally minimal.

- [X] T001 Add a throwaway `import _ "sigs.k8s.io/cluster-api/controlplane/kubeadm/api/v1beta1"` smoke-test in a scratch file, run `go build ./...` in `webserver/` to confirm the subpackage resolves from the already-vendored module (no go.mod version bump expected), then delete the scratch file

---

## Phase 2: User Story 1 - Inspect MachineHealthCheck remediation policy (Priority: P1) 🎯 MVP

**Goal**: A live "Machine Health Checks" list + detail page showing target selector, timeouts,
`maxUnhealthy`, remediation status, full YAML, and "Ask AI about this."

**Independent Test**: Navigate to the new page, confirm live list + detail matching spec.md US1's
Acceptance Scenarios.

### Tests for User Story 1

- [X] T002 [P] [US1] Go test for `ProcessMachineHealthCheck` (selector/timeouts/maxUnhealthy/status mapping) in `webserver/internal/infra/clusterapi/processor/machinehealthcheck_test.go`
- [X] T003 [P] [US1] Jest test for MachineHealthCheck `details.tsx`/`specification.tsx` rendering target selector, `maxUnhealthy`, timeouts, and remediation status (FR-006) in `front/app/ui/dashboard/components/machinehealthchecks/details.test.tsx`

### Implementation for User Story 1

- [X] T004 [US1] Define `models.MachineHealthCheck` in `webserver/internal/infra/models/machinehealthcheck.go` per data-model.md (embed `metav1.ObjectMeta`, `Age`, `Cluster`, `Selector`, `MaxUnhealthy`, `NodeStartupTimeout`, `UnhealthyConditions`, `Status`)
- [X] T005 [US1] Implement `processor.ProcessMachineHealthCheck` in `webserver/internal/infra/clusterapi/processor/machinehealthcheck.go` (depends on T004; T002 must fail first)
- [X] T006 [US1] Implement `watchers.WatchMachineHealthChecks` in `webserver/internal/web/watchers/machinehealthcheck.go`, mirroring `machine.go`'s `WatchMachines` (dynamic-client watch on `machinehealthchecks.cluster.x-k8s.io/v1beta1`, converter via `processor.ProcessMachineHealthCheck`) (depends on T005)
- [X] T007 [US1] Register a new `TypeMachineHealthCheck` `ObjectType` (`"machinehealthcheck"`) and `watchHandlers` entry in `webserver/internal/web/handlers/system/websocket.go` (depends on T006)
- [X] T008 [P] [US1] Add a `machineHealthCheck` entry to `RESOURCE_GVR` in `front/app/lib/resource-gvr.ts` (contracts/watch-types.md)
- [X] T009 [US1] Implement `front/app/ui/dashboard/components/machinehealthchecks/{lister.tsx,table.tsx,details.tsx,specification.tsx}`, mirroring the Machines page's file set, showing target selector/timeouts/`maxUnhealthy`/remediation status on Specification and wiring `useCurrentObjectContext`/`AskAIButton` (FR-006, FR-010) (depends on T003, T007, T008)
- [X] T010 [US1] Add `front/app/dashboard/machinehealthchecks/{layout.tsx,page.tsx}` (depends on T009)
- [X] T011 [US1] Add a "Machine Health Checks" entry to `front/app/ui/dashboard/nav-links.tsx` (FR-011) (depends on T010)

**Checkpoint**: User Story 1 is fully functional and independently testable/demoable.

---

## Phase 3: User Story 2 - Inspect KubeadmControlPlane / etcd health (Priority: P2)

**Goal**: A live "Kubeadm Control Planes" list + detail page showing replica counts, status
conditions (including etcd-related conditions when present), full YAML, and "Ask AI about this."

**Independent Test**: Navigate to the new page, confirm live list + detail matching spec.md US2's
Acceptance Scenarios, including the empty-state case when KCP isn't in use.

### Tests for User Story 2

- [X] T012 [P] [US2] Go test for `ProcessKubeadmControlPlane` (replica/status/condition mapping) in `webserver/internal/infra/clusterapi/processor/kubeadmcontrolplane_test.go`
- [X] T013 [P] [US2] Jest test for KubeadmControlPlane `details.tsx`/`specification.tsx` rendering desired/ready replicas and conditions, including an etcd-condition case (FR-007) in `front/app/ui/dashboard/components/kubeadmcontrolplanes/details.test.tsx`

### Implementation for User Story 2

- [X] T014 [US2] Define `models.KubeadmControlPlane` in `webserver/internal/infra/models/kubeadmcontrolplane.go` per data-model.md, importing `controlplanev1 "sigs.k8s.io/cluster-api/controlplane/kubeadm/api/v1beta1"` (research.md R3) — **deviation**: importing this subpackage transitively pulls in `bootstrap/kubeadm/api/v1beta1` and `feature`, which needed `go mod tidy` (3 new indirect go.sum entries: `k8s.io/cluster-bootstrap`, `k8s.io/component-base/featuregate` deps); research.md R3's "no go.mod version bump expected" held (no version changes), but go.sum did need new entries
- [X] T015 [US2] Implement `processor.ProcessKubeadmControlPlane` in `webserver/internal/infra/clusterapi/processor/kubeadmcontrolplane.go` (depends on T014; T012 must fail first)
- [X] T016 [US2] Implement `watchers.WatchKubeadmControlPlanes` in `webserver/internal/web/watchers/kubeadmcontrolplane.go`, watching GVR `controlplane.cluster.x-k8s.io/v1beta1, Resource=kubeadmcontrolplanes` (research.md R3), converter via `processor.ProcessKubeadmControlPlane` (depends on T015)
- [X] T017 [US2] Register a new `TypeKubeadmControlPlane` `ObjectType` (`"kubeadmcontrolplane"`) and `watchHandlers` entry in `webserver/internal/web/handlers/system/websocket.go` (depends on T016)
- [X] T018 [P] [US2] Add a `kubeadmControlPlane` entry to `RESOURCE_GVR` in `front/app/lib/resource-gvr.ts`
- [X] T019 [US2] Implement `front/app/ui/dashboard/components/kubeadmcontrolplanes/{lister.tsx,table.tsx,details.tsx,specification.tsx}`, showing desired/ready replicas and status conditions on Specification and wiring `useCurrentObjectContext`/`AskAIButton` (FR-007, FR-010) (depends on T013, T017, T018)
- [X] T020 [US2] Add `front/app/dashboard/kubeadmcontrolplanes/{layout.tsx,page.tsx}` (depends on T019)
- [X] T021 [US2] Add a "Kubeadm Control Planes" entry to `front/app/ui/dashboard/nav-links.tsx` (FR-011) (depends on T020)

**Checkpoint**: User Story 2 is fully functional independently of US1/US3/US4.

---

## Phase 4: User Story 3 - Inspect MachineSet rollout state (Priority: P3)

**Goal**: A live "Machine Sets" list + detail page showing replica counts, owning
MachineDeployment, status conditions, full YAML, and "Ask AI about this."

**Independent Test**: Navigate to the new page, confirm replica counts match what the Day-2 Ops
dashboard's stalled-rollout warning already reports for the same MachineSet.

### Tests for User Story 3

- [X] T022 [P] [US3] Go test for `ProcessMachineSet` (replica/owning-MachineDeployment/status mapping) in `webserver/internal/infra/clusterapi/processor/machineset_test.go`
- [X] T023 [P] [US3] Jest test for MachineSet `details.tsx`/`specification.tsx` rendering replica counts and owning MachineDeployment (FR-008) in `front/app/ui/dashboard/components/machinesets/details.test.tsx`

### Implementation for User Story 3

- [X] T024 [US3] Define `models.MachineSet` in `webserver/internal/infra/models/machineset.go` per data-model.md (`MachineDeployment` field from the `cluster.x-k8s.io/deployment-name` label, same convention as 006's `machineSetsFor`)
- [X] T025 [US3] Implement `processor.ProcessMachineSet` in `webserver/internal/infra/clusterapi/processor/machineset.go` (depends on T024; T022 must fail first)
- [X] T026 [US3] Implement `watchers.WatchMachineSets` in `webserver/internal/web/watchers/machineset.go`, watching `machinesets.cluster.x-k8s.io/v1beta1` (depends on T025) — reused the existing `machineSetGVR` package var already declared in `day2ops.go` instead of redeclaring it
- [X] T027 [US3] Register a new `TypeMachineSet` `ObjectType` (`"machineset"`) and `watchHandlers` entry in `webserver/internal/web/handlers/system/websocket.go` (depends on T026)
- [X] T028 [P] [US3] Add a `machineSet` entry to `RESOURCE_GVR` in `front/app/lib/resource-gvr.ts`
- [X] T029 [US3] Implement `front/app/ui/dashboard/components/machinesets/{lister.tsx,table.tsx,details.tsx,specification.tsx}`, showing replicas/ready/available counts, owning MachineDeployment, and conditions on Specification, wiring `useCurrentObjectContext`/`AskAIButton` (FR-008, FR-010) (depends on T023, T027, T028)
- [X] T030 [US3] Add `front/app/dashboard/machinesets/{layout.tsx,page.tsx}` (depends on T029)
- [X] T031 [US3] Add a "Machine Sets" entry to `front/app/ui/dashboard/nav-links.tsx` (FR-011) (depends on T030)

**Checkpoint**: User Story 3 is fully functional independently of US1/US2/US4.

---

## Phase 5: User Story 4 - Browse ClusterClass as a first-class page (Priority: P4)

**Goal**: A live "Cluster Classes" list + detail page reusing the existing `models.ClusterClass`/
`processor.ProcessClusterClass`, alongside (not replacing) the existing main-dashboard widget.

**Independent Test**: Navigate to the new page, confirm live list + detail, and confirm the existing
main-dashboard `ClusterClassLister` widget is unchanged (research.md R5).

### Tests for User Story 4

- [X] T032 [P] [US4] Jest test for ClusterClass `details.tsx`/`specification.tsx` rendering status/reference fields (FR-009) in `front/app/ui/dashboard/components/clusterclasses/details.test.tsx`

### Implementation for User Story 4

- [X] T033 [US4] Implement `watchers.WatchClusterClasses` in `webserver/internal/web/watchers/clusterclass.go`, watching `clusterclasses.cluster.x-k8s.io/v1beta1`, converter via the *existing* `processor.ProcessClusterClass` (research.md R5 — no new model/processor needed) — **deviation**: `models.ClusterClass` has no embedded `metav1.ObjectMeta` (flat `Name`/`Namespace` fields only, matching the pre-existing main-dashboard widget), but `BaseLister`/`ObjectTable` require `metadata.name` for row keys and search. Fixed by wrapping the unchanged processor output in a small `clusterClassWithMeta` struct at the watcher boundary that adds a synthesized `metadata: {name, namespace}` mirror — the model/processor themselves stay untouched per R5
- [X] T034 [US4] Register a new `TypeClusterClass` `ObjectType` (`"clusterclass"`) and `watchHandlers` entry in `webserver/internal/web/handlers/system/websocket.go` (depends on T033)
- [X] T035 [P] [US4] Add a `clusterClass` entry to `RESOURCE_GVR` in `front/app/lib/resource-gvr.ts`
- [X] T036 [US4] Implement `front/app/ui/dashboard/components/clusterclasses/{lister.tsx,table.tsx,details.tsx,specification.tsx}`, wiring `useCurrentObjectContext`/`AskAIButton` (FR-009, FR-010) (depends on T032, T034, T035)
- [X] T037 [US4] Add `front/app/dashboard/clusterclasses/{layout.tsx,page.tsx}` (depends on T036)
- [X] T038 [US4] Add a "Cluster Classes" entry to `front/app/ui/dashboard/nav-links.tsx` (FR-011) (depends on T037)
- [X] T039 [US4] Verify the existing `front/app/ui/dashboard/components/dashboard/clusterclass.tsx` widget on the main dashboard is unmodified and still renders (spec.md US4 Acceptance Scenario 3) — verified via `dashboard.test.tsx` passing unmodified and a `git status` check showing the widget file untouched

**Checkpoint**: All four user stories are independently functional.

---

## Phase 6: Polish & Cross-Cutting Concerns

- [X] T040 [P] Run quickstart.md validation scenarios end-to-end against a `kind-capi-mgmt`-style test cluster — deferred to manual live verification (no test cluster available in this session); all four pages were exercised via unit/component tests and a full production build instead (T041)
- [X] T041 [P] Verify `make build`, `make run-tests-backend`, and `make run-tests-frontend` all pass — all three green: backend 12 packages ok, frontend 40 suites/169 tests passed, `make build` produced a static export with all four new routes (`/dashboard/machinehealthchecks`, `/dashboard/kubeadmcontrolplanes`, `/dashboard/machinesets`, `/dashboard/clusterclasses`) and a compiling backend binary
- [X] T042 Annotate this tasks.md with any deviations discovered mid-implementation (project convention from features 004/005/006)

### Discovered mid-implementation

- **T014 — go.mod/go.sum needed new entries despite research.md R3**: importing
  `controlplanev1 "sigs.k8s.io/cluster-api/controlplane/kubeadm/api/v1beta1"` transitively pulls in
  `bootstrap/kubeadm/api/v1beta1` and `feature`, which reference `k8s.io/cluster-bootstrap` and
  `k8s.io/component-base/featuregate` — not previously in `go.sum`. R3's "no version bump expected"
  held (no module *versions* changed), but `go mod tidy` was required to add the missing indirect
  go.sum entries before the package would build.
- **T033 — `models.ClusterClass` lacks `metadata.name` needed by the shared list/detail shell**:
  the existing model (reused unchanged per R5) has flat `Name`/`Namespace` fields, not an embedded
  `metav1.ObjectMeta`, because it was designed only for the main-dashboard widget. `BaseLister`/
  `ObjectTable` require `metadata.name` for row keys and search filtering (the same convention every
  other kind follows). Rather than changing the shared model (which would also touch the untouched
  widget, violating research.md R5's isolation goal), `watchers.WatchClusterClasses` wraps the
  unchanged `processor.ProcessClusterClass` output in a small `clusterClassWithMeta` struct that adds
  a synthesized `metadata: {name, namespace}` mirror at the JSON-serialization boundary only.
- **Every US1–US3 watcher reused an existing package-level GVR var where one already existed**:
  `machineHealthCheckGVR` and `machineSetGVR` were already declared in `day2ops.go` (006) for
  internal severity-classification use; the new standalone watchers reuse them directly rather than
  redeclaring identical GVR values under new names. `kubeadmControlPlaneGVR` was net-new (KCP was
  never watched before 007) and `clusterClassGVR` was net-new (ClusterClass was previously only
  fetched via the one-shot REST path, not watched).

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — can start immediately
- **User Stories (Phase 2–5)**: Each depends only on Setup; all four are independently
  implementable/testable in any order (no shared Foundational phase — research.md R1)
- **Polish (Phase 6)**: Depends on all desired user stories being complete

### User Story Dependencies

- **US1, US2, US3, US4**: Fully independent of each other. Each touches its own new files plus one
  additive line each in three shared files (`websocket.go`, `resource-gvr.ts`, `nav-links.tsx`) —
  the same low-conflict pattern every existing kind already follows.

### Parallel Opportunities

- All four user-story phases (US1–US4) can be implemented in parallel by different developers once
  Setup completes — none blocks another.
- Within each story, the `[P]`-marked test tasks (processor Go test + Jest test) run in parallel,
  and the `RESOURCE_GVR` task runs in parallel with the backend watcher-registration tasks.

---

## Parallel Example: User Story 1

```bash
# Tests (parallel):
Task: "Go test for ProcessMachineHealthCheck in webserver/internal/infra/clusterapi/processor/machinehealthcheck_test.go"
Task: "Jest test for MachineHealthCheck details.tsx/specification.tsx in front/app/ui/dashboard/components/machinehealthchecks/details.test.tsx"

# RESOURCE_GVR entry in parallel with backend watcher work:
Task: "Add machineHealthCheck entry to RESOURCE_GVR in front/app/lib/resource-gvr.ts"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: User Story 1 (MachineHealthCheck)
3. **STOP and VALIDATE**: confirm against quickstart.md's US1 scenarios
4. Deploy/demo if ready — this alone gives operators direct visibility into the policy behind every
   006/US4 self-healing classification

### Incremental Delivery

1. Setup → ready
2. US1 (MachineHealthCheck) → MVP
3. US2 (KubeadmControlPlane) → closes the biggest blind spot (etcd/control-plane health)
4. US3 (MachineSet) → rollout-state detail
5. US4 (ClusterClass) → consistency with the other three first-class kinds
6. Each story is fully additive — no story's completion depends on, or breaks, another's

### Parallel Team Strategy

With multiple developers, after Setup completes: one developer per user story (US1, US2, US3, US4)
— since each story is a self-contained clone of the same proven pattern, this is the ideal feature
for full parallelization across a team.
