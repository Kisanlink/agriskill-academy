// File: internal/application/service.go

package application

import (
	"fmt"
	"time"
)

type ApplicationService interface {
	Apply(app *Application) error
	GetMyApplications(studentID string) ([]Application, error)
	Remove(appID, studentID string) error
	UpdateStatus(appID, studentID, status string) error
}

type applicationService struct {
	repo ApplicationRepository
}

func NewApplicationService(repo ApplicationRepository) ApplicationService {
	return &applicationService{repo}
}

func (s *applicationService) Apply(app *Application) error {
	exists, _ := s.repo.Exists(app.JobID, app.StudentID)
	if exists {
		return fmt.Errorf("application already exists")
	}
	// Populate job metadata (below)
	app.AppliedAt = time.Now()
	app.Status = "applied"
	job, err := s.repo.GetJobMetadata(app.JobID)
	if err != nil {
		return err
	}
	app.JobTitle = job.Title
	app.Company = job.EmployerName
	app.Location = job.Location
	app.JobType = job.JobType
	app.Experience = job.Experience
	return s.repo.Create(app)
}

func (s *applicationService) GetMyApplications(studentID string) ([]Application, error) {
	return s.repo.GetByStudent(studentID)
}

func (s *applicationService) Remove(appID, studentID string) error {
	return s.repo.Delete(appID, studentID)
}

func (s *applicationService) UpdateStatus(appID, studentID, status string) error {
	return s.repo.UpdateStatus(appID, studentID, status)
}
