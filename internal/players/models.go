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
type PlayerInputPayload struct {
	X      int8   `json:"x"`
	Y      int8   `json:"y"`
	Action string `json:"action"`
}
