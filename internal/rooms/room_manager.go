package rooms

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"

	"github.com/google/uuid"
)

const MAX_PLAYERS = 4

type ROOM_STATE string

const (
	Exploration ROOM_STATE = "exploration"
	Combat      ROOM_STATE = "combat"
)

type Room struct {
	Code    string
	Players map[uuid.UUID]*RoomPlayer
	ctx     context.Context
	stop    context.CancelCauseFunc
	state   ROOM_STATE
	mu      sync.Mutex
	send    func(ids []uuid.UUID, msgType string, data json.RawMessage)
}

type RoomManager struct {
	rooms       map[string]*Room
	playerRooms map[uuid.UUID]string
	mu          sync.RWMutex
	send        func(ids []uuid.UUID, msgType string, data json.RawMessage)
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
		Players: make(map[uuid.UUID]*RoomPlayer, 0),
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

func (rm *RoomManager) JoinOrCreateRoom(roomCode string, client *RoomPlayer) (*Room, error) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	var room *Room
	var exists bool

	if roomCode == "" {
		newCode := generateRoomCode()
		room = &Room{
			Code:    newCode,
			Players: make(map[uuid.UUID]*RoomPlayer, 0),
			ctx:     context.WithTimeout(),
			state:   Exploration,
		}
		rm.rooms[newCode] = room
	} else {
		room, exists = rm.rooms[roomCode]
		if !exists {
			return nil, fmt.Errorf("room code %s does not exists", roomCode)
		}
	}

	room.mu.Lock()

	if len(room.Players) == 4 {
		return nil, fmt.Errorf("Room %s is full.", roomCode)
	}

	room.Players[client.PlayerID] = client
	room.mu.Unlock()

	rm.playerRooms[client.PlayerID] = roomCode

	return room, nil
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
