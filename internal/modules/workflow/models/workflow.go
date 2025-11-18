package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Workflow struct {
	ID              string    `gorm:"type:uuid;primary_key" json:"id"`
	UserID          string    `gorm:"type:uuid;not null" json:"userId"`
	Name            string    `gorm:"not null" json:"name"`
	JSON            string    `gorm:"type:text;not null" json:"json"`
	Active          bool      `gorm:"default:false" json:"active"`
	MaxTimeout      int       `gorm:"default:300" json:"maxTimeout"`
	RetryCount      int       `gorm:"default:3" json:"retryCount"`
	RetryDelay      int       `gorm:"default:60" json:"retryDelay"`
	TriggerType     string    `json:"triggerType"`
	TotalExecutions int       `gorm:"default:0" json:"totalExecutions"`
	SuccessCount    int       `gorm:"default:0" json:"successCount"`
	ErrorCount      int       `gorm:"default:0" json:"errorCount"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

func (w *Workflow) BeforeCreate(tx *gorm.DB) error {
	if w.ID == "" {
		w.ID = uuid.New().String()
	}
	return nil
}

func (Workflow) TableName() string {
	return "workflows"
}
