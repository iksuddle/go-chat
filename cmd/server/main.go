package main

import (
	"log"
	"net/http"

	"github.com/iksuddle/go-chat/internal/server"
)

func main() {
	s := server.NewServer()

	http.HandleFunc("/ws", s.ServeWs)

	log.Println("starting server on port 8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
