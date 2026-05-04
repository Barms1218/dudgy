package main

import (
	"encoding/json"

	"github.com/coder/websocket"
	"github.com/google/uuid"
)

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

type Client struct {
	PlayerID uuid.UUID
	Conn     *websocket.Conn
}

type Envelope struct {
	Type    EnvelopeType    `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type Position struct {
	X float32 `json:"x_position"`
	Y float32 `json:"y_position"`
}

type Direction struct {
	X float32 `json:"x_direction"`
	Y float32 `json:"y_direction"`
}

type Player struct {
	PlayerID       uuid.UUID `json:"player_id"`
	PlayerPosition Position  `json:"player_position"`
	Health         int32     `json:"health"`
}

type RoomPlayer struct {
	PlayerID    uuid.UUID `json:"player_id"`
	Displayname string    `json:"display_name"`
	Ready       bool      `json:"ready"`
}

type Enemy struct {
	EnemyID       uuid.UUID `json:"enemy_id"`
	EnemyPosition Position  `json:"enemy_position"`
	Health        int32     `json:"health"`
}

type Map struct {
	Width  int32  `json:"map_width"`
	Height int32  `json:"map_height"`
	Tiles  []int8 `json:"tiles"`
}

type CreatePlayerReq struct {
	Name string `json:"name"`
}

type CreatePlayerResponse struct {
	PlayerID uuid.UUID `json:"player_id"`
	Name     string    `json:"name"`
}

type JoinRoomPayload struct {
	PlayerID    uuid.UUID `json:"player_id"`
	RoomCode    uuid.UUID `json:"room_code"`
	DisplayName string    `json:"display_name"`
}

type RoomJoinResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type PlayerInputPayload struct {
	Direction Direction `json:"direction"`
	Action    string    `json:"action"`
}

type PlayerInputResponse struct {
}

type DisconnectPaylaod struct {
	Save bool `json:"save"`
}

type RoomJoinedPayload struct {
	RoomCode     uuid.UUID    `json:"room_code"`
	YourPlayerID uuid.UUID    `json:"your_player_id"`
	Players      []RoomPlayer `json:"players"`
}

type SpawnPositions struct {
	PlayerID uuid.UUID `json:"player_id"`
	Position Position  `json:"position"`
}

type RunStartedPayload struct {
	Seed           int32            `json:"seed"`
	Map            map[string]any   `json:"map"`
	SpawnPositions []SpawnPositions `json:"spawn_positions"`
}

type WorldStatePayload struct {
	Tick    int32    `json:"tick"`
	Players []Player `json:"players"`
	Enemies []Enemy  `json:"enemies"`
}

type PlayerDisconnectedPayload struct {
	PlayerID uuid.UUID `json:"player_id"`
	RunSaved bool      `json:"run_saved"`
}

type RunResumedPayload struct {
	Seed    int32    `json:"seed"`
	Map     Map      `json:"map"`
	Players []Player `json:"players"`
	Enemies []Enemy  `json:"enemies"`
}

type GenericError struct {
	Code    string `json:"error_code"`
	Message string `json:"error_message"`
}
