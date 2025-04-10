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

type columnHandler struct {
	columnService  ports.ColumnService
	authMiddleware *middlewares.AuthMiddleware
	hub            *ws.Hub
}

func NewColumnHandler(columnService ports.ColumnService, authMiddleware *middlewares.AuthMiddleware, hub *ws.Hub) *columnHandler {
	return &columnHandler{columnService: columnService, authMiddleware: authMiddleware, hub: hub}
}

func (h *columnHandler) RegisterColumnRouter(r *gin.Engine) {
	r.POST("/columns", h.authMiddleware.Handle, h.CreateColumnHandler)
	r.GET("/columns", h.authMiddleware.Handle, h.GetColumnsHandler)
	r.GET("/columns/:id", h.authMiddleware.Handle, h.GetColumnHandler)
}

func (h *columnHandler) CreateColumnHandler(c *gin.Context) {

	var requestData struct {
		Name      string `json:"name"`
		ProjectID string `json:"project_id"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	column := &domain.Column{
		Name:      requestData.Name,
		ProjectID: requestData.ProjectID,
	}

	err := h.columnService.CreateColumn(c.Request.Context(), column)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create column"})
		return
	}

	h.hub.SendMessage(column.ProjectID, ws.BaseResponse{
		Name: "column_created",
		Data: column,
	})

	c.JSON(http.StatusCreated, gin.H{"message": "Column created successfully"})
}

func (h *columnHandler) GetColumnHandler(c *gin.Context) {
	id := c.Param("id")

	column, err := h.columnService.GetColumnByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get column"})
		return
	}

	c.JSON(http.StatusOK, column)
}

func (h *columnHandler) GetColumnsHandler(c *gin.Context) {
	columns, err := h.columnService.GetColumns(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get columns"})
		return
	}

	c.JSON(http.StatusOK, columns)
}
