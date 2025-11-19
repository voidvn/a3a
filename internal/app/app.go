package app

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"s4s-backend/internal/config"
	"s4s-backend/internal/db"
	"s4s-backend/internal/modules/admin"
	"s4s-backend/internal/modules/api"
	"s4s-backend/internal/modules/api/middleware"
	"syscall"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func Start() {
	// 1. Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// 2. Connect to PostgreSQL
	database, err := db.Connect()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// 3. Run database migrations
	if err := db.RunMigrations(database); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	// 4. Connect to Redis (for cache/sessions)	//_, err = db.ConnectRedis()
	//if err != nil {
	//	log.Fatalf("failed to connect Redis: %v", err)
	//}

	// 5. Connect to RabbitMQ (for queues)
	//rabbitConn, err := db.ConnectRabbitMQ()
	//if err != nil {
	//	log.Fatalf("failed to connect RabbitMQ: %v", err)
	//}
	//defer rabbitConn.Close()

	// 6. Initialize the main Gin router
	r := gin.Default()
	r.Use(middleware.RequestLogger())
	r.Use(middleware.CORSMiddleware())

	// 7. Initialize API routes
	api.SetupRoutes(r, database, cfg)

	// 8. Initialize the admin panel
	adminConfig := admin.GetAdminConfig(
		os.Getenv("DB_URL"),
		os.Getenv("ADMIN_APP_KEY"),
	)
	admin.InitAdmin(r, adminConfig)

	// 9. Configure the HTTP server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 10. Start server in a goroutine
	go func() {
		log.Printf("Server starting on :%s", port)
		if err := http.ListenAndServe(":"+port, r); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 11. Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
}
