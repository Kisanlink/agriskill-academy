// File: internal/application/service.go

package application

import (
	"github.com/Kisanlink/agriskill-academy/internal/jobpost"
	"github.com/Kisanlink/agriskill-academy/internal/middleware"
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	db "github.com/Kisanlink/agriskill-academy/pkg/db"
)

type ApplicationService interface {
	Apply(app *Application) error
	GetMyApplications(studentID string) ([]Application, error)
	GetApplicationsByJob(jobID, employerID string) ([]Application, error)
	GetApplicationByID(appID string) (*Application, error)
	Remove(appID, studentID string) error
	UpdateStatus(appID, studentID, status string) error
	UpdateStatusByEmployer(appID, jobID, employerID, status string) error
	UploadResumeToS3(file *multipart.FileHeader, studentID string) (string, error)
	GetApplicationsCountByJob(jobID string) (int, error)
}

type applicationService struct {
	repo    ApplicationRepository
	jobRepo jobpost.JobPostRepository
	s3      *db.S3Manager
}

func NewApplicationService(repo ApplicationRepository, jobRepo jobpost.JobPostRepository, s3 *db.S3Manager) ApplicationService {
	return &applicationService{repo: repo, jobRepo: jobRepo, s3: s3}
}

func (s *applicationService) Apply(app *Application) error {
	middleware.DebugLog("DEBUG: Service Apply called for JobID: %s, StudentID: %s\n", app.JobID, app.StudentID)

	// Check if application already exists using a custom query
	var count int64
	err := s.repo.(*applicationRepository).db.Model(&Application{}).
		Where("job_id = ? AND student_id = ?", app.JobID, app.StudentID).
		Count(&count).Error
	if err != nil {
		middleware.DebugLog("DEBUG: Error checking if application exists: %v\n", err)
		return err
	}
	if count > 0 {
		middleware.DebugLog("DEBUG: Application already exists\n")
		return fmt.Errorf("application already exists")
	}

	middleware.DebugLog("DEBUG: No existing application found, proceeding...\n")

	// Populate job metadata
	app.AppliedAt = time.Now()
	app.Status = StatusApplied

	middleware.DebugLog("DEBUG: Fetching job metadata for JobID: %s\n", app.JobID)
	job, err := s.repo.GetJobMetadata(app.JobID)
	if err != nil {
		middleware.DebugLog("DEBUG: Error fetching job metadata: %v\n", err)
		return err
	}

	middleware.DebugLog("DEBUG: Job metadata fetched: %+v\n", job)

	app.JobTitle = job.Title
	app.Company = job.EmployerName
	app.Location = job.Location
	app.JobType = job.JobType
	app.Experience = job.Experience

	middleware.DebugLog("DEBUG: Application object before save: %+v\n", app)

	err = s.repo.Create(context.Background(), app)
	if err != nil {
		middleware.DebugLog("DEBUG: Error creating application in database: %v\n", err)
		return err
	}

	middleware.DebugLog("DEBUG: Application created successfully\n")
	return nil
}

func (s *applicationService) GetMyApplications(studentID string) ([]Application, error) {
	return s.repo.GetByStudent(studentID)
}

func (s *applicationService) GetApplicationsByJob(jobID, employerID string) ([]Application, error) {
	middleware.DebugLog("DEBUG: Service GetApplicationsByJob - JobID: %s, EmployerID: %s\n", jobID, employerID)

	// Verify that the job belongs to the employer
	jobEmployerID, err := s.repo.GetJobEmployerID(jobID)
	if err != nil {
		middleware.DebugLog("DEBUG: Error getting job employer ID: %v\n", err)
		return nil, err
	}

	middleware.DebugLog("DEBUG: Job employer ID: %s, Requesting employer ID: %s\n", jobEmployerID, employerID)

	if jobEmployerID != employerID {
		middleware.DebugLog("DEBUG: Authorization failed - job belongs to %s, requesting user is %s\n", jobEmployerID, employerID)
		return nil, errors.New("not authorized to view applications for this job")
	}

	apps, err := s.repo.GetByJob(jobID)
	if err != nil {
		middleware.DebugLog("DEBUG: Error getting applications by job: %v\n", err)
		return nil, err
	}

	middleware.DebugLog("DEBUG: Found %d applications for job %s\n", len(apps), jobID)
	return apps, err
}

func (s *applicationService) GetApplicationByID(appID string) (*Application, error) {
	return s.repo.GetByID(context.Background(), appID, &Application{})
}

func (s *applicationService) Remove(appID, studentID string) error {
	// First get the application to verify ownership
	app, err := s.repo.GetByID(context.Background(), appID, &Application{})
	if err != nil {
		return err
	}

	// Verify the application belongs to the student
	if app.StudentID != studentID {
		return errors.New("not authorized to delete this application")
	}

	return s.repo.Delete(context.Background(), appID, &Application{})
}

func (s *applicationService) UpdateStatus(appID, studentID, status string) error {
	// Validate status
	if !IsValidStatus(status) {
		return fmt.Errorf("invalid status: %s", status)
	}

	// Only allow students to withdraw their own applications
	if status != StatusWithdrawn {
		return errors.New("students can only withdraw their applications")
	}

	return s.repo.UpdateStatus(appID, studentID, status)
}

func (s *applicationService) UpdateStatusByEmployer(appID, jobID, employerID, status string) error {
	// Validate status
	if !IsValidStatus(status) {
		return fmt.Errorf("invalid status: %s", status)
	}

	// Employers cannot withdraw applications (only students can)
	if status == StatusWithdrawn {
		return errors.New("employers cannot withdraw applications")
	}

	// Update the application status
	err := s.repo.UpdateStatusByEmployer(appID, jobID, employerID, status)
	if err != nil {
		return err
	}

	// If the application is accepted, update the job post with the hired candidate
	if status == StatusAccepted {
		// Get the candidate name
		candidateName, err := s.repo.GetCandidateName(appID)
		if err != nil {
			middleware.DebugLog("DEBUG: Error getting candidate name: %v\n", err)
			return fmt.Errorf("failed to get candidate name: %w", err)
		}

		// Update the job post with the hired candidate
		err = s.jobRepo.UpdateHiredCandidate(jobID, candidateName)
		if err != nil {
			middleware.DebugLog("DEBUG: Error updating job with hired candidate: %v\n", err)
			return fmt.Errorf("failed to update job with hired candidate: %w", err)
		}

		middleware.DebugLog("DEBUG: Successfully updated job %s with hired candidate: %s\n", jobID, candidateName)
	}

	return nil
}

func (s *applicationService) UploadResumeToS3(file *multipart.FileHeader, studentID string) (string, error) {
	middleware.DebugLog("DEBUG: Uploading resume for student %s, file: %s", studentID, file.Filename)

	// Open the file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// Generate S3 key
	timestamp := time.Now().UnixNano()
	ext := filepath.Ext(file.Filename)
	baseName := strings.TrimSuffix(file.Filename, ext)
	safeBaseName := strings.ReplaceAll(baseName, " ", "_")
	filename := fmt.Sprintf("%d_%s%s", timestamp, safeBaseName, ext)
	s3Key := fmt.Sprintf("application_resumes/%s_%s", studentID, filename)

	// Get content type
	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Upload to S3
	ctx := context.Background()
	err = s.s3.UploadFile(ctx, s3Key, src, contentType, nil)
	if err != nil {
		middleware.DebugLog("DEBUG: Error uploading file to S3: %v\n", err)
		return "", err
	}

	middleware.DebugLog("DEBUG: File uploaded successfully to S3 with key: %s\n", s3Key)
	return s3Key, nil
}

func (s *applicationService) GetApplicationsCountByJob(jobID string) (int, error) {
	return s.repo.GetApplicationsCountByJob(jobID)
}
