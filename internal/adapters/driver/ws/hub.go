package ws

type BroadcastMessage struct {
	GroupID string
	Data    []byte
}

type Hub struct {
	clients    map[*Client]bool
	groups     map[string][]*Client
	broadcast  chan BroadcastMessage
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan BroadcastMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		groups:     make(map[string][]*Client),
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

				for groupID, clients := range h.groups {
					for i, c := range clients {
						if c == client {
							h.groups[groupID] = append(clients[:i], clients[i+1:]...)
							if len(h.groups[groupID]) == 0 {
								delete(h.groups, groupID)
							}
							break
						}
					}
				}
			}
		case message := <-h.broadcast:
			if message.GroupID != "" {
				if clients, ok := h.groups[message.GroupID]; ok {
					for _, client := range clients {
						select {
						case client.send <- message.Data:
						default:
							close(client.send)
							delete(h.clients, client)
							h.RemoveFromGroup(message.GroupID, client)
						}
					}
				}
			} else {
				for client := range h.clients {
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

func (h *Hub) AddToGroup(groupID string, client *Client) {
	h.groups[groupID] = append(h.groups[groupID], client)
}

func (h *Hub) RemoveFromGroup(groupID string, client *Client) {
	if clients, ok := h.groups[groupID]; ok {
		for i, c := range clients {
			if c == client {
				h.groups[groupID] = append(clients[:i], clients[i+1:]...)
				if len(h.groups[groupID]) == 0 {
					delete(h.groups, groupID)
				}
				break
			}
		}
	}
}

func (h *Hub) Broadcast(groupID string, data []byte) {
	h.broadcast <- BroadcastMessage{
		GroupID: groupID,
		Data:    data,
	}
}
