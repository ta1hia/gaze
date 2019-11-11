package chat

import (
	"fmt"
	"strings"
)

// type CommandDispatcher interface {
// 	SetHandlerFunc(handler func(sender *User, msg *Message, v interface{}))
// }

type Command struct {
	Label string
	Usage string

	Handler func(*Message, interface{})
}

var Help = Command{
	Label: "/help",
	Usage: "Print all available commands in the room",
	Handler: func(msg *Message, v interface{}) {
		room := v.(*Room)

		systemMsg := Message{}
		var b strings.Builder
		for label, cmd := range room.Commands {
			fmt.Fprintf(&b, "Available commands:\n\n%s\t%s\n", label, cmd.Usage)
		}
		systemMsg.Message = b.String()

		// Help msg should go back to the sender
		recipient := room.users[msg.Username]
		recipient.Send(&systemMsg)
	},
}

// More to come!

//var DirectMessage
//var Exit
//var Members

// CommandSet is a collection of commands mapped by the
// command name. It describes all available commands for a room.
type CommandSet map[string]*Command

// Add a command to the set
func (cs CommandSet) Add(cmd *Command) {
	cs[cmd.Label] = cmd
}
