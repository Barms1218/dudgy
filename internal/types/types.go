package types

import (
	"context"
	"github.com/coder/websocket"
)

type ClassType string

const (
	Tank   ClassType = "tank"
	DPS    ClassType = "dps"
	Healer ClassType = "healer"
)

type Account struct {
	ID   string
	Name string
}

type GamePlayer struct {
	PlayerID string    `json:"id"`
	Class    ClassType `json:"class"`
	Position Position  `json:"pos"`
	Velocity Velocity  `json:"vel"`
	Health   int16     `json:"hp"`
}

type Enemy struct {
	EnemyID string   `json:"id"`
	Pos     Position `json:"pos"`
	Health  int8     `json:"hp"`
}

type Player struct {
	PlayerID string          `json:"id"`
	Name     string          `json:"name"`
	Conn     *websocket.Conn `json:"conn"`
}

type Map struct {
	Width  int16  `json:"map_width"`
	Height int16  `json:"map_height"`
	Tiles  []int8 `json:"tiles"`
}

type SpawnPositions struct {
	PlayerID string   `json:"player_id"`
	Pos      Position `json:"pos"`
}

type Position struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

type Velocity struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

type LobbyPlayer struct {
	PlayerID string    `json:"player_id"`
	Class    ClassType `json:"class"`
	Ctx      context.Context
	Cancel   context.CancelFunc
	Ready    bool `json:"ready"`
}

type LobbyInfo struct {
	OwnerID  string `json:"owner"`
	Code     string `json:"code"`
	IsPublic bool   `json:"is_public"`
	Name     string `json:"name"`
}
