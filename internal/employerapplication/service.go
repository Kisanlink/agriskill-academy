// File: internal/employerapplication/service.go

package employerapplication

type EmployerApplicationService interface {
	GetApplicationsForJob(jobID string) ([]JobApplicationWithApplicant, error)
	UpdateStatus(applicationID, status string) error
	GetApplicantProfile(studentID string) (*ApplicantProfile, error)
	SendMessage(msg *Message) error
}

type employerApplicationService struct {
	repo EmployerApplicationRepository
}

func NewEmployerApplicationService(repo EmployerApplicationRepository) EmployerApplicationService {
	return &employerApplicationService{repo}
}

func (s *employerApplicationService) GetApplicationsForJob(jobID string) ([]JobApplicationWithApplicant, error) {
	return s.repo.GetApplicationsForJob(jobID)
}

func (s *employerApplicationService) UpdateStatus(applicationID, status string) error {
	return s.repo.UpdateStatus(applicationID, status)
}

func (s *employerApplicationService) GetApplicantProfile(studentID string) (*ApplicantProfile, error) {
	return s.repo.GetApplicantProfile(studentID)
}

func (s *employerApplicationService) SendMessage(msg *Message) error {
	return s.repo.AddMessage(msg)
}
