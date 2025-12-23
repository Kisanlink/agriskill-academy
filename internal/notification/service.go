// File: internal/notification/service.go

package notification

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Kisanlink/agriskill-academy/internal/middleware"
	"gopkg.in/gomail.v2"
)

type NotificationService interface {
	SendEmail(to, subject, body string) error
	GetPreferences(userID string) (*NotificationPreferences, error)
	UpdatePreferences(userID string, req *UpdateNotificationPreferencesRequest) (*NotificationPreferences, error)
	ShouldSendNotification(userID, notificationType string) (bool, error)
	SendNotificationIfEnabled(userID, notificationType, to, subject, body string) error
	GenerateUnsubscribeToken(userID, notificationType string) (string, error)
	ProcessUnsubscribe(token string) (string, error)
	GetUnsubscribeURL(userID, notificationType string) (string, string, error)
	GetPreferencesByToken(token string) (*NotificationPreferences, error)
	UpdatePreferencesByToken(token string, emailNotifications, pushNotifications, jobAlerts, applicationUpdates bool) error
	GetManagePreferencesURL(userID, notificationType, token string) (string, error)
}

type mailService struct {
	From      string
	SMTPHost  string
	SMTPPort  int
	Password  string
	Timeout   time.Duration
	enabled   bool
	prefsRepo NotificationPreferencesRepository
}

func NewMailService() NotificationService {
	return &mailService{
		From:     os.Getenv("MAIL_FROM"),
		SMTPHost: os.Getenv("MAIL_HOST"),
		SMTPPort: 587, // or os.Getenv("MAIL_PORT")
		Password: os.Getenv("MAIL_PASS"),
	}
}

func NewNotificationService(prefsRepo NotificationPreferencesRepository) NotificationService {
	// Check if email notifications are enabled
	emailEnabled := os.Getenv("EMAIL_NOTIFICATION") == "true"

	if !emailEnabled {
		log.Println("📧 Email notifications disabled (EMAIL_NOTIFICATION=false)")
		return &mailService{
			prefsRepo: prefsRepo,
			enabled:   false,
		}
	}

	// Validate SMTP configuration
	mailFrom := os.Getenv("MAIL_FROM")
	mailHost := os.Getenv("MAIL_HOST")
	mailPass := os.Getenv("MAIL_PASS")

	if mailFrom == "" || mailHost == "" || mailPass == "" {
		log.Println("⚠️  Email notifications enabled but SMTP not fully configured")
		log.Println("⚠️  Missing one or more: MAIL_FROM, MAIL_HOST, MAIL_PASS")
		log.Println("⚠️  Email notifications will be disabled")
		return &mailService{
			prefsRepo: prefsRepo,
			enabled:   false,
		}
	}

	// Parse SMTP port with default
	mailPort := 587
	if mailPortStr := os.Getenv("MAIL_PORT"); mailPortStr != "" {
		if parsed, err := strconv.Atoi(mailPortStr); err == nil {
			mailPort = parsed
		} else {
			log.Printf("⚠️  Invalid MAIL_PORT value '%s', using default: 587", mailPortStr)
		}
	}

	log.Printf("✅ Email notifications enabled: %s:%d", mailHost, mailPort)

	return &mailService{
		From:      mailFrom,
		SMTPHost:  mailHost,
		SMTPPort:  mailPort,
		Password:  mailPass,
		Timeout:   10 * time.Second,
		enabled:   true,
		prefsRepo: prefsRepo,
	}
}

func (s *mailService) SendEmail(to, subject, body string) error {
	// Check if email notifications are enabled
	if !s.enabled {
		log.Printf("⚠️  Skipping email to %s - SMTP not configured", to)
		return nil // Don't fail, just skip silently
	}

	// Validate inputs
	if to == "" {
		return errors.New("email recipient is required")
	}
	if subject == "" {
		return errors.New("email subject is required")
	}

	// Create message
	m := gomail.NewMessage()
	m.SetHeader("From", s.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	// Create dialer
	// Note: gomail v2 uses default timeouts internally
	d := gomail.NewDialer(s.SMTPHost, s.SMTPPort, s.From, s.Password)

	// Retry logic: 3 attempts with exponential backoff
	var lastErr error
	for attempt := 1; attempt <= 3; attempt++ {
		err := d.DialAndSend(m)
		if err == nil {
			log.Printf("✅ Email sent successfully to: %s (attempt %d)", to, attempt)
			return nil
		}

		lastErr = err
		if attempt < 3 {
			backoff := time.Duration(attempt) * time.Second
			log.Printf("⚠️  Email send failed (attempt %d/3), retrying in %v: %v", attempt, backoff, err)
			time.Sleep(backoff)
		}
	}

	log.Printf("❌ Failed to send email to %s after 3 attempts: %v", to, lastErr)
	return fmt.Errorf("failed to send email after 3 attempts: %w", lastErr)
}

func (s *mailService) GetPreferences(userID string) (*NotificationPreferences, error) {
	return s.prefsRepo.GetOrCreate(userID)
}

func (s *mailService) UpdatePreferences(userID string, req *UpdateNotificationPreferencesRequest) (*NotificationPreferences, error) {
	preferences, err := s.prefsRepo.GetOrCreate(userID)
	if err != nil {
		return nil, err
	}

	// Update only provided fields
	if req.EmailNotifications != nil {
		preferences.EmailNotifications = *req.EmailNotifications
	}
	if req.PushNotifications != nil {
		preferences.PushNotifications = *req.PushNotifications
	}
	if req.JobAlerts != nil {
		preferences.JobAlerts = *req.JobAlerts
	}
	if req.ApplicationUpdates != nil {
		preferences.ApplicationUpdates = *req.ApplicationUpdates
	}

	err = s.prefsRepo.Update(preferences)
	if err != nil {
		return nil, err
	}

	return preferences, nil
}

func (s *mailService) ShouldSendNotification(userID, notificationType string) (bool, error) {
	preferences, err := s.prefsRepo.GetOrCreate(userID)
	if err != nil {
		return false, err
	}

	// Check master switch
	if !preferences.EmailNotifications {
		return false, nil
	}

	// Check specific type
	switch notificationType {
	case NotificationTypeJobAlert:
		return preferences.JobAlerts, nil
	case NotificationTypeApplicationUpdate:
		return preferences.ApplicationUpdates, nil
	default:
		return false, nil // Fail-safe: don't send unknown types
	}
}

func (s *mailService) SendNotificationIfEnabled(userID, notificationType, to, subject, body string) error {
	shouldSend, err := s.ShouldSendNotification(userID, notificationType)
	if err != nil {
		return err
	}

	if shouldSend {
		return s.SendEmail(to, subject, body)
	}

	return nil // Notification disabled, don't send
}

// GenerateUnsubscribeToken generates and stores an unsubscribe token for a user
func (s *mailService) GenerateUnsubscribeToken(userID, notificationType string) (string, error) {
	// Validate notification type
	if notificationType != NotificationTypeJobAlert && notificationType != NotificationTypeApplicationUpdate {
		return "", fmt.Errorf("invalid notification type: %s", notificationType)
	}

	token, tokenHash, expiry, err := GenerateUnsubscribeToken()
	if err != nil {
		return "", err
	}

	_, err = s.prefsRepo.GetOrCreate(userID)
	if err != nil {
		return "", err
	}

	// Update type-specific token
	err = s.prefsRepo.UpdateUnsubscribeToken(userID, notificationType, tokenHash, expiry)
	if err != nil {
		return "", err
	}

	return token, nil
}

// ProcessUnsubscribe processes an unsubscribe request using a token
func (s *mailService) ProcessUnsubscribe(token string) (string, error) {
	middleware.DebugLog("🔍 ProcessUnsubscribe - received token: %s (length: %d)", token, len(token))
	
	tokenHash := HashToken(token)

	// Log full hash for debugging
	middleware.DebugLog("🔍 Processing unsubscribe - full token hash: %s", tokenHash)
	
	// Log token processing (show first 16 chars of hash for debugging)
	hashPreview := tokenHash
	if len(tokenHash) > 16 {
		hashPreview = tokenHash[:16] + "..."
	}
	middleware.DebugLog("🔍 Processing unsubscribe - token hash preview: %s", hashPreview)

	// Get preferences and determine type
	preferences, notificationType, err := s.prefsRepo.GetByUnsubscribeToken(tokenHash)
	if err != nil {
		middleware.DebugLog("❌ Token lookup failed: %v", err)
		middleware.DebugLog("❌ Looking for hash: %s", tokenHash)
		return "", fmt.Errorf("invalid or expired token")
	}

	middleware.DebugLog("📋 Found preferences for user: %s, notification type: %s", preferences.UserID, notificationType)

	// Disable only specific type
	err = s.prefsRepo.DisableNotification(preferences.UserID, notificationType)
	if err != nil {
		middleware.DebugLog("❌ Failed to update preferences: %v", err)
		return "", err
	}

	middleware.DebugLog("✅ Preferences updated successfully - %s disabled for user: %s", notificationType, preferences.UserID)
	return notificationType, nil
}

// GetUnsubscribeURL generates an unsubscribe URL for a user
// Returns both the URL and the token so it can be reused for manage preferences
func (s *mailService) GetUnsubscribeURL(userID, notificationType string) (string, string, error) {
	// Generate type-specific token
	token, err := s.GenerateUnsubscribeToken(userID, notificationType)
	if err != nil {
		return "", "", err
	}

	baseURL := os.Getenv("ASA_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	baseURL = strings.TrimSuffix(baseURL, "/")

	unsubscribeURL := fmt.Sprintf("%s/api/notifications/unsubscribe/%s?type=%s",
		baseURL, token, notificationType)

	return unsubscribeURL, token, nil
}

// GetPreferencesByToken retrieves preferences by token without disabling anything
func (s *mailService) GetPreferencesByToken(token string) (*NotificationPreferences, error) {
	tokenHash := HashToken(token)
	
	// Get preferences and determine type (we don't need the type, just the preferences)
	preferences, _, err := s.prefsRepo.GetByUnsubscribeToken(tokenHash)
	if err != nil {
		return nil, fmt.Errorf("invalid or expired token")
	}
	
	return preferences, nil
}

// UpdatePreferencesByToken updates preferences using a token
func (s *mailService) UpdatePreferencesByToken(token string, emailNotifications, pushNotifications, jobAlerts, applicationUpdates bool) error {
	tokenHash := HashToken(token)
	
	// Get preferences by token to find user ID
	preferences, _, err := s.prefsRepo.GetByUnsubscribeToken(tokenHash)
	if err != nil {
		return fmt.Errorf("invalid or expired token")
	}
	
	// Update all preferences
	updates := map[string]interface{}{
		"email_notifications": emailNotifications,
		"push_notifications":  pushNotifications,
		"job_alerts":          jobAlerts,
		"application_updates": applicationUpdates,
	}
	
	err = s.prefsRepo.UpdatePreferencesByUserID(preferences.UserID, updates)
	
	if err != nil {
		return fmt.Errorf("failed to update preferences: %w", err)
	}
	
	middleware.DebugLog("✅ Preferences updated via token for user: %s", preferences.UserID)
	return nil
}

// GetManagePreferencesURL generates a manage preferences URL using an existing token
// If token is empty, it will generate a new one
func (s *mailService) GetManagePreferencesURL(userID, notificationType, token string) (string, error) {
	var err error
	
	// If no token provided, generate one
	if token == "" {
		token, err = s.GenerateUnsubscribeToken(userID, notificationType)
		if err != nil {
			return "", err
		}
	}

	baseURL := os.Getenv("ASA_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	baseURL = strings.TrimSuffix(baseURL, "/")

	manageURL := fmt.Sprintf("%s/api/notifications/manage/%s?type=%s",
		baseURL, token, notificationType)

	return manageURL, nil
}
