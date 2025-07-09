// File: internal/studentprofile/service.go

package studentprofile

type StudentProfileService interface {
	GetProfile(userID string) (*StudentProfile, error)
	UpdateProfile(profile *StudentProfile) error
	CreateProfile(profile *StudentProfile) error
	AddCertificate(cert *Certificate) error
}

type studentProfileService struct {
	repo StudentProfileRepository
}

func NewStudentProfileService(repo StudentProfileRepository) StudentProfileService {
	return &studentProfileService{repo}
}

func (s *studentProfileService) GetProfile(userID string) (*StudentProfile, error) {
	return s.repo.GetByUserID(userID)
}

func (s *studentProfileService) UpdateProfile(profile *StudentProfile) error {
	return s.repo.Update(profile)
}

func (s *studentProfileService) CreateProfile(profile *StudentProfile) error {
	return s.repo.Create(profile)
}

func (s *studentProfileService) AddCertificate(cert *Certificate) error {
	return s.repo.AddCertificate(cert)
}
