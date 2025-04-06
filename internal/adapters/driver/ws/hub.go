package ws

import (
	"context"
	"encoding/json"
	"fmt"
)

type BroadcastMessage struct {
	ProjectID string
	Data      []byte
}

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan BroadcastMessage
	register   chan *Client
	unregister chan *Client
	ctx        context.Context
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan BroadcastMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		ctx:        context.Background(),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				if client.projectID == message.ProjectID {
					select {
					case client.send <- message.Data:
					default:
						close(client.send)
						delete(h.clients, client)
					}
				}
			}
		}
	}
}

func (h *Hub) SendMessage(projectID string, data BaseResponse) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return
	}

	h.broadcast <- BroadcastMessage{
		ProjectID: projectID,
		Data:      jsonData,
	}
}
