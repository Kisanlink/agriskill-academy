// File: internal/application/handler.go

package application

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ApplicationHandler struct {
	service ApplicationService
}

func NewApplicationHandler(s ApplicationService) *ApplicationHandler {
	return &ApplicationHandler{s}
}

// POST /jobs/:jobId/apply
func (h *ApplicationHandler) Apply(c *gin.Context) {
	jobId := c.Param("jobId")
	studentID := c.GetString("user_id")

	coverLetter := c.PostForm("coverLetter")
	file, err := c.FormFile("resume")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Resume file is required"})
		return
	}
	// Save file
	filename := "uploads/resumes/" + studentID + "_" + file.Filename
	if err := c.SaveUploadedFile(file, filename); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to upload resume"})
		return
	}

	app := &Application{
		JobID:       jobId,
		StudentID:   studentID,
		CoverLetter: coverLetter,
		ResumeFile:  filename,
	}

	if err := h.service.Apply(app); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Application submitted"})
}

// GET /applications/my
func (h *ApplicationHandler) GetMyApplications(c *gin.Context) {
	studentID := c.GetString("user_id") // From JWT middleware
	apps, err := h.service.GetMyApplications(studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Could not fetch applications"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "applications": apps})
}

// DELETE /applications/:applicationId
func (h *ApplicationHandler) Remove(c *gin.Context) {
	appID := c.Param("applicationId")
	studentID := c.GetString("user_id") // From JWT middleware
	if err := h.service.Remove(appID, studentID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Could not remove application"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Application removed"})
}

// PUT /applications/:applicationId/status
func (h *ApplicationHandler) UpdateStatus(c *gin.Context) {
	appID := c.Param("applicationId")
	studentID := c.GetString("user_id")
	var req struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid status"})
		return
	}
	if err := h.service.UpdateStatus(appID, studentID, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Update failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Status updated"})
}
