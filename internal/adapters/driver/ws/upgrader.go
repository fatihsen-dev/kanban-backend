package ws

import (
	"net/http"

	"github.com/fatihsen-dev/kanban-backend/config"
	middlewares "github.com/fatihsen-dev/kanban-backend/internal/adapters/driver/http/middleware"
	"github.com/fatihsen-dev/kanban-backend/pkg/jwt"
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
	token := c.Query("token")

	if token == "" {
		zap.L().Error("No token provided")
		return
	}

	user, err := jwt.VerifyToken(token)
	if err != nil {
		zap.L().Error("Failed to validate token", zap.Error(err))
		return
	}

	_, err = middlewares.CheckMemberAccess(user.ID, groupID, c.Request.Context(), hub.projectMemberService)
	if err != nil {
		zap.L().Error("Failed to check member access", zap.Error(err))
		return
	}

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
