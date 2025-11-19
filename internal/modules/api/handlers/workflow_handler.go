package handlers

import (
	"net/http"
	"strconv"
	"time"

	"s4s-backend/internal/modules/workflow/dto"
	"s4s-backend/internal/modules/workflow/services"

	"github.com/gin-gonic/gin"
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

	// Define a custom response struct that will handle the JSON field as a string
	type workflowResponse struct {
		ID              string    `json:"id"`
		UserID          string    `json:"userId"`
		Name            string    `json:"name"`
		JSON            string    `json:"json"` // Changed to string
		Active          bool      `json:"active"`
		MaxTimeout      int       `json:"maxTimeout"`
		RetryCount      int       `json:"retryCount"`
		RetryDelay      int       `json:"retryDelay"`
		TriggerType     string    `json:"triggerType"`
		TotalExecutions int       `json:"totalExecutions"`
		SuccessCount    int       `json:"successCount"`
		ErrorCount      int       `json:"errorCount"`
		CreatedAt       time.Time `json:"createdAt"`
		UpdatedAt       time.Time `json:"updatedAt"`
	}

	// The JSON string as is, without unmarshaling
	jsonStr := `{"edges":[{"id":"xy-edge__2535ztb-mu4umeh","source":"2535ztb","target":"mu4umeh"},{"id":"xy-edge__mu4umeh-0djt8do","source":"mu4umeh","target":"0djt8do"},{"id":"xy-edge__0djt8do-yct4tm0","source":"0djt8do","target":"yct4tm0"}],"nodes":[{"id":"2535ztb","data":{"data":{"type":"webhook","config":{}},"type":"trigger","label":"Trigger Node"},"type":"input","width":75,"height":75,"dragging":false,"measured":{"width":75,"height":75},"position":{"x":122.5,"y":265},"selected":false,"deletable":true},{"id":"mu4umeh","data":{"data":{"type":"http_request","config":{"path":"/api/user","method":"GET","authentication":"JWT Auth"}},"type":"action","label":"Request To User API"},"type":"default","width":75,"height":75,"dragging":false,"measured":{"width":75,"height":75},"position":{"x":117.25034365844411,"y":438.97360127951174},"selected":false,"deletable":true},{"id":"0djt8do","data":{"data":{"type":"delay","config":{}},"type":"action","label":"Wait 2000 ms"},"type":"default","width":75,"height":75,"dragging":false,"measured":{"width":75,"height":75},"position":{"x":215.83062924779745,"y":580.3484294502896},"selected":false,"deletable":true},{"id":"yct4tm0","data":{"data":{"type":"slack","config":{}},"type":"action","label":"Notify Sale Managers"},"type":"default","width":75,"height":75,"dragging":false,"measured":{"width":75,"height":75},"position":{"x":445.0871073625724,"y":702.618551111503},"selected":true,"deletable":true}]}`

	response := workflowResponse{
		ID:              id,
		Name:            "Hardcoded Workflow",
		JSON:            jsonStr, // Directly use the JSON string
		Active:          true,
		MaxTimeout:      300,
		RetryCount:      3,
		RetryDelay:      60,
		TriggerType:     "webhook",
		TotalExecutions: 0,
		SuccessCount:    0,
		ErrorCount:      0,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
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
