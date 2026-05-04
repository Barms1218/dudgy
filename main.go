package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/coder/websocket"
	"github.com/google/uuid"
)

func handleWS(hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, nil)
		if err != nil {
			http.Error(w, "Accept Failed", http.StatusBadRequest)
			return
		}

		client := &Client{
			PlayerID: uuid.New(),
			Conn:     conn,
		}

		hub.register <- client

		defer func() {
			hub.unregister <- client
			conn.CloseNow()
		}()

		for {
			readCtx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
			_, msg, err := conn.Read(readCtx)
			cancel()
			if err != nil {
				return
			}

			var envelope Envelope
			if err := json.Unmarshal(msg, &envelope); err != nil {
				log.Printf("Malformed message from %s: %v", client.PlayerID, err)
				continue
			}

			if err := ParseEnvelope(envelope, client, hub); err != nil {
				log.Printf("%v", err)
				continue
			}

		}
	}

}

func handleJoinRoom(client *Client, payload JoinRoomPayload, hub *Hub) error {
	if err := sendToClient(client, string(RoomJoined), RoomJoinedPayload{
		RoomCode:     payload.RoomCode,
		YourPlayerID: payload.PlayerID,
		Players:      payload.Players,
	}); err != nil {
		return fmt.Errorf("Error handling join room request: %w", err)
	}

	return nil
}

//	func handlePlayerInput[T any](client *Client, payload T, hub *Hub) error {
//		if err := sendToClient(client, string(PlayerInput), PlayerInputPayload{
//			Direction: payload.Direction,
//			Actions:   payload.Action,
//		}); err != nil {
//
//		}
//	}
func ParseEnvelope(envelope Envelope, client *Client, hub *Hub) error {
	switch envelope.Type {
	case JoinRoom:
		var payload JoinRoomPayload
		if err := json.Unmarshal(envelope.Payload, &payload); err != nil {
			return fmt.Errorf("Bad join_room payload: %w", err)
		}
		handleJoinRoom(client, payload, hub)
	case PlayerInput:
		var payload PlayerInputPayload
		if err := json.Unmarshal(envelope.Payload, &payload); err != nil {
			return fmt.Errorf("Bad player_input payload: %w", err)
		}
		handlePlayerInput(client, payload, hub)
	}
	return nil
}

func sendToClient[T any](client *Client, msgType string, payload T) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	envelope := Envelope{
		Type:    EnvelopeType(msgType),
		Payload: json.RawMessage(data),
	}
	out, err := json.Marshal(envelope)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return client.Conn.Write(ctx, websocket.MessageText, out)
}

func main() {
	bgCtx := context.Background()
	ctx, stop := signal.NotifyContext(bgCtx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	hub := NewHub()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: http.HandlerFunc(handleWS(hub)),
	}

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	go hub.Run(ctx)

	<-ctx.Done()
	log.Println("Shutting down, draining connections...")

	shutdownCtx, cancel := context.WithTimeout(bgCtx, 10*time.Second)
	defer cancel()
	srv.Shutdown(shutdownCtx)
}
