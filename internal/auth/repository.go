package auth

import (
	"github.com/Kisanlink/agriskill-academy/internal/middleware"
	"context"
	"errors"

	"github.com/Kisanlink/kisanlink-db/pkg/base"
	"gorm.io/gorm"
)

type UserRepository interface {
	base.Repository[*User]
	FindByEmail(email string) (*User, error)
	FindByUsername(username string) (*User, error)
	CreateWithID(user *User, id string) error
	ListAllUsers() ([]User, error)
}

type userRepository struct {
	*base.BaseRepository[*User]
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		BaseRepository: base.NewBaseRepository[*User](),
		db:             db,
	}
}

func (r *userRepository) Create(ctx context.Context, user *User) error {
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

func (r *userRepository) GetByID(ctx context.Context, id string, user *User) (*User, error) {
	middleware.DebugLog("🔍 Repository.GetByID called with ID: %s\n", id)
	err := r.db.First(user, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		middleware.DebugLog("❌ Repository.GetByID: User not found for ID: %s\n", id)
		return nil, errors.New("user not found")
	}
	if err != nil {
		middleware.DebugLog("❌ Repository.GetByID error: %v\n", err)
		return nil, err
	}
	middleware.DebugLog("✅ Repository.GetByID: Found user - ID: %s, Name: %s, Email: %s\n", user.ID, user.Name, user.Email)
	return user, err
}

func (r *userRepository) Update(ctx context.Context, user *User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(ctx context.Context, id string, user *User) error {
	return r.db.Delete(user, "id = ?", id).Error
}

func (r *userRepository) SoftDelete(ctx context.Context, id string, deletedBy string) error {
	return r.db.Model(&User{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("NOW()")).Error
}

func (r *userRepository) Restore(ctx context.Context, id string) error {
	return r.db.Model(&User{}).Where("id = ?", id).Update("deleted_at", nil).Error
}

func (r *userRepository) List(ctx context.Context, limit, offset int) ([]*User, error) {
	var users []*User
	err := r.db.Limit(limit).Offset(offset).Find(&users).Error
	return users, err
}

func (r *userRepository) ListWithDeleted(ctx context.Context, limit, offset int) ([]*User, error) {
	var users []*User
	err := r.db.Unscoped().Limit(limit).Offset(offset).Find(&users).Error
	return users, err
}

func (r *userRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.Model(&User{}).Count(&count).Error
	return count, err
}

func (r *userRepository) CountWithDeleted(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.Model(&User{}).Unscoped().Count(&count).Error
	return count, err
}

func (r *userRepository) Exists(ctx context.Context, id string) (bool, error) {
	var count int64
	err := r.db.Model(&User{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

func (r *userRepository) ExistsWithDeleted(ctx context.Context, id string) (bool, error) {
	var count int64
	err := r.db.Model(&User{}).Unscoped().Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

func (r *userRepository) GetByCreatedBy(ctx context.Context, createdBy string, limit, offset int) ([]*User, error) {
	var users []*User
	err := r.db.Where("created_by = ?", createdBy).Limit(limit).Offset(offset).Find(&users).Error
	return users, err
}

func (r *userRepository) GetByUpdatedBy(ctx context.Context, updatedBy string, limit, offset int) ([]*User, error) {
	var users []*User
	err := r.db.Where("updated_by = ?", updatedBy).Limit(limit).Offset(offset).Find(&users).Error
	return users, err
}

func (r *userRepository) GetByDeletedBy(ctx context.Context, deletedBy string, limit, offset int) ([]*User, error) {
	var users []*User
	err := r.db.Where("deleted_by = ?", deletedBy).Limit(limit).Offset(offset).Find(&users).Error
	return users, err
}

func (r *userRepository) CreateMany(ctx context.Context, users []*User) error {
	return r.db.Create(users).Error
}

func (r *userRepository) UpdateMany(ctx context.Context, users []*User) error {
	return r.db.Save(users).Error
}

func (r *userRepository) DeleteMany(ctx context.Context, ids []string) error {
	return r.db.Delete(&User{}, ids).Error
}

func (r *userRepository) SoftDeleteMany(ctx context.Context, ids []string, deletedBy string) error {
	return r.db.Model(&User{}).Where("id IN ?", ids).Update("deleted_at", gorm.Expr("NOW()")).Error
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

func (r *userRepository) FindByUsername(username string) (*User, error) {
	middleware.DebugLog("🔍 Repository.FindByUsername called with username: %s\n", username)
	var user User
	err := r.db.Where("username = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		middleware.DebugLog("❌ Repository.FindByUsername: User not found for username: %s\n", username)
		return nil, errors.New("user not found")
	}
	if err != nil {
		middleware.DebugLog("❌ Repository.FindByUsername error: %v\n", err)
		return nil, err
	}
	middleware.DebugLog("✅ Repository.FindByUsername: Found user - ID: %s, Name: %s, Email: %s\n", user.ID, user.Name, user.Email)
	return &user, err
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
