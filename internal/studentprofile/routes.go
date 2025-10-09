// File: internal/studentprofile/routes.go

package studentprofile

import (
	"github.com/Kisanlink/agriskill-academy/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *StudentProfileHandler) {
	students := rg.Group("/students")
	{
		// Public routes (employers can view student profiles)
		students.GET("/:studentId/profile", handler.GetProfile)

		// Student-only routes (require student role)
		studentOnly := students.Group("")
		studentOnly.Use(middleware.RequireRole("student"))
		{
			studentOnly.PUT("/:studentId/profile", handler.UpdateProfile)
			studentOnly.GET("/me/profile", handler.GetMyProfile)
			studentOnly.PUT("/me/profile", handler.UpdateMyProfile)
			studentOnly.PUT("/me/resume", handler.UpdateMyResume)
			studentOnly.POST("/:studentId/certificates", handler.AddCertificate)
			studentOnly.POST("/me/certificates", handler.AddMyCertificate)
			studentOnly.POST("/me/resume", handler.UploadMyResume)
			studentOnly.POST("/me/certificate", handler.UploadMyCertificate)
			studentOnly.POST("/me/certificates/upload", handler.UploadCertificate)
			studentOnly.POST("/me/certificates/add", handler.AddCertificateToProfile)
			studentOnly.DELETE("/me/certificates/:certificateId", handler.DeleteMyCertificate)
		}
	}
}
