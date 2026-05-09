package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	a "github.com/Barms1218/dudgy/internal/accounts"
	l "github.com/Barms1218/dudgy/internal/lobbies"
	n "github.com/Barms1218/dudgy/internal/networking"

	"github.com/coder/websocket"
	"github.com/google/uuid"
)

type App struct {
	logger  *slog.Logger
	l       *l.LobbyManager
	hub     *n.Hub
	am      *a.AccountManager
	funcMap map[n.EnvelopeType]func(id string, payload json.RawMessage) error
}

func NewApp(logger *slog.Logger, manager *l.LobbyManager) *App {
	return &App{
		logger:  logger,
		l:       manager,
		hub:     n.NewHub(),
		funcMap: make(map[n.EnvelopeType]func(id string, payload json.RawMessage) error),
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

		idStr := r.URL.Query().Get("id")
		id, err := a.identifyUser(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := a.resolveIdentity(ctx, conn, id); err != nil {
			http.Error(w, "Unable to Resolve User Identity", http.StatusBadRequest)
			return
		}

		registration := &n.Registration{
			ID:   id,
			Conn: conn,
		}

		a.hub.Register <- registration

		defer func() {
			a.hub.Unregister <- registration
			conn.CloseNow()
		}()

		a.runSession(ctx, conn, registration)
	}

}

func (a *App) identifyUser(idStr string) (string, error) {
	var id string
	var err error
	if idStr == "" {
		id = uuid.NewString()
	} else {
		id = idStr
		if err != nil {
			a.logger.Error("Invalid id", "error", err)
			return "", err
		}

	}

	return id, nil
}

func (a *App) resolveIdentity(ctx context.Context, conn *websocket.Conn, id string) error {
	account := a.am.GetOrCreateAccount(id)
	if account.Name == "" {
		conn.Write(ctx, websocket.MessageText, []byte("Need username"))
		nameCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
		defer cancel()
		_, msg, err := conn.Read(nameCtx)
		if err != nil {
			if websocket.CloseStatus(err) != websocket.StatusNormalClosure {

			}
			return err
		}

		var envelope n.Envelope
		if err := json.Unmarshal(msg, &envelope); err != nil {
			a.logger.Error("Unable to parse name from JSON", "error", err)
			return err
		}

		if envelope.Type == n.Register {
			if err := a.handleRegistration(id, envelope.Payload); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("Invalid request.")
		}
	}

	return nil
}

func (a *App) runSession(ctx context.Context, conn *websocket.Conn, r *n.Registration) {
	for {
		readCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		_, msg, err := conn.Read(readCtx)
		cancel()
		if err != nil {
			if websocket.CloseStatus(err) != websocket.StatusNormalClosure {
				if exists := a.l.PlayerInLobby(r.ID); exists {
					a.l.PreservePlayer(r.ID)
				}
			}
			break
		}

		var envelope n.Envelope
		if err := json.Unmarshal(msg, &envelope); err != nil {
			a.logger.Error("Envelope unmarshaling failed", "error", err)
			continue
		}

		handleFunc, ok := a.funcMap[envelope.Type]
		if !ok {
			a.logger.Error("Command %v does not exist", "error", envelope.Type)
		}

		if err := handleFunc(r.ID, envelope.Payload); err != nil {

		}

	}
}

func (a *App) sendToClient(id string, msgType n.EnvelopeType, data json.RawMessage) error {
	envelope := n.Envelope{
		Type:    msgType,
		Payload: json.RawMessage(data),
	}
	out, err := json.Marshal(envelope)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return a.hub.Clients[id].Write(ctx, websocket.MessageText, out)
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

	ids := make([]string, 0, len(room.Players))
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
	a.funcMap[n.JoinRoom] = a.handleJoinLobby
	a.funcMap[n.PlayerLeft] = a.handleLeaveLobby
	a.funcMap[n.UpdateLobby] = a.handleLobbyVisibility
	a.funcMap[n.Reconnect] = a.handleReconnect
	a.funcMap[n.Register] = a.handleRegistration

	srv := &http.Server{
		Addr:    ":8080",
		Handler: http.HandlerFunc(a.handleWS(ctx)),
	}

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			a.logger.Error("Server Error", "error", err)
		}
	}()

	go a.hub.Run(ctx)

	<-ctx.Done()
	a.logger.Info("Shutting down, draining connections...")

	shutdownCtx, cancel := context.WithTimeout(bgCtx, 10*time.Second)
	defer cancel()
	srv.Shutdown(shutdownCtx)
}
