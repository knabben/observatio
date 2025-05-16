package kubernetes

import (
	"net/http"

	"github.com/knabben/observatio/webserver/internal/web/handlers/system"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/knabben/observatio/webserver/internal/infra/clusterapi"
	"github.com/knabben/observatio/webserver/internal/infra/clusterapi/fetchers"
	"github.com/knabben/observatio/webserver/internal/infra/models"
)

// HandleComponentsVersion returns the cluster components and its versions.
func HandleComponentsVersion(w http.ResponseWriter, r *http.Request) {
	var (
		ctx        = r.Context()
		components []clusterapi.Components
	)

	cli, err := clusterapi.NewClientWithScheme(ctx, system.Scheme)
	if system.HandleError(w, http.StatusInternalServerError, err) {
		return
	}

	components, err = clusterapi.GenerateComponentVersions(ctx, cli)
	if system.HandleError(w, http.StatusInternalServerError, err) {
		return
	}

	err = system.WriteResponse(w, components)
	if system.HandleError(w, http.StatusInternalServerError, err) {
		return
	}
}

// HandleClusterInfo returns the information about the cluster.
func HandleClusterInfo(w http.ResponseWriter, r *http.Request) {
	var (
		ctx       = r.Context()
		namespace = "kube-proxy"
		c         = ctx.Value("client").(client.Client)
	)

	services, err := clusterapi.FindServices(r.Context(), c, namespace)
	if system.HandleError(w, http.StatusInternalServerError, err) {
		return
	}

	err = system.WriteResponse(w, services)
	if system.HandleError(w, http.StatusInternalServerError, err) {
		return
	}
}

// HandleSummaryCluster returns the summary of clusters states.
func HandleSummaryCluster(w http.ResponseWriter, r *http.Request) {
	var (
		ctx     = r.Context()
		summary clusterapi.ClusterSummary
	)

	cli, err := clusterapi.NewClientWithScheme(ctx, system.Scheme)
	if system.HandleError(w, http.StatusInternalServerError, err) {
		return
	}

	summary, err = clusterapi.GenerateClusterSummary(ctx, cli)
	if system.HandleError(w, http.StatusInternalServerError, err) {
		return
	}

	err = system.WriteResponse(w, summary)
	if system.HandleError(w, http.StatusInternalServerError, err) {
		return
	}
}

// HandleClusterClasses returns the available cluster classes in the mgmt cluster.
func HandleClusterClasses(w http.ResponseWriter, r *http.Request) {
	cli, err := clusterapi.NewClientWithScheme(r.Context(), system.Scheme)
	if system.HandleError(w, http.StatusInternalServerError, err) {
		return
	}

	var clusterClasses models.ClusterClassResponse
	clusterClasses, err = fetchers.FetchClusterClass(r.Context(), cli)
	if system.HandleError(w, http.StatusInternalServerError, err) {
		return
	}

	err = system.WriteResponse(w, clusterClasses.ClusterClasses)
	if system.HandleError(w, http.StatusInternalServerError, err) {
		return
	}
}

// HandleClusterTopology returns the cluster topology by owners
func HandleClusterTopology(w http.ResponseWriter, r *http.Request) {
	var (
		ctx      = r.Context()
		topology clusterapi.ClusterTopology
	)

	cli, err := clusterapi.NewClientWithScheme(ctx, system.Scheme)
	if system.HandleError(w, http.StatusInternalServerError, err) {
		return
	}

	topology, err = clusterapi.GenerateClusterTopology(ctx, cli)
	if system.HandleError(w, http.StatusInternalServerError, err) {
		return
	}

	err = system.WriteResponse(w, topology)
	if system.HandleError(w, http.StatusInternalServerError, err) {
		return
	}
}
