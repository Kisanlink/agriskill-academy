package auth

import (
	"time"
)

type User struct {
	ID         string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email" binding:"required,email"`
	Password   string    `json:"password" binding:"required"`
	Role       string    `json:"role" binding:"required"`
	Avatar     []byte    `json:"avatar" gorm:"type:bytea"`
	AvatarName string    `json:"avatar_name"`
	AvatarType string    `json:"avatar_type"`
	AvatarSize int64     `json:"avatar_size"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// TableName specifies the database table name for User
func (User) TableName() string {
	return "users"
}

type SignupRequest struct {
	Name            string `json:"name" binding:"required"`
	Username        string `json:"user_name" binding:"required"` // Username for AAA service
	Email           string `json:"email" binding:"required"`     // Email for our local DB
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
	Role            string `json:"role" binding:"required,oneof=student employer"` // "student" or "employer"
	PhoneNumber     string `json:"phone_number" binding:"required"`                // Changed to string for frontend compatibility
	CountryCode     string `json:"country_code,omitempty"`                         // Optional, defaults to "+91"
	AadhaarNumber   string `json:"aadhaar_number,omitempty"`                       // Optional

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
	Username string `json:"user_name" binding:"required"` // Username for AAA service
	Password string `json:"password" binding:"required"`
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
