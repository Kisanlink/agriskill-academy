// File: internal/notification/routes.go

package notification

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *NotificationHandler) {
	notifications := rg.Group("/notifications")
	{
		notifications.POST("/email", handler.SendEmail)
		notifications.GET("/preferences", handler.GetPreferences)
		notifications.PUT("/preferences", handler.UpdatePreferences)
	}
}
