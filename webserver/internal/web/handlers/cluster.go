package handlers

import (
	"github.com/knabben/observatio/webserver/internal/infra/clusterapi"
	"net/http"
)

// handleClusterList returns the information about the cluster.
func handleClusterList(w http.ResponseWriter, r *http.Request) {
	client, err := clusterapi.NewClientWithScheme(r.Context(), scheme)
	if handleError(w, http.StatusInternalServerError, err) {
		return
	}

	var clusterResponse models.ClusterResponse
	clusterResponse, err = clusterapi.FetchClusters(r.Context(), client)
	if handleError(w, http.StatusInternalServerError, err) {
		return
	}

	if err = writeResponse(w, clusterResponse); handleError(w, http.StatusInternalServerError, err) {
		return
	}
}
