package main

import (
	"encoding/json"
	"fmt"
	n "github.com/Barms1218/dudgy/internal/networking"
	r "github.com/Barms1218/dudgy/internal/rooms"
)

func (a *App) handleJoinRoom(client *n.Client, payload r.JoinRoomPayload) error {
	roomPlayer := r.RoomPlayer{
		PlayerID:    payload.PlayerID,
		Displayname: payload.DisplayName,
		Ready:       false,
	}
	room, err := a.rm.JoinOrCreateRoom(payload.RoomCode.String(), &roomPlayer)
	if err != nil {

	}

	room.Players[client.PlayerID] = &roomPlayer

	response := r.RoomJoinResponse{
		Success: true,
		Message: "Welcome to the dungeon!",
	}
	data, err := json.Marshal(response)
	if err := a.sendToClient(roomPlayer.PlayerID, string(n.RoomJoined), data); err != nil {
		return fmt.Errorf("Error handling join room request: %w", err)
	}

	return nil
}
