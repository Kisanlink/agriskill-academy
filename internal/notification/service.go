// File: internal/notification/service.go

package notification

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"gopkg.in/gomail.v2"
)

type NotificationService interface {
	SendEmail(to, subject, body string) error
	GetPreferences(userID string) (*NotificationPreferences, error)
	UpdatePreferences(userID string, req *UpdateNotificationPreferencesRequest) (*NotificationPreferences, error)
	ShouldSendNotification(userID, notificationType string) (bool, error)
	SendNotificationIfEnabled(userID, notificationType, to, subject, body string) error
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

	// Update only the fields that are provided
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
	if req.CompanyNews != nil {
		preferences.CompanyNews = *req.CompanyNews
	}
	if req.MarketingEmails != nil {
		preferences.MarketingEmails = *req.MarketingEmails
	}
	if req.WeeklyDigest != nil {
		preferences.WeeklyDigest = *req.WeeklyDigest
	}
	if req.DailyJobMatches != nil {
		preferences.DailyJobMatches = *req.DailyJobMatches
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

	// Check if email notifications are enabled
	if !preferences.EmailNotifications {
		return false, nil
	}

	// Check specific notification type
	switch notificationType {
	case NotificationTypeJobAlert:
		return preferences.JobAlerts, nil
	case NotificationTypeApplicationUpdate:
		return preferences.ApplicationUpdates, nil
	case NotificationTypeCompanyNews:
		return preferences.CompanyNews, nil
	case NotificationTypeMarketing:
		return preferences.MarketingEmails, nil
	case NotificationTypeWeeklyDigest:
		return preferences.WeeklyDigest, nil
	case NotificationTypeDailyMatches:
		return preferences.DailyJobMatches, nil
	default:
		return true, nil // Default to sending if type not specified
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
