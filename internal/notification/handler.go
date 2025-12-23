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

// ManageByToken shows a form to manage notification preferences using a token
func (h *NotificationHandler) ManageByToken(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		c.String(http.StatusBadRequest, "Missing token")
		return
	}

	// URL decode the token in case it was encoded
	decodedToken, err := url.QueryUnescape(token)
	if err != nil {
		decodedToken = token
	} else if decodedToken != token {
		middleware.DebugLog("🔍 Token was URL encoded - decoded from: %s to: %s", token, decodedToken)
	}
	token = decodedToken

	middleware.DebugLog("🔍 Manage preferences request received - token: %s", token)

	// Get preferences by token (without disabling anything)
	prefs, err := h.service.GetPreferencesByToken(token)
	if err != nil {
		middleware.DebugLog("❌ Failed to get preferences by token: %v", err)
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusBadRequest, `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Invalid Link - AgriSkill Academy</title>
    <style>
        body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; background: linear-gradient(135deg, #1b5e20 0%%, #2e7d32 100%%); margin: 0; padding: 0; display: flex; justify-content: center; align-items: center; min-height: 100vh; }
        .container { background: white; border-radius: 12px; padding: 40px; max-width: 500px; text-align: center; box-shadow: 0 8px 24px rgba(0,0,0,0.15); }
        h1 { color: #d32f2f; margin-bottom: 20px; }
        p { color: #666; line-height: 1.6; }
    </style>
</head>
<body>
    <div class="container">
        <h1>Invalid or Expired Link</h1>
        <p>This link is invalid or has expired. Please request a new link from your email.</p>
    </div>
</body>
</html>
`)
		return
	}

	// Helper function to generate checkbox attribute
	checkboxAttr := func(enabled bool) string {
		if enabled {
			return "checked"
		}
		return ""
	}

	// Render HTML form
	html := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Manage Notification Preferences - AgriSkill Academy</title>
    <style>
        body { 
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; 
            background: linear-gradient(135deg, #1b5e20 0%%, #2e7d32 100%%); 
            margin: 0; 
            padding: 20px; 
            display: flex; 
            justify-content: center; 
            align-items: center; 
            min-height: 100vh; 
        }
        .container { 
            background: white; 
            border-radius: 12px; 
            padding: 40px; 
            max-width: 600px; 
            width: 100%%; 
            box-shadow: 0 8px 24px rgba(0,0,0,0.15); 
        }
        h1 { 
            color: #1b5e20; 
            margin-bottom: 10px; 
            font-size: 28px;
        }
        .subtitle {
            color: #666;
            margin-bottom: 30px;
            font-size: 14px;
        }
        .form-group {
            margin-bottom: 20px;
            padding: 15px;
            background: #f8f9fa;
            border-radius: 8px;
            border-left: 4px solid #1b5e20;
        }
        label {
            display: flex;
            align-items: center;
            cursor: pointer;
            font-size: 16px;
            color: #333;
            font-weight: 500;
        }
        input[type="checkbox"] {
            width: 20px;
            height: 20px;
            margin-right: 12px;
            cursor: pointer;
            accent-color: #1b5e20;
        }
        .description {
            font-size: 13px;
            color: #666;
            margin-left: 32px;
            margin-top: 5px;
        }
        button {
            background-color: #1b5e20;
            color: white;
            padding: 14px 40px;
            border: none;
            border-radius: 6px;
            font-size: 16px;
            font-weight: bold;
            cursor: pointer;
            width: 100%%;
            margin-top: 20px;
            transition: background-color 0.3s;
        }
        button:hover {
            background-color: #2e7d32;
        }
        button:active {
            background-color: #1b5e20;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Manage Notification Preferences</h1>
        <p class="subtitle">Choose which notifications you'd like to receive</p>
        <form method="POST" action="/api/notifications/manage/%s">
            
            <div class="form-group">
                <label>
                    <input type="checkbox" name="email_notifications" value="1" %s />
                    Email Notifications
                </label>
                <div class="description">Receive notifications via email</div>
            </div>
            
            <div class="form-group">
                <label>
                    <input type="checkbox" name="push_notifications" value="1" %s />
                    Push Notifications
                </label>
                <div class="description">Receive push notifications in your browser</div>
            </div>
            
            <div class="form-group">
                <label>
                    <input type="checkbox" name="job_alerts" value="1" %s />
                    Job Alerts
                </label>
                <div class="description">Get notified about new job postings that match your profile</div>
            </div>
            
            <div class="form-group">
                <label>
                    <input type="checkbox" name="application_updates" value="1" %s />
                    Application Updates
                </label>
                <div class="description">Get notified when your job application status changes</div>
            </div>
            
            <button type="submit">Save Preferences</button>
        </form>
    </div>
</body>
</html>
`, token,
		checkboxAttr(prefs.EmailNotifications),
		checkboxAttr(prefs.PushNotifications),
		checkboxAttr(prefs.JobAlerts),
		checkboxAttr(prefs.ApplicationUpdates),
	)

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, html)
}

// ManageByTokenPost processes the form submission to update preferences
func (h *NotificationHandler) ManageByTokenPost(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		c.String(http.StatusBadRequest, "Missing token")
		return
	}

	// URL decode the token in case it was encoded
	decodedToken, err := url.QueryUnescape(token)
	if err != nil {
		decodedToken = token
	} else if decodedToken != token {
		middleware.DebugLog("🔍 Token was URL encoded - decoded from: %s to: %s", token, decodedToken)
	}
	token = decodedToken

	middleware.DebugLog("🔍 Manage preferences POST request received - token: %s", token)

	// Parse checkboxes (present = true, absent = false)
	emailNotifications := c.PostForm("email_notifications") == "1"
	pushNotifications := c.PostForm("push_notifications") == "1"
	jobAlerts := c.PostForm("job_alerts") == "1"
	applicationUpdates := c.PostForm("application_updates") == "1"

	// Update preferences
	err = h.service.UpdatePreferencesByToken(token, emailNotifications, pushNotifications, jobAlerts, applicationUpdates)
	if err != nil {
		middleware.DebugLog("❌ Failed to update preferences: %v", err)
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusBadRequest, `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Error - AgriSkill Academy</title>
    <style>
        body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; background: linear-gradient(135deg, #1b5e20 0%%, #2e7d32 100%%); margin: 0; padding: 0; display: flex; justify-content: center; align-items: center; min-height: 100vh; }
        .container { background: white; border-radius: 12px; padding: 40px; max-width: 500px; text-align: center; box-shadow: 0 8px 24px rgba(0,0,0,0.15); }
        h1 { color: #d32f2f; margin-bottom: 20px; }
        p { color: #666; line-height: 1.6; }
    </style>
</head>
<body>
    <div class="container">
        <h1>Error</h1>
        <p>Failed to update preferences. The link may be invalid or expired.</p>
    </div>
</body>
</html>
`)
		return
	}

	middleware.DebugLog("✅ Preferences updated successfully via token")

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Preferences Updated - AgriSkill Academy</title>
    <style>
        body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; background: linear-gradient(135deg, #1b5e20 0%%, #2e7d32 100%%); margin: 0; padding: 0; display: flex; justify-content: center; align-items: center; min-height: 100vh; }
        .container { background: white; border-radius: 12px; padding: 40px; max-width: 500px; text-align: center; box-shadow: 0 8px 24px rgba(0,0,0,0.15); }
        h1 { color: #1b5e20; margin-bottom: 20px; }
        p { color: #666; line-height: 1.6; }
        .success-icon { width: 80px; height: 80px; margin: 0 auto 20px; background: #e8f5e9; border-radius: 50%%; display: flex; align-items: center; justify-content: center; font-size: 40px; color: #1b5e20; }
    </style>
</head>
<body>
    <div class="container">
        <div class="success-icon">✓</div>
        <h1>Preferences Updated</h1>
        <p>Your notification preferences have been successfully updated.</p>
        <p style="font-size: 14px; color: #888;">You can close this window.</p>
    </div>
</body>
</html>
`)
}