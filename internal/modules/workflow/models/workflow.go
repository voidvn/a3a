package models

import "time"

type Workflow struct {
	ID            string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID        string    `gorm:"type:uuid;not null;index" json:"user_id"`
	Name          string    `gorm:"size:255;not null" json:"name"`
	JsonGraph     string    `gorm:"type:text;not null" json:"json_graph"` // React Flow JSON
	Active        bool      `gorm:"default:true" json:"active"`
	MaxTimeoutSec int       `gorm:"default:300" json:"max_timeout_sec"`
	RetryCount    int       `gorm:"default:3" json:"retry_count"`
	RetryDelaySec int       `gorm:"default:60" json:"retry_delay_sec"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Executions []Execution `gorm:"foreignKey:WorkflowID" json:"-"`
}
