package handlers

import (
	"net/http"

	"github.com/knabben/observatio/webserver/internal/infra/clusterapi"
	"github.com/knabben/observatio/webserver/internal/infra/clusterapi/fetchers"
)

// handleMachinesDeployment returns the information about the machines deployments
func handleMachineDeployments(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()
	c, err := clusterapi.NewClientWithScheme(ctx, scheme)
	if handleError(w, http.StatusInternalServerError, err) {
		return
	}

	mds, err := fetchers.FetchMachineDeployment(r.Context(), c)
	if handleError(w, http.StatusInternalServerError, err) {
		return
	}

	if err = writeResponse(w, mds); handleError(w, http.StatusInternalServerError, err) {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
}
