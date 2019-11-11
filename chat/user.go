package chat

import (
	"github.com/gorilla/websocket"
)

// User represents a connected user that is actively
// connected. Currently users don't persist beyond a connection
// and are tied to an active websocket connection
type User struct {
	// uid  string

	// A client can change their nickname while in the lobby,
	// but not once they join a room
	nick string

	conn *websocket.Conn

	room *Room
}

func NewUser(nick string, c *websocket.Conn) *User {
	return &User{
		// uid:  uuid.New().String(),
		nick: nick,
		conn: c,
	}
}

func (u *User) Send(msg *Message) error {
	return u.conn.WriteJSON(msg)
}
