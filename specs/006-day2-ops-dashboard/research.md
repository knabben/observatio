# Phase 0 Research: Day-2 Operations Dashboard

## R1: Layers 1–3 of the debugging path require no new Kubernetes API surface

**Decision**: Synthesize the object-conditions, Machine-phase, and provider-infra-resource layers
entirely from data already fetched today.

**Rationale**: `webserver/internal/infra/models.{Cluster,Machine,MachineDeployment}.Status` already
carries the full upstream `clusterv1.*Status` (including `Conditions`) through to the frontend.
`webserver/internal/web/watchers/machine.go` already watches the provider-infra counterpart
(`DockerMachine`, `VSphereMachine`) for every Machine. `isMachineFailed`/`isClusterFailed` in
`webserver/internal/infra/clusterapi/processor/{machine,cluster}.go` already read the readiness
booleans that seed a phase-level verdict. `debugpath.go` only needs to walk conditions already in
memory and pair each Machine with its already-watched provider object — no new watch, no new RBAC.

**Alternatives considered**: Re-deriving conditions via a live `clusterctl describe`-equivalent
call at request time was rejected — it would duplicate data the backend already streams and add a
polling-shaped code path that conflicts with Principle II.

## R2: Layer 4 ("controller reconciliation activity") is sourced from Kubernetes Events, not log tailing

**Decision**: Use `corev1.Event` objects associated (via `involvedObject`) with the affected
resource as the Layer-4 signal, surfaced only when Layers 1–3 don't already explain the failure
(per FR-007).

**Rationale**: Raw controller-pod log tailing/streaming would require a new, effectively unbounded
data path (log volume, retention, multi-replica controller pods, provider-specific pod naming) that
is disproportionate for "was there recent error-level reconciliation activity for this object."
Kubernetes Events are a stable, bounded, already-scoped-to-object API that every controller
(CAPI core and every infra/bootstrap provider) already populates on reconciliation errors. This
matches the spec's own Assumption that Layer 4 data may come from "recent controller log lines **or
Kubernetes Events**."

**Alternatives considered**: Streaming `kubectl logs`-equivalent output for the relevant controller
Deployment was considered and rejected for this iteration — flagged as a possible follow-up if
Events prove insufficiently detailed in practice.

## R3: Drift detection uses `generation`/`observedGeneration`, not cloud-provider introspection

**Decision**: Flag drift when a provider-infra object's `status.observedGeneration` lags its
`metadata.generation` beyond a short grace period, or when its `Ready` condition is `True` while
`observedGeneration` is stale.

**Rationale**: True infrastructure drift detection (comparing live VM/network/storage state against
CAPI's expected state) requires cloud-provider-side, read-only introspection the Kubernetes API
alone cannot provide for every provider (vSphere, AWS, Docker each have different native
inspection surfaces). The `generation`/`observedGeneration` pair is a standard controller-runtime
idiom already present on every object already watched, and reliably signals "the controller has not
finished reconciling the current spec" — a strong, provider-agnostic proxy for drift and other
stuck-reconciliation cases, achievable with zero new API access.

**Alternatives considered**: Provider-specific drift checks (e.g., calling the Docker/vSphere API
directly to compare live VM config) were rejected as out of scope — they would violate Principle
III (no first-class provider-specific logic in the core domain) and multiply the maintenance
surface per provider. This is recorded as a best-effort heuristic in spec.md's Assumptions.

## R4: Certificate expiry reads CAPI-managed Secrets directly

**Decision**: Read the `<cluster-name>-ca`, `<cluster-name>-etcd`, and
`<cluster-name>-proxy` Secrets (standard CAPI-managed cert Secret naming) per Cluster, parse the
`tls.crt` field with the standard library's `crypto/x509`, and compare `NotAfter` against a 30-day
default warning window.

**Rationale**: These Secret names and shapes are a stable CAPI convention, not provider-specific.
`encoding/pem` + `crypto/x509` are standard library — no new dependency. Only the parsed `NotAfter`
timestamp is surfaced to the frontend; raw certificate/key bytes are never returned by the new
detector or logged (Constraints section of plan.md).

**Alternatives considered**: Relying on `kubeadm`'s own certificate-expiration reporting (requires
node/API-server-side exec) was rejected — not accessible read-only via the Kubernetes API from the
management cluster.

## R5: Stalled-rollout detection requires a new MachineSet watcher

**Decision**: Add `webserver/internal/web/watchers/machineset.go` (MachineSet is already a
first-class CAPI type; no watcher exists for it today). A rollout is flagged stalled when a
MachineDeployment has more than one active MachineSet for longer than a grace period, with the new
MachineSet's Machines `Ready` and the old MachineSet not scaling toward zero. The likely blocking
cause is reported as a PodDisruptionBudget (`policy/v1`, already vendored via `k8s.io/api`)
preventing eviction, if one is found covering the relevant workload-cluster Pods, or the presence
of an unexpected finalizer on the stuck Machine/Node, when that information is available.

**Rationale**: MachineSet is the CAPI object that owns the exact "new replica set healthy, old one
won't scale down" signal called out in the source material and in FR-009; it sits directly in the
already-compliant Cluster → MachineDeployment → Machine hierarchy (Principle III). PDB checks use
an already-vendored core Kubernetes API group.

**Alternatives considered**: Inferring rollout stalls purely from MachineDeployment-level
conditions without watching MachineSet was rejected — MachineDeployment conditions alone don't
reliably distinguish "stuck because old MachineSet won't scale down" from other stall causes.

## R6: Version-skew detection extends the existing provider-inventory read

**Decision**: Extend `GenerateComponentVersions`/`GenerateInfrastructureCapability` (already reading
`clusterctlv1.ProviderList`) with a check of each installed provider's owned CRD versions via the
`apiextensions-apiserver` clientset (already an indirect Go module dependency of the backend,
promoted to a direct one — no new external module).

**Rationale**: `clusterctlv1.ProviderList` already tells us which provider version is installed;
cross-referencing the `apiextensions.k8s.io` `CustomResourceDefinition.spec.versions` for that
provider's owned CRDs against what the running controller Deployment declares (via its image tag or
an annotation, where present) is the only read-only, provider-agnostic signal available without
executing into the controller pod. This is explicitly documented (plan.md, Constraints; spec.md,
Assumptions) as best-effort/heuristic — full confidence would require provider-published
compatibility metadata that does not uniformly exist today.

**Alternatives considered**: Requiring operators to manually confirm `clusterctl upgrade plan`
output was rejected as it defeats the purpose of proactive detection, but the plan's heuristic check
is explicitly weaker than running that command directly, and the UI must not overstate its
confidence (see quickstart.md verification notes).

## R7: Severity classification needs two new, generic (non-provider-specific) read paths

**Decision**: Add a `MachineHealthCheck` watcher (first-class CAPI type) to read
`status.remediationsAllowed`/related conditions for the "`maxUnhealthy` breached" case (FR-013), and
read Pod/Deployment status (`apps/v1`, `core/v1`, already-vendored) in well-known controller
namespaces (`capi-system` and provider namespaces such as `capd-system`, already referenced in the
domain description) to detect crash-looping provider controllers (FR-014).

**Rationale**: Both are generic Kubernetes/CAPI reads, keeping provider specifics opaque per
Principle III. `MachineHealthCheck` is the actual CAPI object that encodes the "remediation paused,
possible network partition" signal from the source material.

**Alternatives considered**: Parsing controller logs for crash-loop evidence was rejected in favor
of the structured, already-standard `CrashLoopBackOff`/`RestartCount` signals on `Pod.Status`.

## R8: Management-cluster-degraded and CA-secret-loss are approximated, not exactly diagnosed

**Decision**: "Management cluster degraded" (FR-015) is approximated from the backend's own live
API-server call success/latency/error rate (the backend already holds a persistent client
connection for all existing watchers) rather than true etcd-quorum introspection, which is not
accessible read-only through the Kubernetes API from outside the control plane. "CA secret lost"
(FR-016) is a direct existence/readability check on the `<cluster-name>-ca` Secret from R4.

**Rationale**: This is the most severe failure class in the source material precisely because it
can be undetectable from inside the cluster (etcd quorum loss manifests as API server errors, which
*is* observable read-only). Exact etcd member-count/quorum state would require direct etcd
client access or exec into a control-plane node, both out of scope for a read-only monitoring
dashboard.

**Alternatives considered**: Deferring this detection entirely (since it can't be made exact) was
rejected — an approximate signal (API server unreachable/erroring) still delivers the bulk of the
requested value (FR-015's "top-level banner... all lifecycle operations blocked") and is recorded
transparently as an approximation rather than silently overclaiming precision.

## R9: Delivery stays WebSocket-push via a new aggregator, preserving Principle II

**Decision** (revised after inspecting the actual watcher architecture — see note below): Every
existing resource watcher (`webserver/internal/web/watchers/{cluster,machine,machinedeployment}.go`)
is a *per-connection* 1:1 relay: `HandleWatcher` (`webserver/internal/web/handlers/system/websocket.go`)
upgrades one WebSocket connection, reads a single `{"type": "..."}` message, and dispatches to one
`Watch<Kind>(ctx, conn, objType)` function that opens exactly one K8s watch and streams its events
directly into that one connection for its lifetime (`WatchResourceViaWebSocket`/`streamEvents`).
There is no shared broadcast pool for resource data — `webserver/internal/web/handlers/system/pool.go`
exists solely for the AI chatbot's fan-out, unrelated to resource watching.

The Day2Ops aggregator therefore cannot "subscribe to" existing watcher streams (they're not a
shared bus). Instead, a new `WatchDay2Ops(ctx, conn, objType)` in
`webserver/internal/web/watchers/day2ops.go` is added to the same `watchHandlers` dispatch table
under a new `TypeDay2Ops` (`"day2ops"`) key. On a new connection, it opens its **own** set of K8s
watches (Cluster, Machine, MachineDeployment, MachineSet, MachineHealthCheck, provider-infra
objects) directly via the dynamic client — following the exact same per-GVR `Watch(...)` call
every existing watcher already uses — fans their events into one internal channel, maintains an
in-memory snapshot, and on each event recomputes the full `Day2OpsEvent` (via pure functions in the
`day2ops` package) and writes it to that same connection. This keeps delivery WS-push per
Principle II, using the same primitives (`clusterapi.NewDynamicClient`, `watch.Interface`) already
proven in every other watcher, just fanning in multiple GVRs instead of one. Only on-demand
deep-drill detail (e.g., the full evidence list behind one object's debugging path, fetched only
when an operator expands it) uses a scoped REST endpoint (`GET /api/day2ops/detail`), identical in
shape and justification to the raw-object REST exception already established in feature 005
(WS-triggered, on-demand REST hydration, not independent polling).

*Correction note*: plan.md's original Project Structure described this as a package-level
"aggregator.go... broadcasts over the existing WS connection pool" — that phrasing assumed a shared
pub/sub mechanism that does not exist in this codebase. The actual implementation lives partly in
`webserver/internal/web/watchers/day2ops.go` (the per-connection fan-in loop, matching the existing
watcher convention) and partly in `webserver/internal/infra/clusterapi/day2ops/` (pure compute
functions operating on the fanned-in snapshot) — not a single "aggregator" broadcasting to a pool.

**Rationale**: Keeps the dashboard's primary rollup/severity/path data fully real-time per
Principle II, while avoiding recomputing and re-pushing large per-object evidence payloads that
most operators will never expand.

**Alternatives considered**: A REST-polling summary endpoint (simplest to build) was rejected
outright as a direct Principle II violation.

## R10: The new Logs view streams controller Pod logs via the standard Kubernetes Pod-log subresource

**Decision**: The new "Logs" destination (User Story 5) retrieves logs for the CAPI-core or
infrastructure/bootstrap provider controller Deployment's Pod via the same Pod-log subresource
(`GET /api/v1/namespaces/{ns}/pods/{pod}/log`) that `kubectl logs` itself uses, wrapped in
`webserver/internal/infra/clusterapi/fetchers/controllerlogs.go` and exposed through
`GET /api/logs/controller`. Node/Machine-level log streaming (e.g., a CAPD Machine's underlying
container output) is explicitly out of scope for this feature — tracked as a follow-up TODO, not
built now (per clarification). VM-based-provider node access is limited to static SSH
connection-instruction text (command + node address from `Machine.status.addresses`); Observātiō
never stores, manages, or transmits SSH credentials.

**Rationale**: This keeps the feature entirely inside the existing dependency set — `k8s.io/client-go`
already provides the Pod-log subresource client used by `kubectl logs` itself, so no Docker Engine
API client, no Docker socket mount, and no new external module are required (unlike an earlier,
since-corrected design direction that considered streaming CAPD container logs directly via the
Docker daemon). It also directly matches "Layer 4: Controller logs" from the original feature
request (`kubectl logs -n capi-system deploy/capi-controller-manager`), rather than conflating it
with the separate, harder Layer 5 (node-level) problem.

**Alternatives considered**: Streaming the underlying Docker container's logs for CAPD Machines via
the Docker Engine API was considered and explicitly rejected per clarification — it would have
required a new external Go dependency, a new operational requirement (mounting the Docker socket
into wherever the backend runs), and conflated node-level (Layer 5) with controller-level (Layer 4)
concerns. Node-level log streaming remains a TODO for a future, separately-scoped iteration.
