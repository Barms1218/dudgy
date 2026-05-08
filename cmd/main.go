package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	l "github.com/Barms1218/dudgy/internal/lobbies"
	n "github.com/Barms1218/dudgy/internal/networking"
	t "github.com/Barms1218/dudgy/internal/types"

	"github.com/coder/websocket"
	"github.com/google/uuid"
)

type App struct {
	logger  *slog.Logger
	l       *l.LobbyManager
	hub     *n.Hub
	funcMap map[string]func(client *n.Client, payload json.RawMessage) error
}

func NewApp(logger *slog.Logger, manager *l.LobbyManager) *App {
	return &App{
		logger: logger,
		l:      manager,
		hub:    n.NewHub(),
	}
}

func (a *App) handleWS(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, nil)
		if err != nil {
			http.Error(w, "Accept Failed", http.StatusBadRequest)
			return
		}
		defer conn.CloseNow()

		newAccount := &t.Account{
			ID: uuid.New(),
		}
		client := &n.Client{
			Conn:    conn,
			Account: newAccount,
		}

		a.hub.Register <- client

		defer func() {
			a.hub.Unregister <- client
			conn.CloseNow()
		}()

		a.funcMap[string(n.JoinRoom)] = a.handleJoinLobby
		a.funcMap[string(n.PlayerLeft)] = a.handleLeaveLobby
		a.funcMap[string(n.UpdateLobby)] = a.handleLobbyVisibility
		a.funcMap[string(n.Reconnect)] = a.handleReconnect
		for {
			readCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
			_, msg, err := conn.Read(readCtx)
			cancel()
			if err != nil {
				if websocket.CloseStatus(err) != websocket.StatusNormalClosure {
					if exists := a.l.PlayerInLobby(client.Account.ID); exists {
						a.l.PreservePlayer(client.Account.ID)
					}
				}
				break
			}

			var envelope n.Envelope
			if err := json.Unmarshal(msg, &envelope); err != nil {

				continue
			}

			handleFunc, ok := a.funcMap[string(envelope.Type)]
			if !ok {
				a.logger.Error("Command %v does not exist", "error", envelope.Type)
			}

			if err := handleFunc(client, envelope.Payload); err != nil {

			}

		}
	}

}

func (a *App) handleReconnect(client *n.Client, msg json.RawMessage) error {

	return nil
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

func (a *App) broadcast(roomCode, msgType string, data json.RawMessage) {
	envelope := n.Envelope{
		Type:    n.EnvelopeType(msgType),
		Payload: json.RawMessage(data),
	}
	out, err := json.Marshal(envelope)
	if err != nil {
		return
	}

	room, exists := a.l.GetLobby(roomCode)
	if !exists {
		return
	}

	ids := make([]uuid.UUID, 0, len(room.Players))
	for _, player := range room.Players {
		ids = append(ids, player.PlayerID)
	}
	a.hub.Broadcast <- n.BroadCastMessage{Recipients: ids, Payload: out}
}

func main() {
	bgCtx := context.Background()
	ctx, stop := signal.NotifyContext(bgCtx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	a := NewApp(logger, l.NewLobbyManager(bgCtx))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: http.HandlerFunc(a.handleWS(ctx)),
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
