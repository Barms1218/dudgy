package networking

import (
	"context"
	"time"

	"github.com/coder/websocket"
	"github.com/google/uuid"
)

type Hub struct {
	Clients    map[uuid.UUID]*Client
	Broadcast  chan BroadCastMessage
	Register   chan *Client
	Unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[uuid.UUID]*Client, 0),
		Broadcast:  make(chan BroadCastMessage),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (h *Hub) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			for _, c := range h.Clients {
				c.Conn.CloseNow()
			}
			return
		case client := <-h.Register:
			h.Clients[client.Account.ID] = client
		case client := <-h.Unregister:
			delete(h.Clients, client.Account.ID)
			client.Conn.CloseNow()
		case msg := <-h.Broadcast:
			var failed []*Client
			for i := range msg.Recipients {
				client, ok := h.Clients[msg.Recipients[i]]
				if !ok {
					continue
				}

				writeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				err := client.Conn.Write(writeCtx, websocket.MessageText, msg.Payload)
				cancel()
				if err != nil {
					failed = append(failed, client)
				}
			}
			for _, c := range failed {
				delete(h.Clients, c.Account.ID)
				c.Conn.CloseNow()
			}
		}
	}
}
