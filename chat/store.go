package chat

import (
	"fmt"
)

// Store struct for storing chat-related state,
// such as rooms
type Store struct {
	rooms map[string]*Room

	newRoom chan *Room // channel to listen on for new Room
}

// NewStore initializer
func NewStore() *Store {
	return &Store{
		rooms:   make(map[string]*Room),
		newRoom: make(chan *Room),
	}
}

// Room getter
func (s *Store) Room(name string) *Room {
	room, _ := s.rooms[name]
	return room
}

// AddRoom adds a new room to the rooms map and notifies
func (s *Store) AddRoom(r *Room) error {
	_, ok := s.rooms[r.name]
	if ok {
		return fmt.Errorf("Room %s already exists", r.name)
	}
	s.rooms[r.name] = r
	s.newRoom <- r // does this need a wait
	return nil
}

// StartRoom listens for newly created rooms and
// starts them. Should be run in a goroutine.
func (s *Store) StartRoom() {
	for {
		newRoom := <-s.newRoom
		go newRoom.Run()
	}
}
