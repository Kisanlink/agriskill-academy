package jobpost

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Kisanlink/agriskill-academy/internal/auth"
	"github.com/Kisanlink/agriskill-academy/internal/middleware"
	"github.com/Kisanlink/agriskill-academy/internal/notification"
	"github.com/Kisanlink/agriskill-academy/internal/storage"
	"github.com/Kisanlink/agriskill-academy/pkg/authz"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

type JobPostHandler struct {
	service         JobPostService
	emailSender     *notification.EmailSenderService
	db              *gorm.DB
	notificationSvc notification.NotificationService
	storageSvc      storage.StorageService
}

func NewJobPostHandler(s JobPostService, emailSender *notification.EmailSenderService, db *gorm.DB, notificationSvc notification.NotificationService, storageSvc storage.StorageService) *JobPostHandler {
	return &JobPostHandler{
		service:         s,
		emailSender:     emailSender,
		db:              db,
		notificationSvc: notificationSvc,
		storageSvc:      storageSvc,
	}
}

func getJWT(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		return authHeader[7:]
	}
	return ""
}

// @Summary Create Job Post
// @Description Create a new job post (employer only)
// @Tags Job Posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateJobPostRequest true "Job post data"
// @Success 201 {object} jobpost.JobPostResponse "Job created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Permission denied"
// @Router /api/jobs [post]
// POST /jobs
func (h *JobPostHandler) Create(c *gin.Context) {
	username := c.GetString("username")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckLocalPermission(username, "db_asa_job_posts", "create", "", jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

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
	employerName := c.GetString("username")
	employerEmail := c.GetString("email")

	job, err := h.service.CreateJobPost(c.Request.Context(), &req, employerID, employerName, employerEmail)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Failed to create job: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Job created successfully", "jobPost": job})
}

// @Summary Create Job Draft
// @Description Create a draft job post (employer only)
// @Tags Job Posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateDraftRequest true "Draft job post data"
// @Success 201 {object} jobpost.JobPostResponse "Draft saved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Permission denied"
// @Router /api/jobs/draft [post]
// POST /jobs/draft
func (h *JobPostHandler) CreateDraft(c *gin.Context) {
	username := c.GetString("username")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckLocalPermission(username, "db_asa_job_posts", "create", "", jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	var req CreateDraftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request: " + err.Error()})
		return
	}

	// Debug logging
	middleware.DebugLog("DEBUG: CreateDraft request received - Salary: %+v\n", req.Salary)

	employerID := c.GetString("user_id")
	if employerID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Missing employer ID"})
		return
	}

	// Get employer details from context or database
	employerName := c.GetString("username")
	employerEmail := c.GetString("email")

	job, err := h.service.CreateDraft(&req, employerID, employerName, employerEmail)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Failed to save draft: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Draft saved successfully", "jobPost": job})
}

// @Summary Publish Job Post
// @Description Publish a job post immediately (employer only)
// @Tags Job Posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateJobPostRequest true "Job post data"
// @Success 201 {object} jobpost.JobPostResponse "Job published successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Permission denied"
// @Router /api/jobs/publish [post]
// POST /jobs/publish
func (h *JobPostHandler) Publish(c *gin.Context) {
	username := c.GetString("username")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckLocalPermission(username, "db_asa_job_posts", "create", "", jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

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
	employerName := c.GetString("username")
	employerEmail := c.GetString("email")

	// Create job with published status directly
	job, err := h.service.CreateJobPostWithStatus(c.Request.Context(), &req, employerID, employerName, employerEmail, "published")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Failed to publish job: " + err.Error()})
		return
	}

	// Send email notifications to all students with email notifications enabled
	go func() {
		if h.emailSender != nil {
			middleware.DebugLog("📧 Starting new job email notifications for job: %s", job.ID)

			// Get all students with email notifications enabled
			var students []auth.User
			err := h.db.Joins("LEFT JOIN notification_preferences ON users.id = notification_preferences.user_id").
				Where("users.role = ? AND (notification_preferences.email_notifications = ? OR notification_preferences.email_notifications IS NULL)", "student", true).
				Find(&students).Error

			if err != nil {
				middleware.DebugLog("❌ Failed to fetch students for email notifications: %v", err)
				return
			}

			middleware.DebugLog("📧 Found %d students with email notifications enabled", len(students))

			// Get base URL from environment
			baseURL := os.Getenv("ASA_BASE_URL")
			if baseURL == "" {
				baseURL = "http://localhost:8080"
			}

			// Format salary for display
			salaryStr := ""
			if job.SalaryMin > 0 || job.SalaryMax > 0 {
				if job.SalaryMin == job.SalaryMax {
					salaryStr = fmt.Sprintf("%.0f %s", job.SalaryMin, job.SalaryCurrency)
				} else {
					salaryStr = fmt.Sprintf("%.0f - %.0f %s", job.SalaryMin, job.SalaryMax, job.SalaryCurrency)
				}
			}

			// Truncate description for email
			description := job.RoleOverview
			if len(description) > 200 {
				description = description[:200] + "..."
			}

			// Generate company logo URL if available
			var companyLogoURL string
			if job.CompanyLogoKey != "" && h.storageSvc != nil {
				// Generate a presigned URL valid for 7 days
				if url, err := h.storageSvc.GetPresignedURL(job.CompanyLogoKey, 7*24*time.Hour); err == nil {
					companyLogoURL = url
				} else {
					middleware.DebugLog("⚠️ Failed to generate presigned URL for logo: %v", err)
				}
			}

			// Send email to each student
			emailCount := 0
			for _, student := range students {
				// Check if student has job alerts enabled
				var pref notification.NotificationPreferences
				if err := h.db.Where("user_id = ?", student.ID).First(&pref).Error; err == nil {
					if !pref.JobAlerts {
						continue // Skip if job alerts disabled
					}
				}

				jobData := map[string]interface{}{
					"StudentName": student.Name,
					"JobTitle":    job.Title,
					"Company":     job.EmployerName,
					"Location":    job.Location,
					"JobType":     job.JobType,
					"Experience":  job.Experience,
					"Salary":      salaryStr,
					"Description": description,
					"JobLink":     fmt.Sprintf("%s/jobs/%s", baseURL, job.ID),
					"CompanyLogo": companyLogoURL,
				}

				if err := h.emailSender.SendNewJobEmail(student.Email, jobData); err == nil {
					emailCount++
				}
			}

			middleware.DebugLog("✅ Queued %d new job emails for job: %s", emailCount, job.ID)
		}
	}()

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Job published successfully", "jobPost": job})
}

// @Summary Update Job Post
// @Description Update an existing job post (employer only)
// @Tags Job Posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Job ID"
// @Param request body UpdateJobPostRequest true "Job post update data"
// @Success 200 {object} jobpost.JobPostResponse "Job updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Permission denied"
// @Failure 404 {object} map[string]interface{} "Job not found"
// @Router /api/jobs/{id} [put]
// PUT /jobs/:id
func (h *JobPostHandler) Update(c *gin.Context) {
	username := c.GetString("username")
	jobID := c.Param("id")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckLocalPermission(username, "db_asa_job_posts", "update", jobID, jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

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
	existingJob, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Job not found"})
		return
	}

	if existingJob.EmployerID != employerID {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Not authorized to update this job"})
		return
	}

	job, err := h.service.UpdateJobPost(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Failed to update job: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Job updated successfully", "jobPost": job})
}

// DELETE /jobs/:id
// @Summary Delete Job Post
// @Description Delete a job post (employer only)
// @Tags Job Posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Job ID"
// @Success 200 {object} map[string]interface{} "Job deleted successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Permission denied"
// @Failure 404 {object} map[string]interface{} "Job not found"
// @Failure 500 {object} map[string]interface{} "Failed to delete job"
// @Router /api/jobs/{id} [delete]
// @x-swagger-ui true
func (h *JobPostHandler) Delete(c *gin.Context) {
	username := c.GetString("username")
	jobID := c.Param("id")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckLocalPermission(username, "db_asa_job_posts", "delete", jobID, jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	id := c.Param("id")

	// Verify ownership
	employerID := c.GetString("user_id")
	if employerID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Missing employer ID"})
		return
	}

	// Get existing job to verify ownership
	existingJob, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Job not found"})
		return
	}

	if existingJob.EmployerID != employerID {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Not authorized to delete this job"})
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to delete job"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Job deleted successfully"})
}

// GET /jobs/:id
// GET /jobs/:id
// @Summary Get Job Post by ID
// @Description Retrieve a job post by its ID
// @Tags Job Posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Job ID"
// @Success 200 {object} jobpost.JobPostResponse "Job retrieved successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Permission denied"
// @Failure 404 {object} map[string]interface{} "Job not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/jobs/{id} [get]
// @x-swagger-ui true
func (h *JobPostHandler) GetByID(c *gin.Context) {
	username := c.GetString("username")
	jobID := c.Param("id")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckLocalPermission(username, "db_asa_job_posts", "read", jobID, jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	id := c.Param("id")
	job, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Job not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Job retrieved successfully", "jobPost": job})
}

// GET /jobs/my-posts
// @Summary Get Jobs by Employer
// @Description Retrieve all job posts created by the current employer
// @Tags Job Posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} jobpost.JobPostResponse "Jobs retrieved successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Permission denied"
// @Failure 500 {object} map[string]interface{} "Failed to fetch jobs"
// @Router /api/jobs/my-posts [get]
// @x-swagger-ui true
func (h *JobPostHandler) GetByEmployer(c *gin.Context) {
	username := c.GetString("username")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckLocalPermission(username, "db_asa_job_posts", "read", "", jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

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
// GET /jobs - Get all published jobs (for students)
// @Summary Get all published jobs
// @Description Retrieve all published job posts (for students)
// @Tags Job Posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param limit query int false "Number of jobs per page"
// @Param location query string false "Job location filter"
// @Param jobType query string false "Job type filter"
// @Param experience query string false "Experience filter"
// @Param isRemote query bool false "Remote job filter"
// @Success 200 {object} jobpost.JobPostResponse "Jobs retrieved successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Permission denied"
// @Failure 500 {object} map[string]interface{} "Failed to fetch jobs"
// @Router /api/jobs [get]
// @x-swagger-ui true
func (h *JobPostHandler) GetAllJobs(c *gin.Context) {
	username := c.GetString("username")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckLocalPermission(username, "db_asa_job_posts", "read", "", jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	// Get query parameters for filtering
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "20")
	location := c.Query("location")
	jobType := c.Query("jobType")
	experience := c.Query("experience")
	isRemoteStr := c.Query("isRemote")
	skills := c.QueryArray("skills")
	salaryMinStr := c.Query("salaryMin")
	salaryMaxStr := c.Query("salaryMax")
	postedWithin := c.Query("postedWithin")

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

	// Parse salary range
	var salaryRange *struct {
		Min float64 `json:"min"`
		Max float64 `json:"max"`
	}
	if salaryMinStr != "" || salaryMaxStr != "" {
		salaryRange = &struct {
			Min float64 `json:"min"`
			Max float64 `json:"max"`
		}{}

		if salaryMinStr != "" {
			if min, err := strconv.ParseFloat(salaryMinStr, 64); err == nil {
				salaryRange.Min = min
			}
		}
		if salaryMaxStr != "" {
			if max, err := strconv.ParseFloat(salaryMaxStr, 64); err == nil {
				salaryRange.Max = max
			}
		}
	}

	// Create filter with enhanced parameters
	filter := &JobPostFilter{
		Page:         page,
		Limit:        limit,
		Location:     location,
		JobType:      []string{jobType},
		Experience:   []string{experience},
		IsRemote:     isRemote,
		Skills:       skills,
		SalaryRange:  salaryRange,
		PostedWithin: postedWithin,
	}

	// Remove empty filters
	if jobType == "" {
		filter.JobType = nil
	}
	if experience == "" {
		filter.Experience = nil
	}
	if len(skills) == 0 {
		filter.Skills = nil
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
			// Use CompanyName if available, otherwise fallback to EmployerName
			companyName := job.CompanyName
			if companyName == "" {
				companyName = job.EmployerName
			}
			// Debug: Log the job ID
			middleware.DebugLog("DEBUG GetAllJobs: Job ID='%s', Title='%s', CompanyName='%s', EmployerName='%s'\n",
				job.ID, job.Title, job.CompanyName, job.EmployerName)
			transformedJob := gin.H{
				"id":                  job.ID,
				"title":               job.Title,
				"company":             companyName,
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
					"company": companyName,
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
// GET /jobs/featured
// @Summary Get Featured Jobs
// @Description Retrieve a list of featured job posts
// @Tags Job Posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Maximum number of featured jobs to return (default 10)"
// @Success 200 {object} jobpost.FeaturedJobsResponse "Featured jobs retrieved successfully"
// @Failure 403 {object} map[string]interface{} "Permission denied"
// @Failure 500 {object} map[string]interface{} "Failed to fetch featured jobs"
// @Router /api/jobs/featured [get]
// @x-swagger-ui true
func (h *JobPostHandler) GetFeaturedJobs(c *gin.Context) {
	username := c.GetString("username")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckLocalPermission(username, "db_asa_job_posts", "read", "", jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

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
// @Summary Get Recent Jobs
// @Description Retrieve a list of recent job posts
// @Tags Job Posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Maximum number of recent jobs to return (default 20)"
// @Success 200 {object} jobpost.RecentJobsResponse "Recent jobs retrieved successfully"
// @Failure 403 {object} map[string]interface{} "Permission denied"
// @Failure 500 {object} map[string]interface{} "Failed to fetch recent jobs"
// @Router /api/jobs/recent [get]
// @x-swagger-ui true
func (h *JobPostHandler) GetRecentJobs(c *gin.Context) {
	username := c.GetString("username")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckLocalPermission(username, "db_asa_job_posts", "read", "", jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

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
// @Summary Search Job Posts
// @Description Search for job posts using filters
// @Tags Job Posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body JobPostFilter true "Search filters"
// @Success 200 {object} jobpost.JobPostResponse "Search completed successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 403 {object} map[string]interface{} "Permission denied"
// @Failure 500 {object} map[string]interface{} "Search failed"
// @Router /api/jobs/search [post]
// @x-swagger-ui true
func (h *JobPostHandler) Search(c *gin.Context) {
	username := c.GetString("username")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckLocalPermission(username, "db_asa_job_posts", "read", "", jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

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
// @Summary Advanced Job Search
// @Description Perform an advanced search for job posts using complex filters
// @Tags Job Posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body AdvancedJobSearchRequest true "Advanced search filters"
// @Success 200 {object} jobpost.JobSearchResponse "Advanced search completed successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 403 {object} map[string]interface{} "Permission denied"
// @Failure 500 {object} map[string]interface{} "Advanced search failed"
// @Router /api/jobs/advanced-search [post]
// @x-swagger-ui true
func (h *JobPostHandler) AdvancedSearch(c *gin.Context) {
	username := c.GetString("username")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckLocalPermission(username, "db_asa_job_posts", "read", "", jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

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
// @Summary Get Job Search Filters
// @Description Retrieve available filters for job search (locations, job types, etc.)
// @Tags Job Posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} jobpost.SearchFilters "Search filters retrieved successfully"
// @Failure 403 {object} map[string]interface{} "Permission denied"
// @Failure 500 {object} map[string]interface{} "Failed to get search filters"
// @Router /api/jobs/search-filters [get]
// @x-swagger-ui true
func (h *JobPostHandler) GetSearchFilters(c *gin.Context) {
	username := c.GetString("username")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckLocalPermission(username, "db_asa_job_posts", "read", "", jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	filters, err := h.service.GetSearchFilters()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to get search filters"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Search filters retrieved successfully", "filters": filters})
}

// GET /jobs/trending
// @Summary Get Trending Jobs
// @Description Retrieve a list of trending job posts
// @Tags Job Posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Number of trending jobs to return"
// @Success 200 {object} jobpost.TrendingJobsResponse "Trending jobs retrieved successfully"
// @Failure 403 {object} map[string]interface{} "Permission denied"
// @Failure 500 {object} map[string]interface{} "Failed to fetch trending jobs"
// @Router /api/jobs/trending [get]
// @x-swagger-ui true
func (h *JobPostHandler) GetTrendingJobs(c *gin.Context) {
	username := c.GetString("username")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckLocalPermission(username, "db_asa_job_posts", "read", "", jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

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
// @Summary Get Similar Jobs
// @Description Retrieve a list of jobs similar to the specified job
// @Tags Job Posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Job ID"
// @Param maxResults query int false "Maximum number of similar jobs to return"
// @Success 200 {object} jobpost.SimilarJobsResponse "Similar jobs retrieved successfully"
// @Failure 403 {object} map[string]interface{} "Permission denied"
// @Failure 500 {object} map[string]interface{} "Failed to fetch similar jobs"
// @Router /api/jobs/{id}/similar [get]
// @x-swagger-ui true
func (h *JobPostHandler) GetSimilarJobs(c *gin.Context) {
	username := c.GetString("username")
	jobID := c.Param("id")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckLocalPermission(username, "db_asa_job_posts", "read", jobID, jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	jobID = c.Param("id")
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
// @Summary Get Recommended Jobs
// @Description Retrieve a list of recommended jobs for the user
// @Tags Job Posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body JobRecommendationRequest true "Recommendation request"
// @Success 200 {object} jobpost.JobRecommendationResponse "Recommended jobs retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 403 {object} map[string]interface{} "Permission denied"
// @Failure 500 {object} map[string]interface{} "Failed to get recommendations"
// @Router /api/jobs/recommendations [post]
// @x-swagger-ui true
func (h *JobPostHandler) GetRecommendedJobs(c *gin.Context) {
	username := c.GetString("username")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckLocalPermission(username, "db_asa_job_posts", "read", "", jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

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
// @Summary Create Job Alert
// @Description Create a new job alert for the authenticated user
// @Tags Job Posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body JobAlertRequest true "Job alert data"
// @Success 201 {object} jobpost.JobAlertResponse "Job alert created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "User not authenticated"
// @Failure 403 {object} map[string]interface{} "Permission denied"
// @Failure 500 {object} map[string]interface{} "Failed to create job alert"
// @Router /api/jobs/alerts [post]
// @x-swagger-ui true
func (h *JobPostHandler) CreateJobAlert(c *gin.Context) {
	username := c.GetString("username")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckLocalPermission(username, "db_asa_job_alerts", "create", "", jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

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
// @Summary Update Job Alert
// @Description Update an existing job alert for the authenticated user
// @Tags Job Posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Job Alert ID"
// @Param request body JobAlertRequest true "Job alert data"
// @Success 200 {object} jobpost.JobAlertResponse "Job alert updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "User not authenticated"
// @Failure 403 {object} map[string]interface{} "Permission denied"
// @Failure 404 {object} map[string]interface{} "Job alert not found"
// @Failure 500 {object} map[string]interface{} "Failed to update job alert"
// @Router /api/jobs/alerts/{id} [put]
// @x-swagger-ui true
func (h *JobPostHandler) UpdateJobAlert(c *gin.Context) {
	username := c.GetString("username")
	alertID := c.Param("id")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckLocalPermission(username, "db_asa_job_alerts", "update", alertID, jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	alertID = c.Param("id")
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
// @Summary Delete a job alert
// @Description Deletes a job alert by its ID. Only the owner of the alert can delete it.
// @Tags jobs, alerts
// @Param id path string true "Job Alert ID"
// @Success 200 {object} map[string]interface{} "Job alert deleted successfully"
// @Failure 401 {object} map[string]interface{} "User not authenticated"
// @Failure 403 {object} map[string]interface{} "Permission denied or not authorized to delete this alert"
// @Failure 404 {object} map[string]interface{} "Job alert not found"
// @Failure 500 {object} map[string]interface{} "Failed to delete job alert"
// @Router /api/jobs/alerts/{id} [delete]
// @x-swagger-ui true
func (h *JobPostHandler) DeleteJobAlert(c *gin.Context) {
	username := c.GetString("username")
	alertID := c.Param("id")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckLocalPermission(username, "db_asa_job_alerts", "delete", alertID, jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	alertID = c.Param("id")

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
// @Summary Get a job alert by ID
// @Description Retrieves a job alert by its ID. Only the owner of the alert can view it.
// @Tags jobs, alerts
// @Param id path string true "Job Alert ID"
// @Success 200 {object} jobpost.JobAlertResponse "Job alert retrieved successfully"
// @Failure 401 {object} map[string]interface{} "User not authenticated"
// @Failure 403 {object} map[string]interface{} "Permission denied or not authorized to view this alert"
// @Failure 404 {object} map[string]interface{} "Job alert not found"
// @Failure 500 {object} map[string]interface{} "Failed to retrieve job alert"
// @Router /api/jobs/alerts/{id} [get]
// @x-swagger-ui true
func (h *JobPostHandler) GetJobAlertByID(c *gin.Context) {
	username := c.GetString("username")
	alertID := c.Param("id")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckLocalPermission(username, "db_asa_job_alerts", "read", alertID, jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	alertID = c.Param("id")

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
// @Summary Get job alerts by user
// @Description Retrieves all job alerts for the authenticated user.
// @Tags jobs, alerts
// @Produce json
// @Security BearerAuth
// @Success 200 {object} jobpost.JobAlertResponse "Job alerts retrieved successfully"
// @Failure 401 {object} map[string]interface{} "User not authenticated"
// @Failure 403 {object} map[string]interface{} "Permission denied"
// @Failure 500 {object} map[string]interface{} "Failed to fetch job alerts"
// @Router /api/jobs/alerts [get]
// @x-swagger-ui true
func (h *JobPostHandler) GetJobAlertsByUser(c *gin.Context) {
	username := c.GetString("username")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckLocalPermission(username, "db_asa_job_alerts", "read", "", jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

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
// @Summary Get job drafts by employer
// @Description Retrieves all job drafts for the authenticated employer.
// @Tags Job Posts
// @Produce json
// @Security BearerAuth
// @Success 200 {object} jobpost.JobPostResponse "Drafts retrieved successfully"
// @Failure 401 {object} map[string]interface{} "Missing employer ID"
// @Failure 403 {object} map[string]interface{} "Permission denied"
// @Failure 500 {object} map[string]interface{} "Failed to fetch drafts"
// @Router /api/jobs/drafts [get]
// @x-swagger-ui true
func (h *JobPostHandler) GetDrafts(c *gin.Context) {
	username := c.GetString("username")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckLocalPermission(username, "db_asa_job_posts", "read", "", jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

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
// @Summary Publish a job draft
// @Description Publishes a job draft and makes it visible as a published job post (employer only).
// @Tags Job Posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Job ID"
// @Success 200 {object} jobpost.JobPostResponse "Draft published successfully"
// @Failure 400 {object} map[string]interface{} "Failed to publish draft"
// @Failure 401 {object} map[string]interface{} "Missing employer ID"
// @Failure 403 {object} map[string]interface{} "Permission denied or not authorized to publish this job"
// @Failure 404 {object} map[string]interface{} "Job not found"
// @Router /api/jobs/{id}/publish [post]
// @x-swagger-ui true
func (h *JobPostHandler) PublishDraft(c *gin.Context) {
	username := c.GetString("username")
	jobID := c.Param("id")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckLocalPermission(username, "db_asa_job_posts", "update", jobID, jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	jobID = c.Param("id")
	employerID := c.GetString("user_id")

	if employerID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Missing employer ID"})
		return
	}

	// Verify ownership
	existingJob, err := h.service.GetByID(c.Request.Context(), jobID)
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

	// Send email notifications to all students with email notifications enabled
	go func() {
		if h.emailSender != nil {
			middleware.DebugLog("📧 Starting new job email notifications for job: %s", job.ID)

			// Get all students with email notifications enabled
			var students []auth.User
			err := h.db.Joins("LEFT JOIN notification_preferences ON users.id = notification_preferences.user_id").
				Where("users.role = ? AND (notification_preferences.email_notifications = ? OR notification_preferences.email_notifications IS NULL)", "student", true).
				Find(&students).Error

			if err != nil {
				middleware.DebugLog("❌ Failed to fetch students for email notifications: %v", err)
				return
			}

			middleware.DebugLog("📧 Found %d students with email notifications enabled", len(students))

			// Get base URL from environment
			baseURL := os.Getenv("ASA_BASE_URL")
			if baseURL == "" {
				baseURL = "http://localhost:8080"
			}

			// Format salary for display
			salaryStr := ""
			if job.SalaryMin > 0 || job.SalaryMax > 0 {
				if job.SalaryMin == job.SalaryMax {
					salaryStr = fmt.Sprintf("%.0f %s", job.SalaryMin, job.SalaryCurrency)
				} else {
					salaryStr = fmt.Sprintf("%.0f - %.0f %s", job.SalaryMin, job.SalaryMax, job.SalaryCurrency)
				}
			}

			// Truncate description for email
			description := job.RoleOverview
			if len(description) > 200 {
				description = description[:200] + "..."
			}

			// Send email to each student
			emailCount := 0
			for _, student := range students {
				// Check if student has job alerts enabled
				var pref notification.NotificationPreferences
				if err := h.db.Where("user_id = ?", student.ID).First(&pref).Error; err == nil {
					if !pref.JobAlerts {
						continue // Skip if job alerts disabled
					}
				}

				jobData := map[string]interface{}{
					"StudentName": student.Name,
					"JobTitle":    job.Title,
					"Company":     job.EmployerName,
					"Location":    job.Location,
					"JobType":     job.JobType,
					"Experience":  job.Experience,
					"Salary":      salaryStr,
					"Description": description,
					"JobLink":     fmt.Sprintf("%s/jobs/%s", baseURL, job.ID),
				}

				if err := h.emailSender.SendNewJobEmail(student.Email, jobData); err == nil {
					emailCount++
				}
			}

			middleware.DebugLog("✅ Queued %d new job emails for job: %s", emailCount, job.ID)
		}
	}()

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Draft published successfully", "jobPost": job})
}

// POST /jobs/:id/close
// @Summary Close Job Post
// @Description Manually close a job post by setting status to completed (employer only)
// @Tags Job Posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Job ID"
// @Success 200 {object} map[string]interface{} "Job closed successfully"
// @Failure 403 {object} map[string]interface{} "Not authorized"
// @Failure 404 {object} map[string]interface{} "Job not found"
// @Failure 500 {object} map[string]interface{} "Failed to close job"
// @Router /api/jobs/{id}/close [post]
//
// Expected behavior validation:
// 1. Verify employer authentication and extract user_id from context
// 2. Retrieve job by ID and verify it exists
// 3. Check that the job belongs to the authenticated employer (job.EmployerID == employerID)
// 4. Set job status to "completed" and record completed_at timestamp
// 5. Return success response with job_id
//
// Test scenarios:
// - Success: Employer closes their own job
// - Failure: Non-owner employer tries to close job (403 Forbidden)
// - Failure: Job ID not found (404 Not Found)
// - Failure: Unauthenticated request (401 Unauthorized)
func (h *JobPostHandler) CloseJob(c *gin.Context) {
	employerID := c.GetString("user_id")
	jobID := c.Param("id")

	// Verify job exists and check ownership
	job, err := h.service.GetByID(c.Request.Context(), jobID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Job not found"})
		return
	}

	if job.EmployerID != employerID {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Not authorized"})
		return
	}

	// Close the job
	if err := h.service.CloseJob(jobID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Job closed successfully", "job_id": jobID})
}

// POST /jobs/:id/reopen
// @Summary Reopen Job Post
// @Description Reopen a closed job post by setting status back to published (employer only)
// @Tags Job Posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Job ID"
// @Success 200 {object} map[string]interface{} "Job reopened successfully"
// @Failure 403 {object} map[string]interface{} "Not authorized"
// @Failure 404 {object} map[string]interface{} "Job not found"
// @Failure 500 {object} map[string]interface{} "Failed to reopen job"
// @Router /api/jobs/{id}/reopen [post]
//
// Expected behavior validation:
// 1. Verify employer authentication and extract user_id from context
// 2. Retrieve job by ID and verify it exists
// 3. Check that the job belongs to the authenticated employer (job.EmployerID == employerID)
// 4. Set job status back to "published" (no vacancy_count validation needed)
// 5. Return success response with job_id
//
// Test scenarios:
// - Success: Employer reopens their previously closed job
// - Failure: Non-owner employer tries to reopen job (403 Forbidden)
// - Failure: Job ID not found (404 Not Found)
// - Failure: Unauthenticated request (401 Unauthorized)
func (h *JobPostHandler) ReopenJob(c *gin.Context) {
	employerID := c.GetString("user_id")
	jobID := c.Param("id")

	// Verify job exists and check ownership
	job, err := h.service.GetByID(c.Request.Context(), jobID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Job not found"})
		return
	}

	if job.EmployerID != employerID {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Not authorized"})
		return
	}

	// Reopen the job
	if err := h.service.ReopenJob(jobID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Job reopened successfully", "job_id": jobID})
}
