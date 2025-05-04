package watchers

import (
	"context"

	"github.com/gorilla/websocket"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	capv "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"

	"github.com/knabben/observatio/webserver/internal/infra/clusterapi/processor"
)

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
		cluster, _ := processor.ProcessClusterInfra(vsphereCluster)
		return cluster, nil
	}
	// Process the websocket response and send it back.
	return processWebSocket(ctx, objType, conn, converter, schema.GroupVersionResource{
		Group:    "infrastructure.cluster.x-k8s.io",
		Version:  "v1beta1",
		Resource: "vsphereclusters",
	})
}
