// File: internal/application/repository.go

package application

import (
	"fmt"

	"asa/internal/jobpost"

	"gorm.io/gorm"
)

type ApplicationRepository interface {
	Create(app *Application) error
	GetByStudent(studentID string) ([]Application, error)
	GetByJob(jobID string) ([]Application, error)
	GetByID(appID string) (*Application, error)
	Delete(appID, studentID string) error
	Exists(jobID, studentID string) (bool, error)
	GetJobMetadata(jobID string) (*JobPostMetadata, error)
	UpdateStatus(appID, studentID, status string) error
	UpdateStatusByEmployer(appID, jobID, employerID, status string) error
	GetJobEmployerID(jobID string) (string, error)
}

type applicationRepository struct {
	db *gorm.DB
}

func NewApplicationRepository(db *gorm.DB) ApplicationRepository {
	return &applicationRepository{db}
}

func (r *applicationRepository) Create(app *Application) error {
	return r.db.Create(app).Error
}

func (r *applicationRepository) GetByStudent(studentID string) ([]Application, error) {
	var apps []Application
	err := r.db.Where("student_id = ?", studentID).Order("applied_at DESC").Find(&apps).Error
	return apps, err
}

func (r *applicationRepository) GetByJob(jobID string) ([]Application, error) {
	fmt.Printf("DEBUG: Repository GetByJob - JobID: %s\n", jobID)

	var apps []Application
	err := r.db.Where("job_id = ?", jobID).Order("applied_at DESC").Find(&apps).Error

	fmt.Printf("DEBUG: Repository GetByJob result - Found %d applications, Error: %v\n", len(apps), err)
	return apps, err
}

func (r *applicationRepository) GetByID(appID string) (*Application, error) {
	var app Application
	err := r.db.Where("id = ?", appID).First(&app).Error
	return &app, err
}

func (r *applicationRepository) Delete(appID, studentID string) error {
	return r.db.Where("id = ? AND student_id = ?", appID, studentID).Delete(&Application{}).Error
}

func (r *applicationRepository) Exists(jobID, studentID string) (bool, error) {
	var count int64
	err := r.db.Model(&Application{}).
		Where("job_id = ? AND student_id = ?", jobID, studentID).
		Count(&count).Error
	return count > 0, err
}

type JobPostMetadata struct {
	Title        string
	EmployerName string
	Location     string
	JobType      string
	Experience   string
}

func (r *applicationRepository) GetJobMetadata(jobID string) (*JobPostMetadata, error) {
	var meta JobPostMetadata
	err := r.db.Raw(`
		SELECT title, employer_name, location, job_type, experience
		FROM job_posts
		WHERE id = ?
	`, jobID).Scan(&meta).Error
	return &meta, err
}

func (r *applicationRepository) UpdateStatus(appID, studentID, status string) error {
	return r.db.Model(&Application{}).
		Where("id = ? AND student_id = ?", appID, studentID).
		Update("status", status).Error
}

func (r *applicationRepository) UpdateStatusByEmployer(appID, jobID, employerID, status string) error {
	// Verify that the job belongs to the employer
	var count int64
	err := r.db.Model(&jobpost.JobPost{}).
		Where("id = ? AND employer_id = ?", jobID, employerID).
		Count(&count).Error

	if err != nil {
		return err
	}

	if count == 0 {
		return gorm.ErrRecordNotFound
	}

	// Update the application status
	return r.db.Model(&Application{}).
		Where("id = ? AND job_id = ?", appID, jobID).
		Update("status", status).Error
}

func (r *applicationRepository) GetJobEmployerID(jobID string) (string, error) {
	fmt.Printf("DEBUG: Repository GetJobEmployerID - JobID: %s\n", jobID)

	var employerID string
	err := r.db.Model(&jobpost.JobPost{}).
		Where("id = ?", jobID).
		Select("employer_id").
		Scan(&employerID).Error

	fmt.Printf("DEBUG: Repository GetJobEmployerID result - EmployerID: %s, Error: %v\n", employerID, err)
	return employerID, err
}
