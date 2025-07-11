package employerprofile

import (
	"time"

	"github.com/lib/pq"
)

type EmployerProfile struct {
	ID                 string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID             string         `json:"user_id" binding:"required"`
	CompanyName        string         `json:"company_name"`
	Logo               []byte         `json:"logo" gorm:"type:bytea"`
	LogoName           string         `json:"logo_name"`
	LogoType           string         `json:"logo_type"`
	LogoSize           int64          `json:"logo_size"`
	WebsiteURL         string         `json:"website_url"`
	Industry           string         `json:"industry"`
	CompanySize        string         `json:"company_size"`
	CompanyDescription string         `json:"company_description"`
	RecruiterName      string         `json:"recruiter_name"`
	Designation        string         `json:"designation"`
	OfficialEmail      string         `json:"official_email"`
	PhoneNumber        string         `json:"phone_number"`
	LinkedinProfile    string         `json:"linkedin_profile"`
	JobCategories      pq.StringArray `gorm:"type:text[]" json:"job_categories"`
	HiringLocations    pq.StringArray `gorm:"type:text[]" json:"hiring_locations"`
	HiringTypes        pq.StringArray `gorm:"type:text[]" json:"hiring_types"`
	GSTINNumber        string         `json:"gstin_number"`
	CompanyAddress     string         `json:"company_address"`
	City               string         `json:"city"`
	State              string         `json:"state"`
	Pincode            string         `json:"pincode"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
}

// TableName specifies the database table name for EmployerProfile
func (EmployerProfile) TableName() string {
	return "employer_profiles"
}
