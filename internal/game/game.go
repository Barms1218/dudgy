package game

import (
	"context"
	"fmt"
	"sync"
	"time"

	t "github.com/Barms1218/dudgy/internal/types"
)

type Game struct {
	id             string
	playersID      string
	enemiesID      string
	entities       map[string]*t.Entity
	startPositions map[t.ClassType]t.Position
	checkpoints    []t.Position
	ctx            context.Context
	cancel         context.CancelFunc
	gameMap        t.Map
	mu             sync.Mutex
}

func (g *Game) Run() {
	ticker := time.NewTicker(33 * time.Millisecond)

	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-g.ctx.Done():
				return
			default:

			}
		}
	}()
}

func (g *Game) UpdatePlayerPosition(id string, newPos t.Position) error {

	return nil
}

func (g *Game) PerformAttack(id string, ability *t.Ability) error {
	switch ability.Shape {
	case t.AbilityShape("circle"):

	case t.AbilityShape("rectangle"):

	case t.AbilityShape("line"):

	}

	return nil
}

func (g *Game) TakeDamage(id string, ability *t.Ability) error {
	g.mu.Lock()
	entity, exists := g.entities[id]
	if !exists {
		return fmt.Errorf("No such entity: %s", id)
	}
	g.mu.Unlock()

	entity.Mu.Lock()
	defer entity.Mu.Unlock()

	entity.Health -= ability.Damage

	return nil
}
