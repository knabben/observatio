package handlers

import (
	"net/http"

	"github.com/knabben/observatio/webserver/internal/infra/clusterapi"
	"github.com/knabben/observatio/webserver/internal/infra/clusterapi/fetchers"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// handleComponentsVersion returns the cluster components and its versions.
func handleComponentsVersion(w http.ResponseWriter, r *http.Request) {
	var (
		ctx        = r.Context()
		components []clusterapi.Components
	)

	cli, err := clusterapi.NewClientWithScheme(ctx, scheme)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	if components, err = clusterapi.GenerateComponentVersions(ctx, cli); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	if err = writeResponse(w, components); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
}

// handleClusterInfo returns the information about the cluster.
func handleClusterInfo(w http.ResponseWriter, r *http.Request) {
	var (
		ctx       = r.Context()
		namespace = "kube-proxy"
		c         = ctx.Value("client").(client.Client)
	)

	services, err := clusterapi.FindServices(r.Context(), c, namespace)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	if err = writeResponse(w, services); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
}

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
	var clusterClasses []fetchers.ClusterClass

	cli, err := clusterapi.NewClientWithScheme(r.Context(), scheme)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	if clusterClasses, err = fetchers.FetchClusterClass(r.Context(), cli); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	if err = writeResponse(w, clusterClasses); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
}
