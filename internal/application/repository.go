// File: internal/application/repository.go

package application

import (
	"gorm.io/gorm"
)

type ApplicationRepository interface {
	Create(app *Application) error
	GetByStudent(studentID string) ([]Application, error)
	Delete(appID, studentID string) error
	Exists(jobID, studentID string) (bool, error)
	GetJobMetadata(jobID string) (*JobPostMetadata, error)
	UpdateStatus(appID, studentID, status string) error
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
	err := r.db.Where("student_id = ?", studentID).Find(&apps).Error
	return apps, err
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
