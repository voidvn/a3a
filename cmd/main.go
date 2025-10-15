package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gitlab.com/ncodeGroup/s4s-backend/internal/admin"
	"gitlab.com/ncodeGroup/s4s-backend/internal/config"
	"gitlab.com/ncodeGroup/s4s-backend/internal/db"
)

func main() {
	// Initialize configuration
	config.Init()
	// Initialize database
	db.Connect()

	// Create a new Gin router
	r := gin.Default()

	// Simple health check endpoint
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	// Setup admin interface
	setupAdmin(r)

	// Run migrations
	runMigrations()

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on :%s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func setupAdmin(r *gin.Engine) {
	// Get database connection string from environment
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	if dbHost == "" {
		dbHost = "db" // Default to service name in docker-compose
	}
	if dbPort == "" {
		dbPort = "5432"
	}
	if dbUser == "" {
		dbUser = "postgres"
	}
	if dbPassword == "" {
		dbPassword = "postgres"
	}
	if dbName == "" {
		dbName = "s4sdb"
	}

	dbURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbName, dbPassword)

	// Initialize admin interface
	adminHandler, err := admin.SetupAdmin(dbURL)
	if err != nil {
		log.Printf("Failed to initialize admin interface: %v", err)
		return
	}

	// Mount admin to /admin
	r.Any("/admin/*resources", gin.WrapH(adminHandler))
	log.Println("Admin interface available at /admin")
}

func runMigrations() {
	// Get database configuration from environment
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Create the database URL for migrations
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser,
		dbPass,
		dbHost,
		dbPort,
		dbName,
	)

	log.Printf("Running database migrations on %s\n", dbURL)

	// Run migrations using the database client
	// The migrations will be handled by ent's automatic migration
	// since we already called db.Connect() which runs the migrations
}
