package system

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/log"
)

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
