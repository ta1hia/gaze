package chat

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// Configure the upgrader
var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
	return true
}}

// Gaze is the container for a gaze chat server
type Gaze struct {
	router *mux.Router
	store  *Store
}

func NewGaze() *Gaze {
	server := &Gaze{
		router: mux.NewRouter(),
		store:  NewStore(),
	}

	// Set up the websocket handler
	server.router.HandleFunc("/connect", server.ConnectToRoomHandler)
	return server
}

func (g *Gaze) ConnectToRoomHandler(w http.ResponseWriter, r *http.Request) {

	// Upgrade the HTTP connection to a websocket connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error: %v", err)
		ws.Close()
		return
	}
	defer ws.Close()

	// Create a user, which holds the connection
	u := User{conn: ws}

	// Message handler loop
	//
	// A user is either connected to a room OR in the lobby. If user is
	// connected to a room, let the message be handled by the room. Otherwise
	// let the 'lobby handler' handle the message
	for {
		var msg Message
		err = ws.ReadJSON(&msg)
		if u.room != nil { // Send the msg to the client's room
			u.room.mq <- msg
		} else {
			g.HandleMessage(msg, &u) // Send the msg to the lobby
		}
	}
}

// Would be good to implement this like how routers set handlers on the fly
// This should send out System messages (from the system directly
// to the client)
func (g *Gaze) HandleMessage(msg Message, u *User) {

	// Regex for parsing out opts
	// r, _ := regexp.Compile("(?i)^(\\w+)\\s(.+)")

	systemMsg := Message{Username: "lobby"}

	// TODO add "help"
	switch msg.Command {
	case "/nick": // Changes user's nick name
		// TODO handle from the client side
		u.nick = msg.Message
		systemMsg.Message = fmt.Sprintf("nickname changed to %s", u.nick)
	case "/list": // List all rooms
		// List of all rooms in g.store.rooms
		systemMsg.Message = fmt.Sprintf("this should print all rooms")
	case "/join": // Join a room

		room := g.store.Room(msg.Message) // If the room doesn't exist, create it
		if room == nil {
			room = NewRoom(msg.Message)
			g.store.AddRoom(room)
		}

		// Do a room.AddUser(u) so the user gets added to the room
		room.AddUser(u)
		room.clients[u.conn] = true

		// Set u.room
		u.room = room
		systemMsg.Message = fmt.Sprintf("joining '%s' room", msg.Message)
	default:
		systemMsg = msg
	}

	u.conn.WriteJSON(systemMsg)
}

func (s *Gaze) Serve(bind string) {

	// Listener for starting new rooms. Message broadcasts are
	// handled via each room.
	go s.store.StartRoom()

	// Start the server
	log.Fatal(http.ListenAndServe(bind, s.router))
}
