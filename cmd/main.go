package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"s4s-backend/internal/config"
	"s4s-backend/internal/db"
	"s4s-backend/internal/modules/admin"
	_ "s4s-backend/internal/modules/api"
	"s4s-backend/internal/modules/api/middleware"
	_ "strconv"
	"syscall"

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
	//_, err = db.ConnectRedis()
	//if err != nil {
	//	log.Fatalf("failed to connect Redis: %v", err)
	//}

	// 5. Подключаемся к RabbitMQ (для очередей)
	//rabbitConn, err := db.ConnectRabbitMQ()
	//if err != nil {
	//	log.Fatalf("failed to connect RabbitMQ: %v", err)
	//}
	//defer rabbitConn.Close()

	// 6. Инициализируем основной роутер Gin
	r := gin.Default()
	r.Use(middleware.RequestLogger())

	// 7. Инициализируем API роуты
	//apiRouter := api.NewRouter()
	//api.RegisterRoutes(apiRouter)

	// 8. Инициализируем админ-панель
	adminConfig := admin.GetAdminConfig(
		os.Getenv("DB_URL"),
		os.Getenv("ADMIN_APP_KEY"),
	)
	admin.InitAdmin(r, adminConfig)

	// 9. Подключаем API роуты к основному роутеру
	//r.Any("/api/*any", gin.WrapH(apiRouter))

	// 10. Настраиваем HTTP сервер
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 11. Start HTTP server in a goroutine
	go func() {
		log.Printf("Server starting on :%s", port)
		if err := http.ListenAndServe(":"+port, r); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 12. Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
}
