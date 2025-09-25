package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/ncodeGroup/s4s-backend/internal/config"
	"gitlab.com/ncodeGroup/s4s-backend/internal/db"
)

func main() {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	config.Init()
	db.Connect()
	err := r.Run(":8080")
	if err != nil {
		return
	}
}
