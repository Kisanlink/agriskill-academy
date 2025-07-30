package auth

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *AuthHandler) {
	auth := rg.Group("/auth")
	{
		auth.POST("/login", handler.Login)
		auth.POST("/signup", handler.Signup)
		auth.GET("/verify", handler.Verify)
		auth.POST("/forgot-password", handler.ForgotPassword)
		auth.PUT("/profile", handler.UpdateProfile)
	}
}
