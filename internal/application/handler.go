// File: internal/application/handler.go

package application

import (
	"asa/pkg/authz"
	"fmt"
	"net/http"
	"strings"

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

// POST /jobs/:id/apply
func (h *ApplicationHandler) Apply(c *gin.Context) {
	username := c.GetString("email")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckAAAPermission(username, "db_asa_applications", "create", "", jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	jobId := c.Param("id")
	studentID := c.GetString("user_id")

	fmt.Printf("DEBUG: Apply called - JobID: %s, StudentID: %s\n", jobId, studentID)

	// Debug: Print all form fields
	fmt.Printf("DEBUG: Content-Type: %s\n", c.GetHeader("Content-Type"))
	fmt.Printf("DEBUG: Form fields: %+v\n", c.Request.Form)

	coverLetter := c.PostForm("coverLetter")
	file, err := c.FormFile("resumeFile")
	if err != nil {
		fmt.Printf("DEBUG: Resume file error: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Resume file is required"})
		return
	}

	fmt.Printf("DEBUG: Resume file received - Name: %s, Size: %d\n", file.Filename, file.Size)

	// Save file
	filename := "uploads/resumes/" + studentID + "_" + file.Filename
	if err := c.SaveUploadedFile(file, filename); err != nil {
		fmt.Printf("DEBUG: File save error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to upload resume"})
		return
	}

	fmt.Printf("DEBUG: File saved successfully to: %s\n", filename)

	app := &Application{
		JobID:       jobId,
		StudentID:   studentID,
		CoverLetter: coverLetter,
		ResumeFile:  filename,
	}

	fmt.Printf("DEBUG: Application object created: %+v\n", app)

	if err := h.service.Apply(app); err != nil {
		fmt.Printf("DEBUG: Service Apply error: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	fmt.Printf("DEBUG: Application submitted successfully\n")
	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Application submitted successfully", "application": app})
}

// GET /applications/my
func (h *ApplicationHandler) GetMyApplications(c *gin.Context) {
	username := c.GetString("username")
	studentID := c.GetString("user_id") // From JWT middleware
	jwtToken := getJWT(c)
	allowed, err := authz.CheckAAAPermission(username, "db_asa_applications", "read", studentID, jwtToken)
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

// GET /jobs/:id/applications
func (h *ApplicationHandler) GetApplicationsByJob(c *gin.Context) {
	username := c.GetString("username")
	jobID := c.Param("id")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckAAAPermission(username, "db_asa_applications", "read", jobID, jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	fmt.Printf("DEBUG: ===== REGULAR APPLICATION HANDLER CALLED =====\n")

	employerID := c.GetString("user_id")

	fmt.Printf("DEBUG: GetApplicationsByJob called - JobID: %s, EmployerID: %s\n", jobID, employerID)

	apps, err := h.service.GetApplicationsByJob(jobID, employerID)
	if err != nil {
		fmt.Printf("DEBUG: GetApplicationsByJob error: %v\n", err)
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": err.Error()})
		return
	}

	fmt.Printf("DEBUG: GetApplicationsByJob success - Found %d applications\n", len(apps))
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Applications retrieved successfully", "applications": apps})
}

// GET /applications/:applicationId
func (h *ApplicationHandler) GetApplicationByID(c *gin.Context) {
	username := c.GetString("username")
	appID := c.Param("applicationId")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckAAAPermission(username, "db_asa_applications", "read", appID, jwtToken)
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

// DELETE /applications/:applicationId
func (h *ApplicationHandler) Remove(c *gin.Context) {
	username := c.GetString("username")
	appID := c.Param("applicationId")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckAAAPermission(username, "db_asa_applications", "delete", appID, jwtToken)
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

// PUT /applications/:applicationId/status (Student can only withdraw)
func (h *ApplicationHandler) UpdateStatus(c *gin.Context) {
	username := c.GetString("username")
	appID := c.Param("applicationId")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckAAAPermission(username, "db_asa_applications", "update", appID, jwtToken)
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

// PUT /jobs/:id/applications/:applicationId/status (Employer can update status)
func (h *ApplicationHandler) UpdateStatusByEmployer(c *gin.Context) {
	username := c.GetString("username")
	appID := c.Param("applicationId")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckAAAPermission(username, "db_asa_applications", "update", appID, jwtToken)
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
