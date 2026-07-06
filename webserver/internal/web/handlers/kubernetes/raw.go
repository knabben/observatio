package kubernetes

import (
	"fmt"
	"net/http"
	"net/url"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/knabben/observatio/webserver/internal/infra/clusterapi"
	"github.com/knabben/observatio/webserver/internal/web/handlers/system"
)

// HandleRawObject returns the complete, unmodified Kubernetes object for a given
// group/version/resource/namespace/name, bypassing every curated DTO. Backs the object detail
// screens' "YAML" tree tab (spec 005 FR-011/FR-012).
func HandleRawObject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	gvr, namespace, name, err := parseRawObjectQuery(r.URL.Query())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dyn, err := clusterapi.NewDynamicClient(ctx)
	if system.HandleError(w, http.StatusInternalServerError, err) {
		return
	}

	obj, err := dyn.Resource(gvr).Namespace(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		system.HandleError(w, http.StatusInternalServerError, err)
		return
	}

	if err = system.WriteResponse(w, obj.Object); system.HandleError(w, http.StatusInternalServerError, err) {
		return
	}
}

// parseRawObjectQuery validates and extracts the GVR/namespace/name from the request's query
// parameters (contracts/raw-object-api.md). Extracted from the handler so this validation logic
// is unit-testable without a live Kubernetes client.
func parseRawObjectQuery(values url.Values) (gvr schema.GroupVersionResource, namespace, name string, err error) {
	group := values.Get("group")
	version := values.Get("version")
	resource := values.Get("resource")
	namespace = values.Get("namespace")
	name = values.Get("name")

	var missing []string
	if version == "" {
		missing = append(missing, "version")
	}
	if resource == "" {
		missing = append(missing, "resource")
	}
	if namespace == "" {
		missing = append(missing, "namespace")
	}
	if name == "" {
		missing = append(missing, "name")
	}
	if len(missing) > 0 {
		return gvr, "", "", fmt.Errorf("missing required query parameter(s): %v", missing)
	}

	return schema.GroupVersionResource{Group: group, Version: version, Resource: resource}, namespace, name, nil
}
