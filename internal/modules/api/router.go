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
	connectionRepo "s4s-backend/internal/modules/connection/repository"
	notificationRepo "s4s-backend/internal/modules/notification/repository"
	subscriptionRepo "s4s-backend/internal/modules/subscription/repository"
	workflowRepo "s4s-backend/internal/modules/workflow/repository"
	workflowServices "s4s-backend/internal/modules/workflow/services"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB, cfg *config.Config) {
	// Initialize repositories
	userRepository := authRepo.NewUserRepository(db)
	workflowRepository := workflowRepo.NewWorkflowRepository(db)
	executionRepository := workflowRepo.NewExecutionRepository(db)
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

	// Initialize handlers
	authHandler := authHandlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	workflowHandler := handlers.NewWorkflowHandler(workflowService)
	executionHandler := handlers.NewExecutionHandler(executionService)

	// Apply global middleware
	r.Use(middleware.CORSMiddleware())

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
			auth.POST("/forgot-password", authHandler.ForgotPassword)
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
			}

			// Execution routes
			executions := protected.Group("/executions")
			{
				executions.GET("", executionHandler.ListExecutions)
				executions.GET("/:id", executionHandler.GetExecution)
				executions.DELETE("/:id", executionHandler.DeleteExecution)
			}

			// TODO: Add routes for connections, templates, subscriptions, notifications
		}

		// Admin routes
		//admin := protected.Group("/admin")
		//admin.Use(middleware.AdminMiddleware())
		//{
		//	// TODO: Add admin routes
		//}
	}
}
