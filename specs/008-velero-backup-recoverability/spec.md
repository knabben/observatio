# Feature Specification: Velero Backup Recoverability Awareness

**Feature Branch**: `008-velero-backup-recoverability`
**Created**: 2026-07-09
**Status**: Draft
**Input**: User description: "Observātiō's Day-2 Ops dashboard (006) can tell an operator that a cluster's CA secret is missing or that the management cluster's API server is degraded, and it correctly flags CA loss as the highest-severity, unrecoverable-by-substitution failure — but it has no idea whether that cluster is actually *recoverable*, because it has no visibility into Velero at all. Today "is there a backup" lives entirely outside the tool, in the operator's head or in Velero's own CLI."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Backup Health at a glance (Priority: P1)

An operator opens the Day-2 Ops landing page and, alongside the existing Cluster/Machine/
MachineDeployment rollups, sees a Backup Health summary: whether each configured backup storage
destination is currently reachable, and whether each cluster's most recent successful backup is
within or beyond an acceptable recovery point objective (RPO).

**Why this priority**: This is the foundational visibility the feature description says is
completely missing today ("is there a backup" lives entirely outside the tool). Without it, no
other Velero-aware capability in this feature has data to build on.

**Independent Test**: With Velero installed and at least one Backup and BackupStorageLocation
present on the management cluster, open the Day-2 Ops landing page and confirm a Backup Health
rollup appears showing storage-location reachability and per-cluster backup staleness — without
needing to touch CA-loss detection or any other severity logic.

**Acceptance Scenarios**:

1. **Given** Velero is installed with a reachable BackupStorageLocation and a Backup completed 2
   hours ago for a cluster, **When** the operator opens the Day-2 Ops landing page, **Then** that
   cluster's backup is shown as on-time against the RPO.
2. **Given** a cluster's most recent successful backup completed 3 days ago and the configured RPO
   is 24 hours, **When** the operator views the Backup Health rollup, **Then** that cluster's
   backup is shown as stale.
3. **Given** a configured BackupStorageLocation cannot currently be reached (e.g. expired storage
   credentials), **When** the operator views the Backup Health rollup, **Then** that location is
   shown as unreachable, and backups depending on it are flagged as reduced-confidence.
4. **Given** a cluster has no Backup covering it at all, **When** the operator views the Backup
   Health rollup, **Then** that cluster is shown as having no backup coverage, not silently
   omitted.

---

### User Story 2 - Know whether CA loss is actually recoverable (Priority: P1)

When the dashboard flags a cluster's CA secret as missing — its highest-severity,
unrecoverable-by-substitution failure (006) — the operator sees, in the same place, whether a
backup exists that covers that cluster and how old it is, so they can immediately tell "CA lost,
but a recent backup covers this cluster — recovery is straightforward" apart from "CA lost, no
covering backup exists — this cluster's data is unrecoverable."

**Why this priority**: This is the flagship value of the feature: turning a dead-end,
highest-severity alert into an actionable recoverability signal, closing the exact gap named in
the feature description.

**Independent Test**: Trigger (or simulate) a CA-secret-missing severity for a cluster that has a
recent covering backup, and separately for a cluster with no covering backup; confirm the two
cases are distinguishable from the severity display alone, without consulting Velero's CLI.

**Acceptance Scenarios**:

1. **Given** a cluster is flagged with CA-secret-missing severity and has a covering backup
   completed 3 hours ago, **When** the operator views that cluster's severity, **Then** the
   dashboard indicates the cluster is recoverable via that backup and shows its age.
2. **Given** a cluster is flagged with CA-secret-missing severity and has no covering backup,
   **When** the operator views that cluster's severity, **Then** the dashboard indicates no
   covering backup exists and the cluster's data is not recoverable through Velero.
3. **Given** a cluster is flagged with CA-secret-missing severity and its only covering backup
   completed 30 days ago, **When** the operator views that cluster's severity, **Then** the
   dashboard shows a backup exists but flags its age so the operator can judge how much data would
   be lost.

---

### User Story 3 - See recovery activity and reconciliation-pause state (Priority: P2)

An operator who is actively running (or has just completed) a Velero-based recovery can see, from
the dashboard, whether a Restore is in progress or how the most recent one concluded, and whether
the affected cluster's reconciliation is currently paused — since the documented recovery
procedure requires pausing CAPI reconciliation before a restore and unpausing after, and getting
that step wrong risks the controller fighting the restore.

**Why this priority**: This closes the loop on the recovery procedure itself, not just the
decision to start one. It is secondary to knowing recoverability exists at all (US1/US2), but
prevents a common, costly operational mistake during the recovery itself.

**Independent Test**: With a Restore object present (in progress or completed) and a cluster with
`spec.paused` set, open the affected cluster's view and confirm both the restore outcome and the
paused state are visible without running kubectl.

**Acceptance Scenarios**:

1. **Given** a Restore is currently in progress for a cluster, **When** the operator views that
   cluster, **Then** the dashboard shows a restore is in progress.
2. **Given** a Restore for a cluster most recently completed successfully, **When** the operator
   views that cluster, **Then** the dashboard shows the successful outcome and when it completed.
3. **Given** a Restore for a cluster most recently failed, **When** the operator views that
   cluster, **Then** the dashboard shows the failure so the operator knows recovery did not
   complete as expected.
4. **Given** a cluster has `spec.paused` set to true, **When** the operator views that cluster,
   **Then** the dashboard visibly indicates reconciliation is paused.

---

### Edge Cases

- What happens when Velero is not installed on the management cluster at all (no Backup/Restore/
  BackupStorageLocation/Schedule objects exist)? The dashboard must show an explicit "not
  available" state for Backup Health, not an error or a blank/broken widget — consistent with how
  the dashboard already handles other optional, provider-dependent data (006).
- What happens when a BackupStorageLocation exists but is currently unreachable? It is shown as
  unreachable, and backups whose only storage location is unreachable are flagged as
  reduced-confidence rather than treated as fully verified.
- What happens when a cluster has never been backed up? It is shown as having no backup coverage —
  the same signal as a stale or absent backup, never silently omitted from the rollup.
- What happens when multiple Backups exist that cover the same cluster? The most recent
  successfully completed one determines staleness and recoverability.
- What happens when a Backup's outcome is "partially failed" rather than fully completed? It is
  not counted as a fully successful covering backup for recoverability purposes, though it still
  appears in raw backup visibility so the operator isn't misled either way.
- What happens when a Restore is in progress for a cluster whose CA is also flagged missing? Both
  facts are shown together — an in-progress restore doesn't hide the underlying severity until the
  restore actually completes.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST display a Backup Health summary on the Day-2 Ops landing page, alongside
  the existing Cluster/Machine/MachineDeployment rollups, whenever Velero backup data is available
  on the management cluster.
- **FR-002**: System MUST show, for each configured backup storage destination, whether it is
  currently reachable.
- **FR-003**: System MUST show, for each cluster, the time since its most recent successfully
  completed backup.
- **FR-004**: System MUST classify each cluster's most recent successful backup as on-time or
  stale against a recovery point objective (RPO) threshold, and that threshold MUST be
  configurable rather than hardcoded.
- **FR-005**: System MUST indicate when a cluster has no backup covering it at all, distinctly from
  a stale backup.
- **FR-006**: System MUST determine which backups cover a given cluster using the cluster's
  namespace and/or an association between the backup and the cluster's identity (see Assumptions
  for the exact matching rule) — a match must not depend on the operator manually correlating
  names.
- **FR-007**: System MUST, when a cluster carries a CA-secret-missing severity, additionally
  indicate whether a covering backup exists for that cluster and, if so, its age.
- **FR-008**: System MUST make a cluster with CA-secret-missing severity and a covering backup
  visibly distinguishable from a cluster with CA-secret-missing severity and no covering backup.
- **FR-009**: System MUST show recent Restore activity for a cluster, including whether a restore
  is currently in progress and the outcome (succeeded/failed) of the most recent completed
  restore.
- **FR-010**: System MUST show whether a cluster's reconciliation is currently paused
  (`spec.paused`).
- **FR-011**: System MUST present an explicit "not available" state for Backup Health — not an
  error, not a blank or broken widget — when Velero is not installed or no backup data exists on
  the management cluster.
- **FR-012**: System MUST refresh Backup Health, recoverability, restore-activity, and pause-state
  information live, reflecting changes without requiring a manual page reload, consistent with the
  dashboard's existing rollups.
- **FR-013**: System MUST NOT perform any mutating action against Velero, CAPI, or cluster
  resources as part of this feature — all backup, restore, and pause-state information is
  presented read-only (see Assumptions).

### Key Entities

- **Backup**: A completed or in-progress Velero backup operation. Key attributes: name, namespace,
  outcome/phase, start and completion time, which storage location it used, and enough information
  to determine which cluster(s) it covers.
- **Restore**: A Velero restore operation. Key attributes: name, outcome/phase, start and
  completion time, and which Backup it restored from.
- **BackupStorageLocation**: A configured backup storage destination. Key attributes: name,
  reachability/availability, and which Backups depend on it.
- **Schedule**: A recurring backup policy. Key attributes: name, recurrence, and the most recent
  Backup it produced — used to distinguish "no schedule configured" from "scheduled, but the last
  run failed."
- **Recoverability Assessment**: A derived judgment, not a raw object — for a given cluster, the
  combination of its current severity (from 006) and its backup coverage, classified as
  recoverable, unrecoverable, or unknown (no backup data available).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: An operator can determine whether any cluster's backups are stale relative to its RPO
  within 5 seconds of opening the Day-2 Ops dashboard, without leaving the page.
- **SC-002**: When a cluster's CA secret is reported missing, an operator can determine from the
  same view whether that cluster is recoverable via an existing backup — and the backup's age — 
  without consulting Velero's CLI or reading raw object YAML.
- **SC-003**: 100% of configured backup storage destinations' reachability is reflected on the
  dashboard within the same refresh cadence as the existing rollups, with no separate manual check
  required.
- **SC-004**: An operator can identify a cluster's reconciliation-pause state directly from the
  dashboard without running kubectl.
- **SC-005**: When Velero is not installed on the management cluster, the dashboard communicates
  "not available" for Backup Health in a way that is visibly distinct from an error, on first
  visit, with no configuration required.

## Assumptions

- Backup-to-cluster coverage matching uses the cluster's namespace (Velero backups scoped to that
  namespace) and/or a label association tying the backup to the cluster's identity, since Velero
  has no native concept of a CAPI Cluster. This is a best-effort match against existing Velero
  conventions for CAPI management-cluster backups, not a guarantee that every possible backup
  configuration is recognized; unmatched backups are still shown in raw form so nothing is hidden.
- A backup counts as "successful" for staleness and recoverability purposes only when fully
  completed; a partially-failed backup is shown but not treated as a verified recovery point.
- The RPO threshold defaults to a reasonable interval (e.g. 24 hours) and can be adjusted; it is
  not per-cluster in this iteration unless a single global default proves insufficient during
  implementation.
- This specification scopes Velero awareness to **read-only visibility** — of backups, restores,
  storage-location reachability, recoverability relative to CA loss, and reconciliation-pause
  state. Whether to add an active pause/unpause *control* (which would be this product's first
  mutating action, a deliberate departure from its current read-only design) is an open question
  left to a future decision, not assumed as part of this feature.
- Velero, if installed, is assumed to be installed on the same management cluster Observātiō
  already connects to (consistent with how 006 discovers other optional, provider-dependent data),
  not on a separate backup-only cluster.
- "Covering a cluster" is evaluated only for the management cluster's own Cluster API objects and
  their associated infrastructure/bootstrap objects, not for workload-cluster-internal application
  data, which is out of scope for this feature.
