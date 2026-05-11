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

	"github.com/coder/websocket"
	"github.com/google/uuid"
)

type App struct {
	logger  *slog.Logger
	l       *l.LobbyManager
	hub     *n.Hub
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

		_, exists := a.hub.Clients[id]
		if exists {
			a.hub.Clients[id] = conn
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
	if idStr == "" {
		id = uuid.NewString()
	} else {
		id = idStr
	}

	return id, nil
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

func (a *App) broadcast(roomCode string, msgType n.EnvelopeType, data json.RawMessage) error {
	envelope := n.Envelope{
		Type:    msgType,
		Payload: json.RawMessage(data),
	}
	out, err := json.Marshal(envelope)
	if err != nil {
		return err
	}

	room, exists := a.l.GetLobby(roomCode)
	if !exists {
		return err
	}

	ids := make([]string, 0, len(room.Players))
	for _, player := range room.Players {
		ids = append(ids, player.PlayerID)
	}
	a.hub.Broadcast <- n.BroadCastMessage{Recipients: ids, Payload: out}

	return nil
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
	a.funcMap[n.CreateLobby] = a.handleCreateLobby
	a.funcMap[n.ClassSelected] = a.handleClassSelection

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
