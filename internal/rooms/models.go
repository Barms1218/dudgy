package rooms

import "github.com/google/uuid"

type Position struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
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

type JoinRoomPayload struct {
	PlayerID    uuid.UUID `json:"id"`
	RoomCode    uuid.UUID `json:"code"`
	DisplayName string    `json:"name"`
}

type RoomJoinResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
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
	X        int8      `json:"x"`
	Y        int8      `json:"y"`
}

type RunStartedPayload struct {
	Seed           int32            `json:"seed"`
	Map            map[string]any   `json:"map"`
	SpawnPositions []SpawnPositions `json:"spawn_positions"`
}

type WorldStatePayload struct {
	Tick    int32        `json:"tick"`
	Players []GamePlayer `json:"players"`
	Enemies []Enemy      `json:"enemies"`
}

type PlayerDisconnectedPayload struct {
	PlayerID uuid.UUID `json:"player_id"`
	RunSaved bool      `json:"run_saved"`
}

type RunResumedPayload struct {
	Seed    int32        `json:"seed"`
	Map     Map          `json:"map"`
	Players []GamePlayer `json:"players"`
	Enemies []Enemy      `json:"enemies"`
}

type Enemy struct {
	EnemyID uuid.UUID `json:"id"`
	X       int32     `json:"x"`
	Y       int32     `json:"y"`
	Health  int32     `json:"hp"`
}

type Map struct {
	Width  int32  `json:"map_width"`
	Height int32  `json:"map_height"`
	Tiles  []int8 `json:"tiles"`
}
