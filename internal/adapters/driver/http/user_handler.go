package http

import (
	"fmt"
	"net/http"

	middlewares "github.com/fatihsen-dev/kanban-backend/internal/adapters/driver/http/middleware"
	"github.com/fatihsen-dev/kanban-backend/internal/core/domain"
	ports "github.com/fatihsen-dev/kanban-backend/internal/core/ports/driver"
	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService    ports.UserService
	authMiddleware *middlewares.AuthMiddleware
}

func NewUserHandler(userService ports.UserService, authMiddleware *middlewares.AuthMiddleware) *userHandler {
	return &userHandler{userService: userService, authMiddleware: authMiddleware}
}

func (h *userHandler) RegisterUserRouter(r *gin.Engine) {
	// r.POST("/users", h.authMiddleware.Handle, h.CreateUserHandler)
	// r.GET("/users", h.authMiddleware.Handle, h.GetUsersHandler)
	// r.GET("/users/:id", h.authMiddleware.Handle, h.GetUserHandler)
}

func (h *userHandler) CreateUserHandler(c *gin.Context) {

	var requestData struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	user := &domain.User{
		Name:         requestData.Name,
		Email:        requestData.Email,
		PasswordHash: requestData.Password, // TODO: hash password
	}

	err := h.userService.CreateUser(c.Request.Context(), user)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func (h *userHandler) GetUserHandler(c *gin.Context) {
	id := c.Query("id")

	user, err := h.userService.GetUserByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *userHandler) GetUsersHandler(c *gin.Context) {
	users, err := h.userService.GetUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users"})
		return
	}

	c.JSON(http.StatusOK, users)
}
