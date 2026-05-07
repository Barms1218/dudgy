package main

import (
	"encoding/json"
	"fmt"
	n "github.com/Barms1218/dudgy/internal/networking"
	t "github.com/Barms1218/dudgy/internal/types"
)

func (a *App) handleJoinLobby(client *n.Client, payload json.RawMessage) error {
	var joinedRoom n.JoinRoomPayload
	if err := json.Unmarshal(payload, &joinedRoom); err != nil {
		return fmt.Errorf("Invalid payload: %w", err)
	}

	roomPlayer := t.LobbyPlayer{
		PlayerID:    joinedRoom.PlayerID,
		Displayname: joinedRoom.DisplayName,
	}
	room, err := a.rm.JoinOrCreateLobby(joinedRoom.RoomCode.String(), &roomPlayer)
	if err != nil {
		return err
	}

	a.hub.Register <- client

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

func (a *App) handleLeaveLobby(client *n.Client, payload json.RawMessage) error {
	var disconnected n.PlayerLeftPayload
	if err := json.Unmarshal(payload, &disconnected); err != nil {
		return err
	}

	code, err := a.rm.RemoveFromLobby(disconnected.PlayerID)
	if err != nil {
		return err
	}

	data, err := json.Marshal(payload)

	type broadcast struct {
		Message string `json:"msg"`
	}

	message, err := json.Marshal(broadcast{Message: "You have been disconnected."})
	if err := a.sendToClient(disconnected.PlayerID, string(n.LeaveRoom), message); err != nil {
		return fmt.Errorf("Error handling leave lobby request: %w", err)
	}

	a.broadcast(code, string(n.LeaveRoom), data)

	return nil
}
