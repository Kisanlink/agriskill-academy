// File: internal/auth/routes.go

package auth

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *AuthHandler) {
	auth := rg.Group("/auth")
	{
		auth.POST("/login", handler.Login)
		auth.POST("/signup", handler.Signup)
		// Add additional endpoints: verify, reset-password, etc
	}
}
