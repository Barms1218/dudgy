package main

import (
	"context"
	"encoding/json"
	n "github.com/Barms1218/dudgy/internal/networking"
	r "github.com/Barms1218/dudgy/internal/rooms"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/coder/websocket"
	"github.com/google/uuid"
)

type App struct {
	rm      *r.RoomManager
	hub     *n.Hub
	funcMap map[string]func(*n.Client, json.RawMessage) error
}

func NewApp() *App {
	return &App{
		rm:  r.NewRoomManager(),
		hub: n.NewHub(),
	}
}

func (a *App) handleWS() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, nil)
		if err != nil {
			http.Error(w, "Accept Failed", http.StatusBadRequest)
			return
		}

		client := &n.Client{
			PlayerID: uuid.New(),
			Conn:     conn,
		}

		a.hub.Register <- client

		defer func() {
			a.hub.Unregister <- client
			conn.CloseNow()
		}()

		for {
			readCtx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
			_, msg, err := conn.Read(readCtx)
			cancel()
			if err != nil {
				return
			}

			var envelope n.Envelope
			if err := json.Unmarshal(msg, &envelope); err != nil {
				log.Printf("Malformed message from %s: %v", client.PlayerID, err)
				continue
			}

			handleFunc, ok := a.funcMap[string(envelope.Type)]
			if !ok {
				log.Printf("Command %v does not exist", envelope.Type)
			}

			if err := handleFunc(client, envelope.Payload); err != nil {

			}

		}
	}

}

func (a *App) RouteEnvelope(client *n.Client, msg json.RawMessage, e n.Envelope) {

}

func (a *App) sendToClient(id uuid.UUID, msgType string, data json.RawMessage) error {
	envelope := n.Envelope{
		Type:    n.EnvelopeType(msgType),
		Payload: json.RawMessage(data),
	}
	out, err := json.Marshal(envelope)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return a.hub.Clients[id].Conn.Write(ctx, websocket.MessageText, out)
}

func main() {
	bgCtx := context.Background()
	ctx, stop := signal.NotifyContext(bgCtx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	a := NewApp()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: http.HandlerFunc(a.handleWS()),
	}

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	go a.hub.Run(ctx)

	<-ctx.Done()
	log.Println("Shutting down, draining connections...")

	shutdownCtx, cancel := context.WithTimeout(bgCtx, 10*time.Second)
	defer cancel()
	srv.Shutdown(shutdownCtx)
}
