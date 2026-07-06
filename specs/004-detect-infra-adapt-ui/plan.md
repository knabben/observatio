# Implementation Plan: Infrastructure Provider Detection & Adaptive Listing Screens

**Branch**: `004-detect-infra-adapt-ui` | **Date**: 2026-07-05 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/004-detect-infra-adapt-ui/spec.md`

## Summary

Observātiō hardcodes vSphere as the only infrastructure provider: the Clusters screen always shows a
static "vSphere Clusters" tab and vSphere-typed models (`capv.VSphereCluster`, `capv.VSphereMachine`)
regardless of what the connected management cluster actually has installed, and Docker-backed
environments (common for local/dev CAPI setups) have no equivalent view. This feature adds
server-side detection of which infrastructure provider(s) — Docker, vSphere, or both — are installed
in the connected cluster and which back each individual Cluster/Machine, and makes the listing
screens (Clusters, Machines, Machine Deployments) render accordingly: a per-resource provider (+
version) indicator, a provider-specific infra detail view that only appears when that provider is
actually detected, a Docker-equivalent view where today only vSphere exists, and graceful fallback
for unrecognized providers or an environment with none detected. Detection is derived entirely from
data the backend already has access to — the clusterctl provider inventory
(`sigs.k8s.io/cluster-api/cmd/clusterctl/api/v1alpha3.Provider`, already used by the existing
`/api/clusters/components` endpoint to report `ProviderName`/`Version`) for environment-level
capability, and each resource's standard `infrastructureRef.Kind` for per-resource provider — no new
external configuration, and no new Go module dependency.

## Technical Context

**Language/Version**: Go 1.25 (backend, `webserver/`); TypeScript 5.9, React 19, Next.js 15.3 (App Router, static export) (frontend, `front/`)
**Primary Dependencies**: `sigs.k8s.io/cluster-api` v1.9.6 (`clusterv1`, `clusterctlv1` — both already imported), `sigs.k8s.io/cluster-api-provider-vsphere` v1.12.0 (`capv`, already imported), `sigs.k8s.io/controller-runtime` client, `k8s.io/client-go` v0.32.1 (including its dynamic client, already used by the WebSocket watchers). Docker infra details (`DockerCluster`/`DockerMachine`) are read via the existing dynamic/unstructured client and decoded into small local structs — **not** via `sigs.k8s.io/cluster-api/test/infrastructure/docker`, which was confirmed during implementation to be a separate nested Go module that would force `cluster-api` v1.9.6→v1.13.3 plus a cascade of `k8s.io/*`/`controller-runtime` upgrades (see research.md R3). No go.mod changes. Frontend: existing Mantine UI 7 (`Tabs`, `Badge`) and the shared config-driven components introduced by `003-screen-ui-refactor` (`shared/object-table.tsx`, `shared/status-indicator.tsx`, `shared/empty-state.tsx`)
**Storage**: N/A — stateless; all data read live from the connected cluster's API server
**Testing**: Go `testing` + `testify` with CAPI fake clients (existing pattern in `webserver/internal/infra/clusterapi/dashboard_test.go`); Jest 29 + `@testing-library/react` via the shared `test-render.tsx` helper (existing pattern from `003-screen-ui-refactor`)
**Target Platform**: Same single-binary deployment — Next.js static export embedded in and served by the Go binary
**Project Type**: Web application — both backend (`webserver/`) and frontend (`front/`) are touched
**Performance Goals**: Provider/version detection adds at most one extra `List` call (`clusterctlv1.ProviderList`, already performed today for `/api/clusters/components`) per Clusters-screen load; no measurable added latency budget beyond the existing <2s event-to-UI render target (Constitution II)
**Constraints**: Core domain models (`models.Cluster`, `models.Machine`) may only gain a plain string `Provider` discriminator derived from the standard `infrastructureRef.Kind` field — no proprietary provider type is promoted into the core domain (Constitution III); existing vSphere REST/JSON shapes must not break for current consumers (FR-010); WebSocket remains the transport for live Cluster/Machine/MachineDeployment state — the new capability/version data is near-static installation metadata, following the existing REST-only precedent of `/api/clusters/components`, `/api/clusters/classes`, `/api/clusters/summary` (Constitution II)
**Scale/Scope**: 3 listing screens (Clusters, Machines, Machine Deployments) + 1 new capability endpoint; 2 supported providers (Docker, vSphere) plus an Unknown fallback path; 12 FRs, 6 success criteria

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

Constitution v1.1.0 — five principles:

| Principle | Impact | Verdict |
|-----------|--------|---------|
| **I. Observability & Data Consolidation** | Replaces a hardcoded, single-provider assumption with a consolidated detection layer that surfaces provider + version for every cluster/machine, and an explicit "no supported provider" message instead of a silent empty screen. Strengthens observability. | ✅ PASS |
| **II. Real-Time Visibility** | Cluster/Machine/MachineDeployment state stays WebSocket-delivered, unchanged. The new provider-capability data (which providers are installed, and their version) is environment/installation metadata that only changes when an operator installs/upgrades a provider — not live cluster health state — and is served over REST, matching the existing precedent of `/api/clusters/components`, `/api/clusters/classes`, and `/api/clusters/summary`, all of which are REST-only today. | ✅ PASS |
| **III. ClusterAPI Resource Model Compliance** | The per-resource `Provider` field added to `models.Cluster`/`models.Machine` is a plain string derived from the standard CAPI `infrastructureRef.Kind` field — not a new proprietary type. Provider-specific detail payloads (vSphere today, Docker added here) continue to live in the existing separate `*Infra` models/fetchers layer, preserving the boundary already established for vSphere rather than expanding proprietary types into the core domain. | ✅ PASS (extends pre-existing pattern; see note below) |
| **IV. AI-Augmented Troubleshooting** | Not touched by this feature. | ✅ PASS (N/A) |
| **V. Test-Driven Quality** | New Go fetcher/handler logic (provider/version detection, Docker infra fetcher) ships with `testify` + CAPI fake-client tests; new/changed frontend components (dynamic tabs, provider badge, unknown/empty fallback) ship with Jest tests covering docker-only, vsphere-only, mixed, and neither-detected scenarios. `make run-tests-backend` and `make run-tests-frontend` must pass before merge. | ✅ PASS |

**Result**: No violations. No entries required in Complexity Tracking.

**Note on Principle III**: `models.ClusterInfra` already wraps `capv.VSphereClusterStatus` directly (a
pre-existing condition, not introduced by this feature). This plan does not fix that pre-existing
tension — doing so is out of scope for this feature — but it also does not deepen it: the new
`models.ClusterInfraDocker` mirrors the same existing shape/precedent for Docker, and the field added
to the *core* `Cluster`/`Machine` models is a plain string, not a provider-specific struct.

## Project Structure

### Documentation (this feature)

```text
specs/004-detect-infra-adapt-ui/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md         # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output
│   └── infra-detection-api.md
└── checklists/
    └── requirements.md  # Spec quality checklist (from /speckit-specify)
```

### Source Code (repository root)

```text
webserver/
├── internal/infra/clusterapi/
│   ├── dashboard.go                        # existing GenerateComponentVersions (clusterctlv1.ProviderList) reused for capability detection
│   └── fetchers/
│       ├── cluster.go                      # extend ClusterInput; add ProviderFromKind helper reuse
│       ├── cluster_infra_docker.go         # NEW: ListClusterInfraDocker — dynamic client + unstructured decode (no typed dockerv1 import)
│       └── machine.go                      # add Provider derivation to Machine/MachineInfra
├── internal/infra/models/
│   ├── cluster.go                          # add Provider string field to Cluster
│   ├── machine.go                          # add Provider string field to Machine; NEW MachineInfraDocker
│   └── capability.go                       # NEW: InfrastructureCapability{ Docker, VSphere ProviderStatus }
├── internal/infra/providerkind/
│   └── providerkind.go                     # NEW: shared FromKind(kind string) string helper (docker/vsphere/unknown)
└── internal/web/handlers/
    ├── handlers.go                         # NEW route: GET /api/infra/capabilities; extend /api/clusters/infra/list with ?provider=
    └── kubernetes/
        ├── dashboard.go                    # NEW HandleInfraCapabilities (filters existing GenerateComponentVersions output)
        └── cluster.go                      # extend HandleClusterInfraList to dispatch docker/vsphere by provider

front/
├── app/dashboard/clusters/page.tsx                          # dynamic tabs from /api/infra/capabilities (replace static 2-tab Tabs)
├── app/ui/dashboard/components/clusters/
│   ├── types.tsx                                            # add provider/version fields; NEW DockerInfraType
│   ├── table.tsx                                             # add provider+version badge column
│   └── infra/
│       ├── infra-lister.tsx                                  # generalize to accept a provider config (docker | vsphere)
│       └── infra-table.tsx                                    # column config per provider instead of hardcoded vSphere columns
├── app/ui/dashboard/shared/
│   └── provider-badge.tsx                                     # NEW: shared provider(+version)/unknown indicator, reused by Clusters/Machines rows
└── app/lib/
    └── capabilities.ts                                        # NEW: fetch + type for /api/infra/capabilities
```

**Structure Decision**: Existing single Go binary (`webserver/`) + existing single Next.js frontend
(`front/`) — no new top-level project. Backend work centers on three additions: (1) a
`providerkind` helper package shared by cluster/machine fetchers to turn an `infrastructureRef.Kind`
into `docker`/`vsphere`/`unknown`, (2) a new capability endpoint built on the **already-existing**
`clusterctlv1.ProviderList`-backed `GenerateComponentVersions` (no new provider-inventory mechanism),
and (3) a Docker-mirrored infra fetcher/model alongside the existing vSphere one — populated via the
existing dynamic/unstructured client rather than a typed Docker package (see research.md R3: the
typed package lives in a separate nested Go module and would force a major, unjustified dependency
upgrade). Frontend work replaces the hardcoded 2-tab Clusters page and vSphere-only
infra table with capability-driven, config-based equivalents built on the shared components
`003-screen-ui-refactor` already introduced (`shared/object-table.tsx`, `shared/status-indicator.tsx`,
`shared/empty-state.tsx`), and adds one new shared `provider-badge.tsx`. No backend WebSocket
message-shape changes; no changes to Dashboard-overview cluster topology (`GenerateClusterTopology`,
also currently vSphere-only) — out of scope per this feature's spec, which scopes adaptation to the
Clusters/Machines/Machine-Deployments listing screens only.

## Complexity Tracking

> No Constitution Check violations — this section is intentionally empty.
