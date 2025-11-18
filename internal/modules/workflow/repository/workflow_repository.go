package repository

import (
	"gorm.io/gorm"
	"your-project/internal/models"
)

type WorkflowRepository struct {
	db *gorm.DB
}

func NewWorkflowRepository(db *gorm.DB) *WorkflowRepository {
	return &WorkflowRepository{db: db}
}

func (r *WorkflowRepository) Create(workflow *models.Workflow) error {
	return r.db.Create(workflow).Error
}

func (r *WorkflowRepository) FindByID(id string) (*models.Workflow, error) {
	var workflow models.Workflow
	err := r.db.Preload("User").First(&workflow, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &workflow, nil
}

func (r *WorkflowRepository) FindByUserID(userID string, active *bool, page, limit int) ([]models.Workflow, int64, error) {
	var workflows []models.Workflow
	var total int64

	query := r.db.Model(&models.Workflow{}).Where("user_id = ?", userID)

	if active != nil {
		query = query.Where("active = ?", *active)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&workflows).Error

	return workflows, total, err
}

func (r *WorkflowRepository) Update(workflow *models.Workflow) error {
	return r.db.Save(workflow).Error
}

func (r *WorkflowRepository) Delete(id string) error {
	return r.db.Delete(&models.Workflow{}, "id = ?", id).Error
}

func (r *WorkflowRepository) IncrementExecutionCount(id string, success bool) error {
	updates := map[string]interface{}{
		"total_executions": gorm.Expr("total_executions + 1"),
	}

	if success {
		updates["success_count"] = gorm.Expr("success_count + 1")
	} else {
		updates["error_count"] = gorm.Expr("error_count + 1")
	}

	return r.db.Model(&models.Workflow{}).Where("id = ?", id).Updates(updates).Error
}
