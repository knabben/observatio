# Quick specification statement: Velero backup/restore integration

**Candidate feature branch**: `007-velero-backup-recovery`

Observātiō's Day-2 Ops dashboard (006) can tell an operator that a cluster's CA secret is missing or
that the management cluster's API server is degraded, and it correctly flags CA loss as the
highest-severity, unrecoverable-by-substitution failure — but it has no idea whether that cluster is
actually *recoverable*, because it has no visibility into Velero at all. Today "is there a backup"
lives entirely outside the tool, in the operator's head or in Velero's own CLI.

Add Velero awareness the same way every other CAPI kind is already watched in this codebase (see
`webserver/internal/web/watchers/` for the existing pattern): watch `Backup`, `Schedule`, `Restore`, and
`BackupStorageLocation` objects, and surface a "Backup Health" rollup on the Day-2 Ops landing page
(alongside the existing Cluster/Machine/MachineDeployment rollups from 006) showing last successful
backup time, staleness against a configurable RPO, and whether the configured
`BackupStorageLocation`(s) are reachable. Most importantly, cross-reference this with 006's existing
`severity.go` CA-secret-missing detection so the dashboard can distinguish "CA lost, but a 3-hour-old
backup covers this cluster — recovery is straightforward" from "CA lost, no covering backup exists —
this is unrecoverable." Also surface `Cluster.spec.paused` as visible state, since the guide's recovery
procedure requires pausing CAPI reconciliation before a Velero restore and unpausing after — whether an
actual pause/unpause *control* is built as part of this feature (the product's first mutating action,
a deliberate departure from its current read-only design) or deferred to a later iteration is an open
design question for planning, not something to assume up front.
