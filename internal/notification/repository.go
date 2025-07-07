package notification

import (
	"errors"

	"gorm.io/gorm"
)

type NotificationPreferencesRepository interface {
	GetByUserID(userID string) (*NotificationPreferences, error)
	Create(preferences *NotificationPreferences) error
	Update(preferences *NotificationPreferences) error
	GetOrCreate(userID string) (*NotificationPreferences, error)
}

type notificationPreferencesRepository struct {
	db *gorm.DB
}

func NewNotificationPreferencesRepository(db *gorm.DB) NotificationPreferencesRepository {
	return &notificationPreferencesRepository{db}
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
		preferences = &NotificationPreferences{
			UserID:             userID,
			EmailNotifications: true,
			PushNotifications:  true,
			JobAlerts:          true,
			ApplicationUpdates: true,
			CompanyNews:        false,
			MarketingEmails:    false,
			WeeklyDigest:       true,
			DailyJobMatches:    false,
		}
		err = r.Create(preferences)
		if err != nil {
			return nil, err
		}
	}
	return preferences, nil
}
