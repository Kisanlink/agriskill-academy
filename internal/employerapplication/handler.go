package employerapplication

import (
	"fmt"
	"os"
	"time"

	"net/http"
	"strings"

	"github.com/Kisanlink/agriskill-academy/internal/middleware"
	"github.com/Kisanlink/agriskill-academy/internal/notification"
	"github.com/Kisanlink/agriskill-academy/internal/storage"
	"github.com/Kisanlink/agriskill-academy/pkg/authz"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type EmployerApplicationHandler struct {
	service             EmployerApplicationService
	emailSender         *notification.EmailSenderService
	db                  *gorm.DB
	storage             storage.StorageService
	notificationService notification.NotificationService
}

func NewEmployerApplicationHandler(
	s EmployerApplicationService,
	emailSender *notification.EmailSenderService,
	db *gorm.DB,
	storageService storage.StorageService,
	notificationService notification.NotificationService,
) *EmployerApplicationHandler {
	return &EmployerApplicationHandler{
		service:             s,
		emailSender:         emailSender,
		db:                  db,
		storage:             storageService,
		notificationService: notificationService,
	}
}

func getJWT(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		return authHeader[7:]
	}
	return ""
}

// GetApplicationsForJob godoc
// @Summary Get applications for a specific job
// @Description Retrieve all applications for a job posted by the authenticated employer
// @Tags employer-applications
// @Accept json
// @Produce json
// @Param jobId path string true "Job ID"
// @Param status query string false "Filter by application status (applied, shortlisted, rejected, hired)"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Applications retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden - Not authorized to view applications for this job"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/employer/jobs/{jobId}/applications [get]
func (h *EmployerApplicationHandler) GetApplicationsForJob(c *gin.Context) {
	username := c.GetString("email")
	jobID := c.Param("jobId")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckLocalPermission(username, "db_asa_applications", "read", jobID, jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	middleware.DebugLog("DEBUG: ===== GetApplicationsForJob HANDLER CALLED =====\n")

	status := c.Query("status")
	employerID := c.GetString("user_id")

	middleware.DebugLog("DEBUG: GetApplicationsForJob - JobID: %s, Status: '%s', EmployerID: %s\n", jobID, status, employerID)

	// Verify that the job belongs to the employer
	jobEmployerID, err := h.service.GetJobEmployerID(jobID)
	if err != nil {
		middleware.DebugLog("DEBUG: Error getting job employer ID: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to fetch applications"})
		return
	}

	middleware.DebugLog("DEBUG: Job employer ID: %s, Requesting employer ID: %s\n", jobEmployerID, employerID)

	if jobEmployerID != employerID {
		middleware.DebugLog("DEBUG: Authorization failed - job belongs to %s, requesting user is %s\n", jobEmployerID, employerID)
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Not authorized to view applications for this job"})
		return
	}

	apps, err := h.service.GetApplicationsForJob(jobID, status)
	if err != nil {
		middleware.DebugLog("DEBUG: GetApplicationsForJob error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to fetch applications"})
		return
	}

	middleware.DebugLog("DEBUG: GetApplicationsForJob success - Found %d applications\n", len(apps))
	middleware.DebugLog("DEBUG: Handler returning applications: %+v\n", apps)

	response := gin.H{"success": true, "applications": apps}
	middleware.DebugLog("DEBUG: Handler final response: %+v\n", response)
	c.JSON(http.StatusOK, response)
}

// DebugApplications godoc
// @Summary Debug applications for a job
// @Description Debug endpoint to test database queries for applications
// @Tags employer-applications
// @Accept json
// @Produce json
// @Param jobId path string true "Job ID"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Debug information"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/employer/jobs/{jobId}/applications/debug [get]
func (h *EmployerApplicationHandler) DebugApplications(c *gin.Context) {
	middleware.DebugLog("DEBUG: ===== DebugApplications HANDLER CALLED =====\n")

	jobID := c.Param("jobId")

	// Test simple query first
	var count int64
	err := h.service.(*employerApplicationService).repo.(*employerApplicationRepository).db.Table("applications").
		Where("job_id = ?", jobID).
		Count(&count).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"job_id":             jobID,
		"total_applications": count,
		"message":            "Debug query executed successfully",
	})
}

// UpdateStatus godoc
// @Summary Update application status
// @Description Update the status of a job application (applied, shortlisted, rejected, hired)
// @Tags employer-applications
// @Accept json
// @Produce json
// @Param applicationId path string true "Application ID"
// @Param request body map[string]string true "Status update request"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Status updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad request - Invalid status"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden - Permission denied"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/employer/applications/{applicationId}/status [put]
func (h *EmployerApplicationHandler) UpdateStatus(c *gin.Context) {
	username := c.GetString("email")
	applicationID := c.Param("applicationId")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckLocalPermission(username, "db_asa_applications", "update", applicationID, jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid status"})
		return
	}
	if err := h.service.UpdateStatus(applicationID, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to update status"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Status updated"})

	// Send email notification to student about status update
	go func() {
		middleware.DebugLog("📧 Email notification goroutine started for application: %s, status: %s", applicationID, req.Status)

		if h.emailSender == nil {
			middleware.DebugLog("⚠️  emailSender is nil, cannot send email notification")
			return
		}

		middleware.DebugLog("📧 Checking if status update email should be sent for application: %s", applicationID)

		// Get application with student and job details using flat struct
		var app struct {
			ID             string `gorm:"column:id"`
			StudentID      string `gorm:"column:student_id"`
			JobID          string `gorm:"column:job_id"`
			EmployerID     string `gorm:"column:employer_id"`
			StudentName    string `gorm:"column:student_name"`
			StudentEmail   string `gorm:"column:student_email"`
			JobTitle       string `gorm:"column:job_title"`
			CompanyName    string `gorm:"column:company_name"`
			CompanyWebsite string `gorm:"column:company_website"`
			CompanyLogo    string `gorm:"column:company_logo"`
		}

		err := h.db.Table("applications").
			Select(`
				applications.id,
				applications.student_id,
				applications.job_id,
				job_posts.employer_id,
				COALESCE(student_profiles.name, users.name) as student_name,
				users.email as student_email,
				job_posts.title as job_title,
				COALESCE(employer_profiles.company_name, job_posts.employer_name, '') as company_name,
				COALESCE(employer_profiles.website_url, '') as company_website,
				COALESCE(employer_profiles.logo_key, '') as company_logo
			`).
			Joins("LEFT JOIN users ON applications.student_id = users.id").
			Joins("LEFT JOIN student_profiles ON applications.student_id = student_profiles.user_id").
			Joins("LEFT JOIN job_posts ON applications.job_id = job_posts.id").
			Joins("LEFT JOIN employer_profiles ON job_posts.employer_id = employer_profiles.user_id").
			Where("applications.id = ?", applicationID).
			Scan(&app).Error

		if err != nil {
			middleware.DebugLog("❌ Failed to fetch application details for email: %v", err)
			return
		}

		// Validate that we got the required data
		if app.StudentID == "" || app.StudentEmail == "" {
			middleware.DebugLog("❌ Missing required application data: StudentID=%s, StudentEmail=%s", app.StudentID, app.StudentEmail)
			return
		}

		// Log fetched data for debugging
		middleware.DebugLog("📋 Fetched application data - StudentName: %s, CompanyName: %s, EmployerID: %s, LogoKey: %s",
			app.StudentName, app.CompanyName, app.EmployerID, app.CompanyLogo)

		// Check if student has email notifications enabled using notification service
		// This ensures preferences exist and are properly checked
		shouldSend, err := h.notificationService.ShouldSendNotification(app.StudentID, notification.NotificationTypeApplicationUpdate)
		if err != nil {
			middleware.DebugLog("⚠️  Failed to check notification preferences: %v, skipping email", err)
			return
		}
		if !shouldSend {
			middleware.DebugLog("ℹ️  Student has application updates disabled, skipping email")
			return
		}

		middleware.DebugLog("📧 Student has email notifications enabled, sending status update email")

		// Status messages for each status type (matching application/service.go)
		statusMessages := map[string]string{
			"applied":     "Your application has been received and is under review.",
			"viewed":      "Your application has been viewed by the employer.",
			"reviewing":   "Your application is being reviewed by the employer.",
			"shortlisted": "Congratulations! You've been shortlisted for this position.",
			"interview":   "You've been invited for an interview. The employer will contact you soon.",
			"rejected":    "Thank you for your application. Unfortunately, we've decided to move forward with other candidates.",
			"accepted":    "Congratulations! You've been selected for this position!",
			"hired":       "Congratulations! You've been selected for this position!",
			"withdrawn":   "Your application has been withdrawn.",
		}

		// Get base URL - ensure it's the backend URL, not frontend
		// For emails, we need the backend API URL where /api/files/serve/logo endpoint is hosted
		baseURL := os.Getenv("ASA_BASE_URL")
		if baseURL == "" {
			baseURL = "http://localhost:8080"
		}
		// Ensure baseURL doesn't have trailing slash and points to backend
		baseURL = strings.TrimSuffix(baseURL, "/")
		// If baseURL points to frontend (like localhost:5173), we need backend URL
		// Check if it's a frontend URL and replace with backend
		if strings.Contains(baseURL, ":5173") || strings.Contains(baseURL, ":3000") {
			// Extract host and use backend port
			parts := strings.Split(baseURL, ":")
			if len(parts) >= 2 {
				baseURL = fmt.Sprintf("%s:8080", strings.Join(parts[:len(parts)-1], ":"))
			} else {
				baseURL = "http://localhost:8080"
			}
		}

		statusMessage := statusMessages[req.Status]
		if statusMessage == "" {
			statusMessage = "Your application status has been updated."
		}

		// Build company logo URL if logo_key exists
		// The logo is stored as logo_key (S3 key) in DB (e.g., "employer_logos/1766261179581447700_...jpg")
		// Email clients cannot access API endpoints, so we need to generate presigned S3 URLs
		// Unlike EMAIL_LOGO_URL which is a direct URL, company logos need presigned URLs from S3
		companyLogoURL := ""
		if app.CompanyLogo != "" {
			// Generate presigned S3 URL (valid for 7 days)
			// This creates a direct S3 URL that email clients can access
			presignedURL, err := h.storage.GetPresignedURL(app.CompanyLogo, 7*24*time.Hour)
			if err == nil {
				companyLogoURL = presignedURL
				middleware.DebugLog("📷 Company logo presigned URL generated - LogoKey: %s, URL: %s",
					app.CompanyLogo, companyLogoURL)
			} else {
				middleware.DebugLog("⚠️  Failed to generate presigned URL for logo: %v, LogoKey: %s", err, app.CompanyLogo)
				// Fallback: use API endpoint (may not work in email clients, but better than nothing)
				if app.EmployerID != "" {
					companyLogoURL = fmt.Sprintf("%s/api/files/serve/logo/%s", baseURL, app.EmployerID)
					middleware.DebugLog("⚠️  Using API endpoint as fallback: %s", companyLogoURL)
				}
			}
		} else {
			middleware.DebugLog("⚠️  Company logo missing - LogoKey: '%s'", app.CompanyLogo)
		}

		appData := map[string]interface{}{
			"StudentName":     app.StudentName,
			"JobTitle":        app.JobTitle,
			"Company":         app.CompanyName,
			"CompanyName":     app.CompanyName,
			"CompanyWebsite":  app.CompanyWebsite,
			"CompanyLogo":     companyLogoURL,
			"Status":          req.Status,
			"StatusMessage":   statusMessage,
			"ApplicationLink": fmt.Sprintf("%s/applications/%s", baseURL, applicationID),
			"CurrentYear":     time.Now().Year(),
		}

		if err := h.emailSender.SendStatusUpdateEmail(app.StudentEmail, appData); err != nil {
			middleware.DebugLog("❌ Failed to queue status update email: %v", err)
		} else {
			middleware.DebugLog("✅ Successfully queued status update email to: %s", app.StudentEmail)
		}
	}()
}

// GetApplicantProfile godoc
// @Summary Get applicant profile
// @Description Retrieve detailed profile information for a specific applicant
// @Tags employer-applications
// @Accept json
// @Produce json
// @Param studentId path string true "Student ID"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Applicant profile retrieved successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden - Permission denied"
// @Failure 404 {object} map[string]interface{} "Applicant not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/employer/applicants/{studentId}/profile [get]
func (h *EmployerApplicationHandler) GetApplicantProfile(c *gin.Context) {
	username := c.GetString("email")
	studentID := c.Param("studentId")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckLocalPermission(username, "db_asa_applications", "read", studentID, jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	profile, err := h.service.GetApplicantProfile(studentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Applicant not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "profile": profile})
}

// GetApplicationsByStudent godoc
// @Summary Get applications by student
// @Description Retrieve all applications submitted by the authenticated student
// @Tags employer-applications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Applications retrieved successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden - Permission denied"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/employer/applications/my-applications [get]
func (h *EmployerApplicationHandler) GetApplicationsByStudent(c *gin.Context) {
	username := c.GetString("email")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckLocalPermission(username, "db_asa_applications", "read", "", jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	studentID := c.GetString("user_id")
	apps, err := h.service.GetApplicationsByStudent(studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to fetch applications"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "applications": apps})
}
