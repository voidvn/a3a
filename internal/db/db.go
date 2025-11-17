package db

import (
	"fmt"
	"s4s-backend/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.GetString("DB_HOST", "localhost"),
		config.GetString("DB_USER", "postgres"),
		config.GetString("DB_PASSWORD", "postgres"),
		config.GetString("DB_NAME", "s4s"),
		config.GetString("DB_PORT", "5432"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	DB = db
	return db, nil
}
