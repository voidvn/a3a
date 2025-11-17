package models

import "time"

type Connection struct {
	ID                   string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID               string    `gorm:"type:uuid;not null;index" json:"user_id"`
	Service              string    `gorm:"size:100;not null" json:"service"` // slack, gmail, crm, etc.
	EncryptedCredentials string    `gorm:"type:text;not null" json:"-"`      // AES-256
	CreatedAt            time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt            time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
