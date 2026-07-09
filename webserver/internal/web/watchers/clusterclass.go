package watchers

import (
	"context"

	"github.com/gorilla/websocket"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"

	"github.com/knabben/observatio/webserver/internal/infra/clusterapi/processor"
	"github.com/knabben/observatio/webserver/internal/infra/models"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var clusterClassGVR = schema.GroupVersionResource{
	Group:    "cluster.x-k8s.io",
	Version:  "v1beta1",
	Resource: "clusterclasses",
}

// clusterClassWithMeta adapts models.ClusterClass's flat Name/Namespace fields (research.md R5 —
// no new model/processor needed, reusing the existing main-dashboard widget's model unchanged)
// into the `metadata.name` shape every other first-class object page's BaseLister/ObjectTable
// relies on for row keys and search filtering.
type clusterClassWithMeta struct {
	models.ClusterClass
	Metadata struct {
		Name      string `json:"name"`
		Namespace string `json:"namespace"`
	} `json:"metadata"`
}

// WatchClusterClasses watches Kubernetes ClusterClass resources and streams events through a
// WebSocket connection — the first-class list page for ClusterClass, alongside (not replacing) the
// existing main-dashboard widget.
func WatchClusterClasses(ctx context.Context, conn *websocket.Conn, objType string) error {
	converter := func(event runtime.Object) (any, error) {
		var cc clusterv1.ClusterClass
		unstructuredObj := event.(*unstructured.Unstructured)
		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(
			unstructuredObj.UnstructuredContent(), &cc); err != nil {
			return nil, err
		}
		processed := processor.ProcessClusterClass(cc)
		wrapped := clusterClassWithMeta{ClusterClass: processed}
		wrapped.Metadata.Name = processed.Name
		wrapped.Metadata.Namespace = processed.Namespace
		return wrapped, nil
	}
	return WatchResourceViaWebSocket(ctx, WebSocketWatchConfig{
		ObjectType: objType,
		Conn:       conn,
		Converter:  converter,
		GVR:        clusterClassGVR,
	})
}
