package models

import clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"

// MachineDeploymentResponse returns the list of available machine
// deployments.
type MachineDeploymentResponse struct {
	Total              int                 `json:"total"`
	Failing            int                 `json:"failing"`
	MachineDeployments []MachineDeployment `json:"machineDeployments"`
}

// MachineDeployment stores the details for CAPI machine deployments objects
type MachineDeployment struct {
	Name                string                           `json:"name"`
	Replicas            int32                            `json:"replicas"`
	Cluster             string                           `json:"cluster"`
	Created             string                           `json:"created"`
	ReadyReplicas       int32                            `json:"readyReplicas"`
	UnavailableReplicas int32                            `json:"unavailableReplicas"`
	UpdatedReplicas     int32                            `json:"updatedReplicas"`
	Phase               clusterv1.MachineDeploymentPhase `json:"phase"`
}
