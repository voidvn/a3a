package handlers

import (
	"net/http"

	"s4s-backend/internal/modules/connection/models"
	"s4s-backend/internal/modules/connection/services"

	"github.com/gin-gonic/gin"
)

type ConnectionHandler struct {
	connectionService services.ConnectionService
}

func NewConnectionHandler(service services.ConnectionService) *ConnectionHandler {
	return &ConnectionHandler{
		connectionService: service,
	}
}

type CreateConnectionRequest struct {
	Name        string                 `json:"name" binding:"required"`
	ServiceName string                 `json:"serviceName" binding:"required"`
	Credentials map[string]interface{} `json:"credentials" binding:"required"`
}

func (h *ConnectionHandler) CreateConnection(c *gin.Context) {
	var req CreateConnectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("userID")
	conn := &models.Connection{
		UserID:         userID,
		ConnectionName: req.Name,
		ServiceName:    req.ServiceName,
		Credentials:    req.Credentials,
		IsActive:       true,
	}

	if err := h.connectionService.CreateConnection(c.Request.Context(), conn); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, conn)
}

func (h *ConnectionHandler) GetConnection(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("userID")

	conn, err := h.connectionService.GetConnection(c.Request.Context(), id, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Connection not found"})
		return
	}

	c.JSON(http.StatusOK, conn)
}

func (h *ConnectionHandler) ListConnections(c *gin.Context) {
	userID := c.GetString("userID")

	connections, err := h.connectionService.ListConnections(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, connections)
}

type UpdateConnectionRequest struct {
	Name        string                 `json:"name"`
	Credentials map[string]interface{} `json:"credentials"`
	IsActive    *bool                  `json:"isActive"`
}

func (h *ConnectionHandler) UpdateConnection(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("userID")

	var req UpdateConnectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get existing connection
	conn, err := h.connectionService.GetConnection(c.Request.Context(), id, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Connection not found"})
		return
	}

	// Update fields
	if req.Name != "" {
		conn.ConnectionName = req.Name
	}
	if req.Credentials != nil {
		conn.Credentials = req.Credentials
	}
	if req.IsActive != nil {
		conn.IsActive = *req.IsActive
	}

	if err := h.connectionService.UpdateConnection(c.Request.Context(), conn); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, conn)
}

func (h *ConnectionHandler) DeleteConnection(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("userID")

	if err := h.connectionService.DeleteConnection(c.Request.Context(), id, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *ConnectionHandler) TestConnection(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("userID")

	conn, err := h.connectionService.GetConnection(c.Request.Context(), id, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Connection not found"})
		return
	}

	success, err := h.connectionService.TestConnection(c.Request.Context(), conn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": success})
}

func (h *ConnectionHandler) HandleIncomingWebhook(c *gin.Context) {
	provider := c.Param("provider")

	// Process webhook based on provider
	// This is a basic implementation
	c.JSON(http.StatusOK, gin.H{
		"status":   "received",
		"provider": provider,
	})
}
