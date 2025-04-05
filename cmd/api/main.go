package main

import (
	"fmt"

	"github.com/fatihsen-dev/kanban-backend/config"
	"github.com/fatihsen-dev/kanban-backend/internal/adapters/driver/ws"
	"github.com/gin-gonic/gin"
)

func main() {
	appConfig := config.Read()

	router := gin.Default()
	router.SetTrustedProxies(nil)

	hub := ws.NewHub()
	go hub.Run()

	router.GET("/ws", func(c *gin.Context) {
		ws.ServeWs(hub, c)
	})

	router.Run(fmt.Sprintf(":%s", appConfig.Port))
}
