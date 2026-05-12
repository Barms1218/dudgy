package game

import (
	"context"
	"sync"
	"time"

	t "github.com/Barms1218/dudgy/internal/types"
	"github.com/google/uuid"
)

type Game struct {
	id             string
	players        map[string]*t.GamePlayer
	startPositions map[t.ClassType]t.Position
	checkpoints    []t.Position
	ctx            context.Context
	cancel         context.CancelFunc
	gameMap        t.Map
	mu             sync.Mutex
}

func (g *GameManager) CreateGame(ctx context.Context, players map[string]t.ClassType) *Game {
	gameCtx, gameCancel := context.WithTimeout(ctx, 60*time.Second)
	genMap, mapMiddle := GenerateMap(25, 25, 400)
	game := &Game{
		id:      uuid.NewString(),
		players: make(map[string]*t.GamePlayer, len(players)),
		ctx:     gameCtx,
		cancel:  gameCancel,
		gameMap: genMap,
	}
	game.mu.Lock()
	for id, class := range players {
		game.players[id] = &t.GamePlayer{
			PlayerID: id,
			Class:    class,
			Health:   100,
			Position: mapMiddle,
		}
	}

	game.mu.Unlock()
	g.mu.Lock()
	g.games[game.id] = game

	g.mu.Unlock()
	return game
}

type GameManager struct {
	games  map[string]*Game
	ctx    context.Context
	cancel context.CancelFunc
	mu     sync.Mutex
}

func NewGameManager(ctx context.Context) *GameManager {
	gameCtx, gameCancel := context.WithTimeout(ctx, 120*time.Second)
	return &GameManager{
		games:  make(map[string]*Game),
		ctx:    gameCtx,
		cancel: gameCancel,
	}
}
