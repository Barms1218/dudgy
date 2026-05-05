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

func GenerateMap(width, height, seed int32) (Map, Position) {
	tiles := make([]int8, width*height)
	for i := range tiles {
		tiles[i] = 1
	}

	r := rand.New(rand.NewSource(int64(seed)))
	x, y := int(width/2), int(height/2)
	steps := int(width * height / 3)
	dirs := [][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}

	for i := 0; i < steps; i++ {
		tiles[y*int(width)+x] = 0
		d := dirs[r.Intn(4)]
		nx, ny := x+d[0], y+d[1]
		if nx > 0 && nx < int(width)-1 && ny > 0 && ny < int(height)-1 {
			x, y = nx, ny
		}
	}

	spawn := Position{X: float32(width / 2), Y: float32(height / 2)}
	return Map{Width: width, Height: height, Tiles: tiles}, spawn
}
