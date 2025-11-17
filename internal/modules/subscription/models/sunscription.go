package models

import "time"

type Subscription struct {
	ID               string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID           string     `gorm:"type:uuid;unique;not null" json:"user_id"`
	Plan             string     `gorm:"size:50;default:'freemium'" json:"plan"`
	Status           string     `gorm:"size:50;default:'active'" json:"status"`
	CurrentPeriodEnd *time.Time `json:"current_period_end,omitempty"`
	UsageLimits      string     `gorm:"type:jsonb" json:"usage_limits,omitempty"`
	CreatedAt        time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}
