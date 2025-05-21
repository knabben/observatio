package processor

import (
	"time"

	"github.com/knabben/observatio/webserver/internal/infra/models"
	capv "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

// ProcessMachine returns a Machine object with details extracted from the provided Machine object, including
// machine name, namespace, owner, cluster name, node name, provider ID, version, bootstrap configuration,
// bootstrap and infrastructure readiness, creation time, and machine phase.
func ProcessMachine(m clusterv1.Machine) models.Machine {
	var nodeRef string
	if m.Status.NodeRef != nil {
		nodeRef = m.Status.NodeRef.Name
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
		ObjectMeta: m.ObjectMeta,
		Age:        formatDuration(time.Since(m.CreationTimestamp.Time)),
		Cluster:    m.Spec.ClusterName,
		NodeName:   nodeRef,
		ProviderID: providerId,
		Version:    stringPointer(m.Spec.Version),
		Bootstrap:  bootstrap,
		Status:     m.Status,
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

// isMachineFailed checks if the machine has failed by evaluating infrastructure and bootstrap readiness.
func isMachineFailed(machine clusterv1.Machine) bool {
	return !machine.Status.InfrastructureReady || !machine.Status.BootstrapReady
}

// ProcessMachineInfra maps a capv.VSphereMachine object to a models.MachineInfra object
// extracting relevant infrastructure details such as name, namespace, provider ID, memory,
// disk, failure information, etc.
func ProcessMachineInfra(machine capv.VSphereMachine) models.MachineInfra {
	return models.MachineInfra{
		ObjectMeta:        machine.ObjectMeta,
		Age:               formatDuration(time.Since(machine.CreationTimestamp.Time)),
		ProviderID:        stringPointer(machine.Spec.ProviderID),
		FailureDomain:     stringPointer(machine.Spec.FailureDomain),
		PowerOffMode:      machine.Spec.PowerOffMode,
		Template:          machine.Spec.Template,
		CloneMode:         machine.Spec.CloneMode,
		NumCPUs:           machine.Spec.NumCPUs,
		NumCoresPerSocket: machine.Spec.NumCoresPerSocket,
		MemoryMiB:         machine.Spec.MemoryMiB,
		DiskGiB:           machine.Spec.DiskGiB,
		Status:            machine.Status,
	}
}

// ProcessMachineInfraResponse processes a list of VSphereMachines and returns MachineInfraResponse.
// It processes each machine's infrastructure details and tracks the total count of machines,
// the number of failing machines, and their individual infrastructure information.
func ProcessMachineInfraResponse(machines []capv.VSphereMachine) models.MachineInfraResponse {
	machineList := make([]models.MachineInfra, 0, len(machines))
	failedMachinesCount := countFailedMachines(machines)
	for _, machine := range machines {
		machineList = append(machineList, ProcessMachineInfra(machine))
	}
	return models.MachineInfraResponse{
		Total:    len(machines),
		Machines: machineList,
		Failing:  failedMachinesCount,
	}
}

// countFailedMachines returns the number of machines that are not in ready state.
func countFailedMachines(machines []capv.VSphereMachine) int {
	count := 0
	for _, machine := range machines {
		if !machine.Status.Ready {
			count++
		}
	}
	return count
}

// stringPointer returns an empty string if the input pointer is nil, otherwise returns the dereferenced string.
func stringPointer(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
