// File: internal/studentprofile/handler.go

package studentprofile

import (
	"asa/pkg/authz"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type StudentProfileHandler struct {
	service StudentProfileService
}

func NewStudentProfileHandler(s StudentProfileService) *StudentProfileHandler {
	return &StudentProfileHandler{s}
}

func getJWT(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		return authHeader[7:]
	}
	return ""
}

// GET /students/:studentId/profile
func (h *StudentProfileHandler) GetProfile(c *gin.Context) {
	username := c.GetString("username")
	profileID := c.Param("studentId")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckAAAPermission(username, "db_asa_student_profile", "read", profileID, jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	profile, err := h.service.GetProfile(profileID)
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
func (h *StudentProfileHandler) UpdateProfile(c *gin.Context) {
	username := c.GetString("username")
	profileID := c.Param("studentId")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckAAAPermission(username, "db_asa_student_profile", "update", profileID, jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	userID := c.GetString("user_id")
	if profileID != userID {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "You can only update your own profile"})
		return
	}

	var req StudentProfile
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request"})
		return
	}
	req.UserID = userID
	err = h.service.UpdateProfile(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to update profile"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Profile updated", "data": req})
}

// GET /students/me/profile
func (h *StudentProfileHandler) GetMyProfile(c *gin.Context) {
	username := c.GetString("username")
	userID := c.GetString("user_id")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckAAAPermission(username, "db_asa_student_profile", "read", userID, jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

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
func (h *StudentProfileHandler) AddCertificate(c *gin.Context) {
	username := c.GetString("username")
	studentID := c.Param("studentId")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckAAAPermission(username, "db_asa_certificates", "create", studentID, jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	userID := c.GetString("user_id")
	if studentID != userID {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "You can only add certificates to your own profile"})
		return
	}

	var cert Certificate
	if err := c.ShouldBindJSON(&cert); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid certificate"})
		return
	}

	// Get the student profile first to get the profile ID
	studentProfile, err := h.service.GetProfile(studentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Student profile not found"})
		return
	}

	// Set the student profile ID from the actual profile and clear any invalid ID
	cert.StudentProfileID = studentProfile.ID
	cert.ID = "" // Clear ID to let database generate proper UUID

	err = h.service.AddCertificate(&cert)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to add certificate"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Certificate added", "data": cert})
}

// PUT /students/me/profile
func (h *StudentProfileHandler) UpdateMyProfile(c *gin.Context) {
	username := c.GetString("username")
	userID := c.GetString("user_id")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckAAAPermission(username, "db_asa_student_profile", "update", userID, jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}

	var req UpdateStudentProfileRequest
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

		// Set the student profile ID for each certificate and clear any invalid IDs
		for i := range existingProfile.Certificates {
			existingProfile.Certificates[i].StudentProfileID = existingProfile.ID
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
func (h *StudentProfileHandler) AddMyCertificate(c *gin.Context) {
	username := c.GetString("username")
	userID := c.GetString("user_id")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckAAAPermission(username, "db_asa_certificates", "create", userID, jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

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

	// Get the student profile first to get the profile ID
	studentProfile, err := h.service.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Student profile not found"})
		return
	}

	// Set the student profile ID from the actual profile and clear any invalid ID
	cert.StudentProfileID = studentProfile.ID
	cert.ID = "" // Clear ID to let database generate proper UUID

	err = h.service.AddCertificate(&cert)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to add certificate", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Certificate added successfully", "data": cert})
}
