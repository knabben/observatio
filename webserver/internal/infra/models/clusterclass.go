package models

import clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"

// ClusterClassResponse represents the response containing metadata and a list of cluster classes.
type ClusterClassResponse struct {
	Total          int            `json:"total"`
	Failing        int            `json:"failing"`
	ClusterClasses []ClusterClass `json:"clusterClasses"`
}

// ClusterClass represents a Kubernetes ClusterClass definition including its metadata and status conditions.
type ClusterClass struct {
	Name       string                `json:"name"`
	Namespace  string                `json:"namespace"`
	Generation int64                 `json:"generation"`
	Conditions []clusterv1.Condition `json:"conditions"`
}
