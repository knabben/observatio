package watchers

import (
	"context"

	"github.com/gorilla/websocket"
	controlplanev1 "sigs.k8s.io/cluster-api/controlplane/kubeadm/api/v1beta1"

	"github.com/knabben/observatio/webserver/internal/infra/clusterapi/processor"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var kubeadmControlPlaneGVR = schema.GroupVersionResource{
	Group:    "controlplane.cluster.x-k8s.io",
	Version:  "v1beta1",
	Resource: "kubeadmcontrolplanes",
}

// WatchKubeadmControlPlanes watches Kubernetes KubeadmControlPlane resources and streams events
// through a WebSocket connection — the first-class list page for control-plane replica health and
// etcd conditions.
func WatchKubeadmControlPlanes(ctx context.Context, conn *websocket.Conn, objType string) error {
	converter := func(event runtime.Object) (any, error) {
		var kcp controlplanev1.KubeadmControlPlane
		unstructuredObj := event.(*unstructured.Unstructured)
		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(
			unstructuredObj.UnstructuredContent(), &kcp); err != nil {
			return nil, err
		}
		return processor.ProcessKubeadmControlPlane(kcp), nil
	}
	return WatchResourceViaWebSocket(ctx, WebSocketWatchConfig{
		ObjectType: objType,
		Conn:       conn,
		Converter:  converter,
		GVR:        kubeadmControlPlaneGVR,
	})
}
