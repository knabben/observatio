package watchers

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"

	"github.com/gorilla/websocket"
	"github.com/knabben/observatio/webserver/internal/infra/clusterapi"
)

type EventResponse struct {
	Type  string      `json:"type"`
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

// WebSocketWatchConfig holds configuration for watching Kubernetes resources via WebSocket
type WebSocketWatchConfig struct {
	ObjectType string
	Conn       *websocket.Conn
	Converter  func(runtime.Object) (any, error)
	GVR        schema.GroupVersionResource
}

// WatchResourceViaWebSocket opens a WebSocket connection and streams Kubernetes resource events using a dynamic client.
// ctx is the context for controlling the lifetime of the function.
// objType represents the type of the resource being processed.
// conn is the WebSocket connection through which events are sent.
// converter transforms runtime.Object into a serializable format for WebSocket communication.
// gvr specifies the GroupVersionResource for the Kubernetes resource to watch.
// Returns an error if the WebSocket connection fails, the watch cannot be established, or event conversion fails.
func WatchResourceViaWebSocket(ctx context.Context, config WebSocketWatchConfig) error {
	dynamicClient, err := clusterapi.NewDynamicClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create dynamic client: %w", err)
	}

	watcher, err := dynamicClient.Resource(config.GVR).Namespace("").Watch(ctx, metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to create watcher for %v: %w", config.GVR, err)
	}
	defer watcher.Stop()

	return streamEvents(config, watcher)
}

// streamEvents handles the streaming of events to the WebSocket connection
func streamEvents(config WebSocketWatchConfig, watcher watch.Interface) error {
	for event := range watcher.ResultChan() {
		data, err := config.Converter(event.Object)
		if err != nil {
			return fmt.Errorf("failed to convert event data: %w", err)
		}
		response := EventResponse{
			Type:  string(event.Type),
			Event: config.ObjectType,
			Data:  data,
		}
		if err = config.Conn.WriteJSON(response); err != nil {
			return fmt.Errorf("failed to write to websocket: %w", err)
		}
	}
	return nil
}
