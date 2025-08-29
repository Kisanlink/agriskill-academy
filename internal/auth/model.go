package auth

import (
	"github.com/Kisanlink/kisanlink-db/pkg/base"
	"github.com/Kisanlink/kisanlink-db/pkg/core/hash"
	"gorm.io/gorm"
)

type User struct {
	base.BaseModel
	Name        string `json:"name"`
	Username    string `json:"username" binding:"required" gorm:"uniqueIndex"` // Separate username field with unique constraint
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required"`
	Role        string `json:"role" binding:"required"`
	PhoneNumber string `json:"phone_number,omitempty"`
	AvatarKey   string `json:"avatar_key,omitempty"`
	AvatarName  string `json:"avatar_name,omitempty"`
	AvatarType  string `json:"avatar_type,omitempty"`
	AvatarSize  int64  `json:"avatar_size,omitempty"`
}

// TableName specifies the database table name for User
func (User) TableName() string {
	return "users"
}

// NewUser creates a new User with proper initialization
func NewUser() *User {
	return &User{
		BaseModel: *base.NewBaseModel("USER", hash.Medium),
	}
}

// BeforeCreateGORM is called by GORM before creating a new record
func (u *User) BeforeCreateGORM(tx *gorm.DB) error {
	return u.BeforeCreate()
}

// BeforeUpdateGORM is called by GORM before updating an existing record
func (u *User) BeforeUpdateGORM(tx *gorm.DB) error {
	return u.BeforeUpdate()
}

// BeforeDeleteGORM is called by GORM before hard deleting a record
func (u *User) BeforeDeleteGORM(tx *gorm.DB) error {
	return u.BeforeDelete()
}

type SignupRequest struct {
	Name            string `json:"name" binding:"required"`
	Username        string `json:"user_name" binding:"required"` // Username for local authentication
	Email           string `json:"email" binding:"required"`     // Email for our local DB
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
	Role            string `json:"role" binding:"required,oneof=student employer asa_admin"` // "student" or "employer"
	PhoneNumber     string `json:"phone_number" binding:"required"`                          // Changed to string for frontend compatibility
	CountryCode     string `json:"country_code,omitempty"`                                   // Optional, defaults to "+91"
	AadhaarNumber   string `json:"aadhaar_number,omitempty"`                                 // Optional

	// Employer-only fields (optional)
	CompanyName    string `json:"company_name,omitempty"`
	GstinNumber    string `json:"gstin_number,omitempty"`
	CompanyAddress string `json:"company_address,omitempty"`
	City           string `json:"city,omitempty"`
	State          string `json:"state,omitempty"`
	Pincode        string `json:"pincode,omitempty"`
	IndustryType   string `json:"industry_type,omitempty"`
	CompanySize    string `json:"company_size,omitempty"`
	Website        string `json:"website,omitempty"`
}

type LoginRequest struct {
	Username string `json:"user_name" binding:"required"` // Username for local authentication
	Password string `json:"password" binding:"required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

type UpdateProfileRequest struct {
	// Basic user information
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`

	// Contact information
	PhoneNumber string `json:"phone_number,omitempty"`
	Location    string `json:"location,omitempty"`

	// Profile information
	ProfilePhoto string `json:"profile_photo,omitempty"`
	Bio          string `json:"bio,omitempty"`

	// Social links
	LinkedinProfile string `json:"linkedin_profile,omitempty"`
	Website         string `json:"website,omitempty"`

	// Role-specific fields (will be validated based on user role)
	// For students
	Skills     []string `json:"skills,omitempty"`
	Resume     string   `json:"resume,omitempty"`
	Experience float64  `json:"experience,omitempty"`
	Education  string   `json:"education,omitempty"`
	Portfolio  string   `json:"portfolio,omitempty"`
	Github     string   `json:"github,omitempty"`

	// For employers
	CompanyName        string   `json:"company_name,omitempty"`
	CompanyDescription string   `json:"company_description,omitempty"`
	Industry           string   `json:"industry,omitempty"`
	CompanySize        string   `json:"company_size,omitempty"`
	RecruiterName      string   `json:"recruiter_name,omitempty"`
	Designation        string   `json:"designation,omitempty"`
	OfficialEmail      string   `json:"official_email,omitempty"`
	GstinNumber        string   `json:"gstin_number,omitempty"`
	CompanyAddress     string   `json:"company_address,omitempty"`
	City               string   `json:"city,omitempty"`
	State              string   `json:"state,omitempty"`
	Pincode            string   `json:"pincode,omitempty"`
	JobCategories      []string `json:"job_categories,omitempty"`
	HiringLocations    []string `json:"hiring_locations,omitempty"`
	HiringTypes        []string `json:"hiring_types,omitempty"`
}

type ProfileResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	User    *User       `json:"user,omitempty"`
	Profile interface{} `json:"profile,omitempty"`
}
