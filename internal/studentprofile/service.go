// File: internal/studentprofile/service.go

package studentprofile

import "asa/internal/middleware"

type StudentProfileService interface {
	GetProfile(userID string) (*StudentProfile, error)
	UpdateProfile(profile *StudentProfile) error
	CreateProfile(profile *StudentProfile) error
	AddCertificate(cert *Certificate) error
	DeleteCertificate(certID string, userID string) error
}

type studentProfileService struct {
	repo StudentProfileRepository
}

func NewStudentProfileService(repo StudentProfileRepository) StudentProfileService {
	return &studentProfileService{repo}
}

func (s *studentProfileService) GetProfile(userID string) (*StudentProfile, error) {
	middleware.DebugLog("🔍 DEBUG: Service.GetProfile called for userID: %s\n", userID)
	profile, err := s.repo.GetByUserID(userID)
	if err != nil {
		middleware.DebugLog("❌ DEBUG: Service.GetProfile error: %v\n", err)
		return nil, err
	}
	middleware.DebugLog("✅ DEBUG: Service.GetProfile success: %+v\n", profile)
	return profile, err
}

func (s *studentProfileService) UpdateProfile(profile *StudentProfile) error {
	middleware.DebugLog("🔍 DEBUG: Service.UpdateProfile called - ID: %s, Name: %s, ProfilePhotoSize: %d, ResumeSize: %d\n", profile.ID, profile.Name, profile.ProfilePhotoSize, profile.ResumeSize)
	err := s.repo.Update(profile)
	if err != nil {
		middleware.DebugLog("❌ DEBUG: Service.UpdateProfile error: %v\n", err)
		return err
	}
	middleware.DebugLog("✅ DEBUG: Service.UpdateProfile completed successfully\n")
	return nil
}

func (s *studentProfileService) CreateProfile(profile *StudentProfile) error {
	middleware.DebugLog("🔍 DEBUG: Service.CreateProfile called - UserID: %s, Name: %s, ProfilePhotoSize: %d, ResumeSize: %d\n", profile.UserID, profile.Name, profile.ProfilePhotoSize, profile.ResumeSize)
	err := s.repo.Create(profile)
	if err != nil {
		middleware.DebugLog("❌ DEBUG: Service.CreateProfile error: %v\n", err)
		return err
	}
	middleware.DebugLog("✅ DEBUG: Service.CreateProfile completed successfully\n")
	return nil
}

func (s *studentProfileService) AddCertificate(cert *Certificate) error {
	middleware.DebugLog("🔍 DEBUG: Service.AddCertificate called - StudentProfileID: %s, Name: %s, FileSize: %d bytes, FileType: %s\n", cert.StudentProfileID, cert.Name, cert.FileSize, cert.FileType)
	err := s.repo.AddCertificate(cert)
	if err != nil {
		middleware.DebugLog("❌ DEBUG: Service.AddCertificate error: %v\n", err)
		return err
	}
	middleware.DebugLog("✅ DEBUG: Service.AddCertificate completed successfully\n")
	return nil
}

func (s *studentProfileService) DeleteCertificate(certID string, userID string) error {
	middleware.DebugLog("🔍 DEBUG: Service.DeleteCertificate called - CertID: %s, UserID: %s\n", certID, userID)
	err := s.repo.DeleteCertificate(certID, userID)
	if err != nil {
		middleware.DebugLog("❌ DEBUG: Service.DeleteCertificate error: %v\n", err)
		return err
	}
	middleware.DebugLog("✅ DEBUG: Service.DeleteCertificate completed successfully\n")
	return nil
}
