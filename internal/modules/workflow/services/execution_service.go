package service

import (
	"errors"
	"fmt"
	"time"
	"your-project/internal/models"
	"your-project/internal/repository"
)

type ExecutionService struct {
	executionRepo    *repository.ExecutionRepository
	subscriptionRepo *repository.SubscriptionRepository
}

func NewExecutionService(
	executionRepo *repository.ExecutionRepository,
	subscriptionRepo *repository.SubscriptionRepository,
) *ExecutionService {
	return &ExecutionService{
		executionRepo:    executionRepo,
		subscriptionRepo: subscriptionRepo,
	}
}

func (s *ExecutionService) GetExecution(executionID string) (*models.Execution, error) {
	return s.executionRepo.FindByID(executionID)
}

func (s *ExecutionService) ListExecutions(workflowID, status string, page, limit int) ([]models.Execution, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	return s.executionRepo.FindByWorkflowID(workflowID, status, page, limit)
}

func (s *ExecutionService) DeleteExecution(executionID string) error {
	return s.executionRepo.Delete(executionID)
}

// ============================================
// internal/engine/workflow_engine.go
// ============================================
package engine

import (
"context"
"encoding/json"
"errors"
"fmt"
"time"
"your-project/internal/models"
"your-project/internal/repository"
)

// WorkflowDefinition represents the JSON structure of a workflow
type WorkflowDefinition struct {
	Nodes []Node       `json:"nodes"`
	Edges []Edge       `json:"edges"`
}

type Node struct {
	ID       string                 `json:"id"`
	Type     string                 `json:"type"` // trigger, action, logic, utility
	Data     map[string]interface{} `json:"data"`
	Position Position               `json:"position"`
}

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Edge struct {
	ID     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
}

type WorkflowEngine struct {
	executionRepo *repository.ExecutionRepository
	workflowRepo  *repository.WorkflowRepository
	nodeExecutors map[string]NodeExecutor
}

type NodeExecutor interface {
	Execute(ctx context.Context, node *Node, input map[string]interface{}) (map[string]interface{}, error)
}

func NewWorkflowEngine(
	executionRepo *repository.ExecutionRepository,
	workflowRepo *repository.WorkflowRepository,
) *WorkflowEngine {
	engine := &WorkflowEngine{
		executionRepo: executionRepo,
		workflowRepo:  workflowRepo,
		nodeExecutors: make(map[string]NodeExecutor),
	}

	// Register node executors
	engine.RegisterExecutor("http_request", &HTTPRequestExecutor{})
	engine.RegisterExecutor("email", &EmailExecutor{})
	engine.RegisterExecutor("webhook", &WebhookExecutor{})
	engine.RegisterExecutor("delay", &DelayExecutor{})
	engine.RegisterExecutor("if", &IfExecutor{})

	return engine
}

func (e *WorkflowEngine) RegisterExecutor(nodeType string, executor NodeExecutor) {
	e.nodeExecutors[nodeType] = executor
}

func (e *WorkflowEngine) ExecuteWorkflow(ctx context.Context, workflowID string, isTest bool, testData map[string]interface{}) (string, error) {
	// Get workflow
	workflow, err := e.workflowRepo.FindByID(workflowID)
	if err != nil {
		return "", errors.New("workflow not found")
	}

	// Create execution record
	execution := &models.Execution{
		WorkflowID: workflowID,
		Status:     "pending",
		IsTest:     isTest,
	}

	if err := e.executionRepo.Create(execution); err != nil {
		return "", err
	}

	// Execute asynchronously
	go func() {
		e.runWorkflow(context.Background(), workflow, execution, testData)
	}()

	return execution.ID, nil
}

func (e *WorkflowEngine) runWorkflow(ctx context.Context, workflow *models.Workflow, execution *models.Execution, testData map[string]interface{}) {
	now := time.Now()
	execution.StartedAt = &now
	execution.Status = "running"
	e.executionRepo.Update(execution)

	// Parse workflow JSON
	var workflowDef WorkflowDefinition
	if err := json.Unmarshal([]byte(workflow.JSON), &workflowDef); err != nil {
		e.failExecution(execution, fmt.Sprintf("Failed to parse workflow: %v", err))
		return
	}

	// Build execution graph
	nodeMap := make(map[string]*Node)
	for i := range workflowDef.Nodes {
		nodeMap[workflowDef.Nodes[i].ID] = &workflowDef.Nodes[i]
	}

	// Build adjacency list
	adjacency := make(map[string][]string)
	for _, edge := range workflowDef.Edges {
		adjacency[edge.Source] = append(adjacency[edge.Source], edge.Target)
	}

	// Find trigger node (start node)
	var startNode *Node
	for _, node := range workflowDef.Nodes {
		if node.Type == "trigger" {
			startNode = &node
			break
		}
	}

	if startNode == nil {
		e.failExecution(execution, "No trigger node found")
		return
	}

	// Execute nodes
	logEntries := []string{}
	data := testData
	if data == nil {
		data = make(map[string]interface{})
	}

	// Recursive execution
	if err := e.executeNode(ctx, startNode, nodeMap, adjacency, data, &logEntries); err != nil {
		e.failExecution(execution, fmt.Sprintf("Execution failed: %v", err))
		execution.Log = e.formatLog(logEntries)
		e.executionRepo.Update(execution)
		e.workflowRepo.IncrementExecutionCount(workflow.ID, false)
		return
	}

	// Success
	endTime := time.Now()
	execution.Status = "success"
	execution.EndedAt = &endTime
	execution.DurationSeconds = int(endTime.Sub(*execution.StartedAt).Seconds())
	execution.Log = e.formatLog(logEntries)
	e.executionRepo.Update(execution)
	e.workflowRepo.IncrementExecutionCount(workflow.ID, true)
}

func (e *WorkflowEngine) executeNode(
	ctx context.Context,
	node *Node,
	nodeMap map[string]*Node,
	adjacency map[string][]string,
	data map[string]interface{},
	logEntries *[]string,
) error {
	// Get node type from data if exists, otherwise use type field
	nodeType := node.Type
	if typeStr, ok := node.Data["type"].(string); ok {
		nodeType = typeStr
	}

	// Log node start
	*logEntries = append(*logEntries, fmt.Sprintf("[%s] Node %s started", time.Now().Format("15:04:05"), node.ID))

	// Get executor
	executor, exists := e.nodeExecutors[nodeType]
	if !exists {
		return fmt.Errorf("unknown node type: %s", nodeType)
	}

	// Execute node
	output, err := executor.Execute(ctx, node, data)
	if err != nil {
		*logEntries = append(*logEntries, fmt.Sprintf("[%s] Node %s failed: %v", time.Now().Format("15:04:05"), node.ID, err))
		return err
	}

	*logEntries = append(*logEntries, fmt.Sprintf("[%s] Node %s completed successfully", time.Now().Format("15:04:05"), node.ID))

	// Merge output to data for next nodes
	for k, v := range output {
		data[k] = v
	}

	// Execute next nodes
	nextNodeIDs := adjacency[node.ID]
	for _, nextID := range nextNodeIDs {
		nextNode := nodeMap[nextID]
		if err := e.executeNode(ctx, nextNode, nodeMap, adjacency, data, logEntries); err != nil {
			return err
		}
	}

	return nil
}

func (e *WorkflowEngine) failExecution(execution *models.Execution, errorMsg string) {
	now := time.Now()
	execution.Status = "failed"
	execution.ErrorMessage = errorMsg
	execution.EndedAt = &now
	if execution.StartedAt != nil {
		execution.DurationSeconds = int(now.Sub(*execution.StartedAt).Seconds())
	}
	e.executionRepo.Update(execution)
}

func (e *WorkflowEngine) formatLog(entries []string) string {
	result := ""
	for _, entry := range entries {
		result += entry + "\n"
	}
	return result
}
