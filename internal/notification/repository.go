package notification

import (
	"errors"

	"github.com/Kisanlink/kisanlink-db/pkg/base"
	"gorm.io/gorm"
)

type NotificationPreferencesRepository interface {
	GetByUserID(userID string) (*NotificationPreferences, error)
	Create(preferences *NotificationPreferences) error
	Update(preferences *NotificationPreferences) error
	GetOrCreate(userID string) (*NotificationPreferences, error)
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
		// Create default preferences if not found
		preferences = NewNotificationPreferences() // Use constructor to generate ID
		preferences.UserID = userID
		preferences.EmailNotifications = true
		preferences.PushNotifications = true
		preferences.JobAlerts = true
		preferences.ApplicationUpdates = true
		preferences.CompanyNews = false
		preferences.MarketingEmails = false
		preferences.WeeklyDigest = true
		preferences.DailyJobMatches = false
		err = r.Create(preferences)
		if err != nil {
			return nil, err
		}
	}
	return preferences, nil
}
