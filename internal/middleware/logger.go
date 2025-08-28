// File: internal/middleware/logger.go

package middleware

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// DebugLogWithEmoji logs a message with emoji only if GIN_MODE=debug
func DebugLogWithEmoji(emoji, format string, args ...interface{}) {
	if os.Getenv("GIN_MODE") == "debug" {
		log.Printf("%s "+format, append([]interface{}{emoji}, args...)...)
	}
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)
		status := c.Writer.Status()
		log.Printf("%s %s %d %s", c.Request.Method, c.Request.URL.Path, status, latency)
	}
}

// DebugLog logs debug messages only when GIN_MODE=debug.
// If emoji is non-empty, it is prepended to the log message.
func DebugLog(format string, args ...interface{}) {
	if os.Getenv("GIN_MODE") == "debug" {
		log.Printf("[DEBUG] "+format, args...)
	}
}
