// File: internal/notification/routes.go

package notification

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *NotificationHandler) {
	rg.POST("/notify/email", handler.SendEmail)
}
