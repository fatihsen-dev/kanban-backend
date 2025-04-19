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
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	appConfig := config.Read()
	defer zap.L().Sync()

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{appConfig.ClientUrl},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
	}))
	router.SetTrustedProxies(nil)

	postgresDB := db.NewPostgresRepository(appConfig.DBUrl)

	authnMiddleware := middlewares.NewAuthnMiddleware()

	hub := ws.NewHub()
	go hub.Run()

	router.GET("/ws/:project_id", func(c *gin.Context) {
		ws.ServeWs(hub, c)
	})

	// /users/* routes
	userRepo := db.NewPostgresUserRepo(postgresDB)
	userService := service.NewUserService(userRepo)
	userHandler := httphandler.NewUserHandler(userService, authnMiddleware)
	userHandler.RegisterUserRouter(router)

	// /auth/* routes
	authHandler := httphandler.NewAuthHandler(userService, authnMiddleware)
	authHandler.RegisterAuthRouter(router)

	projectRepo := db.NewPostgresProjectRepo(postgresDB)
	columnRepo := db.NewPostgresColumnRepo(postgresDB)
	taskRepo := db.NewPostgresTaskRepo(postgresDB)
	teamRepo := db.NewPostgresTeamRepo(postgresDB)
	projectMemberRepo := db.NewPostgresProjectMemberRepo(postgresDB)

	// /projects/* routes
	projectService := service.NewProjectService(projectRepo, columnRepo, taskRepo, teamRepo, projectMemberRepo)
	projectHandler := httphandler.NewProjectHandler(projectService, authnMiddleware, hub)
	projectHandler.RegisterProjectRouter(router)

	// /columns/* routes
	columnService := service.NewColumnService(columnRepo, taskRepo)
	columnHandler := httphandler.NewColumnHandler(columnService, authnMiddleware, hub)
	columnHandler.RegisterColumnRouter(router)

	// /tasks/* routes
	taskService := service.NewTaskService(taskRepo)
	taskHandler := httphandler.NewTaskHandler(taskService, authnMiddleware, hub)
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
