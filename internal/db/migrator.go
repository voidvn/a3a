package db

import (
	"s4s-backend/internal/db/migrations"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{})

	m.InitSchema(func(db *gorm.DB) error {
		// создаём extension для uuid если нет
		db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
		db.Exec("CREATE EXTENSION IF NOT EXISTS \"pgcrypto\"") // для gen_random_uuid()
		return nil
	})

	migrationsList := append([]*gormigrate.Migration{}, migrations.InitialSchema...)
	m = gormigrate.New(db, gormigrate.DefaultOptions, migrationsList)

	return m.Migrate()
}
