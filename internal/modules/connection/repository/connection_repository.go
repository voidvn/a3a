package repository

import (
	"gorm.io/gorm"
	"s4s-backend/internal/modules/connection/models"
)

type ConnectionRepository struct {
	db *gorm.DB
}

func NewConnectionRepository(db *gorm.DB) *ConnectionRepository {
	return &ConnectionRepository{db: db}
}

func (r *ConnectionRepository) Create(connection *models.Connection) error {
	return r.db.Create(connection).Error
}

func (r *ConnectionRepository) FindByID(id string) (*models.Connection, error) {
	var connection models.Connection
	err := r.db.First(&connection, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &connection, nil
}

func (r *ConnectionRepository) FindByUserID(userID string) ([]models.Connection, error) {
	var connections []models.Connection
	err := r.db.Where("user_id = ?", userID).Find(&connections).Error
	return connections, err
}

func (r *ConnectionRepository) Update(connection *models.Connection) error {
	return r.db.Save(connection).Error
}

func (r *ConnectionRepository) Delete(id string) error {
	return r.db.Delete(&models.Connection{}, "id = ?", id).Error
}
