package handlers

import (
	"github.com/knabben/observatio/webserver/internal/infra/clusterapi"
	"net/http"
)

// handleMachinesDeployment returns the information about the machines deployments
func handleMachineDeployments(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()
	c, err := clusterapi.NewClientWithScheme(ctx, scheme)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	mds, err := clusterapi.FetchMachineDeployments(r.Context(), c)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	if err = writeResponse(w, mds); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
}
