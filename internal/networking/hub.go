package networking

import (
	"context"
	"sync"
	"time"

	"github.com/coder/websocket"
)

type Client struct {
	ID     string `json:"id"`
	Conn   *websocket.Conn
	Ctx    context.Context
	Cancel context.CancelFunc
	Mu     sync.Mutex
}

type Hub struct {
	Clients    map[string]*Client
	Broadcast  chan BroadCastMessage
	Register   chan *Client
	Unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[string]*Client, 0),
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
			h.Clients[client.ID] = client
		case client := <-h.Unregister:
			delete(h.Clients, client.ID)
			client.Conn.CloseNow()
		case msg := <-h.Broadcast:
			var failed []string
			for i := range msg.Recipients {
				client, ok := h.Clients[msg.Recipients[i]]
				if !ok {
					continue
				}

				writeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				err := client.Conn.Write(writeCtx, websocket.MessageText, msg.Payload)
				cancel()
				if err != nil {
					failed = append(failed, msg.Recipients[i])
				}
			}
			for _, id := range failed {
				delete(h.Clients, id)
				h.Clients[id].Conn.CloseNow()
			}
		}
	}
}
