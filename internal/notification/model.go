package notification

import (
	"time"

	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Notification Preferences Models
type NotificationPreferences struct {
	ID                 string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID             string    `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`
	EmailNotifications bool      `json:"email_notifications" gorm:"default:true"`
	PushNotifications  bool      `json:"push_notifications" gorm:"default:true"`
	JobAlerts          bool      `json:"job_alerts" gorm:"default:true"`
	ApplicationUpdates bool      `json:"application_updates" gorm:"default:true"`
	CompanyNews        bool      `json:"company_news" gorm:"default:false"`
	MarketingEmails    bool      `json:"marketing_emails" gorm:"default:false"`
	WeeklyDigest       bool      `json:"weekly_digest" gorm:"default:true"`
	DailyJobMatches    bool      `json:"daily_job_matches" gorm:"default:false"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// BeforeCreate is a GORM hook that generates UUID for ID if it's empty and validates if not empty
func (n *NotificationPreferences) BeforeCreate(tx *gorm.DB) error {
	if n.ID == "" {
		n.ID = uuid.New().String()
	} else {
		if _, err := uuid.Parse(n.ID); err != nil {
			return fmt.Errorf("invalid UUID format for NotificationPreferences ID: %w", err)
		}
	}
	return nil
}

// Request/Response Models
type UpdateNotificationPreferencesRequest struct {
	EmailNotifications *bool `json:"email_notifications"`
	PushNotifications  *bool `json:"push_notifications"`
	JobAlerts          *bool `json:"job_alerts"`
	ApplicationUpdates *bool `json:"application_updates"`
	CompanyNews        *bool `json:"company_news"`
	MarketingEmails    *bool `json:"marketing_emails"`
	WeeklyDigest       *bool `json:"weekly_digest"`
	DailyJobMatches    *bool `json:"daily_job_matches"`
}

type NotificationPreferencesResponse struct {
	Success bool                     `json:"success"`
	Message string                   `json:"message"`
	Data    *NotificationPreferences `json:"data,omitempty"`
}

// Email Templates
type EmailNotification struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

// Notification Types
const (
	NotificationTypeJobAlert          = "job_alert"
	NotificationTypeApplicationUpdate = "application_update"
	NotificationTypeCompanyNews       = "company_news"
	NotificationTypeMarketing         = "marketing"
	NotificationTypeWeeklyDigest      = "weekly_digest"
	NotificationTypeDailyMatches      = "daily_matches"
)
