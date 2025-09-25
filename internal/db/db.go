package db

import (
	"context"
	"log"

	"github.com/spf13/viper"

	"gitlab.com/ncodeGroup/s4s-backend/internal/models/ent"

	_ "github.com/jackc/pgx/v5/stdlib" // database/sql driver name: "pgx"
	_ "modernc.org/sqlite"             // database/sql driver name: "sqlite"

	"entgo.io/ent/dialect"
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

	var sqlDriverName string
	var migrateDialect string
	switch driver {
	case "postgres":
		sqlDriverName = "pgx"
		migrateDialect = dialect.Postgres
		// DSN пример: "postgres://user:pass@host:5432/dbname?sslmode=disable"
	case "sqlite":
		sqlDriverName = "sqlite"
		migrateDialect = dialect.SQLite
		if dsn == "" {
			dsn = "file:app.db?_pragma=foreign_keys(1)"
		}
	default:
		log.Fatalf("Unsupported DB driver: %s", driver)
	}

	drv, err := entsql.Open(sqlDriverName, dsn)
	if err != nil {
		log.Fatalf("failed opening SQL driver %q: %v", sqlDriverName, err)
	}

	c := ent.NewClient(ent.Driver(drv))

	// ЯВНО указываем диалект для мигратора (нужно для modernc "sqlite")
	if err := c.Schema.Create(context.Background(), schema.WithDialect(migrateDialect)); err != nil {
		log.Fatalf("failed running migrations: %v", err)
	}

	Client = c
	log.Printf("DB connected (%s via %s)", driver, sqlDriverName)
}
