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

Provider-specific detail for Docker-backed clusters, populated from the `DockerCluster` object
(`infrastructure.cluster.x-k8s.io/v1beta1`) read via the dynamic/unstructured client and decoded
field-by-field (see research.md R3 — the typed `dockerv1` package is not used; it lives in a separate
Go module that would force a major dependency upgrade).

| Field | Type | Notes |
|---|---|---|
| `Cluster` | `string` | Owning Cluster name, read from the object's owner reference. |
| `LoadBalancerIP` | `string` | Read from `spec.loadBalancerIP`; Docker-provider equivalent of vSphere's `Server`. |
| `Ready` | `bool` | Read from `status.ready` — the single field needed for the status indicator. |

## MachineInfraDocker (new entity, mirrors existing MachineInfra)

`webserver/internal/infra/models/machine.go`

Provider-specific detail for Docker-backed machines, populated from the `DockerMachine` object via
the same dynamic/unstructured decode approach as `ClusterInfraDocker`.

| Field | Type | Notes |
|---|---|---|
| `ProviderID` | `string` | Read from `spec.providerID`; existing-shape field, reused. |
| `Ready` | `bool` | Read from `status.ready`. |

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
