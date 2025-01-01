package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/gorilla/websocket"
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

	// send and receive messages
	go receiveMessages(conn)
	sendMessages(conn)
}

func sendMessages(conn *websocket.Conn) {
	for {
		msg, err := getMessage()
		if err != nil {
			log.Println(err)
			continue
		}

		err = conn.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			log.Println(err)
			continue
		}
	}
}

func receiveMessages(conn *websocket.Conn) {
	for {
		_, bytes, err := conn.ReadMessage()
		if err != nil {
			log.Println("server disconnected")
			break
		}

		printMessage(string(bytes))
	}

	conn.Close()
}

var reader = bufio.NewReader(os.Stdin)

func getMessage() (string, error) {
	fmt.Print("> ")
	msg, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	msg = strings.TrimSpace(msg)
	if len(msg) == 0 {
		return "", errors.New("cannot send empty message")
	}

	return msg, nil
}

func printMessage(msg string) {
	fmt.Print("\b\b")
	fmt.Println(msg)
	fmt.Print("> ")
}
