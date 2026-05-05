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
	PlayerID uuid.UUID `json:"id"`
	Name     string    `json:"name"`
}

type GamePlayer struct {
	PlayerPosition Position `json:"pos"`
	Health         int32    `json:"hp"`
}

type RoomPlayer struct {
	PlayerID    uuid.UUID `json:"player_id"`
	Displayname string    `json:"display_name"`
	Ready       bool      `json:"ready"`
}

type CreatePlayerReq struct {
	Name string `json:"name"`
}

type CreatePlayerResponse struct {
	PlayerID uuid.UUID `json:"player_id"`
	Name     string    `json:"name"`
}

type GenericError struct {
	Code    string `json:"error_code"`
	Message string `json:"error_message"`
}
