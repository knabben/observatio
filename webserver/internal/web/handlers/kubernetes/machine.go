package kubernetes

import (
	"fmt"
	"net/http"

	"github.com/knabben/observatio/webserver/internal/infra/clusterapi"
	"github.com/knabben/observatio/webserver/internal/infra/clusterapi/fetchers"
	"github.com/knabben/observatio/webserver/internal/infra/providerkind"
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

// HandleMachineInfra returns the infrastructure-specific machine list for the requested
// (`?provider=docker|vsphere`) or auto-detected installed provider.
func HandleMachineInfra(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	c, err := clusterapi.NewClientWithScheme(ctx, system.Scheme)
	if system.HandleError(w, http.StatusInternalServerError, err) {
		return
	}

	capability, err := clusterapi.GenerateInfrastructureCapability(ctx, c)
	if system.HandleError(w, http.StatusInternalServerError, err) {
		return
	}

	requested := r.URL.Query().Get("provider")
	provider, ok := resolveInfraProvider(requested, capability)
	if !ok {
		http.Error(w, fmt.Sprintf("unsupported provider %q: must be %q or %q", requested, providerkind.Docker, providerkind.VSphere), http.StatusBadRequest)
		return
	}
	if provider == "" {
		http.Error(w, "no supported infrastructure provider is installed in the connected environment", http.StatusNotFound)
		return
	}

	switch provider {
	case providerkind.VSphere:
		if !capability.VSphere.Installed {
			http.Error(w, "vsphere infrastructure provider is not installed", http.StatusNotFound)
			return
		}
		machinesInfra, err := fetchers.FetchMachineInfra(ctx, c)
		if system.HandleError(w, http.StatusInternalServerError, err) {
			return
		}
		if err = system.WriteResponse(w, machinesInfra); system.HandleError(w, http.StatusInternalServerError, err) {
			return
		}
	case providerkind.Docker:
		if !capability.Docker.Installed {
			http.Error(w, "docker infrastructure provider is not installed", http.StatusNotFound)
			return
		}
		dyn, err := clusterapi.NewDynamicClient(ctx)
		if system.HandleError(w, http.StatusInternalServerError, err) {
			return
		}
		machinesInfra, err := fetchers.FetchMachineInfraDocker(ctx, dyn)
		if system.HandleError(w, http.StatusInternalServerError, err) {
			return
		}
		if err = system.WriteResponse(w, machinesInfra); system.HandleError(w, http.StatusInternalServerError, err) {
			return
		}
	}
}
