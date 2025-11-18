package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Subscription struct {
	ID                   string     `gorm:"type:uuid;primary_key" json:"id"`
	UserID               string     `gorm:"type:uuid;uniqueIndex;not null" json:"userId"`
	Plan                 string     `gorm:"default:'freemium'" json:"plan"` // freemium, starter, team
	Status               string     `gorm:"default:'active'" json:"status"` // active, inactive, canceled
	WorkflowsLimit       int        `gorm:"default:5" json:"workflowsLimit"`
	ExecutionsLimit      int        `gorm:"default:100" json:"executionsLimit"`
	StripeSubscriptionID string     `json:"stripeSubscriptionId,omitempty"`
	StartedAt            time.Time  `json:"startedAt"`
	ExpiresAt            *time.Time `json:"expiresAt,omitempty"`
	CreatedAt            time.Time  `json:"createdAt"`
	UpdatedAt            time.Time  `json:"updatedAt"`

	// Relations
	User User `gorm:"foreignKey:UserID" json:"-"`
}

func (s *Subscription) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	return nil
}
