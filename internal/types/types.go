package types

import (
	"context"
	"sync"
	"time"
)

type ClassType string

const (
	Tank   ClassType = "tank"
	DPS    ClassType = "dps"
	Healer ClassType = "healer"
)

type AbilityType string

const (
	Auto   AbilityType = "auto"
	Melee  AbilityType = "aoe"
	Ranged AbilityType = "ranged"
)

type AbilityShape string

const (
	Circle    AbilityShape = "circle"
	Rectangle AbilityShape = "rectangle"
	Line      AbilityShape = "line"
)

type EffectType string

const (
	Heal  EffectType = "heal"
	Stun  EffectType = "stun"
	Bleed EffectType = "bleed"
	Burn  EffectType = "burn"
)

type Effect struct {
	ID        string     `json:"id"`
	Type      EffectType `json:"type"`
	Damage    int16      `json:"damage"`
	Start     time.Time  `json:"start"`
	ExpiresAt time.Time  `json:"end"`
}

type Ability struct {
	Name       string            `json:"name"`
	Type       AbilityType       `json:"type"`
	Damage     int16             `json:"damage"`
	TickDamage int16             `json:"tick_damage"`
	Shape      AbilityShape      `json:"shape"`
	Dimensions AbilityDimensions `json:"dimensions"`
	Duration   float32           `json:"duration"`
}

type AbilityDimensions struct {
	Height float32 `json:"height,omitempty"`
	Width  float32 `json:"width,omitempty"`
	Radius float32 `json:"radius,omitempty"`
}

type Entity struct {
	ID        string            `json:"id"`
	TeamID    string            `json:"team_id"`
	Class     ClassType         `json:"class"`
	Abilities []Ability         `json:"abilities"`
	Effects   map[string]Effect `json:"effects"`
	Position  Position          `json:"pos"`
	Velocity  Velocity          `json:"vel"`
	Health    int16             `json:"hp"`
	Mu        sync.Mutex
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
