package repository

import (
	"gorm.io/gorm"
	"s4s-backend/internal/modules/subscription/models"
)

type SubscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) Create(subscription *models.Subscription) error {
	return r.db.Create(subscription).Error
}

func (r *SubscriptionRepository) FindByUserID(userID string) (*models.Subscription, error) {
	var subscription models.Subscription
	err := r.db.First(&subscription, "user_id = ?", userID).Error
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}

func (r *SubscriptionRepository) Update(subscription *models.Subscription) error {
	return r.db.Save(subscription).Error
}
