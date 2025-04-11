package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fatihsen-dev/kanban-backend/internal/adapters/driver/http/datatransfers"
	"github.com/fatihsen-dev/kanban-backend/internal/adapters/driver/http/datatransfers/requests"
	"github.com/fatihsen-dev/kanban-backend/internal/adapters/driver/http/datatransfers/responses"
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

	var requestData requests.ProjectCreateRequest

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.ResponseError("Invalid request data"))
		return
	}

	project := &domain.Project{
		Name: requestData.Name,
	}

	err := h.projectService.CreateProject(c.Request.Context(), project)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, datatransfers.ResponseError("Failed to create project"))
		return
	}

	responseData := responses.ProjectResponse{
		ID:        project.ID,
		Name:      project.Name,
		CreatedAt: project.CreatedAt.Format(time.RFC3339),
	}

	c.JSON(http.StatusCreated, datatransfers.ResponseSuccess("Project created successfully", responseData))
}

func (h *projectHandler) GetProjectHandler(c *gin.Context) {
	id := c.Param("id")

	project, err := h.projectService.GetProjectByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, datatransfers.ResponseError("Failed to get project"))
		return
	}

	responseData := responses.ProjectResponse{
		ID:        project.ID,
		Name:      project.Name,
		CreatedAt: project.CreatedAt.Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, datatransfers.ResponseSuccess("Project fetched successfully", responseData))
}

func (h *projectHandler) GetProjectsHandler(c *gin.Context) {
	projects, err := h.projectService.GetProjects(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, datatransfers.ResponseError("Failed to get projects"))
		return
	}

	projectResponses := make([]responses.ProjectResponse, len(projects))
	for i, project := range projects {
		projectResponses[i] = responses.ProjectResponse{
			ID:        project.ID,
			Name:      project.Name,
			CreatedAt: project.CreatedAt.Format(time.RFC3339),
		}
	}

	c.JSON(http.StatusOK, datatransfers.ResponseSuccess("Projects fetched successfully", projectResponses))
}
