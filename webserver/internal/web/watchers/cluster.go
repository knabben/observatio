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

var (
	clusterGVR = schema.GroupVersionResource{
		Group:    "cluster.x-k8s.io",
		Version:  "v1beta1",
		Resource: "clusters",
	}

	clusterInfraGVR = schema.GroupVersionResource{
		Group:    "infrastructure.cluster.x-k8s.io",
		Version:  "v1beta1",
		Resource: "vsphereclusters",
	}
)

// WatchClusters watches Kubernetes cluster resources and streams events through a WebSocket connection.
func WatchClusters(ctx context.Context, conn *websocket.Conn, objType string) error {
	// Create the converter function for CAPI cluster processing
	converter := func(event runtime.Object) (any, error) {
		var cluster clusterv1.Cluster
		unstructuredObj := event.(*unstructured.Unstructured)
		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(
			unstructuredObj.UnstructuredContent(), &cluster); err != nil {
			return nil, err
		}
		return processor.ProcessCluster(cluster), nil
	}
	// Process the websocket response and send it back.
	return WatchResourceViaWebSocket(ctx, WebSocketWatchConfig{
		ObjectType: objType,
		Conn:       conn,
		Converter:  converter,
		GVR:        clusterGVR,
	})
}

// WatchVSphereClusters streams events of vSphereClusters to a WebSocket connection using a dynamic Kubernetes client.
func WatchVSphereClusters(ctx context.Context, conn *websocket.Conn, objType string) error {
	// Create the converter function for CAPV cluster processing
	converter := func(event runtime.Object) (any, error) {
		var vsphereCluster capv.VSphereCluster
		unstructuredObj := event.(*unstructured.Unstructured)
		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(
			unstructuredObj.UnstructuredContent(), &vsphereCluster); err != nil {
			return nil, err
		}
		return processor.ProcessClusterInfra(vsphereCluster), nil
	}
	// Process the websocket response and send it back.
	return WatchResourceViaWebSocket(ctx, WebSocketWatchConfig{
		ObjectType: objType,
		Conn:       conn,
		Converter:  converter,
		GVR:        clusterInfraGVR,
	})
}
