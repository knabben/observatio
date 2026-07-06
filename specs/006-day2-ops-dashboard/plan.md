# Implementation Plan: Day-2 Operations Dashboard

**Branch**: `006-day2-ops-dashboard` | **Date**: 2026-07-06 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/006-day2-ops-dashboard/spec.md`

## Summary

A new centralized dashboard consolidates management-plane (Cluster, ClusterClass) and
workload-plane (Machine, MachineDeployment, MachineSet) health into one categorized landing view,
and — for any unhealthy object — automatically determines and displays which of the layered CAPI
debugging stages (object conditions → Machine phase → provider-infra resource → controller
reconciliation activity) explains the failure, using only data the backend already has or can read
without new external dependencies. It also proactively detects four known CAPI risk classes
(certificate expiry, stalled MachineDeployment rollouts, provider/CRD version skew, infrastructure
drift) and classifies detected issues by failure severity (self-healing, needs-investigation,
provider-degraded, management-cluster-critical) so operators can prioritize correctly. Live rollup
and severity data is delivered over the existing WebSocket pool via a new server-side aggregator
that subscribes to existing and newly-added resource watchers, keeping Principle II (Real-Time
Visibility) intact. All of the above — layers, status, and rollups — surface directly on the
landing screen itself; a fifth story adds a new "Logs" destination (lateral nav item + debugging-path
deep-dive) streaming the relevant CAPI/provider controller's own Pod log output via the standard
Kubernetes Pod-log API when the controller-activity layer is implicated, plus static SSH connection
instructions for VM-based-provider node access — no live terminal, no stored credentials, and no new
external dependency. Machine/node-level log streaming itself is an explicit TODO for a future
iteration.

## Technical Context

**Language/Version**: Go 1.25 (backend, `webserver/`); TypeScript 5.9, React 19, Next.js 15.3 (frontend, `front/`)
**Primary Dependencies**: No new external module. Backend newly *uses* (rather than adds) API groups already vendored transitively by `k8s.io/client-go`/`k8s.io/api` v0.32.1: `core/v1` (Events, Secrets, Pods — including the Pod-log subresource, the same one `kubectl logs` uses), `apps/v1` (Deployments, for provider-controller health), `policy/v1` (PodDisruptionBudget), and `k8s.io/apiextensions-apiserver` v0.31.3 (promoted from indirect to direct, for CRD version introspection). CAPI-native types already in use (`sigs.k8s.io/cluster-api` `clusterv1`) gain two new watched kinds: `MachineSet` and `MachineHealthCheck` — both first-class CAPI types, not new dependencies. Frontend: no new package; extends the existing `status.ts`/`StatusIndicator` primitives, Mantine components, and the existing `nav-links.tsx` lateral-navigation list already in use.
**Storage**: N/A — stateless; all data is live-derived from the management cluster's API server on each watch event, nothing persisted.
**Testing**: Go `testing` + `testify` (existing pattern) for each new detector/aggregator function; Jest 29 + `@testing-library/react` via `test-render.tsx` (existing pattern) for new components/hooks.
**Target Platform**: Same single-binary deployment — Next.js static export embedded in and served by the Go binary.
**Project Type**: Web application — both backend (`webserver/`) and frontend (`front/`) are touched.
**Performance Goals**: Rollup/severity/debugging-path data is push-driven off existing watch events (no independent polling loop); recomputation on each contributing watch event must stay well under the constitution's 2-second WS render-latency budget for a single management cluster's typical object count (tens to low hundreds of objects).
**Constraints**: New read access (Events, Secrets, Pods/Deployments, PDBs, CRDs, MachineSet, MachineHealthCheck, Pod logs) is read-only everywhere — no new write/mutate calls, no interactive exec. Certificate and CA-secret contents are read only to extract `NotAfter`/existence, never surfaced to the frontend or logs in raw form. Drift and version-skew detection are explicitly best-effort/heuristic (see research.md R3, R6) given no cloud-provider-side introspection is available read-only through the Kubernetes API alone. The Logs view streams controller logs only (Layer 4); Observātiō never stores, manages, or transmits SSH credentials, and Machine/node-level log streaming is explicitly out of scope for this feature (tracked as a follow-up TODO).
**Scale/Scope**: 1 new consolidated dashboard view; 5 user stories (P1 landing view, P1 layered debugging path, P2 proactive risk detection, P3 severity classification, P2 controller-logs deep-dive); 2 new CAPI watchers (MachineSet, MachineHealthCheck); ~5 new backend detectors; 1 new WS event type; 1 new lateral-nav destination; 23 FRs, 7 success criteria.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

Constitution v1.1.0 — five principles:

| Principle | Impact | Verdict |
|-----------|--------|---------|
| **I. Observability & Data Consolidation** | This is the feature's direct purpose: correlating conditions, Machine phase, provider-resource state, (new) Events, and now controller Pod logs across the full CAPI hierarchy into one consolidated, per-category view instead of scattered list pages and manual `kubectl` sessions. | ✅ PASS |
| **II. Real-Time Visibility** | Rollups, debugging paths, and severity banners are recomputed and pushed over the existing WebSocket pool whenever a contributing watched resource changes (new server-side aggregator subscribing to existing + new watchers) — no polling loop. On-demand deep-drill detail (e.g., full evidence list for a single object's path) may use a scoped REST hydration exception identical in shape to the one already justified for the raw-object endpoint in feature 005. | ✅ PASS (aggregator is WS-push; drill-in detail is a scoped, on-demand REST exception, same as 005's R2) |
| **III. ClusterAPI Resource Model Compliance** | The two newly-watched kinds (`MachineSet`, `MachineHealthCheck`) are first-class CAPI types already part of the Cluster → MachineDeployment → Machine hierarchy (MachineSet sits between MachineDeployment and Machine; MachineHealthCheck is a first-class CAPI remediation type) — no proprietary abstraction introduced. Provider-controller health/log access (Pods/Deployments/Pod-logs in provider namespaces) uses only generic Kubernetes types, keeping infra-provider specifics opaque per this principle; the Logs view's node-access instructions are static text, not a provider-specific integration. | ✅ PASS |
| **IV. AI-Augmented Troubleshooting** | Not required by this feature, but the new Debugging Path and Risk Warning data is structured condition/evidence data that naturally strengthens the existing AI panel's auto-context (feature 005) as a future enhancement — no conflict, no action required now. | ✅ PASS (no violation; noted synergy, out of scope here) |
| **V. Test-Driven Quality** | Every new detector (cert-expiry, stalled-rollout, version-skew, drift, severity classifier) and the debugging-path synthesizer get Go unit tests against CAPI fake-client fixtures; new frontend components/hooks get Jest tests, following the existing patterns. | ✅ PASS |

**Result**: No violations. No entries required in Complexity Tracking.

## Project Structure

### Documentation (this feature)

```text
specs/006-day2-ops-dashboard/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md         # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/            # Phase 1 output
│   ├── day2ops-ws-event.md
│   └── day2ops-detail-api.md
└── checklists/
    └── requirements.md   # Spec quality checklist (from /speckit-specify)
```

### Source Code (repository root)

```text
webserver/
├── internal/web/watchers/
│   ├── machineset.go                       # NEW: watch MachineSet (owned by MachineDeployment)
│   └── machinehealthcheck.go               # NEW: watch MachineHealthCheck
├── internal/infra/clusterapi/
│   ├── day2ops/                            # NEW package
│   │   ├── debugpath.go                    # Synthesizes ordered DebugLayer[] for an unhealthy
│   │   │                                    # object from already-fetched Conditions + Machine
│   │   │                                    # Phase + provider-infra object + (new) Events
│   │   ├── risk_certexpiry.go              # Reads CAPI-managed cert Secrets, computes expiry
│   │   ├── risk_stalledrollout.go          # MachineSet scale-down stall + PDB/finalizer heuristic
│   │   ├── risk_versionskew.go             # Extends existing provider-inventory read with CRD
│   │   │                                    # version check (apiextensions-apiserver clientset)
│   │   ├── risk_drift.go                   # metadata.generation vs status.observedGeneration
│   │   │                                    # mismatch on provider-infra objects
│   │   ├── severity.go                     # Self-healing / needs-investigation / provider-degraded
│   │   │                                    # / management-critical classifier (MachineHealthCheck,
│   │   │                                    # provider-controller Pod/Deployment status, API-server
│   │   │                                    # reachability, CA secret presence)
│   │   └── aggregator.go                   # Subscribes to existing + new watcher event streams,
│   │                                        # recomputes the consolidated Day2OpsEvent, broadcasts
│   │                                        # over the existing WS connection pool
│   └── fetchers/
│       ├── secrets.go                      # NEW: read-only Secret fetch scoped to CA/cert lookups
│       └── controllerlogs.go               # NEW: streams a controller Pod's logs via the standard
│                                             # Kubernetes Pod-log subresource (same data/mechanism
│                                             # as `kubectl logs`); resolves CAPI-core/provider
│                                             # controller Deployment→Pod in its known namespace
└── internal/web/handlers/kubernetes/
    ├── day2ops.go                          # NEW: GET /api/day2ops/detail (scoped REST drill-in,
    │                                         # mirrors raw.go's pattern) + WS event registration
    └── logs.go                             # NEW: GET /api/logs/controller (streaming response) +
                                              # GET /api/logs/node-access (static SSH instructions),
                                              # backing the new Logs destination (FR-019–FR-023)

front/
├── app/dashboard/
│   ├── page.tsx                             # MODIFIED: becomes the Day-2 Ops landing view; existing
│   │                                         # ClusterSummary/ClusterHierarchy/ClusterVersions widgets
│   │                                         # are retained as constituent category rollup cards
│   └── logs/
│       ├── layout.tsx                       # NEW: thin wrapper, matches existing dashboard/* routes
│       └── page.tsx                         # NEW: renders <LogsView/>, selectable controller +
│                                             # optional Machine-scoped node-access panel
├── app/ui/dashboard/
│   ├── components/ops/                      # NEW
│   │   ├── ops-dashboard.tsx                # Top-level orchestrator: category rollups + banners
│   │   ├── health-rollup-card.tsx           # Per-category healthy/degraded/failed counts
│   │   ├── debugging-path.tsx               # Ordered DebugLayer[] evidence, shown inline on the
│   │   │                                     # landing screen per FR-004; its controller_activity
│   │   │                                     # layer links to the Logs deep-dive when implicated
│   │   ├── risk-warnings.tsx                # Cert-expiry / stalled-rollout / version-skew / drift
│   │   │                                     # list, grouped by category
│   │   └── severity-banner.tsx              # Top-level provider-degraded / management-critical banner
│   ├── logs/                                # NEW
│   │   ├── logs-view.tsx                    # Controller selector + streamed log pane
│   │   └── node-access-panel.tsx            # Static SSH command/address instructions for a Machine
│   └── shared/
│       ├── status.ts                        # MODIFIED: add 'degraded' StatusState
│       ├── status-indicator.tsx             # MODIFIED: render the new 'degraded' state
│       └── use-day2-ops.ts                  # NEW: WS hook consuming the new Day2Ops event type
├── app/ui/dashboard/nav-links.tsx            # MODIFIED: add a "Logs" entry (`/dashboard/logs`)
└── app/styles/theme.ts                      # MODIFIED: STATUS_COLORS gains 'degraded' (amber)
```

**Structure Decision**: Existing single Go binary (`webserver/`) + existing single Next.js frontend
(`front/`) — no new top-level project. The existing top-level landing page
(`front/app/dashboard/page.tsx`) is redesigned in place to become the consolidated Day-2 Ops view
rather than adding a competing route, since the spec calls for *the* centralized place for Day-2
operations, not an additional one; its current widgets (`ClusterSummary`, `ClusterClassLister`,
`ClusterHierarchy`, `ClusterVersions`) are retained and re-homed as the category rollup building
blocks inside the new `components/ops/` module. Backend work is a new `day2ops` package following
the existing `webserver/internal/infra/clusterapi/` package-per-concern convention (mirroring
`processor/` and `fetchers/`), plus two new watchers following the existing
`webserver/internal/web/watchers/` pattern, with delivery unified through one new aggregator that
reuses the existing WebSocket connection pool (`webserver/internal/web/handlers/system/pool.go`).
The new Logs destination follows the same route-shell convention as the existing `dashboard/clusters`,
`dashboard/machines`, and `dashboard/machinedeployments` routes (`layout.tsx` + `page.tsx`), is added
to `nav-links.tsx` alongside the existing four entries, and its backend is a thin streaming handler
over the standard Kubernetes Pod-log subresource — no new external dependency, no Docker-daemon
access, no stored credentials.

## Complexity Tracking

> No Constitution Check violations — this section is intentionally empty.
