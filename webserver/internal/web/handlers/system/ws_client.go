package system

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/knabben/observatio/webserver/internal/infra/llm"
	mcpaggregator "github.com/knabben/observatio/webserver/internal/infra/mcp"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type WSClient struct {
	ID      string
	conn    *websocket.Conn
	Send    chan []byte
	pool    *ClientPool
	service *llm.ObservationService

	// ctx is canceled once this connection's reader loop exits (client disconnect or read
	// error), so any Anthropic/kubectl call still in flight for this client is aborted instead
	// of running to completion against a connection nobody is listening on anymore.
	ctx    context.Context
	cancel context.CancelFunc

	// chatMu serializes chat turns on this connection: the UI only has one request in flight at
	// a time, but this guards ObservationService's conversation history against a race if a
	// second message ever arrives before the first reply finishes.
	chatMu sync.Mutex
}

// registerClient registers a new WebSocket client with the client pool and initializes its reader and writer goroutines.
// It creates a unique client ID, sets up a communication channel, and binds a new observation service for the client.
// Returns an error if the observation service fails to initialize.
func registerClient(pool *ClientPool, aggregator *mcpaggregator.Aggregator, conn *websocket.Conn) error {
	service, err := llm.NewObservationService(aggregator)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	client := &WSClient{
		ID:      uuid.New().String(),
		pool:    pool,
		conn:    conn,
		Send:    make(chan []byte, 256),
		service: service,
		ctx:     ctx,
		cancel:  cancel,
	}

	client.pool.Register <- client
	go client.reader()
	go client.writer()
	return nil
}

// reader continuously listens for WebSocket messages and hands each chatbot message off to the
// LLM client. Chat turns run in their own goroutine (see handleChat) so this loop keeps reading
// and can notice a client disconnect while a request is in flight.
func (c *WSClient) reader() {
	logger := log.FromContext(c.ctx)
	defer func() {
		// Unregister the client from the pool, cancel any in-flight request, and close the connection.
		c.pool.Unregister <- c
		c.cancel()
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

		go c.handleChat(message)
	}
}

// handleChat runs one streamed chat turn and forwards each chunk to the client as it arrives.
func (c *WSClient) handleChat(message *llm.ChatMessage) {
	c.chatMu.Lock()
	defer c.chatMu.Unlock()

	logger := log.FromContext(c.ctx)

	err := c.service.StreamChatWithAgent(c.ctx, message, c.sendMessage)
	if err != nil {
		// A failed LLM call (e.g. missing/invalid API key, upstream outage) must not drop
		// the WebSocket connection - send a safe, generic error message back to the client
		// (never the raw error, which may contain credential/endpoint details) and keep the
		// session open so the operator can retry once the server is reconfigured.
		logger.Error(err, "error requesting response from AI agent")
		c.sendMessage(llm.ToMessageParam("The AI assistant is not available right now. Please check the server's AI configuration and try again later."))
	}
}

// sendMessage marshals a chat message chunk and queues it for the writer goroutine.
func (c *WSClient) sendMessage(message *llm.ChatMessage) {
	logger := log.FromContext(c.ctx)

	result, err := json.Marshal(message)
	if err != nil {
		logger.Error(err, "error marshalling message")
		return
	}

	select {
	case c.Send <- result:
	default:
		close(c.Send)
		delete(c.pool.Clients, c.ID)
	}
}

// writer manages outgoing WebSocket messages, handles sending data from the Send channel,
// and maintains connection health.
func (c *WSClient) writer() {
	logger := log.FromContext(c.ctx)

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
