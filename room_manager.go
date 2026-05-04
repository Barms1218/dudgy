package main

import (
	"math/rand"
	"sync"
)

type Room struct {
	Code    string
	Players []*Client
	mu      sync.Mutex
}

type RoomManager struct {
	rooms map[string]*Room
	mu    sync.RWMutex
}

func NewRoomManager() *RoomManager {
	return &RoomManager{
		rooms: make(map[string]*Room),
	}
}

func (rm *RoomManager) CreateRoom() *Room {
	code := generateRoomCode()
	room := &Room{
		Code:    code,
		Players: make([]*Client, 0),
	}
	rm.mu.Lock()
	rm.rooms[code] = room
	rm.mu.Unlock()
	return room
}

func (rm *RoomManager) GetRoom(code string) (*Room, bool) {
	rm.mu.RLock()
	defer rm.mu.Unlock()
	room, ok := rm.rooms[code]
	return room, ok
}

func (rm *RoomManager) DeleteRoom(code string) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	delete(rm.rooms, code)
}

func generateRoomCode() string {
	const letters = "ABCDEFGHIJKLMNPQRSTUVWXYZ23456789"
	b := make([]byte, 6)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}
