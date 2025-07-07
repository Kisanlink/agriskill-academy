// File: internal/userprofile/repository.go

package userprofile

import (
	"gorm.io/gorm"
)

type UserProfileRepository interface {
	GetByUserID(userID string) (*UserProfile, error)
	Update(profile *UserProfile) error
	Create(profile *UserProfile) error
	AddCertificate(cert *Certificate) error
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
	// Start a transaction
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Delete existing certificates for this profile
	if err := tx.Where("user_profile_id = ?", profile.ID).Delete(&Certificate{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Update the profile
	if err := tx.Save(profile).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Create new certificates
	for i := range profile.Certificates {
		profile.Certificates[i].UserProfileID = profile.ID
		// Ensure ID is empty so database generates proper UUID
		profile.Certificates[i].ID = ""
		if err := tx.Create(&profile.Certificates[i]).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit the transaction
	return tx.Commit().Error
}

func (r *userProfileRepository) Create(profile *UserProfile) error {
	return r.db.Create(profile).Error
}

func (r *userProfileRepository) AddCertificate(cert *Certificate) error {
	// Ensure ID is empty so database generates proper UUID
	cert.ID = ""
	return r.db.Create(cert).Error
}
