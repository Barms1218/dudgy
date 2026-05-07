package main

import (
	"encoding/json"
	"fmt"
	n "github.com/Barms1218/dudgy/internal/networking"
	t "github.com/Barms1218/dudgy/internal/types"
)

func (a *App) handleJoinLobby(client *n.Client, payload n.JoinRoomPayload) error {
	roomPlayer := t.LobbyPlayer{
		PlayerID:    payload.PlayerID,
		Displayname: payload.DisplayName,
		Ready:       false,
	}
	room, err := a.rm.JoinOrCreateLobby(payload.RoomCode.String(), &roomPlayer)
	if err != nil {

	}

	response := n.RoomJoinResponse{
		Success: true,
		Message: fmt.Sprintf("Welcome to the dungeon, %s!", roomPlayer.Displayname),
	}
	data, err := json.Marshal(response)
	if err := a.sendToClient(roomPlayer.PlayerID, string(n.RoomJoined), data); err != nil {
		return fmt.Errorf("Error handling join room request: %w", err)
	}

	a.broadcast(room.Code, string(n.RoomJoined), data)

	return nil
}

func (a *App) handleLeaveLobby(client *n.Client, payload n.PlayerDisconnectedPayload) error {
	code, err := a.rm.RemoveFromLobby(payload.PlayerID)
	if err != nil {

	}

	data, err := json.Marshal(payload)

	type broadcast struct {
		Message string `json:"msg"`
	}

	message, err := json.Marshal(broadcast{Message: "You have been disconnected."})
	if err := a.sendToClient(payload.PlayerID, string(n.LeaveRoom), message); err != nil {
		return fmt.Errorf("Error handling leave lobby request: %w", err)
	}

	a.broadcast(code, string(n.LeaveRoom), data)

	return nil
}
