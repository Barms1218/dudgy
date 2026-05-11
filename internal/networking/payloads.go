package networking

import (
	"encoding/json"

	t "github.com/Barms1218/dudgy/internal/types"
)

type BroadCastMessage struct {
	Recipients []string `json:"ids"`
	Payload    []byte
}

type EnvelopeType string

const (
	JoinRoom         EnvelopeType = "joined_room"
	PlayerInput      EnvelopeType = "player_input"
	Register         EnvelopeType = "register"
	LeaveRoom        EnvelopeType = "leave_room"
	Reconnect        EnvelopeType = "reconnect"
	ClassSelected    EnvelopeType = "class_selected"
	CreateLobby      EnvelopeType = "create_lobby"
	RunStarted       EnvelopeType = "run_started"
	PlayerReady      EnvelopeType = "player_ready"
	WorldState       EnvelopeType = "world_state"
	PlayerLeft       EnvelopeType = "player_left"
	ToggleVisibility EnvelopeType = "toggle_visilibity"
	RunResumed       EnvelopeType = "run_resumed"
	Error            EnvelopeType = "error"
)

type Envelope struct {
	Type    EnvelopeType    `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type DisconnectPayload struct {
	Save bool `json:"save"`
}

type RoomJoinedPayload struct {
	RoomCode        string         `json:"room_code"`
	YourPlayerID    string         `json:"your_player_id"`
	DisplayName     string         `json:"display_name"`
	ExistingPlayers map[string]any `json:"players"`
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
	PlayerID    string `json:"id"`
	DisplayName string `json:"name"`
	RunSaved    bool   `json:"run_saved"`
}

type LobbyVisibilityPayload struct {
	RoomCode string `json:"code"`
	IsPublic bool   `json:"is_public"`
}

type ReconnectPayload struct {
	ID string `json:"id"`
}

type RunResumedPayload struct {
	Seed     int32  `json:"seed"`
	PlayerID string `json:"id"`
}

type RunResumedResponse struct {
	Seed    int32              `json:"seed"`
	Map     map[string]any     `json:"map"`
	Players map[string][]int16 `json:"players"`
	Enemies map[string][]int16 `json:"enemies"`
}

type JoinLobbyPayload struct {
	PlayerID string `json:"id"`
	RoomCode string `json:"code"`
}

type CreateLobbyPayload struct {
	OwnerID   string `json:"id"`
	LobbyName string `json:"name"`
	IsPublic  bool   `json:"is_public"`
}

type RoomJoinResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type RegisterPayload struct {
	Name string `json:"string"`
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

type SelectClassPayload struct {
	ID    string `json:"id"`
	Room  string `json:"room"`
	Class string `json:"class"`
}

type SelectClassResponse struct {
	Message string `json:"msg"`
	Success bool   `json:"success"`
}

type ToggleReadyPayload struct {
	LobbyCode string `json:"code"`
}

type ToggleReadyResponse struct {
	Ready   bool   `json:"ready"`
	Message string `json:"msg"`
}
