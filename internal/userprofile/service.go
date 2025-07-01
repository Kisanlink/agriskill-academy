// File: internal/userprofile/service.go

package userprofile

type UserProfileService interface {
	GetProfile(userID string) (*UserProfile, error)
	UpdateProfile(profile *UserProfile) error
	CreateProfile(profile *UserProfile) error
}

type userProfileService struct {
	repo UserProfileRepository
}

func NewUserProfileService(repo UserProfileRepository) UserProfileService {
	return &userProfileService{repo}
}

func (s *userProfileService) GetProfile(userID string) (*UserProfile, error) {
	return s.repo.GetByUserID(userID)
}

func (s *userProfileService) UpdateProfile(profile *UserProfile) error {
	return s.repo.Update(profile)
}

func (s *userProfileService) CreateProfile(profile *UserProfile) error {
	return s.repo.Create(profile)
}
