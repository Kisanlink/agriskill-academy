// File: internal/jobpost/model.go

package jobpost

import (
	"time"

	"github.com/lib/pq"
)

type Salary struct {
	Min      float64 `json:"min"`
	Max      float64 `json:"max"`
	Currency string  `json:"currency"`
}

type JobPost struct {
	ID                  string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Title               string         `json:"title"`
	RoleOverview        string         `json:"roleOverview"`
	Requirements        string         `json:"requirements"`
	Location            string         `json:"location"`
	RequiredSkills      pq.StringArray `gorm:"type:text[]" json:"requiredSkills"`
	EmployerID          string         `json:"employerId"` // FK to employer_profiles.id
	EmployerName        string         `json:"employerName"`
	EmployerEmail       string         `json:"employerEmail"`
	Status              string         `json:"status"` // draft, published, closed, completed
	CreatedAt           time.Time      `json:"createdAt"`
	UpdatedAt           time.Time      `json:"updatedAt"`
	ApplicationDeadline string         `json:"applicationDeadline"`
	JobType             string         `json:"jobType"`    // full-time, part-time, contract, internship
	Experience          string         `json:"experience"` // entry, mid, senior
	SalaryMin           float64        `json:"salaryMin"`
	SalaryMax           float64        `json:"salaryMax"`
	SalaryCurrency      string         `json:"salaryCurrency"`
	CompletedAt         *time.Time     `json:"completedAt"`
	HiredCandidateName  *string        `json:"hiredCandidateName"`
}
