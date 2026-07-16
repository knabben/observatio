# Quickstart: Velero Backup Recoverability Awareness

Manual verification scenarios, one per user story, against a `kind-capi-mgmt`-style management
cluster with Velero installed (a local `MinIO` or filesystem `BackupStorageLocation` is sufficient
— no real cloud storage required).

## US1 — Backup Health at a glance

1. Install Velero with one reachable `BackupStorageLocation` and run a Backup for a test cluster's
   namespace; confirm it completes.
2. Open the Day-2 Ops landing page; confirm a Backup Health card appears alongside the existing
   Cluster/Machine/MachineDeployment rollups, showing the storage location as reachable and the
   cluster's backup as on-time.
3. Wait past the configured RPO threshold (or lower the threshold for testing) without running
   another backup; confirm the same cluster now shows as stale.
4. Point the BackupStorageLocation's credentials at something invalid (or otherwise make it
   unreachable); confirm the card reflects it as unreachable.
5. Add a second test cluster with no Backup at all; confirm it's shown as having no backup
   coverage — not silently absent from the list.

## US2 — CA loss vs. recoverability

1. In a disposable test environment, delete a cluster's `<cluster>-ca` Secret after confirming a
   recent covering Backup exists; confirm the resulting CA-secret-missing severity (top banner)
   indicates the cluster is recoverable and shows the backup's age.
2. Repeat against a cluster with no covering Backup; confirm the severity instead indicates no
   covering backup exists / not recoverable.
3. Repeat against a cluster whose only covering Backup is old (e.g. 30 days); confirm the age is
   shown so the operator can judge potential data loss, without claiming full recoverability is
   equivalent to a fresh backup.

## US3 — Restore activity and pause visibility

1. Start a Restore from an existing Backup; while it's running, confirm the Day-2 Ops landing page
   reflects a restore in progress (aggregate count on the Backup Health card).
2. Let the Restore complete successfully; confirm the outcome is reflected (and, if that cluster
   was also CA-lost, that the severity's recovery context reflects the completed restore).
3. Trigger a Restore that fails (e.g. restore into a namespace with a conflicting resource);
   confirm the failure is reflected, not silently treated as success.
4. Set `spec.paused: true` on a Cluster; open that cluster's own detail page (`/dashboard/clusters`
   → select it → Specification tab) and confirm the paused state is visible — this should already
   work today (research.md R6); this step is a confirmation, not new functionality.

## Cross-cutting checks

- With Velero not installed at all on the management cluster, confirm the Backup Health card shows
  an explicit "not available" state — not an error, not a blank card — and that CA-secret-missing
  severities still fire but without a recoverability claim (omitted `recoveryInfo`, not a false
  `false`).
- Confirm Backup Health / recoverability updates arrive over the existing Day-2 Ops WebSocket
  connection (check browser devtools' WS frames) — no new REST polling introduced.
- Confirm no mutating request (POST/PUT/PATCH/DELETE) is ever issued against any Velero, CAPI, or
  cluster resource by this feature — read-only, per spec.md FR-013.
- Run `make run-tests-backend` and `make run-tests-frontend`; confirm both pass with the new
  `day2ops` package tests (coverage matching, backup health) and the new `BackupHealthCard` tests
  included.
