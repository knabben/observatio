package handlers

import (
	"net/http"

	"github.com/knabben/observatio/webserver/internal/infra/clusterapi"
	"github.com/knabben/observatio/webserver/internal/infra/clusterapi/fetchers"
)

// handleMachines returns the information about the cluster.
func handleMachines(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()
	c, err := clusterapi.NewClientWithScheme(ctx, scheme)
	if handleError(w, http.StatusInternalServerError, err) {
		return
	}

	machines, err := fetchers.FetchMachines(r.Context(), c)
	if handleError(w, http.StatusInternalServerError, err) {
		return
	}

	if err = writeResponse(w, machines); handleError(w, http.StatusInternalServerError, err) {
		return
	}
}

// handleMachineInfra returns the information about machine infrastructures in the cluster.
func handleMachineInfra(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()
	c, err := clusterapi.NewClientWithScheme(ctx, scheme)
	if handleError(w, http.StatusInternalServerError, err) {
		return
	}

	machinesInfra, err := fetchers.FetchMachineInfra(r.Context(), c)
	if handleError(w, http.StatusInternalServerError, err) {
		return
	}

	if err = writeResponse(w, machinesInfra); handleError(w, http.StatusInternalServerError, err) {
		return
	}
}
