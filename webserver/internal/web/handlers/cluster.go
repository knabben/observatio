package handlers

import (
	"github.com/knabben/observatio/webserver/internal/infra/clusterapi"
	"net/http"
)

// handleClusterList returns the information about the cluster.
func handleClusterList(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()
	c, err := clusterapi.NewClientWithScheme(ctx, scheme)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	clusters, err := clusterapi.FetchClusters(r.Context(), c)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	if err = writeResponse(w, clusters); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
}
