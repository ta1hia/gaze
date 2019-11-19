package terms

import (
	"strings"

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

// Message struct - repeated
type Message struct {
	// Email    string `json:"email"`
	Username string `json:"username"`
	Command  string `json:"command"` // The command to run, eg "nick", "exit"
	Message  string `json:"message"` // The message body
}

// ParseTerminalMessage parse user input from a terminal into
// a Message struct
// TODO throw out blank msgs, raise error
func ParseTerminalMessage(line, nick string) *Message {
	var msg Message

	if strings.HasPrefix(line, "/") {
		matches := strings.SplitN(line, " ", 2)
		msg = Message{Username: nick, Command: matches[0]}
		if len(matches) > 1 {
			msg.Message = matches[1]
		}
	} else {
		msg = Message{Username: nick, Message: line}
	}
	return &msg
}
