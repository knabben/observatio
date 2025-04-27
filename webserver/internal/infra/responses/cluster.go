package responses

import clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"

// ClusterResponse returns the Cluster list payload and internal formatted
// list of clusters.
type ClusterResponse struct {
	Total    int       `json:"total"`
	Failing  int       `json:"failing"`
	Clusters []Cluster `json:"clusters"`
}

// ClusterClass stores the topology definition for a Cluster
type ClusterClass struct {
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
	Paused              bool                 `json:"paused"`
	ClusterClass        ClusterClass         `json:"clusterClass"`
	PodNetwork          string               `json:"podNetwork"`
	ServiceNetwork      string               `json:"serviceNetwork"`
	Phase               string               `json:"phase"`
	Created             string               `json:"created"`
	Conditions          clusterv1.Conditions `json:"conditions"`
	InfrastructureReady bool                 `json:"infrastructureReady"`
	ControlPlaneReady   bool                 `json:"controlPlaneReady"`
}
