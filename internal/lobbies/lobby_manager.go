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

	// Context
	ctx  context.Context
	stop context.CancelFunc

	mu sync.Mutex
}

type LobbyManager struct {
	lobbies       map[string]*Lobby
	playerLobbies map[string]string
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
	defer l.mu.Unlock()
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
	b := make([]byte, 6)
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

func (l *LobbyManager) JoinOrCreateLobby(info t.LobbyInfo, client *t.LobbyPlayer) (*Lobby, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	var lobby *Lobby
	var exists bool

	if info.Code == "" {
		roomCtx, roomCancel := context.WithTimeout(l.ctx, 300*time.Second)
		newCode := generateLobbyCode()
		lobby = &Lobby{
			Owner:    client.PlayerID,
			Code:     newCode,
			IsPublic: info.IsPublic,
			Players:  make(map[string]*t.LobbyPlayer, 0),
			ctx:      roomCtx,
			stop:     roomCancel,
		}
		l.lobbies[newCode] = lobby
	} else {
		lobby, exists = l.lobbies[info.Code]
		if !exists {
			return nil, fmt.Errorf("lobby code %s does not exists", info.Code)
		} else {
			if client.Cancel != nil {
				client.Cancel()
			}
		}
	}

	lobby.mu.Lock()
	defer lobby.mu.Unlock()

	if len(lobby.Players) == 4 {
		return nil, fmt.Errorf("Lobby %s is full.", info.Code)
	}

	lobby.Players[client.PlayerID] = client

	l.playerLobbies[client.PlayerID] = info.Code

	return lobby, nil
}

func (l *LobbyManager) PreservePlayer(id string) (bool, error) {
	l.mu.Lock()
	lobby, exists := l.GetLobby(l.playerLobbies[id])
	if !exists {
		return exists, fmt.Errorf("No lobby exists for that player.")
	}

	l.mu.Unlock()
	lobby.mu.Lock()

	ctx, cancel := context.WithTimeout(lobby.ctx, 30*time.Second)
	player, exists := lobby.Players[id]
	if !exists {
		cancel()
		return exists, fmt.Errorf("That player does not exist.")
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

	player, exists = lobby.Players[id]
	return exists, nil
}

func (l *LobbyManager) RemoveFromLobby(id string) (string, error) {
	lobby, exists := l.GetLobby(l.playerLobbies[id])
	if !exists {
		return "", fmt.Errorf("Lobby %s does not exist", lobby.Code)
	}

	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.playerLobbies, id)

	lobby.mu.Lock()
	defer lobby.mu.Unlock()
	delete(lobby.Players, id)

	return lobby.Code, nil
}
