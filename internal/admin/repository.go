package admin

import (
	"errors"
	"time"

	"github.com/Kisanlink/agriskill-academy/internal/auth"

	"gorm.io/gorm"
)

// contains checks if a slice of strings contains a specific string
func contains(list []string, val string) bool {
	for _, item := range list {
		if item == val {
			return true
		}
	}
	return false
}

type AdminRepository interface {
	GetJobAnalytics() (*JobAnalytics, error)
	GetUserAnalytics() (*UserAnalytics, error)
	GetApplicationAnalytics() (*ApplicationAnalytics, error)
	GetDashboardAnalytics() (*DashboardAnalytics, error)

	// User Management
	GetUsers(req *UserListRequest) (*UserListResponse, error)
	GetUserByID(userID string) (*UserDetailResponse, error)
	UpdateUser(userID string, req *UpdateUserRequest) error
	DeleteUser(userID string) error
	CreateAdmin(user *auth.User) error

	// Company Management
	GetCompanies(req *CompanyListRequest) (*CompanyListResponse, error)
	GetCompanyByID(companyID string) (*CompanyDetailResponse, error)
	UpdateCompany(companyID string, req *UpdateCompanyRequest) error
	DeleteCompany(companyID string) error
	GetCompanyAnalytics() (*CompanyAnalytics, error)

	// Student/Employer Lists
	GetStudents(req *StudentListRequest) (*StudentListResponse, error)
	GetEmployers(req *EmployerListRequest) (*EmployerListResponse, error)

	// Job Viewing (Admin-only)
	GetAllJobs(req *JobListRequest) (*JobListResponse, error)
	GetJobByID(jobID string) (*JobDetailResponse, error)
	GetJobStatistics() (*JobStatistics, error)
}

type adminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) AdminRepository {
	return &adminRepository{db}
}

func (r *adminRepository) GetJobAnalytics() (*JobAnalytics, error) {
	var analytics JobAnalytics

	// Get total jobs
	var totalJobs int64
	r.db.Model(&JobPost{}).Count(&totalJobs)
	analytics.TotalJobs = int(totalJobs)

	// Get jobs by status
	var publishedJobs, draftJobs, closedJobs int64
	r.db.Model(&JobPost{}).Where("status = ?", "published").Count(&publishedJobs)
	r.db.Model(&JobPost{}).Where("status = ?", "draft").Count(&draftJobs)
	r.db.Model(&JobPost{}).Where("status = ?", "closed").Count(&closedJobs)
	analytics.PublishedJobs = int(publishedJobs)
	analytics.DraftJobs = int(draftJobs)
	analytics.ClosedJobs = int(closedJobs)

	// Get total applications
	var totalApplications int64
	r.db.Model(&Application{}).Count(&totalApplications)
	analytics.TotalApplications = int(totalApplications)

	// Calculate average applications per job
	if analytics.TotalJobs > 0 {
		analytics.AvgApplications = float64(analytics.TotalApplications) / float64(analytics.TotalJobs)
	}

	// Get most popular job
	var mostPopularJob struct {
		Title string
		Count int
	}
	r.db.Model(&JobPost{}).
		Select("title, applications_count as count").
		Order("applications_count DESC").
		Limit(1).
		Scan(&mostPopularJob)
	analytics.MostPopularJob = mostPopularJob.Title

	// Get jobs this month and last month
	now := time.Now()
	thisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	lastMonth := thisMonth.AddDate(0, -1, 0)

	var jobsThisMonth, jobsLastMonth int64
	r.db.Model(&JobPost{}).Where("created_at >= ?", thisMonth).Count(&jobsThisMonth)
	r.db.Model(&JobPost{}).Where("created_at >= ? AND created_at < ?", lastMonth, thisMonth).Count(&jobsLastMonth)
	analytics.JobsThisMonth = int(jobsThisMonth)
	analytics.JobsLastMonth = int(jobsLastMonth)

	// Calculate growth rate
	if analytics.JobsLastMonth > 0 {
		analytics.GrowthRate = float64(analytics.JobsThisMonth-analytics.JobsLastMonth) / float64(analytics.JobsLastMonth) * 100
	}

	// Get jobs by location
	var locationStats []LocationStats
	r.db.Model(&JobPost{}).
		Select("location, COUNT(*) as count").
		Where("location IS NOT NULL AND location != ''").
		Group("location").
		Order("count DESC").
		Limit(10).
		Scan(&locationStats)

	// Calculate percentages
	for i := range locationStats {
		if analytics.TotalJobs > 0 {
			locationStats[i].Percentage = float64(locationStats[i].Count) / float64(analytics.TotalJobs) * 100
		}
	}
	analytics.JobsByLocation = locationStats

	// Get jobs by type
	var jobTypeStats []JobTypeStats
	r.db.Model(&JobPost{}).
		Select("job_type, COUNT(*) as count").
		Where("job_type IS NOT NULL AND job_type != ''").
		Group("job_type").
		Order("count DESC").
		Scan(&jobTypeStats)

	// Calculate percentages
	for i := range jobTypeStats {
		if analytics.TotalJobs > 0 {
			jobTypeStats[i].Percentage = float64(jobTypeStats[i].Count) / float64(analytics.TotalJobs) * 100
		}
	}
	analytics.JobsByType = jobTypeStats

	return &analytics, nil
}

func (r *adminRepository) GetUserAnalytics() (*UserAnalytics, error) {
	var analytics UserAnalytics

	// Get total users
	var totalUsers int64
	r.db.Model(&User{}).Count(&totalUsers)
	analytics.TotalUsers = int(totalUsers)

	// Get users by role
	var totalStudents, totalEmployers int64
	r.db.Model(&User{}).Where("role = ?", "student").Count(&totalStudents)
	r.db.Model(&User{}).Where("role = ?", "employer").Count(&totalEmployers)
	analytics.TotalStudents = int(totalStudents)
	analytics.TotalEmployers = int(totalEmployers)

	// Get active users (users who logged in within last 30 days)
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	var activeUsers int64
	r.db.Model(&User{}).Where("updated_at >= ?", thirtyDaysAgo).Count(&activeUsers)
	analytics.ActiveUsers = int(activeUsers)

	// Get new users this month and last month
	now := time.Now()
	thisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	lastMonth := thisMonth.AddDate(0, -1, 0)

	var newUsersThisMonth, newUsersLastMonth int64
	r.db.Model(&User{}).Where("created_at >= ?", thisMonth).Count(&newUsersThisMonth)
	r.db.Model(&User{}).Where("created_at >= ? AND created_at < ?", lastMonth, thisMonth).Count(&newUsersLastMonth)
	analytics.NewUsersThisMonth = int(newUsersThisMonth)
	analytics.NewUsersLastMonth = int(newUsersLastMonth)

	// Calculate growth rate
	if analytics.NewUsersLastMonth > 0 {
		analytics.UserGrowthRate = float64(analytics.NewUsersThisMonth-analytics.NewUsersLastMonth) / float64(analytics.NewUsersLastMonth) * 100
	}

	// Get top locations (from user profiles)
	var locationStats []LocationStats
	r.db.Model(&StudentProfile{}).
		Select("location, COUNT(*) as count").
		Where("location IS NOT NULL AND location != ''").
		Group("location").
		Order("count DESC").
		Limit(10).
		Scan(&locationStats)

	// Calculate percentages
	for i := range locationStats {
		if analytics.TotalUsers > 0 {
			locationStats[i].Percentage = float64(locationStats[i].Count) / float64(analytics.TotalUsers) * 100
		}
	}
	analytics.TopLocations = locationStats

	return &analytics, nil
}

func (r *adminRepository) GetApplicationAnalytics() (*ApplicationAnalytics, error) {
	var analytics ApplicationAnalytics

	// Get total applications
	var totalApplications int64
	r.db.Model(&Application{}).Count(&totalApplications)
	analytics.TotalApplications = int(totalApplications)

	// Get applications by status
	var pendingApplications, acceptedApplications, rejectedApplications, withdrawnApplications int64
	r.db.Model(&Application{}).Where("status = ?", "applied").Count(&pendingApplications)
	r.db.Model(&Application{}).Where("status = ?", "accepted").Count(&acceptedApplications)
	r.db.Model(&Application{}).Where("status = ?", "rejected").Count(&rejectedApplications)
	r.db.Model(&Application{}).Where("status = ?", "withdrawn").Count(&withdrawnApplications)
	analytics.PendingApplications = int(pendingApplications)
	analytics.AcceptedApplications = int(acceptedApplications)
	analytics.RejectedApplications = int(rejectedApplications)
	analytics.WithdrawnApplications = int(withdrawnApplications)

	// Get applications this month and last month
	now := time.Now()
	thisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	lastMonth := thisMonth.AddDate(0, -1, 0)

	var applicationsThisMonth, applicationsLastMonth int64
	r.db.Model(&Application{}).Where("applied_at >= ?", thisMonth).Count(&applicationsThisMonth)
	r.db.Model(&Application{}).Where("applied_at >= ? AND applied_at < ?", lastMonth, thisMonth).Count(&applicationsLastMonth)
	analytics.ApplicationsThisMonth = int(applicationsThisMonth)
	analytics.ApplicationsLastMonth = int(applicationsLastMonth)

	// Calculate growth rate
	if analytics.ApplicationsLastMonth > 0 {
		analytics.ApplicationGrowthRate = float64(analytics.ApplicationsThisMonth-analytics.ApplicationsLastMonth) / float64(analytics.ApplicationsLastMonth) * 100
	}

	// Get applications by status
	var statusStats []StatusStats
	r.db.Model(&Application{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Order("count DESC").
		Scan(&statusStats)

	// Calculate percentages
	for i := range statusStats {
		if analytics.TotalApplications > 0 {
			statusStats[i].Percentage = float64(statusStats[i].Count) / float64(analytics.TotalApplications) * 100
		}
	}
	analytics.ApplicationsByStatus = statusStats

	return &analytics, nil
}

func (r *adminRepository) GetDashboardAnalytics() (*DashboardAnalytics, error) {
	jobAnalytics, err := r.GetJobAnalytics()
	if err != nil {
		return nil, err
	}

	userAnalytics, err := r.GetUserAnalytics()
	if err != nil {
		return nil, err
	}

	applicationAnalytics, err := r.GetApplicationAnalytics()
	if err != nil {
		return nil, err
	}

	return &DashboardAnalytics{
		Jobs:         *jobAnalytics,
		Users:        *userAnalytics,
		Applications: *applicationAnalytics,
		LastUpdated:  time.Now(),
	}, nil
}

// Local struct definitions for database queries
type JobPost struct {
	ID                string    `gorm:"primaryKey;type:varchar(255)"`
	Title             string    `json:"title"`
	Status            string    `json:"status"`
	Location          string    `json:"location"`
	JobType           string    `json:"job_type"`
	ApplicationsCount int       `json:"applications_count"`
	CreatedAt         time.Time `json:"created_at"`
}

type Application struct {
	ID        string    `gorm:"primaryKey;type:varchar(255)"`
	Status    string    `json:"status"`
	AppliedAt time.Time `json:"applied_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type User struct {
	ID        string    `gorm:"primaryKey;type:varchar(255)"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type EmployerProfile struct {
	ID          string `gorm:"primaryKey;type:varchar(255)"`
	UserID      string `gorm:"column:user_id"`
	CompanyName string `json:"company_name"`
	Industry    string `json:"industry"`
	Location    string `json:"location"`
}

type Bookmark struct {
	ID     string `gorm:"primaryKey;type:varchar(255)"`
	UserID string `json:"user_id"`
}

type StudentProfile struct {
	ID       string `gorm:"primaryKey;type:varchar(255)"`
	Location string `json:"location"`
	Skills   string `json:"skills"`
}

// TableName specifies the database table name for StudentProfile
func (StudentProfile) TableName() string {
	return "student_profiles"
}

func (r *adminRepository) GetUsers(req *UserListRequest) (*UserListResponse, error) {
	var users []UserListItem
	var total int64

	query := r.db.Model(&User{})

	// Apply filters
	if req.Role != "" {
		query = query.Where("role = ?", req.Role)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.Search != "" {
		searchTerm := "%" + req.Search + "%"
		query = query.Where("name ILIKE ? OR email ILIKE ?", searchTerm, searchTerm)
	}

	// Get total count
	query.Count(&total)

	// Apply sorting
	if req.SortBy != "" {
		order := req.SortBy
		if req.SortOrder == "desc" {
			order += " DESC"
		}
		query = query.Order(order)
	} else {
		query = query.Order("created_at DESC")
	}

	// Apply pagination
	offset := (req.Page - 1) * req.Limit
	query = query.Offset(offset).Limit(req.Limit)

	// Execute query
	err := query.Find(&users).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((total + int64(req.Limit) - 1) / int64(req.Limit))

	return &UserListResponse{
		Users: users,
		Pagination: PaginationInfo{
			Page:       req.Page,
			Limit:      req.Limit,
			Total:      int(total),
			TotalPages: totalPages,
		},
	}, nil
}

func (r *adminRepository) GetUserByID(userID string) (*UserDetailResponse, error) {
	var user UserDetailResponse

	err := r.db.Model(&User{}).
		Select("id, name, email, role, status, created_at, updated_at").
		Where("id = ?", userID).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	// Get user profile based on role
	if user.Role == "student" {
		var profile StudentProfile
		err = r.db.Where("user_id = ?", userID).First(&profile).Error
		if err == nil {
			user.Profile = profile
		}
	} else if user.Role == "employer" {
		var profile EmployerProfile
		err = r.db.Where("user_id = ?", userID).First(&profile).Error
		if err == nil {
			user.Profile = profile
		}
	}

	return &user, nil
}

func (r *adminRepository) UpdateUser(userID string, req *UpdateUserRequest) error {
	updates := make(map[string]interface{})

	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}

	if len(updates) == 0 {
		return errors.New("no fields to update")
	}

	updates["updated_at"] = time.Now()

	return r.db.Model(&User{}).Where("id = ?", userID).Updates(updates).Error
}

func (r *adminRepository) DeleteUser(userID string) error {
	// Start a transaction
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Delete user's bookmarks
	if err := tx.Where("user_id = ?", userID).Delete(&Bookmark{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Delete user's applications
	if err := tx.Where("student_id = ?", userID).Delete(&Application{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Delete user's job posts (if employer)
	if err := tx.Where("employer_id = ?", userID).Delete(&JobPost{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Delete user's profile
	if err := tx.Where("user_id = ?", userID).Delete(&StudentProfile{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where("user_id = ?", userID).Delete(&EmployerProfile{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Finally delete the user
	if err := tx.Where("id = ?", userID).Delete(&User{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *adminRepository) GetCompanies(req *CompanyListRequest) (*CompanyListResponse, error) {
	var companies []CompanyListItem
	var total int64

	query := r.db.Model(&EmployerProfile{})

	// Apply filters
	if req.Industry != "" {
		query = query.Where("industry = ?", req.Industry)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.Search != "" {
		searchTerm := "%" + req.Search + "%"
		query = query.Where("company_name ILIKE ? OR industry ILIKE ?", searchTerm, searchTerm)
	}

	// Get total count
	query.Count(&total)

	// Apply sorting
	if req.SortBy != "" {
		order := req.SortBy
		if req.SortOrder == "desc" {
			order += " DESC"
		}
		query = query.Order(order)
	} else {
		query = query.Order("created_at DESC")
	}

	// Apply pagination
	offset := (req.Page - 1) * req.Limit
	query = query.Offset(offset).Limit(req.Limit)

	// Execute query with job and application counts
	err := query.Select(`
		employer_profiles.*,
		COUNT(DISTINCT job_posts.id) as jobs_count,
		COUNT(DISTINCT applications.id) as applications_count
	`).
		Joins("LEFT JOIN job_posts ON job_posts.employer_id = employer_profiles.user_id").
		Joins("LEFT JOIN applications ON applications.job_id = job_posts.id").
		Group("employer_profiles.id").
		Find(&companies).Error

	if err != nil {
		return nil, err
	}

	totalPages := int((total + int64(req.Limit) - 1) / int64(req.Limit))

	return &CompanyListResponse{
		Companies: companies,
		Pagination: PaginationInfo{
			Page:       req.Page,
			Limit:      req.Limit,
			Total:      int(total),
			TotalPages: totalPages,
		},
	}, nil
}

func (r *adminRepository) GetCompanyByID(companyID string) (*CompanyDetailResponse, error) {
	var company CompanyDetailResponse

	err := r.db.Model(&EmployerProfile{}).
		Select(`
			employer_profiles.*,
			COUNT(DISTINCT job_posts.id) as jobs_count,
			COUNT(DISTINCT applications.id) as applications_count
		`).
		Joins("LEFT JOIN job_posts ON job_posts.employer_id = employer_profiles.user_id").
		Joins("LEFT JOIN applications ON applications.job_id = job_posts.id").
		Where("employer_profiles.id = ?", companyID).
		Group("employer_profiles.id").
		First(&company).Error

	if err != nil {
		return nil, err
	}

	// Get associated user info
	var user User
	err = r.db.Where("id = ?", company.UserID).First(&user).Error
	if err == nil {
		company.User = &UserInfo{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		}
	}

	return &company, nil
}

func (r *adminRepository) UpdateCompany(companyID string, req *UpdateCompanyRequest) error {
	updates := make(map[string]interface{})

	if req.Status != "" {
		updates["status"] = req.Status
	}
	if req.Industry != "" {
		updates["industry"] = req.Industry
	}
	if req.CompanySize != "" {
		updates["company_size"] = req.CompanySize
	}
	if req.CompanyDescription != "" {
		updates["company_description"] = req.CompanyDescription
	}
	if req.Location != "" {
		updates["location"] = req.Location
	}
	if req.WebsiteUrl != "" {
		updates["website_url"] = req.WebsiteUrl
	}
	if req.RecruiterName != "" {
		updates["recruiter_name"] = req.RecruiterName
	}
	if req.OfficialEmail != "" {
		updates["official_email"] = req.OfficialEmail
	}
	if req.PhoneNumber != "" {
		updates["phone_number"] = req.PhoneNumber
	}

	if len(updates) == 0 {
		return errors.New("no fields to update")
	}

	updates["updated_at"] = time.Now()

	return r.db.Model(&EmployerProfile{}).Where("id = ?", companyID).Updates(updates).Error
}

func (r *adminRepository) DeleteCompany(companyID string) error {
	// Start a transaction
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Get the company to find the user_id
	var company EmployerProfile
	if err := tx.Where("id = ?", companyID).First(&company).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Delete company's job posts
	if err := tx.Where("employer_id = ?", company.UserID).Delete(&JobPost{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Delete applications for the company's jobs
	if err := tx.Where("job_id IN (SELECT id FROM job_posts WHERE employer_id = ?)", company.UserID).Delete(&Application{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Delete the company profile
	if err := tx.Where("id = ?", companyID).Delete(&EmployerProfile{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Delete the associated user
	if err := tx.Where("id = ?", company.UserID).Delete(&User{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *adminRepository) GetCompanyAnalytics() (*CompanyAnalytics, error) {
	var analytics CompanyAnalytics

	// Get total companies
	var totalCompanies int64
	r.db.Model(&EmployerProfile{}).Count(&totalCompanies)
	analytics.TotalCompanies = int(totalCompanies)

	// Get active companies (companies with at least one job posted in last 30 days)
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	var activeCompanies int64
	r.db.Model(&EmployerProfile{}).
		Joins("JOIN job_posts ON job_posts.employer_id = employer_profiles.user_id").
		Where("job_posts.created_at >= ?", thirtyDaysAgo).
		Distinct("employer_profiles.id").
		Count(&activeCompanies)
	analytics.ActiveCompanies = int(activeCompanies)

	// Get verified companies (companies with complete profile)
	var verifiedCompanies int64
	r.db.Model(&EmployerProfile{}).
		Where("company_name IS NOT NULL AND company_name != '' AND industry IS NOT NULL AND industry != ''").
		Count(&verifiedCompanies)
	analytics.VerifiedCompanies = int(verifiedCompanies)

	// Get new companies this month and last month
	now := time.Now()
	thisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	lastMonth := thisMonth.AddDate(0, -1, 0)

	var newCompaniesThisMonth, newCompaniesLastMonth int64
	r.db.Model(&EmployerProfile{}).Where("created_at >= ?", thisMonth).Count(&newCompaniesThisMonth)
	r.db.Model(&EmployerProfile{}).Where("created_at >= ? AND created_at < ?", lastMonth, thisMonth).Count(&newCompaniesLastMonth)
	analytics.NewCompaniesThisMonth = int(newCompaniesThisMonth)
	analytics.NewCompaniesLastMonth = int(newCompaniesLastMonth)

	// Calculate growth rate
	if analytics.NewCompaniesLastMonth > 0 {
		analytics.CompanyGrowthRate = float64(analytics.NewCompaniesThisMonth-analytics.NewCompaniesLastMonth) / float64(analytics.NewCompaniesLastMonth) * 100
	}

	// Get companies by industry
	var industryStats []IndustryStats
	r.db.Model(&EmployerProfile{}).
		Select("industry, COUNT(*) as count").
		Where("industry IS NOT NULL AND industry != ''").
		Group("industry").
		Order("count DESC").
		Limit(10).
		Scan(&industryStats)

	// Calculate percentages
	for i := range industryStats {
		if analytics.TotalCompanies > 0 {
			industryStats[i].Percentage = float64(industryStats[i].Count) / float64(analytics.TotalCompanies) * 100
		}
	}
	analytics.CompaniesByIndustry = industryStats

	// Get companies by location
	var locationStats []LocationStats
	r.db.Model(&EmployerProfile{}).
		Select("location, COUNT(*) as count").
		Where("location IS NOT NULL AND location != ''").
		Group("location").
		Order("count DESC").
		Limit(10).
		Scan(&locationStats)

	// Calculate percentages
	for i := range locationStats {
		if analytics.TotalCompanies > 0 {
			locationStats[i].Percentage = float64(locationStats[i].Count) / float64(analytics.TotalCompanies) * 100
		}
	}
	analytics.CompaniesByLocation = locationStats

	// Get companies by size
	var sizeStats []CompanySizeStats
	r.db.Model(&EmployerProfile{}).
		Select("company_size as size, COUNT(*) as count").
		Where("company_size IS NOT NULL AND company_size != ''").
		Group("company_size").
		Order("count DESC").
		Scan(&sizeStats)

	// Calculate percentages
	for i := range sizeStats {
		if analytics.TotalCompanies > 0 {
			sizeStats[i].Percentage = float64(sizeStats[i].Count) / float64(analytics.TotalCompanies) * 100
		}
	}
	analytics.CompaniesBySize = sizeStats

	return &analytics, nil
}

func (r *adminRepository) CreateAdmin(user *auth.User) error {
	return r.db.Create(user).Error
}

func (r *adminRepository) GetStudents(req *StudentListRequest) (*StudentListResponse, error) {
	var students []StudentListItem
	var total int64

	query := r.db.Table("student_profiles").
		Select(`
			student_profiles.id,
			student_profiles.user_id,
			student_profiles.name,
			student_profiles.email,
			student_profiles.phone_number,
			student_profiles.location,
			student_profiles.education,
			student_profiles.skills,
			student_profiles.portfolio,
			student_profiles.linkedin,
			student_profiles.created_at,
			student_profiles.updated_at
		`)

	// Apply filters
	if req.Search != "" {
		searchTerm := "%" + req.Search + "%"
		query = query.Where("student_profiles.name ILIKE ? OR student_profiles.email ILIKE ?", searchTerm, searchTerm)
	}
	if req.Location != "" {
		query = query.Where("student_profiles.location ILIKE ?", "%"+req.Location+"%")
	}
	if req.Education != "" {
		query = query.Where("student_profiles.education ILIKE ?", "%"+req.Education+"%")
	}

	// Get total count
	query.Count(&total)

	// Apply sorting
	if req.SortBy != "" {
		order := req.SortBy
		if req.SortOrder == "desc" {
			order += " DESC"
		}
		query = query.Order(order)
	} else {
		query = query.Order("student_profiles.created_at DESC")
	}

	// Apply pagination
	offset := (req.Page - 1) * req.Limit
	query = query.Offset(offset).Limit(req.Limit)

	// Execute query
	err := query.Find(&students).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((total + int64(req.Limit) - 1) / int64(req.Limit))

	return &StudentListResponse{
		Students: students,
		Pagination: PaginationInfo{
			Page:       req.Page,
			Limit:      req.Limit,
			Total:      int(total),
			TotalPages: totalPages,
		},
	}, nil
}

func (r *adminRepository) GetEmployers(req *EmployerListRequest) (*EmployerListResponse, error) {
	var employers []EmployerListItem
	var total int64

	query := r.db.Table("employer_profiles").
		Select(`
			employer_profiles.id,
			employer_profiles.user_id,
			employer_profiles.company_name,
			employer_profiles.industry,
			employer_profiles.company_size,
			employer_profiles.city,
			employer_profiles.state,
			employer_profiles.phone_number,
			employer_profiles.official_email,
			employer_profiles.recruiter_name,
			employer_profiles.official_email AS recruiter_email,
			employer_profiles.company_description,
			employer_profiles.website_url,
			employer_profiles.created_at,
			employer_profiles.updated_at,
			COUNT(DISTINCT CASE WHEN job_posts.status = 'published' THEN job_posts.id END) as active_jobs_count
		`).
		Joins("LEFT JOIN job_posts ON job_posts.employer_id = employer_profiles.user_id").
		Group("employer_profiles.id")

	// Apply filters
	if req.Search != "" {
		searchTerm := "%" + req.Search + "%"
		query = query.Where("employer_profiles.company_name ILIKE ? OR employer_profiles.recruiter_name ILIKE ?", searchTerm, searchTerm)
	}
	if req.Industry != "" {
		query = query.Where("employer_profiles.industry = ?", req.Industry)
	}
	if req.City != "" {
		query = query.Where("employer_profiles.city ILIKE ?", "%"+req.City+"%")
	}
	if req.CompanySize != "" {
		query = query.Where("employer_profiles.company_size = ?", req.CompanySize)
	}

	// Get total count (need to count before group by)
	countQuery := r.db.Table("employer_profiles")
	if req.Search != "" {
		searchTerm := "%" + req.Search + "%"
		countQuery = countQuery.Where("company_name ILIKE ? OR recruiter_name ILIKE ?", searchTerm, searchTerm)
	}
	if req.Industry != "" {
		countQuery = countQuery.Where("industry = ?", req.Industry)
	}
	if req.City != "" {
		countQuery = countQuery.Where("city ILIKE ?", "%"+req.City+"%")
	}
	if req.CompanySize != "" {
		countQuery = countQuery.Where("company_size = ?", req.CompanySize)
	}
	countQuery.Count(&total)

	// Apply sorting
	if req.SortBy != "" {
		order := req.SortBy
		if req.SortOrder == "desc" {
			order += " DESC"
		}
		query = query.Order(order)
	} else {
		query = query.Order("employer_profiles.created_at DESC")
	}

	// Apply pagination
	offset := (req.Page - 1) * req.Limit
	query = query.Offset(offset).Limit(req.Limit)

	// Execute query
	err := query.Find(&employers).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((total + int64(req.Limit) - 1) / int64(req.Limit))

	return &EmployerListResponse{
		Employers: employers,
		Pagination: PaginationInfo{
			Page:       req.Page,
			Limit:      req.Limit,
			Total:      int(total),
			TotalPages: totalPages,
		},
	}, nil
}

func (r *adminRepository) GetAllJobs(req *JobListRequest) (*JobListResponse, error) {
	var jobs []JobListItem
	var total int64

	query := r.db.Table("job_posts").
		Select(`
			job_posts.id,
			job_posts.title,
			job_posts.status,
			job_posts.employer_id,
			COALESCE(employer_profiles.company_name, users.name) as employer_name,
			job_posts.location,
			job_posts.job_type,
			job_posts.applications_count,
			job_posts.created_at,
			job_posts.updated_at
		`).
		Joins("LEFT JOIN employer_profiles ON employer_profiles.user_id = job_posts.employer_id").
		Joins("LEFT JOIN users ON users.id = job_posts.employer_id")

	// Apply filters
	if req.Status != "" {
		query = query.Where("job_posts.status = ?", req.Status)
	}
	if req.EmployerID != "" {
		query = query.Where("job_posts.employer_id = ?", req.EmployerID)
	}

	// Get total count
	query.Count(&total)

	// Apply sorting
	if req.SortBy != "" {
		order := "job_posts." + req.SortBy
		if req.SortOrder == "desc" {
			order += " DESC"
		}
		query = query.Order(order)
	} else {
		query = query.Order("job_posts.created_at DESC")
	}

	// Apply pagination
	offset := (req.Page - 1) * req.Limit
	query = query.Offset(offset).Limit(req.Limit)

	// Execute query
	err := query.Find(&jobs).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((total + int64(req.Limit) - 1) / int64(req.Limit))

	return &JobListResponse{
		Jobs: jobs,
		Pagination: PaginationInfo{
			Page:       req.Page,
			Limit:      req.Limit,
			Total:      int(total),
			TotalPages: totalPages,
		},
	}, nil
}

func (r *adminRepository) GetJobByID(jobID string) (*JobDetailResponse, error) {
	var job JobDetailResponse

	err := r.db.Table("job_posts").
		Select(`
			job_posts.id,
			job_posts.title,
			job_posts.description,
			job_posts.status,
			job_posts.employer_id,
			COALESCE(employer_profiles.company_name, users.name) as employer_name,
			users.email as employer_email,
			job_posts.location,
			job_posts.job_type,
			job_posts.salary,
			job_posts.requirements,
			job_posts.responsibilities,
			job_posts.benefits,
			job_posts.applications_count,
			job_posts.hired_candidate_name,
			job_posts.completed_at,
			job_posts.created_at,
			job_posts.updated_at
		`).
		Joins("LEFT JOIN employer_profiles ON employer_profiles.user_id = job_posts.employer_id").
		Joins("LEFT JOIN users ON users.id = job_posts.employer_id").
		Where("job_posts.id = ?", jobID).
		First(&job).Error

	if err != nil {
		return nil, err
	}

	return &job, nil
}

func (r *adminRepository) GetJobStatistics() (*JobStatistics, error) {
	var stats JobStatistics

	// Get total jobs
	var totalJobs int64
	r.db.Model(&JobPost{}).Count(&totalJobs)
	stats.TotalJobs = int(totalJobs)

	// Get jobs by status
	var draftJobs, publishedJobs, completedJobs int64
	r.db.Model(&JobPost{}).Where("status = ?", "draft").Count(&draftJobs)
	r.db.Model(&JobPost{}).Where("status = ?", "published").Count(&publishedJobs)
	r.db.Model(&JobPost{}).Where("status = ?", "completed").Count(&completedJobs)
	stats.DraftJobs = int(draftJobs)
	stats.PublishedJobs = int(publishedJobs)
	stats.CompletedJobs = int(completedJobs)

	// Get total applications
	var totalApplications int64
	r.db.Model(&Application{}).Count(&totalApplications)
	stats.TotalApplications = int(totalApplications)

	// Get total hires from job_hires table
	var totalHires int64
	r.db.Table("job_hires").Count(&totalHires)
	stats.TotalHires = int(totalHires)

	return &stats, nil
}
