package models

import "time"

type NotificationSetting struct {
	UserID       string    `gorm:"type:uuid;primaryKey" json:"user_id"`
	EmailEnabled bool      `gorm:"default:true" json:"email_enabled"`
	SlackEnabled bool      `gorm:"default:false" json:"slack_enabled"`
	Channels     string    `gorm:"type:jsonb" json:"channels,omitempty"` // ["#leads", "#errors"]
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
