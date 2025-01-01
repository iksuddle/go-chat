package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

var clients map[*websocket.Conn]bool

func main() {
	clients = make(map[*websocket.Conn]bool)

	http.HandleFunc("/ws", serveWs)

	log.Println("starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	log.Println("incoming connection from", r.RemoteAddr)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("new connection", conn.RemoteAddr().String())

	clients[conn] = true

	go receiveMessages(conn)
	// sendMessages(conn)
}

// receive messages from conn
//
// todo: broadcast
func receiveMessages(conn *websocket.Conn) {
	addr := conn.RemoteAddr().String()
	addr = addr[len(addr)-5:]
	for {
		_, bytes, err := conn.ReadMessage()
		if err != nil {
			// conn prolly left
			fmt.Printf("<%s left>\n", addr)
			break
		}

		fmt.Printf("<%s> %s\n", addr, string(bytes))
		broadcast(bytes, conn, addr)
	}

	delete(clients, conn)
	conn.Close()
}

// send message to all connected clients except message sender
func broadcast(msg []byte, source *websocket.Conn, addr string) {
	msg = append([]byte(fmt.Sprintf("<%s> ", addr)), msg...)
	for c := range clients {
		if c == source {
			continue
		}

		err := c.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println(err)
			continue
		}
	}
}
