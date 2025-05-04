package http

import (
	"net/http"

	"github.com/fatihsen-dev/kanban-backend/internal/adapters/driver/http/datatransfers"
	middlewares "github.com/fatihsen-dev/kanban-backend/internal/adapters/driver/http/middleware"
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
	r.Use(h.authMiddleware.Handle(false)).GET("/projects/:project_id/members/online", h.GetOnlineProjectMembersHandler)
}

func (h *projectMemberHandler) GetOnlineProjectMembersHandler(c *gin.Context) {
	projectID := c.Param("project_id")
	onlineUsers := h.hub.GetOnlineUsers(projectID)

	c.JSON(http.StatusOK, datatransfers.ResponseSuccess("Online project members fetched successfully", onlineUsers))
}
