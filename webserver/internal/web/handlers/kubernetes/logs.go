package kubernetes

import (
	"fmt"
	"net/http"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"

	"github.com/knabben/observatio/webserver/internal/infra/clusterapi"
	"github.com/knabben/observatio/webserver/internal/infra/clusterapi/day2ops"
	"github.com/knabben/observatio/webserver/internal/infra/clusterapi/fetchers"
	"github.com/knabben/observatio/webserver/internal/web/handlers/system"
)

// HandleControllerLogs streams a controller's Pod log output via the standard Kubernetes Pod-log
// subresource — the same data and mechanism `kubectl logs` uses (contracts/logs-api.md, FR-019,
// FR-020). Backs the Day-2 Ops dashboard's new Logs destination (User Story 5).
func HandleControllerLogs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	namespace := r.URL.Query().Get("namespace")
	deploymentName := r.URL.Query().Get("deployment")
	follow := r.URL.Query().Get("follow") == "true"

	if namespace == "" || deploymentName == "" {
		http.Error(w, "missing required query parameter(s): namespace, deployment", http.StatusBadRequest)
		return
	}

	dyn, err := clusterapi.NewDynamicClient(ctx)
	if system.HandleError(w, http.StatusInternalServerError, err) {
		return
	}

	podName, err := fetchers.FindControllerPodName(ctx, dyn, namespace, deploymentName)
	if err != nil {
		if apierrors.IsNotFound(err) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		system.HandleError(w, http.StatusInternalServerError, err)
		return
	}
	if podName == "" {
		http.Error(w, fmt.Sprintf("no Pod currently backs deployment %s/%s", namespace, deploymentName), http.StatusNotFound)
		return
	}

	clientset, err := clusterapi.NewClientset(ctx)
	if system.HandleError(w, http.StatusInternalServerError, err) {
		return
	}

	stream, err := fetchers.StreamControllerLogs(ctx, clientset, namespace, podName, follow)
	if err != nil {
		// Logs could not be retrieved (e.g. no retained log history) - the frontend's
		// FR-023 "logs unavailable" state maps to this status.
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer stream.Close() // nolint

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	flusher, canFlush := w.(http.Flusher)
	buf := make([]byte, 4096)
	for {
		n, readErr := stream.Read(buf)
		if n > 0 {
			if _, writeErr := w.Write(buf[:n]); writeErr != nil {
				return
			}
			if canFlush {
				flusher.Flush()
			}
		}
		if readErr != nil {
			return
		}
	}
}

// nodeAccessResponse is the GET /api/logs/node-access response body (contracts/logs-api.md).
type nodeAccessResponse struct {
	ObjectRef day2ops.ObjectRef `json:"objectRef"`
	Command   string            `json:"command"`
	Note      string            `json:"note"`
}

// HandleNodeAccess returns static SSH connection instructions for a Machine's node — never a live
// terminal, never credentials (FR-021, FR-022).
func HandleNodeAccess(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	gvr, namespace, name, err := parseDay2OpsDetailQuery(r.URL.Query())
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

	var m clusterv1.Machine
	if err = runtime.DefaultUnstructuredConverter.FromUnstructured(obj.UnstructuredContent(), &m); system.HandleError(w, http.StatusInternalServerError, err) {
		return
	}

	address := ""
	for _, addr := range m.Status.Addresses {
		if addr.Type == clusterv1.MachineExternalIP || addr.Type == clusterv1.MachineInternalIP {
			address = addr.Address
			break
		}
	}
	if address == "" {
		http.Error(w, fmt.Sprintf("machine %s/%s has no recorded address yet", namespace, name), http.StatusNotFound)
		return
	}

	response := nodeAccessResponse{
		ObjectRef: day2ops.ObjectRef{Group: gvr.Group, Version: gvr.Version, Resource: gvr.Resource, Namespace: namespace, Name: name},
		Command:   fmt.Sprintf("ssh capi@%s", address),
		Note:      "Observātiō does not store or manage SSH credentials. Run this command from your own machine.",
	}
	if err = system.WriteResponse(w, response); system.HandleError(w, http.StatusInternalServerError, err) {
		return
	}
}
