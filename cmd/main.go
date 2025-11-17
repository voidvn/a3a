package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	//"s4s-backend/internal/modules/admin"
	"syscall"
	"time"

	"s4s-backend/internal/config"
	"s4s-backend/internal/db"
	"s4s-backend/internal/modules/api"
	//_ "s4s-backend/internal/modules/auth"
	//_ "s4s-backend/internal/modules/workflow"

	"github.com/GoAdminGroup/go-admin/engine"
	goAdminConfig "github.com/GoAdminGroup/go-admin/modules/config"
	_ "github.com/GoAdminGroup/go-admin/modules/db/drivers/postgres"
	"github.com/GoAdminGroup/go-admin/modules/language"

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

	// Instantiate a GoAdmin engine object.
	eng := engine.Default()

	// GoAdmin global configuration, can also be imported as a json file.
	cfg := goAdminConfig.Config{
		Databases: goAdminConfig.DatabaseList{
			"default": {
				Host:         config.GetString("DB_HOST", "localhost"),
				Port:         config.GetString("DB_PORT", "5432"),
				User:         config.GetString("DB_USER", "postgres"),
				Pwd:          config.GetString("DB_PASSWORD", "postgres"),
				Name:         config.GetString("DB_NAME", "postgres"),
				MaxIdleConns: 50,
				MaxOpenConns: 150,
				Driver:       goAdminConfig.DriverPostgresql,
			},
		},
		UrlPrefix: "admin", // The url prefix of the website.
		// Store must be set and guaranteed to have write access, otherwise new administrator users cannot be added.
		Store: goAdminConfig.Store{
			Path:   "./uploads",
			Prefix: "uploads",
		},
		Language: language.EN,
		Theme:    "adminlte",
	}

	// 9. Add GoAdmin to Gin
	if err := eng.AddConfig(&cfg).Use(router); err != nil {
		log.Fatalf("failed to add GoAdmin config to engine: %v", err)
	}

	_ = r.Run(":9033")

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
