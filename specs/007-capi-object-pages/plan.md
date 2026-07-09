# Implementation Plan: First-Class Pages for MachineHealthCheck, KubeadmControlPlane, MachineSet, and ClusterClass

**Branch**: `007-capi-object-pages` | **Date**: 2026-07-08 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/007-capi-object-pages/spec.md`

## Summary

Four CAPI kinds — MachineHealthCheck, KubeadmControlPlane, MachineSet, ClusterClass — get the exact
same first-class treatment Clusters/Machines/MachineDeployments already have: a live WS-driven list
page, a detail screen (Specification tab with status/conditions + YAML tab with the complete raw
object), and a one-click "Ask AI about this" action. This is a horizontal extension of an existing,
already-generic pattern (`BaseLister`, `ObjectDetails`, `useCurrentObjectContext`/`AskAIButton`,
`/api/raw` + `RESOURCE_GVR`) — no new UI infrastructure, only four new backend watchers and four
small per-kind frontend module sets, following the Machines page file-for-file.

## Technical Context

**Language/Version**: Go 1.25 (backend, `webserver/`); TypeScript 5.9, React 19, Next.js 15.3 (frontend, `front/`)
**Primary Dependencies**: One new Go import: `sigs.k8s.io/cluster-api/controlplane/kubeadm/api/v1beta1` for the `KubeadmControlPlane` type (part of the already-vendored `sigs.k8s.io/cluster-api` module — no new external module, matching Constitution's "no new runtime dependency without justification"). MachineHealthCheck and MachineSet types (`clusterv1`) and ClusterClass (`clusterv1.ClusterClass`) are already imported. No new frontend dependency — reuses `BaseLister`, `ObjectDetails`, `useResourceStream`, `useCurrentObjectContext`, `AskAIButton`, `/api/raw`.
**Storage**: N/A — stateless, live-derived from the management cluster's API server, matching every existing list page.
**Testing**: Go `testing` + `testify` for new processors/watchers; Jest + `@testing-library/react` via `test-render.tsx` for new frontend components — both following the exact existing per-kind test file pattern (e.g. `machine_test.go`, and the Jest tests already covering `BaseLister`/`ObjectDetails` generically, so new per-kind tests only need to cover the thin per-kind wiring, not re-test the shared shell).
**Target Platform**: Same single-binary deployment (Next.js static export embedded in the Go binary).
**Project Type**: Web application — both `webserver/` and `front/` touched.
**Performance Goals**: Same as existing list pages — WS push, no polling; under 2s render latency (Constitution Principle II).
**Constraints**: Read-only throughout — no lifecycle actions (scale, upgrade, delete) for any of the four kinds, consistent with the dashboard's existing observability-only scope. KubeadmControlPlane and MachineHealthCheck CRDs may not exist in every environment (different control-plane provider, or MHC simply not configured) — each new page must degrade to an empty/unavailable state, not an error (FR-012), mirroring 006's fix for optional/provider-specific CRDs.
**Scale/Scope**: 4 new watchers, 4 new WS dispatch entries, 4 new `RESOURCE_GVR` entries, 4 new frontend page module sets (list route + lister + table + details + specification), 1 existing widget left unchanged (ClusterClass); 12 FRs, 5 success criteria.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

Constitution v1.1.0 — five principles:

| Principle | Impact | Verdict |
|-----------|--------|---------|
| **I. Observability & Data Consolidation** | Directly serves this principle: MachineHealthCheck and MachineSet already have structured health data collected internally (006) but not surfaced to operators directly; KubeadmControlPlane's etcd/replica conditions are entirely new observability. All four now expose condition-level health status directly, per the principle's explicit requirement. | ✅ PASS |
| **II. Real-Time Visibility** | Every new list page is WS-driven via the existing generic `useResourceStream`/`WatchResourceViaWebSocket` pattern — no polling. The YAML tab uses the same scoped, already-justified REST exception (`/api/raw`, WS-triggered on `resourceVersion` change) established in 005. | ✅ PASS |
| **III. ClusterAPI Resource Model Compliance** | All four are first-class CAPI types (MachineHealthCheck and MachineSet already part of the Cluster → MachineDeployment → Machine hierarchy; KubeadmControlPlane is a first-class control-plane type; ClusterClass is a first-class CAPI type). No proprietary abstraction introduced. | ✅ PASS |
| **IV. AI-Augmented Troubleshooting** | Each new detail screen wires the existing `useCurrentObjectContext`/`AskAIButton` mechanism from 005 unchanged — the AI panel receives real structured condition data for these kinds for the first time. | ✅ PASS |
| **V. Test-Driven Quality** | New Go tests for each new processor (`ProcessMachineHealthCheck`, `ProcessKubeadmControlPlane`, `ProcessMachineSet`) and watcher wiring; new Jest tests for each per-kind `details.tsx`/`specification.tsx` follow the existing per-kind test pattern. | ✅ PASS |

**Result**: No violations. No entries required in Complexity Tracking.

## Project Structure

### Documentation (this feature)

```text
specs/007-capi-object-pages/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md         # Phase 1 output
├── quickstart.md         # Phase 1 output
├── contracts/            # Phase 1 output
│   └── watch-types.md
└── checklists/
    └── requirements.md   # Spec quality checklist (from /speckit-specify)
```

### Source Code (repository root)

```text
webserver/
├── internal/infra/models/
│   ├── machinehealthcheck.go        # NEW: MachineHealthCheck DTO (mirrors machinedeployment.go)
│   ├── kubeadmcontrolplane.go       # NEW: KubeadmControlPlane DTO
│   ├── machineset.go                # NEW: MachineSet DTO
│   └── clusterclass.go              # EXISTING — reused as-is
├── internal/infra/clusterapi/processor/
│   ├── machinehealthcheck.go        # NEW: ProcessMachineHealthCheck (+ Response variant)
│   ├── kubeadmcontrolplane.go       # NEW: ProcessKubeadmControlPlane
│   └── machineset.go                # NEW: ProcessMachineSet
└── internal/web/watchers/
    ├── machinehealthcheck.go        # NEW: WatchMachineHealthChecks (mirrors machine.go's WatchMachines)
    ├── kubeadmcontrolplane.go       # NEW: WatchKubeadmControlPlanes
    ├── machineset.go                # NEW: WatchMachineSets
    └── clusterclass.go              # NEW: WatchClusterClasses (existing fetcher/REST path for the
                                       # dashboard widget is untouched; this is the new live-list path)

front/
├── app/dashboard/
│   ├── machinehealthchecks/{layout.tsx, page.tsx}   # NEW
│   ├── kubeadmcontrolplanes/{layout.tsx, page.tsx}  # NEW
│   ├── machinesets/{layout.tsx, page.tsx}           # NEW
│   └── clusterclasses/{layout.tsx, page.tsx}        # NEW
├── app/ui/dashboard/components/
│   ├── machinehealthchecks/{lister.tsx, table.tsx, details.tsx, specification.tsx}  # NEW
│   ├── kubeadmcontrolplanes/{lister.tsx, table.tsx, details.tsx, specification.tsx} # NEW
│   ├── machinesets/{lister.tsx, table.tsx, details.tsx, specification.tsx}          # NEW
│   └── clusterclasses/{lister.tsx, table.tsx, details.tsx, specification.tsx}       # NEW
├── app/ui/dashboard/nav-links.tsx    # MODIFIED: 4 new entries
└── app/lib/resource-gvr.ts           # MODIFIED: 4 new RESOURCE_GVR entries
```

**Structure Decision**: Existing single Go binary (`webserver/`) + existing single Next.js frontend
(`front/`) — no new top-level project, no new generic infrastructure. Each of the four kinds gets
exactly the same four-file backend addition (model, processor, watcher, WS dispatch entry) and
five-file frontend addition (route shell, lister, table, details, specification) that Machines
already has, verified file-for-file against `webserver/internal/web/watchers/machine.go` and
`front/app/ui/dashboard/components/machines/*` during Phase 0 research. ClusterClass's existing
fetcher/processor/REST path (used by the main-dashboard widget) is left untouched; only a new
watcher is added for the new page's live list.

## Complexity Tracking

> No Constitution Check violations — this section is intentionally empty.
