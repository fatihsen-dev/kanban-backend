package http

import (
	"net/http"
	"time"

	"github.com/fatihsen-dev/kanban-backend/internal/adapters/driver/http/datatransfers"
	"github.com/fatihsen-dev/kanban-backend/internal/adapters/driver/http/datatransfers/requests"
	"github.com/fatihsen-dev/kanban-backend/internal/adapters/driver/http/datatransfers/responses"
	middlewares "github.com/fatihsen-dev/kanban-backend/internal/adapters/driver/http/middleware"
	"github.com/fatihsen-dev/kanban-backend/internal/adapters/driver/validation"
	"github.com/fatihsen-dev/kanban-backend/internal/core/domain"
	ports "github.com/fatihsen-dev/kanban-backend/internal/core/ports/driver"
	"github.com/fatihsen-dev/kanban-backend/pkg/helpers"
	"github.com/fatihsen-dev/kanban-backend/pkg/jwt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type authHandler struct {
	userService    ports.UserService
	authMiddleware *middlewares.AuthnMiddleware
}

func NewAuthHandler(userService ports.UserService, authMiddleware *middlewares.AuthnMiddleware) *authHandler {
	return &authHandler{userService: userService, authMiddleware: authMiddleware}
}

func (h *authHandler) RegisterAuthRouter(r *gin.Engine) {
	r.POST("/auth/login", h.LoginHandler)
	r.POST("/auth/register", h.RegisterHandler)
	r.GET("/auth/me", h.authMiddleware.Handle, h.AuthUser)
}

func (h *authHandler) LoginHandler(c *gin.Context) {

	var requestData requests.UserLoginRequest

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.ResponseError("Invalid credentials"))
		return
	}

	if err := validation.Validate(requestData); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.ResponseError(err.Error()))
		return
	}

	user, err := h.userService.GetUserByEmail(c.Request.Context(), requestData.Email)
	if err != nil || user == nil {
		c.JSON(http.StatusBadRequest, datatransfers.ResponseError("Invalid email or password"))
		return
	}

	if err := helpers.ValidateHash(user.PasswordHash, requestData.Password); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.ResponseError("Invalid email or password"))
		return
	}

	token, err := jwt.GenerateToken(user.ID, user.Name, user.Email, user.IsAdmin)
	if err != nil {
		zap.L().Error("Failed to generate token", zap.Error(err))
		c.JSON(http.StatusInternalServerError, datatransfers.ResponseError("Something went wrong"))
		return
	}

	responseData := responses.UserLoginResponse{
		Token: token,
		User: responses.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			IsAdmin:   user.IsAdmin,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
		},
	}

	c.JSON(http.StatusOK, datatransfers.ResponseSuccess("Logged in successfully", responseData))
}

func (h *authHandler) RegisterHandler(c *gin.Context) {

	var requestData requests.UserRegisterRequest

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.ResponseError("Invalid request data"))
		return
	}

	if err := validation.Validate(requestData); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.ResponseError(err.Error()))
		return
	}

	_, err := h.userService.GetUserByEmail(c.Request.Context(), requestData.Email)
	if err == nil {
		c.JSON(http.StatusBadRequest, datatransfers.ResponseError("Email already in use"))
		return
	}

	hashedPassword, err := helpers.GenerateHash(requestData.Password)
	if err != nil {
		zap.L().Error("Failed to hash password", zap.Error(err))
		c.JSON(http.StatusInternalServerError, datatransfers.ResponseError("Something went wrong"))
		return
	}

	user := &domain.User{
		Name:         requestData.Name,
		Email:        requestData.Email,
		PasswordHash: string(hashedPassword),
	}

	err = h.userService.CreateUser(c.Request.Context(), user)
	if err != nil {
		zap.L().Error("Failed to create user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, datatransfers.ResponseError("Something went wrong"))
		return
	}

	token, err := jwt.GenerateToken(user.ID, user.Name, user.Email, user.IsAdmin)
	if err != nil {
		zap.L().Error("Failed to generate token", zap.Error(err))
		c.JSON(http.StatusInternalServerError, datatransfers.ResponseError("Something went wrong"))
		return
	}

	responseData := responses.UserRegisterResponse{
		Token: token,
		User: responses.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			IsAdmin:   user.IsAdmin,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
		},
	}

	c.JSON(http.StatusOK, datatransfers.ResponseSuccess("User registered successfully", responseData))
}

func (h *authHandler) AuthUser(c *gin.Context) {
	userClaims := c.MustGet("user").(*jwt.UserClaims)

	user, err := h.userService.GetUserByID(c.Request.Context(), userClaims.UserID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, datatransfers.ResponseError("Unauthorized"))
		return
	}

	c.Set("user", user)

	responseData := responses.UserAuthResponse{
		User: responses.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			IsAdmin:   user.IsAdmin,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
		},
	}

	c.JSON(http.StatusOK, datatransfers.ResponseSuccess("Authenticated user", responseData))
}
