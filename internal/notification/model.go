package notification

import (
	"asa/internal/middleware"

	"github.com/Kisanlink/kisanlink-db/pkg/base"
	"github.com/Kisanlink/kisanlink-db/pkg/core/hash"
	"gorm.io/gorm"
)

// Notification Preferences Models
type NotificationPreferences struct {
	base.BaseModel
	UserID             string `gorm:"type:varchar(255);not null;uniqueIndex" json:"user_id"`
	EmailNotifications bool   `json:"email_notifications" gorm:"default:true"`
	PushNotifications  bool   `json:"push_notifications" gorm:"default:true"`
	JobAlerts          bool   `json:"job_alerts" gorm:"default:true"`
	ApplicationUpdates bool   `json:"application_updates" gorm:"default:true"`
	CompanyNews        bool   `json:"company_news" gorm:"default:false"`
	MarketingEmails    bool   `json:"marketing_emails" gorm:"default:false"`
	WeeklyDigest       bool   `json:"weekly_digest" gorm:"default:true"`
	DailyJobMatches    bool   `json:"daily_job_matches" gorm:"default:false"`
}

// TableName specifies the database table name for NotificationPreferences
func (NotificationPreferences) TableName() string {
	return "notification_preferences"
}

// NewNotificationPreferences creates a new NotificationPreferences with proper initialization
func NewNotificationPreferences() *NotificationPreferences {
	return &NotificationPreferences{
		BaseModel: *base.NewBaseModel("NOTP", hash.Small),
	}
}

func InitializeCounterFromDatabase(db *gorm.DB) {
	var notificationIDs []string
	if err := db.Model(&NotificationPreferences{}).Pluck("id", &notificationIDs).Error; err == nil {
		hash.InitializeGlobalCountersFromDatabase("NOTP", notificationIDs, hash.Small)
		middleware.DebugLog("Initialized NOTP counter with %d existing IDs", len(notificationIDs))
	}
}

// BeforeCreateGORM is called by GORM before creating a new record
func (n *NotificationPreferences) BeforeCreateGORM(tx *gorm.DB) error {
	return n.BeforeCreate()
}

// BeforeUpdateGORM is called by GORM before updating an existing record
func (n *NotificationPreferences) BeforeUpdateGORM(tx *gorm.DB) error {
	return n.BeforeUpdate()
}

// BeforeDeleteGORM is called by GORM before hard deleting a record
func (n *NotificationPreferences) BeforeDeleteGORM(tx *gorm.DB) error {
	return n.BeforeDelete()
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
