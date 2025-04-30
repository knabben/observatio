package handlers

import (
	"context"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/knabben/observatio/webserver/internal/infra/clusterapi"
	"github.com/knabben/observatio/webserver/internal/infra/clusterapi/fetchers"
	"github.com/knabben/observatio/webserver/internal/infra/models"
)

// handleClusterList returns the information about the cluster.
func handleClusterList(w http.ResponseWriter, r *http.Request) {
	fetchClusterData[models.ClusterResponse](r.Context(), w, fetchers.FetchClusters)
}

// handleClusterInfraList returns the information about the vSphere cluster.
func handleClusterInfraList(w http.ResponseWriter, r *http.Request) {
	fetchClusterData[models.ClusterInfraResponse](r.Context(), w, fetchers.FetchClustersInfra)
}

// fetchClusterData write the return of the response based on cluster type.
func fetchClusterData[T any](ctx context.Context, w http.ResponseWriter, fetchFunc func(context.Context, client.Client) (T, error)) {
	var response T
	cli, err := clusterapi.NewClientWithScheme(ctx, scheme)
	if handleError(w, http.StatusInternalServerError, err) {
		return
	}
	response, err = fetchFunc(ctx, cli)
	if handleError(w, http.StatusInternalServerError, err) {
		return
	}
	if err = writeResponse(w, response); handleError(w, http.StatusInternalServerError, err) {
		return
	}
}
