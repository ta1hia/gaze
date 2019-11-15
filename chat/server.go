package chat

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// Websocket upgrader instance
var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
	return true
}}

// Gaze is the container for a gaze chat server
type Gaze struct {
	router   *mux.Router
	store    *Store
	commands CommandSet
}

func NewGaze() *Gaze {

	// Create "lobby" command set. The lobby represents the default
	// channel that users get dropped into when they are not actively
	// connected to a room
	cmdSet := CommandSet{}
	cmdSet.Add(&Command{ // Add join
		Label: "/join",
		Usage: "Join a room",
		Handler: func(msg *Message, sender *User, v interface{}) {
			store := v.(*Store)

			room := store.Room(msg.Message) // If the room doesn't exist, create it
			if room == nil {
				room = NewRoom(msg.Message)
				store.AddRoom(room)
			}
			room.AddUser(sender)
			sender.room = room

			systemMsg := Message{Message: fmt.Sprintf("joining '%s'", msg.Message)}
			sender.Send(&systemMsg)
		},
	})
	cmdSet.Add(&Command{ // Add list
		Label: "/list",
		Usage: "lists all rooms",
		Handler: func(msg *Message, sender *User, v interface{}) {
			store := v.(*Store)

			systemMsg := Message{}
			var b strings.Builder
			for name := range store.rooms {
				fmt.Fprintf(&b, "%s\n", name)
			}
			systemMsg.Message = b.String()
			sender.Send(&systemMsg)
		},
	})
	cmdSet.Add(&Command{ // Add nick
		Label: "/nick",
		Usage: "Change current nickname",
		Handler: func(msg *Message, sender *User, v interface{}) {
			sender.nick = msg.Message
			systemMsg := Message{Message: fmt.Sprintf("Setting nickname to '%s'", msg.Message)}
			sender.Send(&systemMsg)
		},
	})
	cmdSet.Add(&Help) // Add help

	server := &Gaze{
		router:   mux.NewRouter(),
		store:    NewStore(),
		commands: cmdSet,
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
			g.HandleMessage(&msg, &u) // Send the msg to the lobby handler
		}
	}

}

func (g *Gaze) HandleMessage(msg *Message, u *User) {
	if msg.Command == "/help" {
		g.commands["/help"].Handler(msg, u, g.commands)
	} else if msg.Command != "" {
		g.commands.Dispatch(msg, u, g.store)
	} else {
		u.Send(msg)
	}
}

func (s *Gaze) Serve(bind string) {

	// Listener for starting new rooms. Message broadcasts are
	// handled via each room.
	go s.store.StartRoom()

	// Start the server
	log.Fatal(http.ListenAndServe(bind, s.router))
}
