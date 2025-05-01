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
	"github.com/fatihsen-dev/kanban-backend/pkg/jwt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type invitationHandler struct {
	invitationService ports.InvitationService
	authMiddleware    *middlewares.AuthnMiddleware
	hub               *ws.Hub
}

func NewInvitationHandler(invitationService ports.InvitationService, authMiddleware *middlewares.AuthnMiddleware, hub *ws.Hub) *invitationHandler {
	return &invitationHandler{invitationService: invitationService, authMiddleware: authMiddleware, hub: hub}
}

func (h *invitationHandler) RegisterInvitationRouter(r *gin.Engine) {
	r.POST("/invitations", h.authMiddleware.Handle(false), h.CreateInvitationHandler)
	r.GET("/invitations", h.authMiddleware.Handle(false), h.GetInvitationsHandler)
}

func (h *invitationHandler) CreateInvitationHandler(c *gin.Context) {
	user := c.MustGet("user").(*jwt.UserClaims)

	var request requests.InvitationCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.ResponseError(err.Error()))
		return
	}

	request.InviterID = user.ID

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
			Name: "invitation_notification",
			Data: invitation,
		})
	}

	c.JSON(http.StatusCreated, datatransfers.ResponseSuccess("Invitation created successfully", responseData))
}

func (h *invitationHandler) GetInvitationsHandler(c *gin.Context) {
	user := c.MustGet("user").(*jwt.UserClaims)

	invitations, err := h.invitationService.GetInvitations(c.Request.Context(), user.ID)
	if err != nil {
		zap.L().Error("Failed to get invitations", zap.Error(err))
		c.JSON(http.StatusInternalServerError, datatransfers.ResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, datatransfers.ResponseSuccess("Invitations fetched successfully", invitations))
}
