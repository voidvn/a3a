package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"s4s-backend/internal/api"
	"s4s-backend/internal/config"
	"s4s-backend/internal/db"
	"s4s-backend/internal/modules/auth"
	"s4s-backend/internal/modules/workflow"
)

func main() {
	// 1. Загружаем конфиг (.env)
	if err := config.Load(); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// 2. Подключаемся к PostgreSQL
	database, err := db.Connect()
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// 3. Запускаем миграции (один раз — всё встанет)
	if err := db.RunMigrations(database); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	// 4. Инициализируем Gin роутер
	router := api.NewRouter()

	// 5. Регистрируем модули (каждый сам вешает свои роуты)
	auth.RegisterRoutes(router.Group("/auth"))
	workflow.RegisterRoutes(router.Group("/workflows"))
	// connections.RegisterRoutes(router.Group("/connections")) // когда будет готов

	// 6. Глобальные middleware (логгер, recovery, CORS и т.д.)
	api.SetupGlobalMiddleware(router)

	// 7. Запуск сервера
	port := config.GetString("SERVER_PORT", "8080")
	log.Printf("Server starting on :%s", port)

	go func() {
		if err := router.Run(":" + port); err != nil {
			log.Fatalf("server error: %v", err)
		}
	}()

	// 8. Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// тут можно добавить ctx shutdown если захочешь
	log.Println("Server stopped")
}
