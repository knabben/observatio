package system

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/knabben/observatio/webserver/internal/infra/llm"
	mcpaggregator "github.com/knabben/observatio/webserver/internal/infra/mcp"
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
	TypeClusterInfra        ObjectType = "cluster-infra"
	TypeClusterInfraDocker  ObjectType = "cluster-infra-docker"
	TypeCluster             ObjectType = "cluster"
	TypeMachine             ObjectType = "machine"
	TypeMachineInfra        ObjectType = "machine-infra"
	TypeMachineInfraDocker  ObjectType = "machine-infra-docker"
	TypeMachineDeployment   ObjectType = "machine-deployment"
	TypeChatbot             ObjectType = "chatbot"
	TypeDay2Ops             ObjectType = "day2ops"
	TypeMachineHealthCheck  ObjectType = "machinehealthcheck"
	TypeKubeadmControlPlane ObjectType = "kubeadmcontrolplane"
	TypeMachineSet          ObjectType = "machineset"
	TypeClusterClass        ObjectType = "clusterclass"
)

var (
	// watchHandlers maps object types to their respective watch functions
	watchHandlers = map[ObjectType]websocketWatcher{
		TypeClusterInfra:        watchers.WatchVSphereClusters,
		TypeClusterInfraDocker:  watchers.WatchDockerClusters,
		TypeCluster:             watchers.WatchClusters,
		TypeMachine:             watchers.WatchMachines,
		TypeMachineInfra:        watchers.WatchMachinesInfra,
		TypeMachineInfraDocker:  watchers.WatchDockerMachines,
		TypeMachineDeployment:   watchers.WatchMachineDeployments,
		TypeDay2Ops:             watchers.WatchDay2Ops,
		TypeMachineHealthCheck:  watchers.WatchMachineHealthChecks,
		TypeKubeadmControlPlane: watchers.WatchKubeadmControlPlanes,
		TypeMachineSet:          watchers.WatchMachineSets,
		TypeClusterClass:        watchers.WatchClusterClasses,
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
	defer conn.Close() // nolint

	// Past this point the connection has been hijacked for the websocket: errors (a client
	// disconnecting mid-watch, an unreadable first message) can only be logged, never turned
	// into an http.Error - that would write to an already-hijacked ResponseWriter.
	logger := log.FromContext(r.Context())

	message, err := parseMessage(conn)
	if err != nil {
		logger.Error(err, "error parsing websocket message")
		return
	}

	objType := message.Type
	watchHandler, exists := watchHandlers[ObjectType(objType)]
	if !exists {
		return
	}

	if err := watchHandler(r.Context(), conn, objType); err != nil {
		logger.Error(err, "error handling websocket watch")
	}
}

// HandleChatbot opens a connection with the client and allows chat mode. aggregator is the
// shared, process-wide tool source aggregator (specs/009-mcp-server-client-aggregator) - it is
// not rebuilt per connection.
func HandleChatbot(pool *ClientPool, aggregator *mcpaggregator.Aggregator, w http.ResponseWriter, r *http.Request) {
	var wsUpgrader = websocket.Upgrader{
		ReadBufferSize:  websocketBufferSize,
		WriteBufferSize: websocketBufferSize,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if HandleError(w, http.StatusInternalServerError, err) {
		return
	}

	// Past this point the connection has been hijacked for the websocket, so a registration
	// failure can only be logged and the connection closed - not turned into an http.Error,
	// which would write to an already-hijacked ResponseWriter.
	if err := registerClient(pool, aggregator, conn); err != nil {
		log.FromContext(r.Context()).Error(err, "error registering websocket client")
		conn.Close() // nolint
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
