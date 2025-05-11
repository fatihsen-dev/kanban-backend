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
	"github.com/fatihsen-dev/kanban-backend/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type invitationHandler struct {
	invitationService      ports.InvitationService
	authMiddleware         *middlewares.AuthnMiddleware
	projectAuthzMiddleware *middlewares.ProjectAuthzMiddleware
	hub                    *ws.Hub
}

func NewInvitationHandler(invitationService ports.InvitationService, authMiddleware *middlewares.AuthnMiddleware, projectAuthzMiddleware *middlewares.ProjectAuthzMiddleware, hub *ws.Hub) *invitationHandler {
	return &invitationHandler{invitationService: invitationService, authMiddleware: authMiddleware, projectAuthzMiddleware: projectAuthzMiddleware, hub: hub}
}

func (h *invitationHandler) RegisterInvitationRouter(r *gin.Engine) {
	invitationGroup := r.Group("/invitations")

	invitationGroup.Use(h.authMiddleware.Handle(false))

	invitationGroup.GET("", h.GetInvitationsHandler)
	invitationGroup.PUT("/:invitation_id", h.UpdateInvitationStatusHandler)
	invitationGroup.POST("/:project_id", h.projectAuthzMiddleware.Handle(middlewares.Admin), h.CreateInvitationHandler)
}

func (h *invitationHandler) CreateInvitationHandler(c *gin.Context) {
	user := c.MustGet("user").(*jwt.UserClaims)
	projectID := c.Param("project_id")
	var request requests.InvitationCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.ResponseError(err.Error()))
		return
	}

	request.InviterID = user.ID
	request.ProjectID = projectID
	if err := validation.Validate(request); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.ResponseError(err.Error()))
		return
	}

	inviteeIDs := request.InviteeIDs

	invitations := make([]*domain.Invitation, len(inviteeIDs))
	for i, inviteeID := range inviteeIDs {
		invitations[i] = &domain.Invitation{
			InviterID: request.InviterID,
			InviteeID: inviteeID,
			ProjectID: request.ProjectID,
			Message:   request.Message,
			Status:    domain.InvitationStatusPending,
		}
	}

	successInvitations, err := h.invitationService.CreateInvitations(c.Request.Context(), invitations)
	if err != nil {
		c.JSON(http.StatusInternalServerError, datatransfers.ResponseError(err.Error()))
		return
	}

	responseData := make([]responses.InvitationResponse, len(successInvitations))

	for _, invitation := range successInvitations {
		h.hub.SendMessageToUser(invitation.Invitee.ID, ws.BaseResponse{
			Name: ws.EventNameInvitationCreated,
			Data: invitation,
		})
	}

	c.JSON(http.StatusCreated, datatransfers.ResponseSuccess("Invitation created successfully", responseData))
}

func (h *invitationHandler) GetInvitationsHandler(c *gin.Context) {
	user := c.MustGet("user").(*jwt.UserClaims)

	invitations, err := h.invitationService.GetInvitations(c.Request.Context(), user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, datatransfers.ResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, datatransfers.ResponseSuccess("Invitations fetched successfully", invitations))
}

func (h *invitationHandler) UpdateInvitationStatusHandler(c *gin.Context) {
	authUser := c.MustGet("user").(*jwt.UserClaims)

	var request requests.InvitationUpdateStatusRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.ResponseError(err.Error()))
		return
	}
	request.ID = c.Param("invitation_id")
	request.UserID = authUser.ID
	if err := validation.Validate(request); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.ResponseError(err.Error()))
		return
	}

	projectMember, user, err := h.invitationService.UpdateInvitationStatus(c.Request.Context(), request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, datatransfers.ResponseError(err.Error()))
		return
	}

	if projectMember != nil {
		responseData := responses.ProjectMemberWithUserResponse{
			ID:        projectMember.ID,
			UserID:    projectMember.UserID,
			Role:      string(projectMember.Role),
			TeamID:    projectMember.TeamID,
			ProjectID: projectMember.ProjectID,
			CreatedAt: projectMember.CreatedAt.Format(time.RFC3339),
			User: responses.UserResponse{
				ID:        user.ID,
				Name:      user.Name,
				Email:     user.Email,
				IsAdmin:   user.IsAdmin,
				CreatedAt: user.CreatedAt.Format(time.RFC3339),
			},
		}

		h.hub.SendMessageToProject(projectMember.ProjectID, ws.BaseResponse{
			Name: ws.EventNameProjectMemberCreated,
			Data: responseData,
		})

		c.JSON(http.StatusOK, datatransfers.ResponseSuccess("Invitation accepted successfully", responseData))
		return
	}

	c.JSON(http.StatusOK, datatransfers.ResponseSuccess("Invitation rejected successfully", nil))
}
