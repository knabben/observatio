# Quick specification statement: first-class pages for other crucial CAPI objects

Checking `front/app/ui/dashboard/nav-links.tsx` confirms only three CAPI kinds have their own browsable
list/detail pages today: Clusters, MachineDeployments, and Machines. `MachineHealthCheck` (MHC) and
`KubeadmControlPlane` (KCP, see proposal 05) are two of the most operationally important CAPI objects
for exactly the failure modes the companion guide covers — MHC defines the remediation policy (health
check timeouts, `maxUnhealthy` threshold, which Machines it targets) behind every "self-healing vs.
needs-investigation" call the dashboard makes (006/US4), and KCP carries control-plane replica health
and etcd conditions — yet an operator has no page to independently browse or inspect either one. MHC is
currently watched only as an internal signal feeding `severity.go`'s classification (added in 006), never
exposed as its own entity; KCP isn't watched at all today. `ClusterClass` and `MachineSet` are in a
similar position — watched/read for rollups but with no dedicated page either.

Give `MachineHealthCheck`, `KubeadmControlPlane`, `MachineSet`, and `ClusterClass` their own first-class
list + detail pages and per-object AI panel, mirroring the existing
Clusters/Machines/MachineDeployments pattern (list page → detail page with spec/status/conditions →
"Ask AI about this" panel from feature 005) instead of leaving them as backend-only signals that only
surface indirectly through rollups, severity, or debug-path evidence. This makes the underlying object
model fully explorable on its own terms, independent of the Day-2 Ops dashboard's derived views, and
gives operators a place to inspect MHC policy or KCP/etcd status directly when they want to understand
*why* the dashboard reached a given conclusion.
