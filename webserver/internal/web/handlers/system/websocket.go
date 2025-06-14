package system

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/knabben/observatio/webserver/internal/infra/llm"
	"github.com/knabben/observatio/webserver/internal/web/watchers"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/gorilla/websocket"
)

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
	TypeChatbot           ObjectType = "chatbot"
)

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

// websocketWatcher represents a function that watches specific resource types
type websocketWatcher func(context.Context, *websocket.Conn, string) error

// HandleWatcher starts the object listener based on the input object request.
func HandleWatcher(w http.ResponseWriter, r *http.Request) {
	var wsUpgrader = websocket.Upgrader{
		ReadBufferSize:  websocketBufferSize,
		WriteBufferSize: websocketBufferSize,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if handleError(w, http.StatusInternalServerError, err) {
		return
	}

	message, err := parseMessage(conn)
	if handleError(w, http.StatusInternalServerError, err) {
		return
	}

	objType := message.Type
	watchHandler, exists := watchHandlers[ObjectType(objType)]
	if !exists {
		return
	}

	err = watchHandler(r.Context(), conn, objType)
	if handleError(w, http.StatusInternalServerError, err) {
		return
	}
}

// HandleChatbot opens a connection with the client and allows chat mode.
func HandleChatbot(pool *ClientPool, w http.ResponseWriter, r *http.Request) {
	var wsUpgrader = websocket.Upgrader{
		ReadBufferSize:  websocketBufferSize,
		WriteBufferSize: websocketBufferSize,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if HandleError(w, http.StatusInternalServerError, err) {
		return
	}

	err = registerClient(pool, conn)
	if HandleError(w, http.StatusInternalServerError, err) {
		return
	}
}

// parseMessage reads the first WS message
func parseMessage(conn *websocket.Conn) (message *llm.ChatMessage, err error) {
	_, msg, err := conn.ReadMessage()
	if err != nil {
		return message, err
	}
	if err := json.Unmarshal(msg, &message); err != nil {
		return message, err
	}
	return message, nil
}

// handleError write down an error with code to the writer response.
func handleError(w http.ResponseWriter, code int, err error) (hasError bool) {
	var logger = log.FromContext(context.Background())
	hasError = err != nil
	if hasError {
		logger.Error(err, "error handling websocket request")
		http.Error(w, err.Error(), code)
	}
	return hasError
}
