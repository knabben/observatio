package watchers

import (
	"context"

	"github.com/gorilla/websocket"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	capv "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"

	"github.com/knabben/observatio/webserver/internal/infra/clusterapi/processor"
)

// WatchClusters send the websocket request with the serialized payload.
func WatchClusters(ctx context.Context, conn *websocket.Conn, objType string) error {
	// Create the converter function for CAPI cluster processing
	converter := func(event runtime.Object) (any, error) {
		var cluster clusterv1.Cluster
		unstructuredObj := event.(*unstructured.Unstructured)
		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(
			unstructuredObj.UnstructuredContent(), &cluster); err != nil {
			return nil, err
		}
		clusterModel, _ := processor.ProcessCluster(cluster)
		return clusterModel, nil
	}
	// Process the websocket response and send it back.
	return processWebSocket(ctx, objType, conn, converter, schema.GroupVersionResource{
		Group:    "cluster.x-k8s.io",
		Version:  "v1beta1",
		Resource: "clusters",
	})
}

// WatchVSphereClusters send the websocket request with the serialized payload.
func WatchVSphereClusters(ctx context.Context, conn *websocket.Conn, objType string) error {
	// Create the converter function for CAPV cluster processing
	converter := func(event runtime.Object) (any, error) {
		var vsphereCluster capv.VSphereCluster
		unstructuredObj := event.(*unstructured.Unstructured)
		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(
			unstructuredObj.UnstructuredContent(), &vsphereCluster); err != nil {
			return nil, err
		}
		clusterModel, _ := processor.ProcessClusterInfra(vsphereCluster)
		return clusterModel, nil
	}
	// Process the websocket response and send it back.
	return processWebSocket(ctx, objType, conn, converter, schema.GroupVersionResource{
		Group:    "infrastructure.cluster.x-k8s.io",
		Version:  "v1beta1",
		Resource: "vsphereclusters",
	})
}
