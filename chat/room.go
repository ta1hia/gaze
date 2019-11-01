package chat

import (
	"log"

	"github.com/gorilla/websocket"
)

// Room struct
type Room struct {
	name    string
	clients map[*websocket.Conn]bool
	mq      chan Message

	done chan struct{}
}

// NewRoom initializes a new Room struct
func NewRoom(name string) *Room {
	return &Room{
		name:    name,
		clients: make(map[*websocket.Conn]bool), // connected clients
		mq:      make(chan Message),             // broadcast channel
		done:    make(chan struct{}),
	}
}

// RunRoom run dat
func (r *Room) Run() {

	for {
		msg := <-r.mq
		for c := range r.clients {
			err := c.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				c.Close()
				delete(r.clients, c)
			}
		}
	}
}
