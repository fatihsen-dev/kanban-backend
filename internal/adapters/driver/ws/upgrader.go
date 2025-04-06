package ws

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func ServeWs(hub *Hub, c *gin.Context) {
	groupID := c.Param("project_id")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		conn.Close()
		return
	}

	client := NewClient(hub, conn, groupID)

	hub.register <- client

	go client.writePump()
	client.readPump()
}
