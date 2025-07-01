// File: internal/application/service.go

package application

import (
	"time"
)

type ApplicationService interface {
	Apply(app *Application) error
	GetMyApplications(studentID string) ([]Application, error)
	Remove(appID, studentID string) error
}

type applicationService struct {
	repo ApplicationRepository
}

func NewApplicationService(repo ApplicationRepository) ApplicationService {
	return &applicationService{repo}
}

func (s *applicationService) Apply(app *Application) error {
	app.AppliedAt = time.Now()
	app.Status = "applied"
	return s.repo.Create(app)
}

func (s *applicationService) GetMyApplications(studentID string) ([]Application, error) {
	return s.repo.GetByStudent(studentID)
}

func (s *applicationService) Remove(appID, studentID string) error {
	return s.repo.Delete(appID, studentID)
}
