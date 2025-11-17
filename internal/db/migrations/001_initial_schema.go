package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var InitialSchema = []*gormigrate.Migration{
	{
		ID: "20251117_001_initial_schema",
		Migrate: func(db *gorm.DB) error {
			type User struct {
				ID           string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
				FullName     string
				Email        string `gorm:"uniqueIndex;size:255"`
				PasswordHash string `gorm:"size:255"`
				Phone        string `gorm:"size:50"`
				City         string `gorm:"size:100"`
				Role         string `gorm:"size:20;default:'user'"`
				CreatedAt    time.Time
				UpdatedAt    time.Time
				DeletedAt    gorm.DeletedAt `gorm:"index"`
			}

			type Workflow struct {
				ID            string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
				UserID        string `gorm:"type:uuid;index"`
				Name          string `gorm:"size:255;not null"`
				JsonGraph     string `gorm:"type:text;not null"`
				Active        bool   `gorm:"default:true"`
				MaxTimeoutSec int    `gorm:"default:300"`
				RetryCount    int    `gorm:"default:3"`
				RetryDelaySec int    `gorm:"default:60"`
				CreatedAt     time.Time
				UpdatedAt     time.Time
			}

			type Execution struct {
				ID         string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
				WorkflowID string `gorm:"type:uuid;index"`
				UserID     string `gorm:"type:uuid;index"`
				Status     string `gorm:"size:20;default:'pending'"`
				IsTest     bool   `gorm:"default:false"`
				Log        string `gorm:"type:text"`
				StartedAt  time.Time
				EndedAt    *time.Time
			}

			type Connection struct {
				ID                   string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
				UserID               string `gorm:"type:uuid;index"`
				Service              string `gorm:"size:100;not null"`
				EncryptedCredentials string `gorm:"type:text;not null"`
				CreatedAt            time.Time
				UpdatedAt            time.Time
			}

			type Subscription struct {
				ID               string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
				UserID           string `gorm:"type:uuid;uniqueIndex"`
				Plan             string `gorm:"size:50;default:'freemium'"`
				Status           string `gorm:"size:50;default:'active'"`
				CurrentPeriodEnd *time.Time
				UsageLimits      string `gorm:"type:jsonb"`
				CreatedAt        time.Time
				UpdatedAt        time.Time
			}

			type NotificationSetting struct {
				UserID       string `gorm:"type:uuid;primaryKey"`
				EmailEnabled bool   `gorm:"default:true"`
				SlackEnabled bool   `gorm:"default:false"`
				Channels     string `gorm:"type:jsonb"`
				UpdatedAt    time.Time
			}

			return db.Transaction(func(tx *gorm.DB) error {
				return tx.AutoMigrate(
					&User{},
					&Workflow{},
					&Execution{},
					&Connection{},
					&Subscription{},
					&NotificationSetting{},
				)
			})
		},
		Rollback: func(db *gorm.DB) error {
			return db.Migrator().DropTable(
				"executions",
				"workflows",
				"connections",
				"subscriptions",
				"notification_settings",
				"users",
			)
		},
	},
}
