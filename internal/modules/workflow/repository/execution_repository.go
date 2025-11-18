package repository

import (
	"gorm.io/gorm"
	"your-project/internal/models"
)

type ExecutionRepository struct {
	db *gorm.DB
}

func NewExecutionRepository(db *gorm.DB) *ExecutionRepository {
	return &ExecutionRepository{db: db}
}

func (r *ExecutionRepository) Create(execution *models.Execution) error {
	return r.db.Create(execution).Error
}

func (r *ExecutionRepository) FindByID(id string) (*models.Execution, error) {
	var execution models.Execution
	err := r.db.Preload("Workflow").First(&execution, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &execution, nil
}

func (r *ExecutionRepository) FindByWorkflowID(workflowID string, status string, page, limit int) ([]models.Execution, int64, error) {
	var executions []models.Execution
	var total int64

	query := r.db.Model(&models.Execution{}).Where("workflow_id = ?", workflowID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&executions).Error

	return executions, total, err
}

func (r *ExecutionRepository) Update(execution *models.Execution) error {
	return r.db.Save(execution).Error
}

func (r *ExecutionRepository) Delete(id string) error {
	return r.db.Delete(&models.Execution{}, "id = ?", id).Error
}
