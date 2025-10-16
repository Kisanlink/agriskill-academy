package auth

import (
	"github.com/gin-gonic/gin"
)

// RegisterPublicRoutes registers public auth endpoints (no authentication required)
func RegisterPublicRoutes(rg *gin.RouterGroup, handler *AuthHandler) {
	auth := rg.Group("/auth")
	{
		auth.POST("/login", handler.Login)
		auth.POST("/signup", handler.Signup)
		auth.GET("/verify-email", handler.VerifyEmail)       // Email verification endpoint
		auth.POST("/forgot-password", handler.ForgotPassword)
		auth.POST("/reset-password", handler.ResetPassword)
	}
}

// RegisterProtectedRoutes registers protected auth endpoints (authentication required)
func RegisterProtectedRoutes(rg *gin.RouterGroup, handler *AuthHandler) {
	auth := rg.Group("/auth")
	{
		auth.GET("/profile", handler.GetProfile)
		auth.PUT("/profile", handler.UpdateProfile)
	}
}

// RegisterRoutes registers all auth endpoints (legacy function for backward compatibility)
func RegisterRoutes(rg *gin.RouterGroup, handler *AuthHandler) {
	RegisterPublicRoutes(rg, handler)
	RegisterProtectedRoutes(rg, handler)
}
