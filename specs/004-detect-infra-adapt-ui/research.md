# Phase 0 Research: Infrastructure Provider Detection & Adaptive Listing Screens

## R1: How to detect which infrastructure provider(s) are installed, with version

**Decision**: Reuse the existing `clusterapi.GenerateComponentVersions` function
(`webserver/internal/infra/clusterapi/dashboard.go:148-160`), which already lists
`clusterctlv1.ProviderList` (the clusterctl inventory CRs every CAPI provider installs) and returns
`ProviderName`/`Kind`/`Version` per component. Filter that list for entries whose `Kind` equals the
`InfrastructureProvider` provider type and whose `ProviderName` is `docker` or `vsphere` to build the
new `InfrastructureCapability` response.

**Rationale**: This mechanism already exists and is already exercised in production via
`/api/clusters/components` (`HandleComponentsVersion`) — it is the same inventory clusterctl itself
uses to answer "what's installed and at what version." No new dependency, no new discovery-API code
path, and it directly satisfies both FR-002 (environment-level capability) and FR-011 (provider
version) from one existing, already-tested data source.

**Alternatives considered**:
- Raw `discovery.DiscoveryClient` scan (`NewDiscoveryClient` in `client.go`, currently unused/dead
  code) for `infrastructure.cluster.x-k8s.io/v1beta1` group presence — rejected: gives CRD
  existence but not a reliable, provider-attributed version string; duplicates data already available
  via the clusterctl inventory.
- A brand-new label/annotation convention — rejected: no such convention exists in this codebase or
  upstream CAPI; would be a proprietary invention where a standard mechanism already suffices.

## R2: How to resolve which provider backs a specific Cluster/Machine

**Decision**: Add a shared `providerkind.FromKind(kind string) string` helper (`docker` for
`DockerCluster`/`DockerMachine`, `vsphere` for `VSphereCluster`/`VSphereMachine`, `unknown`
otherwise) and apply it to the `Kind` of the resource's existing `infrastructureRef` — already a
populated field on the raw `clusterv1.Cluster`/`clusterv1.Machine` objects, and already exposed on
`models.Cluster.InfrastructureRef` (`webserver/internal/infra/models/cluster.go:35`). Store the
result as a new plain `Provider string` field on `models.Cluster` and `models.Machine`.

**Rationale**: `infrastructureRef.Kind` is a standard CAPI field already present on every Cluster and
Machine — no extra API calls, no proprietary type in the core domain (Constitution III).

**Alternatives considered**:
- Fetching the actual `DockerCluster`/`VSphereCluster` object per row just to read its `Kind` —
  rejected: the owning `Cluster`'s `infrastructureRef.Kind` already carries this; an extra per-row
  fetch would be redundant and slower.

## R3: Docker infrastructure provider Go types

**Decision**: Import Docker infra types from the Docker-provider subpackage of the **already-present**
`sigs.k8s.io/cluster-api` v1.9.6 module (`test/infrastructure/docker/api/v1beta1`, providing
`DockerCluster`/`DockerMachine`) — no new `go.mod` entry, since the module is already a resolved
dependency (used today for `clusterv1`/`clusterctlv1`). Register it in the shared `runtime.Scheme`
(`webserver/internal/web/handlers/system/utils.go`) alongside `clusterctlv1`, `clusterv1`, `capv`.

**Rationale**: Keeps the new `ClusterInfraDocker`/`MachineInfraDocker` fetchers structurally identical
to the existing vSphere ones (`ListClusterInfra` → `capv.VSphereClusterList`), preserving the
established pattern instead of hand-rolling unstructured/dynamic decoding.

**Risk / validation**: This subpackage lives under `test/infrastructure/docker` (CAPI's own reference
Docker provider, "CAPD" — appropriate given the spec's assumption that Docker environments are
local/dev/test-oriented). Confirm during implementation that this subpackage resolves and compiles
cleanly against the pinned `v1.9.6` module before writing the fetcher.
**Fallback** if it does not: decode only the two fields actually needed (`status.ready`,
`spec.loadBalancerIP`) via the dynamic/unstructured client already available
(`clusterapi.NewDynamicClient`), avoiding any new dependency risk.

## R4: REST surface for the new/changed endpoints

**Decision**: Add `GET /api/infra/capabilities` returning `InfrastructureCapability`. Generalize the
existing `GET /api/clusters/infra/list` with an optional `?provider=docker|vsphere` query parameter
(defaulting to the first detected provider when omitted) instead of adding parallel
provider-specific routes.

**Rationale**: Centralizes provider dispatch behind one endpoint + one capability lookup, matching
FR-012 (server-side, single decision point) and keeping the existing vSphere consumer contract
(FR-010) intact by making `provider` additive/optional rather than a breaking rename.

**Alternatives considered**: Separate `/api/clusters/infra/docker/list` and
`/api/clusters/infra/vsphere/list` endpoints — rejected: duplicates dispatch logic on the frontend
that the capability endpoint should already centralize, and multiplies the surface area for every
future provider.

## R5: Transport for the new capability data (REST vs. WebSocket)

**Decision**: REST (`GET /api/infra/capabilities`), not WebSocket.

**Rationale**: Constitution Principle II mandates WebSocket for *live cluster state* (Cluster,
Machine, MachineDeployment health/conditions). Which infrastructure providers are installed, and at
what version, is environment/installation metadata that changes only when an operator
installs/upgrades a provider — not runtime cluster health. This matches the existing precedent of
`/api/clusters/components`, `/api/clusters/classes`, and `/api/clusters/summary`
(`webserver/internal/web/handlers/kubernetes/dashboard.go`), all of which are REST-only today for
the same class of near-static metadata.

**Alternatives considered**: Piggyback capability data onto the existing WebSocket watcher stream —
rejected: adds complexity to a live-update channel for data that essentially never changes within an
operator's session.
