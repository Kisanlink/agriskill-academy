package admin

import (
	"time"
)

// Analytics Response Models
type JobAnalytics struct {
	TotalJobs         int             `json:"total_jobs"`
	PublishedJobs     int             `json:"published_jobs"`
	DraftJobs         int             `json:"draft_jobs"`
	ClosedJobs        int             `json:"closed_jobs"`
	TotalApplications int             `json:"total_applications"`
	AvgApplications   float64         `json:"avg_applications"`
	MostPopularJob    string          `json:"most_popular_job"`
	JobsThisMonth     int             `json:"jobs_this_month"`
	JobsLastMonth     int             `json:"jobs_last_month"`
	GrowthRate        float64         `json:"growth_rate"`
	JobsByLocation    []LocationStats `json:"jobs_by_location"`
	JobsByType        []JobTypeStats  `json:"jobs_by_type"`
}

type UserAnalytics struct {
	TotalUsers        int             `json:"total_users"`
	TotalStudents     int             `json:"total_students"`
	TotalEmployers    int             `json:"total_employers"`
	ActiveUsers       int             `json:"active_users"`
	NewUsersThisMonth int             `json:"new_users_this_month"`
	NewUsersLastMonth int             `json:"new_users_last_month"`
	UserGrowthRate    float64         `json:"user_growth_rate"`
	TopLocations      []LocationStats `json:"top_locations"`
}

type ApplicationAnalytics struct {
	TotalApplications     int           `json:"total_applications"`
	PendingApplications   int           `json:"pending_applications"`
	AcceptedApplications  int           `json:"accepted_applications"`
	RejectedApplications  int           `json:"rejected_applications"`
	WithdrawnApplications int           `json:"withdrawn_applications"`
	ApplicationsThisMonth int           `json:"applications_this_month"`
	ApplicationsLastMonth int           `json:"applications_last_month"`
	ApplicationGrowthRate float64       `json:"application_growth_rate"`
	ApplicationsByStatus  []StatusStats `json:"applications_by_status"`
}

type IndustryStats struct {
	Industry   string  `json:"industry"`
	Count      int     `json:"count"`
	Percentage float64 `json:"percentage"`
}

type LocationStats struct {
	Location   string  `json:"location"`
	Count      int     `json:"count"`
	Percentage float64 `json:"percentage"`
}

type JobTypeStats struct {
	JobType    string  `json:"job_type"`
	Count      int     `json:"count"`
	Percentage float64 `json:"percentage"`
}

type StatusStats struct {
	Status     string  `json:"status"`
	Count      int     `json:"count"`
	Percentage float64 `json:"percentage"`
}

type MonthlyStats struct {
	Month string `json:"month"`
	Count int    `json:"count"`
}

type DashboardAnalytics struct {
	Jobs         JobAnalytics         `json:"jobs"`
	Users        UserAnalytics        `json:"users"`
	Applications ApplicationAnalytics `json:"applications"`
	LastUpdated  time.Time            `json:"last_updated"`
}

type AnalyticsResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// User Management Models
type UserListRequest struct {
	Page      int    `form:"page" binding:"min=1"`
	Limit     int    `form:"limit" binding:"min=1,max=100"`
	Role      string `form:"role"`
	Search    string `form:"search"`
	Status    string `form:"status"`
	SortBy    string `form:"sort_by"`
	SortOrder string `form:"sort_order"`
}

type UserListResponse struct {
	Users      []UserListItem `json:"users"`
	Pagination PaginationInfo `json:"pagination"`
}

type UserListItem struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserDetailResponse struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Email     string      `json:"email"`
	Role      string      `json:"role"`
	Status    string      `json:"status"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	Profile   interface{} `json:"profile,omitempty"`
}

type UpdateUserRequest struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

type PaginationInfo struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type UserResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Company Management Models
type CompanyListRequest struct {
	Page      int    `form:"page" binding:"min=1"`
	Limit     int    `form:"limit" binding:"min=1,max=100"`
	Industry  string `form:"industry"`
	Search    string `form:"search"`
	Status    string `form:"status"`
	SortBy    string `form:"sort_by"`
	SortOrder string `form:"sort_order"`
}

type CompanyListResponse struct {
	Companies  []CompanyListItem `json:"companies"`
	Pagination PaginationInfo    `json:"pagination"`
}

type CompanyListItem struct {
	ID                string    `json:"id"`
	CompanyName       string    `json:"company_name"`
	Industry          string    `json:"industry"`
	Location          string    `json:"location"`
	CompanySize       string    `json:"company_size"`
	Status            string    `json:"status"`
	JobsCount         int       `json:"jobs_count"`
	ApplicationsCount int       `json:"applications_count"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type CompanyDetailResponse struct {
	ID                 string    `json:"id"`
	UserID             string    `json:"user_id"`
	CompanyName        string    `json:"company_name"`
	Industry           string    `json:"industry"`
	CompanySize        string    `json:"company_size"`
	CompanyDescription string    `json:"company_description"`
	Location           string    `json:"location"`
	WebsiteUrl         string    `json:"website_url"`
	RecruiterName      string    `json:"recruiter_name"`
	OfficialEmail      string    `json:"official_email"`
	PhoneNumber        string    `json:"phone_number"`
	Status             string    `json:"status"`
	JobsCount          int       `json:"jobs_count"`
	ApplicationsCount  int       `json:"applications_count"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	User               *UserInfo `json:"user,omitempty"`
}

type UserInfo struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type UpdateCompanyRequest struct {
	Status             string `json:"status"`
	Industry           string `json:"industry"`
	CompanySize        string `json:"company_size"`
	CompanyDescription string `json:"company_description"`
	Location           string `json:"location"`
	WebsiteUrl         string `json:"website_url"`
	RecruiterName      string `json:"recruiter_name"`
	OfficialEmail      string `json:"official_email"`
	PhoneNumber        string `json:"phone_number"`
}

type CompanyAnalytics struct {
	TotalCompanies        int                `json:"total_companies"`
	ActiveCompanies       int                `json:"active_companies"`
	VerifiedCompanies     int                `json:"verified_companies"`
	NewCompaniesThisMonth int                `json:"new_companies_this_month"`
	NewCompaniesLastMonth int                `json:"new_companies_last_month"`
	CompanyGrowthRate     float64            `json:"company_growth_rate"`
	CompaniesByIndustry   []IndustryStats    `json:"companies_by_industry"`
	CompaniesByLocation   []LocationStats    `json:"companies_by_location"`
	CompaniesBySize       []CompanySizeStats `json:"companies_by_size"`
}

type CompanySizeStats struct {
	Size       string  `json:"size"`
	Count      int     `json:"count"`
	Percentage float64 `json:"percentage"`
}

type CompanyResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
