package day2ops

import (
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

// ClusterFailed reports whether a Cluster's readiness booleans indicate failure, matching the
// same semantics already used by processor.ProcessClusterResponse.
func ClusterFailed(cl clusterv1.Cluster) bool {
	return !cl.Status.InfrastructureReady || !cl.Status.ControlPlaneReady
}

// MachineFailed reports whether a Machine's readiness booleans indicate failure, matching the
// same semantics already used by processor.ProcessMachineResponse.
func MachineFailed(m clusterv1.Machine) bool {
	return !m.Status.InfrastructureReady || !m.Status.BootstrapReady
}

// MachineDeploymentFailed reports whether a MachineDeployment's replica counts indicate failure,
// matching the same semantics already used by processor.ProcessMachineDeploymentResponse.
func MachineDeploymentFailed(md clusterv1.MachineDeployment) bool {
	return md.Status.ReadyReplicas != md.Status.Replicas
}

// ComputeRollups produces the per-category HealthRollup[] shown on the Day-2 Ops landing screen
// (FR-002). Degraded stays 0 here — it is populated once risk warnings and debugging-path
// evidence are folded in by the watcher (see day2ops.go's assembleData).
func ComputeRollups(clusters []clusterv1.Cluster, machineDeployments []clusterv1.MachineDeployment, machines []clusterv1.Machine) []HealthRollup {
	clusterFailed := 0
	for _, cl := range clusters {
		if ClusterFailed(cl) {
			clusterFailed++
		}
	}

	mdFailed := 0
	for _, md := range machineDeployments {
		if MachineDeploymentFailed(md) {
			mdFailed++
		}
	}

	machineFailed := 0
	for _, m := range machines {
		if MachineFailed(m) {
			machineFailed++
		}
	}

	return []HealthRollup{
		{
			Category: CategoryCluster,
			Healthy:  len(clusters) - clusterFailed,
			Failed:   clusterFailed,
		},
		{
			Category: CategoryMachineDeployment,
			Healthy:  len(machineDeployments) - mdFailed,
			Failed:   mdFailed,
		},
		{
			Category: CategoryMachine,
			Healthy:  len(machines) - machineFailed,
			Failed:   machineFailed,
		},
	}
}
