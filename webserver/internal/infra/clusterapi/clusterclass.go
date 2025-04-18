package clusterapi

import (
	"context"

	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ClusterClass struct {
	Name       string                `json:"name"`
	Namespace  string                `json:"namespace"`
	Generation int64                 `json:"generation"`
	Conditions []clusterv1.Condition `json:"conditions"`
}

// FetchClusterClass returns the list of ClusterClass objects in the mgmt cluster.
func FetchClusterClass(ctx context.Context, c client.Client) (clusterClasses []ClusterClass, err error) {
	var clusterClassList clusterv1.ClusterClassList
	if err = c.List(ctx, &clusterClassList); err != nil {
		return clusterClasses, err
	}

	for _, class := range clusterClassList.Items {
		clusterClasses = append(clusterClasses, ClusterClass{
			Name:       class.Name,
			Namespace:  class.Namespace,
			Generation: class.Status.ObservedGeneration,
			Conditions: class.Status.Conditions,
		})
	}
	return clusterClasses, nil
}
