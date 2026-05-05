package players

import (
	"github.com/coder/websocket"
	"github.com/google/uuid"
)

type Player struct {
	PlayerID uuid.UUID       `json:"id"`
	Name     string          `json:"name"`
	Conn     *websocket.Conn `json:"conn"`
}

type GamePlayer struct {
	PlayerID uuid.UUID `json:"id"`
	X        int32     `json:"x"`
	Y        int32     `json:"y"`
	Health   int32     `json:"hp"`
}

type RoomPlayer struct {
	PlayerID    uuid.UUID `json:"player_id"`
	Displayname string    `json:"display_name"`
	Ready       bool      `json:"ready"`
}

type PlayerInputPayload struct {
	X      int8   `json:"x"`
	Y      int8   `json:"y"`
	Action string `json:"action"`
}
