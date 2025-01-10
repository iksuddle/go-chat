package clients

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

var reader = bufio.NewReader(os.Stdin)

type Client struct {
	Name string
	// the server
	Conn *websocket.Conn
	Room string
}

func NewClient(conn *websocket.Conn, room string) *Client {
	addr := conn.RemoteAddr().String()
	addr = addr[len(addr)-5:]

	return &Client{
		Name: addr,
		Conn: conn,
		Room: room,
	}
}

func (c *Client) Start() {
	go c.receiveMessages()
	c.sendMessages()
}

// get messages from stdin and send them to the server
func (c *Client) sendMessages() {
	for {
		msg, err := getMessage()
		if err != nil {
			log.Println(err)
			continue
		}

		err = c.Conn.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			log.Println(err)
			continue
		}
	}
}

// print messages sent from server
func (c *Client) receiveMessages() {
	for {
		_, bytes, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println("server disconnected")
			break
		}

		printMessage(string(bytes))
	}

	c.Conn.Close()
}

// get message from stdin
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
