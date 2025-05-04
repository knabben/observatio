package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/knabben/observatio/webserver/internal/web/watchers"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type SubscribeRequest struct {
	Types []string `json:"types"`
}

var TYPE_CLUSTER_INFRA = "cluster-infra"

// handleWebsocket starts the object listener based on input object request.
func handleWebsocket(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if handleError(w, http.StatusInternalServerError, err) {
		log.Println(err)
		return
	}

	// read the first request from the customer to start
	// the specialized watcher.
	_, msg, err := conn.ReadMessage()
	if err != nil {
		log.Println("Error reading message:", err)
		return
	}

	// parse the type of request.
	var subscribeRequest SubscribeRequest
	if err := json.Unmarshal(msg, &subscribeRequest); err != nil {
		log.Println("Error unmarshalling message:", err)
		return
	}

	for _, objType := range subscribeRequest.Types {
		switch objType {
		case TYPE_CLUSTER_INFRA:
			go func() {
				err := watchers.WatchVSphereClusters(r.Context(), conn, objType)
				if err != nil {
					log.Println(err)
					return
				}
			}()
		}
	}
}
