package watchers

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/gorilla/websocket"
	"github.com/knabben/observatio/webserver/internal/infra/clusterapi"
)

type EventResponse struct {
	Type  string      `json:"type"`
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

// processWebSocket send the generic payload back to the websocket client.
func processWebSocket(
	ctx context.Context,
	objType string,
	conn *websocket.Conn,
	converter func(event runtime.Object) (any, error),
	gvr schema.GroupVersionResource,
) error {
	dynamicClient, err := clusterapi.NewDynamicClient(ctx)
	if err != nil {
		return err
	}

	// Start a new dynamic watch to listen for generic objects.
	watcher, err := dynamicClient.Resource(gvr).Namespace("").Watch(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	defer watcher.Stop()

	// Iterate through the result and write back the response
	// with formatted event.
	for event := range watcher.ResultChan() {
		data, err := converter(event.Object)
		if err != nil {
			return err
		}
		if err = conn.WriteJSON(EventResponse{
			Type:  string(event.Type),
			Event: objType,
			Data:  data,
		}); err != nil {
			return err
		}
	}
	return nil
}
