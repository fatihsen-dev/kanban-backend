package http

import (
	"net/http"
	"time"

	"github.com/fatihsen-dev/kanban-backend/internal/adapters/driver/http/datatransfers"
	"github.com/fatihsen-dev/kanban-backend/internal/adapters/driver/http/datatransfers/responses"
	middlewares "github.com/fatihsen-dev/kanban-backend/internal/adapters/driver/http/middleware"
	"github.com/fatihsen-dev/kanban-backend/internal/adapters/driver/validation"
	"github.com/fatihsen-dev/kanban-backend/internal/adapters/driver/ws"
	ports "github.com/fatihsen-dev/kanban-backend/internal/core/ports/driver"
	"github.com/gin-gonic/gin"
)

type projectMemberHandler struct {
	projectMemberService ports.ProjectMemberService
	authMiddleware       *middlewares.AuthnMiddleware
	hub                  *ws.Hub
}

func NewProjectMemberHandler(projectMemberService ports.ProjectMemberService, authMiddleware *middlewares.AuthnMiddleware, hub *ws.Hub) *projectMemberHandler {
	return &projectMemberHandler{projectMemberService: projectMemberService, authMiddleware: authMiddleware, hub: hub}
}

func (h *projectMemberHandler) RegisterProjectMemberRouter(r *gin.Engine) {
	r.Use(h.authMiddleware.Handle(false)).GET("/projects/:project_id/members", h.GetProjectMembersHandler)
}

func (h *projectMemberHandler) GetProjectMembersHandler(c *gin.Context) {
	projectID := c.Param("project_id")

	err := validation.ValidateUUID(projectID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, datatransfers.ResponseError("Invalid project ID"))
		return
	}

	projectMembers, users, err := h.projectMemberService.GetProjectMembersByProjectID(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, datatransfers.ResponseError(err.Error()))
		return
	}

	userMap := make(map[string]responses.UserResponse)
	for _, user := range users {
		userMap[user.ID] = responses.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			IsAdmin:   user.IsAdmin,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
		}
	}

	projectMembersResponse := make([]responses.ProjectMemberWithUserResponse, len(projectMembers))
	for i, projectMember := range projectMembers {
		projectMembersResponse[i] = responses.ProjectMemberWithUserResponse{
			ID:        projectMember.ID,
			UserID:    projectMember.UserID,
			ProjectID: projectMember.ProjectID,
			Role:      string(projectMember.Role),
			TeamID:    projectMember.TeamID,
			CreatedAt: projectMember.CreatedAt.Format(time.RFC3339),
			User:      userMap[projectMember.UserID],
		}
	}

	c.JSON(http.StatusOK, datatransfers.ResponseSuccess("Project members fetched successfully", projectMembersResponse))
}
