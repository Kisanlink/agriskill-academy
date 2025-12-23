package notification

import (
	"fmt"
	"os"
	"strings"

	"github.com/Kisanlink/agriskill-academy/internal/middleware"
	"gorm.io/gorm"
)

// JobEnqueuer is a minimal interface to avoid import cycles
// It matches the Enqueue method from worker.JobService
type JobEnqueuer interface {
	Enqueue(job interface{}) error
}

// JobEnqueuerFunc is a function type that implements JobEnqueuer
type JobEnqueuerFunc func(job interface{}) error

func (f JobEnqueuerFunc) Enqueue(job interface{}) error {
	return f(job)
}

// HandleSendEmail processes email sending jobs from the worker queue
// This function should be called by the worker when processing "send_email" job types
func HandleSendEmail(payload map[string]interface{}, notificationSvc NotificationService) error {
	// Extract payload data
	to, ok := payload["to"].(string)
	if !ok || to == "" {
		return fmt.Errorf("invalid or missing 'to' field in email job payload")
	}

	subject, ok := payload["subject"].(string)
	if !ok || subject == "" {
		return fmt.Errorf("invalid or missing 'subject' field in email job payload")
	}

	htmlContent, ok := payload["html"].(string)
	if !ok || htmlContent == "" {
		return fmt.Errorf("invalid or missing 'html' field in email job payload")
	}

	emailType := "email"
	if et, ok := payload["type"].(string); ok {
		emailType = et
	}

	middleware.DebugLog("📤 Processing email job: type=%s, to=%s", emailType, to)

	// Send email via SMTP (existing NotificationService)
	err := notificationSvc.SendEmail(to, subject, htmlContent)

	if err != nil {
		middleware.DebugLog("❌ Failed to send email to %s: %v", to, err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	middleware.DebugLog("✅ Successfully sent %s email to %s", emailType, to)
	return nil
}

type EmailSenderService struct {
	notificationSvc NotificationService
	jobService      JobEnqueuer
	templateSvc     *EmailTemplateService
	db              *gorm.DB
}

func NewEmailSenderService(notifSvc NotificationService, jobSvc JobEnqueuer, db *gorm.DB) *EmailSenderService {
	return &EmailSenderService{
		notificationSvc: notifSvc,
		jobService:     jobSvc,
		templateSvc:    NewEmailTemplateService(),
		db:             db,
	}
}

// getUserIDByEmail looks up user ID by email address
func (s *EmailSenderService) getUserIDByEmail(email string) (string, error) {
	var user struct {
		ID string
	}
	err := s.db.Table("users").Where("email = ?", email).Select("id").First(&user).Error
	if err != nil {
		return "", err
	}
	return user.ID, nil
}

func (s *EmailSenderService) SendNewJobEmail(studentEmail string, jobData map[string]interface{}) error {
	jobTitle := ""
	if title, ok := jobData["JobTitle"].(string); ok {
		jobTitle = title
	}

	middleware.DebugLog("📧 Sending new job email to: %s, Job: %s", studentEmail, jobTitle)

	// Get logo URL from environment and add to job data
	logoURL := os.Getenv("EMAIL_LOGO_URL")
	if logoURL != "" {
		jobData["LogoURL"] = logoURL
	}

	// Get user ID and generate type-specific unsubscribe URL
	userID, err := s.getUserIDByEmail(studentEmail)
	if err == nil {
		unsubscribeURL, err := s.notificationSvc.GetUnsubscribeURL(userID, NotificationTypeJobAlert)
		if err == nil {
			jobData["UnsubscribeURL"] = unsubscribeURL
		}
	}

	// Get base URL for manage preferences
	baseURL := os.Getenv("ASA_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	baseURL = strings.TrimSuffix(baseURL, "/")
	jobData["ManagePreferencesURL"] = fmt.Sprintf("%s/notifications/preferences", baseURL)

	// Render email HTML
	htmlContent, err := s.templateSvc.RenderNewJobEmail(jobData)
	if err != nil {
		middleware.DebugLog("❌ Failed to render new job email: %v", err)
		return fmt.Errorf("failed to render email template: %w", err)
	}

	subject := fmt.Sprintf("New Job Posted: %s", jobTitle)

	// Queue background job to send email
	emailJobData := map[string]interface{}{
		"to":      studentEmail,
		"subject": subject,
		"html":    htmlContent,
		"type":    "new_job",
	}

	// Create job struct - we need to use reflection or type assertion
	// Since we can't import worker, we'll create a struct that matches BackgroundJob
	type backgroundJob struct {
		Type    string
		Payload map[string]interface{}
	}
	job := &backgroundJob{
		Type:    "send_email",
		Payload: emailJobData,
	}

	if err := s.jobService.Enqueue(job); err != nil {
		middleware.DebugLog("❌ Failed to queue new job email: %v", err)
		return fmt.Errorf("failed to queue email job: %w", err)
	}

	middleware.DebugLog("✅ Queued new job email for: %s", studentEmail)
	return nil
}

func (s *EmailSenderService) SendStatusUpdateEmail(studentEmail string, appData map[string]interface{}) error {
	status := ""
	if st, ok := appData["Status"].(string); ok {
		status = st
	}

	middleware.DebugLog("📧 Sending status update email to: %s, Status: %s", studentEmail, status)

	// Get logo URL from environment and add to app data
	logoURL := os.Getenv("EMAIL_LOGO_URL")
	if logoURL != "" {
		appData["LogoURL"] = logoURL
	}

	// Get user ID and generate type-specific unsubscribe URL
	userID, err := s.getUserIDByEmail(studentEmail)
	if err == nil {
		unsubscribeURL, err := s.notificationSvc.GetUnsubscribeURL(userID, NotificationTypeApplicationUpdate)
		if err == nil {
			appData["UnsubscribeURL"] = unsubscribeURL
		}
	}

	// Get base URL for manage preferences
	baseURL := os.Getenv("ASA_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	baseURL = strings.TrimSuffix(baseURL, "/")
	appData["ManagePreferencesURL"] = fmt.Sprintf("%s/notifications/preferences", baseURL)

	// Render email HTML
	htmlContent, err := s.templateSvc.RenderStatusUpdateEmail(appData)
	if err != nil {
		middleware.DebugLog("❌ Failed to render status update email: %v", err)
		return fmt.Errorf("failed to render email template: %w", err)
	}

	subject := fmt.Sprintf("Application Update: %s", status)

	// Queue background job to send email
	emailJobData := map[string]interface{}{
		"to":      studentEmail,
		"subject": subject,
		"html":    htmlContent,
		"type":    "status_update",
	}

	// Create job struct - we need to use reflection or type assertion
	// Since we can't import worker, we'll create a struct that matches BackgroundJob
	type backgroundJob struct {
		Type    string
		Payload map[string]interface{}
	}
	job := &backgroundJob{
		Type:    "send_email",
		Payload: emailJobData,
	}

	if err := s.jobService.Enqueue(job); err != nil {
		middleware.DebugLog("❌ Failed to queue status update email: %v", err)
		return fmt.Errorf("failed to queue email job: %w", err)
	}

	middleware.DebugLog("✅ Queued status update email for: %s", studentEmail)
	return nil
}

