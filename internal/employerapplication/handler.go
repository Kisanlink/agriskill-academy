// File: internal/employerapplication/handler.go

package employerapplication

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type EmployerApplicationHandler struct {
	service EmployerApplicationService
}

func NewEmployerApplicationHandler(s EmployerApplicationService) *EmployerApplicationHandler {
	return &EmployerApplicationHandler{s}
}

// GET /employer/jobs/:jobId/applications
func (h *EmployerApplicationHandler) GetApplicationsForJob(c *gin.Context) {
	jobID := c.Param("jobId")
	apps, err := h.service.GetApplicationsForJob(jobID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to fetch applications"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "applications": apps})
}

// PUT /employer/applications/:applicationId/status
func (h *EmployerApplicationHandler) UpdateStatus(c *gin.Context) {
	applicationID := c.Param("applicationId")
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

// GET /employer/applicants/:studentId/profile
func (h *EmployerApplicationHandler) GetApplicantProfile(c *gin.Context) {
	studentID := c.Param("studentId")
	profile, err := h.service.GetApplicantProfile(studentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Applicant not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "profile": profile})
}

// POST /employer/applications/:applicationId/message
func (h *EmployerApplicationHandler) SendMessage(c *gin.Context) {
	applicationID := c.Param("applicationId")
	senderID := c.GetString("user_id") // From JWT middleware
	var req struct {
		Message string `json:"message"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Message == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid message"})
		return
	}
	msg := &Message{
		ApplicationID: applicationID,
		SenderID:      senderID,
		Message:       req.Message,
		// SentAt will be set by DB or in service layer if needed
	}
	if err := h.service.SendMessage(msg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to send message"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Message sent"})
}
