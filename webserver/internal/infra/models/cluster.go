package models

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

// Cluster stores the definition of a CAPI Cluster
type Cluster struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Paused indicates whether the cluster is paused, preventing any reconciling actions from being performed.
	Paused bool `json:"paused"`

	// ClusterNetwork represents the network configuration for a CAPI Cluster.
	ClusterNetwork clusterv1.ClusterNetwork `json:"clusterNetwork"`

	// ControlPlaneEndpoint represents the API endpoint of the control plane for the cluster.
	ControlPlaneEndpoint clusterv1.APIEndpoint `json:"controlPlaneEndpoint,omitempty"`

	// ControlPlaneRef references the control plane of a cluster, specifying its configuration and location.
	ControlPlaneRef *corev1.ObjectReference `json:"controlPlaneRef,omitempty"`

	// InfrastructureRef references the infrastructure-specific cluster configuration.
	InfrastructureRef *corev1.ObjectReference `json:"infrastructureRef,omitempty"`

	// Topology defines the cluster topology, including class type, Kubernetes version, and worker/replica configurations.
	Topology ClusterClassType `json:"topology"`

	// Age represents the age of the Cluster in a human-readable format.
	Age string `json:"age"`

	Status clusterv1.ClusterStatus `json:"status"`
}

// ClusterClassType stores the topology definition for a Cluster
type ClusterClassType struct {
	// IsClusterClass indicates whether the Cluster is specified using a ClusterClass topology.
	IsClusterClass bool `json:"isClusterClass,omitempty"`

	// ClassName specifies the name of the topology/cluster class used to create the cluster.
	ClassName string `json:"className,omitempty"`

	// ClassNamespace represents the namespace of the cluster class used in the cluster's topology.
	ClassNamespace string `json:"classNamespace,omitempty"`

	// KubernetesVersion specifies the target Kubernetes version for the cluster topology.
	KubernetesVersion string `json:"kubernetesVersion,omitempty"`

	// ControlPlaneReplicas defines the number of control plane replicas for the cluster topology.
	ControlPlaneReplicas int32 `json:"controlPlaneReplicas,omitempty"`

	// ControlPlaneMHC specifies if a MachineHealthCheck is enabled for the control plane.
	ControlPlaneMHC bool `json:"controlPlaneMHC,omitempty"`

	// WorkersMachineDeployments specifies worker machine deployment topologies for the cluster.
	WorkersMachineDeployments []clusterv1.MachineDeploymentTopology `json:"machineDeployments,omitempty"`
}

// ClusterInfraResponse returns the ClusterInfra list payload and internal formatted
// list of clusters.
type ClusterInfraResponse struct {
	Total    int            `json:"total"`
	Failing  int            `json:"failing"`
	Clusters []ClusterInfra `json:"clusters"`
}

// ClusterInfra stores the definition for CAPV
type ClusterInfra struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Cluster specifies the owning cluster of the ClusterInfra object.
	Cluster string `json:"cluster"`

	// Thumbprint is a string field that stores the unique certificate thumbprint information for the server connection.
	Thumbprint string `json:"thumbprint"`

	// Age represents the creation timestamp of the ClusterInfra resource in a human-readable duration format.
	Age string `json:"age"`

	// Server specifies the URL of the vSphere server associated with the cluster.
	Server string `json:"server"`

	// ControlPlaneEndpoint represents the endpoint details for the control plane of the cluster.
	ControlPlaneEndpoint string `json:"controlPlaneEndpoint"`

	// Modules contain a list of associated ClusterModule objects for tracking and managing cluster modules.
	Modules []capv.ClusterModule `json:"modules"`

	Status capv.VSphereClusterStatus `json:"status"`
}
