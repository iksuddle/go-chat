package server

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/iksuddle/go-chat/internal/clients"
)

// map clients to empty struct for easy deleting and shi
type room struct {
	clients map[*clients.Client]struct{}
}

func (r *room) addClient(c *clients.Client) {
	r.clients[c] = struct{}{}
}

func (r *room) removeClient(c *clients.Client) {
	delete(r.clients, c)
}

// broadcast message to all connected clients except source
func (r *room) broadcast(msg string, source *clients.Client) {
	for c := range r.clients {
		// do not send message to sender
		if c == source {
			continue
		}

		err := c.Conn.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			log.Println(err)
			continue
		}
	}
}
