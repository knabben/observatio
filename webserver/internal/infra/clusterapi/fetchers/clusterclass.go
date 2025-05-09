package fetchers

import (
	"context"

	"github.com/knabben/observatio/webserver/internal/infra/clusterapi/processor"
	"github.com/knabben/observatio/webserver/internal/infra/models"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// FetchClusterClass retrieves a list of ClusterClass objects and processes them into a ClusterClassResponse structure.
func FetchClusterClass(ctx context.Context, c client.Client) (response models.ClusterClassResponse, err error) {
	var clusterClasses []clusterv1.ClusterClass
	if clusterClasses, err = ListClusterClasses(ctx, c); err != nil {
		return response, err
	}
	return processor.ProcessClusterClassResponse(clusterClasses), nil
}

// ListClusterClasses retrieves a list of ClusterClass objects from the Kubernetes cluster using the provided client.
// It returns the list of ClusterClasses and any error encountered during the operation.
func ListClusterClasses(ctx context.Context, c client.Client) (clusterClasses []clusterv1.ClusterClass, err error) {
	var clusterClassList clusterv1.ClusterClassList
	if err = c.List(ctx, &clusterClassList); err != nil {
		return clusterClasses, err
	}
	return clusterClassList.Items, nil
}
