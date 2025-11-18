package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"s4s-backend/internal/modules/workflow/services"
)

type ExecutionHandler struct {
	executionService *services.ExecutionService
}

func NewExecutionHandler(executionService *services.ExecutionService) *ExecutionHandler {
	return &ExecutionHandler{executionService: executionService}
}

func (h *ExecutionHandler) ListExecutions(c *gin.Context) {
	workflowID := c.Query("workflowId")
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	executions, total, err := h.executionService.ListExecutions(workflowID, status, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "code": 500})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  executions,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *ExecutionHandler) GetExecution(c *gin.Context) {
	id := c.Param("id")

	execution, err := h.executionService.GetExecution(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Execution not found", "code": 404})
		return
	}

	c.JSON(http.StatusOK, execution)
}

func (h *ExecutionHandler) DeleteExecution(c *gin.Context) {
	id := c.Param("id")

	err := h.executionService.DeleteExecution(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "code": 400})
		return
	}

	c.Status(http.StatusNoContent)
}
