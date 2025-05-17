package kubernetes

import (
	"context"
	"net/http"

	"github.com/knabben/observatio/webserver/internal/web/handlers/system"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/knabben/observatio/webserver/internal/infra/clusterapi"
	"github.com/knabben/observatio/webserver/internal/infra/clusterapi/fetchers"
	"github.com/knabben/observatio/webserver/internal/infra/models"
)

// HandleClusterList returns the information about the cluster.
func HandleClusterList(w http.ResponseWriter, r *http.Request) {
	fetchClusterData[models.ClusterResponse](r.Context(), w, fetchers.FetchClusters)
}

// HandleClusterInfraList returns the information about the vSphere cluster.
func HandleClusterInfraList(w http.ResponseWriter, r *http.Request) {
	fetchClusterData[models.ClusterInfraResponse](r.Context(), w, fetchers.FetchClustersInfra)
}

// fetchClusterData write the return of the response based on a cluster type.
func fetchClusterData[T any](ctx context.Context, w http.ResponseWriter, fetchFunc func(context.Context, client.Client) (T, error)) {
	var response T
	cli, err := clusterapi.NewClientWithScheme(ctx, system.Scheme)
	if system.HandleError(w, http.StatusInternalServerError, err) {
		return
	}
	if response, err = fetchFunc(ctx, cli); system.HandleError(w, http.StatusInternalServerError, err) {
		return
	}
	if err = system.WriteResponse(w, response); system.HandleError(w, http.StatusInternalServerError, err) {
		return
	}
}
