package employerprofile

type EmployerProfileService interface {
	GetProfile(userID string) (*EmployerProfile, error)
	UpdateProfile(profile *EmployerProfile) error
	CreateProfile(profile *EmployerProfile) error
	DeleteProfile(userID string) error
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

func (s *employerProfileService) DeleteProfile(userID string) error {
	return s.repo.DeleteByUserID(userID)
}
