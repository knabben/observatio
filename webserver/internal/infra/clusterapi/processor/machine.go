package processor

import (
	"time"

	"github.com/knabben/observatio/webserver/internal/infra/models"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

// ProcessMachine returns the list of Machines objects from the mgmt cluster.
func ProcessMachine(m clusterv1.Machine) models.Machine {
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
	var bootstrap string
	if m.Spec.Bootstrap.ConfigRef != nil {
		bootstrap = m.Spec.Bootstrap.ConfigRef.Name
	}
	return models.Machine{
		Name:                m.Name,
		Namespace:           m.Namespace,
		Owner:               machineOwner,
		Cluster:             m.Spec.ClusterName,
		NodeName:            nodeRef,
		ProviderID:          providerId,
		Version:             version,
		BootstrapReady:      m.Status.BootstrapReady,
		InfrastructureReady: m.Status.InfrastructureReady,
		Created:             formatDuration(time.Since(m.ObjectMeta.CreationTimestamp.Time)),
		Bootstrap:           bootstrap,
		Phase:               clusterv1.MachinePhase(m.Status.Phase),
	}
}

// ProcessMachineResponse processes a slice of cluster machines and returns a summarized response
// containing the total count, number of failing machines, and detailed machine information.
func ProcessMachineResponse(machines []clusterv1.Machine) models.MachineResponse {
	failedMachinesCount := 0
	machinesList := make([]models.Machine, 0, len(machines))

	for _, machine := range machines {
		machinesList = append(machinesList, ProcessMachine(machine))
		if isMachineFailed(machine) {
			failedMachinesCount++
		}
	}

	return models.MachineResponse{
		Total:    len(machinesList),
		Failing:  failedMachinesCount,
		Machines: machinesList,
	}
}

func isMachineFailed(machine clusterv1.Machine) bool {
	return !machine.Status.InfrastructureReady || !machine.Status.BootstrapReady
}
