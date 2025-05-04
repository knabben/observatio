package watchers

import (
	"context"
	"log"

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

func processWebSocket(ctx context.Context, conn *websocket.Conn, gvr schema.GroupVersionResource, converter func(event runtime.Object) (any, error), objType string) error {
	cs, err := clusterapi.NewDynamicClient(ctx)
	if err != nil {
		return err
	}

	watcher, err := cs.Resource(gvr).
		Namespace("").
		Watch(context.TODO(),
			metav1.ListOptions{},
		)
	if err != nil {
		return err
	}
	defer watcher.Stop()

	for event := range watcher.ResultChan() {
		data, _ := converter(event.Object)
		response := EventResponse{
			Type:  string(event.Type),
			Event: objType,
			Data:  data,
		}
		if err = conn.WriteJSON(response); err != nil {
			log.Println(err)
			break
		}
	}
	return nil
}
