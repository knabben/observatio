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
		var machineOwner string
		for _, owner := range m.OwnerReferences {
			machineOwner = owner.Name
		}
		var providerId string
		if m.Spec.ProviderID != nil {
			providerId = *m.Spec.ProviderID
		}
		machinesList = append(machinesList, models.Machine{
			Name:                m.Name,
			Namespace:           m.Namespace,
			Owner:               machineOwner,
			Cluster:             m.Spec.ClusterName,
			NodeName:            nodeRef,
			ProviderID:          providerId,
			Version:             version,
			BootstrapReady:      m.Status.BootstrapReady,
			InfrastructureReady: m.Status.InfrastructureReady,
			Phase:               clusterv1.MachinePhase(m.Status.Phase),
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
