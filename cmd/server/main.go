package main

import (
	"log"
	"net/http"

	"github.com/iksuddle/go-chat/internal/server"
)

func main() {
	mux := http.NewServeMux()

	s := server.NewServer(mux)

	mux.HandleFunc("/create", s.CreateRoom)

	log.Println("starting server on port 8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", mux))
}
