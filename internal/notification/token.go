package notification

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"
)

// GenerateUnsubscribeToken generates a cryptographically secure random token
// Returns: token (plain), tokenHash (for storage), expiry (90 days from now)
func GenerateUnsubscribeToken() (token string, tokenHash string, expiry time.Time, err error) {
	// Generate 32 bytes of random data (256 bits of entropy)
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", "", time.Time{}, fmt.Errorf("failed to generate secure token: %w", err)
	}
	
	// Use URL-safe base64 encoding (no padding) for clean URLs
	token = base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(b)
	
	// Hash token before storing (SHA-256)
	hash := sha256.Sum256([]byte(token))
	tokenHash = hex.EncodeToString(hash[:])
	
	// Set expiry to 90 days from now
	expiry = time.Now().Add(90 * 24 * time.Hour)
	
	return token, tokenHash, expiry, nil
}

// HashToken hashes a token using SHA-256 for comparison
func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

