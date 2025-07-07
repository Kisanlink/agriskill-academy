package notification

import (
	"time"
)

// Notification Preferences Models
type NotificationPreferences struct {
	ID                 string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID             string    `gorm:"type:uuid;not null;uniqueIndex" json:"userId"`
	EmailNotifications bool      `json:"emailNotifications" gorm:"default:true"`
	PushNotifications  bool      `json:"pushNotifications" gorm:"default:true"`
	JobAlerts          bool      `json:"jobAlerts" gorm:"default:true"`
	ApplicationUpdates bool      `json:"applicationUpdates" gorm:"default:true"`
	CompanyNews        bool      `json:"companyNews" gorm:"default:false"`
	MarketingEmails    bool      `json:"marketingEmails" gorm:"default:false"`
	WeeklyDigest       bool      `json:"weeklyDigest" gorm:"default:true"`
	DailyJobMatches    bool      `json:"dailyJobMatches" gorm:"default:false"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
}

// Request/Response Models
type UpdatePreferencesRequest struct {
	EmailNotifications *bool `json:"emailNotifications"`
	PushNotifications  *bool `json:"pushNotifications"`
	JobAlerts          *bool `json:"jobAlerts"`
	ApplicationUpdates *bool `json:"applicationUpdates"`
	CompanyNews        *bool `json:"companyNews"`
	MarketingEmails    *bool `json:"marketingEmails"`
	WeeklyDigest       *bool `json:"weeklyDigest"`
	DailyJobMatches    *bool `json:"dailyJobMatches"`
}

type PreferencesResponse struct {
	Success bool                     `json:"success"`
	Message string                   `json:"message"`
	Data    *NotificationPreferences `json:"data,omitempty"`
}

// Email Templates
type EmailTemplate struct {
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
