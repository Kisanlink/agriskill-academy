package employerprofile

import (
	"time"

	"github.com/lib/pq"
)

type EmployerProfile struct {
	ID     string `gorm:"primaryKey;type:varchar(255)" json:"id"`
	UserID string `gorm:"type:varchar(255);not null" json:"user_id" binding:"required"`

	// Required company information
	CompanyName string `gorm:"not null" json:"company_name" binding:"required"`
	Industry    string `gorm:"not null" json:"industry" binding:"required"`
	CompanySize string `gorm:"not null" json:"company_size" binding:"required"`

	// Optional company branding and details
	Logo               []byte `gorm:"type:bytea" json:"logo,omitempty"`
	LogoName           string `json:"logo_name,omitempty"`
	LogoType           string `json:"logo_type,omitempty"`
	LogoSize           int64  `json:"logo_size,omitempty"`
	WebsiteURL         string `json:"website_url,omitempty"`
	CompanyDescription string `json:"company_description,omitempty"`

	// Optional recruiter information
	RecruiterName   string `json:"recruiter_name,omitempty"`
	Designation     string `json:"designation,omitempty"`
	OfficialEmail   string `json:"official_email,omitempty"`
	PhoneNumber     string `json:"phone_number,omitempty"`
	LinkedinProfile string `json:"linkedin_profile,omitempty"`

	// Optional hiring preferences (can be set later)
	JobCategories   pq.StringArray `gorm:"type:text[]" json:"job_categories,omitempty"`
	HiringLocations pq.StringArray `gorm:"type:text[]" json:"hiring_locations,omitempty"`
	HiringTypes     pq.StringArray `gorm:"type:text[]" json:"hiring_types,omitempty"`

	// Optional business information
	GSTINNumber    string `json:"gstin_number,omitempty"`
	CompanyAddress string `json:"company_address,omitempty"`
	City           string `json:"city,omitempty"`
	State          string `json:"state,omitempty"`
	Pincode        string `json:"pincode,omitempty"`

	// System managed fields
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName specifies the database table name for EmployerProfile
func (EmployerProfile) TableName() string {
	return "employer_profiles"
}
