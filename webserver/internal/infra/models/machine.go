package models

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	capv "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

// MachineResponse stores the list of machines in the cluster.
type MachineResponse struct {
	Total    int       `json:"total"`
	Failing  int       `json:"failing"`
	Machines []Machine `json:"machines"`
}

// Machine stores a machine detail with various attributes like name,
// creation time, owner, version, and status.
type Machine struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Age represents the creation time of a Machine, stored as a string in
	// the provided timezone.
	Age string `json:"age,omitempty"`

	// ProcessMachine returns the list of Machines objects from the mgmt cluster.
	Bootstrap string `json:"bootstrap"`

	// Cluster is a field in the Machine struct that represents the name of
	// the cluster the machine belongs to.
	Cluster string `json:"cluster,omitempty"`

	// NodeName is a field in the Machine struct that represents the name
	// of the node associated with the machine.
	NodeName string `json:"nodeName,omitempty"`

	// ProviderID represents the unique identifier provided by the
	// infrastructure provider for the Machine.
	ProviderID string `json:"providerID,omitempty"`

	// ProcessMachine returns the list of Machines objects from the mgmt cluster.
	Version string `json:"version,omitempty"`

	// Status represents the current status details of a machine's infrastructure.
	Status clusterv1.MachineStatus `json:"status"`
}

// MachineInfraResponse stores the list of machines in the cluster.
type MachineInfraResponse struct {
	Total    int            `json:"total"`
	Failing  int            `json:"failing"`
	Machines []MachineInfra `json:"machines"`
}

// MachineInfra represents the infra details of a machine including its name,
// namespace, memory, disk, and failure information.
type MachineInfra struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Age represents the time passed until the object is created.
	Age string `json:"age,omitempty"`

	// ProviderID represents the unique identifier of the machine provider.
	ProviderID string `json:"providerID,omitempty"`

	// Template represents the template used to create the machine.
	Template string `json:"template,omitempty"`

	// FailureDomain represents the failure domain of a machine infrastructure.
	FailureDomain string `json:"failureDomain,omitempty"`

	// PowerOffMode represents the power off mode of a MachineInfra instance.
	PowerOffMode capv.VirtualMachinePowerOpMode `json:"powerOffMode,omitempty"`

	// CloneMode represents the mode for cloning in the machine infrastructure.
	CloneMode capv.CloneMode `json:"cloneMode,omitempty"`

	// NumCPUs represents the number of CPUs of a machine in the infrastructure.
	NumCPUs int32 `json:"numCPUs,omitempty"`

	// NumCoresPerSocket represents the number of CPU cores per socket in
	// the machine infrastructure.
	NumCoresPerSocket int32 `json:"numCoresPerSocket,omitempty"`

	// MemoryMiB represents the amount of memory in mebibytes
	// allocated to a machine instance.
	MemoryMiB int64 `json:"memoryMiB,omitempty"`

	// DiskGiB represents the size of disk in Gibibytes for a machine's
	// infrastructure.
	DiskGiB int32 `json:"diskGiB,omitempty"`

	// Status represents the current status details of a machine's infrastructure.
	Status capv.VSphereMachineStatus `json:"status"`
}
