package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Execution struct {
	ID              string     `gorm:"type:uuid;primary_key" json:"id"`
	WorkflowID      string     `gorm:"type:uuid;not null" json:"workflowId"`
	Status          string     `gorm:"default:'pending'" json:"status"`
	Log             string     `gorm:"type:text" json:"log"`
	IsTest          bool       `gorm:"default:false" json:"isTest"`
	ErrorMessage    string     `gorm:"type:text" json:"errorMessage,omitempty"`
	DurationSeconds int        `json:"durationSeconds"`
	StartedAt       *time.Time `json:"startedAt,omitempty"`
	EndedAt         *time.Time `json:"endedAt,omitempty"`
	CreatedAt       time.Time  `json:"createdAt"`
}

func (e *Execution) BeforeCreate(tx *gorm.DB) error {
	if e.ID == "" {
		e.ID = uuid.New().String()
	}
	return nil
}

func (Execution) TableName() string {
	return "executions"
}
