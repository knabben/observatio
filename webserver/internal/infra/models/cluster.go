package models

import (
	capv "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

// ClusterResponse returns the Cluster list payload and internal formatted
// list of clusters.
type ClusterResponse struct {
	Total    int       `json:"total"`
	Failing  int       `json:"failing"`
	Clusters []Cluster `json:"clusters"`
}

// ClusterInfraResponse returns the ClusterInfra list payload and internal formatted
// list of clusters.
type ClusterInfraResponse struct {
	Total    int            `json:"total"`
	Failing  int            `json:"failing"`
	Clusters []ClusterInfra `json:"clusters"`
}

// ClusterClassType stores the topology definition for a Cluster
type ClusterClassType struct {
	IsClusterClass            bool                                  `json:"isClusterClass,omitempty"`
	ClassName                 string                                `json:"className,omitempty"`
	ClassNamespace            string                                `json:"classNamespace,omitempty"`
	KubernetesVersion         string                                `json:"kubernetesVersion,omitempty"`
	ControlPlaneReplicas      int32                                 `json:"controlPlaneReplicas,omitempty"`
	ControlPlaneMHC           bool                                  `json:"controlPlaneMHC,omitempty"`
	WorkersMachineDeployments []clusterv1.MachineDeploymentTopology `json:"machineDeployments,omitempty"`
}

// Cluster stores the definition of a CAPI Cluster
type Cluster struct {
	Name                string               `json:"name"`
	Namespace           string               `json:"namespace"`
	Paused              bool                 `json:"paused"`
	ClusterClass        ClusterClassType     `json:"clusterClass"`
	PodNetwork          string               `json:"podNetwork"`
	ServiceNetwork      string               `json:"serviceNetwork"`
	Phase               string               `json:"phase"`
	Created             string               `json:"created"`
	Conditions          clusterv1.Conditions `json:"conditions"`
	InfrastructureReady bool                 `json:"infrastructureReady"`
	ControlPlaneReady   bool                 `json:"controlPlaneReady"`
}

// ClusterInfra stores the definition for CAPV
type ClusterInfra struct {
	Name                 string               `json:"name"`
	Namespace            string               `json:"namespace"`
	Cluster              string               `json:"cluster"`
	Thumbprint           string               `json:"thumbprint"`
	Created              string               `json:"created"`
	Server               string               `json:"server"`
	ControlPlaneEndpoint string               `json:"controlPlaneEndpoint"`
	Conditions           clusterv1.Conditions `json:"conditions"`
	Modules              []capv.ClusterModule `json:"modules"`
	Ready                bool                 `json:"ready"`
}
