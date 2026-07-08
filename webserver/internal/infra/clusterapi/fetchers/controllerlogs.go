package fetchers

import (
	"context"
	"fmt"
	"io"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

var deploymentGVR = schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}

// FindControllerPodName resolves a controller Deployment's current Pod name, by reading the
// Deployment's selector and listing Pods matching it. Returns an empty string (no error) when the
// Deployment has no ready Pod yet.
func FindControllerPodName(ctx context.Context, dyn dynamic.Interface, namespace, deploymentName string) (string, error) {
	deployment, err := dyn.Resource(deploymentGVR).Namespace(namespace).Get(ctx, deploymentName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	matchLabels, _, err := unstructured.NestedStringMap(deployment.Object, "spec", "selector", "matchLabels")
	if err != nil {
		return "", err
	}

	selector := make([]string, 0, len(matchLabels))
	for k, v := range matchLabels {
		selector = append(selector, fmt.Sprintf("%s=%s", k, v))
	}

	pods, err := dyn.Resource(podGVR).Namespace(namespace).List(ctx, metav1.ListOptions{LabelSelector: strings.Join(selector, ",")})
	if err != nil {
		return "", err
	}
	if len(pods.Items) == 0 {
		return "", nil
	}
	return pods.Items[0].GetName(), nil
}

// StreamControllerLogs opens the standard Kubernetes Pod-log subresource for a Pod — the same
// mechanism and data `kubectl logs` uses, not a reformatted/filtered view (research.md R10).
// The caller is responsible for closing the returned stream.
func StreamControllerLogs(ctx context.Context, clientset kubernetes.Interface, namespace, podName string, follow bool) (io.ReadCloser, error) {
	return clientset.CoreV1().Pods(namespace).GetLogs(podName, &corev1.PodLogOptions{Follow: follow}).Stream(ctx)
}
