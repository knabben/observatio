# Quickstart: First-Class Pages for MachineHealthCheck, KubeadmControlPlane, MachineSet, and ClusterClass

Manual verification, one section per user story, against a `kind-capi-mgmt`-style test cluster.

## US1 — MachineHealthCheck

1. Ensure at least one MachineHealthCheck exists (most `clusterctl` quickstart templates create one
   for worker nodes by default).
2. Open the new "Machine Health Checks" nav entry; confirm it lists live, matching the existing
   Machines page's look/behavior (search box, table).
3. Select one; confirm the Specification tab shows its target selector, `maxUnhealthy` threshold,
   node-startup/unhealthy-condition timeouts, and current expected-vs-healthy Machine counts.
4. Confirm the YAML tab shows the complete raw object, and "Ask AI about this" opens the global AI
   panel pre-loaded with this MachineHealthCheck's identity.
5. Kill a worker node's kubelet (or otherwise make a targeted Machine unhealthy) and confirm the
   expected-vs-healthy counts update live without a page refresh.

## US2 — KubeadmControlPlane

1. Open the new "Kubeadm Control Planes" nav entry; confirm the management cluster's own KCP object
   (or a workload cluster's, if watched cluster-wide) appears.
2. Select it; confirm desired vs. ready replica counts and status conditions are shown, including
   any etcd-related condition if present.
3. On a cluster that does not use KubeadmControlPlane (or if the CRD is absent), confirm the page
   shows "Unable to load..."/empty state rather than crashing the rest of the dashboard.

## US3 — MachineSet

1. Trigger a MachineDeployment rollout (or use an existing one) so at least one MachineSet exists.
2. Open the new "Machine Sets" nav entry; confirm replica counts and the owning MachineDeployment
   name are visible per row.
3. Select one; confirm replicas/ready/available counts and conditions match what the Day-2 Ops
   dashboard's stalled-rollout warning (006/US3) already reports for the same MachineSet.

## US4 — ClusterClass

1. Confirm the existing ClusterClass widget still appears on the main `/dashboard` page, unchanged.
2. Open the new "Cluster Classes" nav entry; confirm the same ClusterClasses appear, now with live
   updates, a detail screen, YAML tab, and "Ask AI about this."

## Cross-cutting checks

- Confirm all four new pages appear in the lateral navigation in a sensible position (grouped near
  the existing Machines/Machine Deployments/Clusters entries).
- Confirm none of the four new pages introduce any write/mutate action (no scale, delete, or edit
  controls anywhere) — this feature is strictly observability, per plan.md's Constraints.
- Confirm opening any of the four new pages does not regress the existing Clusters/Machines/Machine
  Deployments/Logs pages or the 006 Day-2 Ops dashboard (full regression pass: `make
  run-tests-backend`, `make run-tests-frontend`, `make build`).
