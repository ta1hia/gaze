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

// Room struct
type Room struct {
	name    string
	clients map[*websocket.Conn]bool
	mq      chan Message

	done chan struct{}
}

// User
type User struct {
	nick  string
	token string
}

// InMemoryStore provides an in-memory storage
// TODO implement this through an interface eventually
// so I can swap this out with redis, etc
type InMemoryStore struct {
	rooms map[string]*Room

	// channel to listen on for new Room
	newRoom chan *Room
}

// API for chat. TODO accept a storehandler
type API struct {
	store InMemoryStore
}

// CreateRoomHandler HTTP POST handler
func (api *API) CreateRoomHandler(w http.ResponseWriter, r *http.Request) {
	roomName := mux.Vars(r)["roomName"]
	log.Printf("CreateRoom: %s", roomName)
	if _, ok := api.store.rooms[roomName]; ok {
		http.Error(w, "Room already exists", http.StatusBadRequest)
		return
	}
	room := &Room{
		name:    roomName,
		clients: make(map[*websocket.Conn]bool), // connected clients
		mq:      make(chan Message),             // broadcast channel
		done:    make(chan struct{}),
	}

	api.store.rooms[roomName] = room
	api.store.newRoom <- room

	// var room Room
	// room.clients = make(map[*websocket.Conn]bool) // connected clients
	// room.mq = make(chan Message)                  // broadcast channel
	// api.store.rooms[roomName] = &room
	// api.store.newRoom <- &room
}

// JoinRoomHandler for POST room/join
func (api *API) JoinRoomHandler(w http.ResponseWriter, r *http.Request) {
	roomName := mux.Vars(r)["roomName"]
	log.Printf("JoinRoom: %s", roomName)
	room, ok := api.store.rooms[roomName]

	if !ok {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

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

// RunRoom run dat
// func RunRoom(room *Room) {
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

type Server struct {
	router *mux.Router
	api    *API
}

func NewServer() *Server {
	store := InMemoryStore{}
	store.rooms = make(map[string]*Room)
	store.newRoom = make(chan *Room)
	s := &Server{
		router: mux.NewRouter(),
		api:    &API{store},
	}

	s.router.HandleFunc("/{roomName}", s.api.CreateRoomHandler).Methods("POST")
	s.router.HandleFunc("/{roomName}/join", s.api.JoinRoomHandler)
	return s
}

func (s *Server) Serve(bind string) {
	// Broadcaster
	go func() {
		for {
			newRoom := <-s.api.store.newRoom
			go newRoom.Run() //RunRoom(newChan)
		}
	}()

	// Run WS server
	log.Fatal(http.ListenAndServe(bind, s.router))
}
