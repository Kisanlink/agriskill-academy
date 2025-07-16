package jobpost

import (
	"errors"
	"time"

	"asa/internal/employerprofile"
	"asa/internal/middleware"
)

type JobPostService interface {
	CreateJobPost(req *CreateJobPostRequest, employerID, employerName, employerEmail string) (*JobPost, error)
	CreateJobPostWithStatus(req *CreateJobPostRequest, employerID, employerName, employerEmail, status string) (*JobPost, error)
	UpdateJobPost(id string, req *UpdateJobPostRequest) (*JobPost, error)
	Delete(id string) error
	GetByID(id string) (*JobPost, error)
	GetByEmployer(employerID string) ([]JobPost, error)
	Search(filter *JobPostFilter) ([]JobPost, error)
	IncrementApplicationsCount(jobID string) error
	GetFeaturedJobs(limit int) ([]JobPost, error)
	GetRecentJobs(limit int) ([]JobPost, error)

	// Enhanced search and discovery methods
	AdvancedSearch(request *AdvancedJobSearchRequest) (*JobSearchResponse, error)
	GetSearchFilters() (*SearchFilters, error)
	GetTrendingJobs(limit int) ([]JobPost, error)
	GetSimilarJobs(jobID string, maxResults int) ([]JobPost, error)
	GetRecommendedJobs(request *JobRecommendationRequest) (*JobRecommendationResponse, error)

	// Job alerts methods
	CreateJobAlert(request *JobAlertRequest) (*JobAlert, error)
	UpdateJobAlert(alertID string, request *JobAlertRequest) (*JobAlert, error)
	DeleteJobAlert(alertID string) error
	GetJobAlertByID(alertID string) (*JobAlert, error)
	GetJobAlertsByUser(userID string) ([]JobAlert, error)
	ProcessJobAlerts() error

	// Draft methods
	CreateDraft(req *CreateDraftRequest, employerID, employerName, employerEmail string) (*JobPost, error)
	GetDraftsByEmployer(employerID string) ([]JobPost, error)
	PublishDraft(jobID string) (*JobPost, error)
}

type jobPostService struct {
	repo         JobPostRepository
	employerRepo employerprofile.EmployerProfileRepository
}

func NewJobPostService(repo JobPostRepository, employerRepo employerprofile.EmployerProfileRepository) JobPostService {
	return &jobPostService{repo: repo, employerRepo: employerRepo}
}

func (s *jobPostService) CreateJobPost(req *CreateJobPostRequest, employerID, employerName, employerEmail string) (*JobPost, error) {
	// Validate required fields
	if req.Title == "" {
		return nil, errors.New("title is required")
	}
	if req.RoleOverview == "" {
		return nil, errors.New("role overview is required")
	}
	if req.Requirements == "" {
		return nil, errors.New("requirements are required")
	}
	if req.Location == "" {
		return nil, errors.New("location is required")
	}
	if req.JobType == "" {
		return nil, errors.New("job type is required")
	}
	if req.Experience == "" {
		return nil, errors.New("experience level is required")
	}

	// Validate text field lengths (ensure they can handle substantial content)
	if len(req.RoleOverview) < 10 {
		return nil, errors.New("role overview should be at least 10 characters long")
	}
	if len(req.Requirements) < 10 {
		return nil, errors.New("requirements should be at least 10 characters long")
	}

	// Note: No maximum length limits - PostgreSQL text type can handle unlimited text
	// The backend supports unlimited text for role overview and requirements

	// Fetch employer details if not provided
	if employerName == "" || employerEmail == "" {
		employerProfile, err := s.employerRepo.GetByUserID(employerID)
		if err == nil && employerProfile != nil {
			if employerName == "" {
				employerName = employerProfile.RecruiterName // Use recruiter name (actual employer name)
			}
			if employerEmail == "" {
				employerEmail = employerProfile.OfficialEmail
			}
		}
	}

	// Validate salary
	var salary Salary
	if req.Salary.Salary != nil {
		salary = *req.Salary.Salary
		if salary.Min < 0 || salary.Max < 0 {
			return nil, errors.New("salary values cannot be negative")
		}
		if salary.Min > salary.Max {
			return nil, errors.New("minimum salary cannot be greater than maximum salary")
		}
		if salary.Currency == "" {
			salary.Currency = "USD" // Default currency
		}
	} else {
		// If no salary provided, set default values
		salary = Salary{
			Min:      0,
			Max:      0,
			Currency: "USD",
		}
	}

	// Validate application deadline
	if req.ApplicationDeadline.Before(time.Now()) {
		return nil, errors.New("application deadline cannot be in the past")
	}

	job := &JobPost{
		Title:               req.Title,
		RoleOverview:        req.RoleOverview,
		Requirements:        req.Requirements,
		Location:            req.Location,
		RequiredSkills:      req.RequiredSkills,
		EmployerID:          employerID,
		EmployerName:        employerName,
		EmployerEmail:       employerEmail,
		Status:              "draft", // Default status
		ApplicationDeadline: req.ApplicationDeadline,
		JobType:             req.JobType,
		Experience:          req.Experience,
		Salary:              salary,
		SalaryMin:           salary.Min,
		SalaryMax:           salary.Max,
		SalaryCurrency:      salary.Currency,
		Benefits:            req.Benefits,
		IsRemote:            req.IsRemote,
		ApplicationsCount:   0,
	}

	err := s.repo.Create(job)
	if err != nil {
		return nil, err
	}

	return job, nil
}

func (s *jobPostService) CreateJobPostWithStatus(req *CreateJobPostRequest, employerID, employerName, employerEmail, status string) (*JobPost, error) {
	// Validate required fields
	if req.Title == "" {
		return nil, errors.New("title is required")
	}
	if req.RoleOverview == "" {
		return nil, errors.New("role overview is required")
	}
	if req.Requirements == "" {
		return nil, errors.New("requirements are required")
	}
	if req.Location == "" {
		return nil, errors.New("location is required")
	}
	if req.JobType == "" {
		return nil, errors.New("job type is required")
	}
	if req.Experience == "" {
		return nil, errors.New("experience level is required")
	}

	// Validate text field lengths (ensure they can handle substantial content)
	if len(req.RoleOverview) < 10 {
		return nil, errors.New("role overview should be at least 10 characters long")
	}
	if len(req.Requirements) < 10 {
		return nil, errors.New("requirements should be at least 10 characters long")
	}

	// Fetch employer details if not provided
	if employerName == "" || employerEmail == "" {
		employerProfile, err := s.employerRepo.GetByUserID(employerID)
		if err == nil && employerProfile != nil {
			if employerName == "" {
				employerName = employerProfile.RecruiterName // Use recruiter name (actual employer name)
			}
			if employerEmail == "" {
				employerEmail = employerProfile.OfficialEmail
			}
		}
	}

	// Validate salary
	var salary Salary
	if req.Salary.Salary != nil {
		salary = *req.Salary.Salary
		if salary.Min < 0 || salary.Max < 0 {
			return nil, errors.New("salary values cannot be negative")
		}
		if salary.Min > salary.Max {
			return nil, errors.New("minimum salary cannot be greater than maximum salary")
		}
		if salary.Currency == "" {
			salary.Currency = "USD" // Default currency
		}
	} else {
		// If no salary provided, set default values
		salary = Salary{
			Min:      0,
			Max:      0,
			Currency: "USD",
		}
	}

	// Validate application deadline
	if req.ApplicationDeadline.Before(time.Now()) {
		return nil, errors.New("application deadline cannot be in the past")
	}

	job := &JobPost{
		Title:               req.Title,
		RoleOverview:        req.RoleOverview,
		Requirements:        req.Requirements,
		Location:            req.Location,
		RequiredSkills:      req.RequiredSkills,
		EmployerID:          employerID,
		EmployerName:        employerName,
		EmployerEmail:       employerEmail,
		Status:              status, // Use the provided status
		ApplicationDeadline: req.ApplicationDeadline,
		JobType:             req.JobType,
		Experience:          req.Experience,
		Salary:              salary,
		SalaryMin:           salary.Min,
		SalaryMax:           salary.Max,
		SalaryCurrency:      salary.Currency,
		Benefits:            req.Benefits,
		IsRemote:            req.IsRemote,
		ApplicationsCount:   0,
	}

	err := s.repo.Create(job)
	if err != nil {
		return nil, err
	}

	return job, nil
}

func (s *jobPostService) UpdateJobPost(id string, req *UpdateJobPostRequest) (*JobPost, error) {
	// Get existing job
	existingJob, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Title != nil {
		existingJob.Title = *req.Title
	}
	if req.RoleOverview != nil {
		// Validate role overview length
		if len(*req.RoleOverview) < 10 {
			return nil, errors.New("role overview should be at least 10 characters long")
		}
		existingJob.RoleOverview = *req.RoleOverview
	}
	if req.Requirements != nil {
		// Validate requirements length
		if len(*req.Requirements) < 10 {
			return nil, errors.New("requirements should be at least 10 characters long")
		}
		existingJob.Requirements = *req.Requirements
	}
	if req.Location != nil {
		existingJob.Location = *req.Location
	}
	if req.RequiredSkills != nil {
		existingJob.RequiredSkills = req.RequiredSkills
	}
	if req.ApplicationDeadline != nil {
		if req.ApplicationDeadline.Before(time.Now()) {
			return nil, errors.New("application deadline cannot be in the past")
		}
		existingJob.ApplicationDeadline = *req.ApplicationDeadline
	}
	if req.JobType != nil {
		existingJob.JobType = *req.JobType
	}
	if req.Experience != nil {
		existingJob.Experience = *req.Experience
	}
	if req.Salary != nil {
		// Convert FlexibleSalary to Salary for processing
		var salary Salary
		if req.Salary.Salary != nil {
			salary = *req.Salary.Salary
		} else {
			// Handle case where FlexibleSalary is nil but the field is present
			return nil, errors.New("invalid salary format")
		}

		// Validate salary
		if salary.Min < 0 || salary.Max < 0 {
			return nil, errors.New("salary values cannot be negative")
		}
		if salary.Min > salary.Max {
			return nil, errors.New("minimum salary cannot be greater than maximum salary")
		}
		if salary.Currency == "" {
			salary.Currency = "USD"
		}
		existingJob.Salary = salary
		existingJob.SalaryMin = salary.Min
		existingJob.SalaryMax = salary.Max
		existingJob.SalaryCurrency = salary.Currency
	}
	if req.Benefits != nil {
		existingJob.Benefits = req.Benefits
	}
	if req.IsRemote != nil {
		existingJob.IsRemote = *req.IsRemote
	}

	err = s.repo.Update(existingJob)
	if err != nil {
		return nil, err
	}

	return existingJob, nil
}

func (s *jobPostService) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *jobPostService) GetByID(id string) (*JobPost, error) {
	job, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Populate employer details if missing
	s.populateEmployerDetails(job)

	return job, nil
}

// populateEmployerDetails fetches and populates employer details if they're missing
func (s *jobPostService) populateEmployerDetails(job *JobPost) {
	middleware.DebugLog("DEBUG: populateEmployerDetails called for job %s\n", job.ID)
	middleware.DebugLog("DEBUG: Current employerName: '%s', employerEmail: '%s'\n", job.EmployerName, job.EmployerEmail)

	if job.EmployerName == "" || job.EmployerEmail == "" {
		middleware.DebugLog("DEBUG: Fetching employer profile for user ID: %s\n", job.EmployerID)
		employerProfile, err := s.employerRepo.GetByUserID(job.EmployerID)
		if err != nil {
			middleware.DebugLog("DEBUG: Error fetching employer profile: %v\n", err)
		} else if employerProfile != nil {
			middleware.DebugLog("DEBUG: Found employer profile: %+v\n", employerProfile)
			if job.EmployerName == "" {
				job.EmployerName = employerProfile.RecruiterName // Use recruiter name (actual employer name)
				middleware.DebugLog("DEBUG: Set employerName to: '%s'\n", job.EmployerName)
			}
			if job.EmployerEmail == "" {
				job.EmployerEmail = employerProfile.OfficialEmail
				middleware.DebugLog("DEBUG: Set employerEmail to: '%s'\n", job.EmployerEmail)
			}
		} else {
			middleware.DebugLog("DEBUG: No employer profile found\n")
		}
	} else {
		middleware.DebugLog("DEBUG: Employer details already populated\n")
	}
}

func (s *jobPostService) GetByEmployer(employerID string) ([]JobPost, error) {
	return s.repo.GetByEmployer(employerID)
}

func (s *jobPostService) Search(filter *JobPostFilter) ([]JobPost, error) {
	// Set default pagination if not provided
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 20
	}

	jobs, err := s.repo.Search(filter)
	if err != nil {
		return nil, err
	}

	// Populate employer details for all jobs
	for i := range jobs {
		s.populateEmployerDetails(&jobs[i])
	}

	return jobs, nil
}

func (s *jobPostService) IncrementApplicationsCount(jobID string) error {
	return s.repo.IncrementApplicationsCount(jobID)
}

func (s *jobPostService) GetFeaturedJobs(limit int) ([]JobPost, error) {
	// Set default limit if not provided
	if limit <= 0 {
		limit = 10
	}

	return s.repo.GetFeaturedJobs(limit)
}

func (s *jobPostService) GetRecentJobs(limit int) ([]JobPost, error) {
	// Set default limit if not provided
	if limit <= 0 {
		limit = 20
	}

	return s.repo.GetRecentJobs(limit)
}

// Enhanced search and discovery methods
func (s *jobPostService) AdvancedSearch(request *AdvancedJobSearchRequest) (*JobSearchResponse, error) {
	// Set default values
	if request.Page <= 0 {
		request.Page = 1
	}
	if request.Limit <= 0 {
		request.Limit = 20
	} else if request.Limit > 100 {
		request.Limit = 100
	}

	// Perform search
	jobs, total, err := s.repo.AdvancedSearch(request)
	if err != nil {
		return nil, err
	}

	// Calculate pagination info
	totalPages := (total + request.Limit - 1) / request.Limit
	hasNext := request.Page < totalPages
	hasPrev := request.Page > 1

	pagination := &PaginationInfo{
		Page:       request.Page,
		Limit:      request.Limit,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    hasNext,
		HasPrev:    hasPrev,
	}

	// Get search filters for UI
	filters, _ := s.repo.GetSearchFilters()

	response := &JobSearchResponse{
		Success:    true,
		Message:    "Search completed successfully",
		Jobs:       jobs,
		Filters:    filters,
		Pagination: pagination,
	}

	return response, nil
}

func (s *jobPostService) GetSearchFilters() (*SearchFilters, error) {
	return s.repo.GetSearchFilters()
}

func (s *jobPostService) GetTrendingJobs(limit int) ([]JobPost, error) {
	if limit <= 0 {
		limit = 10
	}

	return s.repo.GetTrendingJobs(limit)
}

func (s *jobPostService) GetSimilarJobs(jobID string, maxResults int) ([]JobPost, error) {
	if maxResults <= 0 {
		maxResults = 5
	}

	return s.repo.GetSimilarJobs(jobID, maxResults)
}

func (s *jobPostService) GetRecommendedJobs(request *JobRecommendationRequest) (*JobRecommendationResponse, error) {
	// Set default values
	if request.MaxResults <= 0 {
		request.MaxResults = 10
	}

	// Get recommended jobs
	jobs, err := s.repo.GetRecommendedJobs(request)
	if err != nil {
		return nil, err
	}

	// Generate recommendation reason
	reason := "Based on your skills and preferences"
	if len(request.UserSkills) > 0 {
		reason = "Based on your skills: " + request.UserSkills[0]
		if len(request.UserSkills) > 1 {
			reason += " and others"
		}
	}

	response := &JobRecommendationResponse{
		Success: true,
		Message: "Recommendations generated successfully",
		Jobs:    jobs,
		Reason:  reason,
	}

	return response, nil
}

// Job alerts methods
func (s *jobPostService) CreateJobAlert(request *JobAlertRequest) (*JobAlert, error) {
	// Validate required fields
	if request.UserID == "" {
		return nil, errors.New("user ID is required")
	}
	if request.Frequency == "" {
		request.Frequency = "weekly"
	}

	// Validate frequency
	validFrequencies := map[string]bool{"daily": true, "weekly": true, "immediate": true}
	if !validFrequencies[request.Frequency] {
		return nil, errors.New("invalid frequency. Must be daily, weekly, or immediate")
	}

	// Create alert
	alert := &JobAlert{
		UserID:     request.UserID,
		Keywords:   request.Keywords,
		Location:   request.Location,
		JobType:    request.JobType,
		Experience: request.Experience,
		Skills:     request.Skills,
		IsRemote:   request.IsRemote,
		Frequency:  request.Frequency,
		IsActive:   request.IsActive,
	}

	// Set salary range if provided
	if request.SalaryRange != nil {
		alert.SalaryMin = &request.SalaryRange.Min
		alert.SalaryMax = &request.SalaryRange.Max
		alert.SalaryCurrency = request.SalaryRange.Currency
	}

	err := s.repo.CreateJobAlert(alert)
	if err != nil {
		return nil, err
	}

	return alert, nil
}

func (s *jobPostService) UpdateJobAlert(alertID string, request *JobAlertRequest) (*JobAlert, error) {
	// Get existing alert
	alert, err := s.repo.GetJobAlertByID(alertID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if request.Keywords != nil {
		alert.Keywords = request.Keywords
	}
	if request.Location != "" {
		alert.Location = request.Location
	}
	if request.JobType != nil {
		alert.JobType = request.JobType
	}
	if request.Experience != nil {
		alert.Experience = request.Experience
	}
	if request.Skills != nil {
		alert.Skills = request.Skills
	}
	if request.SalaryRange != nil {
		alert.SalaryMin = &request.SalaryRange.Min
		alert.SalaryMax = &request.SalaryRange.Max
		alert.SalaryCurrency = request.SalaryRange.Currency
	}
	if request.IsRemote != nil {
		alert.IsRemote = request.IsRemote
	}
	if request.Frequency != "" {
		// Validate frequency
		validFrequencies := map[string]bool{"daily": true, "weekly": true, "immediate": true}
		if !validFrequencies[request.Frequency] {
			return nil, errors.New("invalid frequency. Must be daily, weekly, or immediate")
		}
		alert.Frequency = request.Frequency
	}

	alert.IsActive = request.IsActive

	err = s.repo.UpdateJobAlert(alert)
	if err != nil {
		return nil, err
	}

	return alert, nil
}

func (s *jobPostService) DeleteJobAlert(alertID string) error {
	return s.repo.DeleteJobAlert(alertID)
}

func (s *jobPostService) GetJobAlertByID(alertID string) (*JobAlert, error) {
	return s.repo.GetJobAlertByID(alertID)
}

func (s *jobPostService) GetJobAlertsByUser(userID string) ([]JobAlert, error) {
	return s.repo.GetJobAlertsByUser(userID)
}

func (s *jobPostService) ProcessJobAlerts() error {
	// Get all active job alerts
	alerts, err := s.repo.GetActiveJobAlerts()
	if err != nil {
		return err
	}

	// Process each alert
	for _, alert := range alerts {
		// Get jobs matching this alert
		jobs, err := s.repo.GetJobsMatchingAlert(&alert)
		if err != nil {
			continue // Skip this alert if there's an error
		}

		// If there are matching jobs, send notification
		if len(jobs) > 0 {
			// TODO: Send notification to user
			// This would integrate with the notification service
			// For now, we just log it
			// log.Printf("Found %d matching jobs for alert %s (user: %s)", len(jobs), alert.ID, alert.UserID)
		}
	}

	return nil
}

// Draft methods
func (s *jobPostService) CreateDraft(req *CreateDraftRequest, employerID, employerName, employerEmail string) (*JobPost, error) {
	// Fetch employer details if not provided
	if employerName == "" || employerEmail == "" {
		employerProfile, err := s.employerRepo.GetByUserID(employerID)
		if err == nil && employerProfile != nil {
			if employerName == "" {
				employerName = employerProfile.RecruiterName // Use recruiter name (actual employer name)
			}
			if employerEmail == "" {
				employerEmail = employerProfile.OfficialEmail
			}
		}
	}

	// For drafts, all fields are optional
	job := &JobPost{
		EmployerID:        employerID,
		EmployerName:      employerName,
		EmployerEmail:     employerEmail,
		Status:            "draft",
		ApplicationsCount: 0,
	}

	// Set fields if provided
	if req.Title != nil {
		job.Title = *req.Title
	}
	if req.RoleOverview != nil {
		// Validate role overview length if provided
		if len(*req.RoleOverview) > 0 && len(*req.RoleOverview) < 10 {
			return nil, errors.New("role overview should be at least 10 characters long")
		}
		job.RoleOverview = *req.RoleOverview
	}
	if req.Requirements != nil {
		// Validate requirements length if provided
		if len(*req.Requirements) > 0 && len(*req.Requirements) < 10 {
			return nil, errors.New("requirements should be at least 10 characters long")
		}
		job.Requirements = *req.Requirements
	}
	if req.Location != nil {
		job.Location = *req.Location
	}
	if req.RequiredSkills != nil {
		job.RequiredSkills = req.RequiredSkills
	}
	if req.ApplicationDeadline != nil {
		job.ApplicationDeadline = *req.ApplicationDeadline
	}
	if req.JobType != nil {
		job.JobType = *req.JobType
	}
	if req.Experience != nil {
		job.Experience = *req.Experience
	}
	if req.Salary != nil {
		// Convert FlexibleSalary to Salary for processing
		var salary Salary
		if req.Salary.Salary != nil {
			salary = *req.Salary.Salary
		} else {
			// Handle case where FlexibleSalary is nil but the field is present
			salary = Salary{
				Min:      0,
				Max:      0,
				Currency: "USD",
			}
		}

		job.Salary = salary
		// Also set the individual fields for database storage
		job.SalaryMin = salary.Min
		job.SalaryMax = salary.Max
		job.SalaryCurrency = salary.Currency

		// Debug logging
		middleware.DebugLog("DEBUG: Salary data received - Min: %f, Max: %f, Currency: %s\n",
			salary.Min, salary.Max, salary.Currency)
		middleware.DebugLog("DEBUG: Job salary fields set - Min: %f, Max: %f, Currency: %s\n",
			job.SalaryMin, job.SalaryMax, job.SalaryCurrency)
	} else {
		middleware.DebugLog("DEBUG: No salary data provided in request\n")
	}
	if req.Benefits != nil {
		job.Benefits = req.Benefits
	}
	if req.IsRemote != nil {
		job.IsRemote = *req.IsRemote
	}

	err := s.repo.CreateDraft(job)
	if err != nil {
		return nil, err
	}

	return job, nil
}

func (s *jobPostService) GetDraftsByEmployer(employerID string) ([]JobPost, error) {
	return s.repo.GetDraftsByEmployer(employerID)
}

func (s *jobPostService) PublishDraft(jobID string) (*JobPost, error) {
	// Get the draft first
	draft, err := s.repo.GetByID(jobID)
	if err != nil {
		return nil, err
	}

	// Validate that it's actually a draft
	if draft.Status != "draft" {
		return nil, errors.New("only draft jobs can be published")
	}

	// Validate required fields before publishing
	if draft.Title == "" {
		return nil, errors.New("title is required to publish a job")
	}
	if draft.RoleOverview == "" {
		return nil, errors.New("role overview is required to publish a job")
	}
	if draft.Requirements == "" {
		return nil, errors.New("requirements are required to publish a job")
	}
	if draft.Location == "" {
		return nil, errors.New("location is required to publish a job")
	}
	if draft.JobType == "" {
		return nil, errors.New("job type is required to publish a job")
	}
	if draft.Experience == "" {
		return nil, errors.New("experience level is required to publish a job")
	}

	// Validate text field lengths
	if len(draft.RoleOverview) < 10 {
		return nil, errors.New("role overview should be at least 10 characters long")
	}
	if len(draft.Requirements) < 10 {
		return nil, errors.New("requirements should be at least 10 characters long")
	}

	// Publish the draft
	err = s.repo.PublishDraft(jobID)
	if err != nil {
		return nil, err
	}

	// Return the updated job
	return s.repo.GetByID(jobID)
}
