package kubernetes

import (
	"net/http"

	"github.com/knabben/observatio/webserver/internal/infra/clusterapi"
	"github.com/knabben/observatio/webserver/internal/infra/clusterapi/fetchers"
	"github.com/knabben/observatio/webserver/internal/web/handlers/system"
)

// HandleMachines returns the information about the cluster.
func HandleMachines(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()
	c, err := clusterapi.NewClientWithScheme(ctx, system.Scheme)
	if system.HandleError(w, http.StatusInternalServerError, err) {
		return
	}

	machines, err := fetchers.FetchMachines(r.Context(), c)
	if system.HandleError(w, http.StatusInternalServerError, err) {
		return
	}

	if err = system.WriteResponse(w, machines); system.HandleError(w, http.StatusInternalServerError, err) {
		return
	}
}

// HandleMachineInfra returns the information about machine infrastructures in the cluster.
func HandleMachineInfra(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()
	c, err := clusterapi.NewClientWithScheme(ctx, system.Scheme)
	if system.HandleError(w, http.StatusInternalServerError, err) {
		return
	}

	machinesInfra, err := fetchers.FetchMachineInfra(r.Context(), c)
	if system.HandleError(w, http.StatusInternalServerError, err) {
		return
	}

	if err = system.WriteResponse(w, machinesInfra); system.HandleError(w, http.StatusInternalServerError, err) {
		return
	}
}
