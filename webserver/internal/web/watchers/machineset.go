package watchers

import (
	"context"

	"github.com/gorilla/websocket"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"

	"github.com/knabben/observatio/webserver/internal/infra/clusterapi/processor"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

// WatchMachineSets watches Kubernetes MachineSet resources and streams events through a
// WebSocket connection — the first-class list page for rollout-state detail behind the Day-2 Ops
// dashboard's stalled-rollout warning (006). Reuses machineSetGVR, already declared in day2ops.go.
func WatchMachineSets(ctx context.Context, conn *websocket.Conn, objType string) error {
	converter := func(event runtime.Object) (any, error) {
		var ms clusterv1.MachineSet
		unstructuredObj := event.(*unstructured.Unstructured)
		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(
			unstructuredObj.UnstructuredContent(), &ms); err != nil {
			return nil, err
		}
		return processor.ProcessMachineSet(ms), nil
	}
	return WatchResourceViaWebSocket(ctx, WebSocketWatchConfig{
		ObjectType: objType,
		Conn:       conn,
		Converter:  converter,
		GVR:        machineSetGVR,
	})
}
