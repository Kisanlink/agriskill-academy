package employerprofile

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, handler *EmployerProfileHandler) {
	employers := rg.Group("/employers")
	{
		employers.GET("/:employerId/profile", handler.GetProfile)
		employers.GET("/me/profile", handler.GetMyProfile)
		employers.PUT("/:employerId/profile", handler.UpdateProfile)
		employers.POST("/profile", handler.CreateProfile)
		employers.DELETE("/:employerId/profile", handler.DeleteProfile)
	}
}
