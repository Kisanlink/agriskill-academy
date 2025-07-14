package jobpost

import (
	"encoding/json"
	"time"

	"github.com/lib/pq"
)

type Salary struct {
	Min      float64 `json:"min"`
	Max      float64 `json:"max"`
	Currency string  `json:"currency"`
}

// FlexibleSalary handles both object and number formats
type FlexibleSalary struct {
	*Salary
}

// UnmarshalJSON handles both object and number formats for salary
func (fs *FlexibleSalary) UnmarshalJSON(data []byte) error {
	// First try to unmarshal as a number
	var num float64
	if err := json.Unmarshal(data, &num); err == nil {
		// If it's a number, create a salary object with the number as both min and max
		fs.Salary = &Salary{
			Min:      num,
			Max:      num,
			Currency: "USD", // Default currency
		}
		return nil
	}

	// If it's not a number, try to unmarshal as a salary object
	var salary Salary
	if err := json.Unmarshal(data, &salary); err != nil {
		return err
	}
	fs.Salary = &salary
	return nil
}

// MarshalJSON ensures we always output the salary object format
func (fs *FlexibleSalary) MarshalJSON() ([]byte, error) {
	if fs.Salary == nil {
		return []byte("null"), nil
	}
	return json.Marshal(fs.Salary)
}

type JobPost struct {
	ID                  string         `gorm:"primaryKey;type:varchar(255)" json:"id"`
	Title               string         `json:"title"`
	RoleOverview        string         `json:"role_overview"`
	Requirements        string         `json:"requirements"`
	Location            string         `json:"location"`
	RequiredSkills      pq.StringArray `gorm:"type:text[]" json:"required_skills"`
	EmployerID          string         `gorm:"type:varchar(255)" json:"employer_id"`
	EmployerName        string         `json:"employer_name"`
	EmployerEmail       string         `json:"employer_email"`
	Status              string         `json:"status"` // draft, published, closed, completed
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
	ApplicationDeadline time.Time      `json:"application_deadline"`
	JobType             string         `json:"job_type"`                                      // full-time, part-time, contract, internship
	Experience          string         `json:"experience"`                                    // entry, mid, senior
	SalaryMin           float64        `json:"salary_min" gorm:"column:salary_min"`           // Database column
	SalaryMax           float64        `json:"salary_max" gorm:"column:salary_max"`           // Database column
	SalaryCurrency      string         `json:"salary_currency" gorm:"column:salary_currency"` // Database column
	Salary              Salary         `json:"salary" gorm:"-"`                               // Virtual field for JSON
	Benefits            pq.StringArray `gorm:"type:text[]" json:"benefits"`
	IsRemote            bool           `json:"is_remote"`
	ApplicationsCount   int            `json:"applications_count" gorm:"default:0"`
	CompletedAt         *time.Time     `json:"completed_at"`
	HiredCandidateName  *string        `json:"hired_candidate_name"`
}

// Request Models
type CreateJobPostRequest struct {
	Title               string         `json:"title" binding:"required"`
	RoleOverview        string         `json:"role_overview" binding:"required"`
	Requirements        string         `json:"requirements" binding:"required"`
	Location            string         `json:"location" binding:"required"`
	RequiredSkills      []string       `json:"required_skills"`
	ApplicationDeadline time.Time      `json:"application_deadline"`
	JobType             string         `json:"job_type" binding:"required"`   // full-time, part-time, contract, internship
	Experience          string         `json:"experience" binding:"required"` // entry, mid, senior
	Salary              FlexibleSalary `json:"salary"`
	Benefits            []string       `json:"benefits"`
	IsRemote            bool           `json:"is_remote"`
}

// CreateDraftRequest - For creating draft jobs (all fields optional)
type CreateDraftRequest struct {
	Title               *string         `json:"title,omitempty"`
	RoleOverview        *string         `json:"role_overview,omitempty"`
	Requirements        *string         `json:"requirements,omitempty"`
	Location            *string         `json:"location,omitempty"`
	RequiredSkills      []string        `json:"required_skills,omitempty"`
	ApplicationDeadline *time.Time      `json:"application_deadline,omitempty"`
	JobType             *string         `json:"job_type,omitempty"`
	Experience          *string         `json:"experience,omitempty"`
	Salary              *FlexibleSalary `json:"salary,omitempty"`
	Benefits            []string        `json:"benefits,omitempty"`
	IsRemote            *bool           `json:"is_remote,omitempty"`
}

type UpdateJobPostRequest struct {
	Title               *string         `json:"title"`
	RoleOverview        *string         `json:"role_overview"`
	Requirements        *string         `json:"requirements"`
	Location            *string         `json:"location"`
	RequiredSkills      []string        `json:"required_skills"`
	ApplicationDeadline *time.Time      `json:"application_deadline"`
	JobType             *string         `json:"job_type"`
	Experience          *string         `json:"experience"`
	Salary              *FlexibleSalary `json:"salary"`
	Benefits            []string        `json:"benefits"`
	IsRemote            *bool           `json:"is_remote"`
}

// Response Models
type JobPostResponse struct {
	Success  bool      `json:"success"`
	Message  string    `json:"message"`
	JobPost  *JobPost  `json:"job_post,omitempty"`
	JobPosts []JobPost `json:"job_posts,omitempty"`
}

// Search Filter Models
type JobPostFilter struct {
	Location    string   `json:"location"`
	JobType     []string `json:"job_type"`
	Experience  []string `json:"experience"`
	SalaryRange *struct {
		Min float64 `json:"min"`
		Max float64 `json:"max"`
	} `json:"salary_range"`
	IsRemote     *bool    `json:"is_remote"`
	Skills       []string `json:"skills"`
	PostedWithin string   `json:"posted_within"` // all, 24h, 7d, 30d
	Page         int      `json:"page"`
	Limit        int      `json:"limit"`
}

// Enhanced Search Models
type AdvancedJobSearchRequest struct {
	// Basic filters
	Keywords    string   `json:"keywords"`     // Search in title, description, requirements
	Location    string   `json:"location"`     // City, state, or remote
	JobType     []string `json:"job_type"`     // full-time, part-time, contract, internship
	Experience  []string `json:"experience"`   // entry, mid, senior
	Skills      []string `json:"skills"`       // Required skills
	Industry    []string `json:"industry"`     // Company industry
	CompanySize []string `json:"company_size"` // Company size filters

	// Salary and benefits
	SalaryRange *struct {
		Min      float64 `json:"min"`
		Max      float64 `json:"max"`
		Currency string  `json:"currency"`
	} `json:"salary_range"`
	Benefits []string `json:"benefits"` // Health insurance, 401k, etc.

	// Work preferences
	IsRemote *bool `json:"is_remote"` // Remote work preference
	IsHybrid *bool `json:"is_hybrid"` // Hybrid work preference
	IsOnsite *bool `json:"is_onsite"` // On-site work preference

	// Timing filters
	PostedWithin string `json:"posted_within"` // all, 24h, 7d, 30d, 90d
	Urgent       *bool  `json:"urgent"`        // Urgent hiring positions

	// Sorting and pagination
	SortBy    string `json:"sort_by"`    // relevance, date, salary, applications
	SortOrder string `json:"sort_order"` // asc, desc
	Page      int    `json:"page"`       // Default: 1
	Limit     int    `json:"limit"`      // Default: 20, Max: 100
}

type JobSearchResponse struct {
	Success    bool            `json:"success"`
	Message    string          `json:"message"`
	Jobs       []JobPost       `json:"jobs"`
	Filters    *SearchFilters  `json:"filters,omitempty"`
	Pagination *PaginationInfo `json:"pagination,omitempty"`
}

type SearchFilters struct {
	AvailableLocations    []string      `json:"available_locations"`
	AvailableJobTypes     []string      `json:"available_job_types"`
	AvailableExperience   []string      `json:"available_experience"`
	AvailableSkills       []string      `json:"available_skills"`
	AvailableIndustries   []string      `json:"available_industries"`
	AvailableCompanySizes []string      `json:"available_company_sizes"`
	SalaryRanges          []SalaryRange `json:"salary_ranges"`
	AvailableBenefits     []string      `json:"available_benefits"`
}

type SalaryRange struct {
	Min      float64 `json:"min"`
	Max      float64 `json:"max"`
	Currency string  `json:"currency"`
	Label    string  `json:"label"` // e.g., "$40k-$60k"
}

type PaginationInfo struct {
	Page       int  `json:"page"`
	Limit      int  `json:"limit"`
	Total      int  `json:"total"`
	TotalPages int  `json:"total_pages"`
	HasNext    bool `json:"has_next"`
	HasPrev    bool `json:"has_prev"`
}

// Job Recommendation Models
type JobRecommendationRequest struct {
	UserID            string   `json:"user_id"`
	UserSkills        []string `json:"user_skills"`
	UserLocation      string   `json:"user_location"`
	UserExperience    string   `json:"user_experience"`
	PreferredJobTypes []string `json:"preferred_job_types"`
	MaxResults        int      `json:"max_results"` // Default: 10
}

type JobRecommendationResponse struct {
	Success bool      `json:"success"`
	Message string    `json:"message"`
	Jobs    []JobPost `json:"jobs"`
	Reason  string    `json:"reason"` // Why these jobs were recommended
}

// Job Alert Models
type JobAlertRequest struct {
	UserID      string   `json:"user_id"`
	Keywords    []string `json:"keywords"`
	Location    string   `json:"location"`
	JobType     []string `json:"job_type"`
	Experience  []string `json:"experience"`
	Skills      []string `json:"skills"`
	SalaryRange *struct {
		Min      float64 `json:"min"`
		Max      float64 `json:"max"`
		Currency string  `json:"currency"`
	} `json:"salary_range"`
	IsRemote  *bool  `json:"is_remote"`
	Frequency string `json:"frequency"` // daily, weekly, immediate
	IsActive  bool   `json:"is_active"`
}

type JobAlertResponse struct {
	Success bool       `json:"success"`
	Message string     `json:"message"`
	Alert   *JobAlert  `json:"alert,omitempty"`
	Alerts  []JobAlert `json:"alerts,omitempty"`
}

// JobAlert represents the job_alerts table structure
type JobAlert struct {
	ID             string         `gorm:"primaryKey;type:varchar(255)" json:"id"`
	UserID         string         `gorm:"type:varchar(255);not null" json:"user_id"`
	Keywords       pq.StringArray `gorm:"type:text[]" json:"keywords"`
	Location       string         `gorm:"type:varchar(255)" json:"location"`
	JobType        pq.StringArray `gorm:"type:text[]" json:"job_type"`
	Experience     pq.StringArray `gorm:"type:text[]" json:"experience"`
	Skills         pq.StringArray `gorm:"type:text[]" json:"skills"`
	SalaryMin      *float64       `gorm:"type:numeric(10,2)" json:"salary_min"`
	SalaryMax      *float64       `gorm:"type:numeric(10,2)" json:"salary_max"`
	SalaryCurrency string         `gorm:"type:varchar(10);default:'USD'" json:"salary_currency"`
	IsRemote       *bool          `json:"is_remote"`
	Frequency      string         `gorm:"type:varchar(20);not null;default:'weekly';check:frequency IN ('daily', 'weekly', 'immediate')" json:"frequency"`
	IsActive       bool           `gorm:"not null;default:true" json:"is_active"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName specifies the database table name for JobAlert
func (JobAlert) TableName() string {
	return "job_alerts"
}

// Job Discovery Models
type FeaturedJobsResponse struct {
	Success bool      `json:"success"`
	Message string    `json:"message"`
	Jobs    []JobPost `json:"jobs"`
}

type RecentJobsResponse struct {
	Success bool      `json:"success"`
	Message string    `json:"message"`
	Jobs    []JobPost `json:"jobs"`
}

type TrendingJobsResponse struct {
	Success bool      `json:"success"`
	Message string    `json:"message"`
	Jobs    []JobPost `json:"jobs"`
}

type SimilarJobsRequest struct {
	JobID      string `json:"job_id"`
	MaxResults int    `json:"max_results"` // Default: 5
}

type SimilarJobsResponse struct {
	Success bool      `json:"success"`
	Message string    `json:"message"`
	Jobs    []JobPost `json:"jobs"`
}
