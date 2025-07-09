package auth

import "time"

type User struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name      string    `json:"name"`
	Email     string    `gorm:"uniqueIndex" json:"email"`
	Password  string    `json:"-"` // Stored locally but not exposed in JSON
	Role      string    `json:"-"` // Stored locally but not exposed in JSON
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type SignupRequest struct {
	Name            string `json:"name" binding:"required"`
	Username        string `json:"username" binding:"required"` // Username for AAA service
	Email           string `json:"email" binding:"required"`    // Email for our local DB
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirmPassword" binding:"required"`
	Role            string `json:"role" binding:"required,oneof=student employer"` // "student" or "employer"
	PhoneNumber     int64  `json:"phoneNumber" binding:"required"`                 // Changed to int64 for AAA service
	CountryCode     string `json:"countryCode,omitempty"`                          // Optional, defaults to "+91"
	AadhaarNumber   string `json:"aadhaarNumber,omitempty"`                        // Optional

	// Employer-only fields (optional)
	CompanyName    string `json:"companyName,omitempty"`
	GstinNumber    string `json:"gstinNumber,omitempty"`
	CompanyAddress string `json:"companyAddress,omitempty"`
	City           string `json:"city,omitempty"`
	State          string `json:"state,omitempty"`
	Pincode        string `json:"pincode,omitempty"`
	IndustryType   string `json:"industryType,omitempty"`
	CompanySize    string `json:"companySize,omitempty"`
	Website        string `json:"website,omitempty"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"` // Username for AAA service
	Password string `json:"password" binding:"required"`
}

type UpdateProfileRequest struct {
	// Basic user information
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`

	// Contact information
	PhoneNumber string `json:"phoneNumber,omitempty"`
	Location    string `json:"location,omitempty"`

	// Profile information
	ProfilePhoto string `json:"profilePhoto,omitempty"`
	Bio          string `json:"bio,omitempty"`

	// Social links
	LinkedinProfile string `json:"linkedinProfile,omitempty"`
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
	CompanyName        string   `json:"companyName,omitempty"`
	CompanyDescription string   `json:"companyDescription,omitempty"`
	Industry           string   `json:"industry,omitempty"`
	CompanySize        string   `json:"companySize,omitempty"`
	RecruiterName      string   `json:"recruiterName,omitempty"`
	Designation        string   `json:"designation,omitempty"`
	OfficialEmail      string   `json:"officialEmail,omitempty"`
	GstinNumber        string   `json:"gstinNumber,omitempty"`
	CompanyAddress     string   `json:"companyAddress,omitempty"`
	City               string   `json:"city,omitempty"`
	State              string   `json:"state,omitempty"`
	Pincode            string   `json:"pincode,omitempty"`
	JobCategories      []string `json:"jobCategories,omitempty"`
	HiringLocations    []string `json:"hiringLocations,omitempty"`
	HiringTypes        []string `json:"hiringTypes,omitempty"`
}

type ProfileResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	User    *User       `json:"user,omitempty"`
	Profile interface{} `json:"profile,omitempty"`
}
