// File: internal/userprofile/repository.go

package userprofile

import (
	"gorm.io/gorm"
)

type UserProfileRepository interface {
	GetByUserID(userID string) (*UserProfile, error)
	Update(profile *UserProfile) error
	Create(profile *UserProfile) error
}

type userProfileRepository struct {
	db *gorm.DB
}

func NewUserProfileRepository(db *gorm.DB) UserProfileRepository {
	return &userProfileRepository{db}
}

func (r *userProfileRepository) GetByUserID(userID string) (*UserProfile, error) {
	var profile UserProfile
	err := r.db.Preload("Certificates").Where("user_id = ?", userID).First(&profile).Error
	return &profile, err
}

func (r *userProfileRepository) Update(profile *UserProfile) error {
	return r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(profile).Error
}

func (r *userProfileRepository) Create(profile *UserProfile) error {
	return r.db.Create(profile).Error
}
