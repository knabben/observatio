package kubernetes

import (
	"context"
	"fmt"
	"net/http"

	"github.com/knabben/observatio/webserver/internal/web/handlers/system"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/knabben/observatio/webserver/internal/infra/clusterapi"
	"github.com/knabben/observatio/webserver/internal/infra/clusterapi/fetchers"
	"github.com/knabben/observatio/webserver/internal/infra/models"
	"github.com/knabben/observatio/webserver/internal/infra/providerkind"
)

// HandleClusterList returns the information about the cluster.
func HandleClusterList(w http.ResponseWriter, r *http.Request) {
	fetchClusterData[models.ClusterResponse](r.Context(), w, fetchers.FetchClusters)
}

// HandleClusterInfraList returns the infrastructure-specific cluster list for the requested
// (`?provider=docker|vsphere`) or auto-detected installed provider.
func HandleClusterInfraList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cli, err := clusterapi.NewClientWithScheme(ctx, system.Scheme)
	if system.HandleError(w, http.StatusInternalServerError, err) {
		return
	}

	capability, err := clusterapi.GenerateInfrastructureCapability(ctx, cli)
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
		response, err := fetchers.FetchClustersInfra(ctx, cli)
		if system.HandleError(w, http.StatusInternalServerError, err) {
			return
		}
		if err = system.WriteResponse(w, response); system.HandleError(w, http.StatusInternalServerError, err) {
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
		response, err := fetchers.FetchClusterInfraDocker(ctx, dyn)
		if system.HandleError(w, http.StatusInternalServerError, err) {
			return
		}
		if err = system.WriteResponse(w, response); system.HandleError(w, http.StatusInternalServerError, err) {
			return
		}
	}
}

// resolveInfraProvider validates an explicit `?provider=` value, or auto-selects the first
// installed provider (docker, then vsphere) when none was requested. Returns ok=false only
// for an unrecognized explicit value; an empty provider with ok=true means none is installed.
func resolveInfraProvider(requested string, capability models.InfrastructureCapability) (provider string, ok bool) {
	switch requested {
	case "":
		if capability.Docker.Installed {
			return providerkind.Docker, true
		}
		if capability.VSphere.Installed {
			return providerkind.VSphere, true
		}
		return "", true
	case providerkind.Docker, providerkind.VSphere:
		return requested, true
	default:
		return "", false
	}
}

// fetchClusterData write the return of the response based on a cluster type.
func fetchClusterData[T any](ctx context.Context, w http.ResponseWriter, fetchFunc func(context.Context, client.Client) (T, error)) {
	var response T
	cli, err := clusterapi.NewClientWithScheme(ctx, system.Scheme)
	if system.HandleError(w, http.StatusInternalServerError, err) {
		return
	}
	if response, err = fetchFunc(ctx, cli); system.HandleError(w, http.StatusInternalServerError, err) {
		return
	}
	if err = system.WriteResponse(w, response); system.HandleError(w, http.StatusInternalServerError, err) {
		return
	}
}
