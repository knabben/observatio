package fetchers

import (
	"context"

	capv "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/knabben/observatio/webserver/internal/infra/clusterapi/processor"
	"github.com/knabben/observatio/webserver/internal/infra/models"
)

type input interface {
	clusterv1.Cluster | capv.VSphereCluster
}

type response interface {
	models.ClusterResponse | models.ClusterInfraResponse
}

type ClusterFetcher[T input] func(ctx context.Context, c client.Client) ([]T, error)
type ProcessCluster[T input, R response] func(clusters []T) R

// FetchClusters returns and processes the cluster list response
func FetchClusters(ctx context.Context, c client.Client) (response models.ClusterResponse, err error) {
	return ProcessClusters[clusterv1.Cluster, models.ClusterResponse](
		ctx, c, ListClusters, processor.ProcessClusterResponse,
	)
}

// FetchClustersInfra returns and process the capv cluster list response
func FetchClustersInfra(ctx context.Context, c client.Client) (response models.ClusterInfraResponse, err error) {
	return ProcessClusters[capv.VSphereCluster, models.ClusterInfraResponse](
		ctx, c, ListClusterInfra, processor.ProcessClusterInfraResponse,
	)
}

// ProcessClusters is a generic function that fetches clusters and processes them
func ProcessClusters[T input, R response](
	ctx context.Context, c client.Client, fetcher ClusterFetcher[T], process ProcessCluster[T, R],
) (res R, err error) {
	var clusters []T
	if clusters, err = fetcher(ctx, c); err != nil {
		return res, err
	}
	return process(clusters), nil
}

func ListClusters(ctx context.Context, c client.Client) (clusters []clusterv1.Cluster, err error) {
	var clusterList clusterv1.ClusterList
	if err = c.List(ctx, &clusterList); err != nil {
		return nil, err
	}
	return clusterList.Items, nil
}

func ListClusterInfra(ctx context.Context, c client.Client) (clusters []capv.VSphereCluster, err error) {
	var vsphereClusters capv.VSphereClusterList
	if err = c.List(ctx, &vsphereClusters); err != nil {
		return nil, err
	}
	return vsphereClusters.Items, nil
}
