package http

import (
	"fmt"
	"net/http"

	middlewares "github.com/fatihsen-dev/kanban-backend/internal/adapters/driver/http/middleware"
	"github.com/fatihsen-dev/kanban-backend/internal/adapters/driver/ws"
	"github.com/fatihsen-dev/kanban-backend/internal/core/domain"
	ports "github.com/fatihsen-dev/kanban-backend/internal/core/ports/driver"
	"github.com/gin-gonic/gin"
)

type taskHandler struct {
	taskService    ports.TaskService
	authMiddleware *middlewares.AuthMiddleware
	hub            *ws.Hub
}

func NewTaskHandler(taskService ports.TaskService, authMiddleware *middlewares.AuthMiddleware, hub *ws.Hub) *taskHandler {
	return &taskHandler{taskService: taskService, authMiddleware: authMiddleware, hub: hub}
}

func (h *taskHandler) RegisterTaskRouter(r *gin.Engine) {
	r.POST("/tasks", h.authMiddleware.Handle, h.CreateTaskHandler)
	r.GET("/tasks", h.authMiddleware.Handle, h.GetTasksHandler)
	r.GET("/tasks/:id", h.authMiddleware.Handle, h.GetTaskHandler)
}

func (h *taskHandler) CreateTaskHandler(c *gin.Context) {

	var requestData struct {
		Title     string `json:"title"`
		ProjectID string `json:"project_id"`
		ColumnID  string `json:"column_id"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	task := &domain.Task{
		Title:     requestData.Title,
		ProjectID: requestData.ProjectID,
		ColumnID:  requestData.ColumnID,
	}

	err := h.taskService.CreateTask(c.Request.Context(), task)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	h.hub.SendMessage(task.ProjectID, ws.BaseResponse{
		Name: "task_created",
		Data: task,
	})

	c.JSON(http.StatusCreated, gin.H{"message": "Project created successfully"})
}

func (h *taskHandler) GetTaskHandler(c *gin.Context) {
	id := c.Param("id")

	task, err := h.taskService.GetTaskByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get task"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *taskHandler) GetTasksHandler(c *gin.Context) {
	tasks, err := h.taskService.GetTasks(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tasks"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}
