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

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{appConfig.ClientUrl},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
	}))
	router.SetTrustedProxies(nil)

	// repositories
	postgresDB := db.NewPostgresRepository(appConfig.DBUrl)
	userRepo := db.NewPostgresUserRepo(postgresDB)
	projectRepo := db.NewPostgresProjectRepo(postgresDB)
	columnRepo := db.NewPostgresColumnRepo(postgresDB)
	taskRepo := db.NewPostgresTaskRepo(postgresDB)
	teamRepo := db.NewPostgresTeamRepo(postgresDB)
	projectMemberRepo := db.NewPostgresProjectMemberRepo(postgresDB)
	invitationRepo := db.NewPostgresInvitationRepo(postgresDB)

	// services
	userService := service.NewUserService(userRepo)
	projectService := service.NewProjectService(projectRepo, columnRepo, taskRepo, teamRepo, projectMemberRepo, userRepo)
	columnService := service.NewColumnService(columnRepo, taskRepo)
	taskService := service.NewTaskService(taskRepo)
	projectMemberService := service.NewProjectMemberService(projectMemberRepo, userRepo)
	teamService := service.NewTeamService(teamRepo, projectMemberRepo)
	invitationService := service.NewInvitationService(invitationRepo, userRepo, projectRepo, projectMemberRepo)

	hub := ws.NewHub(projectMemberService)
	go hub.Run()

	router.GET("/ws", func(c *gin.Context) {
		ws.ServeWs(hub, c)
	})

	// middlewares
	authnMiddleware := middlewares.NewAuthnMiddleware()
	projectAuthzMiddleware := middlewares.NewProjectAuthzMiddleware(projectMemberService, teamService)

	// /users/* routes
	userHandler := httphandler.NewUserHandler(userService, authnMiddleware)
	userHandler.RegisterUserRouter(router)

	// /auth/* routes
	authHandler := httphandler.NewAuthHandler(userService, authnMiddleware)
	authHandler.RegisterAuthRouter(router)

	// /invitations/* routes
	invitationHandler := httphandler.NewInvitationHandler(invitationService, authnMiddleware, projectAuthzMiddleware, hub)
	invitationHandler.RegisterInvitationRouter(router)

	// /projects/* routes
	projectHandler := httphandler.NewProjectHandler(projectService, authnMiddleware, projectAuthzMiddleware, hub)
	projectHandler.RegisterProjectRouter(router)

	// /projects/:project_id/members/* routes
	projectMemberHandler := httphandler.NewProjectMemberHandler(projectMemberService, authnMiddleware, projectAuthzMiddleware, hub)
	projectMemberHandler.RegisterProjectMemberRouter(router)

	// /projects/:project_id/teams/* routes
	teamHandler := httphandler.NewTeamHandler(teamService, authnMiddleware, projectAuthzMiddleware, hub)
	teamHandler.RegisterTeamRouter(router)

	// /projects/:project_id/columns/* routes
	columnHandler := httphandler.NewColumnHandler(columnService, authnMiddleware, projectAuthzMiddleware, hub)
	columnHandler.RegisterColumnRouter(router)

	// /projects/:project_id/tasks/* routes
	taskHandler := httphandler.NewTaskHandler(taskService, authnMiddleware, projectAuthzMiddleware, hub)
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
