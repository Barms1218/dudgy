package lobbies

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	t "github.com/Barms1218/dudgy/internal/types"
)

const MAX_PLAYERS = 4

type Lobby struct {
	Owner    string
	Code     string
	IsPublic bool
	Players  map[string]*t.LobbyPlayer
	Name     string

	// Context
	ctx  context.Context
	stop context.CancelFunc

	mu sync.Mutex
}

type LobbyManager struct {
	lobbies       map[string]*Lobby
	playerLobbies map[string]string // Maps a player's id to a lobby code
	mu            sync.RWMutex
	ctx           context.Context
}

func NewLobbyManager(c context.Context) *LobbyManager {
	return &LobbyManager{
		ctx:     c,
		lobbies: make(map[string]*Lobby),
	}
}

func (l *LobbyManager) CreateLobbies() *Lobby {
	code := generateLobbyCode()
	lobby := &Lobby{
		Code:    code,
		Players: make(map[string]*t.LobbyPlayer, 0),
	}
	l.mu.Lock()
	l.lobbies[code] = lobby
	l.mu.Unlock()
	return lobby
}

func (l *LobbyManager) GetLobby(code string) (*Lobby, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	lobby, ok := l.lobbies[code]
	return lobby, ok
}

func (l *LobbyManager) PlayerInLobby(id string) bool {
	_, exists := l.playerLobbies[id]
	return exists
}

func (l *LobbyManager) DeleteLobby(code string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.lobbies, code)
}

func generateLobbyCode() string {
	const letters = "ABCDEFGHIJKLMNPQRSTUVWXYZ23456789"
	b := make([]byte, 10)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

func (l *LobbyManager) ToggleLobbyVisibility(roomCode string, isPublic bool) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	lobby, exists := l.GetLobby(roomCode)
	if !exists {
		return fmt.Errorf("Room %s does not exist", roomCode)
	}

	lobby.IsPublic = isPublic

	return nil
}

func (l *LobbyManager) GetPublicLobbies() []string {
	l.mu.RLock()
	defer l.mu.RUnlock()

	var public []string
	for code, lobby := range l.lobbies {
		if lobby.IsPublic {
			public = append(public, code)
		}
	}
	return public
}

func (l *LobbyManager) CreateLobby(info t.LobbyInfo, client *t.LobbyPlayer) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	roomCtx, roomCancel := context.WithTimeout(l.ctx, 300*time.Second)
	lobby := &Lobby{
		Owner:    client.PlayerID,
		Code:     generateLobbyCode(),
		IsPublic: info.IsPublic,
		Players:  make(map[string]*t.LobbyPlayer, 0),
		ctx:      roomCtx,
		stop:     roomCancel,
	}

	l.lobbies[lobby.Code] = lobby

	lobby.mu.Lock()
	defer lobby.mu.Unlock()

	lobby.Players[client.PlayerID] = client

	return nil
}

func (l *LobbyManager) RemoveFromLobby(id string) error {
	lobby, exists := l.GetLobby(l.playerLobbies[id])
	if !exists {
		return fmt.Errorf("Lobby %s does not exist", lobby.Code)
	}

	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.playerLobbies, id)

	lobby.mu.Lock()
	defer lobby.mu.Unlock()
	delete(lobby.Players, id)

	return nil
}

func (l *LobbyManager) SelectClass(id, code string, class t.ClassType) error {
	lobby, exists := l.GetLobby(l.playerLobbies[id])
	if !exists {
		return fmt.Errorf("Lobby %s does not exist.", code)
	}

	lobby.mu.Lock()
	defer lobby.mu.Unlock()

	var claimed bool
	requestingPlayer, exists := lobby.Players[id]
	if !exists {
		return fmt.Errorf("No valid id sent with request")
	}

	for _, player := range lobby.Players {
		if player.Class == class && lobby.IsPublic {
			claimed = true
			break
		}
	}

	if !claimed {
		requestingPlayer.Class = class
	}

	return nil
}
