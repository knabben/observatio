package models

// MachineDeploymentResponse returns the list of available machine
// deployments.
type MachineDeploymentResponse struct {
	Total              int                 `json:"total"`
	Failing            int                 `json:"failing"`
	MachineDeployments []MachineDeployment `json:"clusters"`
}

// MachineDeployment stores the details for CAPI machine deployments objects
type MachineDeployment struct {
	Name                string `json:"name"`
	Replicas            int    `json:"replicas"`
	Cluster             string `json:"cluster"`
	Created             string `json:"created"`
	ReadyReplicas       int    `json:"readyReplicas"`
	UnavailableReplicas int    `json:"unavailableReplicas"`
	UpdatedReplicas     int    `json:"updatedReplicas"`
	Phase               string `json:"phase"`
}
