package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"unicode"

	"github.com/gorilla/websocket"
	"github.com/iksuddle/go-chat/internal/clients"
	"github.com/iksuddle/go-chat/internal/messages"
)

type Server struct {
	Mux      *http.ServeMux
	upgrader *websocket.Upgrader
	rooms    map[string]*room
}

func NewServer(mux *http.ServeMux) *Server {
	return &Server{
		Mux:      mux,
		upgrader: &websocket.Upgrader{},
		rooms:    make(map[string]*room),
	}
}

// create a new room
func (s *Server) CreateRoom(w http.ResponseWriter, r *http.Request) {
	roomName := r.FormValue("roomName")

	if roomName == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("room name must be alphabetical"))
		return
	}

	if !isAlpha(roomName) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("room name must be alphabetical"))
		return
	}

	_, exists := s.rooms[roomName]
	if exists {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("room name must be alphabetical"))
		return
	}

	s.rooms[roomName] = &room{
		clients: make(map[*clients.Client]struct{}),
	}
	s.Mux.HandleFunc(fmt.Sprintf("/%s", roomName), func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(context.WithValue(r.Context(), "roomName", roomName))
		s.ServeWs(w, r)
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("room %s created", roomName)))
}

func (s *Server) ServeWs(w http.ResponseWriter, r *http.Request) {
	log.Println("incoming connection from", r.RemoteAddr)

	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("new connection", conn.RemoteAddr().String())

	roomName := r.Context().Value("roomName").(string)
	c := clients.NewClient(conn, roomName)

	s.rooms[roomName].addClient(c)

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
	s.rooms[c.Room].removeClient(c)
	c.Conn.Close()
}

// send message to all connected clients except message sender
func (s *Server) broadcast(msg string, source *clients.Client) {
	s.rooms[source.Room].broadcast(msg, source)
}

func isAlpha(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}
