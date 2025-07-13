// File: internal/studentprofile/routes.go

package studentprofile

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *StudentProfileHandler) {
	students := rg.Group("/students")
	{
		students.GET("/:studentId/profile", handler.GetProfile)
		students.PUT("/:studentId/profile", handler.UpdateProfile)
		students.GET("/me/profile", handler.GetMyProfile)
		students.PUT("/me/profile", handler.UpdateMyProfile)
		students.PUT("/me/resume", handler.UpdateMyResume)
		students.POST("/:studentId/certificates", handler.AddCertificate)
		students.POST("/me/certificates", handler.AddMyCertificate)
		students.POST("/me/resume", handler.UploadMyResume)
		students.POST("/me/certificate", handler.UploadMyCertificate)
		students.POST("/me/certificates/upload", handler.UploadCertificate)
		students.POST("/me/certificates/add", handler.AddCertificateToProfile)
		students.DELETE("/me/certificates/:certificateId", handler.DeleteMyCertificate)
	}
}
