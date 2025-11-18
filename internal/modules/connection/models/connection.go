package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Connection struct {
	ID             string         `gorm:"type:uuid;primary_key" json:"id"`
	UserID         string         `gorm:"type:uuid;not null" json:"userId"`
	ConnectionName string         `gorm:"not null" json:"name"`
	ServiceName    string         `gorm:"not null" json:"service"`
	Credentials    JSONB          `gorm:"type:jsonb;not null" json:"credentials"`
	IsActive       bool           `gorm:"default:true" json:"isActive"`
	LastTestedAt   *time.Time     `json:"lastTestedAt,omitempty"`
	LastTestStatus *string        `gorm:"type:varchar(20)" json:"lastTestStatus,omitempty"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &j)
}

func (c *Connection) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return nil
}

func (Connection) TableName() string {
	return "connections"
}
