// File: internal/userprofile/handler.go

package userprofile

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserProfileHandler struct {
	service UserProfileService
}

func NewUserProfileHandler(s UserProfileService) *UserProfileHandler {
	return &UserProfileHandler{s}
}

// GET /users/:userId/profile
func (h *UserProfileHandler) GetProfile(c *gin.Context) {
	userID := c.Param("userId")
	profile, err := h.service.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Profile not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Profile fetched", "data": profile})
}

// PUT /users/:userId/profile
func (h *UserProfileHandler) UpdateProfile(c *gin.Context) {
	userID := c.Param("userId")
	var req UserProfile
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request"})
		return
	}
	req.UserID = userID
	err := h.service.UpdateProfile(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to update profile"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Profile updated", "data": req})
}
