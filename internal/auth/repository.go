package auth

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindByEmail(email string) (*User, error)
	FindByID(id string) (*User, error)
	Create(user *User) error
	Update(user *User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) FindByEmail(email string) (*User, error) {
	var user User
	err := r.db.Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}
	return &user, err
}

func (r *userRepository) FindByID(id string) (*User, error) {
	var user User
	err := r.db.First(&user, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}
	return &user, err
}

func (r *userRepository) Create(user *User) error {
	fmt.Printf("🔍 Repository.Create called with user: %+v\n", user)
	err := r.db.Create(user).Error
	if err != nil {
		fmt.Printf("❌ Repository.Create failed: %v\n", err)
	} else {
		fmt.Printf("✅ Repository.Create successful, user ID: %s\n", user.ID)
	}
	return err
}

func (r *userRepository) Update(user *User) error {
	return r.db.Save(user).Error
}
