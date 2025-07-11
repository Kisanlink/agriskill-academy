// File: internal/studentprofile/repository.go

package studentprofile

import (
	"gorm.io/gorm"
)

type StudentProfileRepository interface {
	GetByUserID(userID string) (*StudentProfile, error)
	Update(profile *StudentProfile) error
	Create(profile *StudentProfile) error
	AddCertificate(cert *Certificate) error
	DeleteCertificate(certID string, userID string) error
}

type studentProfileRepository struct {
	db *gorm.DB
}

func NewStudentProfileRepository(db *gorm.DB) StudentProfileRepository {
	return &studentProfileRepository{db}
}

func (r *studentProfileRepository) GetByUserID(userID string) (*StudentProfile, error) {
	var profile StudentProfile
	err := r.db.Preload("Certificates").Where("user_id = ?", userID).First(&profile).Error
	return &profile, err
}

func (r *studentProfileRepository) Update(profile *StudentProfile) error {
	// Start a transaction
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Delete existing certificates for this profile
	if err := tx.Where("student_profile_id = ?", profile.ID).Delete(&Certificate{}).Error; err != nil {
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
		profile.Certificates[i].StudentProfileID = profile.ID
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

func (r *studentProfileRepository) Create(profile *StudentProfile) error {
	return r.db.Create(profile).Error
}

func (r *studentProfileRepository) AddCertificate(cert *Certificate) error {
	// Ensure ID is empty so database generates proper UUID
	cert.ID = ""
	return r.db.Create(cert).Error
}

func (r *studentProfileRepository) DeleteCertificate(certID string, userID string) error {
	// First verify that the certificate belongs to the user
	var cert Certificate
	err := r.db.Joins("JOIN student_profiles ON certificates.student_profile_id = student_profiles.id").
		Where("certificates.id = ? AND student_profiles.user_id = ?", certID, userID).
		First(&cert).Error

	if err != nil {
		return err
	}

	// Delete the certificate
	return r.db.Delete(&cert).Error
}
