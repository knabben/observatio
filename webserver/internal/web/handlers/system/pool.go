package system

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
	"github.com/knabben/observatio/webserver/internal/infra/llm"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type WSClient struct {
	ID        string
	conn      *websocket.Conn
	Send      chan []byte
	pool      *ClientPool
	LLMClient *llm.Client
}

func (c *WSClient) reader() {
	ctx := context.Background()
	var logger = log.FromContext(ctx)
	defer func() {
		c.pool.Unregister <- c
		c.conn.Close() // nolint
	}()

	for {
		message, err := parseMessage(c.conn)
		if err != nil {
			logger.Error(err, "error parsing message")
			break
		}

		logger.Info("Received message", "msg", message)
		if message.Type != string(TypeChatbot) {
			logger.Error(nil, "ERROR: wrong type of request, expected chatbot", "msg", message.Type)
			break
		}

		logger.Info("Received message", "msg", message)
		response, err := (*c.LLMClient).SendMessage(ctx, message.Content)
		if err != nil {
			logger.Error(err, "error writing close message")
			return
		}

		response.AgentID = c.ID
		result, err := json.Marshal(response)
		if err != nil {
			logger.Error(err, "error writing close message")
			return
		}
		select {
		case c.Send <- result:
		default:
			close(c.Send)
			delete(c.pool.Clients, c.ID)
		}
	}
}

func (c *WSClient) writer() {
	var logger = log.FromContext(context.Background())
	ticker := time.NewTicker(60 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close() // nolint
	}()

	for {
		select {
		case message, ok := <-c.Send:
			logger.Info("Sending message", "msg", message, "ok", ok)
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

type ClientPool struct {
	Clients    map[string]*WSClient
	Broadcast  chan []byte
	Register   chan *WSClient
	Unregister chan *WSClient
}

func (c *ClientPool) Run(ctx context.Context) {
	var logging = log.FromContext(ctx)
	for {
		select {
		case client := <-c.Register:
			c.Clients[client.ID] = client
			logging.Info("Client connected.", "client", client.ID)
		case client := <-c.Unregister:
			if _, ok := c.Clients[client.ID]; ok {
				delete(c.Clients, client.ID)
				close(client.Send)
				logging.Info("Client disconnected", "client", client.ID)
			}
		case message := <-c.Broadcast:
			for clientID, client := range c.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(c.Clients, clientID)
				}
			}
		}
	}
}
