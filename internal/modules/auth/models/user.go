package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID              string    `gorm:"type:uuid;primary_key" json:"id"`
	FullName        string    `gorm:"not null" json:"fullName"`
	Email           string    `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash    string    `gorm:"not null" json:"-"`
	Phone           string    `json:"phone"`
	City            string    `json:"city"`
	Role            string    `gorm:"default:'user'" json:"role"` // user, team, admin
	IsActive        bool      `gorm:"default:true" json:"isActive"`
	EmailVerifiedAt bool      `gorm:"default:false" json:"emailVerifiedAt"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}

func (User) TableName() string {
	return "users"
}
