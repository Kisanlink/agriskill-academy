package employerprofile

import (
	"time"

	"github.com/lib/pq"
)

type EmployerProfile struct {
	ID                 string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID             string         `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`
	CompanyName        string         `gorm:"column:company_name" json:"company_name"`
	Logo               string         `json:"logo"`
	WebsiteUrl         string         `gorm:"column:website_url" json:"website_url"`
	Industry           string         `json:"industry"`
	CompanySize        string         `gorm:"column:company_size" json:"company_size"`
	CompanyDescription string         `gorm:"column:company_description" json:"company_description"`
	RecruiterName      string         `gorm:"column:recruiter_name" json:"recruiter_name"`
	Designation        string         `json:"designation"`
	OfficialEmail      string         `gorm:"column:official_email" json:"official_email"`
	PhoneNumber        string         `gorm:"column:phone_number" json:"phone_number"`
	LinkedinProfile    string         `gorm:"column:linkedin_profile" json:"linkedin_profile"`
	GstinNumber        string         `gorm:"column:gstin_number" json:"gstin_number"`
	CompanyAddress     string         `gorm:"column:company_address" json:"company_address"`
	City               string         `gorm:"column:city" json:"city"`
	State              string         `gorm:"column:state" json:"state"`
	Pincode            string         `gorm:"column:pincode" json:"pincode"`
	JobCategories      pq.StringArray `gorm:"type:text[];column:job_categories" json:"job_categories"`
	HiringLocations    pq.StringArray `gorm:"type:text[];column:hiring_locations" json:"hiring_locations"`
	HiringTypes        pq.StringArray `gorm:"type:text[];column:hiring_types" json:"hiring_types"`
	CreatedAt          time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt          time.Time      `gorm:"column:updated_at" json:"updated_at"`
}
