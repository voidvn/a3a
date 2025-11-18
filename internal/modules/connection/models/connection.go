package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Connection struct {
	ID             string    `gorm:"type:uuid;primary_key" json:"id"`
	UserID         string    `gorm:"type:uuid;not null" json:"userId"`
	ServiceName    string    `gorm:"not null" json:"service"`
	ConnectionName string    `json:"connectionName"`
	Credentials    string    `gorm:"type:text;not null" json:"-"` // Encrypted JSON
	IsActive       bool      `gorm:"default:true" json:"isActive"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`

	// Relations
	User User `gorm:"foreignKey:UserID" json:"-"`
}

func (c *Connection) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return nil
}
