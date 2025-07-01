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
	studentID := c.GetString("user_id") // From JWT middleware

	var req struct {
		CoverLetter string `form:"coverLetter"`
		ResumeFile  string `form:"resume"` // In real code, use file upload handling
	}

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid application"})
		return
	}

	app := &Application{
		JobID:       jobId,
		StudentID:   studentID,
		CoverLetter: req.CoverLetter,
		ResumeFile:  req.ResumeFile,
		// Additional fields like JobTitle, Company, etc. could be set by querying job info
	}

	if err := h.service.Apply(app); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to apply"})
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
