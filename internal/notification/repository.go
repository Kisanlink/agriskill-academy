package notification

import (
	"errors"
	"fmt"
	"time"

	"github.com/Kisanlink/agriskill-academy/internal/middleware"
	"github.com/Kisanlink/kisanlink-db/pkg/base"
	"gorm.io/gorm"
)

type NotificationPreferencesRepository interface {
	GetByUserID(userID string) (*NotificationPreferences, error)
	Create(preferences *NotificationPreferences) error
	Update(preferences *NotificationPreferences) error
	GetOrCreate(userID string) (*NotificationPreferences, error)
	// Returns preferences, notification type, and error
	GetByUnsubscribeToken(tokenHash string) (*NotificationPreferences, string, error)
	// Now takes notification type parameter
	UpdateUnsubscribeToken(userID, notificationType, tokenHash string, expiry time.Time) error
	// DisableNotification disables a specific notification type for a user
	DisableNotification(userID, notificationType string) error
}

type notificationPreferencesRepository struct {
	*base.BaseRepository[*NotificationPreferences]
	db *gorm.DB
}

func NewNotificationPreferencesRepository(db *gorm.DB) NotificationPreferencesRepository {
	return &notificationPreferencesRepository{
		BaseRepository: base.NewBaseRepository[*NotificationPreferences](),
		db:             db,
	}
}

func (r *notificationPreferencesRepository) GetByUserID(userID string) (*NotificationPreferences, error) {
	var preferences NotificationPreferences
	err := r.db.Where("user_id = ?", userID).First(&preferences).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("notification preferences not found")
		}
		return nil, err
	}
	return &preferences, nil
}

func (r *notificationPreferencesRepository) Create(preferences *NotificationPreferences) error {
	return r.db.Create(preferences).Error
}

func (r *notificationPreferencesRepository) Update(preferences *NotificationPreferences) error {
	return r.db.Save(preferences).Error
}

func (r *notificationPreferencesRepository) GetOrCreate(userID string) (*NotificationPreferences, error) {
	preferences, err := r.GetByUserID(userID)
	if err != nil {
		// Create with all enabled by default
		preferences = NewNotificationPreferences()
		preferences.UserID = userID
		preferences.EmailNotifications = true
		preferences.PushNotifications = true
		preferences.JobAlerts = true
		preferences.ApplicationUpdates = true
		err = r.Create(preferences)
		if err != nil {
			return nil, err
		}
	}
	return preferences, nil
}

// GetByUnsubscribeToken retrieves preferences by unsubscribe token hash
func (r *notificationPreferencesRepository) GetByUnsubscribeToken(tokenHash string) (*NotificationPreferences, string, error) {
	middleware.DebugLog("🔍 GetByUnsubscribeToken - searching for hash: %s", tokenHash)
	
	var preferences NotificationPreferences

	// Check job_alerts token
	err := r.db.Where(
		"job_alerts_unsubscribe_token_hash = ? AND job_alerts_unsubscribe_token_expiry > ?",
		tokenHash, time.Now(),
	).First(&preferences).Error

	if err == nil {
		middleware.DebugLog("✅ Found job_alerts token match for user: %s", preferences.UserID)
		return &preferences, NotificationTypeJobAlert, nil
	}
	middleware.DebugLog("ℹ️  No job_alerts token match: %v", err)

	// Check application_updates token
	err = r.db.Where(
		"application_updates_unsubscribe_token_hash = ? AND application_updates_unsubscribe_token_expiry > ?",
		tokenHash, time.Now(),
	).First(&preferences).Error

	if err == nil {
		middleware.DebugLog("✅ Found application_updates token match for user: %s", preferences.UserID)
		return &preferences, NotificationTypeApplicationUpdate, nil
	}
	middleware.DebugLog("ℹ️  No application_updates token match: %v", err)

	return nil, "", fmt.Errorf("invalid or expired token")
}

// UpdateUnsubscribeToken updates the unsubscribe token for a user
func (r *notificationPreferencesRepository) UpdateUnsubscribeToken(userID, notificationType, tokenHash string, expiry time.Time) error {
	var hashColumn, expiryColumn string

	switch notificationType {
	case NotificationTypeJobAlert:
		hashColumn = "job_alerts_unsubscribe_token_hash"
		expiryColumn = "job_alerts_unsubscribe_token_expiry"
	case NotificationTypeApplicationUpdate:
		hashColumn = "application_updates_unsubscribe_token_hash"
		expiryColumn = "application_updates_unsubscribe_token_expiry"
	default:
		return fmt.Errorf("invalid notification type: %s", notificationType)
	}

	return r.db.Model(&NotificationPreferences{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			hashColumn:   tokenHash,
			expiryColumn: expiry,
		}).Error
}

// DisableNotification disables a specific notification type for a user
func (r *notificationPreferencesRepository) DisableNotification(userID, notificationType string) error {
	updates := map[string]interface{}{}

	switch notificationType {
	case NotificationTypeJobAlert:
		updates["job_alerts"] = false
		updates["job_alerts_unsubscribe_token_hash"] = ""
		updates["job_alerts_unsubscribe_token_expiry"] = nil
	case NotificationTypeApplicationUpdate:
		updates["application_updates"] = false
		updates["application_updates_unsubscribe_token_hash"] = ""
		updates["application_updates_unsubscribe_token_expiry"] = nil
	default:
		return fmt.Errorf("invalid notification type: %s", notificationType)
	}

	return r.db.Model(&NotificationPreferences{}).Where("user_id = ?", userID).Updates(updates).Error
}