// File: internal/application/handler.go

package application

import (
	"github.com/Kisanlink/agriskill-academy/internal/middleware"
	"github.com/Kisanlink/agriskill-academy/pkg/authz"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"mime/multipart"

	"github.com/gin-gonic/gin"
)

type ApplicationHandler struct {
	service ApplicationService
}

func NewApplicationHandler(s ApplicationService) *ApplicationHandler {
	return &ApplicationHandler{s}
}

func getJWT(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		return authHeader[7:]
	}
	return ""
}

// @Summary Apply for Job
// @Description Apply for a job with resume upload (student only)
// @Tags Applications
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path string true "Job ID"
// @Param cover_letter formData string false "Cover letter text"
// @Param resume_file formData file true "Resume file (PDF, DOC, DOCX, max 10MB)"
// @Success 201 {object} map[string]interface{} "Application submitted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request or file format"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Permission denied"
// @Router /api/jobs/{id}/apply [post]
// POST /jobs/:id/apply
func (h *ApplicationHandler) Apply(c *gin.Context) {
	username := c.GetString("email")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckLocalPermission(username, "db_asa_applications", "create", "", jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	jobId := c.Param("id")
	studentID := c.GetString("user_id")

	middleware.DebugLog("DEBUG: Apply called - JobID: %s, StudentID: %s\n", jobId, studentID)
	middleware.DebugLog("DEBUG: Content-Type: %s\n", c.GetHeader("Content-Type"))

	// Parse multipart form data
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil { // 32MB max
		middleware.DebugLog("DEBUG: ParseMultipartForm error: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Failed to parse form data"})
		return
	}

	middleware.DebugLog("DEBUG: Form fields: %+v\n", c.Request.Form)

	coverLetter := c.PostForm("cover_letter")
	middleware.DebugLog("DEBUG: Cover letter: %s\n", coverLetter)

	// Try multiple possible field names for the resume file
	var file *multipart.FileHeader
	var fileErr error

	// Try different field names (prioritize snake_case)
	file, fileErr = c.FormFile("resume_file")
	if fileErr != nil {
		middleware.DebugLog("DEBUG: resume_file not found, trying resumeFile\n")
		file, fileErr = c.FormFile("resumeFile")
		if fileErr != nil {
			middleware.DebugLog("DEBUG: resumeFile not found, trying resume\n")
			file, fileErr = c.FormFile("resume")
			if fileErr != nil {
				middleware.DebugLog("DEBUG: resume not found, trying file\n")
				file, fileErr = c.FormFile("file")
				if fileErr != nil {
					middleware.DebugLog("DEBUG: Resume file error - tried resume_file, resumeFile, resume, file: %v\n", fileErr)
					c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Resume file is required"})
					return
				}
			}
		}
	}

	middleware.DebugLog("DEBUG: Resume file received - Name: %s, Size: %d\n", file.Filename, file.Size)

	// Validate file type
	if !IsValidResumeFile(file.Filename) {
		middleware.DebugLog("DEBUG: Invalid file type: %s\n", file.Filename)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid file type. Allowed: PDF, DOC, DOCX",
		})
		return
	}

	// Validate file size (10MB max)
	if file.Size > 10*1024*1024 {
		middleware.DebugLog("DEBUG: File too large: %d bytes\n", file.Size)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "File size exceeds maximum allowed size (10MB)",
		})
		return
	}

	// Get file metadata
	fileName := file.Filename
	fileType := file.Header.Get("Content-Type")
	if fileType == "" {
		fileType = getMimeTypeFromExtension(fileName)
	}
	fileSize := file.Size

	// Upload file to S3 and get the key
	resumeKey, err := h.service.UploadResumeToS3(file, studentID)
	if err != nil {
		middleware.DebugLog("DEBUG: Failed to upload resume to S3: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to upload resume"})
		return
	}

	app := NewApplication() // Use constructor to generate ID
	app.JobID = jobId
	app.StudentID = studentID
	app.CoverLetter = coverLetter
	app.ResumeKey = resumeKey
	app.ResumeFileName = fileName
	app.ResumeFileType = fileType
	app.ResumeFileSize = fileSize

	middleware.DebugLog("DEBUG: Application object created - JobID: %s, StudentID: %s, ResumeFileName: %s, ResumeFileSize: %d\n", app.JobID, app.StudentID, app.ResumeFileName, app.ResumeFileSize)

	if err := h.service.Apply(app); err != nil {
		middleware.DebugLog("DEBUG: Service Apply error: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	middleware.DebugLog("DEBUG: Application submitted successfully with ID: %s\n", app.ID)
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Application submitted successfully",
		"application": gin.H{
			"id":               app.ID,
			"job_id":           app.JobID,
			"student_id":       app.StudentID,
			"resume_key":       app.ResumeKey,
			"resume_file_name": app.ResumeFileName,
			"resume_file_type": app.ResumeFileType,
			"resume_file_size": app.ResumeFileSize,
			"file_url":         fmt.Sprintf("/api/files/serve/application-resume/%s", app.ID),
		},
	})
}

// Helper function to validate resume file type
func IsValidResumeFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	allowedTypes := []string{".pdf", ".doc", ".docx"}
	for _, allowedExt := range allowedTypes {
		if ext == allowedExt {
			return true
		}
	}
	return false
}

// Helper function to get MIME type from file extension
func getMimeTypeFromExtension(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".pdf":
		return "application/pdf"
	case ".doc":
		return "application/msword"
	case ".docx":
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	default:
		return "application/octet-stream"
	}
}

// @Summary Get My Applications
// @Description Get all applications submitted by the current student
// @Tags Applications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Applications fetched successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Permission denied"
// @Router /api/applications/my [get]
// GET /applications/my
func (h *ApplicationHandler) GetMyApplications(c *gin.Context) {
	username := c.GetString("username")
	studentID := c.GetString("user_id") // From JWT middleware
	jwtToken := getJWT(c)
	allowed, err := authz.CheckLocalPermission(username, "db_asa_applications", "read", studentID, jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}
	apps, err := h.service.GetMyApplications(studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Could not fetch applications"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Applications retrieved successfully", "applications": apps})
}

// @Summary Get Applications By Job
// @Description Get all applications submitted for a specific job (employer only)
// @Tags Applications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Job ID"
// @Success 200 {object} map[string]interface{} "Applications fetched successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Permission denied"
// @Router /api/jobs/{id}/applications [get]
// @x-swagger-ui true
// GET /jobs/:id/applications
func (h *ApplicationHandler) GetApplicationsByJob(c *gin.Context) {
	username := c.GetString("username")
	jobID := c.Param("id")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckLocalPermission(username, "db_asa_applications", "read", jobID, jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	middleware.DebugLog("DEBUG: ===== REGULAR APPLICATION HANDLER CALLED =====\n")

	employerID := c.GetString("user_id")

	middleware.DebugLog("DEBUG: GetApplicationsByJob called - JobID: %s, EmployerID: %s\n", jobID, employerID)

	apps, err := h.service.GetApplicationsByJob(jobID, employerID)
	if err != nil {
		middleware.DebugLog("DEBUG: GetApplicationsByJob error: %v\n", err)
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": err.Error()})
		return
	}

	middleware.DebugLog("DEBUG: GetApplicationsByJob success - Found %d applications\n", len(apps))
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Applications retrieved successfully", "applications": apps})
}

// @Summary Get Application By ID
// @Description Get a job application by its ID (student or employer)
// @Tags Applications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param applicationId path string true "Application ID"
// @Success 200 {object} map[string]interface{} "Application retrieved successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Permission denied"
// @Router /api/applications/{applicationId} [get]
// @x-swagger-ui true
func (h *ApplicationHandler) GetApplicationByID(c *gin.Context) {
	username := c.GetString("username")
	appID := c.Param("applicationId")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckLocalPermission(username, "db_asa_applications", "read", appID, jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	userID := c.GetString("user_id")

	app, err := h.service.GetApplicationByID(appID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Application not found"})
		return
	}

	// Check if user is authorized to view this application
	if app.StudentID != userID {
		// Check if user is the employer for this job
		jobEmployerID, err := h.service.(*applicationService).repo.GetJobEmployerID(app.JobID)
		if err != nil || jobEmployerID != userID {
			c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Not authorized to view this application"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Application retrieved successfully", "application": app})
}

// @Summary Remove Application
// @Description Delete a job application (student only)
// @Tags Applications
// @Security BearerAuth
// @Param applicationId path string true "Application ID"
// @Success 200 {object} map[string]interface{} "Application removed successfully"
// @Failure 403 {object} map[string]interface{} "Permission denied"
// @Failure 500 {object} map[string]interface{} "Could not remove application"
// @Router /api/applications/{applicationId} [delete]
// @x-swagger-ui true
func (h *ApplicationHandler) Remove(c *gin.Context) {
	username := c.GetString("username")
	appID := c.Param("applicationId")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckLocalPermission(username, "db_asa_applications", "delete", appID, jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	studentID := c.GetString("user_id") // From JWT middleware
	if err := h.service.Remove(appID, studentID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Could not remove application"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Application removed successfully"})
}

// @Summary Update Application Status
// @Description Update the status of a job application (student only)
// @Tags Applications
// @Security BearerAuth
// @Param applicationId path string true "Application ID"
// @Param status body UpdateStatusRequest true "Status update request"
// @Success 200 {object} map[string]interface{} "Application status updated successfully"
// @Failure 403 {object} map[string]interface{} "Permission denied"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Router /api/applications/{applicationId}/status [put]
// @x-swagger-ui true
func (h *ApplicationHandler) UpdateStatus(c *gin.Context) {
	username := c.GetString("username")
	appID := c.Param("applicationId")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckLocalPermission(username, "db_asa_applications", "update", appID, jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	studentID := c.GetString("user_id")

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request: " + err.Error()})
		return
	}

	if err := h.service.UpdateStatus(appID, studentID, req.Status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Application status updated successfully"})
}

// @Summary Update Application Status by Employer
// @Description Update the status of a job application (employer only)
// @Tags Applications
// @Security BearerAuth
// @Param id path string true "Job ID"
// @Param applicationId path string true "Application ID"
// @Param status body UpdateStatusRequest true "Status update request"
// @Success 200 {object} map[string]interface{} "Application status updated successfully"
// @Failure 403 {object} map[string]interface{} "Permission denied"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Router /api/jobs/{id}/applications/{applicationId}/status [put]
// @x-swagger-ui true
func (h *ApplicationHandler) UpdateStatusByEmployer(c *gin.Context) {
	username := c.GetString("username")
	appID := c.Param("applicationId")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckLocalPermission(username, "db_asa_applications", "update", appID, jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	jobID := c.Param("id")
	employerID := c.GetString("user_id")

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request: " + err.Error()})
		return
	}

	if err := h.service.UpdateStatusByEmployer(appID, jobID, employerID, req.Status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Application status updated successfully"})
}
