package fetchers

import (
	"context"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"

	"github.com/knabben/observatio/webserver/internal/infra/models"
)

// dockerClusterGVR identifies the DockerCluster resource, read via the dynamic client since
// no typed Go package for it is imported (see specs/004-detect-infra-adapt-ui/research.md R3).
var dockerClusterGVR = schema.GroupVersionResource{
	Group:    "infrastructure.cluster.x-k8s.io",
	Version:  "v1beta1",
	Resource: "dockerclusters",
}

// FetchClusterInfraDocker retrieves and processes all DockerCluster resources into a
// ClusterInfraDockerResponse, mirroring FetchClustersInfra for vSphere.
func FetchClusterInfraDocker(ctx context.Context, dyn dynamic.Interface) (models.ClusterInfraDockerResponse, error) {
	clusters, err := ListClusterInfraDocker(ctx, dyn)
	if err != nil {
		return models.ClusterInfraDockerResponse{}, err
	}

	failing := 0
	for _, c := range clusters {
		if !c.Ready {
			failing++
		}
	}
	return models.ClusterInfraDockerResponse{
		Total:    len(clusters),
		Failing:  failing,
		Clusters: clusters,
	}, nil
}

// ListClusterInfraDocker lists all DockerCluster resources across namespaces via the dynamic
// client and decodes only the fields the Docker infra view needs.
func ListClusterInfraDocker(ctx context.Context, dyn dynamic.Interface) ([]models.ClusterInfraDocker, error) {
	list, err := dyn.Resource(dockerClusterGVR).Namespace("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	clusters := make([]models.ClusterInfraDocker, 0, len(list.Items))
	for _, item := range list.Items {
		clusters = append(clusters, ProcessDockerCluster(item))
	}
	return clusters, nil
}

// ProcessDockerCluster decodes a single unstructured DockerCluster object into a
// models.ClusterInfraDocker.
func ProcessDockerCluster(obj unstructured.Unstructured) models.ClusterInfraDocker {
	var clusterOwner string
	for _, owner := range obj.GetOwnerReferences() {
		clusterOwner = owner.Name
	}

	loadBalancerIP, _, _ := unstructured.NestedString(obj.Object, "spec", "loadBalancerIP")
	ready, _, _ := unstructured.NestedBool(obj.Object, "status", "ready")

	return models.ClusterInfraDocker{
		ObjectMeta: metav1.ObjectMeta{
			Name:              obj.GetName(),
			Namespace:         obj.GetNamespace(),
			CreationTimestamp: obj.GetCreationTimestamp(),
		},
		Cluster:        clusterOwner,
		Age:            formatDuration(time.Since(obj.GetCreationTimestamp().Time)),
		LoadBalancerIP: loadBalancerIP,
		Ready:          ready,
	}
}

// formatDuration mirrors processor.formatDuration so Age renders identically across
// the vSphere and Docker infra views (that helper is unexported in another package).
func formatDuration(d time.Duration) string {
	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh", days, hours)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}
