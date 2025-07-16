package employerprofile

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, handler *EmployerProfileHandler) {
	employers := rg.Group("/employers")
	{
		employers.GET("/:employerId/profile", handler.GetProfile)
		employers.GET("/me/profile", handler.GetMyProfile)
		employers.PUT("/:employerId/profile", handler.UpdateProfile)
		employers.PUT("/me/profile", handler.UpdateMyProfile)
		employers.POST("/profile", handler.CreateProfile)
		employers.DELETE("/:employerId/profile", handler.DeleteProfile)
	}
}
