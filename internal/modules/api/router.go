package api

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"s4s-backend/internal/config"
	"s4s-backend/internal/modules/api/handlers"
	"s4s-backend/internal/modules/api/middleware"
	authHandlers "s4s-backend/internal/modules/auth/handlers"
	authRepo "s4s-backend/internal/modules/auth/repository"
	authServices "s4s-backend/internal/modules/auth/services"
	connectionHandlers "s4s-backend/internal/modules/connection/handlers"
	connectionRepo "s4s-backend/internal/modules/connection/repository"
	connectionServices "s4s-backend/internal/modules/connection/services"
	notificationHandlers "s4s-backend/internal/modules/notification/handlers"
	notificationRepo "s4s-backend/internal/modules/notification/repository"
	notificationServices "s4s-backend/internal/modules/notification/services"
	subscriptionHandlers "s4s-backend/internal/modules/subscription/handlers"
	subscriptionRepo "s4s-backend/internal/modules/subscription/repository"
	subscriptionServices "s4s-backend/internal/modules/subscription/services"
	workflowRepo "s4s-backend/internal/modules/workflow/repository"
	workflowServices "s4s-backend/internal/modules/workflow/services"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB, cfg *config.Config) {
	// Initialize repositories
	userRepository := authRepo.NewUserRepository(db)
	workflowRepository := workflowRepo.NewWorkflowRepository(db)
	executionRepository := workflowRepo.NewExecutionRepository(db)
	templateRepository := workflowRepo.NewTemplateRepository(db)
	subscriptionRepository := subscriptionRepo.NewSubscriptionRepository(db)
	notificationRepository := notificationRepo.NewNotificationRepository(db)
	connectionRepository := connectionRepo.NewConnectionRepository(db)

	// Initialize services
	authService := authServices.NewAuthService(
		userRepository,
		subscriptionRepository,
		notificationRepository,
		cfg.JWTSecret,
	)
	userService := authServices.NewUserService(userRepository)
	workflowService := workflowServices.NewWorkflowService(
		workflowRepository,
		executionRepository,
		subscriptionRepository,
	)
	executionService := workflowServices.NewExecutionService(executionRepository)
	templateService := workflowServices.NewTemplateService(templateRepository)
	subscriptionService := subscriptionServices.NewSubscriptionService(subscriptionRepository, userRepository)
	notificationService := notificationServices.NewNotificationService(notificationRepository)
	connectionService := connectionServices.NewConnectionService(connectionRepository)

	// Initialize handlers
	authHandler := authHandlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	workflowHandler := handlers.NewWorkflowHandler(workflowService)
	executionHandler := handlers.NewExecutionHandler(executionService)
	templateHandler := handlers.NewTemplateHandler(templateService)
	subscriptionHandler := subscriptionHandlers.NewSubscriptionHandler(subscriptionService)
	notificationHandler := notificationHandlers.NewNotificationHandler(notificationService)
	connectionHandler := connectionHandlers.NewConnectionHandler(connectionService)

	// Apply global middleware
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.RequestLogger())

	// API routes
	api := r.Group("/api/v1")
	{
		// Health check
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok", "message": "Server is running"})
		})

		// Auth routes (public)
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.POST("/forgot-password", authHandler.ForgotPassword)
			auth.POST("/reset-password", authHandler.ResetPassword)
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			// User routes
			users := protected.Group("/users")
			{
				users.GET("/me", userHandler.GetCurrentUser)
				users.PUT("/me", userHandler.UpdateCurrentUser)
				users.DELETE("/me", userHandler.DeleteCurrentUser)
				users.POST("/invite", userHandler.InviteUser)
				users.PUT("/password", userHandler.ChangePassword)
			}

			// Workflow routes
			workflows := protected.Group("/workflows")
			{
				workflows.GET("", workflowHandler.ListWorkflows)
				workflows.POST("", workflowHandler.CreateWorkflow)
				workflows.GET("/:id", workflowHandler.GetWorkflow)
				workflows.PUT("/:id", workflowHandler.UpdateWorkflow)
				workflows.DELETE("/:id", workflowHandler.DeleteWorkflow)
				workflows.POST("/:id/test", workflowHandler.TestWorkflow)
				workflows.POST("/:id/run", workflowHandler.RunWorkflow)
				workflows.GET("/:id/executions", executionHandler.ListWorkflowExecutions)
			}

			// Execution routes
			executions := protected.Group("/executions")
			{
				executions.GET("", executionHandler.ListExecutions)
				executions.GET("/:id", executionHandler.GetExecution)
				executions.DELETE("/:id", executionHandler.DeleteExecution)
				executions.GET("/:id/logs", executionHandler.GetExecutionLogs)
			}

			// Connection routes
			connections := protected.Group("/connections")
			{
				connections.GET("", connectionHandler.ListConnections)
				connections.POST("", connectionHandler.CreateConnection)
				connections.GET("/:id", connectionHandler.GetConnection)
				connections.PUT("/:id", connectionHandler.UpdateConnection)
				connections.DELETE("/:id", connectionHandler.DeleteConnection)
				connections.POST("/:id/test", connectionHandler.TestConnection)
			}

			// Template routes
			templates := protected.Group("/templates")
			{
				templates.GET("", templateHandler.ListTemplates)
				templates.GET("/categories", templateHandler.ListTemplateCategories)
				templates.GET("/:id", templateHandler.GetTemplate)
				templates.POST("/:id/use", templateHandler.UseTemplate)
			}

			// Subscription routes
			subscription := protected.Group("/subscription")
			{
				subscription.GET("", subscriptionHandler.GetSubscription)
				subscription.GET("/plans", subscriptionHandler.ListPlans)
				subscription.POST("", subscriptionHandler.CreateSubscription)
				subscription.PUT("/cancel", subscriptionHandler.CancelSubscription)
				subscription.GET("/invoices", subscriptionHandler.ListInvoices)
			}

			// Notification routes
			notifications := protected.Group("/notifications")
			{
				notifications.GET("", notificationHandler.ListNotifications)
				notifications.GET("/unread", notificationHandler.GetUnreadCount)
				notifications.PUT("/:id/read", notificationHandler.MarkAsRead)
				notifications.PUT("/read-all", notificationHandler.MarkAllAsRead)
				notifications.GET("/settings", notificationHandler.GetSettings)
				notifications.PUT("/settings", notificationHandler.UpdateSettings)
			}

			// Admin routes
			//admin := protected.Group("/admin")
			//admin.Use(middleware.AdminMiddleware())
			//{
			//	// User management
			//	admin.GET("/users", userHandler.AdminListUsers)
			//	admin.POST("/users", userHandler.AdminCreateUser)
			//	admin.GET("/users/:id", userHandler.AdminGetUser)
			//	admin.PUT("/users/:id", userHandler.AdminUpdateUser)
			//	admin.DELETE("/users/:id", userHandler.AdminDeleteUser)
			//	admin.POST("/users/:id/impersonate", userHandler.AdminImpersonateUser)
			//}
		}

		// Webhook endpoints (public)
		webhooks := api.Group("/webhooks")
		{
			webhooks.POST("/stripe", subscriptionHandler.HandleStripeWebhook)
			webhooks.POST("/:provider", connectionHandler.HandleIncomingWebhook)
		}
	}
}
