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
	r.POST("/projects/:project_id/columns", h.authMiddleware.Handle, h.CreateColumnHandler)
	r.GET("/projects/:project_id/columns", h.authMiddleware.Handle, h.GetColumnsHandler)
	r.GET("/projects/:project_id/columns/:column_id", h.authMiddleware.Handle, h.GetColumnHandler)
	r.PUT("/projects/:project_id/columns/:column_id", h.authMiddleware.Handle, h.UpdateColumnHandler)
	r.DELETE("/projects/:project_id/columns/:column_id", h.authMiddleware.Handle, h.DeleteColumnHandler)
}

func (h *columnHandler) CreateColumnHandler(c *gin.Context) {
	projectID := c.Param("project_id")

	var requestData requests.ColumnCreateRequest

	requestData.ProjectID = projectID

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
		Color:     requestData.Color,
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
		Color:     column.Color,
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
	id := c.Param("column_id")

	err := validation.ValidateUUID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.ResponseError("Invalid column ID"))
		return
	}

	column, tasks, err := h.columnService.GetColumnWithDetails(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, datatransfers.ResponseError("Failed to get column with tasks"))
		return
	}

	responseData := responses.ColumnWithDetailsResponse{
		ID:        column.ID,
		Name:      column.Name,
		Color:     column.Color,
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
	projectID := c.Param("project_id")

	columns, err := h.columnService.GetColumnsByProjectID(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, datatransfers.ResponseError("Failed to get columns"))
		return
	}

	responseData := make([]responses.ColumnResponse, len(columns))
	for i, column := range columns {
		responseData[i] = responses.ColumnResponse{
			ID:        column.ID,
			Name:      column.Name,
			Color:     column.Color,
			ProjectID: column.ProjectID,
			CreatedAt: column.CreatedAt.Format(time.RFC3339),
		}
	}

	c.JSON(http.StatusOK, datatransfers.ResponseSuccess("Columns fetched successfully", responseData))
}

func (h *columnHandler) UpdateColumnHandler(c *gin.Context) {
	id := c.Param("column_id")
	projectID := c.Param("project_id")

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
		ID:    id,
		Color: requestData.Color,
	}

	if requestData.Name != nil {
		column.Name = *requestData.Name
	}

	err = h.columnService.UpdateColumn(c.Request.Context(), column)
	if err != nil {
		c.JSON(http.StatusInternalServerError, datatransfers.ResponseError("Failed to update column"))
		return
	}

	responseData := responses.ColumnUpdateResponse{
		ID:    column.ID,
		Name:  column.Name,
		Color: column.Color,
	}

	h.hub.SendMessage(projectID, ws.BaseResponse{
		Name: "column_updated",
		Data: responseData,
	})

	c.JSON(http.StatusOK, datatransfers.ResponseSuccess("Column updated successfully", responseData))
}

func (h *columnHandler) DeleteColumnHandler(c *gin.Context) {
	id := c.Param("column_id")
	projectID := c.Param("project_id")

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
