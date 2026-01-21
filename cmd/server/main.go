package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var (
	version   = "dev"
	commit    = "unknown"
	buildTime = "unknown"
)

func main() {
	r := gin.New()
	r.Use(gin.Recovery(), gin.Logger())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "ok",
			"version":   version,
			"commit":    commit,
			"buildTime": buildTime,
		})
	})

	r.POST("/compress", func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{
			"error": "compression not implemented yet",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Println("listening on :" + port)
	log.Fatal(r.Run(":" + port))
}
