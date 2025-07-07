package employerprofile

import (
	"time"

	"github.com/lib/pq"
)

type EmployerProfile struct {
	ID                 string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID             string         `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`
	CompanyName        string         `gorm:"column:company_name" json:"companyName"`
	Logo               string         `json:"logo"`
	WebsiteUrl         string         `gorm:"column:website_url" json:"websiteUrl"`
	Industry           string         `json:"industry"`
	CompanySize        string         `gorm:"column:company_size" json:"companySize"`
	CompanyDescription string         `gorm:"column:company_description" json:"companyDescription"`
	RecruiterName      string         `gorm:"column:recruiter_name" json:"recruiterName"`
	Designation        string         `json:"designation"`
	OfficialEmail      string         `gorm:"column:official_email" json:"officialEmail"`
	PhoneNumber        string         `gorm:"column:phone_number" json:"phoneNumber"`
	LinkedinProfile    string         `gorm:"column:linkedin_profile" json:"linkedinProfile"`
	GstinNumber        string         `gorm:"column:gstin_number" json:"gstinNumber"`
	CompanyAddress     string         `gorm:"column:company_address" json:"companyAddress"`
	City               string         `gorm:"column:city" json:"city"`
	State              string         `gorm:"column:state" json:"state"`
	Pincode            string         `gorm:"column:pincode" json:"pincode"`
	JobCategories      pq.StringArray `gorm:"type:text[];column:job_categories" json:"jobCategories"`
	HiringLocations    pq.StringArray `gorm:"type:text[];column:hiring_locations" json:"hiringLocations"`
	HiringTypes        pq.StringArray `gorm:"type:text[];column:hiring_types" json:"hiringTypes"`
	CreatedAt          time.Time      `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt          time.Time      `gorm:"column:updated_at" json:"updatedAt"`
}
