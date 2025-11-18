package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"your-project/internal/service"
)

type ExecutionHandler struct {
	executionService *service.ExecutionService
}

func NewExecutionHandler(executionService *service.ExecutionService) *ExecutionHandler {
	return &ExecutionHandler{executionService: executionService}
}

// ListExecutions godoc
// @Summary List executions
// @Tags executions
// @Security BearerAuth
// @Produce json
// @Param workflowId query string false "Filter by workflow ID"
// @Param status query string false "Filter by status"
// @Success 200 {array} models.Execution
// @Router /executions [get]
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

// GetExecution godoc
// @Summary Get execution
// @Tags executions
// @Security BearerAuth
// @Produce json
// @Param id path string true "Execution ID"
// @Success 200 {object} models.Execution
// @Router /executions/{id} [get]
func (h *ExecutionHandler) GetExecution(c *gin.Context) {
	id := c.Param("id")

	execution, err := h.executionService.GetExecution(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Execution not found", "code": 404})
		return
	}

	c.JSON(http.StatusOK, execution)
}

// DeleteExecution godoc
// @Summary Delete execution
// @Tags executions
// @Security BearerAuth
// @Param id path string true "Execution ID"
// @Success 204
// @Router /executions/{id} [delete]
func (h *ExecutionHandler) DeleteExecution(c *gin.Context) {
	id := c.Param("id")

	err := h.executionService.DeleteExecution(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "code": 400})
		return
	}

	c.Status(http.StatusNoContent)
}
