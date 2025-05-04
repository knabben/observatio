package watchers

import (
	"context"

	"github.com/gorilla/websocket"
	"github.com/knabben/observatio/webserver/internal/infra/clusterapi/processor"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	capv "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
)

func WatchVSphereClusters(ctx context.Context, conn *websocket.Conn, objType string) error {
	gvr := schema.GroupVersionResource{
		Group:    "infrastructure.cluster.x-k8s.io",
		Version:  "v1beta1",
		Resource: "vsphereclusters",
	}
	converter := func(event runtime.Object) (any, error) {
		var vsphereCluster capv.VSphereCluster
		unstructuredObj := event.(*unstructured.Unstructured)
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(
			unstructuredObj.UnstructuredContent(), &vsphereCluster)
		if err != nil {
			return nil, err
		}
		cluster, _ := processor.ProcessClusterInfra(vsphereCluster)
		return cluster, nil
	}
	return processWebSocket(ctx, conn, gvr, converter, objType)
}
