# Contract: `day2ops` WebSocket Event

## Transport

Delivered over the existing WebSocket connection pool
(`webserver/internal/web/handlers/system/pool.go`), using the same envelope shape as every other
resource event (`webserver/internal/web/watchers/processor.go`'s `EventResponse`).

## Trigger

Broadcast whenever any resource contributing to the rollup/risk/severity computation changes:
Cluster, MachineDeployment, Machine, MachineSet, MachineHealthCheck, provider-infra objects
(DockerMachine/DockerCluster, VSphereMachine/VSphereCluster), or on the aggregator's own periodic
re-check of time-based conditions that don't correspond to a watch event (certificate expiry
crossing the warning threshold, a rollout crossing the stall grace period, generation-mismatch
drift persisting past its grace period).

## Payload

```json
{
  "Type": "MODIFIED",
  "Event": "day2ops",
  "Data": {
    "Rollups": [
      {"Category": "cluster", "Healthy": 2, "Degraded": 1, "Failed": 0, "Unavailable": false},
      {"Category": "machine_deployment", "Healthy": 3, "Degraded": 0, "Failed": 1, "Unavailable": false},
      {"Category": "machine", "Healthy": 8, "Degraded": 1, "Failed": 1, "Unavailable": false}
    ],
    "DebugPaths": [
      {
        "ObjectRef": {"group": "cluster.x-k8s.io", "version": "v1beta1", "resource": "machines", "namespace": "default", "name": "worker-0"},
        "Layers": [
          {"Layer": "conditions", "Status": "implicated", "Evidence": ["Ready=False: WaitingForInfrastructure"], "Source": "Machine/worker-0"},
          {"Layer": "phase", "Status": "implicated", "Evidence": ["Phase=Provisioning"], "Source": "Machine/worker-0"},
          {"Layer": "provider_resource", "Status": "implicated", "Evidence": ["Ready=False: VM creation failed"], "Source": "DockerMachine/worker-0"},
          {"Layer": "controller_activity", "Status": "inconclusive", "Evidence": [], "Source": ""}
        ],
        "Summary": "Waiting on infrastructure provisioning (DockerMachine: VM creation failed)"
      }
    ],
    "Risks": [
      {
        "ObjectRef": {"group": "cluster.x-k8s.io", "version": "v1beta1", "resource": "clusters", "namespace": "default", "name": "prod-1"},
        "Kind": "cert_expiry",
        "Detail": "prod-1-ca expires 2026-08-01",
        "LikelyCause": "",
        "CheckStatus": "evaluated"
      }
    ],
    "Severities": [
      {
        "ObjectRef": null,
        "Level": "management_critical",
        "Reason": "API server unreachable"
      }
    ],
    "SourceUnavailable": false
  }
}
```

## Consumer contract

- The frontend `use-day2-ops.ts` hook treats every `day2ops` event as a full-state replace, not a
  patch — `Data` always represents the complete current rollup/risk/severity/debug-path set,
  mirroring how `dashboard.go`'s existing `GenerateClusterSummary` returns a complete snapshot today.
- `Data.DebugPaths` evidence strings are truncated to one line per layer to bound WS payload size;
  `ops-dashboard.tsx` renders these directly (FR-004) and only calls `GET /api/day2ops/detail` when
  an operator expands a specific object for the full, uncapped evidence list.
- `Data.SourceUnavailable: true` means the aggregator itself lost its connection to the management
  cluster's API server; when true, `Rollups`/`Risks`/`Severities` MUST be treated as stale and the
  UI MUST show the FR-017 "data unavailable" state rather than the last-known values.
