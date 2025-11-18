package repository

import (
	"errors"

	"gorm.io/gorm"
	"s4s-backend/internal/modules/connection/models"
)

type ConnectionRepository interface {
	Create(conn *models.Connection) error
	GetByID(id, userID string) (*models.Connection, error)
	FindByUserID(userID string) ([]*models.Connection, error)
	Update(conn *models.Connection) error
	Delete(id, userID string) error
}

type connectionRepository struct {
	db *gorm.DB
}

func NewConnectionRepository(db *gorm.DB) ConnectionRepository {
	return &connectionRepository{db: db}
}

func (r *connectionRepository) Create(conn *models.Connection) error {
	return r.db.Create(conn).Error
}

func (r *connectionRepository) GetByID(id, userID string) (*models.Connection, error) {
	var conn models.Connection
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&conn).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("connection not found")
		}
		return nil, err
	}
	return &conn, nil
}

func (r *connectionRepository) FindByUserID(userID string) ([]*models.Connection, error) {
	var connections []*models.Connection
	err := r.db.Where("user_id = ?", userID).Find(&connections).Error
	if err != nil {
		return nil, err
	}
	return connections, nil
}

func (r *connectionRepository) Update(conn *models.Connection) error {
	return r.db.Save(conn).Error
}

func (r *connectionRepository) Delete(id, userID string) error {
	result := r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Connection{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("connection not found")
	}
	return nil
}
