package system

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type WSClient struct {
	ID   string
	conn *websocket.Conn
	Send chan []byte
	pool *ClientPool
	//Messages []AnthropicMessage
}

func (c *WSClient) reader() {
	var logger = log.FromContext(context.Background())
	defer func() {
		c.pool.Unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Error(err, "error reading from websocket")
			}
			break
		}

		var wsMsg WSMessage
		if err := json.Unmarshal(message, &wsMsg); err != nil {
			logger.Error(err, "error unmarshaling message")
			continue
		}

		// Process the message
		//go c.handleMessage(wsMsg)
		logger.Info("messsage ", "msg", wsMsg)
		c.sendMessage("chatbot", "Hello, I'm a bot!")
	}
}

func (c *WSClient) writer() {
	var logger = log.FromContext(context.Background())
	ticker := time.NewTicker(10 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				if err := c.conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					logger.Error(err, "error writing close message")
					return
				}
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *WSClient) sendMessage(msgType string, content string) {
	var logger = log.FromContext(context.Background())
	msg := WSMessage{
		ID:        uuid.New().String(),
		Type:      msgType,
		Content:   content,
		Timestamp: time.Now().Unix(),
	}

	data, err := json.Marshal(msg)
	if err != nil {
		logger.Error(err, "error marshaling message")
		return
	}

	select {
	case c.Send <- data:
	default:
		close(c.Send)
		delete(c.pool.Clients, c.ID)
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
