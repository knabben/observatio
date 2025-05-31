package llm

import "github.com/gorilla/websocket"

type WSClient struct {
	ID   string
	conn *websocket.Conn
	send chan []byte
	pool *ClientPool
	//messages []AnthropicMessage
}

type ClientPool struct {
	Clients    map[*WSClient]bool
	Broadcast  chan []byte
	Register   chan *WSClient
	Unregister chan *WSClient
}

func (c *ClientPool) Run() {

}
