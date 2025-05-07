package watchers

import (
	"context"

	"github.com/gorilla/websocket"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"

	"github.com/knabben/observatio/webserver/internal/infra/clusterapi/processor"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var machineDeploymentGVR = schema.GroupVersionResource{
	Group:    "cluster.x-k8s.io",
	Version:  "v1beta1",
	Resource: "machinedeployments",
}

// WatchMachineDeployments watches Kubernetes cluster resources and streams events through a WebSocket connection.
func WatchMachineDeployments(ctx context.Context, conn *websocket.Conn, objType string) error {
	// Create the converter function for CAPI machine deployment processing
	converter := func(event runtime.Object) (any, error) {
		var machineDeployment clusterv1.MachineDeployment
		unstructuredObj := event.(*unstructured.Unstructured)
		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(
			unstructuredObj.UnstructuredContent(), &machineDeployment); err != nil {
			return nil, err
		}
		return processor.ProcessMachineDeployment(machineDeployment), nil
	}
	// Process the websocket response and send it back.
	return WatchResourceViaWebSocket(ctx, WebSocketWatchConfig{
		ObjectType: objType,
		Conn:       conn,
		Converter:  converter,
		GVR:        machineDeploymentGVR,
	})
}
