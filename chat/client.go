package chat

import (
	"github.com/gorilla/websocket"
	"log"
)

type Client struct {
	ws     *websocket.Conn
	sendCh chan []byte
}

func NewClient(ws *websocket.Conn) *Client {
	return &Client{
		ws:     ws,
		sendCh: make(chan []byte),
	}
}

func (c *Client) ReadLoop(hub *Hub) {
	defer func() {
		hub.UnregisterCh <- c
		c.ws.Close()
	}()
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			log.Printf("error: %v", err)
			break
		}
		hub.BroadcastCh <- message
	}
}

func (c *Client) WriteLoop() {
	defer c.ws.Close()
	for {
		message, ok := <-c.sendCh
		if !ok {
			c.ws.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}
		c.ws.WriteMessage(websocket.TextMessage, message)
	}
}
