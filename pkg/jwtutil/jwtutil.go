package jwtutil

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(getSecret()) // Use the same secret as AAA service

// getSecret returns the JWT secret with fallback
func getSecret() string {
	secret := os.Getenv("SECRET_KEY")
	if secret == "" {
		return "secret" // Fallback secret
	}
	return secret
}

// GenerateToken creates a JWT token using user info and expiration duration
func GenerateToken(userID, email, role string, duration time.Duration) (string, error) {
	log.Printf("🔑 === GENERATE TOKEN START ===")
	log.Printf("🔑 UserID: %s", userID)
	log.Printf("🔑 Email: %s", email)
	log.Printf("🔑 Role: %s", role)
	log.Printf("🔑 Duration: %v", duration)
	log.Printf("🔑 JWT Secret length: %d", len(jwtSecret))

	// Convert single role to roles array for consistency with middleware
	roles := []string{role}

	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"role":    role,  // Keep single role for backward compatibility
		"roles":   roles, // Add roles array for middleware
		"exp":     time.Now().Add(duration).Unix(),
	}

	log.Printf("🔑 Claims: %+v", claims)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	log.Printf("🔑 Token created with signing method: %s", jwt.SigningMethodHS256)

	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		log.Printf("❌ Failed to sign token: %v", err)
		return "", err
	}

	log.Printf("✅ Token generated successfully, length: %d", len(signedToken))
	log.Printf("🔑 === GENERATE TOKEN COMPLETE ===")
	return signedToken, nil
}

// ParseToken validates and parses the JWT token
func ParseToken(tokenStr string) (jwt.MapClaims, error) {
	log.Printf("🔑 === PARSE TOKEN START ===")
	log.Printf("🔑 Token string length: %d", len(tokenStr))
	log.Printf("🔑 Token preview: %s...", tokenStr[:min(50, len(tokenStr))])
	log.Printf("🔑 JWT Secret length: %d", len(jwtSecret))

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		log.Printf("🔑 Parsing token with method: %s", token.Method.Alg())
		return jwtSecret, nil // Use shared secret
	})

	if err != nil {
		log.Printf("❌ Failed to parse JWT token: %v", err)
		return nil, err
	}

	if !token.Valid {
		log.Printf("❌ Token is not valid")
		return nil, jwt.ErrTokenMalformed
	}

	log.Printf("✅ Token is valid")

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Printf("❌ Failed to extract claims from token")
		return nil, jwt.ErrTokenMalformed
	}

	log.Printf("✅ Claims extracted successfully: %+v", claims)
	log.Printf("🔑 === PARSE TOKEN COMPLETE ===")
	return claims, nil
}

// Helper function for min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
