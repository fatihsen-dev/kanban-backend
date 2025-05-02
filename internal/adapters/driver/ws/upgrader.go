package ws

import (
	"net/http"

	"github.com/fatihsen-dev/kanban-backend/config"
	middlewares "github.com/fatihsen-dev/kanban-backend/internal/adapters/driver/http/middleware"
	"github.com/fatihsen-dev/kanban-backend/pkg/jwt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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
	projectID := c.Query("project_id")
	token := c.Query("token")

	if token == "" {
		return
	}

	user, err := jwt.VerifyToken(token)
	if err != nil {
		return
	}

	if projectID != "" {
		_, err = middlewares.CheckAccess(user.ID, projectID, c.Request.Context(), hub.projectMemberService)
		if err != nil {
			return
		}
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := NewClient(hub, conn, &projectID, user.ID)

	hub.register <- client

	go client.writePump()
	client.readPump()
}
