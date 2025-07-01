// File: internal/employerprofile/service.go

package employerprofile

type EmployerProfileService interface {
	GetProfile(userID string) (*EmployerProfile, error)
	UpdateProfile(profile *EmployerProfile) error
	CreateProfile(profile *EmployerProfile) error
}

type employerProfileService struct {
	repo EmployerProfileRepository
}

func NewEmployerProfileService(repo EmployerProfileRepository) EmployerProfileService {
	return &employerProfileService{repo}
}

func (s *employerProfileService) GetProfile(userID string) (*EmployerProfile, error) {
	return s.repo.GetByUserID(userID)
}

func (s *employerProfileService) UpdateProfile(profile *EmployerProfile) error {
	return s.repo.Update(profile)
}

func (s *employerProfileService) CreateProfile(profile *EmployerProfile) error {
	return s.repo.Create(profile)
}
