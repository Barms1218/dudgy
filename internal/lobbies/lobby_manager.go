package lobbies

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

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
}

func NewLobbyManager(c context.Context) *LobbyManager {
	return &LobbyManager{
		ctx:     c,
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

func (l *LobbyManager) PlayerInLobby(id uuid.UUID) bool {
	_, exists := l.playerLobbies[id]
	return exists
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

	if roomCode == "" {
		roomCtx, roomCancel := context.WithTimeout(rm.ctx, 300*time.Second)
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
			return nil, fmt.Errorf("lobby code %s does not exists", roomCode)
		} else {
			if client.Cancel != nil {
				client.Cancel()
			}
		}
	}

	lobby.mu.Lock()
	defer lobby.mu.Unlock()

	if len(lobby.Players) == 4 {
		return nil, fmt.Errorf("Lobby %s is full.", roomCode)
	}

	lobby.Players[client.PlayerID] = client

	rm.playerLobbies[client.PlayerID] = roomCode

	return lobby, nil
}

func (l *LobbyManager) PreservePlayer(id uuid.UUID) {
	l.mu.Lock()
	lobby, exists := l.GetLobby(l.playerLobbies[id])
	if !exists {
		return
	}

	l.mu.Unlock()
	lobby.mu.Lock()

	ctx, cancel := context.WithTimeout(lobby.ctx, 30*time.Second)
	player, exists := lobby.Players[id]
	if !exists {
		cancel()
		return
	}

	player.Ctx = ctx
	player.Cancel = cancel

	lobby.mu.Unlock()

	go func() {
		<-ctx.Done()
		lobby.mu.Lock()
		defer lobby.mu.Unlock()
		if ctx.Err() == context.DeadlineExceeded {
			delete(lobby.Players, id)
		}
		lobby.mu.Unlock()
	}()
}

func (rm *LobbyManager) RemoveFromLobby(id uuid.UUID) (string, error) {
	lobby, exists := rm.GetLobby(rm.playerLobbies[id])
	if !exists {
		return "", fmt.Errorf("Lobby %s does not exist", lobby.Code)
	}

	rm.mu.Lock()
	defer rm.mu.Unlock()
	delete(rm.playerLobbies, id)

	lobby.mu.Lock()
	defer lobby.mu.Unlock()
	delete(lobby.Players, id)

	return lobby.Code, nil
}
