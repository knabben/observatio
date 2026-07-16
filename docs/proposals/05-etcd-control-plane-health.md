# Quick specification statement: etcd/control-plane-specific health detection

The companion guide's "Level 3: Management Cluster Degraded" failure mode is specifically about etcd
quorum loss in the management cluster (e.g., two of three etcd members lost, making the API server
read-only) — a distinct and more precise condition than "the API server is unreachable or erroring,"
which is all FR-015/006 currently detects. Observātiō's watched CAPI kinds
(`webserver/internal/web/watchers/`: Cluster, ClusterClass, Machine, MachineDeployment, MachineSet,
MachineHealthCheck per 006) do not include `KubeadmControlPlane` (KCP) at all, so there's no path to
attributing a degraded management cluster to etcd quorum specifically versus any other kind of
API-server trouble.

Add a `KubeadmControlPlane` watcher, following the existing watcher pattern (mirrors
`machineset.go`/`machinehealthcheck.go` added in 006), and read KCP's status conditions
(`EtcdClusterHealthy`, control-plane replica counts vs. desired) to give the existing
management-cluster-degraded banner (FR-015) a more specific, actionable message — "etcd quorum lost (1/3
members healthy)" instead of only "API server unreachable." This directly strengthens 006's existing
Level-3 severity classification rather than introducing a new category, and should be scoped as an
enhancement to `severity.go`'s existing management-critical path.
