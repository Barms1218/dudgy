package main

import (
	"encoding/json"
	"fmt"

	n "github.com/Barms1218/dudgy/internal/networking"
	t "github.com/Barms1218/dudgy/internal/types"
	"github.com/google/uuid"
)

func (a *App) handleJoinLobby(id uuid.UUID, payload json.RawMessage) error {
	var joinedLobby n.JoinLobbyPayload
	if err := json.Unmarshal(payload, &joinedLobby); err != nil {
		return fmt.Errorf("Invalid payload: %w", err)
	}

	lobbyPlayer := t.LobbyPlayer{
		Displayname: joinedLobby.DisplayName,
	}

	var isPublic bool
	if joinedLobby.RoomCode.String() != "" {
		lobby, exists := a.l.GetLobby(joinedLobby.RoomCode.String())
		if !exists {
			return fmt.Errorf("Lobby %s does not exist", joinedLobby.RoomCode.String())
		}
		isPublic = lobby.IsPublic
	} else {
		isPublic = false
	}

	room, err := a.l.JoinOrCreateLobby(t.LobbyInfo{
		Code:     joinedLobby.RoomCode.String(),
		IsPublic: isPublic,
	}, &lobbyPlayer)
	if err != nil {
		return err
	}

	response := n.RoomJoinResponse{
		Success: true,
		Message: fmt.Sprintf("Welcome to the dungeon, %s!", lobbyPlayer.Displayname),
	}
	data, err := json.Marshal(response)
	if err := a.sendToClient(id, n.RoomJoined, data); err != nil {
		return fmt.Errorf("Error handling join room request: %w", err)
	}

	a.broadcast(room.Code, string(n.RoomJoined), data)

	return nil
}

func (a *App) handleLobbyVisibility(id uuid.UUID, payload json.RawMessage) error {
	var visibilityToggle n.LobbyVisibilityPayload
	if err := json.Unmarshal(payload, &visibilityToggle); err != nil {
		return err
	}

	lobby, exists := a.l.GetLobby(visibilityToggle.RoomCode)
	if !exists {
		return fmt.Errorf("Lobby %s does not exist", visibilityToggle.RoomCode)
	}

	if id != lobby.Owner {
		return fmt.Errorf("User not authorized to change this lobby")
	}

	if err := a.l.ToggleLobbyVisibility(visibilityToggle.RoomCode, visibilityToggle.IsPublic); err != nil {
		return err
	}

	return nil
}

func (a *App) handleLeaveLobby(id uuid.UUID, payload json.RawMessage) error {
	var disconnected n.PlayerLeftPayload
	if err := json.Unmarshal(payload, &disconnected); err != nil {
		return err
	}

	code, err := a.l.RemoveFromLobby(disconnected.PlayerID)
	if err != nil {
		return err
	}

	data, err := json.Marshal(payload)

	type broadcast struct {
		Message string `json:"msg"`
	}

	message, err := json.Marshal(broadcast{Message: "You have been disconnected."})
	if err := a.sendToClient(id, n.LeaveRoom, message); err != nil {
		return fmt.Errorf("Error handling leave lobby request: %w", err)
	}

	a.broadcast(code, string(n.LeaveRoom), data)

	return nil
}
