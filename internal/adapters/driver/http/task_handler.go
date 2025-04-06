package http

import (
	"fmt"
	"net/http"

	"github.com/fatihsen-dev/kanban-backend/internal/adapters/driver/ws"
	"github.com/fatihsen-dev/kanban-backend/internal/core/domain"
	ports "github.com/fatihsen-dev/kanban-backend/internal/core/ports/driver"
	"github.com/gin-gonic/gin"
)

type taskHandler struct {
	taskService ports.TaskService
	hub         *ws.Hub
}

func NewTaskHandler(taskService ports.TaskService, hub *ws.Hub) *taskHandler {
	return &taskHandler{taskService: taskService, hub: hub}
}

func (h *taskHandler) RegisterTaskRouter(r *gin.Engine) {
	r.POST("/tasks", h.CreateTaskHandler)
	r.GET("/tasks", h.GetTasksHandler)
	r.GET("/tasks/:id", h.GetTaskHandler)
}

func (h *taskHandler) CreateTaskHandler(c *gin.Context) {

	var requestData struct {
		Title string `json:"title"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	task := &domain.Task{
		Title: requestData.Title,
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
	id := c.Query("id")

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
