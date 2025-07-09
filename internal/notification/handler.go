// File: internal/notification/handler.go

package notification

import (
	"asa/pkg/authz"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	service NotificationService
}

func NewNotificationHandler(s NotificationService) *NotificationHandler {
	return &NotificationHandler{s}
}

func getJWT(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		return authHeader[7:]
	}
	return ""
}

// POST /notify/email
func (h *NotificationHandler) SendEmail(c *gin.Context) {
	// Typically, email notification sending is system/admin only, but you could check permission:
	username := c.GetString("username")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckAAAPermission(username, "db_asa_notification_preferences", "create", "", jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	var req struct {
		To      string `json:"to" binding:"required,email"`
		Subject string `json:"subject" binding:"required"`
		Body    string `json:"body" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request"})
		return
	}
	if err := h.service.SendEmail(req.To, req.Subject, req.Body); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to send email"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Email sent"})
}

// GET /notifications/preferences
func (h *NotificationHandler) GetPreferences(c *gin.Context) {
	username := c.GetString("username")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckAAAPermission(username, "db_asa_notification_preferences", "read", "", jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}

	preferences, err := h.service.GetPreferences(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to get preferences"})
		return
	}

	c.JSON(http.StatusOK, PreferencesResponse{
		Success: true,
		Message: "Preferences retrieved successfully",
		Data:    preferences,
	})
}

// PUT /notifications/preferences
func (h *NotificationHandler) UpdatePreferences(c *gin.Context) {
	username := c.GetString("username")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckAAAPermission(username, "db_asa_notification_preferences", "update", "", jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}

	var req UpdatePreferencesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request body"})
		return
	}

	preferences, err := h.service.UpdatePreferences(userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to update preferences"})
		return
	}

	c.JSON(http.StatusOK, PreferencesResponse{
		Success: true,
		Message: "Preferences updated successfully",
		Data:    preferences,
	})
}
