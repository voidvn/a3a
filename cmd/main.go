package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"s4s-backend/internal/modules/admin"

	//"s4s-backend/internal/modules/admin"
	"syscall"
	"time"

	"s4s-backend/internal/config"
	"s4s-backend/internal/db"
	"s4s-backend/internal/modules/api"
	//_ "s4s-backend/internal/modules/auth"
	//_ "s4s-backend/internal/modules/workflow"

	_ "github.com/GoAdminGroup/go-admin/modules/db/drivers/postgres"
	_ "github.com/GoAdminGroup/themes/adminlte"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	// 1. Загружаем конфиг
	if err := config.Load(); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// 2. Подключаемся к PostgreSQL
	database, err := db.Connect()
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// 3. Запускаем миграции (AutoMigrate через GORM)
	if err := db.RunMigrations(database); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	// 4. Подключаемся к Redis (для кэша/сессий)
	_, err = db.ConnectRedis()
	if err != nil {
		log.Fatalf("failed to connect Redis: %v", err)
	}
	// Здесь инициализируй репозитории, если нужно: e.g., authRepo := repositories.NewAuthRepo(database, redisClient)

	// 5. Подключаемся к RabbitMQ (для очередей)
	rabbitConn, err := db.ConnectRabbitMQ()
	if err != nil {
		log.Fatalf("failed to connect RabbitMQ: %v", err)
	}
	defer rabbitConn.Close() // Закрываем в shutdown

	// 6. Инициализируем Gin роутер
	router := api.NewRouter()
	api.RegisterRoutes(router)

	// 7. Регистрируем модули с роутами (заглушки внутри модулей)
	//auth.RegisterRoutes(router.Group("/api/v1"))
	//workflow.RegisterRoutes(router.Group("/api/v1"))
	// Другие модули: connections, subscription и т.д.

	// 8. Запуск сервера
	port := config.GetString("SERVER_PORT", "8080")
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		log.Printf("Server starting on :%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	r := gin.Default()

	adminConfig := admin.GetAdminConfig(
		os.Getenv("DB_URL"),
		os.Getenv("ADMIN_APP_KEY"),
	)

	adminEngine := admin.InitAdmin(r, adminConfig)

	// Регистрация дашборда
	adminEngine.HTML("GET", "/admin", admin.GetDashboard)

	// 9. Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server stopped gracefully")
}
