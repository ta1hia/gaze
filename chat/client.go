package chat

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

// TerminalUI for interfacing with gaze chat
type TerminalUI interface {

	// This should be a blocking listener that reads user input from the
	// terminal and then writes the input to the websocket conn that is
	// provided.
	ListenShell(conn *websocket.Conn, done chan bool)

	WriteShell(buf []byte) error
}

type GazeClient struct {
	term TerminalUI
	conn *websocket.Conn
	done chan bool
}

func NewGazeClient(conn *websocket.Conn, term TerminalUI) *GazeClient {
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
		var msg Message
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			log.Println("read:", err)
			return
		}
		s := fmt.Sprintf("%s: %s\n", msg.Username, msg.Message)
		c.term.WriteShell([]byte(s)) // Write the websocket msg to the terminal
	}
}

// Run the gaze client. This starts a go routine that listens for incoming
// messages on the websocket, and runs a blocking listener on the terminal shell
func (c *GazeClient) Run() {
	go c.ListenConnection()
	c.term.ListenShell(c.conn, c.done)
}
