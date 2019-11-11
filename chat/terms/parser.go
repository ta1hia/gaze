package terms

import "strings"

// Message struct
type Message struct {
	// Email    string `json:"email"`
	Username string `json:"username"`
	Command  string `json:"command"` // The command to run, eg "nick", "exit"
	Message  string `json:"message"` // The message body
}

// ParseTerminalMessage parse user input from a terminal into
// a Message struct
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
