package main

import (
	"log"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/iksuddle/go-chat/internal/clients"
)

func main() {
	u := url.URL{
		Scheme: "ws",
		Host:   "localhost:8080",
		Path:   "/ws",
	}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal(err)
	}

	c := clients.NewClient(conn)
	c.Start()
}
