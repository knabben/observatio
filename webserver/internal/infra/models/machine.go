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

// MachineInfraResponse stores the list of machines in the cluster.
type MachineInfraResponse struct {
	Total    int            `json:"total"`
	Failing  int            `json:"failing"`
	Machines []MachineInfra `json:"machines"`
}

// Machine stores a machine detail with various attributes like name,
// creation time, owner, version, and status.
type Machine struct {
	// ProcessMachine returns the list of Machines objects from the mgmt cluster.
	Name string `json:"name"`

	// Namespace stores the namespace of a Machine object, representing the scope
	// of the machine within a cluster.
	Namespace string `json:"namespace"`

	// Created represents the creation time of a Machine, stored as a string in
	// the provided timezone.
	Created string `json:"created"`

	// Owner represents the owner of a Machine instance, extracted from its
	// OwnerReferences.
	Owner string `json:"owner"`

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

	// BootstrapReady represents whether the machine bootstrap process
	// is completed and ready.
	BootstrapReady bool `json:"bootstrapReady"`

	// InfrastructureReady indicates if the infrastructure for the machine is ready.
	InfrastructureReady bool `json:"infrastructureReady"`

	// Phase represents the current phase of a Machine, indicating its
	// state in the cluster lifecycle.
	Phase clusterv1.MachinePhase `json:"phase"`
}

type MachineInfraVCenter struct {
	// Template is a string field representing a template for MachineInfra.
	Template string `json:"template,omitempty"`

	// Server is the IP address or FQDN of the vSphere server on which
	Server string `json:"server,omitempty"`

	// Thumbprint is toohe colon-separated SHA-1 checksum of the given vCenter server's host certificate
	Thumbprint string `json:"thumbprint,omitempty"`

	// Folder is the name, inventory path, managed object reference or the managed
	// object ID of the folder in which the virtual machine is created/located.
	// +optional
	Folder string `json:"folder,omitempty"`

	// Datastore is the name, inventory path, managed object reference or the managed
	// object ID of the datastore in which the virtual machine is created/located.
	// +optional
	Datastore string `json:"datastore,omitempty"`

	// StoragePolicyName of the storage policy to use with this
	// Virtual Machine
	// +optional
	StoragePolicyName string `json:"storagePolicyName,omitempty"`

	// ResourcePool is the name, inventory path, managed object reference or the managed
	// object ID in which the virtual machine is created/located.
	// +optional
	ResourcePool string `json:"resourcePool,omitempty"`

	// Network is the network configuration for this machine's VM.
	Network capv.NetworkSpec `json:"network"`
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
