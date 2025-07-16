package auth

import (
	"asa/internal/middleware"
	"errors"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindByEmail(email string) (*User, error)
	FindByID(id string) (*User, error)
	Create(user *User) error
	CreateWithID(user *User, id string) error
	Update(user *User) error
	ListAllUsers() ([]User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) FindByEmail(email string) (*User, error) {
	middleware.DebugLog("🔍 Repository.FindByEmail called with email: %s\n", email)
	var user User
	err := r.db.Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		middleware.DebugLog("❌ Repository.FindByEmail: User not found for email: %s\n", email)
		return nil, errors.New("user not found")
	}
	if err != nil {
		middleware.DebugLog("❌ Repository.FindByEmail error: %v\n", err)
		return nil, err
	}
	middleware.DebugLog("✅ Repository.FindByEmail: Found user - ID: %s, Name: %s, Email: %s\n", user.ID, user.Name, user.Email)
	return &user, err
}

func (r *userRepository) FindByID(id string) (*User, error) {
	middleware.DebugLog("🔍 Repository.FindByID called with ID: %s\n", id)
	var user User
	err := r.db.First(&user, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		middleware.DebugLog("❌ Repository.FindByID: User not found for ID: %s\n", id)
		return nil, errors.New("user not found")
	}
	if err != nil {
		middleware.DebugLog("❌ Repository.FindByID error: %v\n", err)
		return nil, err
	}
	middleware.DebugLog("✅ Repository.FindByID: Found user - ID: %s, Name: %s, Email: %s\n", user.ID, user.Name, user.Email)
	return &user, err
}

func (r *userRepository) Create(user *User) error {
	middleware.DebugLog("🔍 Repository.Create called with user: %+v\n", user)
	middleware.DebugLog("🔍 User ID before create: %s\n", user.ID)
	err := r.db.Create(user).Error
	if err != nil {
		middleware.DebugLog("❌ Repository.Create failed: %v\n", err)
	} else {
		middleware.DebugLog("✅ Repository.Create successful, user ID: %s\n", user.ID)
	}
	return err
}

func (r *userRepository) CreateWithID(user *User, id string) error {
	middleware.DebugLog("🔍 Repository.CreateWithID called with user: %+v, ID: %s\n", user, id)
	user.ID = id
	middleware.DebugLog("🔍 Setting user ID to: %s\n", user.ID)
	err := r.db.Create(user).Error
	if err != nil {
		middleware.DebugLog("❌ Repository.CreateWithID failed: %v\n", err)
	} else {
		middleware.DebugLog("✅ Repository.CreateWithID successful, user ID: %s\n", user.ID)
	}
	return err
}

func (r *userRepository) Update(user *User) error {
	return r.db.Save(user).Error
}

// Debug method to list all users (for debugging only)
func (r *userRepository) ListAllUsers() ([]User, error) {
	middleware.DebugLog("🔍 Repository.ListAllUsers called\n")
	var users []User
	err := r.db.Find(&users).Error
	if err != nil {
		middleware.DebugLog("❌ Repository.ListAllUsers error: %v\n", err)
		return nil, err
	}
	middleware.DebugLog("✅ Repository.ListAllUsers: Found %d users\n", len(users))
	for i, user := range users {
		middleware.DebugLog("   User %d: ID=%s, Name=%s, Email=%s, Role=%s\n", i+1, user.ID, user.Name, user.Email, user.Role)
	}
	return users, nil
}
