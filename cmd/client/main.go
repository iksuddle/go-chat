package main

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/gorilla/websocket"
	"github.com/iksuddle/go-chat/internal/clients"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("provide room name")
		os.Exit(1)
	}
	room := os.Args[1]

	u := url.URL{
		Scheme: "ws",
		// Host:   "192.168.100.46:8080",
		Host: "localhost:8080",
		Path: fmt.Sprintf("/%s", room),
	}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal(err)
	}

	// todo: dont do ""
	c := clients.NewClient(conn, "")
	c.Start()
}
