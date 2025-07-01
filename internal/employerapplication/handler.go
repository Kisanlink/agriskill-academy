package employerapplication

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type EmployerApplicationHandler struct {
	service EmployerApplicationService
}

func NewEmployerApplicationHandler(s EmployerApplicationService) *EmployerApplicationHandler {
	return &EmployerApplicationHandler{s}
}

func (h *EmployerApplicationHandler) GetApplicationsForJob(c *gin.Context) {
	jobID := c.Param("jobId")
	status := c.Query("status")
	fmt.Println("Status passed to repo:", status)
	apps, err := h.service.GetApplicationsForJob(jobID, status)
	if err != nil {
		log.Println("GetApplicationsForJob error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to fetch applications"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "applications": apps})
}

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

func (h *EmployerApplicationHandler) GetApplicantProfile(c *gin.Context) {
	studentID := c.Param("studentId")
	profile, err := h.service.GetApplicantProfile(studentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Applicant not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "profile": profile})
}

func (h *EmployerApplicationHandler) SendMessage(c *gin.Context) {
	applicationID := c.Param("applicationId")
	senderID := c.GetString("user_id")
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
	}
	if err := h.service.SendMessage(msg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to send message"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Message sent"})
}

func (h *EmployerApplicationHandler) GetMessages(c *gin.Context) {
	applicationID := c.Param("applicationId")
	messages, err := h.service.GetMessages(applicationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Could not fetch messages"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "messages": messages})
}

func (h *EmployerApplicationHandler) GetApplicationsByStudent(c *gin.Context) {
	studentID := c.GetString("user_id")
	apps, err := h.service.GetApplicationsByStudent(studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to fetch applications"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "applications": apps})
}
