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
				ID              string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
				FullName        string `gorm:"size:255;not null"`
				Email           string `gorm:"uniqueIndex;size:255;not null"`
				PasswordHash    string `gorm:"size:255;not null"`
				Phone           string `gorm:"size:50"`
				City            string `gorm:"size:100"`
				Role            string `gorm:"size:20;default:'user';not null"`
				IsActive        bool   `gorm:"default:true;not null"`
				EmailVerifiedAt *time.Time
				LastLoginAt     *time.Time
				CreatedAt       time.Time
				UpdatedAt       time.Time
				DeletedAt       gorm.DeletedAt `gorm:"index"`
			}

			type Workflow struct {
				ID            string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
				UserID        string `gorm:"type:uuid;index;not null"`
				Name          string `gorm:"size:255;not null"`
				Description   string `gorm:"type:text"`
				GraphConfig   string `gorm:"type:jsonb;not null"`
				Status        string `gorm:"size:20;default:'draft';not null"`
				MaxTimeoutSec int    `gorm:"default:300;not null"`
				RetryCount    int    `gorm:"default:3;not null"`
				RetryDelaySec int    `gorm:"default:60;not null"`
				Version       int    `gorm:"default:1;not null"`
				IsPublic      bool   `gorm:"default:false;not null"`
				Tags          string `gorm:"type:jsonb;default:'[]'"`
				CreatedAt     time.Time
				UpdatedAt     time.Time
			}

			type Execution struct {
				ID         string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
				WorkflowID string `gorm:"type:uuid;index;not null"`
				UserID     string `gorm:"type:uuid;index;not null"`
				Status     string `gorm:"size:20;default:'pending';not null;index"`
				IsTest     bool   `gorm:"default:false;not null"`
				Input      string `gorm:"type:jsonb;default:'{}'"`
				Output     string `gorm:"type:jsonb"`
				Error      string `gorm:"type:text"`
				Logs       string `gorm:"type:text"`
				StartedAt  *time.Time
				EndedAt    *time.Time
				DurationMs int64 `gorm:"-"`
				CreatedAt  time.Time
				UpdatedAt  time.Time
			}

			type Connection struct {
				ID             string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
				UserID         string `gorm:"type:uuid;index;not null"`
				Name           string `gorm:"size:255;not null"`
				ServiceType    string `gorm:"size:100;not null;index"`
				Config         string `gorm:"type:jsonb;not null"`
				IsActive       bool   `gorm:"default:true;not null"`
				LastTestedAt   *time.Time
				LastTestStatus string `gorm:"size:20"`
				Version        int    `gorm:"default:1;not null"`
				CreatedAt      time.Time
				UpdatedAt      time.Time
			}

			type Subscription struct {
				ID                string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
				UserID            string     `gorm:"type:uuid;uniqueIndex;not null"`
				PlanID            string     `gorm:"size:50;not null;index"`
				Status            string     `gorm:"size:20;default:'active';not null;index"`
				CurrentPeriodEnd  *time.Time `gorm:"index"`
				Metadata          string     `gorm:"type:jsonb;default:'{}'"`
				CancelAtPeriodEnd bool       `gorm:"default:false;not null"`
				CreatedAt         time.Time
				UpdatedAt         time.Time
			}

			type Notification struct {
				ID        string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
				UserID    string `gorm:"type:uuid;index;not null"`
				Title     string `gorm:"size:255;not null"`
				Message   string `gorm:"type:text;not null"`
				Type      string `gorm:"size:50;not null;index"`
				IsRead    bool   `gorm:"default:false;not null;index"`
				ReadAt    *time.Time
				Metadata  string     `gorm:"type:jsonb;default:'{}'"`
				ExpiresAt *time.Time `gorm:"index"`
				CreatedAt time.Time
			}

			type Template struct {
				ID          string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
				Name        string `gorm:"size:255;not null"`
				Description string `gorm:"type:text"`
				Category    string `gorm:"size:100;index"`
				PreviewURL  string `gorm:"size:255"`
				Config      string `gorm:"type:jsonb;not null"`
				IsActive    bool   `gorm:"default:true;not null;index"`
				Version     string `gorm:"size:20;not null"`
				Tags        string `gorm:"type:jsonb;default:'[]'"`
				CreatedAt   time.Time
				UpdatedAt   time.Time
			}

			return db.Transaction(func(tx *gorm.DB) error {
				// Enable required extensions
				if err := tx.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
					return err
				}
				if err := tx.Exec("CREATE EXTENSION IF NOT EXISTS \"pgcrypto\"").Error; err != nil {
					return err
				}

				return tx.AutoMigrate(
					&User{},
					&Workflow{},
					&Execution{},
					&Connection{},
					&Subscription{},
					&Notification{},
					&Template{},
				)
			})
		},
		Rollback: func(db *gorm.DB) error {
			return db.Transaction(func(tx *gorm.DB) error {
				tables := []string{
					"notifications",
					"templates",
					"executions",
					"workflows",
					"connections",
					"subscriptions",
					"users",
				}

				for _, table := range tables {
					if err := tx.Migrator().DropTable(table); err != nil {
						return err
					}
				}

				// Drop extensions
				tx.Exec("DROP EXTENSION IF EXISTS \"pgcrypto\"")
				return tx.Exec("DROP EXTENSION IF EXISTS \"uuid-ossp\"").Error
			})
		},
	},
}
