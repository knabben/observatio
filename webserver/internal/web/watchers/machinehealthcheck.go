package watchers

import (
	"context"

	"github.com/gorilla/websocket"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"

	"github.com/knabben/observatio/webserver/internal/infra/clusterapi/processor"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

// WatchMachineHealthChecks watches Kubernetes MachineHealthCheck resources and streams events
// through a WebSocket connection — the first-class list page for the remediation policy behind
// the Day-2 Ops dashboard's self-healing/needs-investigation severity classification (006/US4).
func WatchMachineHealthChecks(ctx context.Context, conn *websocket.Conn, objType string) error {
	converter := func(event runtime.Object) (any, error) {
		var mhc clusterv1.MachineHealthCheck
		unstructuredObj := event.(*unstructured.Unstructured)
		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(
			unstructuredObj.UnstructuredContent(), &mhc); err != nil {
			return nil, err
		}
		return processor.ProcessMachineHealthCheck(mhc), nil
	}
	return WatchResourceViaWebSocket(ctx, WebSocketWatchConfig{
		ObjectType: objType,
		Conn:       conn,
		Converter:  converter,
		GVR:        machineHealthCheckGVR,
	})
}
