// File: internal/studentprofile/repository.go

package studentprofile

import (
	"asa/internal/middleware"

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
	middleware.DebugLog("🔍 DEBUG: Repository.GetByUserID called for userID: %s\n", userID)

	var profile StudentProfile
	err := r.db.Preload("Certificates").Where("user_id = ?", userID).First(&profile).Error

	if err != nil {
		middleware.DebugLog("❌ DEBUG: Repository.GetByUserID error: %v\n", err)
		return &profile, err
	}

	middleware.DebugLog("✅ DEBUG: Repository.GetByUserID success: %+v\n", profile)
	return &profile, err
}

func (r *studentProfileRepository) Update(profile *StudentProfile) error {
	middleware.DebugLog("🔍 DEBUG: Repository.Update called - ID: %s, Name: %s, ProfilePhotoSize: %d, ResumeSize: %d\n", profile.ID, profile.Name, profile.ProfilePhotoSize, profile.ResumeSize)

	// Check if profile exists
	var existingProfile StudentProfile
	err := r.db.Where("id = ?", profile.ID).First(&existingProfile).Error
	if err != nil {
		middleware.DebugLog("❌ DEBUG: Repository.Update - Profile not found: %v\n", err)
		return err
	}

	middleware.DebugLog("✅ DEBUG: Repository.Update - Profile found, updating...\n")

	result := r.db.Save(profile)
	if result.Error != nil {
		middleware.DebugLog("❌ DEBUG: Repository.Update error: %v\n", result.Error)
		return result.Error
	}

	middleware.DebugLog("✅ DEBUG: Repository.Update completed successfully - Rows affected: %d\n", result.RowsAffected)
	return nil
}

func (r *studentProfileRepository) Create(profile *StudentProfile) error {
	middleware.DebugLog("🔍 DEBUG: Repository.Create called - UserID: %s, Name: %s, ProfilePhotoSize: %d, ResumeSize: %d\n", profile.UserID, profile.Name, profile.ProfilePhotoSize, profile.ResumeSize)

	result := r.db.Create(profile)
	if result.Error != nil {
		middleware.DebugLog("❌ DEBUG: Repository.Create error: %v\n", result.Error)
		return result.Error
	}

	middleware.DebugLog("✅ DEBUG: Repository.Create completed successfully - New ID: %s, Rows affected: %d\n", profile.ID, result.RowsAffected)
	return nil
}

func (r *studentProfileRepository) AddCertificate(cert *Certificate) error {
	middleware.DebugLog("🔍 DEBUG: Repository.AddCertificate called - Name: %s, FileSize: %d bytes, FileType: %s\n", cert.Name, cert.FileSize, cert.FileType)

	// Ensure ID is empty so database generates proper UUID
	cert.ID = ""
	err := r.db.Create(cert).Error
	if err != nil {
		middleware.DebugLog("❌ DEBUG: Repository.AddCertificate error: %v\n", err)
		return err
	}

	middleware.DebugLog("✅ DEBUG: Repository.AddCertificate success, new certificate ID: %s\n", cert.ID)
	return nil
}

func (r *studentProfileRepository) DeleteCertificate(certID string, userID string) error {
	middleware.DebugLog("🔍 DEBUG: Repository.DeleteCertificate called for certID: %s, userID: %s\n", certID, userID)

	// First verify that the certificate belongs to the user
	var cert Certificate
	err := r.db.Joins("JOIN student_profiles ON certificates.student_profile_id = student_profiles.id").
		Where("certificates.id = ? AND student_profiles.user_id = ?", certID, userID).
		First(&cert).Error

	if err != nil {
		middleware.DebugLog("❌ DEBUG: Repository.DeleteCertificate - certificate not found or access denied: %v\n", err)
		return err
	}
	middleware.DebugLog("✅ DEBUG: Certificate found and verified\n")

	// Delete the certificate
	err = r.db.Delete(&cert).Error
	if err != nil {
		middleware.DebugLog("❌ DEBUG: Repository.DeleteCertificate - failed to delete: %v\n", err)
		return err
	}

	middleware.DebugLog("✅ DEBUG: Repository.DeleteCertificate success\n")
	return nil
}
