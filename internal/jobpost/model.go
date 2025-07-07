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
	ID                  string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Title               string         `json:"title"`
	RoleOverview        string         `json:"roleOverview"`
	Requirements        string         `json:"requirements"`
	Location            string         `json:"location"`
	RequiredSkills      pq.StringArray `gorm:"type:text[]" json:"requiredSkills"`
	EmployerID          string         `json:"employerId"`
	EmployerName        string         `json:"employerName"`
	EmployerEmail       string         `json:"employerEmail"`
	Status              string         `json:"status"` // draft, published, closed, completed
	CreatedAt           time.Time      `json:"createdAt"`
	UpdatedAt           time.Time      `json:"updatedAt"`
	ApplicationDeadline time.Time      `json:"applicationDeadline"`
	JobType             string         `json:"jobType"`                                      // full-time, part-time, contract, internship
	Experience          string         `json:"experience"`                                   // entry, mid, senior
	SalaryMin           float64        `json:"salaryMin" gorm:"column:salary_min"`           // Database column
	SalaryMax           float64        `json:"salaryMax" gorm:"column:salary_max"`           // Database column
	SalaryCurrency      string         `json:"salaryCurrency" gorm:"column:salary_currency"` // Database column
	Salary              Salary         `json:"salary" gorm:"-"`                              // Virtual field for JSON
	Benefits            pq.StringArray `gorm:"type:text[]" json:"benefits"`
	IsRemote            bool           `json:"isRemote"`
	ApplicationsCount   int            `json:"applicationsCount" gorm:"default:0"`
	CompletedAt         *time.Time     `json:"completedAt"`
	HiredCandidateName  *string        `json:"hiredCandidateName"`
}

// Request Models
type CreateJobPostRequest struct {
	Title               string         `json:"title" binding:"required"`
	RoleOverview        string         `json:"roleOverview" binding:"required"`
	Requirements        string         `json:"requirements" binding:"required"`
	Location            string         `json:"location" binding:"required"`
	RequiredSkills      []string       `json:"requiredSkills"`
	ApplicationDeadline time.Time      `json:"applicationDeadline"`
	JobType             string         `json:"jobType" binding:"required"`    // full-time, part-time, contract, internship
	Experience          string         `json:"experience" binding:"required"` // entry, mid, senior
	Salary              FlexibleSalary `json:"salary"`
	Benefits            []string       `json:"benefits"`
	IsRemote            bool           `json:"isRemote"`
}

// CreateDraftRequest - For creating draft jobs (all fields optional)
type CreateDraftRequest struct {
	Title               *string         `json:"title,omitempty"`
	RoleOverview        *string         `json:"roleOverview,omitempty"`
	Requirements        *string         `json:"requirements,omitempty"`
	Location            *string         `json:"location,omitempty"`
	RequiredSkills      []string        `json:"requiredSkills,omitempty"`
	ApplicationDeadline *time.Time      `json:"applicationDeadline,omitempty"`
	JobType             *string         `json:"jobType,omitempty"`
	Experience          *string         `json:"experience,omitempty"`
	Salary              *FlexibleSalary `json:"salary,omitempty"`
	Benefits            []string        `json:"benefits,omitempty"`
	IsRemote            *bool           `json:"isRemote,omitempty"`
}

type UpdateJobPostRequest struct {
	Title               *string         `json:"title"`
	RoleOverview        *string         `json:"roleOverview"`
	Requirements        *string         `json:"requirements"`
	Location            *string         `json:"location"`
	RequiredSkills      []string        `json:"requiredSkills"`
	ApplicationDeadline *time.Time      `json:"applicationDeadline"`
	JobType             *string         `json:"jobType"`
	Experience          *string         `json:"experience"`
	Salary              *FlexibleSalary `json:"salary"`
	Benefits            []string        `json:"benefits"`
	IsRemote            *bool           `json:"isRemote"`
}

// Response Models
type JobPostResponse struct {
	Success  bool      `json:"success"`
	Message  string    `json:"message"`
	JobPost  *JobPost  `json:"jobPost,omitempty"`
	JobPosts []JobPost `json:"jobPosts,omitempty"`
}

// Search Filter Models
type JobPostFilter struct {
	Location    string   `json:"location"`
	JobType     []string `json:"jobType"`
	Experience  []string `json:"experience"`
	SalaryRange *struct {
		Min float64 `json:"min"`
		Max float64 `json:"max"`
	} `json:"salaryRange"`
	IsRemote     *bool    `json:"isRemote"`
	Skills       []string `json:"skills"`
	PostedWithin string   `json:"postedWithin"` // all, 24h, 7d, 30d
	Page         int      `json:"page"`
	Limit        int      `json:"limit"`
}

// Enhanced Search Models
type AdvancedJobSearchRequest struct {
	// Basic filters
	Keywords    string   `json:"keywords"`    // Search in title, description, requirements
	Location    string   `json:"location"`    // City, state, or remote
	JobType     []string `json:"jobType"`     // full-time, part-time, contract, internship
	Experience  []string `json:"experience"`  // entry, mid, senior
	Skills      []string `json:"skills"`      // Required skills
	Industry    []string `json:"industry"`    // Company industry
	CompanySize []string `json:"companySize"` // Company size filters

	// Salary and benefits
	SalaryRange *struct {
		Min      float64 `json:"min"`
		Max      float64 `json:"max"`
		Currency string  `json:"currency"`
	} `json:"salaryRange"`
	Benefits []string `json:"benefits"` // Health insurance, 401k, etc.

	// Work preferences
	IsRemote *bool `json:"isRemote"` // Remote work preference
	IsHybrid *bool `json:"isHybrid"` // Hybrid work preference
	IsOnsite *bool `json:"isOnsite"` // On-site work preference

	// Timing filters
	PostedWithin string `json:"postedWithin"` // all, 24h, 7d, 30d, 90d
	Urgent       *bool  `json:"urgent"`       // Urgent hiring positions

	// Sorting and pagination
	SortBy    string `json:"sortBy"`    // relevance, date, salary, applications
	SortOrder string `json:"sortOrder"` // asc, desc
	Page      int    `json:"page"`      // Default: 1
	Limit     int    `json:"limit"`     // Default: 20, Max: 100
}

type JobSearchResponse struct {
	Success    bool            `json:"success"`
	Message    string          `json:"message"`
	Jobs       []JobPost       `json:"jobs"`
	Filters    *SearchFilters  `json:"filters,omitempty"`
	Pagination *PaginationInfo `json:"pagination,omitempty"`
}

type SearchFilters struct {
	AvailableLocations    []string      `json:"availableLocations"`
	AvailableJobTypes     []string      `json:"availableJobTypes"`
	AvailableExperience   []string      `json:"availableExperience"`
	AvailableSkills       []string      `json:"availableSkills"`
	AvailableIndustries   []string      `json:"availableIndustries"`
	AvailableCompanySizes []string      `json:"availableCompanySizes"`
	SalaryRanges          []SalaryRange `json:"salaryRanges"`
	AvailableBenefits     []string      `json:"availableBenefits"`
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
	TotalPages int  `json:"totalPages"`
	HasNext    bool `json:"hasNext"`
	HasPrev    bool `json:"hasPrev"`
}

// Job Recommendation Models
type JobRecommendationRequest struct {
	UserID            string   `json:"userId"`
	UserSkills        []string `json:"userSkills"`
	UserLocation      string   `json:"userLocation"`
	UserExperience    string   `json:"userExperience"`
	PreferredJobTypes []string `json:"preferredJobTypes"`
	MaxResults        int      `json:"maxResults"` // Default: 10
}

type JobRecommendationResponse struct {
	Success bool      `json:"success"`
	Message string    `json:"message"`
	Jobs    []JobPost `json:"jobs"`
	Reason  string    `json:"reason"` // Why these jobs were recommended
}

// Job Alert Models
type JobAlertRequest struct {
	UserID      string   `json:"userId"`
	Keywords    []string `json:"keywords"`
	Location    string   `json:"location"`
	JobType     []string `json:"jobType"`
	Experience  []string `json:"experience"`
	Skills      []string `json:"skills"`
	SalaryRange *struct {
		Min      float64 `json:"min"`
		Max      float64 `json:"max"`
		Currency string  `json:"currency"`
	} `json:"salaryRange"`
	IsRemote  *bool  `json:"isRemote"`
	Frequency string `json:"frequency"` // daily, weekly, immediate
	IsActive  bool   `json:"isActive"`
}

type JobAlertResponse struct {
	Success bool       `json:"success"`
	Message string     `json:"message"`
	Alert   *JobAlert  `json:"alert,omitempty"`
	Alerts  []JobAlert `json:"alerts,omitempty"`
}

type JobAlert struct {
	ID          string   `json:"id"`
	UserID      string   `json:"userId"`
	Keywords    []string `json:"keywords"`
	Location    string   `json:"location"`
	JobType     []string `json:"jobType"`
	Experience  []string `json:"experience"`
	Skills      []string `json:"skills"`
	SalaryRange *struct {
		Min      float64 `json:"min"`
		Max      float64 `json:"max"`
		Currency string  `json:"currency"`
	} `json:"salaryRange"`
	IsRemote  *bool     `json:"isRemote"`
	Frequency string    `json:"frequency"`
	IsActive  bool      `json:"isActive"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
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
	JobID      string `json:"jobId"`
	MaxResults int    `json:"maxResults"` // Default: 5
}

type SimilarJobsResponse struct {
	Success bool      `json:"success"`
	Message string    `json:"message"`
	Jobs    []JobPost `json:"jobs"`
}
