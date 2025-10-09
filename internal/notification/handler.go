// File: internal/notification/handler.go

package notification

import (
	"github.com/Kisanlink/agriskill-academy/pkg/authz"
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

// @Summary Get Notification Preferences
// @Description Get notification preferences for the current user
// @Tags Notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} NotificationPreferencesResponse "Preferences retrieved successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Permission denied"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/notifications/preferences [get]
// GET /notifications/preferences
func (h *NotificationHandler) GetPreferences(c *gin.Context) {
	username := c.GetString("email")
	userID := c.GetString("user_id")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckLocalPermission(username, "db_asa_notifications", "read", userID, jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	preferences, err := h.service.GetPreferences(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to get preferences"})
		return
	}

	c.JSON(http.StatusOK, NotificationPreferencesResponse{
		Success: true,
		Message: "Preferences retrieved successfully",
		Data:    preferences,
	})
}

// @Summary Update Notification Preferences
// @Description Update notification preferences for the current user
// @Tags Notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UpdateNotificationPreferencesRequest true "Notification preferences"
// @Success 200 {object} NotificationPreferencesResponse "Preferences updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Permission denied"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/notifications/preferences [put]
// PUT /notifications/preferences
func (h *NotificationHandler) UpdatePreferences(c *gin.Context) {
	username := c.GetString("email")
	userID := c.GetString("user_id")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckLocalPermission(username, "db_asa_notifications", "update", userID, jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	var req UpdateNotificationPreferencesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request body"})
		return
	}

	preferences, err := h.service.UpdatePreferences(userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to update preferences"})
		return
	}

	c.JSON(http.StatusOK, NotificationPreferencesResponse{
		Success: true,
		Message: "Preferences updated successfully",
		Data:    preferences,
	})
}

// @Summary Send Email
// @Description Send an email notification
// @Tags Notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body map[string]interface{} true "Email data"
// @Success 200 {object} map[string]interface{} "Email sent successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Permission denied"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/notifications/email [post]
// POST /notifications/send-email
func (h *NotificationHandler) SendEmail(c *gin.Context) {
	username := c.GetString("email")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckLocalPermission(username, "db_asa_notifications", "create", "", jwtToken)
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
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request body"})
		return
	}

	err = h.service.SendEmail(req.To, req.Subject, req.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to send email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Email sent successfully"})
}
