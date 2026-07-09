# Phase 0 Research: First-Class Pages for MachineHealthCheck, KubeadmControlPlane, MachineSet, and ClusterClass

## R1: The existing Machines page is the exact template to mirror

**Decision**: Every one of the four new kinds follows the Machines page file-for-file:
`webserver/internal/web/watchers/machine.go`'s `WatchMachines` (dynamic-client watch + typed
convert via `processor.ProcessMachine`, registered in `websocket.go`'s `watchHandlers` map) on the
backend; `front/app/ui/dashboard/components/machines/{lister.tsx, table.tsx, details.tsx,
specification.tsx}` plus `front/app/dashboard/machines/{layout.tsx, page.tsx}` on the frontend, with
`BaseLister` (`front/app/ui/dashboard/base/lister.tsx`) and `ObjectDetails`
(`front/app/ui/dashboard/base/details.tsx`) doing all the generic list/detail/tab/search/error-state
work already.

**Rationale**: Confirmed via direct code inspection (Phase 0 research agent) that `BaseLister` and
`ObjectDetails` are already 100% generic — no Machine-specific logic lives in either. A new kind
needs zero changes to shared infrastructure, only the same thin per-kind wiring every existing kind
already has.

**Alternatives considered**: A single generic "any CAPI kind" page driven entirely by
config/metadata (no per-kind files at all) was considered but rejected as premature abstraction —
the existing three kinds don't do this, each has its own (thin) per-kind files, and introducing a
config-driven mega-abstraction for four more kinds would be a larger, riskier change than repeating
the same proven five-file pattern once more.

## R2: MachineHealthCheck and MachineSet need new standalone watchers, separate from 006's inline ones

**Decision**: `webserver/internal/web/watchers/day2ops.go` (006) already tracks both kinds, but only
as private fields inside `day2opsStore`, fed by watches opened inline inside `WatchDay2Ops` — there
is no reusable, independently-dispatchable watcher function for either kind today. New
`WatchMachineHealthChecks`/`WatchMachineSets` functions are added (mirroring `machine.go`), each
registered under a new `ObjectType` (`"machinehealthcheck"`, `"machineset"`) in `websocket.go`'s
`watchHandlers` map.

**Rationale**: `day2opsStore`'s internal state (keyed by namespace/name, minimal fields, no
`processor.Process*`-style DTO conversion) is shaped for 006's internal risk/severity computation,
not for a general-purpose, fully-detailed list/detail page. Reusing it directly would couple this
feature's UI needs to 006's internal aggregation state. A separate, standard watcher (same pattern
as every other kind) keeps the two concerns independent, at the cost of the management cluster being
watched for these two kinds via two separate connections when both 006's dashboard and one of these
new pages are open simultaneously — an accepted, existing characteristic of this architecture (every
kind already has this property; see 006 research.md R9's note on redundant per-tab watches).

## R3: KubeadmControlPlane type comes from `sigs.k8s.io/cluster-api/controlplane/kubeadm/api/v1beta1`

**Decision**: Import `sigs.k8s.io/cluster-api/controlplane/kubeadm/api/v1beta1` (aliased `kcpv1` by
convention) for the typed `KubeadmControlPlane` struct, watched via GVR
`{Group: "controlplane.cluster.x-k8s.io", Version: "v1beta1", Resource: "kubeadmcontrolplanes"}` —
confirmed live against the real `kind-capi-mgmt` test cluster during 006's topology fix validation
(`capi-workload-xg594controlplane.cluster.x-k8s.io/v1beta1, Resource=kubeadmcontrolplanes` appeared
in the topology graph's node IDs).

**Rationale**: This subpackage is part of the already-vendored `sigs.k8s.io/cluster-api` module
(same module already providing `clusterv1`), so importing it is not a new external dependency,
consistent with the Constitution's Technology Stack constraint — only a new import path within an
already-approved module.

**Alternatives considered**: Reading KCP generically via `unstructured.Unstructured` (as 006 did for
provider-infra objects, per that feature's Principle III scoping) was considered, but rejected here:
KCP is a first-class, always-present-when-kubeadm-is-used CAPI type (not a provider-specific,
possibly-absent CRD like `VSphereMachine`), so there's no Principle III reason to keep it generic,
and a typed struct gives cleaner access to its replica/etcd status fields for the Specification tab.

## R4: No new REST list endpoints

**Decision**: None of the four new kinds get a `GET /api/<kind>/list`-style REST handler. Only the
WS watcher (live list) and the existing generic `/api/raw` (YAML tab, keyed by a new `RESOURCE_GVR`
entry) are added.

**Rationale**: Phase 0 research confirmed the existing Machines/Clusters/MachineDeployments REST
list handlers (e.g. `HandleMachines`) exist but are *not* what powers their live list pages —
`BaseLister` uses `useResourceStream` (WS) exclusively. Adding unused REST list endpoints for the
four new kinds would be dead code the spec doesn't require (FR-001–FR-004 only require a "live list
page," not a REST list API) — a premature-abstraction violation of the "no code beyond what the task
requires" guidance.

**Alternatives considered**: Adding REST handlers "for consistency" with the older three kinds was
considered and rejected — those REST handlers are themselves unused legacy surface from before the
WS-only pattern was established, not something newer kinds should be obligated to replicate.

## R5: ClusterClass gets a new watcher; its existing REST/widget path is untouched

**Decision**: The existing `fetchers/clusterclass.go` → `processor/clusterclass.go` →
`getClusterClasses` REST path, consumed only by the main-dashboard `ClusterClassLister` widget
(`front/app/ui/dashboard/components/dashboard/clusterclass.tsx`), is left exactly as-is. A new
`WatchClusterClasses` watcher is added purely for the new dedicated page's live list, following the
same pattern as the other three new kinds.

**Rationale**: Spec FR-004 and Acceptance Scenario 3 (User Story 4) explicitly require the existing
widget to remain, alongside (not replaced by) the new page — matching how Clusters/Machines/
MachineDeployments already have both a main-dashboard rollup and a dedicated page. Since the widget
already works and isn't being changed, there's no reason to migrate its REST-based data path to WS;
only the *new* page needs live updates.
