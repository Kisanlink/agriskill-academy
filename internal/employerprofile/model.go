package employerprofile

import (
	"github.com/Kisanlink/kisanlink-db/pkg/base"
	"github.com/Kisanlink/kisanlink-db/pkg/core/hash"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type EmployerProfile struct {
	base.BaseModel
	UserID string `gorm:"type:varchar(255);not null" json:"user_id" binding:"required"`

	// Required company information
	CompanyName string `gorm:"not null" json:"company_name" binding:"required"`
	Industry    string `gorm:"not null" json:"industry" binding:"required"`
	CompanySize string `gorm:"not null" json:"company_size" binding:"required"`

	// Optional company branding and details
	LogoKey            string `json:"logo_key,omitempty"`
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
}

// TableName specifies the database table name for EmployerProfile
func (EmployerProfile) TableName() string {
	return "employer_profiles"
}

// NewEmployerProfile creates a new EmployerProfile with proper initialization
func NewEmployerProfile() *EmployerProfile {
	return &EmployerProfile{
		BaseModel: *base.NewBaseModel("EMPL", hash.Medium),
	}
}

// BeforeCreateGORM is called by GORM before creating a new record
func (e *EmployerProfile) BeforeCreateGORM(tx *gorm.DB) error {
	return e.BaseModel.BeforeCreate()
}

// BeforeUpdateGORM is called by GORM before updating an existing record
func (e *EmployerProfile) BeforeUpdateGORM(tx *gorm.DB) error {
	return e.BaseModel.BeforeUpdate()
}

// BeforeDeleteGORM is called by GORM before hard deleting a record
func (e *EmployerProfile) BeforeDeleteGORM(tx *gorm.DB) error {
	return e.BaseModel.BeforeDelete()
}

// UpdateEmployerProfileRequest is used for update operations where UserID comes from URL parameter
type UpdateEmployerProfileRequest struct {
	// Company information (all optional for updates)
	CompanyName string `json:"company_name,omitempty"`
	Industry    string `json:"industry,omitempty"`
	CompanySize string `json:"company_size,omitempty"`

	// Optional company branding and details
	LogoKey            string `json:"logo_key,omitempty"`
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
	JobCategories   pq.StringArray `json:"job_categories,omitempty"`
	HiringLocations pq.StringArray `json:"hiring_locations,omitempty"`
	HiringTypes     pq.StringArray `json:"hiring_types,omitempty"`

	// Optional business information
	GSTINNumber    string `json:"gstin_number,omitempty"`
	CompanyAddress string `json:"company_address,omitempty"`
	City           string `json:"city,omitempty"`
	State          string `json:"state,omitempty"`
	Pincode        string `json:"pincode,omitempty"`
}
