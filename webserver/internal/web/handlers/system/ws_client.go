package system

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/knabben/observatio/webserver/internal/infra/llm"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type WSClient struct {
	ID      string
	conn    *websocket.Conn
	Send    chan []byte
	pool    *ClientPool
	service *llm.ObservationService
}

// registerClient registers a new WebSocket client with the client pool and initializes its reader and writer goroutines.
// It creates a unique client ID, sets up a communication channel, and binds a new observation service for the client.
// Returns an error if the observation service fails to initialize.
func registerClient(pool *ClientPool, conn *websocket.Conn) error {
	service, err := llm.NewObservationService()
	if err != nil {
		return err
	}

	client := &WSClient{
		ID:      uuid.New().String(),
		pool:    pool,
		conn:    conn,
		Send:    make(chan []byte, 256),
		service: service,
	}

	client.pool.Register <- client
	go client.reader()
	go client.writer()
	return nil
}

// reader continuously listens for WebSocket messages, processes them,
// and handles message exchange with the LLM client.
func (c *WSClient) reader() {
	var (
		ctx    = context.Background()
		logger = log.FromContext(ctx)
	)
	defer func() {
		// Unregister the client from the pool and close the connection.
		c.pool.Unregister <- c
		c.conn.Close() // nolint
	}()

	for {
		message, err := parseMessage(c.conn)
		if err != nil {
			logger.Error(err, "error parsing message")
			break
		}

		logger.Info("Received message", "msg", message, "client", c.ID)

		if message.Type != string(TypeChatbot) {
			logger.Error(nil, "Wrong type of request, expected chatbot", "type", message.Type)
			break
		}

		// Start to chat with the bot agent, sending the first message.
		response, err := (*c.service).ChatWithAgent(ctx, message, c.ID)
		if err != nil {
			logger.Error(err, "error writing close message")
			return
		}

		result, err := json.Marshal(response)
		if err != nil {
			logger.Error(err, "error writing close message")
			return
		}

		select {
		// Send the bot response to the client writer goroutine.
		case c.Send <- result:
		default:
			close(c.Send)
			delete(c.pool.Clients, c.ID)
		}
	}
}

// writer manages outgoing WebSocket messages, handles sending data from the Send channel,
// and maintains connection health.
func (c *WSClient) writer() {
	var (
		ctx    = context.Background()
		logger = log.FromContext(ctx)
	)

	ticker := time.NewTicker(60 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close() // nolint
	}()

	for {
		select {
		case message, ok := <-c.Send:
			logger.Info("Sending message", "msg", message, "client", c.ID)

			c.conn.SetWriteDeadline(time.Now().Add(60 * time.Second)) // nolint
			if !ok {
				if err := c.conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					logger.Error(err, "error when trying to write a close message")
					return
				}
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				logger.Error(err, "error writing the next text message")
				return
			}

			w.Write(message) // nolint
			if err := w.Close(); err != nil {
				logger.Error(err, "error closing the writter.")
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(60 * time.Second)) // nolint
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
