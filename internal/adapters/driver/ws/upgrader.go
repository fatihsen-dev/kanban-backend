package ws

import (
	"net/http"

	"github.com/fatihsen-dev/kanban-backend/config"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var appConfig = config.Read()
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		allowedOrigins := []string{appConfig.ClientUrl}
		origin := r.Header.Get("Origin")
		for _, allowedOrigin := range allowedOrigins {
			if allowedOrigin == origin {
				return true
			}
		}
		return false
	},
}

func ServeWs(hub *Hub, c *gin.Context) {
	groupID := c.Param("project_id")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		zap.L().Error("Failed to upgrade connection", zap.Error(err))
		return
	}

	client := NewClient(hub, conn, groupID)

	hub.register <- client

	go client.writePump()
	client.readPump()
}
