package watchers

import (
	"context"

	"github.com/gorilla/websocket"
	capv "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"

	"github.com/knabben/observatio/webserver/internal/infra/clusterapi/processor"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	machineGVR = schema.GroupVersionResource{
		Group:    "cluster.x-k8s.io",
		Version:  "v1beta1",
		Resource: "machines",
	}
	machineInfraGVR = schema.GroupVersionResource{
		Group:    "infrastructure.cluster.x-k8s.io",
		Version:  "v1beta1",
		Resource: "vspheremachines",
	}
)

// WatchMachines watches Kubernetes cluster resources and streams events through a WebSocket connection.
func WatchMachines(ctx context.Context, conn *websocket.Conn, objType string) error {
	// Create the converter function for CAPI machine processing
	converter := func(event runtime.Object) (any, error) {
		var machine clusterv1.Machine
		unstructuredObj := event.(*unstructured.Unstructured)
		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(
			unstructuredObj.UnstructuredContent(), &machine); err != nil {
			return nil, err
		}
		return processor.ProcessMachine(machine), nil
	}
	// Process the websocket response and send it back.
	return WatchResourceViaWebSocket(ctx, WebSocketWatchConfig{
		ObjectType: objType,
		Conn:       conn,
		Converter:  converter,
		GVR:        machineGVR,
	})
}

// WatchMachinesInfra watches Kubernetes cluster resources and streams events through a WebSocket connection.
func WatchMachinesInfra(ctx context.Context, conn *websocket.Conn, objType string) error {
	// Create the converter function for CAPI machine processing
	converter := func(event runtime.Object) (any, error) {
		var machine capv.VSphereMachine
		unstructuredObj := event.(*unstructured.Unstructured)
		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(
			unstructuredObj.UnstructuredContent(), &machine); err != nil {
			return nil, err
		}
		return processor.ProcessMachineInfra(machine), nil
	}
	// Process the websocket response and send it back.
	return WatchResourceViaWebSocket(ctx, WebSocketWatchConfig{
		ObjectType: objType,
		Conn:       conn,
		Converter:  converter,
		GVR:        machineInfraGVR,
	})
}
