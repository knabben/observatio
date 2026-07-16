---

description: "Task list for feature implementation"
---

# Tasks: Velero Backup Recoverability Awareness

**Input**: Design documents from `/specs/008-velero-backup-recoverability/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md

**Tests**: Constitution Principle V (Test-Driven Quality) mandates test coverage for all backend
and frontend changes. Per research.md, all new logic lives in small, pure, already-tested-pattern
`day2ops` functions (mirroring `severity.go`/`risk_certexpiry.go`), each getting its own test file.

**Organization**: Tasks are grouped by user story (US1–US3, priorities from spec.md). Unlike
feature 007 (where each kind was fully independent), all three stories here extend the *same*
existing Day-2 Ops aggregator (`day2ops.go`/`assembleData`) with the *same* underlying Velero
watch data (research.md R2) — so watching the four Velero GVRs and storing their decoded state is
a genuine Foundational dependency shared by every story, not story-specific work.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1, US2, US3)

## Path Conventions

Web app: `webserver/` (Go backend), `front/` (Next.js frontend) — per plan.md Project Structure.
This feature extends the existing Day-2 Ops vertical slice; it adds no new frontend routes and no
new WS `ObjectType`.

---

## Phase 1: Setup

**Purpose**: Confirm research.md R1 holds (no new Go dependency) before touching the aggregator.

- [X] T001 Run `go build ./...` in `webserver/` to record the current baseline (zero Velero
      dependency in `go.mod`/`go.sum`) before any Velero-decoding code is added, confirming
      research.md R1's unstructured-decode approach is the only change needed to support it

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Watch the four Velero GVRs and hold their decoded state in the existing
`day2opsStore` — every user story below reads from this same data (research.md R2, R8).

**⚠️ CRITICAL**: No user story work can begin until this phase is complete.

- [X] T002 [P] Add `RecoveryInfo`, `BackupStorageLocationStatus`, `ClusterBackupCoverage`, and
      `BackupHealth` types to `webserver/internal/infra/clusterapi/day2ops/types.go` per
      data-model.md; add `RecoveryInfo *RecoveryInfo` to the existing `FailureSeverity` struct and
      `BackupHealth BackupHealth` to the existing `Data` struct (both additive, backward compatible)
- [X] T003 Add `backupGVR`, `restoreGVR`, `scheduleGVR`, `backupStorageLocationGVR` package-level
      vars (group `velero.io`, version `v1`) to `webserver/internal/web/watchers/day2ops.go`
- [X] T004 Add a Velero-installed detection check (existence of the `backups.velero.io` CRD via
      the `apiextClient` already constructed in `WatchDay2Ops`) and gate inclusion of the four
      Velero GVRs in `day2opsWatchedGVRs` on it, in `webserver/internal/web/watchers/day2ops.go`
      (research.md R8; depends on T003)
- [X] T005 Add `day2opsStore` fields for `backups`, `restores`, `schedules`,
      `backupStorageLocations` (decoded lightweight structs per data-model.md, extracted via
      `unstructured.Nested*` — research.md R1, not raw `unstructured.Unstructured`), `apply()`
      cases for the four new GVRs (add/modify/delete), and snapshot accessors, in
      `webserver/internal/web/watchers/day2ops.go` (depends on T003) — **deviation**: the
      decode-from-unstructured functions (`ExtractBackupInfo` etc.) live in the `day2ops` package
      itself, not `watchers`, mirroring the existing `day2ops.ExtractProviderResourceStatus`
      precedent; this also avoids an import-direction problem, since `ComputeBackupHealth` (US1)
      needs to consume these same decoded types and `day2ops` cannot import `watchers` back
- [X] T006 [P] Go test for the store's decode-and-apply logic for the four new Velero kinds
      (add/modify/delete, and a case confirming an unparseable/partial object doesn't crash the
      store) in `webserver/internal/web/watchers/day2ops_test.go` (depends on T005)

**Checkpoint**: Foundation ready — US1, US2, and US3 can all now read Velero state from the store.

---

## Phase 3: User Story 1 - Backup Health at a glance (Priority: P1) 🎯 MVP

**Goal**: A Backup Health card on the Day-2 Ops landing page showing storage-location
reachability and per-cluster backup staleness against a configurable RPO.

**Independent Test**: With Velero installed and at least one Backup and BackupStorageLocation
present, open the Day-2 Ops landing page and confirm the card appears with correct reachability
and staleness — without touching CA-loss severity or restore-activity logic at all.

### Tests for User Story 1

- [X] T007 [P] [US1] Go test for `ComputeBackupHealth` — storage-location reachability, per-cluster
      on-time/stale/no-coverage classification against an RPO threshold, and the
      `Available: false` case when Velero isn't installed — in
      `webserver/internal/infra/clusterapi/day2ops/backuphealth_test.go`
- [X] T008 [P] [US1] Go test for `ComputeClusterBackupCoverage`'s namespace/label-selector matching
      heuristic (research.md R5) — namespace match, label match, no match, `PartiallyFailed`
      backup not counted as covering — in
      `webserver/internal/infra/clusterapi/day2ops/backuphealth_test.go`
- [X] T009 [P] [US1] Jest test for `BackupHealthCard` rendering reachable/unreachable storage
      locations, on-time/stale/no-coverage clusters, and the "not available" state, in
      `front/app/ui/dashboard/components/ops/backup-health-card.test.tsx` — confirmed: 6/6 tests
      pass (checkbox was left unmarked despite the file already being complete)

### Implementation for User Story 1

- [X] T010 [US1] Implement `ComputeClusterBackupCoverage` in
      `webserver/internal/infra/clusterapi/day2ops/backuphealth.go` (depends on T002, T008 failing
      first) — **deviation**: built with its full US3 restore-outcome/in-progress fields already
      included (see T021 note) rather than a US1-only subset, since splitting it into two edits
      would have meant rewriting the same function; T008's tests cover the US1-relevant behavior
      only, T020 covers the restore-specific behavior
- [X] T011 [US1] Implement `ComputeBackupHealth` (storage-location reachability + calls
      `ComputeClusterBackupCoverage` per known cluster + the configurable RPO threshold default) in
      `webserver/internal/infra/clusterapi/day2ops/backuphealth.go` (depends on T010, T007 failing
      first)
- [X] T012 [US1] Wire `ComputeBackupHealth` into `assembleData` in
      `webserver/internal/web/watchers/day2ops.go`, populating `Data.BackupHealth` on every
      recompute (depends on T005, T011) — `assembleData` gained a `veleroInstalled bool` parameter,
      checked once per connection at `WatchDay2Ops` setup rather than on every recompute
- [X] T013 [P] [US1] Add `BackupHealth`, `ClusterBackupCoverage`, `BackupStorageLocationStatus`
      frontend types and extend the existing `Data` type in
      `front/app/ui/dashboard/shared/use-day2-ops.ts` per data-model.md — also updated two
      pre-existing fixtures (`ops-dashboard.test.tsx`, `use-day2-ops.test.tsx`) that construct a
      full `Day2OpsData` literal, since the new field is required, not optional
- [X] T014 [US1] Implement `BackupHealthCard` in
      `front/app/ui/dashboard/components/ops/backup-health-card.tsx` (depends on T009, T013)
- [X] T015 [US1] Render `BackupHealthCard` on the Day-2 Ops landing page in
      `front/app/ui/dashboard/components/ops/ops-dashboard.tsx`, alongside the existing rollup
      cards (depends on T014) — only shown when `filter === 'all'` (no `backup` entry was added to
      the category-filter buttons, since Backup Health isn't `Category`-based)

**Checkpoint**: User Story 1 is fully functional and independently testable/demoable.

---

## Phase 4: User Story 2 - Know whether CA loss is actually recoverable (Priority: P1)

**Goal**: A cluster's CA-secret-missing severity indicates whether a covering backup exists and
its age, distinguishing "recoverable" from "unrecoverable" in the same severity banner.

**Independent Test**: Trigger CA-secret-missing severity for a cluster with a recent covering
backup, and separately for one with none; confirm the two cases are distinguishable from the
severity alone.

### Tests for User Story 2

- [X] T016 [P] [US2] Go test for `ComputeCASecretMissingSeverity`'s new recovery-info behavior:
      recoverable-with-age, no-covering-backup, and the case where coverage data isn't available
      (Velero not installed) omits `RecoveryInfo` entirely rather than reporting `false`, in
      `webserver/internal/infra/clusterapi/day2ops/severity_test.go`

### Implementation for User Story 2

- [X] T017 [US2] Extend `ComputeCASecretMissingSeverity`'s signature to accept a
      `*ClusterBackupCoverage` and populate `FailureSeverity.RecoveryInfo` + enrich `Reason` with
      the covering-backup age (or its absence), in
      `webserver/internal/infra/clusterapi/day2ops/severity.go` (depends on T002, T016 failing
      first)
- [X] T018 [US2] Update `clusterCertRisksAndSeverities` in
      `webserver/internal/web/watchers/day2ops.go` to look up each cluster's
      `ClusterBackupCoverage` (from the same computation T011 already performs) and pass it to
      `ComputeCASecretMissingSeverity` (depends on T012, T017) — **deviation**: required reordering
      `assembleData` so `ComputeBackupHealth` runs before `clusterCertRisksAndSeverities` (previously
      the first thing computed), and threading a `veleroInstalled bool` + a
      `map[string]day2ops.ClusterBackupCoverage` (keyed by `namespace/name`, reusing the existing
      `objectKey` helper) through both functions' signatures
- [X] T019 [P] [US2] Jest regression test confirming `SeverityBanner` renders the enriched
      recoverable/unrecoverable `reason` prose correctly (no frontend code change expected per
      research.md R4 — this guards against the enrichment silently breaking existing rendering) in
      `front/app/ui/dashboard/components/ops/severity-banner.test.tsx` — confirmed: 0 lines changed
      in `severity-banner.tsx` itself, only new test cases added

**Checkpoint**: User Stories 1 AND 2 both work independently; CA-loss severities now carry
recoverability.

---

## Phase 5: User Story 3 - See recovery activity and reconciliation-pause state (Priority: P2)

**Goal**: Restore-in-progress/outcome is visible per cluster and in aggregate; reconciliation-pause
state is visible (confirmed already implemented — research.md R6).

**Independent Test**: With a Restore in progress or recently completed, confirm its state is
reflected on the Day-2 Ops landing page; confirm `spec.paused` is visible on the Cluster detail
page (pre-existing).

### Tests for User Story 3

- [X] T020 [P] [US3] Go test for restore-derived fields — `RestoresInProgress` aggregate count on
      `BackupHealth`, and `RestoreInProgress`/`LastRestoreOutcome` per `ClusterBackupCoverage`
      (in-progress, succeeded, failed cases) — in
      `webserver/internal/infra/clusterapi/day2ops/backuphealth_test.go` — **deviation**: already
      written as part of T008 (`Test_ComputeClusterBackupCoverage_restoreInProgress`,
      `_restoreOutcomes`, `_restoreForUnrelatedBackupIgnored`) and T007
      (`Test_ComputeBackupHealth_perClusterCoverageAndAggregateRestoresInProgress`), per T010's
      noted deviation of building the complete function upfront rather than in two passes

### Implementation for User Story 3

- [X] T021 [US3] Extend `ComputeBackupHealth`/`ComputeClusterBackupCoverage` to populate
      `RestoresInProgress`, `RestoreInProgress`, and `LastRestoreOutcome` from the store's restore
      data in `webserver/internal/infra/clusterapi/day2ops/backuphealth.go` (depends on T020
      failing first, T010, T011) — already implemented as part of T010/T011 (see deviation note
      there); no additional code needed here
- [X] T022 [P] [US3] Extend `BackupHealthCard` to render the restores-in-progress count, with a
      Jest test case added to `backup-health-card.test.tsx`, in
      `front/app/ui/dashboard/components/ops/backup-health-card.tsx` (depends on T021, T014) —
      already implemented as part of T014/T009 ("shows the restores-in-progress count when greater
      than zero" test); no additional code needed here
- [X] T023 [US3] Confirm `Cluster.spec.paused` visibility on the existing Cluster detail page
      (`front/app/ui/dashboard/components/clusters/specification.tsx`) needs no changes
      (research.md R6) — verification-only, no code; check off after re-running
      `clusters/table.test.tsx`/existing Cluster tests and confirming they still pass unmodified —
      confirmed: all 13 tests in `components/clusters` + `details-tabs.test.tsx` pass unmodified

**Checkpoint**: All three user stories are independently functional.

---

## Phase 6: Polish & Cross-Cutting Concerns

- [X] T024 [P] Run quickstart.md validation scenarios end-to-end against a Velero-enabled
      `kind-capi-mgmt`-style test cluster — deferred to manual live verification (no Velero
      installation available in this session); all logic was exercised via unit/component tests
      instead (T007-T023)
- [X] T025 [P] Verify `make build`, `make run-tests-backend`, and `make run-tests-frontend` all
      pass — **`make run-tests-backend` and `make run-tests-frontend` both green** (12 Go packages
      ok, 33 Jest suites / 150 tests passed). **`make build` could not be verified clean in this
      session**: this repo's `front/next.config.ts` sets `distDir: './output'`, and a `next dev`
      process was already running in another terminal (not started by this session, `pnpm run dev`,
      started independently by the user) sharing that same directory — every `make build` attempt
      raced against it and picked up stale, dev-server-regenerated route-type declarations for
      pages that don't exist on this branch (`kubeadmcontrolplanes`, etc., left over from a
      different branch state). Cleared `.next`/`output` and retried once; the race reproduced
      identically. Did not stop the user's dev server to force a clean run. As independent
      evidence the code itself is sound: `go build ./...` is clean, and a scoped `tsc --noEmit`
      check confirmed zero errors in every file this feature touched (only the stale/unrelated
      route files under `output/types` errored). Recommend re-running `make build` once the
      conflicting dev server is stopped
- [X] T026 Annotate this tasks.md with any deviations discovered mid-implementation (project
      convention from features 004–007)

---

### Discovered mid-implementation

- **T010 — `ComputeClusterBackupCoverage` built complete (US1+US3) in one pass**: rather than
  implementing US1's on-time/stale/coverage fields first and revisiting the same function for
  US3's restore fields later, both were built together, since splitting them would have meant
  rewriting the same switch/sort logic twice. T007/T008's tests cover the US1-relevant behavior;
  T020's restore-specific tests already existed by the time US3's phase began, so T020-T022 needed
  no additional code — only verification that the coverage was already there (see those tasks'
  notes).
- **T018 — `assembleData` reordering**: wiring US2's CA-loss cross-reference required computing
  `BackupHealth` (and thus each cluster's coverage) *before* `clusterCertRisksAndSeverities` ran,
  reversing their original order. `clusterCertRisksAndSeverities` and `assembleData` both gained
  new parameters (`veleroInstalled bool`, a coverage lookup map) as a result.
- **Research finding confirmed correct (R6)**: `Cluster.spec.paused` visibility was already fully
  implemented by prior work (feature predates 008) — no code was needed for that half of US3,
  only a confirmation that its tests still pass.
- **`make build` could not be independently verified in this session** due to a live `next dev`
  process (started by the user outside this session) sharing this repo's `distDir: './output'`
  with `next build`, causing every build attempt to race against the dev server's own regeneration
  of that directory and pick up stale type declarations for routes that don't exist on this
  branch. `go build`, `go test ./...` (both modules), and a scoped `tsc --noEmit` check (zero
  errors in any file this feature touched) were used as substitute verification. See T025.

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — can start immediately
- **Foundational (Phase 2)**: Depends on Setup — BLOCKS all user stories (unlike 007, this
  feature's stories share one underlying data source, so this gate is real, not a formality)
- **User Stories (Phase 3–5)**: All depend on Foundational (Phase 2) completion
  - US1 and US2 are both P1 and independent of each other once Foundational is done
  - US3 reuses US1's `ComputeBackupHealth`/`ComputeClusterBackupCoverage` (extends the same
    functions with restore fields) — sequenced after US1 for that reason, though its Cluster-pause
    confirmation (T023) has no dependency at all
- **Polish (Phase 6)**: Depends on all desired user stories being complete

### User Story Dependencies

- **US1 (P1)**: Depends only on Foundational.
- **US2 (P1)**: Depends only on Foundational for its own severity/coverage logic (T016-T017); T018
  additionally depends on US1's T012 (the `assembleData` wiring) since it hooks into the same call.
- **US3 (P2)**: T021-T022 extend US1's `ComputeBackupHealth`/`BackupHealthCard` (T010-T014)
  directly, so US3 is sequenced after US1 in practice even though nothing in spec.md makes it
  strictly blocked; T023 has no dependency on anything in this feature.

### Parallel Opportunities

- T002 (types) can run in parallel with T001 (Setup).
- Within Foundational, T006's test can be written in parallel with T003/T004 (different files)
  but depends on T005 to pass.
- Within each user story, `[P]`-marked test tasks run in parallel with each other.
- US1 and US2's test-writing (T007/T008 and T016) can happen in parallel — they touch different
  files (`backuphealth_test.go` vs `severity_test.go`).

---

## Parallel Example: User Story 1

```bash
# Tests (parallel):
Task: "Go test for ComputeBackupHealth in webserver/internal/infra/clusterapi/day2ops/backuphealth_test.go"
Task: "Go test for ComputeClusterBackupCoverage matching heuristic in webserver/internal/infra/clusterapi/day2ops/backuphealth_test.go"
Task: "Jest test for BackupHealthCard in front/app/ui/dashboard/components/ops/backup-health-card.test.tsx"

# Frontend types in parallel with backend implementation:
Task: "Add BackupHealth frontend types in front/app/ui/dashboard/shared/use-day2-ops.ts"
```

---

## Implementation Strategy

### MVP First (Setup + Foundational + User Story 1)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL — watches the four Velero GVRs, nothing else works
   without this)
3. Complete Phase 3: User Story 1
4. **STOP and VALIDATE**: confirm against quickstart.md's US1 scenarios — this alone gives
   operators the visibility the feature description says is completely missing today
5. Deploy/demo if ready

### Incremental Delivery

1. Setup + Foundational → ready
2. US1 (Backup Health rollup) → MVP
3. US2 (CA-loss recoverability cross-reference) → the feature's flagship value, closes the exact
   gap named in the feature description
4. US3 (restore activity + pause confirmation) → completes the recovery-procedure loop
5. Each story is additive; US3 extends US1's structs but never breaks US1's existing behavior

### Parallel Team Strategy

Foundational must land first (it's a genuine shared blocker here, unlike 007). After that, one
developer on US1, one on US2 — both read the same store, write to different files
(`backuphealth.go` vs `severity.go`), and integrate at T018's single line in `day2ops.go`. US3
should follow US1 since it directly extends the same functions.
