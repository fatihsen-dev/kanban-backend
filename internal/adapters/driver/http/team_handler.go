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
)

type teamHandler struct {
	teamService    ports.TeamService
	authMiddleware *middlewares.AuthnMiddleware
	hub            *ws.Hub
}

func NewTeamHandler(teamService ports.TeamService, authMiddleware *middlewares.AuthnMiddleware, hub *ws.Hub) *teamHandler {
	return &teamHandler{teamService: teamService, authMiddleware: authMiddleware, hub: hub}
}

func (h *teamHandler) RegisterTeamRouter(r *gin.Engine) {
	r.POST("/teams", h.authMiddleware.Handle, h.CreateTeamHandler)
	r.GET("/teams", h.authMiddleware.Handle, h.GetTeamsHandler)
	r.GET("/teams/:id", h.authMiddleware.Handle, h.GetTeamHandler)
	r.PUT("/teams/:id", h.authMiddleware.Handle, h.UpdateTeamHandler)
	r.DELETE("/teams/:id", h.authMiddleware.Handle, h.DeleteTeamHandler)
	r.POST("/teams/:id/members", h.authMiddleware.Handle, h.CreateTeamMemberHandler)
}

func (h *teamHandler) CreateTeamHandler(c *gin.Context) {
	var requestData requests.CreateTeamRequest
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validation.Validate(requestData); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.ResponseError(err.Error()))
		return
	}

	team := &domain.Team{
		Name:      requestData.Name,
		ProjectID: requestData.ProjectID,
	}

	err := h.teamService.CreateTeam(c.Request.Context(), team)
	if err != nil {
		c.JSON(http.StatusInternalServerError, datatransfers.ResponseError("Failed to create team"))
		return
	}

	responseData := responses.TeamResponse{
		ID:        team.ID,
		Name:      team.Name,
		ProjectID: team.ProjectID,
		CreatedAt: team.CreatedAt.Format(time.RFC3339),
	}

	h.hub.SendMessage(team.ProjectID, ws.BaseResponse{
		Name: "team_created",
		Data: responseData,
	})

	c.JSON(http.StatusCreated, datatransfers.ResponseSuccess("Team created successfully", responseData))
}

func (h *teamHandler) GetTeamsHandler(c *gin.Context) {
	projectID := c.Query("project_id")

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
	id := c.Param("id")

	err := validation.ValidateUUID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.ResponseError("Invalid team ID"))
		return
	}

	team, teamMembers, err := h.teamService.GetTeamWithMembersByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, datatransfers.ResponseError("Failed to get team"))
		return
	}

	responseData := responses.TeamWithMembersResponse{
		ID:        team.ID,
		Name:      team.Name,
		ProjectID: team.ProjectID,
		CreatedAt: team.CreatedAt.Format(time.RFC3339),
		Members:   make([]responses.ProjectMemberResponse, len(teamMembers)),
	}

	for i, teamMember := range teamMembers {
		responseData.Members[i] = responses.ProjectMemberResponse{
			ID:        teamMember.ID,
			TeamID:    *teamMember.TeamID,
			UserID:    teamMember.UserID,
			CreatedAt: teamMember.CreatedAt.Format(time.RFC3339),
		}
	}

	c.JSON(http.StatusOK, datatransfers.ResponseSuccess("Team fetched successfully", responseData))
}

func (h *teamHandler) UpdateTeamHandler(c *gin.Context) {
	id := c.Param("id")

	err := validation.ValidateUUID(id)
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
		ID: id,
	}

	if requestData.Name != nil {
		team.Name = *requestData.Name
	}

	if requestData.Role != nil {
		team.Role = *requestData.Role
	}

	err = h.teamService.UpdateTeam(c.Request.Context(), team)
	if err != nil {
		c.JSON(http.StatusInternalServerError, datatransfers.ResponseError("Failed to update team"))
		return
	}

	responseData := responses.UpdateTeamResponse{
		ID:        team.ID,
		Name:      team.Name,
		Role:      team.Role,
		ProjectID: team.ProjectID,
	}

	h.hub.SendMessage(team.ProjectID, ws.BaseResponse{
		Name: "team_updated",
		Data: responseData,
	})

	c.JSON(http.StatusOK, datatransfers.ResponseSuccess("Team updated successfully", responseData))
}

func (h *teamHandler) DeleteTeamHandler(c *gin.Context) {
	id := c.Param("id")
	projectID := c.Query("project_id")

	err := validation.ValidateUUID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.ResponseError("Invalid team ID"))
		return
	}

	err = validation.ValidateUUID(projectID)
	if err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.ResponseError("Invalid project ID"))
		return
	}

	err = h.teamService.DeleteTeamByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, datatransfers.ResponseError("Failed to delete team"))
		return
	}

	responseData := responses.DeleteTeamResponse{
		ID: id,
	}

	h.hub.SendMessage(projectID, ws.BaseResponse{
		Name: "team_deleted",
		Data: responseData,
	})

	c.JSON(http.StatusOK, datatransfers.ResponseSuccess("Team deleted successfully", responseData))
}

func (h *teamHandler) CreateTeamMemberHandler(c *gin.Context) {
	id := c.Param("id")
	projectID := c.Query("project_id")

	var requestData requests.CreateProjectMemberRequest

	requestData.TeamID = id
	requestData.ProjectID = projectID

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validation.Validate(requestData); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.ResponseError(err.Error()))
		return
	}

	projectMember := &domain.ProjectMember{
		TeamID:    &id,
		UserID:    requestData.UserID,
		ProjectID: projectID,
	}

	err := h.teamService.CreateTeamMember(c.Request.Context(), projectMember)
	if err != nil {
		c.JSON(http.StatusInternalServerError, datatransfers.ResponseError("Failed to add team member"))
		return
	}

	responseData := responses.ProjectMemberResponse{
		ID:        projectMember.ID,
		TeamID:    *projectMember.TeamID,
		UserID:    projectMember.UserID,
		CreatedAt: projectMember.CreatedAt.Format(time.RFC3339),
	}

	h.hub.SendMessage(projectID, ws.BaseResponse{
		Name: "team_member_created",
		Data: responseData,
	})

	c.JSON(http.StatusCreated, datatransfers.ResponseSuccess("Team member added successfully", responseData))
}
