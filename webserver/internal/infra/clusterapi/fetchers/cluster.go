package fetchers

import (
	"context"

	capv "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/knabben/observatio/webserver/internal/infra/clusterapi/processor"
	"github.com/knabben/observatio/webserver/internal/infra/models"
)

// ClusterInput is a type constraint for clusterv1.Cluster or capv.VSphereCluster, representing input cluster types.
type ClusterInput interface {
	clusterv1.Cluster | capv.VSphereCluster
}

// ClusterResponse represents a generic response type for cluster-related data, supporting multiple implementations.
type ClusterResponse interface {
	models.ClusterResponse | models.ClusterInfraResponse
}

// ClusterFetcher defines a function type that retrieves a list of ClusterInput objects using a context and client.
type ClusterFetcher[T ClusterInput] func(ctx context.Context, c client.Client) ([]T, error)

// ProcessCluster defines a function that processes a slice of ClusterInput values and returns a ClusterResponse.
type ProcessCluster[T ClusterInput, R ClusterResponse] func(clusters []T) R

// FetchClusters retrieves and processes a list of clusters from the client, returning a formatted cluster response.
func FetchClusters(ctx context.Context, c client.Client) (response models.ClusterResponse, err error) {
	return ProcessClusters[clusterv1.Cluster, models.ClusterResponse](
		ctx, c, ListClusters, processor.ProcessClusterResponse,
	)
}

// FetchClustersInfra retrieves a structured response of vSphere cluster infrastructure data by processing cluster objects.
func FetchClustersInfra(ctx context.Context, c client.Client) (response models.ClusterInfraResponse, err error) {
	return ProcessClusters[capv.VSphereCluster, models.ClusterInfraResponse](
		ctx, c, ListClusterInfra, processor.ProcessClusterInfraResponse,
	)
}

// ProcessClusters fetches clusters using the provided fetcher and processes them into a response using the process function.
func ProcessClusters[T ClusterInput, R ClusterResponse](
	ctx context.Context, c client.Client, fetcher ClusterFetcher[T], process ProcessCluster[T, R],
) (res R, err error) {
	var clusters []T
	if clusters, err = fetcher(ctx, c); err != nil {
		return res, err
	}
	return process(clusters), nil
}

// ListClusters retrieves the list of clusters from the Kubernetes API using the provided client and context.
func ListClusters(ctx context.Context, c client.Client) (clusters []clusterv1.Cluster, err error) {
	var clusterList clusterv1.ClusterList
	if err = c.List(ctx, &clusterList); err != nil {
		return nil, err
	}
	return clusterList.Items, nil
}

// ListClusterInfra retrieves all VSphereCluster resources from the cluster and returns them as a list along with any errors.
func ListClusterInfra(ctx context.Context, c client.Client) (clusters []capv.VSphereCluster, err error) {
	var vsphereClusters capv.VSphereClusterList
	if err = c.List(ctx, &vsphereClusters); err != nil {
		return nil, err
	}
	return vsphereClusters.Items, nil
}
