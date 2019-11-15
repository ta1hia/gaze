package chat

import (
	"fmt"
	"log"
	"strings"
)

// TODO: room clean up
// TODO: make the lobby a room
// TODO: support private rooms

// A room is a named group of one or more users which will all
// receive messages addressed to that room.  The room is created
// implicitly when the first client joins it, and the room ceases to
// exist when the last client leaves it (TODO).  While rooms exists, any
// client can reference the room using the name of the room.
type Room struct {
	name string

	// Message queue
	mq chan Message

	// Map of all active connections. Maps user nick's to User.
	// TODO add/remove lock
	users map[string]*User

	done chan struct{}

	Commands CommandSet
}

// NewRoom initializes a new Room struct
func NewRoom(name string) *Room {

	cmdSet := CommandSet{}
	cmdSet.Add(&Command{
		Label: "/members",
		Usage: "List all members currently connected to this room",
		Handler: func(msg *Message, sender *User, v interface{}) {
			room := v.(*Room)
			systemMsg := Message{}
			var b strings.Builder
			for name := range room.users {
				fmt.Fprintf(&b, "%s\n", name)
			}
			systemMsg.Message = b.String()
			sender.Send(&systemMsg)
		},
	})
	cmdSet.Add(&Command{
		Label: "/exit",
		Usage: "Exit the current room",
		Handler: func(msg *Message, sender *User, v interface{}) {
			room := v.(*Room)
			room.RemoveUser(sender)
			systemMsg := Message{Message: fmt.Sprintf("Exiting '%s'", room.name)}
			sender.Send(&systemMsg)
		},
	})
	cmdSet.Add(&Help)

	return &Room{
		name:  name,
		users: make(map[string]*User),
		mq:    make(chan Message), // broadcast channel
		done:  make(chan struct{}),

		Commands: cmdSet,
	}
}

// AddUser adds a user to the map of currently connected users
func (r *Room) AddUser(u *User) error {
	r.users[u.nick] = u
	return nil
}

// RemoveUser removes a users from the room
func (r *Room) RemoveUser(u *User) error {
	delete(r.users, u.nick)
	u.room = nil
	// u.channel = lobby
	return nil
}

// RunRoom run dat
func (r *Room) Run() {
	for {
		msg := <-r.mq

		// If its a commands (ie /something), run its handler
		if msg.Command == "/help" {
			r.Commands["/help"].Handler(&msg, r.users[msg.Username], r.Commands)
		} else if msg.Command != "" {
			r.Commands.Dispatch(&msg, r.users[msg.Username], r)
		} else { // Otherwise its a regular message - broadcast as usual
			r.Broadcast(&msg)
		}
	}
}

// Broadcast a message to all users currently connected to the room
func (r *Room) Broadcast(msg *Message) {
	for nick, u := range r.users {
		err := u.conn.WriteJSON(msg)
		if err != nil {
			log.Printf("error: %v", err)
			u.conn.Close()
			delete(r.users, nick)
		}
	}
}
