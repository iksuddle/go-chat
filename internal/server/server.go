package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/iksuddle/go-chat/internal/clients"
	"github.com/iksuddle/go-chat/internal/messages"
)

type Server struct {
	clients  map[*clients.Client]bool
	upgrader *websocket.Upgrader
}

func NewServer() *Server {
	return &Server{
		clients:  make(map[*clients.Client]bool),
		upgrader: &websocket.Upgrader{},
	}
}

func (s *Server) ServeWs(w http.ResponseWriter, r *http.Request) {
	log.Println("incoming connection from", r.RemoteAddr)

	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("new connection", conn.RemoteAddr().String())

	c := clients.NewClient(conn)

	s.clients[c] = true

	go s.receiveMessages(c)
}

// receive messages from conn and broadcast to all others
func (s *Server) receiveMessages(c *clients.Client) {
	for {
		_, bytes, err := c.Conn.ReadMessage()
		if err != nil {
			// conn prolly left
			fmt.Println(messages.GetLeaveMessage(c.Name))
			break
		}

		msg := messages.GetMessageFrom(bytes, c.Name)
		fmt.Println(msg)
		s.broadcast(msg, c)
	}

	// client left
	delete(s.clients, c)
	c.Conn.Close()
}

// send message to all connected clients except message sender
func (s *Server) broadcast(msg string, source *clients.Client) {
	for c := range s.clients {
		// do not broadcast a message to the sender
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
