// File: internal/application/repository.go

package application

import (
	"asa/internal/jobpost"
	"asa/internal/middleware"
	"context"

	"github.com/Kisanlink/kisanlink-db/pkg/base"
	"gorm.io/gorm"
)

type ApplicationRepository interface {
	base.Repository[*Application]
	GetByStudent(studentID string) ([]Application, error)
	GetByJob(jobID string) ([]Application, error)
	GetJobMetadata(jobID string) (*JobPostMetadata, error)
	UpdateStatus(appID, studentID, status string) error
	UpdateStatusByEmployer(appID, jobID, employerID, status string) error
	GetJobEmployerID(jobID string) (string, error)
	GetApplicationsCountByJob(jobID string) (int, error)
	GetCandidateName(applicationID string) (string, error)
}

type applicationRepository struct {
	*base.BaseRepository[*Application]
	db *gorm.DB
}

func NewApplicationRepository(db *gorm.DB) ApplicationRepository {
	return &applicationRepository{
		BaseRepository: base.NewBaseRepository[*Application](),
		db:             db,
	}
}

func (r *applicationRepository) Create(ctx context.Context, app *Application) error {
	return r.db.Create(app).Error
}

func (r *applicationRepository) GetByID(ctx context.Context, id string, app *Application) (*Application, error) {
	err := r.db.First(app, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return app, nil
}

func (r *applicationRepository) Update(ctx context.Context, app *Application) error {
	return r.db.Save(app).Error
}

func (r *applicationRepository) Delete(ctx context.Context, id string, app *Application) error {
	return r.db.Delete(app, "id = ?", id).Error
}

func (r *applicationRepository) SoftDelete(ctx context.Context, id string, deletedBy string) error {
	return r.db.Model(&Application{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("NOW()")).Error
}

func (r *applicationRepository) Restore(ctx context.Context, id string) error {
	return r.db.Model(&Application{}).Where("id = ?", id).Update("deleted_at", nil).Error
}

func (r *applicationRepository) List(ctx context.Context, limit, offset int) ([]*Application, error) {
	var apps []*Application
	err := r.db.Limit(limit).Offset(offset).Find(&apps).Error
	return apps, err
}

func (r *applicationRepository) ListWithDeleted(ctx context.Context, limit, offset int) ([]*Application, error) {
	var apps []*Application
	err := r.db.Unscoped().Limit(limit).Offset(offset).Find(&apps).Error
	return apps, err
}

func (r *applicationRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.Model(&Application{}).Count(&count).Error
	return count, err
}

func (r *applicationRepository) CountWithDeleted(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.Model(&Application{}).Unscoped().Count(&count).Error
	return count, err
}

func (r *applicationRepository) Exists(ctx context.Context, id string) (bool, error) {
	var count int64
	err := r.db.Model(&Application{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

func (r *applicationRepository) ExistsWithDeleted(ctx context.Context, id string) (bool, error) {
	var count int64
	err := r.db.Model(&Application{}).Unscoped().Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

func (r *applicationRepository) GetByCreatedBy(ctx context.Context, createdBy string, limit, offset int) ([]*Application, error) {
	var apps []*Application
	err := r.db.Where("created_by = ?", createdBy).Limit(limit).Offset(offset).Find(&apps).Error
	return apps, err
}

func (r *applicationRepository) GetByUpdatedBy(ctx context.Context, updatedBy string, limit, offset int) ([]*Application, error) {
	var apps []*Application
	err := r.db.Where("updated_by = ?", updatedBy).Limit(limit).Offset(offset).Find(&apps).Error
	return apps, err
}

func (r *applicationRepository) GetByDeletedBy(ctx context.Context, deletedBy string, limit, offset int) ([]*Application, error) {
	var apps []*Application
	err := r.db.Where("deleted_by = ?", deletedBy).Limit(limit).Offset(offset).Find(&apps).Error
	return apps, err
}

func (r *applicationRepository) CreateMany(ctx context.Context, apps []*Application) error {
	return r.db.Create(apps).Error
}

func (r *applicationRepository) UpdateMany(ctx context.Context, apps []*Application) error {
	return r.db.Save(apps).Error
}

func (r *applicationRepository) DeleteMany(ctx context.Context, ids []string) error {
	return r.db.Delete(&Application{}, ids).Error
}

func (r *applicationRepository) SoftDeleteMany(ctx context.Context, ids []string, deletedBy string) error {
	return r.db.Model(&Application{}).Where("id IN ?", ids).Update("deleted_at", gorm.Expr("NOW()")).Error
}

func (r *applicationRepository) GetByStudent(studentID string) ([]Application, error) {
	var apps []Application
	err := r.db.Where("student_id = ?", studentID).Order("applied_at DESC").Find(&apps).Error
	return apps, err
}

func (r *applicationRepository) GetByJob(jobID string) ([]Application, error) {
	middleware.DebugLog("DEBUG: Repository GetByJob - JobID: %s\n", jobID)

	var apps []Application
	err := r.db.Raw(`
		SELECT 
			a.id,
			a.job_id,
			a.student_id,
			a.applied_at,
			a.status,
			a.cover_letter,
			a.resume_key,
			a.resume_file_name,
			a.resume_file_type,
			a.resume_file_size,
			a.job_title,
			a.company,
			a.location,
			a.job_type,
			a.experience,
			a.updated_at,
			COALESCE(sp.phone_number, '') as student_phone_number
		FROM applications a
		LEFT JOIN student_profiles sp ON a.student_id = sp.user_id
		WHERE a.job_id = ?
		ORDER BY a.applied_at DESC
	`, jobID).Scan(&apps).Error

	middleware.DebugLog("DEBUG: Repository GetByJob result - Found %d applications, Error: %v\n", len(apps), err)

	// Debug log each application to see the ID
	for i, app := range apps {
		middleware.DebugLog("DEBUG: Application %d - ID: %s, JobID: %s, StudentID: %s\n", i, app.ID, app.JobID, app.StudentID)
	}
	return apps, err
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
	middleware.DebugLog("DEBUG: Repository GetJobEmployerID - JobID: %s\n", jobID)

	var employerID string
	err := r.db.Model(&jobpost.JobPost{}).
		Where("id = ?", jobID).
		Select("employer_id").
		Scan(&employerID).Error

	middleware.DebugLog("DEBUG: Repository GetJobEmployerID result - EmployerID: %s, Error: %v\n", employerID, err)
	return employerID, err
}

func (r *applicationRepository) GetApplicationsCountByJob(jobID string) (int, error) {
	var count int64
	err := r.db.Model(&Application{}).
		Where("job_id = ?", jobID).
		Count(&count).Error
	return int(count), err
}

func (r *applicationRepository) GetCandidateName(applicationID string) (string, error) {
	var candidateName string
	err := r.db.Raw(`
		SELECT u.name 
		FROM applications a
		JOIN users u ON u.id = a.student_id
		WHERE a.id = ?
	`, applicationID).Scan(&candidateName).Error
	return candidateName, err
}
