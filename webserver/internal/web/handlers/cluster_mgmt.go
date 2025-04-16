package handlers

import (
	"github.com/knabben/observatio/webserver/internal/infra/clusterapi"
	"net/http"
)

// handleSummaryCluster returns the summary of clusters states.
func handleSummaryCluster(w http.ResponseWriter, r *http.Request) {
	var (
		ctx     = r.Context()
		summary clusterapi.ClusterSummary
	)

	cli, err := clusterapi.NewClientWithScheme(ctx, scheme)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	if summary, err = clusterapi.GenerateClusterSummary(ctx, cli); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	if err = writeResponse(w, summary); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
}
