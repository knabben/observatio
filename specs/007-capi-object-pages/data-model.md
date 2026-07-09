# Phase 1 Data Model: First-Class Pages for MachineHealthCheck, KubeadmControlPlane, MachineSet, and ClusterClass

Each entity below is a new `models.*` DTO, built the same way every existing one is: embed
`metav1.ObjectMeta`, add an `Age` string (via the existing `processor.formatDuration` helper), pass
through the upstream `Status` struct verbatim, and surface only the Spec fields the Specification
tab actually needs (mirroring `models.MachineDeployment`'s shape from 006/005).

## MachineHealthCheck

| Field | Type | Notes |
|-------|------|-------|
| `ObjectMeta` | `metav1.ObjectMeta` | embedded |
| `Age` | `string` | |
| `Cluster` | `string` | from `Spec.ClusterName` |
| `Selector` | `metav1.LabelSelector` | from `Spec.Selector` — the target-Machine selector (FR-006) |
| `MaxUnhealthy` | `string` | from `Spec.MaxUnhealthy` (stringified `intstr.IntOrString`), empty if unset |
| `NodeStartupTimeout` | `string` | formatted duration, from `Spec.NodeStartupTimeout` |
| `UnhealthyConditions` | `[]clusterv1.UnhealthyCondition` | pass-through, rendered as a small table (FR-006) |
| `Status` | `clusterv1.MachineHealthCheckStatus` | `ExpectedMachines`/`CurrentHealthy`/`RemediationsAllowed`/`Conditions` (already used internally by 006's `severity.go`) |

## KubeadmControlPlane

| Field | Type | Notes |
|-------|------|-------|
| `ObjectMeta` | `metav1.ObjectMeta` | embedded |
| `Age` | `string` | |
| `Cluster` | `string` | derived from owner reference or a `cluster.x-k8s.io/cluster-name` label, matching the existing Machine/MachineDeployment convention |
| `Replicas` | `*int32` | from `Spec.Replicas` (desired) |
| `Version` | `string` | from `Spec.Version` (Kubernetes version) |
| `Status` | `controlplanev1.KubeadmControlPlaneStatus` | pass-through — `Replicas`/`ReadyReplicas`/`UpdatedReplicas`/`Conditions` (etcd-related conditions surface here when present, FR-007) |

## MachineSet

| Field | Type | Notes |
|-------|------|-------|
| `ObjectMeta` | `metav1.ObjectMeta` | embedded |
| `Age` | `string` | |
| `Cluster` | `string` | from `Spec.ClusterName` |
| `MachineDeployment` | `string` | from the `cluster.x-k8s.io/deployment-name` label (same convention 006's `machineSetsFor` already reads) |
| `Replicas` | `*int32` | from `Spec.Replicas` (desired) |
| `Status` | `clusterv1.MachineSetStatus` | pass-through — `Replicas`/`ReadyReplicas`/`AvailableReplicas`/`Conditions` (FR-008) |

## ClusterClass

Reuses the existing `models.ClusterClass` type unchanged (already defined for the main-dashboard
widget) — no new fields needed for FR-009's Specification tab, which shows the same status/reference
data the widget already has access to.

## Frontend type mapping

Each backend DTO gets a matching TypeScript interface colocated with its `lister.tsx` (or a small
`types.ts`), following the exact convention already used by `front/app/ui/dashboard/components/
machines/`. `front/app/lib/resource-gvr.ts`'s `RESOURCE_GVR` gains four new entries:
`machineHealthCheck`, `kubeadmControlPlane`, `machineSet`, `clusterClass`, each
`{group, version, resource}` matching the corresponding backend GVR (contracts/watch-types.md).
