// File: pkg/firebase/auth_adapter.go

package firebase

import (
	"fmt"
)

// AuthAdapter adapts the AuthClient to match the interface expected by auth service
type AuthAdapter struct {
	client *AuthClient
}

// NewAuthAdapter creates a new adapter for Firebase authentication
func NewAuthAdapter(apiKey string) *AuthAdapter {
	return &AuthAdapter{
		client: NewAuthClient(apiKey),
	}
}

// SignUpWithPassword creates a new user in Firebase and returns the Firebase UID
func (a *AuthAdapter) SignUpWithPassword(email, password string) (string, error) {
	response, err := a.client.SignUpWithPassword(email, password)
	if err != nil {
		return "", fmt.Errorf("firebase signup failed: %w", err)
	}
	return response.LocalID, nil
}

// SignInWithPassword authenticates a user and returns Firebase UID and email verification status
func (a *AuthAdapter) SignInWithPassword(email, password string) (string, bool, error) {
	response, err := a.client.SignInWithPassword(email, password)
	if err != nil {
		return "", false, fmt.Errorf("firebase signin failed: %w", err)
	}
	return response.LocalID, response.EmailVerified, nil
}

// UpdatePassword updates a user's password in Firebase
// Note: This requires idToken, which we don't have in password reset flow
// This method is not currently used - we use EmailService.UpdateFirebasePassword instead
func (a *AuthAdapter) UpdatePassword(email, newPassword string) error {
	// This would require getting an idToken first, which we can't do without the old password
	// Instead, we use the Firebase Admin SDK via EmailService.UpdateFirebasePassword
	return fmt.Errorf("UpdatePassword via REST API not implemented - use EmailService.UpdateFirebasePassword instead")
}
