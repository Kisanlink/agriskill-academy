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
		students.POST("/:studentId/certificates", handler.AddCertificate)
		students.POST("/me/certificates", handler.AddMyCertificate)
	}
}
