# Data Model: Infrastructure Provider Detection & Adaptive Listing Screens

## Cluster (existing entity, extended)

`webserver/internal/infra/models/cluster.go`

| Field | Type | Notes |
|---|---|---|
| `InfrastructureRef` | `*corev1.ObjectReference` | Existing. Source of truth for provider derivation (`.Kind`). |
| `Provider` | `string` | **NEW**. One of `docker`, `vsphere`, `unknown`. Derived via `providerkind.FromKind(InfrastructureRef.Kind)` at fetch time; not stored in the cluster's own CRD. |

## Machine (existing entity, extended)

`webserver/internal/infra/models/machine.go`

| Field | Type | Notes |
|---|---|---|
| `Provider` | `string` | **NEW**. Same derivation as Cluster, from the raw `clusterv1.Machine.Spec.InfrastructureRef.Kind` (not previously surfaced on this DTO at all). |

## InfrastructureCapability (new entity)

`webserver/internal/infra/models/capability.go`

Represents environment-level detection: which infrastructure providers are installed, and their
version, derived from the clusterctl provider inventory (`clusterctlv1.ProviderList`).

| Field | Type | Notes |
|---|---|---|
| `Docker` | `ProviderStatus` | Detection result for the Docker infrastructure provider. |
| `VSphere` | `ProviderStatus` | Detection result for the vSphere infrastructure provider. |

`ProviderStatus`:

| Field | Type | Notes |
|---|---|---|
| `Installed` | `bool` | Whether this provider's clusterctl inventory entry was found. |
| `Version` | `string` | Empty when `Installed` is `false`. |

## ClusterInfraDocker (new entity, mirrors existing ClusterInfra)

`webserver/internal/infra/models/cluster.go` (alongside existing `ClusterInfra`)

Provider-specific detail for Docker-backed clusters, populated from `dockerv1.DockerCluster`
(analogous to how `ClusterInfra` wraps `capv.VSphereClusterStatus` today).

| Field | Type | Notes |
|---|---|---|
| `Cluster` | `string` | Owning Cluster name. |
| `LoadBalancerIP` | `string` | Docker-provider equivalent of vSphere's `Server`. |
| `Status` | `dockerv1.DockerClusterStatus` | Raw upstream status, same pattern as `ClusterInfra.Status`. |

## MachineInfraDocker (new entity, mirrors existing MachineInfra)

`webserver/internal/infra/models/machine.go`

Provider-specific detail for Docker-backed machines, populated from `dockerv1.DockerMachine`.

| Field | Type | Notes |
|---|---|---|
| `ProviderID` | `string` | Existing-shape field, reused. |
| `Status` | `dockerv1.DockerMachineStatus` | Raw upstream status. |

Note: Docker machines have no direct equivalent of vSphere's `NumCPUs`/`MemoryMiB`/`DiskGiB`/
`CloneMode`/`PowerOffMode` (those are vSphere-specific virtualization concepts); the Docker view
surfaces only the fields Docker actually exposes rather than padding with empty vSphere-only columns.

## Relationships / Derivation Flow

```text
clusterv1.Cluster.Spec.InfrastructureRef.Kind ──┐
clusterv1.Machine.Spec.InfrastructureRef.Kind ──┼─► providerkind.FromKind() ─► "docker" | "vsphere" | "unknown"
                                                 │
clusterctlv1.ProviderList (Type=InfrastructureProvider) ─► InfrastructureCapability{Docker, VSphere}
```

No new persistent storage — all values are derived per-request from the connected cluster's live
API objects, consistent with the project's existing stateless model.
