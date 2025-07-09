package employerapplication

import (
	"asa/pkg/authz"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type EmployerApplicationHandler struct {
	service EmployerApplicationService
}

func NewEmployerApplicationHandler(s EmployerApplicationService) *EmployerApplicationHandler {
	return &EmployerApplicationHandler{s}
}

func getJWT(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		return authHeader[7:]
	}
	return ""
}

func (h *EmployerApplicationHandler) GetApplicationsForJob(c *gin.Context) {
	username := c.GetString("email")
	jobID := c.Param("jobId")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckAAAPermission(username, "db_asa_applications", "read", jobID, jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	fmt.Printf("DEBUG: ===== GetApplicationsForJob HANDLER CALLED =====\n")

	status := c.Query("status")
	employerID := c.GetString("user_id")

	fmt.Printf("DEBUG: GetApplicationsForJob - JobID: %s, Status: '%s', EmployerID: %s\n", jobID, status, employerID)

	// Verify that the job belongs to the employer
	jobEmployerID, err := h.service.GetJobEmployerID(jobID)
	if err != nil {
		fmt.Printf("DEBUG: Error getting job employer ID: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to fetch applications"})
		return
	}

	fmt.Printf("DEBUG: Job employer ID: %s, Requesting employer ID: %s\n", jobEmployerID, employerID)

	if jobEmployerID != employerID {
		fmt.Printf("DEBUG: Authorization failed - job belongs to %s, requesting user is %s\n", jobEmployerID, employerID)
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Not authorized to view applications for this job"})
		return
	}

	apps, err := h.service.GetApplicationsForJob(jobID, status)
	if err != nil {
		log.Printf("DEBUG: GetApplicationsForJob error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to fetch applications"})
		return
	}

	fmt.Printf("DEBUG: GetApplicationsForJob success - Found %d applications\n", len(apps))
	fmt.Printf("DEBUG: Handler returning applications: %+v\n", apps)

	response := gin.H{"success": true, "applications": apps}
	fmt.Printf("DEBUG: Handler final response: %+v\n", response)
	c.JSON(http.StatusOK, response)
}

// Debug endpoint to test database queries
func (h *EmployerApplicationHandler) DebugApplications(c *gin.Context) {
	fmt.Printf("DEBUG: ===== DebugApplications HANDLER CALLED =====\n")

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

func (h *EmployerApplicationHandler) UpdateStatus(c *gin.Context) {
	username := c.GetString("email")
	applicationID := c.Param("applicationId")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckAAAPermission(username, "db_asa_applications", "update", applicationID, jwtToken)
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
}

func (h *EmployerApplicationHandler) GetApplicantProfile(c *gin.Context) {
	username := c.GetString("email")
	studentID := c.Param("studentId")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckAAAPermission(username, "db_asa_applications", "read", studentID, jwtToken)
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

func (h *EmployerApplicationHandler) SendMessage(c *gin.Context) {
	username := c.GetString("email")
	applicationID := c.Param("applicationId")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckAAAPermission(username, "db_asa_messages", "create", applicationID, jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	senderID := c.GetString("user_id")
	var req struct {
		Message string `json:"message"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Message == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid message"})
		return
	}

	// Verify user is authorized to send message for this application
	authorized, err := h.service.IsUserAuthorizedForApplication(applicationID, senderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to verify authorization"})
		return
	}
	if !authorized {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Not authorized to send message for this application"})
		return
	}

	msg := &Message{
		ApplicationID: applicationID,
		SenderID:      senderID,
		Message:       req.Message,
	}

	fmt.Printf("DEBUG: Creating message - database will set timestamp automatically\n")

	if err := h.service.SendMessage(msg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to send message"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Message sent"})
}

func (h *EmployerApplicationHandler) GetMessages(c *gin.Context) {
	username := c.GetString("email")
	applicationID := c.Param("applicationId")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckAAAPermission(username, "db_asa_messages", "read", applicationID, jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	userID := c.GetString("user_id")

	// Verify user is authorized to view messages for this application
	authorized, err := h.service.IsUserAuthorizedForApplication(applicationID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to verify authorization"})
		return
	}
	if !authorized {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Not authorized to view messages for this application"})
		return
	}

	messages, err := h.service.GetMessagesWithSenderInfo(applicationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Could not fetch messages"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "messages": messages})
}

func (h *EmployerApplicationHandler) GetApplicationsByStudent(c *gin.Context) {
	username := c.GetString("email")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckAAAPermission(username, "db_asa_applications", "read", "", jwtToken)
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
