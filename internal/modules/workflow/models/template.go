package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Template struct {
	ID        string    `gorm:"type:uuid;primary_key" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Category  string    `gorm:"not null" json:"category"`
	JSON      string    `gorm:"type:text;not null" json:"json"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (t *Template) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	return nil
}

func (Template) TableName() string {
	return "templates"
}
