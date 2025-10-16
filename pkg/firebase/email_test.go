// File: pkg/firebase/email_test.go

package firebase

import (
	"testing"
)

// TestNewEmailService tests Firebase email service initialization
func TestNewEmailService(t *testing.T) {
	// Test case 1: Missing credentials
	_, err := NewEmailService("", "", "http://localhost:3000")
	if err == nil {
		t.Error("Expected error when no credentials provided, got nil")
	}

	// Test case 2: Invalid credentials path
	_, err = NewEmailService("/nonexistent/path/to/credentials.json", "", "http://localhost:3000")
	if err == nil {
		t.Error("Expected error for nonexistent credentials file, got nil")
	}

	// Note: Real Firebase integration tests require actual credentials
	// and should be run separately with proper test fixtures
}

// TestEmailServiceConfiguration tests configuration validation
func TestEmailServiceConfiguration(t *testing.T) {
	tests := []struct {
		name            string
		credentialsPath string
		credentialsJSON string
		frontendURL     string
		expectError     bool
	}{
		{
			name:            "No credentials provided",
			credentialsPath: "",
			credentialsJSON: "",
			frontendURL:     "http://localhost:3000",
			expectError:     true,
		},
		{
			name:            "Only path provided (invalid path)",
			credentialsPath: "/invalid/path.json",
			credentialsJSON: "",
			frontendURL:     "http://localhost:3000",
			expectError:     true,
		},
		{
			name:            "Only JSON provided (invalid JSON)",
			credentialsPath: "",
			credentialsJSON: "invalid-json",
			frontendURL:     "http://localhost:3000",
			expectError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewEmailService(tt.credentialsPath, tt.credentialsJSON, tt.frontendURL)
			if (err != nil) != tt.expectError {
				t.Errorf("Expected error: %v, got: %v", tt.expectError, err != nil)
			}
		})
	}
}

// Note: Integration tests for SendVerificationEmail and SendPasswordResetEmail
// require actual Firebase credentials and should be run in a separate test suite
// with proper mocking or real Firebase test project setup.
//
// Example integration test structure:
//
// func TestSendVerificationEmailIntegration(t *testing.T) {
//     if testing.Short() {
//         t.Skip("Skipping integration test")
//     }
//
//     credPath := os.Getenv("FIREBASE_CREDENTIALS_PATH")
//     if credPath == "" {
//         t.Skip("FIREBASE_CREDENTIALS_PATH not set")
//     }
//
//     service, err := NewEmailService(credPath, "", "http://localhost:3000")
//     if err != nil {
//         t.Fatalf("Failed to create email service: %v", err)
//     }
//
//     err = service.SendVerificationEmail(context.Background(), "test@example.com", "test-token")
//     if err != nil {
//         t.Errorf("Failed to send verification email: %v", err)
//     }
// }
