// File: internal/middleware/cors.go

package middleware

import (
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func CORSMiddleware() gin.HandlerFunc {
	err := godotenv.Load()
	if err != nil {
		DebugLog("Error loading .env file:", err)
	}
	allowOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")

	// Default origins if not set
	if allowOrigins == "" {
		allowOrigins = "http://localhost:3000,http://localhost:3001,http://localhost:5173,http://localhost:8080,http://127.0.0.1:3000,http://127.0.0.1:3001,http://127.0.0.1:5173,http://127.0.0.1:8080"
	}

	origins := strings.Split(allowOrigins, ",")
	DebugLog("CORS Allow Origins: %v", origins)

	return cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "Accept", "x-user-role"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
