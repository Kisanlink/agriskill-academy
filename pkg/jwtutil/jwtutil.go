package jwtutil

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte
var jwtSecretOnce sync.Once

// getSecret returns the JWT secret with fallback
func getSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT_SECRET environment variable is required")
	}
	return secret
}

// getJWTSecret returns the JWT secret, initializing it once
func getJWTSecret() []byte {
	jwtSecretOnce.Do(func() {
		jwtSecret = []byte(getSecret())
	})
	return jwtSecret
}

func DebugLog(format string, args ...interface{}) {
	if os.Getenv("GIN_MODE") == "debug" {
		log.Printf("[DEBUG] "+format, args...)
	}
}

// GenerateToken creates a JWT token using user info and expiration duration
func GenerateToken(userID, username, email, role string, duration time.Duration) (string, error) {
	DebugLog("🔑 === GENERATE TOKEN START ===")
	DebugLog("🔑 UserID: %s", userID)
	DebugLog("🔑 Username: %s", username)
	DebugLog("🔑 Email: %s", email)
	DebugLog("🔑 Role: %s", role)
	DebugLog("🔑 Duration: %v", duration)

	secret := getJWTSecret()
	DebugLog("🔑 JWT Secret length: %d", len(secret))

	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"email":    email,
		"role":     role,
		"exp":      time.Now().Add(duration).Unix(),
	}

	DebugLog("🔑 Claims: %+v", claims)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	DebugLog("🔑 Token created with signing method: %s", jwt.SigningMethodHS256)

	signedToken, err := token.SignedString(secret)
	if err != nil {
		DebugLog("❌ Failed to sign token: %v", err)
		return "", err
	}

	DebugLog("✅ Token generated successfully, length: %d", len(signedToken))
	DebugLog("🔑 === GENERATE TOKEN COMPLETE ===")
	return signedToken, nil
}

// ParseToken validates and parses the JWT token
func ParseToken(tokenStr string) (jwt.MapClaims, error) {
	DebugLog("🔑 === PARSE TOKEN START ===")
	DebugLog("🔑 Token string length: %d", len(tokenStr))
	DebugLog("🔑 Token preview: %s...", tokenStr[:min(50, len(tokenStr))])

	secret := getJWTSecret()
	DebugLog("🔑 JWT Secret length: %d", len(secret))

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		DebugLog("🔑 Parsing token with method: %s", token.Method.Alg())
		return secret, nil // Use shared secret
	})

	if err != nil {
		DebugLog("❌ Failed to parse JWT token: %v", err)
		return nil, err
	}

	if !token.Valid {
		DebugLog("❌ Token is not valid")
		return nil, jwt.ErrTokenMalformed
	}

	DebugLog("✅ Token is valid")

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		DebugLog("❌ Failed to extract claims from token")
		return nil, jwt.ErrTokenMalformed
	}

	DebugLog("✅ Claims extracted successfully: %+v", claims)
	DebugLog("🔑 === PARSE TOKEN COMPLETE ===")
	return claims, nil
}

// Helper function for min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
