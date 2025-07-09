package employerprofile

import (
	"asa/pkg/authz"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type EmployerProfileHandler struct {
	service EmployerProfileService
}

func NewEmployerProfileHandler(s EmployerProfileService) *EmployerProfileHandler {
	return &EmployerProfileHandler{s}
}

func getJWT(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		return authHeader[7:]
	}
	return ""
}

// GET /employers/:employerId/profile
func (h *EmployerProfileHandler) GetProfile(c *gin.Context) {
	username := c.GetString("email")
	employerID := c.Param("employerId")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckAAAPermission(username, "db_asa_employer_profiles", "read", employerID, jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	userID := c.Param("employerId")
	profile, err := h.service.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Employer profile not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": profile})
}

// GET /employers/me/profile
func (h *EmployerProfileHandler) GetMyProfile(c *gin.Context) {
	username := c.GetString("email")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckAAAPermission(username, "db_asa_employer_profiles", "read", "", jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}
	profile, err := h.service.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Profile not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": profile})
}

// PUT /employers/:employerId/profile
func (h *EmployerProfileHandler) UpdateProfile(c *gin.Context) {
	username := c.GetString("email")
	employerID := c.Param("employerId")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckAAAPermission(username, "db_asa_employer_profiles", "update", employerID, jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	userID := c.Param("employerId")
	var req EmployerProfile
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request"})
		return
	}
	req.UserID = userID
	if err := h.service.UpdateProfile(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Update failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": req})
}

// POST /employers/profile
func (h *EmployerProfileHandler) CreateProfile(c *gin.Context) {
	username := c.GetString("email")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckAAAPermission(username, "db_asa_employer_profiles", "create", "", jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	var req EmployerProfile
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request"})
		return
	}

	// ✅ Get user ID from JWT context (set by AuthMiddleware)
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "User ID not found in token"})
		return
	}

	// ✅ Inject user ID into the profile model before saving
	req.UserID = userID.(string)

	if err := h.service.CreateProfile(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to create profile"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": req})
}

// DELETE /employers/:employerId/profile
func (h *EmployerProfileHandler) DeleteProfile(c *gin.Context) {
	username := c.GetString("email")
	employerID := c.Param("employerId")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckAAAPermission(username, "db_asa_employer_profiles", "delete", employerID, jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	userID := c.Param("employerId")
	if err := h.service.DeleteProfile(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to delete profile"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Profile deleted"})
}
