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

	// Room that the user is currently connected to. If nil, we
	// assume that the user is connected to the lobby
	room *Room
}

func NewUser(nick string, c *websocket.Conn) *User {
	return &User{
		nick: nick,
		conn: c,
	}
}

func (u *User) Send(msg *Message) error {
	return u.conn.WriteJSON(msg)
}
