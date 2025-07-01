// File: internal/employerapplication/repository.go

package employerapplication

import (
	"gorm.io/gorm"
)

type EmployerApplicationRepository interface {
	GetApplicationsForJob(jobID string) ([]JobApplicationWithApplicant, error)
	UpdateStatus(applicationID, status string) error
	GetApplicantProfile(studentID string) (*ApplicantProfile, error)
	AddMessage(msg *Message) error
}

type employerApplicationRepository struct {
	db *gorm.DB
}

func NewEmployerApplicationRepository(db *gorm.DB) EmployerApplicationRepository {
	return &employerApplicationRepository{db}
}

func (r *employerApplicationRepository) GetApplicationsForJob(jobID string) ([]JobApplicationWithApplicant, error) {
	// For demo, you should join "applications" and "users" and build JobApplicationWithApplicant
	return []JobApplicationWithApplicant{}, nil
}

func (r *employerApplicationRepository) UpdateStatus(applicationID, status string) error {
	return r.db.Table("applications").Where("id = ?", applicationID).Update("status", status).Error
}

func (r *employerApplicationRepository) GetApplicantProfile(studentID string) (*ApplicantProfile, error) {
	// Join user_profiles, users, etc.
	return &ApplicantProfile{}, nil
}

func (r *employerApplicationRepository) AddMessage(msg *Message) error {
	return r.db.Create(msg).Error
}
