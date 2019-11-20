package terms

import (
	"strings"

	"github.com/gorilla/websocket"
	"github.com/tahia-khan/gaze/chat"
)

// TerminalUI for interfacing with gaze chat
type TerminalUI interface {

	// This should be a blocking listener that reads user input from the
	// terminal and then writes the input to the websocket conn that is
	// provided.
	ListenShell(conn *websocket.Conn, done chan bool)

	WriteShell(buf []byte) error
}

// ParseTerminalMessage parse user input from a terminal into
// a Message struct
// TODO throw out blank msgs, raise error
func ParseTerminalMessage(line, nick string) *chat.Message {
	if strings.HasPrefix(line, "/") {
		matches := strings.SplitN(line, " ", 2)
		msg := chat.Message{Username: nick, Command: matches[0]}
		if len(matches) > 1 {
			msg.Message = matches[1]
		}
		return &msg
	}
	return &chat.Message{Username: nick, Message: line}
}
