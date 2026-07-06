# Phase 1 Data Model: Day-2 Operations Dashboard

## DebugLayer

One stage of evidence in an object's debugging path.

| Field | Type | Notes |
|-------|------|-------|
| `Layer` | enum: `conditions` \| `phase` \| `provider_resource` \| `controller_activity` | Fixed ordering, matches FR-005 |
| `Status` | enum: `ok` \| `implicated` \| `inconclusive` | Whether this layer explains the failure |
| `Evidence` | string[] | Human-readable evidence lines (condition reason/message, phase name, provider condition/event, Kubernetes Event message) |
| `Source` | string | Which underlying object the evidence came from (e.g. `Machine/worker-0`, `DockerMachine/worker-0`) |

## DebugPath

The full ordered result for one unhealthy object (FR-004–FR-007).

| Field | Type | Notes |
|-------|------|-------|
| `ObjectRef` | GVR + namespace/name | The object being explained |
| `Layers` | `DebugLayer[]` | Always all four layers, in fixed order; `controller_activity` is only populated when earlier layers are `inconclusive` (FR-007) |
| `Summary` | string | One-line rollup, e.g. "Waiting on infrastructure provisioning (DockerMachine)" |

## HealthRollup

Per-category summary shown on the landing view (FR-002).

| Field | Type | Notes |
|-------|------|-------|
| `Category` | enum: `cluster` \| `machine_deployment` \| `machine` | Extends today's `ClusterSummary` shape |
| `Healthy` | int | |
| `Degraded` | int | New tri-state addition — object has a Risk Warning or an `implicated`/`inconclusive` DebugLayer but is not failed |
| `Failed` | int | Existing "Failing" concept from `dashboard.go` |
| `Unavailable` | bool | True when the underlying watch/source is unreachable (FR-017) — distinct from all-zero-counts |

## RiskWarning

A proactively detected issue (FR-008–FR-011, US3).

| Field | Type | Notes |
|-------|------|-------|
| `ObjectRef` | GVR + namespace/name | Object the warning is attached to |
| `Kind` | enum: `cert_expiry` \| `stalled_rollout` \| `version_skew` \| `drift` | |
| `Detail` | string | e.g. expiry date, blocking PDB name, CRD version pair, generation mismatch |
| `LikelyCause` | string, optional | Populated when determinable (FR-009); empty string means "not determinable", never omitted (FR-018) |
| `CheckStatus` | enum: `evaluated` \| `not_evaluable` | Distinguishes "checked, no issue" from "could not check" (FR-018) |

## FailureSeverity

Classification attached to a detected issue (FR-012–FR-016, US4).

| Field | Type | Notes |
|-------|------|-------|
| `ObjectRef` | GVR + namespace/name, optional | Empty for cluster-wide severities (management-critical) |
| `Level` | enum: `self_healing` \| `needs_investigation` \| `provider_degraded` \| `management_critical` | Strictly increasing urgency |
| `Reason` | string | e.g. "MachineHealthCheck remediating NotReady node", "maxUnhealthy threshold breached", "capd-controller-manager crash-looping", "API server unreachable", "cluster CA secret missing" |

## Day2OpsEvent (WS payload)

The single event type broadcast by the new aggregator (research.md R9), matching the existing
`EventResponse{Type, Event, Data}` envelope from `webserver/internal/web/watchers/processor.go`.

| Field | Type | Notes |
|-------|------|-------|
| `Type` | `"MODIFIED"` | Always a delta recompute; no meaningful ADDED/DELETED distinction at aggregate level |
| `Event` | `"day2ops"` | New `Event` discriminator value |
| `Data.Rollups` | `HealthRollup[]` | One per category |
| `Data.DebugPaths` | `DebugPath[]` | One per currently-unhealthy object, so layers/status render on the landing screen itself (FR-004) without a drill-in round trip; `Evidence` arrays here are capped to their first line each — the uncapped, full evidence list is fetched on demand from `GET /api/day2ops/detail` only when the operator expands one (research.md R9) |
| `Data.Risks` | `RiskWarning[]` | All currently-detected risks across categories |
| `Data.Severities` | `FailureSeverity[]` | All currently-classified severities, cluster-wide and per-object |
| `Data.SourceUnavailable` | bool | Top-level "data unavailable" flag (FR-017), independent of per-category `Unavailable` |

## ControllerRef

Identifies which controller a `debug_layer: controller_activity` entry or a Logs-view request
refers to (User Story 5).

| Field | Type | Notes |
|-------|------|-------|
| `Namespace` | string | e.g. `capi-system`, `capd-system` |
| `DeploymentName` | string | e.g. `capi-controller-manager`, `capd-controller-manager` |
| `PodName` | string, optional | Resolved at request time from the Deployment's current Pod; not persisted |

## NodeAccessInfo

Static SSH connection guidance for a Machine on a VM-based provider (FR-021, FR-022). Never
includes credentials.

| Field | Type | Notes |
|-------|------|-------|
| `ObjectRef` | GVR + namespace/name | The Machine |
| `Command` | string | e.g. `ssh capi@10.0.1.23` — address sourced from `Machine.status.addresses` |
| `Note` | string | Fixed disclaimer that Observātiō does not manage credentials/connections |

## Frontend type mapping

TypeScript mirrors of the above live alongside the existing shared types
(`front/app/ui/dashboard/shared/status.ts` gains the `'degraded'` `StatusState` member; new types
`DebugLayer`, `DebugPath`, `HealthRollup`, `RiskWarning`, `FailureSeverity`, `Day2OpsEvent` are
added in `front/app/ui/dashboard/shared/use-day2-ops.ts` alongside the hook that consumes them,
following the existing pattern of colocating a resource's WS hook with its types (see
`resource-hooks.ts`). `ControllerRef` and `NodeAccessInfo` are defined alongside
`front/app/ui/dashboard/logs/logs-view.tsx` and `node-access-panel.tsx` respectively.
