package networking

import (
	"context"
	"time"

	"github.com/coder/websocket"
	"github.com/google/uuid"
)

type Registration struct {
	ID   uuid.UUID `json:"id"`
	Conn *websocket.Conn
}

type Hub struct {
	Clients    map[uuid.UUID]*websocket.Conn
	Broadcast  chan BroadCastMessage
	Register   chan *Registration
	Unregister chan *Registration
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[uuid.UUID]*websocket.Conn, 0),
		Broadcast:  make(chan BroadCastMessage),
		Register:   make(chan *Registration),
		Unregister: make(chan *Registration),
	}
}

func (h *Hub) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			for _, c := range h.Clients {
				c.CloseNow()
			}
			return
		case client := <-h.Register:
			h.Clients[client.ID] = client.Conn
		case client := <-h.Unregister:
			delete(h.Clients, client.ID)
			client.Conn.CloseNow()
		case msg := <-h.Broadcast:
			var failed []uuid.UUID
			for i := range msg.Recipients {
				client, ok := h.Clients[msg.Recipients[i]]
				if !ok {
					continue
				}

				writeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				err := client.Write(writeCtx, websocket.MessageText, msg.Payload)
				cancel()
				if err != nil {
					failed = append(failed, msg.Recipients[i])
				}
			}
			for _, id := range failed {
				delete(h.Clients, id)
				h.Clients[id].CloseNow()
			}
		}
	}
}
