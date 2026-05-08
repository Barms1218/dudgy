package types

import (
	"context"
	"github.com/coder/websocket"
	"github.com/google/uuid"
)

type Account struct {
	ID   uuid.UUID
	Name string
}

type Enemy struct {
	EnemyID uuid.UUID `json:"id"`
	Pos     Position  `json:"pos"`
	Health  int8      `json:"hp"`
}

type Player struct {
	PlayerID uuid.UUID       `json:"id"`
	Name     string          `json:"name"`
	Conn     *websocket.Conn `json:"conn"`
}

type Map struct {
	Width  int16  `json:"map_width"`
	Height int16  `json:"map_height"`
	Tiles  []int8 `json:"tiles"`
}

type SpawnPositions struct {
	PlayerID uuid.UUID `json:"player_id"`
	Pos      Position  `json:"pos"`
}

type Position struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

type Velocity struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

type GamePlayer struct {
	PlayerID uuid.UUID `json:"id"`
	Pos      Position  `json:"pos"`
	Health   int8      `json:"hp"`
}

type LobbyPlayer struct {
	PlayerID    uuid.UUID `json:"player_id"`
	Displayname string    `json:"display_name"`
	Ctx         context.Context
	Cancel      context.CancelFunc
}

type LobbyInfo struct {
	Code     string `json:"code"`
	IsPublic bool   `json:"public"`
}
