package models

import clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"

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

// Machine stores a machine detail.
type Machine struct {
	Name                string                 `json:"name"`
	Namespace           string                 `json:"namespace"`
	Created             string                 `json:"created"`
	Owner               string                 `json:"owner"`
	Bootstrap           string                 `json:"bootstrap"`
	Cluster             string                 `json:"cluster,omitempty"`
	NodeName            string                 `json:"nodeName,omitempty"`
	ProviderID          string                 `json:"providerID,omitempty"`
	Version             string                 `json:"version,omitempty"`
	BootstrapReady      bool                   `json:"bootstrapReady"`
	InfrastructureReady bool                   `json:"infrastructureReady"`
	Phase               clusterv1.MachinePhase `json:"phase"`
}

type MachineInfra struct {
	Name string `json:"name"`
}
