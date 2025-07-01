// File: internal/employerprofile/handler.go

package employerprofile

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type EmployerProfileHandler struct {
	service EmployerProfileService
}

func NewEmployerProfileHandler(s EmployerProfileService) *EmployerProfileHandler {
	return &EmployerProfileHandler{s}
}

// GET /employers/:employerId/profile
func (h *EmployerProfileHandler) GetProfile(c *gin.Context) {
	userID := c.Param("employerId")
	profile, err := h.service.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Employer profile not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Profile fetched",
		"data":    profile,
	})
}

// PUT /employers/:employerId/profile
func (h *EmployerProfileHandler) UpdateProfile(c *gin.Context) {
	userID := c.Param("employerId")
	var req EmployerProfile
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request"})
		return
	}
	req.UserID = userID
	err := h.service.UpdateProfile(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Update failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Profile updated", "data": req})
}
