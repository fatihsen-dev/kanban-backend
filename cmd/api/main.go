package main

import (
	"fmt"

	"github.com/fatihsen-dev/kanban-backend/config"
	db "github.com/fatihsen-dev/kanban-backend/internal/adapters/driven/db/postgres"
	"github.com/fatihsen-dev/kanban-backend/internal/adapters/driver/http"
	"github.com/fatihsen-dev/kanban-backend/internal/adapters/driver/ws"
	"github.com/fatihsen-dev/kanban-backend/internal/core/service"
	"github.com/gin-gonic/gin"
)

func main() {
	appConfig := config.Read()

	router := gin.Default()
	router.SetTrustedProxies(nil)

	postgresDB := db.NewPostgresRepository(appConfig.DBUrl)

	hub := ws.NewHub()
	go hub.Run()

	router.GET("/ws/:project_id", func(c *gin.Context) {
		ws.ServeWs(hub, c)
	})

	// /projects/* routes
	projectRepo := db.NewPostgresProjectRepo(postgresDB)
	projectService := service.NewProjectService(projectRepo)
	projectHandler := http.NewProjectHandler(projectService, hub)
	projectHandler.RegisterProjectRouter(router)

	// /columns/* routes
	columnRepo := db.NewPostgresColumnRepo(postgresDB)
	columnService := service.NewColumnService(columnRepo)
	columnHandler := http.NewColumnHandler(columnService, hub)
	columnHandler.RegisterColumnRouter(router)

	// /tasks/* routes
	taskRepo := db.NewPostgresTaskRepo(postgresDB)
	taskService := service.NewTaskService(taskRepo)
	taskHandler := http.NewTaskHandler(taskService, hub)
	taskHandler.RegisterTaskRouter(router)

	router.Run(fmt.Sprintf(":%s", appConfig.Port))
}
