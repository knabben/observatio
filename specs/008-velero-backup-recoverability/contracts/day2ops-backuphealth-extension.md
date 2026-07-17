# Contract: Day-2 Ops WS event — Backup Health extension

Extends the existing `day2ops` WebSocket event contract
(`specs/006-day2-ops-dashboard/contracts/day2ops-ws-event.md`) — same connection, same `ObjectType`
(`"day2ops"`), no new WS subscription type. Backward compatible: a client that ignores the two new
fields continues to work unchanged.

## Additive changes to `Data`

```jsonc
{
  "rollups": [ /* unchanged */ ],
  "debugPaths": [ /* unchanged */ ],
  "risks": [ /* unchanged */ ],
  "severities": [
    {
      "objectRef": {"group": "cluster.x-k8s.io", "version": "v1beta1", "resource": "clusters", "namespace": "default", "name": "capi-workload"},
      "level": "management_critical",
      "reason": "Cluster capi-workload's CA secret is missing or inaccessible — certificate issuance/rotation is blocked for its nodes; the original CA cannot be substituted. A backup completed 3h ago covers this cluster — recovery is straightforward.",
      // NEW, additive, omitted (not null) for every non-CA-related severity:
      "recoveryInfo": {"recoverable": true, "coveringBackupAge": "3h0m0s"}
    }
  ],
  "sourceUnavailable": false,
  // NEW:
  "backupHealth": {
    "available": true,
    "storageLocations": [
      {"name": "default", "namespace": "velero", "reachable": true, "default": true}
    ],
    "clusterCoverage": [
      {
        "clusterRef": {"group": "cluster.x-k8s.io", "version": "v1beta1", "resource": "clusters", "namespace": "default", "name": "capi-workload"},
        "covered": true,
        "mostRecentBackupAge": "3h0m0s",
        "mostRecentBackupName": "capi-workload-nightly-20260709",
        "stale": false,
        "restoreInProgress": false,
        "lastRestoreOutcome": ""
      }
    ],
    "rpoThresholdSeconds": 86400
  }
}
```

## Velero-not-installed state

```jsonc
{
  // ... rollups/severities/etc. unaffected — CA-secret-missing severities still fire, just
  // without recoveryInfo (recoverability genuinely unknown, not false):
  "backupHealth": {
    "available": false,
    "storageLocations": [],
    "clusterCoverage": [],
    "rpoThresholdSeconds": 86400
  }
}
```

When `available` is `false`, a CA-secret-missing severity's `recoveryInfo` is omitted entirely
(not `{"recoverable": false}`) — the distinction matters: "no Velero installed, recoverability
unknown" must not be presented the same as "Velero installed, definitively no covering backup"
(spec.md Edge Cases via FR-011).

## Consumer contract

- Frontend `use-day2-ops.ts` decodes `backupHealth` the same way it already decodes `rollups` —
  no new WS subscription, no new REST endpoint.
- `BackupHealthCard` renders `available === false` as an explicit "Velero not installed" state
  (mirroring `HealthRollupCard`'s existing `unavailable` branch), never a blank card or error.
- `SeverityBanner` requires no code change to pick up the enriched `reason` prose; a future
  iteration MAY read `recoveryInfo` for a more structured visual treatment (research.md R4).
