package admin

import (
	"asa/internal/employerprofile"
	"errors"
	"time"

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

	// Company Management
	GetCompanies(req *CompanyListRequest) (*CompanyListResponse, error)
	GetCompanyByID(companyID string) (*CompanyDetailResponse, error)
	UpdateCompany(companyID string, req *UpdateCompanyRequest) error
	DeleteCompany(companyID string) error
	GetCompanyAnalytics() (*CompanyAnalytics, error)
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
	var pendingApps, acceptedApps, rejectedApps, withdrawnApps int64
	r.db.Model(&Application{}).Where("status = ?", "applied").Count(&pendingApps)
	r.db.Model(&Application{}).Where("status = ?", "accepted").Count(&acceptedApps)
	r.db.Model(&Application{}).Where("status = ?", "rejected").Count(&rejectedApps)
	r.db.Model(&Application{}).Where("status = ?", "withdrawn").Count(&withdrawnApps)
	analytics.PendingApplications = int(pendingApps)
	analytics.AcceptedApplications = int(acceptedApps)
	analytics.RejectedApplications = int(rejectedApps)
	analytics.WithdrawnApplications = int(withdrawnApps)

	// Get applications this month and last month
	now := time.Now()
	thisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	lastMonth := thisMonth.AddDate(0, -1, 0)

	var appsThisMonth, appsLastMonth int64
	r.db.Model(&Application{}).Where("applied_at >= ?", thisMonth).Count(&appsThisMonth)
	r.db.Model(&Application{}).Where("applied_at >= ? AND applied_at < ?", lastMonth, thisMonth).Count(&appsLastMonth)
	analytics.ApplicationsThisMonth = int(appsThisMonth)
	analytics.ApplicationsLastMonth = int(appsLastMonth)

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

	dashboard := &DashboardAnalytics{
		Jobs:         *jobAnalytics,
		Users:        *userAnalytics,
		Applications: *applicationAnalytics,
		LastUpdated:  time.Now(),
	}

	return dashboard, nil
}

// Model structs for database queries
type JobPost struct {
	ID                string    `gorm:"primaryKey;type:uuid"`
	Title             string    `json:"title"`
	Status            string    `json:"status"`
	Location          string    `json:"location"`
	JobType           string    `json:"jobType"`
	ApplicationsCount int       `json:"applicationsCount"`
	CreatedAt         time.Time `json:"createdAt"`
}

type Application struct {
	ID        string    `gorm:"primaryKey;type:uuid"`
	Status    string    `json:"status"`
	AppliedAt time.Time `json:"appliedAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type User struct {
	ID        string    `gorm:"primaryKey;type:uuid"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type EmployerProfile struct {
	ID          string `gorm:"primaryKey;type:uuid"`
	CompanyName string `json:"companyName"`
	Industry    string `json:"industry"`
	Location    string `json:"location"`
}

type Bookmark struct {
	ID     string `gorm:"primaryKey;type:uuid"`
	UserID string `json:"userId"`
}

type StudentProfile struct {
	ID       string `gorm:"primaryKey;type:uuid"`
	Location string `json:"location"`
	Skills   string `json:"skills"`
}

// TableName specifies the database table name for StudentProfile
func (StudentProfile) TableName() string {
	return "student_profiles"
}

// User Management Methods
func (r *adminRepository) GetUsers(req *UserListRequest) (*UserListResponse, error) {
	var users []UserListItem
	var total int64

	// Build query
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
		} else {
			order += " ASC"
		}
		query = query.Order(order)
	} else {
		query = query.Order("created_at DESC")
	}

	// Apply pagination
	offset := (req.Page - 1) * req.Limit
	query = query.Offset(offset).Limit(req.Limit)

	// Execute query
	err := query.Select("id, name, email, role, status, created_at, updated_at").Scan(&users).Error
	if err != nil {
		return nil, err
	}

	// Calculate pagination info
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

	// Get basic user info
	err := r.db.Model(&User{}).
		Select("id, name, email, role, status, created_at, updated_at").
		Where("id = ?", userID).
		First(&user).Error
	if err != nil {
		return nil, err
	}

	// Get role-specific profile based on profile existence
	// Try employer profile first, then student profile
	var employerProfile struct {
		CompanyName string `json:"companyName"`
		Industry    string `json:"industry"`
		Location    string `json:"location"`
	}
	err = r.db.Model(&EmployerProfile{}).
		Select("company_name, industry, location").
		Where("user_id = ?", userID).
		First(&employerProfile).Error

	if err == nil {
		user.Profile = employerProfile
	} else {
		// Try student profile if employer profile doesn't exist
		var studentProfile struct {
			Location string `json:"location"`
			Skills   string `json:"skills"`
		}
		err = r.db.Model(&StudentProfile{}).
			Select("location, skills").
			Where("user_id = ?", userID).
			First(&studentProfile).Error
		if err == nil {
			user.Profile = studentProfile
		}
	}

	return &user, nil
}

func (r *adminRepository) UpdateUser(userID string, req *UpdateUserRequest) error {
	// Get existing user
	var user User
	err := r.db.First(&user, "id = ?", userID).Error
	if err != nil {
		return err
	}

	// Update fields if provided
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		// Check if email is already taken by another user
		var existingUser User
		err := r.db.Where("email = ? AND id != ?", req.Email, userID).First(&existingUser).Error
		if err == nil {
			return errors.New("email already taken by another user")
		}
		user.Email = req.Email
	}
	// Role updates are now handled by AAA service
	// We don't store roles in local DB anymore
	if req.Status != "" {
		user.Status = req.Status
	}

	// Save changes
	return r.db.Save(&user).Error
}

func (r *adminRepository) DeleteUser(userID string) error {
	// Start a transaction
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Delete user's profile based on role
	var user User
	err := tx.First(&user, "id = ?", userID).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// Delete role-specific profile based on profile existence
	// Delete both profile types to ensure cleanup
	tx.Where("user_id = ?", userID).Delete(&EmployerProfile{})
	tx.Where("user_id = ?", userID).Delete(&StudentProfile{})

	// Delete related data
	tx.Where("student_id = ? OR employer_id = ?", userID, userID).Delete(&Application{})
	tx.Where("user_id = ?", userID).Delete(&Bookmark{})
	tx.Where("employer_id = ?", userID).Delete(&JobPost{})

	// Finally delete the user
	err = tx.Delete(&user).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// Company Management Methods
func (r *adminRepository) GetCompanies(req *CompanyListRequest) (*CompanyListResponse, error) {
	var companies []CompanyListItem
	var total int64

	// Build query with joins to get job and application counts
	query := r.db.Table("employer_profiles").
		Select(`
			employer_profiles.id,
			employer_profiles.company_name,
			employer_profiles.industry,
			employer_profiles.city || ', ' || employer_profiles.state as location,
			employer_profiles.company_size,
			'active' as status,
			employer_profiles.created_at,
			employer_profiles.updated_at,
			COALESCE(job_counts.jobs_count, 0) as jobs_count,
			COALESCE(app_counts.applications_count, 0) as applications_count
		`).
		Joins("LEFT JOIN (SELECT employer_id, COUNT(*) as jobs_count FROM job_posts GROUP BY employer_id) job_counts ON employer_profiles.id = job_counts.employer_id").
		Joins("LEFT JOIN (SELECT j.employer_id, COUNT(a.id) as applications_count FROM job_posts j LEFT JOIN applications a ON j.id = a.job_id GROUP BY j.employer_id) app_counts ON employer_profiles.id = app_counts.employer_id")

	// Apply filters
	if req.Industry != "" {
		query = query.Where("employer_profiles.industry = ?", req.Industry)
	}
	if req.Search != "" {
		searchTerm := "%" + req.Search + "%"
		query = query.Where("employer_profiles.company_name ILIKE ? OR employer_profiles.industry ILIKE ?", searchTerm, searchTerm)
	}

	// Get total count
	query.Count(&total)

	// Apply sorting
	if req.SortBy != "" {
		order := req.SortBy
		if req.SortOrder == "desc" {
			order += " DESC"
		} else {
			order += " ASC"
		}
		query = query.Order(order)
	} else {
		query = query.Order("employer_profiles.created_at DESC")
	}

	// Apply pagination
	offset := (req.Page - 1) * req.Limit
	query = query.Offset(offset).Limit(req.Limit)

	// Execute query
	err := query.Scan(&companies).Error
	if err != nil {
		return nil, err
	}

	// Calculate pagination info
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

	// Get company details with job and application counts
	err := r.db.Table("employer_profiles").
		Select(`
			employer_profiles.id,
			employer_profiles.company_name,
			employer_profiles.industry,
			employer_profiles.company_size,
			employer_profiles.company_description,
			employer_profiles.city || ', ' || employer_profiles.state as location,
			employer_profiles.website_url,
			employer_profiles.recruiter_name,
			employer_profiles.official_email,
			employer_profiles.phone_number,
			'active' as status,
			employer_profiles.created_at,
			employer_profiles.updated_at,
			COALESCE(job_counts.jobs_count, 0) as jobs_count,
			COALESCE(app_counts.applications_count, 0) as applications_count
		`).
		Joins("LEFT JOIN (SELECT employer_id, COUNT(*) as jobs_count FROM job_posts GROUP BY employer_id) job_counts ON employer_profiles.id = job_counts.employer_id").
		Joins("LEFT JOIN (SELECT j.employer_id, COUNT(a.id) as applications_count FROM job_posts j LEFT JOIN applications a ON j.id = a.job_id GROUP BY j.employer_id) app_counts ON employer_profiles.id = app_counts.employer_id").
		Where("employer_profiles.id = ?", companyID).
		First(&company).Error

	if err != nil {
		return nil, err
	}

	// Get associated user info
	var user UserInfo
	err = r.db.Table("users").
		Select("id, name, email, role").
		Joins("JOIN employer_profiles ON users.id = employer_profiles.user_id").
		Where("employer_profiles.id = ?", companyID).
		First(&user).Error

	if err == nil {
		company.User = &user
	}

	return &company, nil
}

func (r *adminRepository) UpdateCompany(companyID string, req *UpdateCompanyRequest) error {
	// Get existing company
	var company employerprofile.EmployerProfile
	err := r.db.First(&company, "id = ?", companyID).Error
	if err != nil {
		return err
	}

	// Update fields if provided
	if req.CompanyName != "" {
		company.CompanyName = req.CompanyName
	}
	if req.Industry != "" {
		company.Industry = req.Industry
	}
	if req.CompanySize != "" {
		company.CompanySize = req.CompanySize
	}
	if req.CompanyDescription != "" {
		company.CompanyDescription = req.CompanyDescription
	}
	if req.Location != "" {
		// Parse location into city and state
		// This is a simplified approach - you might want to enhance this
		company.City = req.Location
	}
	if req.WebsiteUrl != "" {
		company.WebsiteUrl = req.WebsiteUrl
	}
	if req.RecruiterName != "" {
		company.RecruiterName = req.RecruiterName
	}
	if req.OfficialEmail != "" {
		company.OfficialEmail = req.OfficialEmail
	}
	if req.PhoneNumber != "" {
		company.PhoneNumber = req.PhoneNumber
	}

	// Save changes
	return r.db.Save(&company).Error
}

func (r *adminRepository) DeleteCompany(companyID string) error {
	// Start a transaction
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get company details
	var company employerprofile.EmployerProfile
	err := tx.First(&company, "id = ?", companyID).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// Delete related data
	tx.Where("employer_id = ?", companyID).Delete(&JobPost{})

	// Delete applications for jobs from this company
	tx.Exec("DELETE FROM applications WHERE job_id IN (SELECT id FROM job_posts WHERE employer_id = ?)", companyID)

	// Finally delete the company profile
	err = tx.Delete(&company).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *adminRepository) GetCompanyAnalytics() (*CompanyAnalytics, error) {
	var analytics CompanyAnalytics

	// Get total companies
	var totalCompanies int64
	r.db.Model(&employerprofile.EmployerProfile{}).Count(&totalCompanies)
	analytics.TotalCompanies = int(totalCompanies)

	// Get active companies (companies with at least one job)
	var activeCompanies int64
	r.db.Model(&employerprofile.EmployerProfile{}).
		Joins("JOIN job_posts ON employer_profiles.id = job_posts.employer_id").
		Distinct("employer_profiles.id").
		Count(&activeCompanies)
	analytics.ActiveCompanies = int(activeCompanies)

	// Get verified companies (companies with complete profile)
	var verifiedCompanies int64
	r.db.Model(&employerprofile.EmployerProfile{}).
		Where("company_name IS NOT NULL AND company_name != '' AND industry IS NOT NULL AND industry != ''").
		Count(&verifiedCompanies)
	analytics.VerifiedCompanies = int(verifiedCompanies)

	// Get new companies this month and last month
	now := time.Now()
	thisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	lastMonth := thisMonth.AddDate(0, -1, 0)

	var newCompaniesThisMonth, newCompaniesLastMonth int64
	r.db.Model(&employerprofile.EmployerProfile{}).Where("created_at >= ?", thisMonth).Count(&newCompaniesThisMonth)
	r.db.Model(&employerprofile.EmployerProfile{}).Where("created_at >= ? AND created_at < ?", lastMonth, thisMonth).Count(&newCompaniesLastMonth)
	analytics.NewCompaniesThisMonth = int(newCompaniesThisMonth)
	analytics.NewCompaniesLastMonth = int(newCompaniesLastMonth)

	// Calculate growth rate
	if analytics.NewCompaniesLastMonth > 0 {
		analytics.CompanyGrowthRate = float64(analytics.NewCompaniesThisMonth-analytics.NewCompaniesLastMonth) / float64(analytics.NewCompaniesLastMonth) * 100
	}

	// Get companies by industry
	var industryStats []IndustryStats
	r.db.Model(&employerprofile.EmployerProfile{}).
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
	r.db.Model(&employerprofile.EmployerProfile{}).
		Select("city || ', ' || state as location, COUNT(*) as count").
		Where("city IS NOT NULL AND city != '' AND state IS NOT NULL AND state != ''").
		Group("city, state").
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
	r.db.Model(&employerprofile.EmployerProfile{}).
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
