package system

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/knabben/observatio/webserver/internal/infra/llm"
	"github.com/knabben/observatio/webserver/internal/web/watchers"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/gorilla/websocket"
)

type WSMessage struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
}

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

// HandleChatBot opens a connection with the client and allows chat mode.
func HandleChatBot(pool *ClientPool, w http.ResponseWriter, r *http.Request) {
	var wsUpgrader = websocket.Upgrader{
		ReadBufferSize:  websocketBufferSize,
		WriteBufferSize: websocketBufferSize,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if HandleError(w, http.StatusInternalServerError, err) {
		return
	}

	registerClient(pool, conn)
}

func registerClient(pool *ClientPool, conn *websocket.Conn) {
	llmClient, err := llm.NewClient()
	if err != nil {
		log.FromContext(context.Background()).Error(err, "error creating llm client")
	}

	client := &WSClient{
		ID:        uuid.New().String(),
		pool:      pool,
		conn:      conn,
		Send:      make(chan []byte, 256),
		LLMClient: &llmClient,
	}

	client.pool.Register <- client

	go client.reader()
	go client.writer()
}

// parseMessage reads the first WS message
func parseMessage(conn *websocket.Conn) (message WSMessage, err error) {
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
