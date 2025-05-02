package http

import (
	"net/http"
	"time"

	"github.com/fatihsen-dev/kanban-backend/internal/adapters/driver/http/datatransfers"
	"github.com/fatihsen-dev/kanban-backend/internal/adapters/driver/http/datatransfers/responses"
	middlewares "github.com/fatihsen-dev/kanban-backend/internal/adapters/driver/http/middleware"
	ports "github.com/fatihsen-dev/kanban-backend/internal/core/ports/driver"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type userHandler struct {
	userService    ports.UserService
	authMiddleware *middlewares.AuthnMiddleware
}

func NewUserHandler(userService ports.UserService, authMiddleware *middlewares.AuthnMiddleware) *userHandler {
	return &userHandler{userService: userService, authMiddleware: authMiddleware}
}

func (h *userHandler) RegisterUserRouter(r *gin.Engine) {
	r.GET("/users", h.authMiddleware.Handle(false), h.GetUsersHandler)
}

func (h *userHandler) GetUsersHandler(c *gin.Context) {
	query := c.Query("query")

	if len(query) > 0 {
		users, err := h.userService.GetUsersByQuery(c.Request.Context(), query)
		if err != nil {
			zap.L().Error("Failed to get users by query", zap.Error(err))
			c.JSON(http.StatusInternalServerError, datatransfers.ResponseError(err.Error()))
			return
		}

		userResponses := make([]*responses.UserResponse, len(users))
		for i, user := range users {
			userResponses[i] = &responses.UserResponse{
				ID:        user.ID,
				Name:      user.Name,
				Email:     user.Email,
				CreatedAt: user.CreatedAt.Format(time.RFC3339),
			}
		}

		c.JSON(http.StatusOK, datatransfers.ResponseSuccess("Users fetched successfully", userResponses))
		return
	}

	users, err := h.userService.GetUsers(c.Request.Context())

	if err != nil {
		c.JSON(http.StatusInternalServerError, datatransfers.ResponseError(err.Error()))
		return
	}

	userResponses := make([]*responses.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = &responses.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
		}
	}

	c.JSON(http.StatusOK, datatransfers.ResponseSuccess("Users fetched successfully", userResponses))
}
