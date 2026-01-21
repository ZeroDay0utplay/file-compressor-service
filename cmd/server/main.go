package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/ZeroDay0utplay/file-compressor-service/internal/api"
	"github.com/ZeroDay0utplay/file-compressor-service/internal/compressor/gs"
	"github.com/ZeroDay0utplay/file-compressor-service/internal/limiter"
	"github.com/ZeroDay0utplay/file-compressor-service/internal/storage"
)

func envInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil && i > 0 {
			return i
		}
	}
	return fallback
}

func envSeconds(key string, fallback int) time.Duration {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil && i > 0 {
			return time.Duration(i) * time.Second
		}
	}
	return time.Duration(fallback) * time.Second
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	maxConcurrency := envInt("MAX_CONCURRENCY", 1)
	maxUploadMB := envInt("MAX_UPLOAD_MB", 100)
	requestTimeout := envSeconds("REQUEST_TIMEOUT_SEC", 120)

	pool := limiter.New(maxConcurrency)
	tempStore := storage.New()

	compressor := gs.New(gs.Config{
		DefaultPreset: "/ebook",
		Timeout:       90 * time.Second,
	})

	handler := api.New(api.Dependencies{
		Compressor:     compressor,
		WorkerPool:     pool,
		TempStore:      tempStore,
		MaxUploadBytes: int64(maxUploadMB) * 1024 * 1024,
		RequestTimeout: requestTimeout,
	})

	router := gin.New()
	router.Use(gin.Recovery(), gin.Logger())

	router.GET("/health", handler.Health)
	router.POST("/compress", handler.Compress)

	log.Printf("starting on :%s maxConcurrency=%d maxUploadMB=%d", port, maxConcurrency, maxUploadMB)
	log.Fatal(router.Run(":" + port))
}
