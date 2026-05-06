package lobbies

import (
	"context"
	"fmt"
	"math/rand"
	"sync"

	t "github.com/Barms1218/dudgy/internal/types"
	"github.com/google/uuid"
)

const MAX_PLAYERS = 4

type Lobby struct {
	Code    string
	Players map[uuid.UUID]*t.LobbyPlayer
	ctx     context.Context
	stop    context.CancelFunc
	mu      sync.Mutex
}

type LobbyManager struct {
	lobbies       map[string]*Lobby
	playerLobbies map[uuid.UUID]string
	mu            sync.RWMutex
	ctx           context.Context
	cancel        context.CancelFunc
}

func NewLobbyManager(c context.Context) *LobbyManager {
	ctx, cancel := context.WithCancel(c)
	return &LobbyManager{
		ctx:     ctx,
		cancel:  cancel,
		lobbies: make(map[string]*Lobby),
	}
}

func (rm *LobbyManager) CreateLobbies() *Lobby {
	code := generateLobbyCode()
	lobby := &Lobby{
		Code:    code,
		Players: make(map[uuid.UUID]*t.LobbyPlayer, 0),
	}
	rm.mu.Lock()
	rm.lobbies[code] = lobby
	rm.mu.Unlock()
	return lobby
}

func (rm *LobbyManager) GetLobby(code string) (*Lobby, bool) {
	rm.mu.RLock()
	defer rm.mu.Unlock()
	lobby, ok := rm.lobbies[code]
	return lobby, ok
}

func (rm *LobbyManager) DeleteLobby(code string) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	delete(rm.lobbies, code)
}

func generateLobbyCode() string {
	const letters = "ABCDEFGHIJKLMNPQRSTUVWXYZ23456789"
	b := make([]byte, 6)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

func (rm *LobbyManager) JoinOrCreateLobby(roomCode string, client *t.LobbyPlayer) (*Lobby, error) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	var lobby *Lobby
	var exists bool

	roomCtx, roomCancel := context.WithCancel(rm.ctx)
	if roomCode == "" {
		newCode := generateLobbyCode()
		lobby = &Lobby{
			Code:    newCode,
			Players: make(map[uuid.UUID]*t.LobbyPlayer, 0),
			ctx:     roomCtx,
			stop:    roomCancel,
		}
		rm.lobbies[newCode] = lobby
	} else {
		lobby, exists = rm.lobbies[roomCode]
		if !exists {
			roomCancel()
			return nil, fmt.Errorf("lobby code %s does not exists", roomCode)
		}
	}

	lobby.mu.Lock()

	if len(lobby.Players) == 4 {
		return nil, fmt.Errorf("Lobby %s is full.", roomCode)
	}

	lobby.Players[client.PlayerID] = client
	lobby.mu.Unlock()

	rm.playerLobbies[client.PlayerID] = roomCode

	return lobby, nil
}
