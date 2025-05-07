package http

import (
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
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type teamHandler struct {
	teamService            ports.TeamService
	authMiddleware         *middlewares.AuthnMiddleware
	projectAuthzMiddleware *middlewares.ProjectAuthzMiddleware
	hub                    *ws.Hub
}

func NewTeamHandler(teamService ports.TeamService, authMiddleware *middlewares.AuthnMiddleware, projectAuthzMiddleware *middlewares.ProjectAuthzMiddleware, hub *ws.Hub) *teamHandler {
	return &teamHandler{teamService: teamService, authMiddleware: authMiddleware, projectAuthzMiddleware: projectAuthzMiddleware, hub: hub}
}

func (h *teamHandler) RegisterTeamRouter(r *gin.Engine) {

	teamGroup := r.Group("/projects/:project_id/teams")

	teamGroup.Use(h.authMiddleware.Handle(false))

	teamGroup.POST("", h.projectAuthzMiddleware.Handle(middlewares.Owner), h.CreateTeamHandler)
	teamGroup.GET("", h.projectAuthzMiddleware.Handle(middlewares.Member), h.GetTeamsHandler)
	teamGroup.GET("/:team_id", h.projectAuthzMiddleware.Handle(middlewares.Member), h.GetTeamHandler)
	teamGroup.PUT("/:team_id", h.projectAuthzMiddleware.Handle(middlewares.Owner), h.UpdateTeamHandler)
	teamGroup.DELETE("/:team_id", h.projectAuthzMiddleware.Handle(middlewares.Owner), h.DeleteTeamHandler)
}

func (h *teamHandler) CreateTeamHandler(c *gin.Context) {
	var requestData requests.CreateTeamRequest
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	requestData.ProjectID = c.Param("project_id")
	if err := validation.Validate(requestData); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.ResponseError(err.Error()))
		return
	}

	projectMember, ok := c.MustGet("project_member").(domain.ProjectMember)
	if !ok {
		c.JSON(http.StatusInternalServerError, datatransfers.ResponseError("Project member not found"))
		return
	}

	if projectMember.Role != domain.AccessOwnerRole && requestData.Role == string(domain.AccessAdminRole) {
		c.JSON(http.StatusForbidden, datatransfers.ResponseError("You are not authorized to create a team"))
		return
	}

	team := &domain.Team{
		Name:      requestData.Name,
		Role:      domain.AccessRole(requestData.Role),
		ProjectID: requestData.ProjectID,
	}

	err := h.teamService.CreateTeam(c.Request.Context(), team)
	if err != nil {
		zap.L().Error("Failed to create team", zap.Error(err))
		c.JSON(http.StatusInternalServerError, datatransfers.ResponseError(err.Error()))
		return
	}

	responseData := responses.TeamWithMembersResponse{
		ID:        team.ID,
		Name:      team.Name,
		Role:      string(team.Role),
		ProjectID: team.ProjectID,
		Members:   make([]responses.ProjectMemberResponse, 0),
		CreatedAt: team.CreatedAt.Format(time.RFC3339),
	}

	h.hub.SendMessageToProject(team.ProjectID, ws.BaseResponse{
		Name: ws.EventNameTeamCreated,
		Data: responseData,
	})

	c.JSON(http.StatusCreated, datatransfers.ResponseSuccess("Team created successfully", responseData))
}

func (h *teamHandler) GetTeamsHandler(c *gin.Context) {
	projectID := c.Param("project_id")

	err := validation.ValidateUUID(projectID)
	if err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.ResponseError("Invalid project ID"))
		return
	}

	teams, err := h.teamService.GetTeamsByProjectID(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, datatransfers.ResponseError("Failed to get teams"))
		return
	}

	responseData := make([]responses.TeamResponse, len(teams))
	for i, team := range teams {
		responseData[i] = responses.TeamResponse{
			ID:        team.ID,
			Name:      team.Name,
			ProjectID: team.ProjectID,
			CreatedAt: team.CreatedAt.Format(time.RFC3339),
		}
	}

	c.JSON(http.StatusOK, datatransfers.ResponseSuccess("Teams fetched successfully", responseData))
}

func (h *teamHandler) GetTeamHandler(c *gin.Context) {
	teamID := c.Param("team_id")

	err := validation.ValidateUUID(teamID)
	if err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.ResponseError("Invalid team ID"))
		return
	}

	team, err := h.teamService.GetTeamByID(c.Request.Context(), teamID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, datatransfers.ResponseError("Failed to get team"))
		return
	}

	responseData := responses.TeamResponse{
		ID:        team.ID,
		Name:      team.Name,
		ProjectID: team.ProjectID,
		CreatedAt: team.CreatedAt.Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, datatransfers.ResponseSuccess("Team fetched successfully", responseData))
}

func (h *teamHandler) UpdateTeamHandler(c *gin.Context) {
	teamID := c.Param("team_id")
	projectID := c.Param("project_id")

	err := validation.ValidateUUID(teamID)
	if err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.ResponseError("Invalid team ID"))
		return
	}

	var requestData requests.UpdateTeamRequest
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validation.Validate(requestData); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.ResponseError(err.Error()))
		return
	}

	team := &domain.Team{
		ID: teamID,
	}

	if requestData.Name != nil {
		team.Name = *requestData.Name
	}

	if requestData.Role != nil {
		team.Role = domain.AccessRole(*requestData.Role)
	}

	err = h.teamService.UpdateTeam(c.Request.Context(), team)
	if err != nil {
		c.JSON(http.StatusInternalServerError, datatransfers.ResponseError("Failed to update team"))
		return
	}

	responseData := responses.UpdateTeamResponse{
		ID:        team.ID,
		Name:      team.Name,
		Role:      string(team.Role),
		ProjectID: projectID,
	}

	h.hub.SendMessageToProject(projectID, ws.BaseResponse{
		Name: ws.EventNameTeamUpdated,
		Data: responseData,
	})

	c.JSON(http.StatusOK, datatransfers.ResponseSuccess("Team updated successfully", responseData))
}

func (h *teamHandler) DeleteTeamHandler(c *gin.Context) {
	teamID := c.Param("team_id")
	projectID := c.Param("project_id")

	err := validation.ValidateUUID(teamID)
	if err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.ResponseError("Invalid team ID"))
		return
	}

	err = h.teamService.DeleteTeamByID(c.Request.Context(), teamID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, datatransfers.ResponseError("Failed to delete team"))
		return
	}

	responseData := responses.DeleteTeamResponse{
		ID: teamID,
	}

	h.hub.SendMessageToProject(projectID, ws.BaseResponse{
		Name: ws.EventNameTeamDeleted,
		Data: responseData,
	})

	c.JSON(http.StatusOK, datatransfers.ResponseSuccess("Team deleted successfully", responseData))
}
