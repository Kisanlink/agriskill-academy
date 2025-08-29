package middleware

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// SecurityConfig holds security-related configuration
type SecurityConfig struct {
	RateLimitRequests int           `json:"rate_limit_requests"`
	RateLimitWindow   time.Duration `json:"rate_limit_window"`
	AllowedOrigins    []string      `json:"allowed_origins"`
	MaxRequestSize    int64         `json:"max_request_size"`
	EnableCORS        bool          `json:"enable_cors"`
}

// RateLimiter implements IP-based rate limiting
type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	r        rate.Limit
	b        int
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		r:        r,
		b:        b,
	}
}

// GetLimiter returns the rate limiter for the given key
func (rl *RateLimiter) GetLimiter(key string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.limiters[key]
	if !exists {
		limiter = rate.NewLimiter(rl.r, rl.b)
		rl.limiters[key] = limiter
	}

	return limiter
}

// Cleanup removes old limiters to prevent memory leaks
func (rl *RateLimiter) Cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// In a production environment, you might want to implement
	// a more sophisticated cleanup strategy based on usage patterns
	if len(rl.limiters) > 10000 { // Arbitrary limit
		rl.limiters = make(map[string]*rate.Limiter)
	}
}

// RateLimitMiddleware implements IP-based rate limiting
func RateLimitMiddleware(requests int, window time.Duration) gin.HandlerFunc {
	limiter := NewRateLimiter(rate.Limit(float64(requests)/window.Seconds()), requests)

	// Start cleanup goroutine
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			limiter.Cleanup()
		}
	}()

	return func(c *gin.Context) {
		key := getClientIP(c)
		if !limiter.GetLimiter(key).Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "Rate limit exceeded",
				"message": "Too many requests, please try again later",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// getClientIP extracts the real client IP from various headers
func getClientIP(c *gin.Context) string {
	// Check for forwarded headers (common with proxies)
	if ip := c.GetHeader("X-Forwarded-For"); ip != "" {
		// X-Forwarded-For can contain multiple IPs, take the first one
		if commaIndex := strings.Index(ip, ","); commaIndex != -1 {
			return strings.TrimSpace(ip[:commaIndex])
		}
		return strings.TrimSpace(ip)
	}

	if ip := c.GetHeader("X-Real-IP"); ip != "" {
		return strings.TrimSpace(ip)
	}

	if ip := c.GetHeader("X-Client-IP"); ip != "" {
		return strings.TrimSpace(ip)
	}

	// Fallback to remote address
	return c.ClientIP()
}

// InputSanitizationMiddleware sanitizes input to prevent XSS and injection attacks
func InputSanitizationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Sanitize query parameters
		for key, values := range c.Request.URL.Query() {
			for i, value := range values {
				values[i] = sanitizeInput(value)
			}
			c.Request.URL.Query()[key] = values
		}

		// Sanitize form data
		if err := c.Request.ParseForm(); err == nil {
			for key, values := range c.Request.PostForm {
				for i, value := range values {
					values[i] = sanitizeInput(value)
				}
				c.Request.PostForm[key] = values
			}
		}

		c.Next()
	}
}

// sanitizeInput removes potentially dangerous characters and patterns
func sanitizeInput(input string) string {
	// Remove null bytes
	input = strings.ReplaceAll(input, "\x00", "")

	// Remove control characters (except newlines and tabs)
	re := regexp.MustCompile(`[\x00-\x08\x0B\x0C\x0E-\x1F\x7F]`)
	input = re.ReplaceAllString(input, "")

	// Basic XSS prevention - encode common dangerous patterns
	dangerousPatterns := map[string]string{
		"<script":     "&lt;script",
		"javascript:": "javascript&#58;",
		"onload=":     "onload&#61;",
		"onerror=":    "onerror&#61;",
		"onclick=":    "onclick&#61;",
	}

	for pattern, replacement := range dangerousPatterns {
		input = strings.ReplaceAll(strings.ToLower(input), pattern, replacement)
	}

	return strings.TrimSpace(input)
}

// RequestSizeLimitMiddleware limits the size of incoming requests
func RequestSizeLimitMiddleware(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)

		c.Next()

		// Check if the request was too large
		if c.Request.ContentLength > maxSize {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"error":   "Request too large",
				"message": fmt.Sprintf("Request size exceeds maximum allowed size of %d bytes", maxSize),
			})
			c.Abort()
			return
		}
	}
}

// CORSValidationMiddleware validates CORS configuration
func CORSValidationMiddleware(allowedOrigins []string) gin.HandlerFunc {
	originMap := make(map[string]bool)
	for _, origin := range allowedOrigins {
		originMap[origin] = true
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		// Allow requests without Origin header (same-origin requests)
		if origin == "" {
			c.Next()
			return
		}

		// Check if origin is allowed
		if !originMap[origin] && !originMap["*"] {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "CORS policy violation",
				"message": "Origin not allowed",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// SecurityHeadersMiddleware adds security-related HTTP headers
func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Prevent clickjacking
		c.Header("X-Frame-Options", "DENY")

		// Prevent MIME type sniffing
		c.Header("X-Content-Type-Options", "nosniff")

		// Enable XSS protection
		c.Header("X-XSS-Protection", "1; mode=block")

		// Strict transport security (HTTPS only)
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		// Content security policy
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' https:; connect-src 'self' https:;")

		// Referrer policy
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// Permissions policy
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		c.Next()
	}
}

// SQLInjectionProtectionMiddleware provides basic SQL injection protection
func SQLInjectionProtectionMiddleware() gin.HandlerFunc {
	// Common SQL injection patterns
	sqlPatterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)(union|select|insert|update|delete|drop|create|alter|exec|execute)`),
		regexp.MustCompile(`(?i)(--|/\*|\*/|xp_|sp_)`),
		regexp.MustCompile(`(?i)(or\s+\d+\s*=\s*\d+|and\s+\d+\s*=\s*\d+)`),
		regexp.MustCompile(`(?i)(union\s+select|select\s+union)`),
	}

	return func(c *gin.Context) {
		// Check query parameters
		for _, values := range c.Request.URL.Query() {
			for _, value := range values {
				if containsSQLInjection(value, sqlPatterns) {
					c.JSON(http.StatusBadRequest, gin.H{
						"error":   "Invalid input detected",
						"message": "Request contains potentially malicious content",
					})
					c.Abort()
					return
				}
			}
		}

		// Check form data
		if err := c.Request.ParseForm(); err == nil {
			for _, values := range c.Request.PostForm {
				for _, value := range values {
					if containsSQLInjection(value, sqlPatterns) {
						c.JSON(http.StatusBadRequest, gin.H{
							"error":   "Invalid input detected",
							"message": "Request contains potentially malicious content",
						})
						c.Abort()
						return
					}
				}
			}
		}

		c.Next()
	}
}

// containsSQLInjection checks if input contains SQL injection patterns
func containsSQLInjection(input string, patterns []*regexp.Regexp) bool {
	for _, pattern := range patterns {
		if pattern.MatchString(input) {
			return true
		}
	}
	return false
}

// FileUploadSecurityMiddleware validates file uploads
func FileUploadSecurityMiddleware(allowedTypes []string, maxFileSize int64) gin.HandlerFunc {
	allowedTypesMap := make(map[string]bool)
	for _, t := range allowedTypes {
		allowedTypesMap[strings.ToLower(t)] = true
	}

	return func(c *gin.Context) {
		// Check file size
		if c.Request.ContentLength > maxFileSize {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"error":   "File too large",
				"message": fmt.Sprintf("File size exceeds maximum allowed size of %d bytes", maxFileSize),
			})
			c.Abort()
			return
		}

		// Check content type
		contentType := c.GetHeader("Content-Type")
		if !strings.HasPrefix(contentType, "multipart/form-data") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid content type",
				"message": "Only multipart/form-data is allowed for file uploads",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// ContextTimeoutMiddleware adds timeout to request context
func ContextTimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)

		// Create a channel to signal completion
		done := make(chan struct{})

		go func() {
			c.Next()
			close(done)
		}()

		select {
		case <-done:
			// Request completed successfully
		case <-ctx.Done():
			c.JSON(http.StatusRequestTimeout, gin.H{
				"error":   "Request timeout",
				"message": "Request took too long to process",
			})
			c.Abort()
		}
	}
}
