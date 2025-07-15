package jobpost

import (
	"fmt"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type JobPostRepository interface {
	Create(job *JobPost) error
	Update(job *JobPost) error
	Delete(id string) error
	GetByID(id string) (*JobPost, error)
	GetByEmployer(employerID string) ([]JobPost, error)
	Search(filter *JobPostFilter) ([]JobPost, error)
	GetJobsByIDs(ids []string) ([]JobPost, error)
	IncrementApplicationsCount(jobID string) error
	GetFeaturedJobs(limit int) ([]JobPost, error)
	GetRecentJobs(limit int) ([]JobPost, error)

	// Enhanced search and discovery methods
	AdvancedSearch(request *AdvancedJobSearchRequest) ([]JobPost, int, error)
	GetSearchFilters() (*SearchFilters, error)
	GetTrendingJobs(limit int) ([]JobPost, error)
	GetSimilarJobs(jobID string, maxResults int) ([]JobPost, error)
	GetRecommendedJobs(request *JobRecommendationRequest) ([]JobPost, error)

	// Job alerts methods
	CreateJobAlert(alert *JobAlert) error
	UpdateJobAlert(alert *JobAlert) error
	DeleteJobAlert(alertID string) error
	GetJobAlertByID(alertID string) (*JobAlert, error)
	GetJobAlertsByUser(userID string) ([]JobAlert, error)
	GetActiveJobAlerts() ([]JobAlert, error)
	GetJobsMatchingAlert(alert *JobAlert) ([]JobPost, error)

	// Draft methods
	CreateDraft(job *JobPost) error
	GetDraftsByEmployer(employerID string) ([]JobPost, error)
	PublishDraft(jobID string) error
}

type jobPostRepository struct {
	db *gorm.DB
}

func NewJobPostRepository(db *gorm.DB) JobPostRepository {
	return &jobPostRepository{db}
}

// BeforeCreate hook to map salary fields
func (job *JobPost) BeforeCreate(tx *gorm.DB) error {
	// Map nested salary to individual fields for database storage
	job.SalaryMin = job.Salary.Min
	job.SalaryMax = job.Salary.Max
	job.SalaryCurrency = job.Salary.Currency

	// Debug logging
	fmt.Printf("DEBUG: BeforeCreate hook - Salary: %+v, Min: %f, Max: %f, Currency: %s\n",
		job.Salary, job.SalaryMin, job.SalaryMax, job.SalaryCurrency)

	return nil
}

// BeforeUpdate hook to map salary fields
func (job *JobPost) BeforeUpdate(tx *gorm.DB) error {
	// Map nested salary to individual fields for database storage
	job.SalaryMin = job.Salary.Min
	job.SalaryMax = job.Salary.Max
	job.SalaryCurrency = job.Salary.Currency
	return nil
}

// AfterFind hook to map database fields to nested salary structure
func (job *JobPost) AfterFind(tx *gorm.DB) error {
	// Map individual fields to nested salary structure for JSON response
	job.Salary = Salary{
		Min:      job.SalaryMin,
		Max:      job.SalaryMax,
		Currency: job.SalaryCurrency,
	}
	return nil
}

func (r *jobPostRepository) Create(job *JobPost) error {
	return r.db.Create(job).Error
}

func (r *jobPostRepository) Update(job *JobPost) error {
	return r.db.Save(job).Error
}

func (r *jobPostRepository) Delete(id string) error {
	return r.db.Delete(&JobPost{}, "id = ?", id).Error
}

func (r *jobPostRepository) GetByID(id string) (*JobPost, error) {
	var job JobPost
	err := r.db.First(&job, "id = ?", id).Error
	return &job, err
}

func (r *jobPostRepository) GetByEmployer(employerID string) ([]JobPost, error) {
	var jobs []JobPost
	err := r.db.Where("employer_id = ?", employerID).Find(&jobs).Error
	return jobs, err
}

func (r *jobPostRepository) Search(filter *JobPostFilter) ([]JobPost, error) {
	var jobs []JobPost
	query := r.db.Model(&JobPost{}).Where("status = ?", "published") // Only published jobs

	fmt.Printf("🔍 Job Search Debug - Filters applied:\n")
	fmt.Printf("  - Location: %s\n", filter.Location)
	fmt.Printf("  - JobType: %v\n", filter.JobType)
	fmt.Printf("  - Experience: %v\n", filter.Experience)
	fmt.Printf("  - Skills: %v\n", filter.Skills)
	fmt.Printf("  - IsRemote: %v\n", filter.IsRemote)
	fmt.Printf("  - PostedWithin: %s\n", filter.PostedWithin)
	if filter.SalaryRange != nil {
		fmt.Printf("  - SalaryRange: %v - %v\n", filter.SalaryRange.Min, filter.SalaryRange.Max)
	}

	if filter.Location != "" {
		query = query.Where("location ILIKE ?", "%"+filter.Location+"%")
	}

	if len(filter.JobType) > 0 {
		query = query.Where("job_type IN ?", filter.JobType)
	}

	if len(filter.Experience) > 0 {
		query = query.Where("experience IN ?", filter.Experience)
	}

	if filter.SalaryRange != nil {
		query = query.Where("salary_min <= ? AND salary_max >= ?", filter.SalaryRange.Max, filter.SalaryRange.Min)
	}

	if filter.IsRemote != nil {
		query = query.Where("is_remote = ?", *filter.IsRemote)
	}

	if len(filter.Skills) > 0 {
		// Search for jobs that have any of the required skills
		// Use a more robust approach for PostgreSQL array search
		for _, skill := range filter.Skills {
			query = query.Where("required_skills @> ? OR required_skills ILIKE ?",
				pq.Array([]string{skill}), "%"+skill+"%")
		}
	}

	if filter.PostedWithin != "" && filter.PostedWithin != "all" {
		var cutoffDate time.Time
		now := time.Now()

		switch filter.PostedWithin {
		case "24h":
			cutoffDate = now.AddDate(0, 0, -1)
		case "7d":
			cutoffDate = now.AddDate(0, 0, -7)
		case "30d":
			cutoffDate = now.AddDate(0, 0, -30)
		default:
			// No filtering
		}

		if !cutoffDate.IsZero() {
			query = query.Where("created_at >= ?", cutoffDate)
		}
	}

	// Add pagination
	if filter.Page > 0 && filter.Limit > 0 {
		offset := (filter.Page - 1) * filter.Limit
		query = query.Offset(offset).Limit(filter.Limit)
	}

	// Order by most recent first
	query = query.Order("created_at DESC")

	err := query.Find(&jobs).Error
	fmt.Printf("🔍 Job Search Debug - Found %d jobs\n", len(jobs))
	return jobs, err
}

func (r *jobPostRepository) GetJobsByIDs(ids []string) ([]JobPost, error) {
	var jobs []JobPost
	err := r.db.Where("id IN ?", ids).Find(&jobs).Error
	return jobs, err
}

func (r *jobPostRepository) IncrementApplicationsCount(jobID string) error {
	return r.db.Model(&JobPost{}).
		Where("id = ?", jobID).
		UpdateColumn("applications_count", gorm.Expr("applications_count + 1")).
		Error
}

func (r *jobPostRepository) GetFeaturedJobs(limit int) ([]JobPost, error) {
	var jobs []JobPost

	// Featured jobs criteria:
	// 1. Published status
	// 2. High applications count (popular jobs)
	// 3. Recent posts (within last 30 days)
	// 4. Order by applications count, then by recency
	query := r.db.Model(&JobPost{}).
		Where("status = ?", "published").
		Where("created_at >= ?", time.Now().AddDate(0, 0, -30)).
		Order("applications_count DESC, created_at DESC").
		Limit(limit)

	err := query.Find(&jobs).Error
	return jobs, err
}

func (r *jobPostRepository) GetRecentJobs(limit int) ([]JobPost, error) {
	var jobs []JobPost

	// Recent jobs criteria:
	// 1. Published status
	// 2. Order by most recent first
	query := r.db.Model(&JobPost{}).
		Where("status = ?", "published").
		Order("created_at DESC").
		Limit(limit)

	err := query.Find(&jobs).Error
	return jobs, err
}

// Enhanced search and discovery methods
func (r *jobPostRepository) AdvancedSearch(request *AdvancedJobSearchRequest) ([]JobPost, int, error) {
	var jobs []JobPost
	var total int64

	query := r.db.Model(&JobPost{}).Where("status = ?", "published")

	// Keywords search (title, role_overview, requirements)
	if request.Keywords != "" {
		keywords := "%" + request.Keywords + "%"
		query = query.Where("title ILIKE ? OR role_overview ILIKE ? OR requirements ILIKE ?",
			keywords, keywords, keywords)
	}

	// Location filter
	if request.Location != "" {
		if request.Location == "remote" {
			query = query.Where("is_remote = ?", true)
		} else {
			query = query.Where("location ILIKE ?", "%"+request.Location+"%")
		}
	}

	// Job type filter
	if len(request.JobType) > 0 {
		query = query.Where("job_type IN ?", request.JobType)
	}

	// Experience filter
	if len(request.Experience) > 0 {
		query = query.Where("experience IN ?", request.Experience)
	}

	// Skills filter
	if len(request.Skills) > 0 {
		for _, skill := range request.Skills {
			query = query.Where("required_skills @> ? OR required_skills ILIKE ?",
				pq.Array([]string{skill}), "%"+skill+"%")
		}
	}

	// Salary range filter
	if request.SalaryRange != nil {
		query = query.Where("salary_min <= ? AND salary_max >= ?",
			request.SalaryRange.Max, request.SalaryRange.Min)
	}

	// Benefits filter
	if len(request.Benefits) > 0 {
		for _, benefit := range request.Benefits {
			query = query.Where("benefits @> ?", pq.Array([]string{benefit}))
		}
	}

	// Work preference filters
	if request.IsRemote != nil {
		query = query.Where("is_remote = ?", *request.IsRemote)
	}

	// Posted within filter
	if request.PostedWithin != "" && request.PostedWithin != "all" {
		var cutoffDate time.Time
		now := time.Now()

		switch request.PostedWithin {
		case "24h":
			cutoffDate = now.AddDate(0, 0, -1)
		case "7d":
			cutoffDate = now.AddDate(0, 0, -7)
		case "30d":
			cutoffDate = now.AddDate(0, 0, -30)
		case "90d":
			cutoffDate = now.AddDate(0, 0, -90)
		}

		if !cutoffDate.IsZero() {
			query = query.Where("created_at >= ?", cutoffDate)
		}
	}

	// Get total count before pagination
	query.Count(&total)

	// Sorting
	sortBy := request.SortBy
	if sortBy == "" {
		sortBy = "relevance"
	}

	sortOrder := request.SortOrder
	if sortOrder == "" {
		sortOrder = "desc"
	}

	switch sortBy {
	case "date":
		query = query.Order("created_at " + sortOrder)
	case "salary":
		query = query.Order("salary_min " + sortOrder)
	case "applications":
		query = query.Order("applications_count " + sortOrder)
	default: // relevance
		query = query.Order("applications_count DESC, created_at DESC")
	}

	// Pagination
	page := request.Page
	if page <= 0 {
		page = 1
	}

	limit := request.Limit
	if limit <= 0 {
		limit = 20
	} else if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit
	query = query.Offset(offset).Limit(limit)

	err := query.Find(&jobs).Error
	return jobs, int(total), err
}

func (r *jobPostRepository) GetSearchFilters() (*SearchFilters, error) {
	filters := &SearchFilters{}

	// Get available locations
	var locations []string
	r.db.Model(&JobPost{}).
		Where("status = ?", "published").
		Distinct().
		Pluck("location", &locations)
	filters.AvailableLocations = locations

	// Get available job types
	var jobTypes []string
	r.db.Model(&JobPost{}).
		Where("status = ?", "published").
		Distinct().
		Pluck("job_type", &jobTypes)
	filters.AvailableJobTypes = jobTypes

	// Get available experience levels
	var experience []string
	r.db.Model(&JobPost{}).
		Where("status = ?", "published").
		Distinct().
		Pluck("experience", &experience)
	filters.AvailableExperience = experience

	// Get available skills (from all jobs)
	var allSkills []string
	r.db.Model(&JobPost{}).
		Where("status = ?", "published").
		Pluck("required_skills", &allSkills)

	// Flatten and deduplicate skills
	skillSet := make(map[string]bool)
	for _, skills := range allSkills {
		// This is a simplified approach - in production you'd want to properly parse the array
		if skills != "" {
			skillSet[skills] = true
		}
	}

	for skill := range skillSet {
		filters.AvailableSkills = append(filters.AvailableSkills, skill)
	}

	// Get salary ranges
	filters.SalaryRanges = []SalaryRange{
		{Min: 0, Max: 30000, Currency: "USD", Label: "$0-$30k"},
		{Min: 30000, Max: 50000, Currency: "USD", Label: "$30k-$50k"},
		{Min: 50000, Max: 75000, Currency: "USD", Label: "$50k-$75k"},
		{Min: 75000, Max: 100000, Currency: "USD", Label: "$75k-$100k"},
		{Min: 100000, Max: 150000, Currency: "USD", Label: "$100k-$150k"},
		{Min: 150000, Max: 0, Currency: "USD", Label: "$150k+"},
	}

	// Get available benefits
	var allBenefits []string
	r.db.Model(&JobPost{}).
		Where("status = ?", "published").
		Pluck("benefits", &allBenefits)

	benefitSet := make(map[string]bool)
	for _, benefits := range allBenefits {
		if benefits != "" {
			benefitSet[benefits] = true
		}
	}

	for benefit := range benefitSet {
		filters.AvailableBenefits = append(filters.AvailableBenefits, benefit)
	}

	return filters, nil
}

func (r *jobPostRepository) GetTrendingJobs(limit int) ([]JobPost, error) {
	var jobs []JobPost

	// Trending jobs criteria:
	// 1. Published status
	// 2. High applications count (popular)
	// 3. Recent posts (within last 7 days)
	// 4. Order by applications count, then by recency
	query := r.db.Model(&JobPost{}).
		Where("status = ?", "published").
		Where("created_at >= ?", time.Now().AddDate(0, 0, -7)).
		Where("applications_count > ?", 0).
		Order("applications_count DESC, created_at DESC").
		Limit(limit)

	err := query.Find(&jobs).Error
	return jobs, err
}

func (r *jobPostRepository) GetSimilarJobs(jobID string, maxResults int) ([]JobPost, error) {
	// First get the reference job
	var referenceJob JobPost
	err := r.db.First(&referenceJob, "id = ?", jobID).Error
	if err != nil {
		return nil, err
	}

	var jobs []JobPost

	// Find similar jobs based on:
	// 1. Same job type
	// 2. Same experience level
	// 3. Similar skills
	// 4. Similar location
	// 5. Similar salary range
	query := r.db.Model(&JobPost{}).
		Where("status = ?", "published").
		Where("id != ?", jobID).
		Where("job_type = ?", referenceJob.JobType).
		Where("experience = ?", referenceJob.Experience)

	// Add location similarity
	if referenceJob.Location != "" {
		query = query.Where("location ILIKE ?", "%"+referenceJob.Location+"%")
	}

	// Add salary range similarity (within 20% range)
	if referenceJob.Salary.Min > 0 && referenceJob.Salary.Max > 0 {
		salaryMin := referenceJob.Salary.Min * 0.8
		salaryMax := referenceJob.Salary.Max * 1.2
		query = query.Where("salary_min >= ? AND salary_max <= ?", salaryMin, salaryMax)
	}

	// Order by skill similarity, then by recency
	query = query.Order("created_at DESC").Limit(maxResults)

	err = query.Find(&jobs).Error
	return jobs, err
}

func (r *jobPostRepository) GetRecommendedJobs(request *JobRecommendationRequest) ([]JobPost, error) {
	var jobs []JobPost

	query := r.db.Model(&JobPost{}).
		Where("status = ?", "published")

	// Filter by user's preferred job types
	if len(request.PreferredJobTypes) > 0 {
		query = query.Where("job_type IN ?", request.PreferredJobTypes)
	}

	// Filter by user's experience level
	if request.UserExperience != "" {
		query = query.Where("experience = ?", request.UserExperience)
	}

	// Filter by user's location
	if request.UserLocation != "" {
		query = query.Where("location ILIKE ?", "%"+request.UserLocation+"%")
	}

	// Filter by user's skills
	if len(request.UserSkills) > 0 {
		for _, skill := range request.UserSkills {
			query = query.Where("required_skills @> ?", pq.Array([]string{skill}))
		}
	}

	// Order by relevance (skill match, then recency)
	query = query.Order("created_at DESC").Limit(request.MaxResults)

	err := query.Find(&jobs).Error
	return jobs, err
}

// Job alerts methods
func (r *jobPostRepository) CreateJobAlert(alert *JobAlert) error {
	return r.db.Create(alert).Error
}

func (r *jobPostRepository) UpdateJobAlert(alert *JobAlert) error {
	return r.db.Save(alert).Error
}

func (r *jobPostRepository) DeleteJobAlert(alertID string) error {
	return r.db.Delete(&JobAlert{}, "id = ?", alertID).Error
}

func (r *jobPostRepository) GetJobAlertByID(alertID string) (*JobAlert, error) {
	var alert JobAlert
	err := r.db.First(&alert, "id = ?", alertID).Error
	return &alert, err
}

func (r *jobPostRepository) GetJobAlertsByUser(userID string) ([]JobAlert, error) {
	var alerts []JobAlert
	err := r.db.Where("user_id = ?", userID).Find(&alerts).Error
	return alerts, err
}

func (r *jobPostRepository) GetActiveJobAlerts() ([]JobAlert, error) {
	var alerts []JobAlert
	err := r.db.Where("is_active = ?", true).Find(&alerts).Error
	return alerts, err
}

func (r *jobPostRepository) GetJobsMatchingAlert(alert *JobAlert) ([]JobPost, error) {
	var jobs []JobPost

	query := r.db.Model(&JobPost{}).
		Where("status = ?", "published")

	// Keywords filter
	if len(alert.Keywords) > 0 {
		for _, keyword := range alert.Keywords {
			query = query.Where("title ILIKE ? OR role_overview ILIKE ? OR requirements ILIKE ?",
				"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
		}
	}

	// Location filter
	if alert.Location != "" {
		query = query.Where("location ILIKE ?", "%"+alert.Location+"%")
	}

	// Job type filter
	if len(alert.JobType) > 0 {
		query = query.Where("job_type IN ?", alert.JobType)
	}

	// Experience filter
	if len(alert.Experience) > 0 {
		query = query.Where("experience IN ?", alert.Experience)
	}

	// Skills filter
	if len(alert.Skills) > 0 {
		for _, skill := range alert.Skills {
			query = query.Where("required_skills @> ?", pq.Array([]string{skill}))
		}
	}

	// Salary range filter
	if alert.SalaryMin != nil && alert.SalaryMax != nil {
		query = query.Where("salary_min >= ? AND salary_max <= ?",
			*alert.SalaryMin, *alert.SalaryMax)
	}

	// Remote filter
	if alert.IsRemote != nil {
		query = query.Where("is_remote = ?", *alert.IsRemote)
	}

	// Only recent jobs (within last 7 days for alerts)
	query = query.Where("created_at >= ?", time.Now().AddDate(0, 0, -7))

	// Order by recency
	query = query.Order("created_at DESC")

	err := query.Find(&jobs).Error
	return jobs, err
}

// Draft methods
func (r *jobPostRepository) CreateDraft(job *JobPost) error {
	// Set status to draft
	job.Status = "draft"
	return r.db.Create(job).Error
}

func (r *jobPostRepository) GetDraftsByEmployer(employerID string) ([]JobPost, error) {
	var jobs []JobPost
	err := r.db.Where("employer_id = ? AND status = ?", employerID, "draft").Order("updated_at DESC").Find(&jobs).Error
	return jobs, err
}

func (r *jobPostRepository) PublishDraft(jobID string) error {
	return r.db.Model(&JobPost{}).Where("id = ?", jobID).Update("status", "published").Error
}
