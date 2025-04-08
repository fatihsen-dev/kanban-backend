package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fatihsen-dev/kanban-backend/config"
	db "github.com/fatihsen-dev/kanban-backend/internal/adapters/driven/db/postgres"
	httphandler "github.com/fatihsen-dev/kanban-backend/internal/adapters/driver/http"
	middlewares "github.com/fatihsen-dev/kanban-backend/internal/adapters/driver/http/middleware"
	"github.com/fatihsen-dev/kanban-backend/internal/adapters/driver/ws"
	"github.com/fatihsen-dev/kanban-backend/internal/core/service"
	_ "github.com/fatihsen-dev/kanban-backend/pkg/log"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	appConfig := config.Read()
	defer zap.L().Sync()

	router := gin.Default()
	router.SetTrustedProxies(nil)

	postgresDB := db.NewPostgresRepository(appConfig.DBUrl)

	authMiddleware := middlewares.NewAuthMiddleware(false)

	hub := ws.NewHub()
	go hub.Run()

	router.GET("/ws/:project_id", func(c *gin.Context) {
		ws.ServeWs(hub, c)
	})

	// /users/* routes
	userRepo := db.NewPostgresUserRepo(postgresDB)
	userService := service.NewUserService(userRepo)
	userHandler := httphandler.NewUserHandler(userService, authMiddleware)
	userHandler.RegisterUserRouter(router)

	// /auth/* routes
	authHandler := httphandler.NewAuthHandler(userService, authMiddleware)
	authHandler.RegisterAuthRouter(router)

	// /projects/* routes
	projectRepo := db.NewPostgresProjectRepo(postgresDB)
	projectService := service.NewProjectService(projectRepo)
	projectHandler := httphandler.NewProjectHandler(projectService, authMiddleware, hub)
	projectHandler.RegisterProjectRouter(router)

	// /columns/* routes
	columnRepo := db.NewPostgresColumnRepo(postgresDB)
	columnService := service.NewColumnService(columnRepo)
	columnHandler := httphandler.NewColumnHandler(columnService, authMiddleware, hub)
	columnHandler.RegisterColumnRouter(router)

	// /tasks/* routes
	taskRepo := db.NewPostgresTaskRepo(postgresDB)
	taskService := service.NewTaskService(taskRepo)
	taskHandler := httphandler.NewTaskHandler(taskService, authMiddleware, hub)
	taskHandler.RegisterTaskRouter(router)

	router.Run(fmt.Sprintf(":%s", appConfig.Port))

	gracefulShutdown(router)
}

func gracefulShutdown(router *gin.Engine) {
	srv := &http.Server{
		Handler: router,
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	zap.L().Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Error("Error during server shutdown", zap.Error(err))
	}

	zap.L().Info("Server gracefully stopped")
}
