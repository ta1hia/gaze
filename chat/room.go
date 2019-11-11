package chat

import (
	"log"

	"github.com/gorilla/websocket"
)

// TODO: room clean up
// TODO: make the lobby a room
// TODO: support private rooms
type Channel interface {
	SetMessageHandler(handler func(sender *User, msg *Message, v interface{}))
}

// A room is a named group of one or more users which will all
// receive messages addressed to that room.  The room is created
// implicitly when the first client joins it, and the room ceases to
// exist when the last client leaves it (TODO).  While rooms exists, any
// client can reference the room using the name of the room.
type Room struct {
	name    string
	clients map[*websocket.Conn]bool
	users   map[string]*User // map of user nick's to User object, ie all active websocket connections
	mq      chan Message

	done chan struct{}

	Commands CommandSet //map[string]*Command
}

// NewRoom initializes a new Room struct
func NewRoom(name string) *Room {

	// Initialize the set of commands for this room
	// TODO probably make this a global in commands.go since
	// basic rooms will all have the same commands
	cmdSet := CommandSet{}
	cmdSet.Add(&Help)

	return &Room{
		name:  name,
		users: make(map[string]*User),
		mq:    make(chan Message), // broadcast channel
		done:  make(chan struct{}),

		Commands: cmdSet,
	}
}

// func (r *Room) AddUser(nick string, c *websocket.Conn) error {
func (r *Room) AddUser(u *User) error {
	r.users[u.nick] = u
	return nil
}

// RunRoom run dat
func (r *Room) Run() {
	for {
		msg := <-r.mq

		// If its a commands (ie /something), run its handler
		if msg.Command != "" {
			cmd, ok := r.Commands[msg.Command]
			if !ok { // TODO error for invalid command
				continue
			}
			cmd.Handler(&msg, r)
		} else { // Otherwise broadcast as usual
			r.Broadcast(&msg)
		}
	}
}

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
