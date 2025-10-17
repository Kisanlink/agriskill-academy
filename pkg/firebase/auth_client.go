// File: pkg/firebase/auth_client.go

package firebase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// AuthClient handles Firebase Authentication via REST API
type AuthClient struct {
	apiKey     string
	httpClient *http.Client
}

// SignUpResponse represents the response from Firebase signUp endpoint
type SignUpResponse struct {
	IDToken      string `json:"idToken"`
	Email        string `json:"email"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    string `json:"expiresIn"`
	LocalID      string `json:"localId"` // This is the Firebase UID
}

// SignInResponse represents the response from Firebase signIn endpoint
type SignInResponse struct {
	IDToken      string `json:"idToken"`
	Email        string `json:"email"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    string `json:"expiresIn"`
	LocalID      string `json:"localId"`      // Firebase UID
	Registered   bool   `json:"registered"`   // True if user exists
	DisplayName  string `json:"displayName"`  // User's display name if set
	EmailVerified bool  `json:"emailVerified"` // Email verification status
}

// UpdatePasswordResponse represents the response from Firebase password update
type UpdatePasswordResponse struct {
	LocalID      string `json:"localId"`
	Email        string `json:"email"`
	IDToken      string `json:"idToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    string `json:"expiresIn"`
}

// FirebaseError represents an error response from Firebase
type FirebaseError struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Errors  []struct {
			Message string `json:"message"`
			Domain  string `json:"domain"`
			Reason  string `json:"reason"`
		} `json:"errors"`
	} `json:"error"`
}

// NewAuthClient creates a new Firebase Authentication client
func NewAuthClient(apiKey string) *AuthClient {
	return &AuthClient{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SignUpWithPassword creates a new Firebase user with email and password
// Firebase will automatically send a verification email
func (c *AuthClient) SignUpWithPassword(email, password string) (*SignUpResponse, error) {
	url := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/accounts:signUp?key=%s", c.apiKey)

	payload := map[string]interface{}{
		"email":             email,
		"password":          password,
		"returnSecureToken": true,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal signup request: %w", err)
	}

	resp, err := c.httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to make signup request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read signup response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, c.parseFirebaseError(body, resp.StatusCode)
	}

	var result SignUpResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse signup response: %w", err)
	}

	log.Printf("✅ Firebase user created successfully: %s (UID: %s)", email, result.LocalID)
	return &result, nil
}

// SignInWithPassword authenticates a user with email and password
func (c *AuthClient) SignInWithPassword(email, password string) (*SignInResponse, error) {
	url := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=%s", c.apiKey)

	payload := map[string]interface{}{
		"email":             email,
		"password":          password,
		"returnSecureToken": true,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal signin request: %w", err)
	}

	resp, err := c.httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to make signin request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read signin response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, c.parseFirebaseError(body, resp.StatusCode)
	}

	var result SignInResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse signin response: %w", err)
	}

	log.Printf("✅ Firebase authentication successful: %s (UID: %s, Verified: %v)",
		email, result.LocalID, result.EmailVerified)
	return &result, nil
}

// UpdatePassword updates the password for a user (requires idToken)
func (c *AuthClient) UpdatePassword(idToken, newPassword string) (*UpdatePasswordResponse, error) {
	url := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/accounts:update?key=%s", c.apiKey)

	payload := map[string]interface{}{
		"idToken":           idToken,
		"password":          newPassword,
		"returnSecureToken": true,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal update password request: %w", err)
	}

	resp, err := c.httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to make update password request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read update password response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, c.parseFirebaseError(body, resp.StatusCode)
	}

	var result UpdatePasswordResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse update password response: %w", err)
	}

	log.Printf("✅ Firebase password updated successfully for UID: %s", result.LocalID)
	return &result, nil
}

// parseFirebaseError parses Firebase error responses into readable error messages
func (c *AuthClient) parseFirebaseError(body []byte, statusCode int) error {
	var fbErr FirebaseError
	if err := json.Unmarshal(body, &fbErr); err != nil {
		return fmt.Errorf("firebase API error (%d): %s", statusCode, string(body))
	}

	// Map common Firebase error messages to user-friendly messages
	errMsg := fbErr.Error.Message
	switch errMsg {
	case "EMAIL_EXISTS":
		return fmt.Errorf("email already exists in Firebase")
	case "INVALID_EMAIL":
		return fmt.Errorf("invalid email address")
	case "WEAK_PASSWORD : Password should be at least 6 characters":
		return fmt.Errorf("password must be at least 6 characters")
	case "EMAIL_NOT_FOUND":
		return fmt.Errorf("user not found in Firebase")
	case "INVALID_PASSWORD":
		return fmt.Errorf("invalid password")
	case "USER_DISABLED":
		return fmt.Errorf("user account has been disabled")
	case "INVALID_ID_TOKEN":
		return fmt.Errorf("invalid authentication token")
	default:
		return fmt.Errorf("firebase error (%d): %s", statusCode, errMsg)
	}
}
