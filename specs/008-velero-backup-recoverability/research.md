# Research: Velero Backup Recoverability Awareness

## R1: Decode Velero objects via the dynamic client — no new typed Go dependency

**Decision**: Watch `Backup`, `Restore`, `Schedule`, and `BackupStorageLocation` (group
`velero.io`, version `v1`) via the existing `k8s.io/client-go` dynamic client and
`unstructured.Nested*` accessors, extracting only the handful of fields this feature needs.

**Rationale**: `github.com/vmware-tanzu/velero` is not currently a dependency anywhere in the repo
(confirmed: absent from `webserver/go.mod`, `go.sum`, and the local module cache). Importing its
typed API package would pull in a real external module purely to read `phase`, a few timestamps,
and a label selector — the exact situation the codebase already has a precedent for rejecting:
`webserver/internal/infra/clusterapi/fetchers/machine_infra_docker.go`'s `ProcessDockerMachine`
decodes `DockerMachine` via `unstructured.Nested*` specifically because "no typed Go package for it
is imported." The same reasoning applies here, more strongly, since Velero's types are further
from this project's actual domain (CAPI) than CAPI's own Docker provider is.

**Alternatives considered**:
- Add `github.com/vmware-tanzu/velero` as a dependency and use its typed `velerov1.Backup` etc. —
  rejected: unnecessary weight for four field-reads, and would need a Complexity Tracking
  justification the constitution's "no new runtime dependency without justification" gate doesn't
  clearly warrant here.
- Vendor just the CRD type definitions by hand-copying Velero's struct definitions into this repo
  — rejected: duplicates upstream API surface we'd have to keep in sync by hand; unstructured
  decoding of a handful of fields is simpler and has zero drift risk (it only ever reads, never
  assumes a field it doesn't decode).

## R2: Extend the existing Day-2 Ops aggregator, not a new pipeline

**Decision**: Add the four Velero GVRs into `webserver/internal/web/watchers/day2ops.go`'s existing
fan-in watch (`day2opsWatchedGVRs`, `day2opsStore`, `assembleData`), the same way `MachineSet` and
`MachineHealthCheck` were added in 006 — not as standalone watchers exposed via a new WS
`ObjectType`, and not as a second parallel aggregator.

**Rationale**: Every signal this feature needs (backup staleness, BSL reachability, cluster
coverage, CA-loss cross-reference) exists purely to enrich the Day-2 Ops dashboard's existing
rollups and severities — there is no requirement in spec.md for a standalone "Backups" list page
the way feature 007 built first-class pages for MachineHealthCheck/KubeadmControlPlane/MachineSet/
ClusterClass. 006's own tasks.md documents the identical precedent: MachineSet and
MachineHealthCheck GVRs/store-tracking were added directly inside `day2ops.go` "rather than
creating unused standalone watcher functions," specifically because they're internal-only signals.
Velero's four kinds are the same shape of dependency here. Extending the existing aggregator also
means the CA-loss cross-reference (US2) is computed in the same `assembleData` call that already
computes `ComputeCASecretMissingSeverity`, avoiding a second cross-cutting read of cluster cert
Secrets or a fragile cross-goroutine handoff between two independent watchers.

**Alternatives considered**:
- A second fan-in aggregator dedicated to Velero, feeding its own WS `ObjectType` — rejected: would
  require the frontend to open and reconcile two independent live streams to answer one question
  ("is this CA-lost cluster recoverable?"), and duplicates `day2ops.go`'s fan-in/store/apply
  machinery for no benefit.
- First-class Backup/Restore/Schedule/BackupStorageLocation list+detail pages (the 007 pattern) —
  rejected for this feature: not requested by spec.md, and the primary value (US1/US2) is the
  cross-referenced rollup, not independent browsing of raw Velero objects. Nothing here precludes
  adding first-class pages in a later feature if operators want to browse raw Backups directly —
  this plan does not foreclose that option, it just doesn't build it now.

## R3: Backup Health is a new dedicated type, not squeezed into the existing 3-bucket rollup

**Decision**: Introduce `day2ops.BackupHealth` as its own struct (storage-location reachability
list + per-cluster coverage list + restores-in-progress count + an `Available bool` for the
Velero-not-installed case), rendered by a new `BackupHealthCard` frontend component placed
alongside the existing `HealthRollupCard`s in the same rollup row — not shoehorned into the
existing `Category`/`HealthRollup{Healthy,Degraded,Failed}` shape.

**Rationale**: The existing `HealthRollup` type is a clean fit for "how many objects of this kind
are healthy/degraded/failed" (Cluster, MachineDeployment, Machine). Backup Health's actual content
— which storage locations are reachable, which clusters are stale vs. covered vs. never backed up,
whether recoverability is currently in question — doesn't reduce to three counts without losing the
information spec.md's acceptance scenarios need to be checkable (e.g. "no covering backup exists"
must remain distinguishable from "stale," not collapsed into one "degraded" bucket). spec.md's
FR-001 asks for the summary to appear "alongside" the existing rollups, not to *be* one — the plan
satisfies that literally: same landing page, same visual row, its own tailored shape underneath.

**Alternatives considered**: Add a fourth `Category` (`"backup"`) and force-fit healthy=on-time,
degraded=stale, failed=no-coverage into the existing `HealthRollupCard` — rejected: BSL
reachability and per-cluster ages have no home in that shape, and a card that changes meaning per
category (three generic buckets meaning different things for backups vs. objects) is worse UX than
a purpose-built card, for one extra frontend file.

## R4: CA-loss recoverability is a structured, backward-compatible addition to `FailureSeverity`

**Decision**: Add an optional field to the existing `day2ops.FailureSeverity` struct —
`RecoveryInfo *RecoveryInfo` (nil for every severity kind except CA-secret-missing) — carrying
`{Recoverable bool; CoveringBackupAge string}`, and enrich `ComputeCASecretMissingSeverity`'s
`Reason` text with the same information in prose (mirroring the feature description's own exact
phrasing: "recovery is straightforward" / "this is unrecoverable"). The existing `SeverityBanner`
frontend component (which already renders `Reason` as free text) picks up the enriched prose with
no frontend change required; the structured `RecoveryInfo` field is available for a future,
more visually distinct treatment (e.g. a colored badge) without another backend change.

**Rationale**: `FailureSeverity` is already consumed by exactly one frontend surface today
(`SeverityBanner`, which shows only the single highest-urgency severity system-wide, as prose).
Enriching `Reason` is the minimal change that satisfies FR-007/SC-002 today. Adding the structured
`RecoveryInfo` field alongside (rather than only enriching prose) keeps the distinction
machine-readable per FR-008 ("visibly distinguishable") without over-building a new UI component
this iteration doesn't clearly need — `SeverityBanner`'s existing `isCritical` color-branching logic
can trivially check `RecoveryInfo` in the same pass once a designer wants that.

**Alternatives considered**: Encode recoverability only in `Reason` prose, no structured field —
rejected: would make FR-008's "visibly distinguishable" dependent on string-parsing in the
frontend, which is fragile and untestable in the way a typed field isn't. Build a full new
per-cluster severity list UI (beyond the single top banner) as part of this feature — rejected as
out of scope: the existing single-banner limitation (only the highest severity across all clusters
is shown) is a pre-existing 006 design choice, not something this feature's spec asks to fix.

## R5: Backup-to-cluster coverage matching

**Decision**: A Backup covers a Cluster when either:
1. The Backup's `spec.includedNamespaces` is empty/`["*"]` (Velero's "all namespaces" convention)
   or contains the Cluster's own namespace, **or**
2. The Backup's `spec.labelSelector` matches the `cluster.x-k8s.io/cluster-name=<name>` label.

The most recent `Completed`-phase covering Backup determines staleness/recoverability; a covering
Backup that exists but has never completed (`Failed`, `PartiallyFailed`, `InProgress`) does not
count as a verified recovery point (spec.md Edge Cases), though it is not hidden from raw
visibility (`day2ops.BackupHealth.ClusterCoverage` still lists it).

**Rationale**: Velero has no native concept of a CAPI Cluster — this is a best-effort heuristic
against two real, common operator conventions for backing up a CAPI management cluster's
Cluster-owned resources (documented in spec.md Assumptions): namespace-scoped backups (CAPI's own
convention of keeping a Cluster and its owned Machines/infra objects in one namespace) and
label-selector-scoped backups. Matching via *either* signal, rather than requiring both, is
deliberately permissive: a false negative here (failing to recognize a real covering backup) is
worse than a false positive, since spec.md's CA-loss cross-reference is meant to prevent operators
from believing "unrecoverable" when a backup actually exists.

**Alternatives considered**: Require an explicit, project-specific annotation on each Backup
naming which Cluster it covers — rejected: would only work for backups created after this feature
ships, providing zero value for existing backup history, which is precisely the situation an
operator facing CA loss today is in.

## R6: `Cluster.spec.paused` visibility is already fully implemented — no new work

**Finding** (not a decision — a verification result): `models.Cluster.Paused` already exists
(`webserver/internal/infra/models/cluster.go`), is already populated by `processor.ProcessCluster`,
and is already rendered on the existing Cluster detail page's Specification tab
(`front/app/ui/dashboard/components/clusters/specification.tsx` line 19-20, `types.tsx` line 22,
`details.tsx` line 23). Confirmed via direct inspection, not assumed.

**Impact**: spec.md's FR-010 ("System MUST show whether a cluster's reconciliation is currently
paused") and US3's pause-visibility acceptance scenario are already satisfied by prior work. No
tasks are needed for this part of the feature; `quickstart.md`'s corresponding verification step
is a confirmation check, not new functionality to build.

## R7: Restore activity surfaces as an aggregate + per-cluster enrichment, not a new list UI

**Decision**: `day2ops.BackupHealth` carries a `RestoresInProgress int` aggregate count (shown on
the new card), and the same per-cluster coverage computation (R5) also looks up that cluster's most
recent Restore, threading its outcome into the CA-loss severity's enriched `Reason`/`RecoveryInfo`
when a restore is relevant to that cluster's current state (in progress, or most recently failed).
No new frontend list of individual Restore objects is built this iteration.

**Rationale**: spec.md's US3 acceptance scenarios ask that restore-in-progress and last-restore-
outcome be visible per cluster, not that every historical Restore be browsable — the aggregate
count plus the targeted enrichment where it matters most (a cluster already flagged CA-lost)
satisfies the acceptance scenarios without a new page. This keeps the feature's frontend footprint
to one new card, matching R3's scope decision.

**Alternatives considered**: A dedicated Restore list/detail panel akin to `risk-warnings.tsx`'s
list style — deferred, not rejected outright: if a future need for browsing restore history
independent of CA-loss context emerges, this is a natural follow-up; nothing in this plan's data
model precludes it (the same `day2opsStore.restores` map would back it).

## R8: Detecting whether Velero is installed

**Decision**: Reuse the `apiextensionsclientset.Clientset` already constructed in `WatchDay2Ops`
(`apiextClient`, currently used only for version-skew checks) to check for the existence of the
`backups.velero.io` CustomResourceDefinition before including any Velero GVR in
`day2opsWatchedGVRs`. Absence (a `NotFound` error, or `apiextClient == nil`) means Velero isn't
installed: no Velero GVRs are watched, and `BackupHealth.Available` is `false`, driving the
explicit "not available" card state (FR-011) — the same non-fatal, detect-before-watching pattern
006 was corrected to use for provider-specific GVRs, rather than attempting the watch and treating
the resulting error as fatal.

**Rationale**: Velero is not a `clusterctl` "infrastructure provider" (`clusterctlv1.ProviderList`
has no entry for it), so the existing `clusterapi.GenerateInfrastructureCapability` mechanism used
for Docker/vSphere detection doesn't cover it. A direct CRD-existence check is the smallest correct
primitive, and the `apiextClient` needed for it is already constructed on every `WatchDay2Ops` call.

**Alternatives considered**: Attempt the watch unconditionally and rely on `day2opsWatchedGVRs`'s
existing non-fatal-per-GVR-failure fallback (log and skip) — technically sufficient given that
fallback already exists, but rejected as the *sole* mechanism: it would log four errors on every
connection to a non-Velero-installed cluster (the common case for many deployments) instead of
zero, and wouldn't let `BackupHealth.Available` be set proactively before the first watch attempt
even completes.
