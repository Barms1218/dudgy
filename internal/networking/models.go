package networking

import (
	"encoding/json"

	"github.com/coder/websocket"
	"github.com/google/uuid"
)

type Client struct {
	PlayerID uuid.UUID `json:"id"`
	Conn     *websocket.Conn
}

type BroadCastMessage struct {
	Recipients []uuid.UUID `json:"ids"`
	Payload    []byte
}

type EnvelopeType string

const (
	JoinRoom           EnvelopeType = "joined_room"
	PlayerInput        EnvelopeType = "player_input"
	LeaveRoom          EnvelopeType = "leave_room"
	RoomJoined         EnvelopeType = "room_joined"
	PlayerJoined       EnvelopeType = "player_joined"
	RunStarted         EnvelopeType = "run_started"
	WorldState         EnvelopeType = "world_state"
	PlayerDisconnected EnvelopeType = "player_disconnected"
	RunResumed         EnvelopeType = "run_resumed"
	Error              EnvelopeType = "error"
)

type Envelope struct {
	Type    EnvelopeType    `json:"type"`
	Payload json.RawMessage `json:"payload"`
}
