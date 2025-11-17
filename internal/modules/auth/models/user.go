package models

import "time"

type User struct {
	ID           string         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	FullName     string         `gorm:"size:255;not null" json:"full_name"`
	Email        string         `gorm:"size:255;unique;not null" json:"email"`
	PasswordHash string         `gorm:"size:255;not null" json:"-"`
	Phone        string         `gorm:"size:50" json:"phone,omitempty"`
	City         string         `gorm:"size:100" json:"city,omitempty"`
	Role         string         `gorm:"size:20;default:'user'" json:"role"` // user, team, admin
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	Subscription        Subscription        `gorm:"foreignKey:UserID" json:"-"`
	NotificationSetting NotificationSetting `gorm:"foreignKey:UserID" json:"-"`
}
