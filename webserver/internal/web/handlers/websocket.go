package handlers

import (
	"context"
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
const (
	websocketBufferSize = 1024
)

// ObjectType represents supported websocket subscription types
type ObjectType string

const (
	TypeClusterInfra      ObjectType = "cluster-infra"
	TypeCluster           ObjectType = "cluster"
	TypeMachine           ObjectType = "machine"
	TypeMachineInfra      ObjectType = "machine-infra"
	TypeMachineDeployment ObjectType = "machine-deployment"
)

// websocketWatcher represents a function that watches specific resource types
type websocketWatcher func(context.Context, *websocket.Conn, string) error

var (
	// watchHandlers maps object types to their respective watch functions
	watchHandlers = map[ObjectType]websocketWatcher{
		TypeClusterInfra:      watchers.WatchVSphereClusters,
		TypeCluster:           watchers.WatchClusters,
		TypeMachine:           watchers.WatchMachines,
		TypeMachineInfra:      watchers.WatchMachinesInfra,
		TypeMachineDeployment: watchers.WatchMachineDeployments,
	}
)

// handleWebsocket starts the object listener based on the input object request.
func handleWebsocket(w http.ResponseWriter, r *http.Request) {
	var wsUpgrader = websocket.Upgrader{
		ReadBufferSize:  websocketBufferSize,
		WriteBufferSize: websocketBufferSize,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if handleError(w, http.StatusInternalServerError, err) {
		return
	}

	subscribeRequest, err := parseMessage(conn)
	if handleError(w, http.StatusInternalServerError, err) {
		return
	}

	for _, objType := range subscribeRequest.Types {
		watchHandler, exists := watchHandlers[ObjectType(objType)]
		if !exists {
			continue
		}
		err := watchHandler(r.Context(), conn, objType)
		if handleError(w, http.StatusInternalServerError, err) {
			return
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
