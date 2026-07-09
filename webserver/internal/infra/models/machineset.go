package models

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

// MachineSet stores the details for a CAPI MachineSet — replica counts, owning
// MachineDeployment, and status conditions behind the Day-2 Ops dashboard's stalled-rollout
// warning (006).
type MachineSet struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Age represents the creation time of the MachineSet, formatted as a duration string.
	Age string `json:"age"`

	// Cluster represents the name of the cluster this MachineSet belongs to.
	Cluster string `json:"cluster"`

	// MachineDeployment is the name of the owning MachineDeployment, read from the
	// cluster.x-k8s.io/deployment-name label. Empty for a standalone MachineSet.
	MachineDeployment string `json:"machineDeployment,omitempty"`

	// Replicas is the desired number of Machines. Nil when unset.
	Replicas *int32 `json:"replicas,omitempty"`

	// Status represents the current status of the MachineSet (replica counts and conditions).
	Status clusterv1.MachineSetStatus `json:"status"`
}
