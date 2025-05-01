package ws

import (
	"context"
	"encoding/json"

	ports "github.com/fatihsen-dev/kanban-backend/internal/core/ports/driver"
)

type BroadcastMessage struct {
	ProjectID string
	Data      []byte
}

type Hub struct {
	clients              map[string]*Client
	broadcast            chan BroadcastMessage
	register             chan *Client
	unregister           chan *Client
	ctx                  context.Context
	projectMemberService ports.ProjectMemberService
}

func NewHub(projectMemberService ports.ProjectMemberService) *Hub {
	return &Hub{
		broadcast:            make(chan BroadcastMessage),
		register:             make(chan *Client),
		unregister:           make(chan *Client),
		clients:              make(map[string]*Client),
		ctx:                  context.Background(),
		projectMemberService: projectMemberService,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.userID] = client
		case client := <-h.unregister:
			if client, ok := h.clients[client.userID]; ok {
				delete(h.clients, client.userID)
				close(client.send)
			}
		case message := <-h.broadcast:
			if message.ProjectID == "" {
				continue
			}

			for userID, client := range h.clients {
				if *client.projectID != message.ProjectID {
					continue
				}

				select {
				case client.send <- message.Data:
				default:
					close(client.send)
					delete(h.clients, userID)
				}
			}
		}
	}
}

func (h *Hub) SendMessageToProject(projectID string, data BaseResponse) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}

	h.broadcast <- BroadcastMessage{
		ProjectID: projectID,
		Data:      jsonData,
	}
}

func (h *Hub) SendMessageToUser(userID string, data BaseResponse) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}

	if client, ok := h.clients[userID]; ok {
		client.send <- jsonData
	}
}
