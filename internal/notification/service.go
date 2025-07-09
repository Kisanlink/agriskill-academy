// File: internal/notification/service.go

package notification

import (
	"os"

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
	return &mailService{
		From:      os.Getenv("MAIL_FROM"),
		SMTPHost:  os.Getenv("MAIL_HOST"),
		SMTPPort:  587,
		Password:  os.Getenv("MAIL_PASS"),
		prefsRepo: prefsRepo,
	}
}

func (s *mailService) SendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(s.SMTPHost, s.SMTPPort, s.From, s.Password)
	return d.DialAndSend(m)
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
