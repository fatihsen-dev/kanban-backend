package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fatihsen-dev/kanban-backend/internal/adapters/driver/http/datatransfers"
	"github.com/fatihsen-dev/kanban-backend/internal/adapters/driver/http/datatransfers/requests"
	"github.com/fatihsen-dev/kanban-backend/internal/adapters/driver/http/datatransfers/responses"
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

	var requestData requests.TaskCreateRequest

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.ResponseError("Invalid request data"))
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
		c.JSON(http.StatusInternalServerError, datatransfers.ResponseError("Failed to create task"))
		return
	}

	h.hub.SendMessage(task.ProjectID, ws.BaseResponse{
		Name: "task_created",
		Data: task,
	})

	responseData := responses.TaskResponse{
		ID:        task.ID,
		Title:     task.Title,
		ProjectID: task.ProjectID,
		ColumnID:  task.ColumnID,
		CreatedAt: task.CreatedAt.Format(time.RFC3339),
	}

	c.JSON(http.StatusCreated, datatransfers.ResponseSuccess("Task created successfully", responseData))
}

func (h *taskHandler) GetTaskHandler(c *gin.Context) {
	id := c.Param("id")

	task, err := h.taskService.GetTaskByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, datatransfers.ResponseError("Failed to get task"))
		return
	}

	responseData := responses.TaskResponse{
		ID:        task.ID,
		Title:     task.Title,
		ProjectID: task.ProjectID,
		ColumnID:  task.ColumnID,
		CreatedAt: task.CreatedAt.Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, datatransfers.ResponseSuccess("Task fetched successfully", responseData))
}

func (h *taskHandler) GetTasksHandler(c *gin.Context) {
	tasks, err := h.taskService.GetTasks(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, datatransfers.ResponseError("Failed to get tasks"))
		return
	}

	responseData := make([]responses.TaskResponse, len(tasks))
	for i, task := range tasks {
		responseData[i] = responses.TaskResponse{
			ID:        task.ID,
			Title:     task.Title,
			ProjectID: task.ProjectID,
			ColumnID:  task.ColumnID,
			CreatedAt: task.CreatedAt.Format(time.RFC3339),
		}
	}

	c.JSON(http.StatusOK, datatransfers.ResponseSuccess("Tasks fetched successfully", responseData))
}
