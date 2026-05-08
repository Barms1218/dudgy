package main

import (
	"encoding/json"

	n "github.com/Barms1218/dudgy/internal/networking"
	"github.com/google/uuid"
)

func (a *App) handleRegistration(id uuid.UUID, payload json.RawMessage) error {
	var registration n.RegisterPayload
	if err := json.Unmarshal(payload, &registration); err != nil {
		return err
	}

	if err := a.am.SetAccountName(id, registration.Name); err != nil {
		return err
	}

	return nil
}

func (a *App) handleReconnect(id uuid.UUID, msg json.RawMessage) error {

	return nil
}
