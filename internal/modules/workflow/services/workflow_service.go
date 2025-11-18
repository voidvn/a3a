package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	subscriptionRepo "s4s-backend/internal/modules/subscription/repository"
	"s4s-backend/internal/modules/workflow/dto"
	"s4s-backend/internal/modules/workflow/models"
	"s4s-backend/internal/modules/workflow/repository"
	"s4s-backend/internal/modules/workflow/services/engine"
)

type WorkflowService struct {
	workflowRepo     *repository.WorkflowRepository
	executionRepo    *repository.ExecutionRepository
	subscriptionRepo *subscriptionRepo.SubscriptionRepository
}

func NewWorkflowService(
	workflowRepo *repository.WorkflowRepository,
	executionRepo *repository.ExecutionRepository,
	subscriptionRepo *subscriptionRepo.SubscriptionRepository,
) *WorkflowService {
	return &WorkflowService{
		workflowRepo:     workflowRepo,
		executionRepo:    executionRepo,
		subscriptionRepo: subscriptionRepo,
	}
}

func (s *WorkflowService) CreateWorkflow(userID string, req *dto.CreateWorkflowRequest) (*models.Workflow, error) {
	// Check subscription limits
	subscription, err := s.subscriptionRepo.FindByUserID(userID)
	if err != nil {
		return nil, errors.New("subscription not found")
	}

	workflows, _, err := s.workflowRepo.FindByUserID(userID, nil, 1, 1000)
	if err != nil {
		return nil, err
	}

	if len(workflows) >= subscription.WorkflowsLimit {
		return nil, errors.New("workflow limit reached, please upgrade your plan")
	}

	workflow := &models.Workflow{
		UserID:     userID,
		Name:       req.Name,
		JSON:       req.JSON,
		Active:     false,
		MaxTimeout: req.MaxTimeout,
		RetryCount: req.RetryCount,
		RetryDelay: req.RetryDelay,
	}

	if workflow.MaxTimeout == 0 {
		workflow.MaxTimeout = 300
	}
	if workflow.RetryCount == 0 {
		workflow.RetryCount = 3
	}
	if workflow.RetryDelay == 0 {
		workflow.RetryDelay = 60
	}

	if err := s.workflowRepo.Create(workflow); err != nil {
		return nil, err
	}

	return workflow, nil
}

func (s *WorkflowService) GetWorkflow(workflowID string) (*models.Workflow, error) {
	return s.workflowRepo.FindByID(workflowID)
}

func (s *WorkflowService) ListWorkflows(userID string, active *bool, page, limit int) ([]models.Workflow, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	return s.workflowRepo.FindByUserID(userID, active, page, limit)
}

func (s *WorkflowService) UpdateWorkflow(workflowID string, req *dto.UpdateWorkflowRequest) (*models.Workflow, error) {
	workflow, err := s.workflowRepo.FindByID(workflowID)
	if err != nil {
		return nil, errors.New("workflow not found")
	}

	if req.Name != "" {
		workflow.Name = req.Name
	}
	if req.JSON != "" {
		workflow.JSON = req.JSON
	}
	if req.Active != nil {
		workflow.Active = *req.Active
	}
	if req.MaxTimeout > 0 {
		workflow.MaxTimeout = req.MaxTimeout
	}
	if req.RetryCount >= 0 {
		workflow.RetryCount = req.RetryCount
	}
	if req.RetryDelay > 0 {
		workflow.RetryDelay = req.RetryDelay
	}

	if err := s.workflowRepo.Update(workflow); err != nil {
		return nil, err
	}

	return workflow, nil
}

func (s *WorkflowService) DeleteWorkflow(workflowID string) error {
	return s.workflowRepo.Delete(workflowID)
}

func (s *WorkflowService) ExecuteWorkflow(ctx context.Context, workflowID string, isTest bool, testData map[string]interface{}) (string, error) {
	workflow, err := s.workflowRepo.FindByID(workflowID)
	if err != nil {
		return "", errors.New("workflow not found")
	}

	execution := &models.Execution{
		WorkflowID: workflowID,
		Status:     "pending",
		IsTest:     isTest,
	}

	if err := s.executionRepo.Create(execution); err != nil {
		return "", err
	}

	// Execute asynchronously
	go s.runWorkflow(context.Background(), workflow, execution, testData)

	return execution.ID, nil
}

func (s *WorkflowService) runWorkflow(ctx context.Context, workflow *models.Workflow, execution *models.Execution, testData map[string]interface{}) {
	now := time.Now()
	execution.StartedAt = &now
	execution.Status = "running"
	s.executionRepo.Update(execution)

	// Parse workflow JSON
	var workflowDef WorkflowDefinition
	if err := json.Unmarshal([]byte(workflow.JSON), &workflowDef); err != nil {
		s.failExecution(execution, fmt.Sprintf("Failed to parse workflow: %v", err))
		return
	}

	// Build execution graph
	nodeMap := make(map[string]*engine.Node)
	for i := range workflowDef.Nodes {
		nodeMap[workflowDef.Nodes[i].ID] = &workflowDef.Nodes[i]
	}

	adjacency := make(map[string][]string)
	for _, edge := range workflowDef.Edges {
		adjacency[edge.Source] = append(adjacency[edge.Source], edge.Target)
	}

	// Find trigger node
	var startNode *engine.Node
	for _, node := range workflowDef.Nodes {
		if node.Type == "trigger" {
			startNode = &node
			break
		}
	}

	if startNode == nil {
		s.failExecution(execution, "No trigger node found")
		return
	}

	// Initialize executors
	executors := map[string]engine.NodeExecutor{
		"http_request": &engine.HTTPRequestExecutor{},
		"email":        &engine.EmailExecutor{},
		"webhook":      &engine.WebhookExecutor{},
		"delay":        &engine.DelayExecutor{},
		"if":           &engine.IfExecutor{},
	}

	// Execute nodes
	logEntries := []string{}
	data := testData
	if data == nil {
		data = make(map[string]interface{})
	}

	if err := s.executeNode(ctx, startNode, nodeMap, adjacency, data, &logEntries, executors); err != nil {
		s.failExecution(execution, fmt.Sprintf("Execution failed: %v", err))
		execution.Log = s.formatLog(logEntries)
		s.executionRepo.Update(execution)
		s.workflowRepo.IncrementExecutionCount(workflow.ID, false)
		return
	}

	// Success
	endTime := time.Now()
	execution.Status = "success"
	execution.EndedAt = &endTime
	execution.DurationSeconds = int(endTime.Sub(*execution.StartedAt).Seconds())
	execution.Log = s.formatLog(logEntries)
	s.executionRepo.Update(execution)
	s.workflowRepo.IncrementExecutionCount(workflow.ID, true)
}

func (s *WorkflowService) executeNode(
	ctx context.Context,
	node *engine.Node,
	nodeMap map[string]*engine.Node,
	adjacency map[string][]string,
	data map[string]interface{},
	logEntries *[]string,
	executors map[string]engine.NodeExecutor,
) error {
	nodeType := node.Type
	if typeStr, ok := node.Data["type"].(string); ok {
		nodeType = typeStr
	}

	*logEntries = append(*logEntries, fmt.Sprintf("[%s] Node %s started", time.Now().Format("15:04:05"), node.ID))

	executor, exists := executors[nodeType]
	if !exists {
		return fmt.Errorf("unknown node type: %s", nodeType)
	}

	output, err := executor.Execute(ctx, node, data)
	if err != nil {
		*logEntries = append(*logEntries, fmt.Sprintf("[%s] Node %s failed: %v", time.Now().Format("15:04:05"), node.ID, err))
		return err
	}

	*logEntries = append(*logEntries, fmt.Sprintf("[%s] Node %s completed successfully", time.Now().Format("15:04:05"), node.ID))

	for k, v := range output {
		data[k] = v
	}

	nextNodeIDs := adjacency[node.ID]
	for _, nextID := range nextNodeIDs {
		nextNode := nodeMap[nextID]
		if err := s.executeNode(ctx, nextNode, nodeMap, adjacency, data, logEntries, executors); err != nil {
			return err
		}
	}

	return nil
}

func (s *WorkflowService) failExecution(execution *models.Execution, errorMsg string) {
	now := time.Now()
	execution.Status = "failed"
	execution.ErrorMessage = errorMsg
	execution.EndedAt = &now
	if execution.StartedAt != nil {
		execution.DurationSeconds = int(now.Sub(*execution.StartedAt).Seconds())
	}
	s.executionRepo.Update(execution)
}

func (s *WorkflowService) formatLog(entries []string) string {
	result := ""
	for _, entry := range entries {
		result += entry + "\n"
	}
	return result
}

type WorkflowDefinition struct {
	Nodes []engine.Node `json:"nodes"`
	Edges []Edge        `json:"edges"`
}

type Edge struct {
	ID     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
}
