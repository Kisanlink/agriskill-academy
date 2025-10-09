package employerprofile

import (
	"github.com/Kisanlink/agriskill-academy/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *EmployerProfileHandler) {
	employers := rg.Group("/employers")
	{
		// Public routes (anyone can view employer profiles)
		employers.GET("/:employerId/profile", handler.GetProfile)

		// Employer-only routes (require employer role)
		employerOnly := employers.Group("")
		employerOnly.Use(middleware.RequireRole("employer"))
		{
			employerOnly.GET("/me/profile", handler.GetMyProfile)
			employerOnly.PUT("/:employerId/profile", handler.UpdateProfile)
			employerOnly.PUT("/me/profile", handler.UpdateMyProfile)
			employerOnly.POST("/me/logo", handler.UploadMyLogo)
			employerOnly.PUT("/me/logo", handler.UpdateMyLogo)
			employerOnly.POST("/profile", handler.CreateProfile)
			employerOnly.DELETE("/:employerId/profile", handler.DeleteProfile)
		}
	}
}
