package services

import (
	"s4s-backend/internal/modules/workflow/models"
	"s4s-backend/internal/modules/workflow/repository"
)

type ExecutionService struct {
	executionRepo *repository.ExecutionRepository
}

func NewExecutionService(executionRepo *repository.ExecutionRepository) *ExecutionService {
	return &ExecutionService{executionRepo: executionRepo}
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
