# Data Model: Velero Backup Recoverability Awareness

All new types live in `webserver/internal/infra/clusterapi/day2ops` (the existing pure-domain
package for Day-2 Ops), following the file-per-concept convention already used by
`severity.go`/`risk_certexpiry.go`. None of these are CAPI or Velero API types — they are this
project's own normalized, minimal views, decoded from unstructured Velero objects.

## BackupStorageLocationStatus

Normalized summary of one `BackupStorageLocation`'s reachability.

| Field | Type | Notes |
|---|---|---|
| `Name` | `string` | |
| `Namespace` | `string` | |
| `Reachable` | `bool` | From `status.phase == "Available"` |
| `Default` | `bool` | From `spec.default` |

## ClusterBackupCoverage

Per-cluster recoverability computed from the current set of known Backups/Restores (research.md
R5). One entry per known Cluster, always present (even when no covering backup exists — spec.md
Edge Cases: "never omitted").

| Field | Type | Notes |
|---|---|---|
| `ClusterRef` | `ObjectRef` | Existing `day2ops.ObjectRef` type |
| `Covered` | `bool` | Whether any `Completed`-phase Backup covers this cluster (R5) |
| `MostRecentBackupAge` | `*time.Duration` | Nil when `Covered` is false |
| `MostRecentBackupName` | `string` | Empty when `Covered` is false |
| `Stale` | `bool` | `Covered && MostRecentBackupAge > RPO threshold` |
| `RestoreInProgress` | `bool` | A `Restore` referencing a covering Backup is currently `InProgress` |
| `LastRestoreOutcome` | `string` | `""` (none), `"succeeded"`, or `"failed"` — most recent completed Restore for this cluster, if any |

## BackupHealth

The full Backup Health payload for the Day-2 Ops landing page (FR-001–FR-005, FR-009, FR-011).

| Field | Type | Notes |
|---|---|---|
| `Available` | `bool` | `false` when Velero isn't installed (research.md R8) — drives the "not available" state |
| `StorageLocations` | `[]BackupStorageLocationStatus` | Empty when `Available` is false |
| `ClusterCoverage` | `[]ClusterBackupCoverage` | Empty when `Available` is false |
| `RPOThreshold` | `time.Duration` | The threshold used to compute `Stale`, echoed back so the frontend can display it (FR-004: configurable, not hardcoded in the UI layer) |

JSON tags mirror the existing `day2ops.Data` struct's `json:"..."` convention (camelCase).

## RecoveryInfo (addition to existing `FailureSeverity`)

```go
// RecoveryInfo augments a CA-secret-missing FailureSeverity with backup-based recoverability
// (spec.md US2). Nil for every other severity kind.
type RecoveryInfo struct {
    Recoverable        bool   `json:"recoverable"`
    CoveringBackupAge   string `json:"coveringBackupAge,omitempty"` // formatted duration, empty if not recoverable
}
```

`FailureSeverity` gains one new field:

```go
type FailureSeverity struct {
    ObjectRef    *ObjectRef    `json:"objectRef"`
    Level        SeverityLevel `json:"level"`
    Reason       string        `json:"reason"`
    RecoveryInfo *RecoveryInfo `json:"recoveryInfo,omitempty"` // NEW, additive
}
```

## day2opsStore additions (`webserver/internal/web/watchers/day2ops.go`)

Following the existing store's map-keyed-by-`namespace/name` convention:

| New field | Type | Populated from GVR |
|---|---|---|
| `backups` | `map[string]unstructured.Unstructured` (or a small decoded struct — see note) | `backups.velero.io/v1` |
| `restores` | `map[string]unstructured.Unstructured` | `restores.velero.io/v1` |
| `schedules` | `map[string]unstructured.Unstructured` | `schedules.velero.io/v1` |
| `backupStorageLocations` | `map[string]unstructured.Unstructured` | `backupstoragelocations.velero.io/v1` |

**Note**: unlike the existing store's CAPI maps (typed `clusterv1.X`), these are kept as decoded
lightweight internal structs (e.g. an unexported `veleroBackup{Name, Namespace, Phase string;
IncludedNamespaces []string; LabelSelector *metav1.LabelSelector; StorageLocation string;
CompletionTimestamp *time.Time}`) built once at `apply()` time via `unstructured.Nested*` reads
(research.md R1) — not raw `unstructured.Unstructured`, so downstream computation
(`ComputeBackupHealth`, `ComputeClusterBackupCoverage`) stays fully typed and testable without
re-parsing on every recomputation.

## Frontend types (`front/app/ui/dashboard/shared/use-day2-ops.ts`)

Mirrors the backend JSON shape exactly (existing convention — this file already mirrors
`day2ops.Data`):

```ts
export type BackupStorageLocationStatus = {
  name: string; namespace: string; reachable: boolean; default: boolean;
};

export type ClusterBackupCoverage = {
  clusterRef: ObjectRef; covered: boolean; mostRecentBackupAge?: string;
  mostRecentBackupName?: string; stale: boolean; restoreInProgress: boolean;
  lastRestoreOutcome: '' | 'succeeded' | 'failed';
};

export type BackupHealth = {
  available: boolean; storageLocations: BackupStorageLocationStatus[];
  clusterCoverage: ClusterBackupCoverage[]; rpoThresholdSeconds: number;
};

export type RecoveryInfo = { recoverable: boolean; coveringBackupAge?: string };
// FailureSeverity gains: recoveryInfo?: RecoveryInfo
```

`Data` (the existing top-level WS payload type) gains `backupHealth: BackupHealth`.
