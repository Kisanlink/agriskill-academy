// File: internal/middleware/logger.go

package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
)

// LoggerConfig holds logging configuration
type LoggerConfig struct {
	Level       string `json:"level"`
	OutputPath  string `json:"output_path"`
	Format      string `json:"format"` // json or console
	Development bool   `json:"development"`
}

// Global logger instance
var logger *zap.Logger

// InitLogger initializes the global logger
func InitLogger(config LoggerConfig) error {
	// Parse log level
	level, err := zapcore.ParseLevel(config.Level)
	if err != nil {
		return fmt.Errorf("invalid log level: %w", err)
	}

	// Create encoder config
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	// Choose encoder
	var encoder zapcore.Encoder
	if config.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// Create output
	var output zapcore.WriteSyncer
	if config.OutputPath != "" {
		file, err := os.OpenFile(config.OutputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to open log file: %w", err)
		}
		output = zapcore.AddSync(file)
	} else {
		output = zapcore.AddSync(os.Stdout)
	}

	// Create core
	core := zapcore.NewCore(encoder, output, level)

	// Create logger
	if config.Development {
		logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	} else {
		logger = zap.New(core, zap.AddCaller())
	}

	return nil
}

// GetLogger returns the global logger instance
func GetLogger() *zap.Logger {
	if logger == nil {
		// Fallback to default logger if not initialized
		logger, _ = zap.NewProduction()
	}
	return logger
}

// DebugLogWithEmoji logs a message with emoji only if GIN_MODE=debug
func DebugLogWithEmoji(emoji, format string, args ...interface{}) {
	if os.Getenv("GIN_MODE") == "debug" {
		log.Printf("%s "+format, append([]interface{}{emoji}, args...)...)
	}
}

// DebugLog logs debug messages only when GIN_MODE=debug.
// If emoji is non-empty, it is prepended to the log message.
func DebugLog(format string, args ...interface{}) {
	if os.Getenv("GIN_MODE") == "debug" {
		log.Printf("[DEBUG] "+format, args...)
	}
}

// Legacy Logger middleware (keeping for compatibility)
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)
		status := c.Writer.Status()
		log.Printf("%s %s %d %s", c.Request.Method, c.Request.URL.Path, status, latency)
	}
}

// StructuredLoggingMiddleware provides structured logging for all requests
func StructuredLoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get client IP
		clientIP := getClientIP(c)

		// Get user agent
		userAgent := c.Request.UserAgent()

		// Get status code
		statusCode := c.Writer.Status()

		// Get request size
		requestSize := c.Request.ContentLength

		// Get response size
		responseSize := c.Writer.Size()

		// Create structured log entry
		logEntry := GetLogger().With(
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", raw),
			zap.String("ip", clientIP),
			zap.String("user_agent", userAgent),
			zap.Int("status", statusCode),
			zap.Duration("latency", latency),
			zap.Int64("request_size", requestSize),
			zap.Int("response_size", responseSize),
			zap.String("request_id", getRequestID(c)),
		)

		// Log based on status code
		switch {
		case statusCode >= 500:
			logEntry.Error("Server error")
		case statusCode >= 400:
			logEntry.Warn("Client error")
		default:
			logEntry.Info("Request completed")
		}
	}
}

// PerformanceMonitoringMiddleware tracks performance metrics
func PerformanceMonitoringMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Add performance tracking to context
		ctx := context.WithValue(c.Request.Context(), "start_time", start)
		c.Request = c.Request.WithContext(ctx)

		// Process request
		c.Next()

		// Calculate metrics
		latency := time.Since(start)

		// Log performance metrics
		GetLogger().Info("Performance metrics",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Duration("latency", latency),
			zap.Int("status", c.Writer.Status()),
			zap.String("request_id", getRequestID(c)),
		)

		// Alert on slow requests (over 5 seconds)
		if latency > 5*time.Second {
			GetLogger().Warn("Slow request detected",
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.Duration("latency", latency),
				zap.String("request_id", getRequestID(c)),
			)
		}
	}
}

// ErrorLoggingMiddleware logs errors with context
func ErrorLoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Process request
		c.Next()

		// Check for errors
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				GetLogger().Error("Request error",
					zap.String("method", c.Request.Method),
					zap.String("path", c.Request.URL.Path),
					zap.String("error", err.Error()),
					zap.String("request_id", getRequestID(c)),
					zap.String("client_ip", getClientIP(c)),
					zap.String("user_agent", c.Request.UserAgent()),
				)
			}
		}
	}
}

// HealthCheckMiddleware provides health check endpoint
func HealthCheckMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Basic health check
		health := map[string]interface{}{
			"status":    "healthy",
			"timestamp": time.Now().UTC(),
			"version":   "1.0.0", // This should come from config
		}

		// Add additional health checks here
		// - Database connectivity
		// - External service health
		// - Memory usage
		// - Disk space

		c.JSON(http.StatusOK, health)
	}
}

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if request ID already exists
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			// Generate new request ID
			requestID = generateRequestID()
		}

		// Add to context
		c.Set("request_id", requestID)

		// Add to response headers
		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}

// getRequestID retrieves request ID from context
func getRequestID(c *gin.Context) string {
	if requestID, exists := c.Get("request_id"); exists {
		return requestID.(string)
	}
	return "unknown"
}

// generateRequestID creates a unique request ID
func generateRequestID() string {
	return fmt.Sprintf("req_%d", time.Now().UnixNano())
}

// DatabaseHealthCheck checks database connectivity
func DatabaseHealthCheck(db interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Convert interface{} to *gorm.DB
		gormDB, ok := db.(*gorm.DB)
		if !ok {
			c.JSON(http.StatusServiceUnavailable, map[string]interface{}{
				"database":  "unhealthy",
				"error":     "invalid database connection",
				"timestamp": time.Now().UTC(),
			})
			return
		}

		// Get underlying sql.DB
		sqlDB, err := gormDB.DB()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, map[string]interface{}{
				"database":  "unhealthy",
				"error":     "failed to get database connection: " + err.Error(),
				"timestamp": time.Now().UTC(),
			})
			return
		}

		// Test database connectivity with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := sqlDB.PingContext(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, map[string]interface{}{
				"database":  "unhealthy",
				"error":     "database ping failed: " + err.Error(),
				"timestamp": time.Now().UTC(),
			})
			return
		}

		// Get database stats
		stats := sqlDB.Stats()
		health := map[string]interface{}{
			"database":      "healthy",
			"timestamp":     time.Now().UTC(),
			"max_open":      stats.MaxOpenConnections,
			"open":          stats.OpenConnections,
			"in_use":        stats.InUse,
			"idle":          stats.Idle,
			"wait_count":    stats.WaitCount,
			"wait_duration": stats.WaitDuration.String(),
		}

		c.JSON(http.StatusOK, health)
	}
}

// MetricsMiddleware collects basic metrics
func MetricsMiddleware() gin.HandlerFunc {
	// Simple in-memory metrics (in production, use Prometheus or similar)
	var (
		requestCount int64
		errorCount   int64
		responseTime time.Duration
		lastRequest  time.Time
	)

	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		// Update metrics
		requestCount++
		if c.Writer.Status() >= 400 {
			errorCount++
		}
		responseTime = time.Since(start)
		lastRequest = time.Now()

		// Log metrics periodically
		if requestCount%100 == 0 {
			GetLogger().Info("Metrics snapshot",
				zap.Int64("total_requests", requestCount),
				zap.Int64("error_count", errorCount),
				zap.Duration("avg_response_time", responseTime),
				zap.Time("last_request", lastRequest),
			)
		}
	}
}

// InfoLog logs info level messages
func InfoLog(format string, args ...interface{}) {
	GetLogger().Sugar().Infof(format, args...)
}

// WarnLog logs warning level messages
func WarnLog(format string, args ...interface{}) {
	GetLogger().Sugar().Warnf(format, args...)
}

// ErrorLog logs error level messages
func ErrorLog(format string, args ...interface{}) {
	GetLogger().Sugar().Errorf(format, args...)
}

// FatalLog logs fatal level messages and exits
func FatalLog(format string, args ...interface{}) {
	GetLogger().Sugar().Fatalf(format, args...)
}
