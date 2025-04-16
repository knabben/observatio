package handlers

import (
	"github.com/knabben/observatio/webserver/internal/infra/clusterapi"
	"net/http"
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
