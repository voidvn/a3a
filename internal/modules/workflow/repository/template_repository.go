package repository

import (
	"s4s-backend/internal/modules/workflow/models"

	"gorm.io/gorm"
)

type TemplateRepository struct {
	db *gorm.DB
}

func NewTemplateRepository(db *gorm.DB) *TemplateRepository {
	return &TemplateRepository{db: db}
}

func (r *TemplateRepository) Create(template *models.Template) error {
	return r.db.Create(template).Error
}

func (r *TemplateRepository) FindByID(id string) (*models.Template, error) {
	var template models.Template
	err := r.db.First(&template, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

func (r *TemplateRepository) FindByCategory(category string) ([]models.Template, error) {
	var templates []models.Template
	query := r.db.Model(&models.Template{})

	if category != "" {
		query = query.Where("category = ?", category)
	}

	err := query.Find(&templates).Error
	return templates, err
}

func (r *TemplateRepository) Update(template *models.Template) error {
	return r.db.Save(template).Error
}

func (r *TemplateRepository) Delete(id string) error {
	return r.db.Delete(&models.Template{}, "id = ?", id).Error
}
