package models

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	controlplanev1 "sigs.k8s.io/cluster-api/controlplane/kubeadm/api/v1beta1"
)

// KubeadmControlPlane stores the details for a CAPI KubeadmControlPlane — control-plane replica
// health and status conditions (including etcd-related conditions when present).
type KubeadmControlPlane struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Age represents the creation time of the KubeadmControlPlane, formatted as a duration string.
	Age string `json:"age"`

	// Cluster represents the name of the cluster this KubeadmControlPlane belongs to, read from
	// the cluster.x-k8s.io/cluster-name label (KubeadmControlPlaneSpec has no ClusterName field).
	Cluster string `json:"cluster"`

	// Version is the desired Kubernetes version for the control plane.
	Version string `json:"version"`

	// Replicas is the desired number of control plane machines. Nil when unset.
	Replicas *int32 `json:"replicas,omitempty"`

	// Status represents the current status of the control plane (replica counts, initialized/ready
	// flags, and conditions).
	Status controlplanev1.KubeadmControlPlaneStatus `json:"status"`
}
