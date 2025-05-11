package http

import (
	"fmt"
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

type taskHandler struct {
	taskService            ports.TaskService
	authMiddleware         *middlewares.AuthnMiddleware
	projectAuthzMiddleware *middlewares.ProjectAuthzMiddleware
	hub                    *ws.Hub
}

func NewTaskHandler(taskService ports.TaskService, authMiddleware *middlewares.AuthnMiddleware, projectAuthzMiddleware *middlewares.ProjectAuthzMiddleware, hub *ws.Hub) *taskHandler {
	return &taskHandler{taskService: taskService, authMiddleware: authMiddleware, projectAuthzMiddleware: projectAuthzMiddleware, hub: hub}
}

func (h *taskHandler) RegisterTaskRouter(r *gin.Engine) {

	taskGroup := r.Group("/projects/:project_id/tasks")

	taskGroup.Use(h.authMiddleware.Handle(false))

	taskGroup.POST("", h.projectAuthzMiddleware.Handle(middlewares.Member), h.CreateTaskHandler)
	taskGroup.GET("", h.projectAuthzMiddleware.Handle(middlewares.Member), h.GetTasksHandler)
	taskGroup.GET("/:task_id", h.projectAuthzMiddleware.Handle(middlewares.Member), h.GetTaskHandler)
	taskGroup.PUT("/:task_id", h.projectAuthzMiddleware.Handle(middlewares.Member), h.UpdateTaskHandler)
	taskGroup.DELETE("/:task_id", h.projectAuthzMiddleware.Handle(middlewares.Member), h.DeleteTaskHandler)
}

func (h *taskHandler) CreateTaskHandler(c *gin.Context) {
	projectID := c.Param("project_id")

	var requestData requests.TaskCreateRequest

	requestData.ProjectID = projectID

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.ResponseError("Invalid request data"))
		return
	}

	if err := validation.Validate(requestData); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.ResponseError(err.Error()))
		return
	}

	task := &domain.Task{
		Title:     requestData.Title,
		Content:   requestData.Content,
		ProjectID: requestData.ProjectID,
		ColumnID:  requestData.ColumnID,
	}

	err := h.taskService.CreateTask(c.Request.Context(), task)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, datatransfers.ResponseError("Failed to create task"))
		return
	}

	responseData := responses.TaskResponse{
		ID:        task.ID,
		Title:     task.Title,
		Content:   task.Content,
		ProjectID: task.ProjectID,
		ColumnID:  task.ColumnID,
		CreatedAt: task.CreatedAt.Format(time.RFC3339),
	}

	h.hub.SendMessageToProject(task.ProjectID, ws.BaseResponse{
		Name: ws.EventNameTaskCreated,
		Data: responseData,
	})

	c.JSON(http.StatusCreated, datatransfers.ResponseSuccess("Task created successfully", responseData))
}

func (h *taskHandler) GetTaskHandler(c *gin.Context) {
	id := c.Param("task_id")

	err := validation.ValidateUUID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.ResponseError("Invalid task ID"))
		return
	}

	task, err := h.taskService.GetTaskByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, datatransfers.ResponseError("Task not found"))
		return
	}

	responseData := responses.TaskResponse{
		ID:        task.ID,
		Title:     task.Title,
		Content:   task.Content,
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
			Content:   task.Content,
			ProjectID: task.ProjectID,
			ColumnID:  task.ColumnID,
			CreatedAt: task.CreatedAt.Format(time.RFC3339),
		}
	}

	c.JSON(http.StatusOK, datatransfers.ResponseSuccess("Tasks fetched successfully", responseData))
}

func (h *taskHandler) UpdateTaskHandler(c *gin.Context) {
	id := c.Param("task_id")
	projectID := c.Param("project_id")

	err := validation.ValidateUUID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.ResponseError("Invalid task ID"))
		return
	}

	err = validation.ValidateUUID(projectID)
	if err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.ResponseError("Invalid project ID"))
		return
	}

	var requestData requests.TaskUpdateRequest
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.ResponseError("Invalid request data"))
		return
	}

	if err := validation.Validate(requestData); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.ResponseError(err.Error()))
		return
	}

	task := &domain.Task{
		ID: id,
	}

	responseData := responses.TaskUpdateResponse{
		ID: task.ID,
	}

	if requestData.Title != nil {
		task.Title = *requestData.Title
		responseData.Title = task.Title
	}

	if requestData.ColumnID != nil {
		task.ColumnID = *requestData.ColumnID
		responseData.ColumnID = task.ColumnID
	}

	if requestData.Content != nil {
		task.Content = requestData.Content
		responseData.Content = task.Content
	}

	err = h.taskService.UpdateTask(c.Request.Context(), task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, datatransfers.ResponseError("Failed to update task"))
		return
	}

	if requestData.ColumnID != nil {
		h.hub.SendMessageToProject(projectID, ws.BaseResponse{
			Name: ws.EventNameTaskMoved,
			Data: responseData,
		})
	} else {
		h.hub.SendMessageToProject(projectID, ws.BaseResponse{
			Name: ws.EventNameTaskUpdated,
			Data: responseData,
		})
	}

	c.JSON(http.StatusOK, datatransfers.ResponseSuccess("Task updated successfully", responseData))
}

func (h *taskHandler) DeleteTaskHandler(c *gin.Context) {
	id := c.Param("task_id")
	projectID := c.Param("project_id")
	err := validation.ValidateUUID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.ResponseError("Invalid task ID"))
		return
	}

	err = validation.ValidateUUID(projectID)
	if err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.ResponseError("Invalid project ID"))
		return
	}

	err = h.taskService.DeleteTask(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, datatransfers.ResponseError("Failed to delete task"))
		return
	}

	responseData := responses.TaskDeleteResponse{
		ID: id,
	}

	h.hub.SendMessageToProject(projectID, ws.BaseResponse{
		Name: ws.EventNameTaskDeleted,
		Data: responseData,
	})

	c.JSON(http.StatusOK, datatransfers.ResponseSuccess("Task deleted successfully", responseData))
}
