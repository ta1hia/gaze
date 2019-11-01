package chat

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// Configure the upgrader
var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
	return true
}}

// // User
// type User struct {
// 	nick  string
// 	token string
// }

// Gaze is the container for a gaze chat server
type Gaze struct {
	router *mux.Router
	store  Store
}

func NewGaze() *Gaze {
	server := &Gaze{
		router: mux.NewRouter(),
		store:  NewInMemoryStore(),
	}

	// Set up the websocket handler
	server.router.HandleFunc("/{roomName}/join", server.ConnectToRoomHandler)
	return server
}

func (g *Gaze) ConnectToRoomHandler(w http.ResponseWriter, r *http.Request) {
	roomName := mux.Vars(r)["roomName"]
	room := g.store.Room(roomName)
	log.Printf("ConnectToRoom: %s", roomName)

	if room == nil { // If the room doesn't exist, create it
		room = NewRoom(roomName)
		g.store.AddRoom(room)
	}

	// Upgrade the HTTP connection to a websocket connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error: %v", err)
		ws.Close()
		return
	}
	defer ws.Close()
	room.clients[ws] = true

	// Message handler loop
	for {
		var msg Message
		err = ws.ReadJSON(&msg)
		log.Printf("JoinRoom: %v: %v", msg.Username, msg.Message)
		if len(msg.Message) > 0 {
			room.mq <- msg
		}
	}
}

func (s *Gaze) Serve(bind string) {

	// Listener for starting new rooms. Message broadcasts are
	// handled via each room.
	go s.store.StartRoom()

	// Start the server
	log.Fatal(http.ListenAndServe(bind, s.router))
}
