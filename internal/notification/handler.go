// File: internal/notification/handler.go

package notification

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/Kisanlink/agriskill-academy/internal/middleware"
	"github.com/Kisanlink/agriskill-academy/pkg/authz"
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

// @Summary Unsubscribe from Email Notifications
// @Description Unsubscribe from specific email notification type using a token from email
// @Tags Notifications
// @Accept json
// @Produce json
// @Param token path string true "Unsubscribe token from email"
// @Success 200 {object} map[string]interface{} "Unsubscribed successfully"
// @Failure 400 {object} map[string]interface{} "Invalid or expired token"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/notifications/unsubscribe/{token} [get]
// GET /notifications/unsubscribe/:token
func (h *NotificationHandler) Unsubscribe(c *gin.Context) {
	token := c.Param("token")
	
	// Add detailed logging for debugging
	middleware.DebugLog("🔍 Raw token from URL param: %s", token)
	middleware.DebugLog("🔍 Token length: %d", len(token))
	
	if token == "" {
		middleware.DebugLog("❌ Token is empty")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Unsubscribe token required",
		})
		return
	}

	// URL decode the token in case it was encoded
	decodedToken, err := url.QueryUnescape(token)
	if err != nil {
		middleware.DebugLog("⚠️  Failed to URL decode token, using raw token: %v", err)
		decodedToken = token
	} else if decodedToken != token {
		middleware.DebugLog("🔍 Token was URL encoded - decoded from: %s to: %s", token, decodedToken)
	}
	token = decodedToken

	// Log full token for debugging
	middleware.DebugLog("🔔 Unsubscribe request received - full token: %s", token)

	notificationType, err := h.service.ProcessUnsubscribe(token)
	if err != nil {
		middleware.DebugLog("❌ Unsubscribe failed: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	middleware.DebugLog("✅ Successfully unsubscribed from: %s", notificationType)

	// Map type to human-readable name
	typeLabel := "Notifications"
	switch notificationType {
	case NotificationTypeJobAlert:
		typeLabel = "Job Alerts"
	case NotificationTypeApplicationUpdate:
		typeLabel = "Application Updates"
	}

	// HTML confirmation page
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Unsubscribed - AgriSkill Academy</title>
    <style>
        body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; background: linear-gradient(135deg, #1b5e20 0%%, #2e7d32 100%%); margin: 0; padding: 0; display: flex; justify-content: center; align-items: center; min-height: 100vh; }
        .container { background: white; border-radius: 12px; padding: 40px; max-width: 500px; text-align: center; box-shadow: 0 8px 24px rgba(0,0,0,0.15); }
        h1 { color: #1b5e20; margin-bottom: 20px; }
        p { color: #666; line-height: 1.6; }
        .success-icon { width: 80px; height: 80px; margin: 0 auto 20px; background: #e8f5e9; border-radius: 50%%; display: flex; align-items: center; justify-content: center; font-size: 40px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="success-icon">✓</div>
        <h1>Successfully Unsubscribed</h1>
        <p>You have been unsubscribed from <strong>%s</strong>.</p>
        <p>You can manage your preferences anytime by logging into your account.</p>
    </div>
</body>
</html>
	`, typeLabel))
}