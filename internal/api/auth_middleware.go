package api

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	raw := os.Getenv("API_KEYS")
	if raw == "" {
		log.Fatal("API_KEYS is required but not set")
	}

	allowed := map[string]struct{}{}
	for _, k := range strings.Split(raw, ",") {
		k = strings.TrimSpace(k)
		if k != "" {
			allowed[k] = struct{}{}
		}
	}

	if len(allowed) == 0 {
		log.Fatal("API_KEYS is empty")
	}

	return func(c *gin.Context) {
		key := c.GetHeader("X-API-Key")
		if _, ok := allowed[key]; !ok {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "unauthorized"},
			)
			return
		}
		c.Next()
	}
}
