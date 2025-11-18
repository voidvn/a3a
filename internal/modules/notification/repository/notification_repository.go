package repository

import (
	"gorm.io/gorm"
	"s4s-backend/internal/modules/notification/models"
)

type NotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(settings *models.NotificationSettings) error {
	return r.db.Create(settings).Error
}

func (r *NotificationRepository) FindByUserID(userID string) (*models.NotificationSettings, error) {
	var settings models.NotificationSettings
	err := r.db.First(&settings, "user_id = ?", userID).Error
	if err != nil {
		return nil, err
	}
	return &settings, nil
}

func (r *NotificationRepository) Update(settings *models.NotificationSettings) error {
	return r.db.Save(settings).Error
}
