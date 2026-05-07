package networking

import (
	"encoding/json"

	t "github.com/Barms1218/dudgy/internal/types"
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
	JoinRoom     EnvelopeType = "joined_room"
	PlayerInput  EnvelopeType = "player_input"
	LeaveRoom    EnvelopeType = "leave_room"
	RoomJoined   EnvelopeType = "room_joined"
	PlayerJoined EnvelopeType = "player_joined"
	RunStarted   EnvelopeType = "run_started"
	WorldState   EnvelopeType = "world_state"
	PlayerLeft   EnvelopeType = "player_left"
	RunResumed   EnvelopeType = "run_resumed"
	Error        EnvelopeType = "error"
)

type Envelope struct {
	Type    EnvelopeType    `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type DisconnectPayload struct {
	Save bool `json:"save"`
}

type RoomJoinedPayload struct {
	RoomCode        uuid.UUID         `json:"room_code"`
	YourPlayerID    uuid.UUID         `json:"your_player_id"`
	DisplayName     string            `json:"display_name"`
	ExistingPlayers map[uuid.UUID]any `json:"players"`
}

type RunStartedPayload struct {
	Seed int32          `json:"seed"`
	Map  map[string]any `json:"map"`
	X    int8           `json:"x"`
	Y    int8           `json:"y"`
}

type WorldStatePayload struct {
	Tick    int32          `json:"tick"`
	Players []t.GamePlayer `json:"players"`
	Enemies []t.Enemy      `json:"enemies"`
}

type PlayerLeftPayload struct {
	PlayerID    uuid.UUID `json:"id"`
	DisplayName string    `json:"name"`
	RunSaved    bool      `json:"run_saved"`
}

type RunResumedPayload struct {
	Seed     int32     `json:"seed"`
	PlayerID uuid.UUID `json:"id"`
}

type RunResumedResponse struct {
	Seed    int32                 `json:"seed"`
	Map     map[string]any        `json:"map"`
	Players map[uuid.UUID][]int16 `json:"players"` // Players maps player ID to [x, y, hp]
	Enemies map[uuid.UUID][]int16 `json:"enemies"` // Enemies maps enemy ID to [x, y, hp]
}

type JoinRoomPayload struct {
	PlayerID    uuid.UUID `json:"id"`
	RoomCode    uuid.UUID `json:"code"`
	DisplayName string    `json:"name"`
}

type RoomJoinResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type PlayerInputPayload struct {
	X      int8   `json:"x"`
	Y      int8   `json:"y"`
	Action string `json:"action"`
}

type GenericError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}
