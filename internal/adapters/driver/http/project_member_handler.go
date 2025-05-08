package http

import (
	"net/http"

	"github.com/fatihsen-dev/kanban-backend/internal/adapters/driver/http/datatransfers"
	"github.com/fatihsen-dev/kanban-backend/internal/adapters/driver/http/datatransfers/requests"
	"github.com/fatihsen-dev/kanban-backend/internal/adapters/driver/http/datatransfers/responses"
	middlewares "github.com/fatihsen-dev/kanban-backend/internal/adapters/driver/http/middleware"
	"github.com/fatihsen-dev/kanban-backend/internal/adapters/driver/validation"
	"github.com/fatihsen-dev/kanban-backend/internal/adapters/driver/ws"
	"github.com/fatihsen-dev/kanban-backend/internal/core/domain"
	ports "github.com/fatihsen-dev/kanban-backend/internal/core/ports/driver"
	"github.com/gin-gonic/gin"
)

type projectMemberHandler struct {
	projectMemberService   ports.ProjectMemberService
	authMiddleware         *middlewares.AuthnMiddleware
	projectAuthzMiddleware *middlewares.ProjectAuthzMiddleware
	hub                    *ws.Hub
}

func NewProjectMemberHandler(projectMemberService ports.ProjectMemberService, authMiddleware *middlewares.AuthnMiddleware, projectAuthzMiddleware *middlewares.ProjectAuthzMiddleware, hub *ws.Hub) *projectMemberHandler {
	return &projectMemberHandler{projectMemberService: projectMemberService, authMiddleware: authMiddleware, projectAuthzMiddleware: projectAuthzMiddleware, hub: hub}
}

func (h *projectMemberHandler) RegisterProjectMemberRouter(r *gin.Engine) {

	projectMemberGroup := r.Group("/projects/:project_id/members")

	projectMemberGroup.Use(h.authMiddleware.Handle(false))
	projectMemberGroup.GET("/", h.projectAuthzMiddleware.Handle(middlewares.Member), h.GetProjectMembersHandler)
	projectMemberGroup.PUT("/:member_id", h.projectAuthzMiddleware.Handle(middlewares.Owner), h.UpdateProjectMemberHandler)
	projectMemberGroup.GET("/online", h.projectAuthzMiddleware.Handle(middlewares.Member), h.GetOnlineProjectMembersHandler)
}

func (h *projectMemberHandler) GetProjectMembersHandler(c *gin.Context) {
	projectID := c.Param("project_id")
	query := c.Query("query")

	members, users, err := h.projectMemberService.GetProjectMembersByProjectID(c.Request.Context(), projectID, &query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, datatransfers.ResponseError(err.Error()))
		return
	}

	userMap := make(map[string]*domain.User, len(users))
	for _, u := range users {
		userMap[u.ID] = u
	}

	response := make([]responses.ProjectMemberWithUserResponse, len(members))
	for i, member := range members {
		user := userMap[member.UserID]
		response[i] = responses.ProjectMemberWithUserResponse{
			ID:        member.ID,
			ProjectID: member.ProjectID,
			UserID:    member.UserID,
			Role:      string(member.Role),
			TeamID:    member.TeamID,
			CreatedAt: member.CreatedAt.String(),
			User: responses.UserResponse{
				ID:    user.ID,
				Email: user.Email,
			},
		}
	}

	c.JSON(http.StatusOK, datatransfers.ResponseSuccess("Project members fetched successfully", response))
}

func (h *projectMemberHandler) UpdateProjectMemberHandler(c *gin.Context) {
	projectID := c.Param("project_id")
	memberID := c.Param("member_id")

	var requestData requests.UpdateProjectMemberRequest
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.ResponseError(err.Error()))
		return
	}

	requestData.ID = memberID

	if err := validation.Validate(requestData); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.ResponseError(err.Error()))
		return
	}

	member := &domain.ProjectMember{
		ID: memberID,
	}

	response := responses.UpdateProjectMemberResponse{
		ID: memberID,
	}

	if requestData.Role != nil {
		member.Role = domain.AccessRole(*requestData.Role)
		response.Role = *requestData.Role
	}

	if requestData.TeamID != nil {
		member.TeamID = requestData.TeamID
		response.TeamID = *requestData.TeamID
	}

	err := h.projectMemberService.UpdateProjectMember(c.Request.Context(), member)
	if err != nil {
		c.JSON(http.StatusInternalServerError, datatransfers.ResponseError(err.Error()))
		return
	}

	h.hub.SendMessageToProject(projectID, ws.BaseResponse{
		Name: ws.EventNameProjectMemberUpdated,
		Data: response,
	})

	c.JSON(http.StatusOK, datatransfers.ResponseSuccess("Project member updated successfully", response))
}

func (h *projectMemberHandler) GetOnlineProjectMembersHandler(c *gin.Context) {
	projectID := c.Param("project_id")
	onlineUsers := h.hub.GetOnlineUsers(projectID)

	c.JSON(http.StatusOK, datatransfers.ResponseSuccess("Online project members fetched successfully", onlineUsers))
}
