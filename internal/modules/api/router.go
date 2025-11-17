package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	// CORS для фронта
	r.Use(cors.Default())

	return r
}

func SetupGlobalMiddleware(r *gin.Engine) {
	//
}

func RegisterRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	{
		// === AUTH MODULE ===
		auth := v1.Group("/auth")
		{
			auth.POST("/register", handlers.Register)       // 201 + JWT
			auth.POST("/login", handlers.Login)             // 200 + JWT
			auth.POST("/refresh", handlers.RefreshToken)    // новый токен
			auth.POST("/logout", handlers.Logout)           // blacklist (опционально)
			auth.GET("/me", middleware.Auth(), handlers.Me) // профиль юзера
		}

		// === WORKFLOWS MODULE ===
		workflows := v1.Group("/workflows")
		workflows.Use(middleware.Auth())
		{
			workflows.GET("", handlers.ListWorkflows)         // список всех
			workflows.POST("", handlers.CreateWorkflow)       // создать новый
			workflows.GET("/:id", handlers.GetWorkflow)       // получить по ID
			workflows.PUT("/:id", handlers.UpdateWorkflow)    // обновить граф
			workflows.DELETE("/:id", handlers.DeleteWorkflow) // удалить

			workflows.POST("/:id/run", handlers.RunWorkflow)                 // запуск (async)
			workflows.POST("/:id/test", handlers.TestWorkflow)               // тестовый запуск с payload
			workflows.GET("/:id/executions", handlers.ListExecutions)        // история запусков
			workflows.GET("/:id/executions/:exec_id", handlers.GetExecution) // лог одного запуска
		}

		// === CONNECTIONS (интеграции) ===
		connections := v1.Group("/connections")
		connections.Use(middleware.Auth())
		{
			connections.GET("", handlers.ListConnections)
			connections.POST("", handlers.CreateConnection) // Slack, Gmail, CRM и т.д.
			connections.GET("/:id", handlers.GetConnection)
			connections.PUT("/:id", handlers.UpdateConnection)
			connections.DELETE("/:id", handlers.DeleteConnection)
			connections.POST("/:id/test", handlers.TestConnection) // проверка соединения
		}

		// === SUBSCRIPTION & BILLING ===
		billing := v1.Group("/billing")
		billing.Use(middleware.Auth())
		{
			billing.GET("/plans", handlers.ListPlans) // freemium, starter, team
			billing.GET("/subscription", handlers.GetSubscription)
			billing.POST("/subscription", handlers.CreateSubscription) // Stripe Checkout
			billing.POST("/webhook/stripe", handlers.StripeWebhook)    // raw body!
			billing.PUT("/subscription/cancel", handlers.CancelSubscription)
		}

		// === USER SETTINGS ===
		settings := v1.Group("/settings")
		settings.Use(middleware.Auth())
		{
			settings.GET("/notifications", handlers.GetNotificationSettings)
			settings.PUT("/notifications", handlers.UpdateNotificationSettings)
			settings.GET("/usage", handlers.GetUsageStats) // лимиты, executions count и т.д.
		}

		// === HEALTH & METRICS ===
		r.GET("/health", handlers.HealthCheck)
		r.GET("/metrics", handlers.PrometheusHandler) // если подключишь promhttp
	}
}
