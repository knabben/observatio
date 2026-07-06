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

## R3: Docker infrastructure provider Go types — SUPERSEDED, see below

**Decision (revised after implementation spike)**: Do **not** import
`sigs.k8s.io/cluster-api/test/infrastructure/docker/api/v1beta1`. Confirmed via `go get` during
implementation: `test/infrastructure/docker` is a **separate nested Go module**
(`sigs.k8s.io/cluster-api/test`), not just a subpackage of the already-present
`sigs.k8s.io/cluster-api` module. Importing it forces upgrading `sigs.k8s.io/cluster-api` from the
pinned `v1.9.6` to `v1.13.3`, plus a cascade of transitive bumps (`k8s.io/client-go` v0.32.1→v0.35.4,
`sigs.k8s.io/controller-runtime` v0.19.7→v0.23.3, and others) — exactly the invasive, unjustified
dependency change this plan set out to avoid.

**Revised decision**: Fetch Docker infra details (`DockerCluster`/`DockerMachine`,
`infrastructure.cluster.x-k8s.io/v1beta1`) via the existing dynamic/unstructured client
(`clusterapi.NewDynamicClient`, already used by the WebSocket watchers), decoding only the fields the
Docker infra view actually needs (`status.ready`, `status.conditions`, `spec.loadBalancerIP` for
Cluster; `status.ready`, `spec.providerID` for Machine) via
`runtime.DefaultUnstructuredConverter.FromUnstructured` into small local structs — the same
conversion pattern the watchers package already uses for typed CAPI objects. No `runtime.Scheme`
registration needed for these two kinds since they're never typed-decoded through the controller-runtime
client. `models.ClusterInfraDocker`/`models.MachineInfraDocker` (data-model.md) are unaffected — only
how they're populated changes.

**Original risk note (for record)**: this subpackage lives under `test/infrastructure/docker` (CAPI's
own reference Docker provider, "CAPD"); the module-boundary issue above was exactly the risk flagged
before implementation, and the fallback documented then is now the actual, confirmed approach.

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
