package kubernetes

import (
	"net/http"

	"github.com/knabben/observatio/webserver/internal/infra/clusterapi"
	"github.com/knabben/observatio/webserver/internal/infra/clusterapi/fetchers"
	"github.com/knabben/observatio/webserver/internal/web/handlers/system"
)

// HandleMachineDeployments returns the information about the machine deployments
func HandleMachineDeployments(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()
	c, err := clusterapi.NewClientWithScheme(ctx, system.Scheme)
	if system.HandleError(w, http.StatusInternalServerError, err) {
		return
	}

	mds, err := fetchers.FetchMachineDeployment(r.Context(), c)
	if system.HandleError(w, http.StatusInternalServerError, err) {
		return
	}

	if err = system.WriteResponse(w, mds); system.HandleError(w, http.StatusInternalServerError, err) {
		return
	}
}
