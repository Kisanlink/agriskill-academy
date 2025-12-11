// File: internal/studentprofile/repository.go

package studentprofile

import (
	"github.com/Kisanlink/agriskill-academy/internal/middleware"
	"context"

	"github.com/Kisanlink/kisanlink-db/pkg/base"
	"gorm.io/gorm"
)

type StudentProfileRepository interface {
	base.Repository[*StudentProfile]
	GetByUserID(userID string) (*StudentProfile, error)
	AddCertificate(cert *Certificate) error
	DeleteCertificate(certID string, userID string) error
}

type studentProfileRepository struct {
	*base.BaseRepository[*StudentProfile]
	db *gorm.DB
}

func NewStudentProfileRepository(db *gorm.DB) StudentProfileRepository {
	return &studentProfileRepository{
		BaseRepository: base.NewBaseRepository[*StudentProfile](),
		db:             db,
	}
}

func (r *studentProfileRepository) Create(ctx context.Context, profile *StudentProfile) error {
	middleware.DebugLog("🔍 DEBUG: Repository.Create called - UserID: %s, Name: %s, Skills: %v\n", profile.UserID, profile.Name, profile.Skills)

	result := r.db.Create(profile)
	if result.Error != nil {
		middleware.DebugLog("❌ DEBUG: Repository.Create error: %v\n", result.Error)
		return result.Error
	}

	middleware.DebugLog("✅ DEBUG: Repository.Create completed successfully - New ID: %s, Rows affected: %d\n", profile.ID, result.RowsAffected)
	return nil
}

func (r *studentProfileRepository) GetByID(ctx context.Context, id string, profile *StudentProfile) (*StudentProfile, error) {
	err := r.db.Preload("Certificates").First(profile, "id = ?", id).Error
	return profile, err
}

func (r *studentProfileRepository) Update(ctx context.Context, profile *StudentProfile) error {
	middleware.DebugLog("🔍 DEBUG: Repository.Update called - ID: %s, Name: %s, Skills: %v\n", profile.ID, profile.Name, profile.Skills)

	// Check if profile exists
	var existingProfile StudentProfile
	err := r.db.Where("id = ?", profile.ID).First(&existingProfile).Error
	if err != nil {
		middleware.DebugLog("❌ DEBUG: Repository.Update - Profile not found: %v\n", err)
		return err
	}

	middleware.DebugLog("✅ DEBUG: Repository.Update - Profile found, updating...\n")
	middleware.DebugLog("🔍 DEBUG: Repository.Update - ProfilePhotoKey: %s, ResumeKey: %s\n", profile.ProfilePhotoKey, profile.ResumeKey)

	// Use map to force GORM to update ALL fields including zero values
	// GORM's struct-based Updates() skips zero-value fields even with Select("*")
	updateMap := map[string]interface{}{
		"name":              profile.Name,
		"email":             profile.Email,
		"location":          profile.Location,
		"phone_number":      profile.PhoneNumber,
		"profile_photo_key": profile.ProfilePhotoKey,
		"resume_key":        profile.ResumeKey,
		"skills":            profile.Skills,
		"experience":        profile.Experience,
		"education":         profile.Education,
		"portfolio":         profile.Portfolio,
		"linkedin":          profile.Linkedin,
		"github":            profile.Github,
	}

	result := r.db.Model(&StudentProfile{}).Where("id = ?", profile.ID).Updates(updateMap)
	if result.Error != nil {
		middleware.DebugLog("❌ DEBUG: Repository.Update error: %v\n", result.Error)
		return result.Error
	}

	middleware.DebugLog("✅ DEBUG: Repository.Update completed successfully - Rows affected: %d\n", result.RowsAffected)
	return nil
}

func (r *studentProfileRepository) Delete(ctx context.Context, id string, profile *StudentProfile) error {
	return r.db.Delete(profile, "id = ?", id).Error
}

func (r *studentProfileRepository) SoftDelete(ctx context.Context, id string, deletedBy string) error {
	return r.db.Model(&StudentProfile{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("NOW()")).Error
}

func (r *studentProfileRepository) Restore(ctx context.Context, id string) error {
	return r.db.Model(&StudentProfile{}).Where("id = ?", id).Update("deleted_at", nil).Error
}

func (r *studentProfileRepository) List(ctx context.Context, limit, offset int) ([]*StudentProfile, error) {
	var profiles []*StudentProfile
	err := r.db.Preload("Certificates").Limit(limit).Offset(offset).Find(&profiles).Error
	return profiles, err
}

func (r *studentProfileRepository) ListWithDeleted(ctx context.Context, limit, offset int) ([]*StudentProfile, error) {
	var profiles []*StudentProfile
	err := r.db.Preload("Certificates").Unscoped().Limit(limit).Offset(offset).Find(&profiles).Error
	return profiles, err
}

func (r *studentProfileRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.Model(&StudentProfile{}).Count(&count).Error
	return count, err
}

func (r *studentProfileRepository) CountWithDeleted(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.Model(&StudentProfile{}).Unscoped().Count(&count).Error
	return count, err
}

func (r *studentProfileRepository) Exists(ctx context.Context, id string) (bool, error) {
	var count int64
	err := r.db.Model(&StudentProfile{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

func (r *studentProfileRepository) ExistsWithDeleted(ctx context.Context, id string) (bool, error) {
	var count int64
	err := r.db.Model(&StudentProfile{}).Unscoped().Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

func (r *studentProfileRepository) GetByCreatedBy(ctx context.Context, createdBy string, limit, offset int) ([]*StudentProfile, error) {
	var profiles []*StudentProfile
	err := r.db.Preload("Certificates").Where("created_by = ?", createdBy).Limit(limit).Offset(offset).Find(&profiles).Error
	return profiles, err
}

func (r *studentProfileRepository) GetByUpdatedBy(ctx context.Context, updatedBy string, limit, offset int) ([]*StudentProfile, error) {
	var profiles []*StudentProfile
	err := r.db.Preload("Certificates").Where("updated_by = ?", updatedBy).Limit(limit).Offset(offset).Find(&profiles).Error
	return profiles, err
}

func (r *studentProfileRepository) GetByDeletedBy(ctx context.Context, deletedBy string, limit, offset int) ([]*StudentProfile, error) {
	var profiles []*StudentProfile
	err := r.db.Preload("Certificates").Where("deleted_by = ?", deletedBy).Limit(limit).Offset(offset).Find(&profiles).Error
	return profiles, err
}

func (r *studentProfileRepository) CreateMany(ctx context.Context, profiles []*StudentProfile) error {
	return r.db.Create(profiles).Error
}

func (r *studentProfileRepository) UpdateMany(ctx context.Context, profiles []*StudentProfile) error {
	return r.db.Save(profiles).Error
}

func (r *studentProfileRepository) DeleteMany(ctx context.Context, ids []string) error {
	return r.db.Delete(&StudentProfile{}, ids).Error
}

func (r *studentProfileRepository) SoftDeleteMany(ctx context.Context, ids []string, deletedBy string) error {
	return r.db.Model(&StudentProfile{}).Where("id IN ?", ids).Update("deleted_at", gorm.Expr("NOW()")).Error
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

func (r *studentProfileRepository) AddCertificate(cert *Certificate) error {
	middleware.DebugLog("🔍 DEBUG: Repository.AddCertificate called - Name: %s, FileSize: %d bytes, FileType: %s\n", cert.Name, cert.FileSize, cert.FileType)

	// Ensure ID is empty so database generates proper ID
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
