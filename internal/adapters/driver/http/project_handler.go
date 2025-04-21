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
	"go.uber.org/zap"
)

type projectHandler struct {
	projectService         ports.ProjectService
	authMiddleware         *middlewares.AuthnMiddleware
	projectAuthzMiddleware *middlewares.ProjectAuthzMiddleware
	hub                    *ws.Hub
}

func NewProjectHandler(projectService ports.ProjectService, authMiddleware *middlewares.AuthnMiddleware, projectAuthzMiddleware *middlewares.ProjectAuthzMiddleware, hub *ws.Hub) *projectHandler {
	return &projectHandler{projectService: projectService, authMiddleware: authMiddleware, projectAuthzMiddleware: projectAuthzMiddleware, hub: hub}
}

func (h *projectHandler) RegisterProjectRouter(r *gin.Engine) {
	r.Use(h.authMiddleware.Handle(false)).POST("/projects", h.CreateProjectHandler)
	r.Use(h.authMiddleware.Handle(false)).GET("/projects", h.GetProjectsHandler)
	r.Use(h.projectAuthzMiddleware.Handle(middlewares.Admin)).GET("/projects/:project_id", h.GetProjectHandler)
}

func (h *projectHandler) CreateProjectHandler(c *gin.Context) {
	userClaims := c.MustGet("user").(*jwt.UserClaims)

	var requestData requests.ProjectCreateRequest

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.ResponseError("Invalid request data"))
		return
	}

	if err := validation.Validate(requestData); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.ResponseError(err.Error()))
		return
	}

	project := &domain.Project{
		Name:    requestData.Name,
		OwnerID: userClaims.ID,
	}

	err := h.projectService.CreateProject(c.Request.Context(), project)

	if err != nil {
		zap.L().Error("Failed to create project", zap.Error(err))
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
	projectID := c.Param("project_id")

	err := validation.ValidateUUID(projectID)
	if err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.ResponseError("Invalid project ID"))
		return
	}

	project, columns, tasksByColumn, teams, projectMembers, err := h.projectService.GetProjectWithDetails(c.Request.Context(), projectID)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusNotFound, datatransfers.ResponseError("Project not found"))
		return
	}

	teamResponses := make([]responses.TeamResponse, len(teams))
	for i, team := range teams {
		teamResponses[i] = responses.TeamResponse{
			ID:        team.ID,
			Name:      team.Name,
			Role:      string(team.Role),
			ProjectID: team.ProjectID,
			CreatedAt: team.CreatedAt.Format(time.RFC3339),
		}
	}

	memberResponses := make([]responses.ProjectMemberResponse, len(projectMembers))
	for i, member := range projectMembers {
		memberResponses[i] = responses.ProjectMemberResponse{
			ID:        member.ID,
			UserID:    member.UserID,
			Role:      string(member.Role),
			ProjectID: member.ProjectID,
			CreatedAt: member.CreatedAt.Format(time.RFC3339),
		}

		if member.TeamID != nil {
			memberResponses[i].TeamID = *member.TeamID
		}
	}

	columnResponses := make([]responses.ColumnWithDetailsResponse, len(columns))
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

		columnResponses[i] = responses.ColumnWithDetailsResponse{
			ID:        column.ID,
			Name:      column.Name,
			Color:     column.Color,
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
		Teams:     teamResponses,
		Members:   memberResponses,
	}

	c.JSON(http.StatusOK, datatransfers.ResponseSuccess("Project details fetched successfully", response))
}

func (h *projectHandler) GetProjectsHandler(c *gin.Context) {
	userClaims := c.MustGet("user").(*jwt.UserClaims)
	projects, err := h.projectService.GetUserProjects(c.Request.Context(), userClaims.ID)
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
