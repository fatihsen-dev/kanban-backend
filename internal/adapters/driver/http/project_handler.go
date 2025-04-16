package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fatihsen-dev/kanban-backend/internal/adapters/driver/http/datatransfers"
	"github.com/fatihsen-dev/kanban-backend/internal/adapters/driver/http/datatransfers/requests"
	"github.com/fatihsen-dev/kanban-backend/internal/adapters/driver/http/datatransfers/responses"
	middlewares "github.com/fatihsen-dev/kanban-backend/internal/adapters/driver/http/middleware"
	"github.com/fatihsen-dev/kanban-backend/internal/adapters/driver/validation"
	"github.com/fatihsen-dev/kanban-backend/internal/adapters/driver/ws"
	"github.com/fatihsen-dev/kanban-backend/internal/core/domain"
	ports "github.com/fatihsen-dev/kanban-backend/internal/core/ports/driver"
	"github.com/fatihsen-dev/kanban-backend/pkg/jwt"
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
	r.GET("/projects/:id/columns", h.authMiddleware.Handle, h.GetProjectWithColumnsHandler)
}

func (h *projectHandler) CreateProjectHandler(c *gin.Context) {
	userClaims := c.MustGet("user").(*jwt.UserClaims)

	var requestData requests.ProjectCreateRequest

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.ResponseError("Invalid request data"))
		return
	}

	project := &domain.Project{
		Name:    requestData.Name,
		OwnerID: userClaims.UserID,
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
		OwnerID:   project.OwnerID,
		CreatedAt: project.CreatedAt.Format(time.RFC3339),
	}

	c.JSON(http.StatusCreated, datatransfers.ResponseSuccess("Project created successfully", responseData))
}

func (h *projectHandler) GetProjectHandler(c *gin.Context) {
	id := c.Param("id")

	err := validation.ValidateUUID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.ResponseError("Invalid project ID"))
		return
	}

	project, err := h.projectService.GetProjectByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, datatransfers.ResponseError("Project not found"))
		return
	}

	responseData := responses.ProjectResponse{
		ID:        project.ID,
		Name:      project.Name,
		CreatedAt: project.CreatedAt.Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, datatransfers.ResponseSuccess("Project fetched successfully", responseData))
}

func (h *projectHandler) GetProjectWithColumnsHandler(c *gin.Context) {
	projectID := c.Param("id")

	err := validation.ValidateUUID(projectID)
	if err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.ResponseError("Invalid project ID"))
		return
	}

	project, columns, tasksByColumn, err := h.projectService.GetProjectWithColumns(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(http.StatusNotFound, datatransfers.ResponseError("Project not found"))
		return
	}

	columnResponses := make([]responses.ColumnWithTasksResponse, len(columns))
	for i, column := range columns {
		tasks := tasksByColumn[column.ID]
		taskResponses := make([]responses.TaskResponse, len(tasks))
		for j, task := range tasks {
			taskResponses[j] = responses.TaskResponse{
				ID:        task.ID,
				Title:     task.Title,
				ProjectID: task.ProjectID,
				ColumnID:  task.ColumnID,
				CreatedAt: task.CreatedAt.Format(time.RFC3339),
			}
		}

		columnResponses[i] = responses.ColumnWithTasksResponse{
			ID:        column.ID,
			Name:      column.Name,
			CreatedAt: column.CreatedAt.Format(time.RFC3339),
			Tasks:     taskResponses,
		}
	}

	response := responses.ProjectWithDetailsResponse{
		ID:        project.ID,
		Name:      project.Name,
		OwnerID:   project.OwnerID,
		CreatedAt: project.CreatedAt.Format(time.RFC3339),
		Columns:   columnResponses,
	}

	c.JSON(http.StatusOK, datatransfers.ResponseSuccess("Project details fetched successfully", response))
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
