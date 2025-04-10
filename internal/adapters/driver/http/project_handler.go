package http

import (
	"fmt"
	"net/http"

	middlewares "github.com/fatihsen-dev/kanban-backend/internal/adapters/driver/http/middleware"
	"github.com/fatihsen-dev/kanban-backend/internal/adapters/driver/ws"
	"github.com/fatihsen-dev/kanban-backend/internal/core/domain"
	ports "github.com/fatihsen-dev/kanban-backend/internal/core/ports/driver"
	"github.com/gin-gonic/gin"
)

type projectHandler struct {
	projectService ports.ProjectService
	authMiddleware *middlewares.AuthMiddleware
	hub            *ws.Hub
}

func NewProjectHandler(projectService ports.ProjectService, authMiddleware *middlewares.AuthMiddleware, hub *ws.Hub) *projectHandler {
	return &projectHandler{projectService: projectService, authMiddleware: authMiddleware, hub: hub}
}

func (h *projectHandler) RegisterProjectRouter(r *gin.Engine) {
	r.POST("/projects", h.authMiddleware.Handle, h.CreateProjectHandler)
	r.GET("/projects", h.authMiddleware.Handle, h.GetProjectsHandler)
	r.GET("/projects/:id", h.authMiddleware.Handle, h.GetProjectHandler)
}

func (h *projectHandler) CreateProjectHandler(c *gin.Context) {

	var requestData struct {
		Name string `json:"name"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	project := &domain.Project{
		Name: requestData.Name,
	}

	err := h.projectService.CreateProject(c.Request.Context(), project)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Project created successfully"})
}

func (h *projectHandler) GetProjectHandler(c *gin.Context) {
	id := c.Param("id")

	project, err := h.projectService.GetProjectByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get project"})
		return
	}

	c.JSON(http.StatusOK, project)
}

func (h *projectHandler) GetProjectsHandler(c *gin.Context) {
	projects, err := h.projectService.GetProjects(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get projects"})
		return
	}
	c.JSON(http.StatusOK, projects)
}
