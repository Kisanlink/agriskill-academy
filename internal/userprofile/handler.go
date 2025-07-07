// File: internal/userprofile/handler.go

package userprofile

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserProfileHandler struct {
	service UserProfileService
}

func NewUserProfileHandler(s UserProfileService) *UserProfileHandler {
	return &UserProfileHandler{s}
}

// GET /students/:studentId/profile
func (h *UserProfileHandler) GetProfile(c *gin.Context) {
	studentID := c.Param("studentId")
	profile, err := h.service.GetProfile(studentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Profile not found"})
		return
	}

	// Rewrite resume field to always use /api/files/serve/
	if profile.Resume != "" {
		if !strings.HasPrefix(profile.Resume, "http") {
			trimmed := profile.Resume
			if strings.HasPrefix(trimmed, "/") {
				trimmed = trimmed[1:]
			}
			if strings.HasPrefix(trimmed, "uploads/") {
				trimmed = trimmed[len("uploads/"):]
			}
			profile.Resume = "http://localhost:3000/api/files/serve/" + trimmed
		}
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Profile fetched", "data": profile})
}

// PUT /students/:studentId/profile
func (h *UserProfileHandler) UpdateProfile(c *gin.Context) {
	studentID := c.Param("studentId")
	var req UserProfile
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request"})
		return
	}
	req.UserID = studentID
	err := h.service.UpdateProfile(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to update profile"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Profile updated", "data": req})
}

// GET /students/me/profile
func (h *UserProfileHandler) GetMyProfile(c *gin.Context) {
	userID := c.GetString("user_id") // from JWT
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}
	profile, err := h.service.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Profile not found"})
		return
	}

	// Rewrite resume field to always use /api/files/serve/
	if profile.Resume != "" {
		if !strings.HasPrefix(profile.Resume, "http") {
			trimmed := profile.Resume
			if strings.HasPrefix(trimmed, "/") {
				trimmed = trimmed[1:]
			}
			if strings.HasPrefix(trimmed, "uploads/") {
				trimmed = trimmed[len("uploads/"):]
			}
			profile.Resume = "http://localhost:3000/api/files/serve/" + trimmed
		}
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Profile fetched", "data": profile})
}

// POST /students/:studentId/certificates
func (h *UserProfileHandler) AddCertificate(c *gin.Context) {
	studentID := c.Param("studentId")

	var cert Certificate
	if err := c.ShouldBindJSON(&cert); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid certificate"})
		return
	}

	// Get the user profile first to get the profile ID
	userProfile, err := h.service.GetProfile(studentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "User profile not found"})
		return
	}

	// Set the user profile ID from the actual profile and clear any invalid ID
	cert.UserProfileID = userProfile.ID
	cert.ID = "" // Clear ID to let database generate proper UUID

	err = h.service.AddCertificate(&cert)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to add certificate"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Certificate added", "data": cert})
}

// PUT /students/me/profile
func (h *UserProfileHandler) UpdateMyProfile(c *gin.Context) {
	userID := c.GetString("user_id") // from JWT
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}

	var req UpdateUserProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request",
			"error":   err.Error(),
			"details": "Check that all required fields are present and properly formatted",
		})
		return
	}

	// Get existing profile
	existingProfile, err := h.service.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Profile not found"})
		return
	}

	// Update only provided fields
	if req.Name != "" {
		existingProfile.Name = req.Name
	}
	if req.Email != "" {
		existingProfile.Email = req.Email
	}
	if req.Location != "" {
		existingProfile.Location = req.Location
	}
	if req.ProfilePhoto != "" {
		existingProfile.ProfilePhoto = req.ProfilePhoto
	}
	if req.Resume != "" {
		existingProfile.Resume = req.Resume
	}
	if req.Skills != nil {
		existingProfile.Skills = req.Skills
	}

	// Handle certificates update
	if req.Certificates != nil {
		// Clear existing certificates and add new ones
		existingProfile.Certificates = req.Certificates

		// Set the user profile ID for each certificate and clear any invalid IDs
		for i := range existingProfile.Certificates {
			existingProfile.Certificates[i].UserProfileID = existingProfile.ID
			// Clear the ID to let the database generate a proper UUID
			existingProfile.Certificates[i].ID = ""
		}
	}

	// Set the user ID
	existingProfile.UserID = userID

	err = h.service.UpdateProfile(existingProfile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to update profile", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Profile updated successfully", "data": existingProfile})
}

// POST /students/me/certificates
func (h *UserProfileHandler) AddMyCertificate(c *gin.Context) {
	userID := c.GetString("user_id") // from JWT
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}

	var cert Certificate
	if err := c.ShouldBindJSON(&cert); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid certificate data",
			"error":   err.Error(),
		})
		return
	}

	// Get the user profile first to get the profile ID
	userProfile, err := h.service.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "User profile not found"})
		return
	}

	// Set the user profile ID from the actual profile and clear any invalid ID
	cert.UserProfileID = userProfile.ID
	cert.ID = "" // Clear ID to let database generate proper UUID

	err = h.service.AddCertificate(&cert)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to add certificate", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Certificate added successfully", "data": cert})
}
