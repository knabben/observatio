package fetchers

import (
	"context"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"

	"github.com/knabben/observatio/webserver/internal/infra/clusterapi/day2ops"
)

var podGVR = schema.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}

// ControllerNamespaces are the well-known namespaces CAPI core and its providers run controllers
// in, checked for crash-looping/not-ready Pods (FR-014, research.md R7).
var ControllerNamespaces = []string{"capi-system", "capd-system", "capv-system"}

// FetchControllerPodStatuses lists Pods in a controller namespace and reports any that aren't
// ready, including a CrashLoopBackOff-style waiting reason when present. A namespace that doesn't
// exist (that provider isn't installed) is skipped, not an error.
func FetchControllerPodStatuses(ctx context.Context, dyn dynamic.Interface, namespace string) ([]day2ops.ControllerPodStatus, error) {
	list, err := dyn.Resource(podGVR).Namespace(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	var statuses []day2ops.ControllerPodStatus
	for _, item := range list.Items {
		if isPodReady(item) {
			continue
		}
		statuses = append(statuses, day2ops.ControllerPodStatus{
			Namespace:     namespace,
			PodName:       item.GetName(),
			Ready:         false,
			WaitingReason: podWaitingReason(item),
		})
	}
	return statuses, nil
}

func isPodReady(pod unstructured.Unstructured) bool {
	conditions, _, _ := unstructured.NestedSlice(pod.Object, "status", "conditions")
	for _, raw := range conditions {
		cond, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		condType, _, _ := unstructured.NestedString(cond, "type")
		if condType != "Ready" {
			continue
		}
		condStatus, _, _ := unstructured.NestedString(cond, "status")
		return condStatus == "True"
	}
	return false
}

func podWaitingReason(pod unstructured.Unstructured) string {
	statuses, _, _ := unstructured.NestedSlice(pod.Object, "status", "containerStatuses")
	for _, raw := range statuses {
		cs, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		reason, _, _ := unstructured.NestedString(cs, "state", "waiting", "reason")
		if reason != "" {
			return reason
		}
	}
	return ""
}
