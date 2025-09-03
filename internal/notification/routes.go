// File: internal/notification/routes.go

package notification

import (
	"asa/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *NotificationHandler) {
	notifications := rg.Group("/notifications")
	{
		// Admin-only route for sending emails
		notifications.POST("/email", middleware.RequireRole("asa_admin"), handler.SendEmail)
		// User preferences - any authenticated user can manage their preferences
		notifications.GET("/preferences", handler.GetPreferences)
		notifications.PUT("/preferences", handler.UpdatePreferences)
	}
}
