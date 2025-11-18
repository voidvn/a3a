package api

import (
	"net/http"
	"time"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"s4s-backend/internal/config"
	"s4s-backend/internal/modules/api/handlers"
	"s4s-backend/internal/modules/api/middleware"
	authHandlers "s4s-backend/internal/modules/auth/handlers"
	authRepo "s4s-backend/internal/modules/auth/repository"
	authServices "s4s-backend/internal/modules/auth/services"
	workflowRepo "s4s-backend/internal/modules/workflow/repository"
	workflowServices "s4s-backend/internal/modules/workflow/services"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB, cfg *config.Config) {
	// Initialize repositories
	userRepository := authRepo.NewUserRepository(db)
	workflowRepository := workflowRepo.NewWorkflowRepository(db)
	executionRepository := workflowRepo.NewExecutionRepository(db)

	// Initialize services
	authService := authServices.NewAuthService(
		userRepository,
		nil, // subscription repo not needed for demo
		nil, // notification repo not needed for demo
		cfg.JWT.Secret,
	)
	userService := authServices.NewUserService(userRepository)
	workflowService := workflowServices.NewWorkflowService(
		workflowRepository,
		executionRepository,
		nil, // subscription service not needed for demo
	)

	// Initialize handlers
	authHandler := authHandlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	workflowHandler := handlers.NewWorkflowHandler(workflowService)

	// Apply global middleware
	r.Use(
		middleware.CORSMiddleware(),
		requestid.New(),
		middleware.RequestLogger(),
		//middleware.Recovery(),
	)

	// API routes
	api := r.Group("/api/v1")
	{
		// Health check
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status": "ok",
				//"version": cfg.Version,
				"time": time.Now().UTC().Format(time.RFC3339),
			})
		})

		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			//auth.POST("/refresh", authHandler.RefreshToken)
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
		{
			// User routes
			users := protected.Group("/users")
			{
				users.GET("/me", userHandler.GetCurrentUser)
				users.PUT("/me", userHandler.UpdateCurrentUser)
			}

			// Workflow routes
			workflows := protected.Group("/workflows")
			{
				workflows.GET("", workflowHandler.ListWorkflows)
				workflows.POST("", workflowHandler.CreateWorkflow)
				workflows.GET("/:id", workflowHandler.GetWorkflow)
				workflows.POST("/:id/run", workflowHandler.RunWorkflow)
			}
		}
	}

	// 404 handler
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "The requested resource was not found",
		})
	})
}
