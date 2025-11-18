package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"s4s-backend/internal/modules/workflow/dto"
	"s4s-backend/internal/modules/workflow/services"
)

type WorkflowHandler struct {
	workflowService *services.WorkflowService
}

func NewWorkflowHandler(workflowService *services.WorkflowService) *WorkflowHandler {
	return &WorkflowHandler{workflowService: workflowService}
}

func (h *WorkflowHandler) ListWorkflows(c *gin.Context) {
	userID := c.GetString("userID")

	var active *bool
	if activeStr := c.Query("active"); activeStr != "" {
		activeBool := activeStr == "true"
		active = &activeBool
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	workflows, total, err := h.workflowService.ListWorkflows(userID, active, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "code": 500})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  workflows,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *WorkflowHandler) CreateWorkflow(c *gin.Context) {
	userID := c.GetString("userID")

	var req dto.CreateWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "code": 400})
		return
	}

	workflow, err := h.workflowService.CreateWorkflow(userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "code": 400})
		return
	}

	c.JSON(http.StatusCreated, workflow)
}

func (h *WorkflowHandler) GetWorkflow(c *gin.Context) {
	id := c.Param("id")

	workflow, err := h.workflowService.GetWorkflow(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Workflow not found", "code": 404})
		return
	}

	c.JSON(http.StatusOK, workflow)
}

func (h *WorkflowHandler) UpdateWorkflow(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdateWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "code": 400})
		return
	}

	workflow, err := h.workflowService.UpdateWorkflow(id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "code": 400})
		return
	}

	c.JSON(http.StatusOK, workflow)
}

func (h *WorkflowHandler) DeleteWorkflow(c *gin.Context) {
	id := c.Param("id")

	err := h.workflowService.DeleteWorkflow(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "code": 400})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *WorkflowHandler) TestWorkflow(c *gin.Context) {
	id := c.Param("id")

	var req dto.TestWorkflowRequest
	c.ShouldBindJSON(&req)

	executionID, err := h.workflowService.ExecuteWorkflow(c.Request.Context(), id, true, req.TestData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "code": 400})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"executionId": executionID,
		"message":     "Test execution started",
	})
}

func (h *WorkflowHandler) RunWorkflow(c *gin.Context) {
	id := c.Param("id")

	executionID, err := h.workflowService.ExecuteWorkflow(c.Request.Context(), id, false, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "code": 400})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"executionId": executionID,
		"message":     "Workflow queued for execution",
	})
}
