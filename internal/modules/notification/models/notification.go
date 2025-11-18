package models

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type StringArray []string

func (s StringArray) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *StringArray) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, s)
}

type NotificationSettings struct {
	ID        string      `gorm:"type:uuid;primary_key" json:"id"`
	UserID    string      `gorm:"type:uuid;uniqueIndex;not null" json:"userId"`
	Email     bool        `gorm:"default:true" json:"email"`
	Slack     bool        `gorm:"default:false" json:"slack"`
	Channels  StringArray `gorm:"type:jsonb" json:"channels"` // errors, completions
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`

	// Relations
	User User `gorm:"foreignKey:UserID" json:"-"`
}

func (n *NotificationSettings) BeforeCreate(tx *gorm.DB) error {
	if n.ID == "" {
		n.ID = uuid.New().String()
	}
	return nil
}
