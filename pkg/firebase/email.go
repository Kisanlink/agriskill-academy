// File: pkg/firebase/email.go

package firebase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

// EmailService handles sending emails via Firebase
type EmailService struct {
	client      *auth.Client
	app         *firebase.App
	apiKey      string
	frontendURL string
}

// NewEmailService initializes Firebase Admin SDK for email sending
func NewEmailService(credentialsPath, credentialsJSON, apiKey, frontendURL string) (*EmailService, error) {
	ctx := context.Background()

	var opt option.ClientOption

	// Priority: JSON credentials (for containers) > file path (for local dev)
	if credentialsJSON != "" {
		log.Println("Initializing Firebase with base64 credentials from environment")
		opt = option.WithCredentialsJSON([]byte(credentialsJSON))
	} else if credentialsPath != "" {
		log.Printf("Initializing Firebase with credentials file: %s", credentialsPath)
		if _, err := os.Stat(credentialsPath); os.IsNotExist(err) {
			return nil, fmt.Errorf("firebase credentials file not found: %s", credentialsPath)
		}
		opt = option.WithCredentialsFile(credentialsPath)
	} else {
		return nil, fmt.Errorf("firebase credentials not provided: set FIREBASE_CREDENTIALS_PATH or FIREBASE_CREDENTIALS_JSON")
	}

	// Initialize Firebase App
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Firebase app: %w", err)
	}

	// Get Auth client
	client, err := app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get Firebase Auth client: %w", err)
	}

	if apiKey == "" {
		return nil, fmt.Errorf("firebase Web API key not provided: set FIREBASE_WEB_API_KEY")
	}

	log.Println("✅ Firebase email service initialized successfully")

	return &EmailService{
		client:      client,
		app:         app,
		apiKey:      apiKey,
		frontendURL: frontendURL,
	}, nil
}

// SendVerificationEmail sends an email verification email to the user via Firebase REST API
func (s *EmailService) SendVerificationEmail(ctx context.Context, email, token string) error {
	log.Printf("📧 Sending verification email to: %s", email)

	// First, get or create Firebase user
	userRecord, err := s.client.GetUserByEmail(ctx, email)
	if err != nil {
		// User doesn't exist in Firebase, create one
		params := (&auth.UserToCreate{}).Email(email).EmailVerified(false)
		userRecord, err = s.client.CreateUser(ctx, params)
		if err != nil {
			return fmt.Errorf("failed to create Firebase user for email: %w", err)
		}
		log.Printf("Created Firebase user for email: %s", email)
	}

	// Generate a custom token for the user
	customToken, err := s.client.CustomToken(ctx, userRecord.UID)
	if err != nil {
		return fmt.Errorf("failed to create custom token: %w", err)
	}

	// Exchange custom token for ID token using Firebase REST API
	idToken, err := s.exchangeCustomToken(customToken)
	if err != nil {
		return fmt.Errorf("failed to exchange custom token: %w", err)
	}

	// Send verification email using Firebase REST API
	err = s.sendOobCode(idToken, "VERIFY_EMAIL")
	if err != nil {
		return fmt.Errorf("failed to send verification email: %w", err)
	}

	log.Printf("✅ Verification email sent successfully to: %s", email)
	return nil
}

// SendPasswordResetEmail sends a password reset email to the user via Firebase REST API
func (s *EmailService) SendPasswordResetEmail(ctx context.Context, email, token string) error {
	log.Printf("📧 Sending password reset email to: %s", email)

	// Get or create Firebase user
	_, err := s.client.GetUserByEmail(ctx, email)
	if err != nil {
		// User doesn't exist in Firebase, create one
		params := (&auth.UserToCreate{}).Email(email).EmailVerified(true)
		_, err = s.client.CreateUser(ctx, params)
		if err != nil {
			return fmt.Errorf("failed to create Firebase user for password reset: %w", err)
		}
		log.Printf("Created Firebase user for password reset: %s", email)
	}

	// Send password reset email directly using Firebase REST API
	err = s.sendPasswordResetByEmail(email)
	if err != nil {
		return fmt.Errorf("failed to send password reset email: %w", err)
	}

	log.Printf("✅ Password reset email sent successfully to: %s", email)
	return nil
}

// exchangeCustomToken exchanges a custom token for an ID token
func (s *EmailService) exchangeCustomToken(customToken string) (string, error) {
	url := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/accounts:signInWithCustomToken?key=%s", s.apiKey)

	payload := map[string]interface{}{
		"token":             customToken,
		"returnSecureToken": true,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("firebase API error (%d): %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	idToken, ok := result["idToken"].(string)
	if !ok {
		return "", fmt.Errorf("idToken not found in response")
	}

	return idToken, nil
}

// sendOobCode sends an out-of-band code (verification email)
func (s *EmailService) sendOobCode(idToken, requestType string) error {
	url := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/accounts:sendOobCode?key=%s", s.apiKey)

	payload := map[string]interface{}{
		"requestType": requestType,
		"idToken":     idToken,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("firebase API error (%d): %s", resp.StatusCode, string(body))
	}

	log.Printf("📧 Firebase sendOobCode response: %s", string(body))
	return nil
}

// sendPasswordResetByEmail sends a password reset email using Firebase REST API
func (s *EmailService) sendPasswordResetByEmail(email string) error {
	url := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/accounts:sendOobCode?key=%s", s.apiKey)

	payload := map[string]interface{}{
		"requestType": "PASSWORD_RESET",
		"email":       email,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("firebase API error (%d): %s", resp.StatusCode, string(body))
	}

	log.Printf("📧 Firebase password reset response: %s", string(body))
	return nil
}

// DeleteFirebaseUser deletes a Firebase user (cleanup helper)
func (s *EmailService) DeleteFirebaseUser(ctx context.Context, email string) error {
	userRecord, err := s.client.GetUserByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("user not found in Firebase: %w", err)
	}

	if err := s.client.DeleteUser(ctx, userRecord.UID); err != nil {
		return fmt.Errorf("failed to delete Firebase user: %w", err)
	}

	log.Printf("Deleted Firebase user: %s", email)
	return nil
}
