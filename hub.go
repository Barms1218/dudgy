package main

import (
	"context"
	"slices"
	"time"

	"github.com/coder/websocket"
)

type Hub struct {
	clients    []*Client
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:    make([]*Client, 0),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			for _, c := range h.clients {
				c.Conn.CloseNow()
			}
			return
		case client := <-h.register:
			h.clients = append(h.clients, client)
		case client := <-h.unregister:
			for i, c := range h.clients {
				if c == client {
					h.clients = slices.Delete(h.clients, i, i+1)
					c.Conn.CloseNow()
					break
				}
			}
		case msg := <-h.broadcast:
			var failed []*Client
			for _, c := range h.clients {
				writeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				err := c.Conn.Write(writeCtx, websocket.MessageText, msg)
				cancel()
				if err != nil {
					failed = append(failed, c)
				}
			}
			for _, c := range failed {
				h.clients = slices.DeleteFunc(h.clients, func(x *Client) bool {
					return x == c
				})
				c.Conn.CloseNow()
			}
		}
	}
}
