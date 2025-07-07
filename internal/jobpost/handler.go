package jobpost

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type JobPostHandler struct {
	service JobPostService
}

func NewJobPostHandler(s JobPostService) *JobPostHandler {
	return &JobPostHandler{s}
}

// POST /jobs
func (h *JobPostHandler) Create(c *gin.Context) {
	var req CreateJobPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request: " + err.Error()})
		return
	}

	employerID := c.GetString("user_id")
	if employerID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Missing employer ID"})
		return
	}

	// Get employer details from context or database
	employerName := c.GetString("user_name")
	employerEmail := c.GetString("user_email")

	job, err := h.service.CreateJobPost(&req, employerID, employerName, employerEmail)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Failed to create job: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Job created successfully", "jobPost": job})
}

// POST /jobs/draft
func (h *JobPostHandler) CreateDraft(c *gin.Context) {
	var req CreateDraftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request: " + err.Error()})
		return
	}

	// Debug logging
	fmt.Printf("DEBUG: CreateDraft request received - Salary: %+v\n", req.Salary)

	employerID := c.GetString("user_id")
	if employerID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Missing employer ID"})
		return
	}

	// Get employer details from context or database
	employerName := c.GetString("user_name")
	employerEmail := c.GetString("user_email")

	job, err := h.service.CreateDraft(&req, employerID, employerName, employerEmail)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Failed to save draft: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Draft saved successfully", "jobPost": job})
}

// POST /jobs/publish
func (h *JobPostHandler) Publish(c *gin.Context) {
	var req CreateJobPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request: " + err.Error()})
		return
	}

	employerID := c.GetString("user_id")
	if employerID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Missing employer ID"})
		return
	}

	// Get employer details from context or database
	employerName := c.GetString("user_name")
	employerEmail := c.GetString("user_email")

	job, err := h.service.CreateJobPost(&req, employerID, employerName, employerEmail)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Failed to publish job: " + err.Error()})
		return
	}

	// Set status to published
	job.Status = "published"

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Job published successfully", "jobPost": job})
}

// PUT /jobs/:id
func (h *JobPostHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req UpdateJobPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request: " + err.Error()})
		return
	}

	// Verify ownership
	employerID := c.GetString("user_id")
	if employerID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Missing employer ID"})
		return
	}

	// Get existing job to verify ownership
	existingJob, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Job not found"})
		return
	}

	if existingJob.EmployerID != employerID {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Not authorized to update this job"})
		return
	}

	job, err := h.service.UpdateJobPost(id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Failed to update job: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Job updated successfully", "jobPost": job})
}

// DELETE /jobs/:id
func (h *JobPostHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	// Verify ownership
	employerID := c.GetString("user_id")
	if employerID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Missing employer ID"})
		return
	}

	// Get existing job to verify ownership
	existingJob, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Job not found"})
		return
	}

	if existingJob.EmployerID != employerID {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Not authorized to delete this job"})
		return
	}

	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to delete job"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Job deleted successfully"})
}

// GET /jobs/:id
func (h *JobPostHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	job, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Job not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Job retrieved successfully", "jobPost": job})
}

// GET /jobs/my-posts
func (h *JobPostHandler) GetByEmployer(c *gin.Context) {
	employerID := c.GetString("user_id")
	if employerID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Missing employer ID"})
		return
	}

	jobs, err := h.service.GetByEmployer(employerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to fetch jobs"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Jobs retrieved successfully", "jobPosts": jobs})
}

// GET /jobs - Get all published jobs (for students)
func (h *JobPostHandler) GetAllJobs(c *gin.Context) {
	// Get query parameters for filtering
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "20")
	location := c.Query("location")
	jobType := c.Query("jobType")
	experience := c.Query("experience")
	isRemoteStr := c.Query("isRemote")

	// Parse pagination
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100 // Max limit
	}

	// Parse isRemote
	var isRemote *bool
	if isRemoteStr != "" {
		remote, err := strconv.ParseBool(isRemoteStr)
		if err == nil {
			isRemote = &remote
		}
	}

	// Create filter
	filter := &JobPostFilter{
		Page:       page,
		Limit:      limit,
		Location:   location,
		JobType:    []string{jobType},
		Experience: []string{experience},
		IsRemote:   isRemote,
	}

	// Remove empty filters
	if jobType == "" {
		filter.JobType = nil
	}
	if experience == "" {
		filter.Experience = nil
	}

	jobs, err := h.service.Search(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to fetch jobs"})
		return
	}

	// Transform jobs to match the expected response format
	var transformedJobs []gin.H
	for _, job := range jobs {
		// Only include published jobs
		if job.Status == "published" {
			transformedJob := gin.H{
				"id":                  job.ID,
				"title":               job.Title,
				"company":             job.EmployerName,
				"location":            job.Location,
				"jobType":             job.JobType,
				"experience":          job.Experience,
				"description":         job.RoleOverview,
				"requirements":        job.Requirements,
				"skills":              job.RequiredSkills,
				"postedAt":            job.CreatedAt,
				"applicationDeadline": job.ApplicationDeadline,
				"salary":              job.Salary,
				"recruiter": gin.H{
					"name":    job.EmployerName,
					"email":   job.EmployerEmail,
					"company": job.EmployerName,
					"avatar":  nil,
				},
				"benefits":          job.Benefits,
				"isRemote":          job.IsRemote,
				"applicationsCount": job.ApplicationsCount,
				"status":            "active", // All published jobs are considered active
			}
			transformedJobs = append(transformedJobs, transformedJob)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"jobs":    transformedJobs,
	})
}

// GET /jobs/featured
func (h *JobPostHandler) GetFeaturedJobs(c *gin.Context) {
	// Get limit from query parameter, default to 10
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	jobs, err := h.service.GetFeaturedJobs(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to fetch featured jobs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Featured jobs retrieved successfully", "jobPosts": jobs})
}

// GET /jobs/recent
func (h *JobPostHandler) GetRecentJobs(c *gin.Context) {
	// Get limit from query parameter, default to 20
	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20
	}

	jobs, err := h.service.GetRecentJobs(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to fetch recent jobs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Recent jobs retrieved successfully", "jobPosts": jobs})
}

// POST /jobs/search
func (h *JobPostHandler) Search(c *gin.Context) {
	var filter JobPostFilter
	if err := c.ShouldBindJSON(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request: " + err.Error()})
		return
	}

	jobs, err := h.service.Search(&filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Search failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Search completed successfully", "jobPosts": jobs})
}

// Enhanced Search and Discovery Endpoints

// POST /jobs/advanced-search
func (h *JobPostHandler) AdvancedSearch(c *gin.Context) {
	var request AdvancedJobSearchRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request: " + err.Error()})
		return
	}

	response, err := h.service.AdvancedSearch(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Advanced search failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GET /jobs/search-filters
func (h *JobPostHandler) GetSearchFilters(c *gin.Context) {
	filters, err := h.service.GetSearchFilters()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to get search filters"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Search filters retrieved successfully", "filters": filters})
}

// GET /jobs/trending
func (h *JobPostHandler) GetTrendingJobs(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	jobs, err := h.service.GetTrendingJobs(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to fetch trending jobs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Trending jobs retrieved successfully", "jobPosts": jobs})
}

// GET /jobs/:id/similar
func (h *JobPostHandler) GetSimilarJobs(c *gin.Context) {
	jobID := c.Param("id")
	maxResultsStr := c.DefaultQuery("maxResults", "5")
	maxResults, err := strconv.Atoi(maxResultsStr)
	if err != nil || maxResults <= 0 {
		maxResults = 5
	}

	jobs, err := h.service.GetSimilarJobs(jobID, maxResults)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to fetch similar jobs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Similar jobs retrieved successfully", "jobPosts": jobs})
}

// POST /jobs/recommendations
func (h *JobPostHandler) GetRecommendedJobs(c *gin.Context) {
	var request JobRecommendationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request: " + err.Error()})
		return
	}

	// Set user ID from context if not provided
	if request.UserID == "" {
		request.UserID = c.GetString("user_id")
	}

	response, err := h.service.GetRecommendedJobs(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to get recommendations: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Job Alerts Endpoints

// POST /jobs/alerts
func (h *JobPostHandler) CreateJobAlert(c *gin.Context) {
	var request JobAlertRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request: " + err.Error()})
		return
	}

	// Set user ID from context
	request.UserID = c.GetString("user_id")
	if request.UserID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "User not authenticated"})
		return
	}

	alert, err := h.service.CreateJobAlert(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Failed to create job alert: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Job alert created successfully", "alert": alert})
}

// PUT /jobs/alerts/:id
func (h *JobPostHandler) UpdateJobAlert(c *gin.Context) {
	alertID := c.Param("id")
	var request JobAlertRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request: " + err.Error()})
		return
	}

	// Verify ownership
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "User not authenticated"})
		return
	}

	// Get existing alert to verify ownership
	existingAlert, err := h.service.GetJobAlertByID(alertID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Job alert not found"})
		return
	}

	if existingAlert.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Not authorized to update this alert"})
		return
	}

	alert, err := h.service.UpdateJobAlert(alertID, &request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Failed to update job alert: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Job alert updated successfully", "alert": alert})
}

// DELETE /jobs/alerts/:id
func (h *JobPostHandler) DeleteJobAlert(c *gin.Context) {
	alertID := c.Param("id")

	// Verify ownership
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "User not authenticated"})
		return
	}

	// Get existing alert to verify ownership
	existingAlert, err := h.service.GetJobAlertByID(alertID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Job alert not found"})
		return
	}

	if existingAlert.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Not authorized to delete this alert"})
		return
	}

	if err := h.service.DeleteJobAlert(alertID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to delete job alert"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Job alert deleted successfully"})
}

// GET /jobs/alerts/:id
func (h *JobPostHandler) GetJobAlertByID(c *gin.Context) {
	alertID := c.Param("id")

	// Verify ownership
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "User not authenticated"})
		return
	}

	alert, err := h.service.GetJobAlertByID(alertID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Job alert not found"})
		return
	}

	if alert.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Not authorized to view this alert"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Job alert retrieved successfully", "alert": alert})
}

// GET /jobs/alerts
func (h *JobPostHandler) GetJobAlertsByUser(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "User not authenticated"})
		return
	}

	alerts, err := h.service.GetJobAlertsByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to fetch job alerts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Job alerts retrieved successfully", "alerts": alerts})
}

// Draft-specific endpoints

// GET /jobs/drafts
func (h *JobPostHandler) GetDrafts(c *gin.Context) {
	employerID := c.GetString("user_id")
	if employerID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Missing employer ID"})
		return
	}

	drafts, err := h.service.GetDraftsByEmployer(employerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to fetch drafts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Drafts retrieved successfully", "jobPosts": drafts})
}

// POST /jobs/:id/publish
func (h *JobPostHandler) PublishDraft(c *gin.Context) {
	jobID := c.Param("id")
	employerID := c.GetString("user_id")

	if employerID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Missing employer ID"})
		return
	}

	// Verify ownership
	existingJob, err := h.service.GetByID(jobID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Job not found"})
		return
	}

	if existingJob.EmployerID != employerID {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Not authorized to publish this job"})
		return
	}

	job, err := h.service.PublishDraft(jobID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Failed to publish draft: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Draft published successfully", "jobPost": job})
}
