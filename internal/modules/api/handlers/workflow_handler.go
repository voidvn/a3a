package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"your-project/internal/dto"
	"your-project/internal/engine"
	"your-project/internal/service"
)

type WorkflowHandler struct {
	workflowService *service.WorkflowService
	workflowEngine  *engine.WorkflowEngine
}

func NewWorkflowHandler(workflowService *service.WorkflowService, workflowEngine *engine.WorkflowEngine) *WorkflowHandler {
	return &WorkflowHandler{
		workflowService: workflowService,
		workflowEngine:  workflowEngine,
	}
}

// ListWorkflows godoc
// @Summary List workflows
// @Tags workflows
// @Security BearerAuth
// @Produce json
// @Param active query bool false "Filter by active status"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {array} models.Workflow
// @Router /workflows [get]
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

// CreateWorkflow godoc
// @Summary Create workflow
// @Tags workflows
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateWorkflowRequest true "Workflow data"
// @Success 201 {object} models.Workflow
// @Router /workflows [post]
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

// GetWorkflow godoc
// @Summary Get workflow
// @Tags workflows
// @Security BearerAuth
// @Produce json
// @Param id path string true "Workflow ID"
// @Success 200 {object} models.Workflow
// @Router /workflows/{id} [get]
func (h *WorkflowHandler) GetWorkflow(c *gin.Context) {
	id := c.Param("id")

	workflow, err := h.workflowService.GetWorkflow(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Workflow not found", "code": 404})
		return
	}

	c.JSON(http.StatusOK, workflow)
}

// UpdateWorkflow godoc
// @Summary Update workflow
// @Tags workflows
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Workflow ID"
// @Param request body dto.UpdateWorkflowRequest true "Update data"
// @Success 200 {object} models.Workflow
// @Router /workflows/{id} [put]
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

// DeleteWorkflow godoc
// @Summary Delete workflow
// @Tags workflows
// @Security BearerAuth
// @Param id path string true "Workflow ID"
// @Success 204
// @Router /workflows/{id} [delete]
func (h *WorkflowHandler) DeleteWorkflow(c *gin.Context) {
	id := c.Param("id")

	err := h.workflowService.DeleteWorkflow(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "code": 400})
		return
	}

	c.Status(http.StatusNoContent)
}

// TestWorkflow godoc
// @Summary Test workflow
// @Tags workflows
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Workflow ID"
// @Param request body dto.TestWorkflowRequest false "Test data"
// @Success 200 {object} models.Execution
// @Router /workflows/{id}/test [post]
func (h *WorkflowHandler) TestWorkflow(c *gin.Context) {
	id := c.Param("id")

	var req dto.TestWorkflowRequest
	c.ShouldBindJSON(&req)

	executionID, err := h.workflowEngine.ExecuteWorkflow(c.Request.Context(), id, true, req.TestData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "code": 400})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"executionId": executionID,
		"message":     "Test execution started",
	})
}

// RunWorkflow godoc
// @Summary Run workflow
// @Tags workflows
// @Security BearerAuth
// @Produce json
// @Param id path string true "Workflow ID"
// @Success 202 {object} map[string]interface{}
// @Router /workflows/{id}/run [post]
func (h *WorkflowHandler) RunWorkflow(c *gin.Context) {
	id := c.Param("id")

	executionID, err := h.workflowEngine.ExecuteWorkflow(c.Request.Context(), id, false, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "code": 400})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"executionId": executionID,
		"message":     "Workflow queued for execution",
	})
}
