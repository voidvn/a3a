package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip logging for these paths
		if c.Request.URL.Path == "/health" || c.Request.URL.Path == "/favicon.ico" {
			c.Next()
			return
		}

		// Log request
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// Process request
		c.Next()

		// Log details
		latency := time.Since(start)
		status := c.Writer.Status()

		// Basic log line
		log.Printf("[%s] %s | %d | %v | %s | %s",
			time.Now().Format("2006/01/02 15:04:05"),
			method,
			status,
			latency,
			c.ClientIP(),
			path,
		)

		// Log form data for POST requests
		if method == "POST" || method == "PUT" || method == "PATCH" {
			if err := c.Request.ParseForm(); err == nil && len(c.Request.PostForm) > 0 {
				log.Printf("Form data: %v", c.Request.PostForm)
			}
		}

		// Log errors if any
		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				log.Printf("Error: %v", e)
			}
		}
	}
}
