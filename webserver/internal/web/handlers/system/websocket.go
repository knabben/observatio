package system

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/knabben/observatio/webserver/internal/web/watchers"

	"github.com/gorilla/websocket"
)

type SubscribeRequest struct {
	Type string `json:"types"`
	Data string `json:"data"`
}

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

	TypeChatbot ObjectType = "chatbot"
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

// HandleWebsocket starts the object listener based on the input object request.
func HandleWebsocket(w http.ResponseWriter, r *http.Request) {
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

	objType := subscribeRequest.Type
	watchHandler, exists := watchHandlers[ObjectType(objType)]
	if !exists {
		return
	}
	err = watchHandler(r.Context(), conn, objType)
	if handleError(w, http.StatusInternalServerError, err) {
		return
	}
}

// HandleLLMWebsocket opens a connection with the client and allows
// chat mode.
func HandleLLMWebsocket(pool *ClientPool, w http.ResponseWriter, r *http.Request) {
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

	//client, err := llm.NewClient()
	//if handleError(w, http.StatusInternalServerError, err) {
	//	return
	//}
	//
	//for {
	//	fmt.Println("Waiting for message...")
	//	subscribeRequest, err := parseMessage(conn)
	//	if err != nil {
	//		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
	//			handleError(w, http.StatusInternalServerError, err)
	//		}
	//		break
	//	}
	//
	//	if subscribeRequest.Type != string(TypeChatbot) {
	//		fmt.Println("ERROR: wrong type of request, expected chatbot, got: ", subscribeRequest.Type)
	//		break
	//	}
	//
	//	response, err := client.SendMessage(r.Context(), subscribeRequest.Data)
	//	if handleError(w, http.StatusInternalServerError, err) {
	//		break
	//	}
	//
	//	if err := conn.WriteJSON(response); handleError(w, http.StatusInternalServerError, err) {
	//		break
	//	}
	//}

	//defer conn.Close()
}

func registerClient(pool *ClientPool, conn *websocket.Conn) {
	client := &WSClient{
		ID:   uuid.New().String(),
		pool: pool,
		conn: conn,
		Send: make(chan []byte, 256),
		//Messages: []AnthropicMessage{},
	}

	client.pool.Register <- client
	go client.reader()
	go client.writer()
}

// parseMessage reads the first WS message
func parseMessage(conn *websocket.Conn) (subscribeRequest SubscribeRequest, err error) {
	_, msg, err := conn.ReadMessage()
	if err != nil {
		return subscribeRequest, err
	}

	if err := json.Unmarshal(msg, &subscribeRequest); err != nil {
		return subscribeRequest, err
	}

	return subscribeRequest, nil
}

// handleError write down an error with code to the writer response.
func handleError(w http.ResponseWriter, code int, err error) (hasError bool) {
	hasError = err != nil
	if hasError {
		log.Println("ERROR: ", err)
		http.Error(w, err.Error(), code)
	}
	return hasError
}
