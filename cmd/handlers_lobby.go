package main

import (
	"encoding/json"
	"fmt"

	n "github.com/Barms1218/dudgy/internal/networking"
	t "github.com/Barms1218/dudgy/internal/types"
)

func (a *App) handleJoinLobby(id string, payload json.RawMessage) error {
	var joinedLobby n.JoinLobbyPayload
	if err := json.Unmarshal(payload, &joinedLobby); err != nil {
		return fmt.Errorf("Invalid payload: %w", err)
	}
	if joinedLobby.RoomCode != "" {
		lobby := a.l.GetLobby(joinedLobby.RoomCode)
		if lobby == nil {
			return fmt.Errorf("Lobby %s does not exist", joinedLobby.RoomCode)
		}
	}

	response := n.RoomJoinResponse{
		Success: true,
		Message: "Welcome to the lobby!",
	}
	data, err := json.Marshal(&response)
	if err != nil {
		return err
	}
	return a.sendToClient(id, n.JoinRoom, json.RawMessage(data))
}

func (a *App) handleCreateLobby(id string, payload json.RawMessage) error {
	var createdLobby n.CreateLobbyPayload
	if err := json.Unmarshal(payload, &createdLobby); err != nil {
		return err
	}

	if err := a.l.CreateLobby(t.LobbyInfo{
		IsPublic: createdLobby.IsPublic,
		OwnerID:  id,
		Name:     createdLobby.LobbyName,
	}, &t.LobbyPlayer{PlayerID: id}); err != nil {
		return err
	}

	response := n.RoomJoinResponse{
		Success: true,
		Message: "Welcome to the lobby!",
	}
	data, err := json.Marshal(&response)
	if err != nil {
		return err
	}

	return a.sendToClient(id, n.JoinRoom, json.RawMessage(data))
}

func (a *App) handleClassSelection(id string, payload json.RawMessage) error {
	var info n.SelectClassPayload
	if err := json.Unmarshal(payload, &info); err != nil {
		return err
	}

	err := a.l.SelectClass(id, info.Room, t.ClassType(info.Class))
	if err != nil {
		msg := n.SelectClassResponse{
			Success: false,
			Message: fmt.Sprintf("%s is already claimed.", info.Class),
		}

		data, err := json.Marshal(&msg)
		if err != nil {
			return err
		}

		return a.sendToClient(id, n.ClassSelected, data)
	}

	broadcast := n.SelectClassResponse{
		Message: fmt.Sprintf("%s has been claimed.", info.Class),
		Success: err == nil,
	}
	data, err := json.Marshal(&broadcast)
	if err != nil {
		return err
	}

	return a.broadcast(info.Room, n.ClassSelected, json.RawMessage(data))

}

func (a *App) handleLobbyVisibility(id string, payload json.RawMessage) error {
	var visibilityToggle n.LobbyVisibilityPayload
	if err := json.Unmarshal(payload, &visibilityToggle); err != nil {
		return err
	}

	lobby := a.l.GetLobby(visibilityToggle.RoomCode)
	if lobby == nil {
		return fmt.Errorf("Lobby %s does not exist", visibilityToggle.RoomCode)
	}

	if id != lobby.Owner {
		return fmt.Errorf("User not authorized to change this lobby")
	}

	return a.l.ToggleLobbyVisibility(visibilityToggle.RoomCode, visibilityToggle.IsPublic)
}

func (a *App) handleLeaveLobby(id string, payload json.RawMessage) error {
	var disconnected n.PlayerLeftPayload
	if err := json.Unmarshal(payload, &disconnected); err != nil {
		return err
	}

	err := a.l.RemoveFromLobby(disconnected.PlayerID)
	if err != nil {
		return err
	}

	return nil
}
