package handlers

import (
	"github.com/knabben/observatio/webserver/internal/infra/clusterapi"
	"net/http"
)

// handleMachine returns the information about the cluster.
func handleMachines(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()
	c, err := clusterapi.NewClientWithScheme(ctx, scheme)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	machines, err := clusterapi.FetchMachine(r.Context(), c)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	if err = writeResponse(w, machines); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
}
