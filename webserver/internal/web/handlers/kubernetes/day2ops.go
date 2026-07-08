package kubernetes

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"

	"github.com/knabben/observatio/webserver/internal/infra/clusterapi"
	"github.com/knabben/observatio/webserver/internal/infra/clusterapi/day2ops"
	"github.com/knabben/observatio/webserver/internal/infra/clusterapi/fetchers"
	"github.com/knabben/observatio/webserver/internal/web/handlers/system"
)

// machineGVR identifies the Machine resource; the debugging-path detail endpoint currently only
// supports Machines, matching every Acceptance Scenario in spec 006's User Story 2 (stuck
// provisioning/bootstrap, provider-resource errors — all Machine-level cases). Cluster/
// MachineDeployment-level debugging paths are a natural follow-on, not built in this pass.
var machineGVR = schema.GroupVersionResource{Group: "cluster.x-k8s.io", Version: "v1beta1", Resource: "machines"}

// day2opsDetailResponse is the GET /api/day2ops/detail response body
// (contracts/day2ops-detail-api.md).
type day2opsDetailResponse struct {
	ObjectRef day2ops.ObjectRef `json:"objectRef"`
	Path      day2ops.DebugPath `json:"path"`
}

// HandleDay2OpsDetail returns the full, uncapped debugging path for one object, on demand — the
// scoped REST drill-in exception described in specs/006-day2-ops-dashboard/research.md R9. The
// live WS `day2ops` event already carries a capped summary of every unhealthy object's path
// (FR-004); this endpoint is only for expanding one specific object's full evidence.
func HandleDay2OpsDetail(w http.ResponseWriter, r *http.Request) {
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

	objectRef := day2ops.ObjectRef{
		Group: gvr.Group, Version: gvr.Version, Resource: gvr.Resource,
		Namespace: namespace, Name: name,
	}

	path := day2ops.DebugPath{ObjectRef: objectRef, Layers: []day2ops.DebugLayer{}}
	if gvr == machineGVR {
		if path, err = computeMachineDetailPath(ctx, dyn, objectRef, obj); system.HandleError(w, http.StatusInternalServerError, err) {
			return
		}
	}

	response := day2opsDetailResponse{ObjectRef: objectRef, Path: path}
	if err = system.WriteResponse(w, response); system.HandleError(w, http.StatusInternalServerError, err) {
		return
	}
}

// providerResourceGVRs maps a Machine's Spec.InfrastructureRef.Kind to the GVR used to fetch that
// provider-infra object generically (Constitution Principle III: no per-provider Go type needed).
var providerResourceGVRs = map[string]schema.GroupVersionResource{
	"DockerMachine":  {Group: "infrastructure.cluster.x-k8s.io", Version: "v1beta1", Resource: "dockermachines"},
	"VSphereMachine": {Group: "infrastructure.cluster.x-k8s.io", Version: "v1beta1", Resource: "vspheremachines"},
}

// computeMachineDetailPath returns the full, uncapped debugging path for one Machine, fetching its
// provider-infra object and (only when the higher layers are inconclusive) recent Events on demand.
func computeMachineDetailPath(ctx context.Context, dyn dynamic.Interface, objectRef day2ops.ObjectRef, obj *unstructured.Unstructured) (day2ops.DebugPath, error) {
	var m clusterv1.Machine
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.UnstructuredContent(), &m); err != nil {
		return day2ops.DebugPath{}, err
	}

	var provider *day2ops.ProviderResourceStatus
	if m.Spec.InfrastructureRef.Name != "" {
		if providerGVR, ok := providerResourceGVRs[m.Spec.InfrastructureRef.Kind]; ok {
			namespace := m.Spec.InfrastructureRef.Namespace
			if namespace == "" {
				namespace = m.Namespace
			}
			if providerObj, err := dyn.Resource(providerGVR).Namespace(namespace).Get(ctx, m.Spec.InfrastructureRef.Name, metav1.GetOptions{}); err == nil {
				status := day2ops.ExtractProviderResourceStatus(providerObj)
				provider = &status
			}
		}
	}

	var events []string
	if day2ops.ShouldFetchControllerActivityEvents(m) {
		events, _ = fetchers.FetchInvolvedObjectEvents(ctx, dyn, m.Namespace, m.Name, "Machine")
	}

	return day2ops.ComputeMachineDebugPath(objectRef, m, provider, events), nil
}

// parseDay2OpsDetailQuery validates and extracts the GVR/namespace/name from the request's query
// parameters, same shape as the raw-object endpoint's parseRawObjectQuery.
func parseDay2OpsDetailQuery(values url.Values) (gvr schema.GroupVersionResource, namespace, name string, err error) {
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
