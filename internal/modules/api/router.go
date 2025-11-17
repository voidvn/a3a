package api

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var handlers = struct {
	// Auth
	Register     gin.HandlerFunc
	Login        gin.HandlerFunc
	RefreshToken gin.HandlerFunc
	Logout       gin.HandlerFunc
	Me           gin.HandlerFunc

	// Workflows
	ListWorkflows  gin.HandlerFunc
	CreateWorkflow gin.HandlerFunc
	GetWorkflow    gin.HandlerFunc
	UpdateWorkflow gin.HandlerFunc
	DeleteWorkflow gin.HandlerFunc
	RunWorkflow    gin.HandlerFunc
	TestWorkflow   gin.HandlerFunc
	ListExecutions gin.HandlerFunc
	GetExecution   gin.HandlerFunc

	// Connections
	ListConnections  gin.HandlerFunc
	CreateConnection gin.HandlerFunc
	GetConnection    gin.HandlerFunc
	UpdateConnection gin.HandlerFunc
	DeleteConnection gin.HandlerFunc
	TestConnection   gin.HandlerFunc

	// Billing
	ListPlans          gin.HandlerFunc
	GetSubscription    gin.HandlerFunc
	CreateSubscription gin.HandlerFunc
	StripeWebhook      gin.HandlerFunc
	CancelSubscription gin.HandlerFunc

	// Settings
	GetNotificationSettings    gin.HandlerFunc
	UpdateNotificationSettings gin.HandlerFunc
	GetUsageStats              gin.HandlerFunc

	// System
	HealthCheck       gin.HandlerFunc
	PrometheusHandler gin.HandlerFunc
}{}

func init() {
	// === AUTH ===
	handlers.Register = func(c *gin.Context) { c.JSON(201, gin.H{"message": "stub: register ok", "user_id": "usr_123"}) }
	handlers.Login = func(c *gin.Context) {
		c.JSON(200, gin.H{"access_token": "jwt.stub.token", "refresh_token": "refresh.stub", "expires_in": 3600})
	}
	handlers.RefreshToken = func(c *gin.Context) { c.JSON(200, gin.H{"access_token": "jwt.stub.refreshed"}) }
	handlers.Logout = func(c *gin.Context) { c.JSON(200, gin.H{"message": "logged out"}) }
	handlers.Me = func(c *gin.Context) {
		c.JSON(200, gin.H{
			"id":        "usr_123",
			"email":     "test@example.com",
			"full_name": "Test User",
			"role":      "user",
		})
	}

	// === WORKFLOWS ===
	handlers.ListWorkflows = func(c *gin.Context) { c.JSON(200, gin.H{"data": []any{}, "total": 0}) }
	handlers.CreateWorkflow = func(c *gin.Context) { c.JSON(201, gin.H{"id": "wf_123", "message": "workflow created"}) }
	handlers.GetWorkflow = func(c *gin.Context) { c.JSON(200, gin.H{"id": c.Param("id"), "name": "Test Workflow"}) }
	handlers.UpdateWorkflow = func(c *gin.Context) { c.JSON(200, gin.H{"message": "updated"}) }
	handlers.DeleteWorkflow = func(c *gin.Context) { c.JSON(200, gin.H{"message": "deleted"}) }
	handlers.RunWorkflow = func(c *gin.Context) { c.JSON(202, gin.H{"execution_id": "exec_123", "status": "queued"}) }
	handlers.TestWorkflow = func(c *gin.Context) { c.JSON(202, gin.H{"execution_id": "exec_test_123", "status": "test_queued"}) }
	handlers.ListExecutions = func(c *gin.Context) { c.JSON(200, gin.H{"executions": []any{}}) }
	handlers.GetExecution = func(c *gin.Context) { c.JSON(200, gin.H{"log": "success", "status": "success"}) }

	// === CONNECTIONS ===
	handlers.ListConnections = func(c *gin.Context) { c.JSON(200, gin.H{"connections": []any{}}) }
	handlers.CreateConnection = func(c *gin.Context) { c.JSON(201, gin.H{"id": "conn_123", "service": "slack"}) }
	handlers.GetConnection = func(c *gin.Context) { c.JSON(200, gin.H{"id": c.Param("id"), "service": "gmail"}) }
	handlers.UpdateConnection = func(c *gin.Context) { c.JSON(200, gin.H{"message": "updated"}) }
	handlers.DeleteConnection = func(c *gin.Context) { c.JSON(200, gin.H{"message": "deleted"}) }
	handlers.TestConnection = func(c *gin.Context) { c.JSON(200, gin.H{"status": "connected"}) }

	// === BILLING ===
	handlers.ListPlans = func(c *gin.Context) {
		c.JSON(200, gin.H{"plans": []gin.H{{"id": "freemium", "price": 0}, {"id": "starter", "price": 29}}})
	}
	handlers.GetSubscription = func(c *gin.Context) { c.JSON(200, gin.H{"plan": "freemium", "status": "active"}) }
	handlers.CreateSubscription = func(c *gin.Context) { c.JSON(201, gin.H{"checkout_url": "https://checkout.stripe.test/success"}) }

	// Stripe Webhook — РАБОЧАЯ заглушка с RAW body
	handlers.StripeWebhook = func(c *gin.Context) {
		rawBody, err := c.GetRawData()
		if err != nil {
			c.JSON(400, gin.H{"error": "cannot read body"})
			return
		}
		log.Printf("[STRIPE WEBHOOK] received %d bytes", len(rawBody))
		c.Status(http.StatusOK) // ← обязательно 200, иначе Stripe будет ретраить
	}

	handlers.CancelSubscription = func(c *gin.Context) { c.JSON(200, gin.H{"message": "canceled"}) }

	// === SETTINGS ===
	handlers.GetNotificationSettings = func(c *gin.Context) { c.JSON(200, gin.H{"email_enabled": true}) }
	handlers.UpdateNotificationSettings = func(c *gin.Context) { c.JSON(200, gin.H{"message": "saved"}) }
	handlers.GetUsageStats = func(c *gin.Context) { c.JSON(200, gin.H{"executions_used": 42, "limit": 1000}) }

	// === SYSTEM ===
	handlers.HealthCheck = func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) }
	handlers.PrometheusHandler = func(c *gin.Context) { c.String(200, "# stub metrics") }
}

func NewRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(cors.Default())
	return r
}

func RegisterRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	{
		// Auth
		auth := v1.Group("/auth")
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
			auth.POST("/refresh", handlers.RefreshToken)
			auth.POST("/logout", handlers.Logout)
			auth.GET("/me", handlers.Me)
		}

		// Workflows
		wf := v1.Group("/workflows")
		{
			wf.GET("", handlers.ListWorkflows)
			wf.POST("", handlers.CreateWorkflow)
			wf.GET("/:id", handlers.GetWorkflow)
			wf.PUT("/:id", handlers.UpdateWorkflow)
			wf.DELETE("/:id", handlers.DeleteWorkflow)
			wf.POST("/:id/run", handlers.RunWorkflow)
			wf.POST("/:id/test", handlers.TestWorkflow)
			wf.GET("/:id/executions", handlers.ListExecutions)
			wf.GET("/:id/executions/:exec_id", handlers.GetExecution)
		}

		// Connections
		conn := v1.Group("/connections")
		{
			conn.GET("", handlers.ListConnections)
			conn.POST("", handlers.CreateConnection)
			conn.GET("/:id", handlers.GetConnection)
			conn.PUT("/:id", handlers.UpdateConnection)
			conn.DELETE("/:id", handlers.DeleteConnection)
			conn.POST("/:id/test", handlers.TestConnection)
		}

		// Billing
		billing := v1.Group("/billing")
		{
			billing.GET("/plans", handlers.ListPlans)
			billing.GET("/subscription", handlers.GetSubscription)
			billing.POST("/subscription", handlers.CreateSubscription)
			billing.POST("/webhook/stripe", handlers.StripeWebhook) // ← теперь всё ок
			billing.PUT("/subscription/cancel", handlers.CancelSubscription)
		}

		// Settings
		settings := v1.Group("/settings")
		{
			settings.GET("/notifications", handlers.GetNotificationSettings)
			settings.PUT("/notifications", handlers.UpdateNotificationSettings)
			settings.GET("/usage", handlers.GetUsageStats)
		}
	}

	// Health
	r.GET("/health", handlers.HealthCheck)
	r.GET("/metrics", handlers.PrometheusHandler)
}
