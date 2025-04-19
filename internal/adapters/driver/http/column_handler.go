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

type columnHandler struct {
	columnService  ports.ColumnService
	authMiddleware *middlewares.AuthnMiddleware
	hub            *ws.Hub
}

func NewColumnHandler(columnService ports.ColumnService, authMiddleware *middlewares.AuthnMiddleware, hub *ws.Hub) *columnHandler {
	return &columnHandler{columnService: columnService, authMiddleware: authMiddleware, hub: hub}
}

func (h *columnHandler) RegisterColumnRouter(r *gin.Engine) {
	r.POST("/columns", h.authMiddleware.Handle, h.CreateColumnHandler)
	r.GET("/columns", h.authMiddleware.Handle, h.GetColumnsHandler)
	r.GET("/columns/:id", h.authMiddleware.Handle, h.GetColumnHandler)
	r.GET("/columns/:id/tasks", h.authMiddleware.Handle, h.GetColumnWithTasksHandler)
	r.PUT("/columns/:id", h.authMiddleware.Handle, h.UpdateColumnHandler)
	r.DELETE("/columns/:id", h.authMiddleware.Handle, h.DeleteColumnHandler)
}

func (h *columnHandler) CreateColumnHandler(c *gin.Context) {

	var requestData requests.ColumnCreateRequest

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.ResponseError("Invalid request data"))
		return
	}

	if err := validation.Validate(requestData); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.ResponseError(err.Error()))
		return
	}

	column := &domain.Column{
		Name:      requestData.Name,
		ProjectID: requestData.ProjectID,
	}

	err := h.columnService.CreateColumn(c.Request.Context(), column)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, datatransfers.ResponseError("Failed to create column"))
		return
	}

	responseData := responses.ColumnResponse{
		ID:        column.ID,
		Name:      column.Name,
		ProjectID: column.ProjectID,
		CreatedAt: column.CreatedAt.Format(time.RFC3339),
	}

	h.hub.SendMessage(column.ProjectID, ws.BaseResponse{
		Name: "column_created",
		Data: responseData,
	})

	c.JSON(http.StatusCreated, datatransfers.ResponseSuccess("Column created successfully", responseData))
}

func (h *columnHandler) GetColumnHandler(c *gin.Context) {
	id := c.Param("id")

	err := validation.ValidateUUID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.ResponseError("Invalid column ID"))
		return
	}

	column, err := h.columnService.GetColumnByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, datatransfers.ResponseError("Failed to get column"))
		return
	}

	responseData := responses.ColumnResponse{
		ID:        column.ID,
		Name:      column.Name,
		ProjectID: column.ProjectID,
		CreatedAt: column.CreatedAt.Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, datatransfers.ResponseSuccess("Column fetched successfully", responseData))
}

func (h *columnHandler) GetColumnWithTasksHandler(c *gin.Context) {
	id := c.Param("id")

	err := validation.ValidateUUID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.ResponseError("Invalid column ID"))
		return
	}

	column, tasks, err := h.columnService.GetColumnWithTasks(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, datatransfers.ResponseError("Failed to get column with tasks"))
		return
	}

	responseData := responses.ColumnWithDetailsResponse{
		ID:        column.ID,
		Name:      column.Name,
		CreatedAt: column.CreatedAt.Format(time.RFC3339),
		Tasks:     make([]responses.TaskResponse, len(tasks)),
	}

	for i, task := range tasks {
		responseData.Tasks[i] = responses.TaskResponse{
			ID:        task.ID,
			Title:     task.Title,
			ProjectID: task.ProjectID,
			ColumnID:  task.ColumnID,
			CreatedAt: task.CreatedAt.Format(time.RFC3339),
		}
	}

	c.JSON(http.StatusOK, datatransfers.ResponseSuccess("Column with tasks fetched successfully", responseData))
}

func (h *columnHandler) GetColumnsHandler(c *gin.Context) {
	columns, err := h.columnService.GetColumns(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, datatransfers.ResponseError("Failed to get columns"))
		return
	}

	responseData := make([]responses.ColumnResponse, len(columns))
	for i, column := range columns {
		responseData[i] = responses.ColumnResponse{
			ID:        column.ID,
			Name:      column.Name,
			ProjectID: column.ProjectID,
			CreatedAt: column.CreatedAt.Format(time.RFC3339),
		}
	}

	c.JSON(http.StatusOK, datatransfers.ResponseSuccess("Columns fetched successfully", responseData))
}

func (h *columnHandler) UpdateColumnHandler(c *gin.Context) {
	id := c.Param("id")
	projectID := c.Query("project_id")

	err := validation.ValidateUUID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.ResponseError("Invalid column ID"))
		return
	}

	err = validation.ValidateUUID(projectID)
	if err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.ResponseError("Invalid project ID"))
		return
	}

	var requestData requests.ColumnUpdateRequest
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.ResponseError("Invalid request data"))
		return
	}

	if err := validation.Validate(requestData); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.ResponseError(err.Error()))
		return
	}

	column := &domain.Column{
		ID:   id,
		Name: requestData.Name,
	}

	err = h.columnService.UpdateColumn(c.Request.Context(), column)
	if err != nil {
		c.JSON(http.StatusInternalServerError, datatransfers.ResponseError("Failed to update column"))
		return
	}

	responseData := responses.ColumnUpdateResponse{
		ID:   column.ID,
		Name: column.Name,
	}

	h.hub.SendMessage(projectID, ws.BaseResponse{
		Name: "column_updated",
		Data: responseData,
	})

	c.JSON(http.StatusOK, datatransfers.ResponseSuccess("Column updated successfully", responseData))
}

func (h *columnHandler) DeleteColumnHandler(c *gin.Context) {
	id := c.Param("id")
	projectID := c.Query("project_id")

	err := validation.ValidateUUID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.ResponseError("Invalid column ID"))
		return
	}

	err = validation.ValidateUUID(projectID)
	if err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.ResponseError("Invalid project ID"))
		return
	}

	err = h.columnService.DeleteColumn(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, datatransfers.ResponseError("Failed to delete column"))
		return
	}

	responseData := responses.ColumnDeleteResponse{
		ID: id,
	}

	h.hub.SendMessage(projectID, ws.BaseResponse{
		Name: "column_deleted",
		Data: responseData,
	})

	c.JSON(http.StatusOK, datatransfers.ResponseSuccess("Column deleted successfully", responseData))
}
