package handlers

import (
	"github.com/knabben/observatio/webserver/internal/infra/clusterapi"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

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

// handleClusterList returns the information about the cluster.
func handleClusterList(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()
	c, err := clusterapi.NewClientWithScheme(ctx, scheme)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	clusters, err := clusterapi.GenerateClusterList(r.Context(), c)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	if err = writeResponse(w, clusters); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
}
