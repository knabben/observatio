package handlers

import (
	"net/http"

	"github.com/knabben/observatio/webserver/internal/infra/clusterapi"
	"github.com/knabben/observatio/webserver/internal/infra/clusterapi/fetchers"
	"github.com/knabben/observatio/webserver/internal/infra/models"
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
	cli, err := clusterapi.NewClientWithScheme(r.Context(), scheme)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	var clusterClasses models.ClusterClassResponse
	if clusterClasses, err = fetchers.FetchClusterClass(r.Context(), cli); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	if err = writeResponse(w, clusterClasses.ClusterClasses); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
}

// handleClusterTopology returns the cluster topology by owners
func handleClusterTopology(w http.ResponseWriter, r *http.Request) {
	var (
		ctx      = r.Context()
		topology clusterapi.ClusterTopology
	)

	cli, err := clusterapi.NewClientWithScheme(ctx, scheme)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	if topology, err = clusterapi.GenerateClusterTopology(ctx, cli); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	if err = writeResponse(w, topology); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
}
