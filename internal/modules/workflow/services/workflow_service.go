package service

import (
	"errors"
	"your-project/internal/dto"
	"your-project/internal/models"
	"your-project/internal/repository"
)

type WorkflowService struct {
	workflowRepo     *repository.WorkflowRepository
	subscriptionRepo *repository.SubscriptionRepository
}

func NewWorkflowService(
	workflowRepo *repository.WorkflowRepository,
	subscriptionRepo *repository.SubscriptionRepository,
) *WorkflowService {
	return &WorkflowService{
		workflowRepo:     workflowRepo,
		subscriptionRepo: subscriptionRepo,
	}
}

func (s *WorkflowService) CreateWorkflow(userID string, req *dto.CreateWorkflowRequest) (*models.Workflow, error) {
	// Check subscription limits
	subscription, err := s.subscriptionRepo.FindByUserID(userID)
	if err != nil {
		return nil, errors.New("subscription not found")
	}

	// Count existing workflows
	workflows, _, err := s.workflowRepo.FindByUserID(userID, nil, 1, 1000)
	if err != nil {
		return nil, err
	}

	if len(workflows) >= subscription.WorkflowsLimit {
		return nil, errors.New("workflow limit reached, please upgrade your plan")
	}

	// Create workflow
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
