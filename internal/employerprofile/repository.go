package employerprofile

import "gorm.io/gorm"

type EmployerProfileRepository interface {
	GetByUserID(userID string) (*EmployerProfile, error)
	Update(profile *EmployerProfile) error
	Create(profile *EmployerProfile) error
	DeleteByUserID(userID string) error
}

type employerProfileRepository struct {
	db *gorm.DB
}

func NewEmployerProfileRepository(db *gorm.DB) EmployerProfileRepository {
	return &employerProfileRepository{db}
}

func (r *employerProfileRepository) GetByUserID(userID string) (*EmployerProfile, error) {
	var profile EmployerProfile
	err := r.db.Where("user_id = ?", userID).First(&profile).Error
	return &profile, err
}

func (r *employerProfileRepository) Update(profile *EmployerProfile) error {
	return r.db.Save(profile).Error
}

func (r *employerProfileRepository) Create(profile *EmployerProfile) error {
	return r.db.Create(profile).Error
}

func (r *employerProfileRepository) DeleteByUserID(userID string) error {
	return r.db.Where("user_id = ?", userID).Delete(&EmployerProfile{}).Error
}
