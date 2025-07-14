// File: internal/application/service.go

package application

import (
	"errors"
	"fmt"
	"time"
)

type ApplicationService interface {
	Apply(app *Application) error
	GetMyApplications(studentID string) ([]Application, error)
	GetApplicationsByJob(jobID, employerID string) ([]Application, error)
	GetApplicationByID(appID string) (*Application, error)
	Remove(appID, studentID string) error
	UpdateStatus(appID, studentID, status string) error
	UpdateStatusByEmployer(appID, jobID, employerID, status string) error
}

type applicationService struct {
	repo ApplicationRepository
}

func NewApplicationService(repo ApplicationRepository) ApplicationService {
	return &applicationService{repo}
}

func (s *applicationService) Apply(app *Application) error {
	fmt.Printf("DEBUG: Service Apply called for JobID: %s, StudentID: %s\n", app.JobID, app.StudentID)

	// (UUID generation removed; handled by BeforeCreate hook)

	exists, err := s.repo.Exists(app.JobID, app.StudentID)
	if err != nil {
		fmt.Printf("DEBUG: Error checking if application exists: %v\n", err)
		return err
	}
	if exists {
		fmt.Printf("DEBUG: Application already exists\n")
		return fmt.Errorf("application already exists")
	}

	fmt.Printf("DEBUG: No existing application found, proceeding...\n")

	// Populate job metadata
	app.AppliedAt = time.Now()
	app.Status = StatusApplied

	fmt.Printf("DEBUG: Fetching job metadata for JobID: %s\n", app.JobID)
	job, err := s.repo.GetJobMetadata(app.JobID)
	if err != nil {
		fmt.Printf("DEBUG: Error fetching job metadata: %v\n", err)
		return err
	}

	fmt.Printf("DEBUG: Job metadata fetched: %+v\n", job)

	app.JobTitle = job.Title
	app.Company = job.EmployerName
	app.Location = job.Location
	app.JobType = job.JobType
	app.Experience = job.Experience

	fmt.Printf("DEBUG: Application object before save: %+v\n", app)

	err = s.repo.Create(app)
	if err != nil {
		fmt.Printf("DEBUG: Error creating application in database: %v\n", err)
		return err
	}

	fmt.Printf("DEBUG: Application created successfully\n")
	return nil
}

func (s *applicationService) GetMyApplications(studentID string) ([]Application, error) {
	return s.repo.GetByStudent(studentID)
}

func (s *applicationService) GetApplicationsByJob(jobID, employerID string) ([]Application, error) {
	fmt.Printf("DEBUG: Service GetApplicationsByJob - JobID: %s, EmployerID: %s\n", jobID, employerID)

	// Verify that the job belongs to the employer
	jobEmployerID, err := s.repo.GetJobEmployerID(jobID)
	if err != nil {
		fmt.Printf("DEBUG: Error getting job employer ID: %v\n", err)
		return nil, err
	}

	fmt.Printf("DEBUG: Job employer ID: %s, Requesting employer ID: %s\n", jobEmployerID, employerID)

	if jobEmployerID != employerID {
		fmt.Printf("DEBUG: Authorization failed - job belongs to %s, requesting user is %s\n", jobEmployerID, employerID)
		return nil, errors.New("not authorized to view applications for this job")
	}

	apps, err := s.repo.GetByJob(jobID)
	if err != nil {
		fmt.Printf("DEBUG: Error getting applications by job: %v\n", err)
		return nil, err
	}

	fmt.Printf("DEBUG: Found %d applications for job %s\n", len(apps), jobID)
	return apps, err
}

func (s *applicationService) GetApplicationByID(appID string) (*Application, error) {
	return s.repo.GetByID(appID)
}

func (s *applicationService) Remove(appID, studentID string) error {
	return s.repo.Delete(appID, studentID)
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

	return s.repo.UpdateStatusByEmployer(appID, jobID, employerID, status)
}
