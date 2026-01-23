// File: internal/notification/routes.go

package notification

import (
	"github.com/Kisanlink/agriskill-academy/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterPublicRoutes registers public notification routes (no auth required)
func RegisterPublicRoutes(rg *gin.RouterGroup, handler *NotificationHandler) {
	notifications := rg.Group("/notifications")
	{
		// Public unsubscribe endpoint (no auth required)
		notifications.GET("/unsubscribe/:token", handler.Unsubscribe)
		// Public manage preferences endpoints (no auth required)
		notifications.GET("/manage/:token", handler.ManageByToken)
		notifications.POST("/manage/:token", handler.ManageByTokenPost)
	}
}

// RegisterRoutes registers authenticated notification routes
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
