package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/ZeroDay0utplay/file-compressor-service/internal/api"
	"github.com/ZeroDay0utplay/file-compressor-service/internal/compressor"
	"github.com/ZeroDay0utplay/file-compressor-service/internal/compressor/gs"
	"github.com/ZeroDay0utplay/file-compressor-service/internal/limiter"
	"github.com/ZeroDay0utplay/file-compressor-service/internal/storage"
)

func getenvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil && i > 0 {
			return i
		}
	}
	return fallback
}

func getenvDuration(key string, fallback int) time.Duration {
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

	println("PORT: ", os.Getenv("PORT"))

	maxConcurrency := getenvInt("MAX_CONCURRENCY", 1)
	maxUploadMB := getenvInt("MAX_UPLOAD_MB", 100)
	requestTimeout := getenvDuration("REQUEST_TIMEOUT_SEC", 120)
	gsTimeout := getenvDuration("GS_TIMEOUT_SEC", 90)
	gsPreset := os.Getenv("GS_DEFAULT_PRESET")
	if gsPreset == "" {
		gsPreset = "/ebook"
	}

	pool := limiter.New(maxConcurrency)
	temp := storage.New()
	reg := compressor.NewRegistry()

	pdfBackend := gs.New(gs.Config{
		DefaultPreset: gsPreset,
		Timeout:       gsTimeout,
	})
	reg.Register("application/pdf", pdfBackend)

	handler := api.New(api.Dependencies{
		Registry:       reg,
		Limiter:        pool,
		TempStore:      temp,
		MaxUploadBytes: int64(maxUploadMB) * 1024 * 1024,
		RequestTimeout: requestTimeout,
	})

	router := gin.New()
	router.Use(
		gin.Recovery(),
		gin.Logger(),
		api.AuthMiddleware(),
	)
	router.GET("/health", handler.Health)
	router.POST("/compress", handler.Compress)

	log.Printf("listening on :%s concurrency=%d maxUploadMB=%d", port, maxConcurrency, maxUploadMB)
	log.Fatal(router.Run(":" + port))
}
