package admin

import (
	"time"
)

// Analytics Response Models
type JobAnalytics struct {
	TotalJobs         int             `json:"totalJobs"`
	PublishedJobs     int             `json:"publishedJobs"`
	DraftJobs         int             `json:"draftJobs"`
	ClosedJobs        int             `json:"closedJobs"`
	TotalApplications int             `json:"totalApplications"`
	AvgApplications   float64         `json:"avgApplications"`
	MostPopularJob    string          `json:"mostPopularJob"`
	JobsThisMonth     int             `json:"jobsThisMonth"`
	JobsLastMonth     int             `json:"jobsLastMonth"`
	GrowthRate        float64         `json:"growthRate"`
	JobsByLocation    []LocationStats `json:"jobsByLocation"`
	JobsByType        []JobTypeStats  `json:"jobsByType"`
}

type UserAnalytics struct {
	TotalUsers        int             `json:"totalUsers"`
	TotalStudents     int             `json:"totalStudents"`
	TotalEmployers    int             `json:"totalEmployers"`
	ActiveUsers       int             `json:"activeUsers"`
	NewUsersThisMonth int             `json:"newUsersThisMonth"`
	NewUsersLastMonth int             `json:"newUsersLastMonth"`
	UserGrowthRate    float64         `json:"userGrowthRate"`
	TopLocations      []LocationStats `json:"topLocations"`
}

type ApplicationAnalytics struct {
	TotalApplications     int           `json:"totalApplications"`
	PendingApplications   int           `json:"pendingApplications"`
	AcceptedApplications  int           `json:"acceptedApplications"`
	RejectedApplications  int           `json:"rejectedApplications"`
	WithdrawnApplications int           `json:"withdrawnApplications"`
	ApplicationsThisMonth int           `json:"applicationsThisMonth"`
	ApplicationsLastMonth int           `json:"applicationsLastMonth"`
	ApplicationGrowthRate float64       `json:"applicationGrowthRate"`
	ApplicationsByStatus  []StatusStats `json:"applicationsByStatus"`
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
	JobType    string  `json:"jobType"`
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
	LastUpdated  time.Time            `json:"lastUpdated"`
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
	SortBy    string `form:"sortBy"`
	SortOrder string `form:"sortOrder"`
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
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UserDetailResponse struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Email     string      `json:"email"`
	Role      string      `json:"role"`
	Status    string      `json:"status"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
	Profile   interface{} `json:"profile,omitempty"`
}

type UpdateUserRequest struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	Status string `json:"status"`
}

type PaginationInfo struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"totalPages"`
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
	SortBy    string `form:"sortBy"`
	SortOrder string `form:"sortOrder"`
}

type CompanyListResponse struct {
	Companies  []CompanyListItem `json:"companies"`
	Pagination PaginationInfo    `json:"pagination"`
}

type CompanyListItem struct {
	ID                string    `json:"id"`
	CompanyName       string    `json:"companyName"`
	Industry          string    `json:"industry"`
	Location          string    `json:"location"`
	CompanySize       string    `json:"companySize"`
	Status            string    `json:"status"`
	JobsCount         int       `json:"jobsCount"`
	ApplicationsCount int       `json:"applicationsCount"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

type CompanyDetailResponse struct {
	ID                 string    `json:"id"`
	CompanyName        string    `json:"companyName"`
	Industry           string    `json:"industry"`
	CompanySize        string    `json:"companySize"`
	CompanyDescription string    `json:"companyDescription"`
	Location           string    `json:"location"`
	WebsiteUrl         string    `json:"websiteUrl"`
	RecruiterName      string    `json:"recruiterName"`
	OfficialEmail      string    `json:"officialEmail"`
	PhoneNumber        string    `json:"phoneNumber"`
	Status             string    `json:"status"`
	JobsCount          int       `json:"jobsCount"`
	ApplicationsCount  int       `json:"applicationsCount"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
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
	CompanyName        string `json:"companyName"`
	Industry           string `json:"industry"`
	CompanySize        string `json:"companySize"`
	CompanyDescription string `json:"companyDescription"`
	Location           string `json:"location"`
	WebsiteUrl         string `json:"websiteUrl"`
	RecruiterName      string `json:"recruiterName"`
	OfficialEmail      string `json:"officialEmail"`
	PhoneNumber        string `json:"phoneNumber"`
}

type CompanyAnalytics struct {
	TotalCompanies        int                `json:"totalCompanies"`
	ActiveCompanies       int                `json:"activeCompanies"`
	VerifiedCompanies     int                `json:"verifiedCompanies"`
	NewCompaniesThisMonth int                `json:"newCompaniesThisMonth"`
	NewCompaniesLastMonth int                `json:"newCompaniesLastMonth"`
	CompanyGrowthRate     float64            `json:"companyGrowthRate"`
	CompaniesByIndustry   []IndustryStats    `json:"companiesByIndustry"`
	CompaniesByLocation   []LocationStats    `json:"companiesByLocation"`
	CompaniesBySize       []CompanySizeStats `json:"companiesBySize"`
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
