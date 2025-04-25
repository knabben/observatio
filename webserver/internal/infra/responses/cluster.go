package responses

import clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"

// ClusterResponse returns the Cluster list payload and format
type ClusterResponse struct {
	Total   int `json:"total"`
	Failing int `json:"failing"`

	Clusters []Cluster `json:"clusters"`
}

type Cluster struct {
	Name           string `json:"name"`
	IsClusterClass bool   `json:"isClusterClass"`

	Conditions clusterv1.Conditions `json:"conditions"`
}
