package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/knabben/observatio/webserver/internal/web/watchers"

	"github.com/gorilla/websocket"
)

type SubscribeRequest struct {
	Types []string `json:"types"`
}

var (
	TYPE_CLUSTER_INFRA = "cluster-infra"
	TYPE_CLUSTER       = "cluster"
)

// handleWebsocket starts the object listener based on input object request.
func handleWebsocket(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if handleError(w, http.StatusInternalServerError, err) {
		return
	}

	var subscribeRequest SubscribeRequest
	subscribeRequest, err = parseMessage(conn)
	if handleError(w, http.StatusInternalServerError, err) {
		return
	}

	for _, objType := range subscribeRequest.Types {
		switch objType {
		case TYPE_CLUSTER_INFRA:
			go func() {
				err := watchers.WatchVSphereClusters(ctx, conn, objType)
				if handleError(w, http.StatusInternalServerError, err) {
					return
				}
			}()
		case TYPE_CLUSTER:
			go func() {
				err := watchers.WatchClusters(ctx, conn, objType)
				if handleError(w, http.StatusInternalServerError, err) {
					return
				}
			}()
		}
	}
}

// parseMessage reads the first WS message
func parseMessage(conn *websocket.Conn) (subscribeRequest SubscribeRequest, err error) {
	// read the first request from the customer to start
	// the specialized watcher.
	_, msg, err := conn.ReadMessage()
	if err != nil {
		return subscribeRequest, err
	}

	// parse the type of request to the datastruct
	if err := json.Unmarshal(msg, &subscribeRequest); err != nil {
		return subscribeRequest, err
	}

	return subscribeRequest, nil
}
