package auth

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *AuthHandler) {
	auth := rg.Group("/auth")
	{
		auth.POST("/login", handler.Login)
		auth.POST("/signup", handler.Signup)
		auth.POST("/forgot-password", handler.ForgotPassword)
		auth.POST("/reset-password", handler.ResetPassword)
		auth.GET("/profile", handler.GetProfile)
		auth.PUT("/profile", handler.UpdateProfile)
	}
}
