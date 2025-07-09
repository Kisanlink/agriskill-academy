package auth

import "time"

type User struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name      string    `json:"name"`
	Email     string    `gorm:"uniqueIndex" json:"email"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type SignupRequest struct {
	Name            string `json:"name" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirmPassword" binding:"required"`
	Role            string `json:"role" binding:"required"`        // "student" or "employer"
	PhoneNumber     string `json:"phoneNumber" binding:"required"` // Added phone number for signup

	// Employer-only fields
	CompanyName    string `json:"companyName"`
	GstinNumber    string `json:"gstinNumber"`
	CompanyAddress string `json:"companyAddress"`
	City           string `json:"city"`
	State          string `json:"state"`
	Pincode        string `json:"pincode"`
	IndustryType   string `json:"industryType"`
	CompanySize    string `json:"companySize"`
	Website        string `json:"website"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"` //(required,email)
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
