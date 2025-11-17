package models

import "time"

type Execution struct {
	ID         string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	WorkflowID string     `gorm:"type:uuid;not null;index" json:"workflow_id"`
	UserID     string     `gorm:"type:uuid;not null;index" json:"user_id"`
	Status     string     `gorm:"size:20;default:'pending'" json:"status"` // pending, running, success, failed
	IsTest     bool       `gorm:"default:false" json:"is_test"`
	Log        string     `gorm:"type:text" json:"log,omitempty"`
	StartedAt  time.Time  `json:"started_at,omitempty"`
	EndedAt    *time.Time `json:"ended_at,omitempty"`
}
