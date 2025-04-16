package handlers

import (
	"net/http"

	"github.com/knabben/observatio/webserver/internal/infra/clusterapi"
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

// handleClusterClass returns the available cluster classes in the mgmt cluster.
func handleClusterClasses(w http.ResponseWriter, r *http.Request) {
	var clusterClasses []clusterapi.ClusterClass

	cli, err := clusterapi.NewClientWithScheme(r.Context(), scheme)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	if clusterClasses, err = clusterapi.FetchClusterClass(r.Context(), cli); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	if err = writeResponse(w, clusterClasses); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
}
