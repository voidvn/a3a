package db

import (
	"context"
	"database/sql"
	"log"

	"github.com/spf13/viper"

	"gitlab.com/ncodeGroup/s4s-backend/internal/models/ent"

	_ "github.com/jackc/pgx/v5/stdlib" // database/sql driver name: "pgx"
	_ "modernc.org/sqlite"             // database/sql driver name: "sqlite"

	entsql "entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/schema"
)

// ... existing code ...

var Client *ent.Client

func Connect() {
	driver := viper.GetString("DB_DRIVER")
	dsn := viper.GetString("DB_URL")

	if driver == "" {
		driver = "sqlite"
	}

	var (
		sqlDriverName string
		dialect       string
	)

	switch driver {
	case "postgres":
		sqlDriverName = "pgx"
		dialect = "postgres"
		if dsn == "" {
			dsn = "postgres://postgres:postgres@db:5432/postgres?sslmode=disable"
		}
	case "sqlite":
		sqlDriverName = "sqlite3"
		dialect = "sqlite3"
		if dsn == "" {
			dsn = "file:app.db?_pragma=foreign_keys(1)"
		}
	default:
		log.Fatalf("Unsupported DB driver: %s", driver)
	}

	// Create ent.Client and run the schema migration.
	db, err := sql.Open(driver, dsn)
	if err != nil {
		log.Fatalf("failed opening connection to %s: %v", driver, err)
	}

	// Create an ent.Driver from `db`.
	drv := entsql.OpenDB(dialect, db)

	// Create the ent client with the driver.
	c := ent.NewClient(ent.Driver(drv))

	// Run migrations with the appropriate dialect
	if err := c.Schema.Create(
		context.Background(),
		schema.WithDropIndex(true),   // Drop indexes if they were removed in the schema
		schema.WithForeignKeys(true), // Enable foreign keys
		schema.WithDialect(dialect),  // Use the appropriate dialect
	); err != nil {
		log.Fatalf("failed running migrations: %v", err)
	}

	Client = c
	log.Printf("DB connected (%s via %s)", driver, sqlDriverName)
}
