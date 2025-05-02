package processor

import (
	"github.com/knabben/observatio/webserver/internal/infra/models"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

// ProcessMachine returns the list of Machines objects from the mgmt cluster.
func ProcessMachine(machines []clusterv1.Machine) models.MachineResponse {
	var failed int
	var machinesList []models.Machine
	for _, m := range machines {
		var nodeRef string
		if m.Status.NodeRef != nil {
			nodeRef = m.Status.NodeRef.Name
		}
		var version string
		if m.Spec.Version != nil {
			version = *m.Spec.Version
		}
		machinesList = append(machinesList, models.Machine{
			Name:      m.Name,
			Namespace: m.Namespace,
			Cluster:   m.Spec.ClusterName,
			NodeName:  nodeRef,
			Version:   version,
			Phase:     clusterv1.MachinePhase(m.Status.Phase),
		})
		if !m.Status.InfrastructureReady || !m.Status.BootstrapReady {
			failed += 1
		}
	}
	return models.MachineResponse{
		Total:    len(machinesList),
		Failing:  failed,
		Machines: machinesList,
	}
}
