package client

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/tahia-khan/gaze/chat"
	"github.com/tahia-khan/gaze/client/terms"
)

type GazeClient struct {
	term terms.TerminalUI
	conn *websocket.Conn
	done chan bool
}

func NewGazeClient(conn *websocket.Conn, term terms.TerminalUI) *GazeClient {
	return &GazeClient{
		term: term,
		conn: conn,
		done: make(chan bool),
	}
}

// ListenConnection listens on the connection and writes any
// incoming lines to the terminal
func (c *GazeClient) ListenConnection() {
	// Display msgs
	for {
		var msg chat.Message
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			log.Println("read:", err)
			return
		}

		if msg.Username != "" {
			s := fmt.Sprintf("%s: %s\n", msg.Username, msg.Message)
			c.term.WriteShell([]byte(s)) // Write the websocket msg to the terminal
		} else {
			s := fmt.Sprintf("%s\n", msg.Message)
			c.term.WriteShell([]byte(s)) // Write the websocket msg to the terminal
		}
	}
}

// Run the gaze client. This starts a go routine that listens for incoming
// messages on the websocket, and runs a blocking listener on the terminal shell
func (c *GazeClient) Run() {
	go c.ListenConnection()
	c.term.ListenShell(c.conn, c.done)
}
