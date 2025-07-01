// File: internal/employerprofile/model.go

package employerprofile

import (
	"time"

	"github.com/lib/pq"
)

type EmployerProfile struct {
	ID                 string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID             string         `gorm:"type:uuid;not null" json:"user_id"` // FK to users
	CompanyName        string         `json:"companyName"`
	Logo               string         `json:"logo"`
	WebsiteUrl         string         `json:"websiteUrl"`
	Industry           string         `json:"industry"`
	CompanySize        string         `json:"companySize"`
	CompanyDescription string         `json:"companyDescription"`
	RecruiterName      string         `json:"recruiterName"`
	Designation        string         `json:"designation"`
	OfficialEmail      string         `json:"officialEmail"`
	PhoneNumber        string         `json:"phoneNumber"`
	LinkedinProfile    string         `json:"linkedinProfile"`
	JobCategories      pq.StringArray `gorm:"type:text[]" json:"jobCategories"`
	HiringLocations    pq.StringArray `gorm:"type:text[]" json:"hiringLocations"`
	HiringTypes        pq.StringArray `gorm:"type:text[]" json:"hiringTypes"`
	CreatedAt          time.Time      `json:"createdAt"`
	UpdatedAt          time.Time      `json:"updatedAt"`
}
