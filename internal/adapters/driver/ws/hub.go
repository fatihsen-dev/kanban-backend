package ws

import (
	"context"
	"encoding/json"

	"github.com/fatihsen-dev/kanban-backend/internal/adapters/driver/http/datatransfers/responses"
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
			go client.SendUserStatusEvent("online")
		case client := <-h.unregister:
			if client, ok := h.clients[client.userID]; ok {
				go client.SendUserStatusEvent("offline")
				delete(h.clients, client.userID)
				close(client.send)
			}
		case message := <-h.broadcast:
			if message.ProjectID == "" {
				continue
			}

			for userID, client := range h.clients {
				if client.projectID == nil || *client.projectID != message.ProjectID {
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
		select {
		case client.send <- jsonData:
		default:
			close(client.send)
			delete(h.clients, userID)
		}
	}
}

func (h *Hub) GetOnlineUsers(projectID string) []responses.OnlineProjectMembersResponse {
	onlineUsers := []responses.OnlineProjectMembersResponse{}
	for _, client := range h.clients {
		if client.projectID != nil && *client.projectID == projectID {
			onlineUsers = append(onlineUsers, responses.OnlineProjectMembersResponse{
				ID:     client.userID,
				Status: "online",
			})
		}
	}

	return onlineUsers
}
