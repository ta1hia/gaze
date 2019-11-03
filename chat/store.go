package chat

import (
	"fmt"
)

// Store interface for storing chat-related state,
// such as rooms
type Store interface {
	Room(string) *Room
	AddRoom(r *Room) error
	StartRoom()
}

// InMemoryStore implements Store and provides an in-memory storage
type InMemoryStore struct {
	rooms map[string]*Room

	newRoom chan *Room // channel to listen on for new Room
}

// NewInMemoryStore initializer
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		rooms:   make(map[string]*Room),
		newRoom: make(chan *Room),
	}
}

// Room getter
func (s *InMemoryStore) Room(name string) *Room {
	room, _ := s.rooms[name]
	return room
}

// AddRoom adds a new room to the rooms map and notifies
func (s *InMemoryStore) AddRoom(r *Room) error {
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
func (s *InMemoryStore) StartRoom() {
	for {
		newRoom := <-s.newRoom
		go newRoom.Run()
	}
}
