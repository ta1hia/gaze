package chat

import (
	"fmt"
	"strings"
)

type Command struct {
	Label string
	Usage string

	Handler func(*Message, *User, interface{})
}

var Help = Command{
	Label: "/help",
	Usage: "Print all available commands in the room",
	Handler: func(msg *Message, sender *User, v interface{}) {
		commands := v.(CommandSet)
		systemMsg := Message{}
		var b strings.Builder
		fmt.Fprintf(&b, "Available commands:\n")
		for label, cmd := range commands {
			fmt.Fprintf(&b, "\n%10s\t%s", label, cmd.Usage)
		}
		systemMsg.Message = b.String()
		sender.Send(&systemMsg)
	},
}

// CommandSet is a collection of commands mapped by the
// command name. It describes all available commands for a room.
type CommandSet map[string]*Command

// Add a command to the set
func (cs CommandSet) Add(cmd *Command) {
	cs[cmd.Label] = cmd
}

// Dispatch the message command to the correct handler within the CommandSet
func (cs CommandSet) Dispatch(msg *Message, sender *User, v interface{}) {
	cmd, ok := cs[msg.Command]
	if !ok {
		systemMsg := Message{Message: fmt.Sprintf("Unrecognized command %s", msg.Command)}
		sender.Send(&systemMsg)
	} else {
		cmd.Handler(msg, sender, v)
	}
}
