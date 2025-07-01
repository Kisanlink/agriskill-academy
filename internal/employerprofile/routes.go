// File: internal/employerprofile/routes.go

package employerprofile

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *EmployerProfileHandler) {
	employers := rg.Group("/employers")
	{
		employers.GET("/:employerId/profile", handler.GetProfile)
		employers.PUT("/:employerId/profile", handler.UpdateProfile)
	}
}
