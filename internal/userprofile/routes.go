// File: internal/userprofile/routes.go

package userprofile

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *UserProfileHandler) {
	users := rg.Group("/users")
	{
		users.GET("/:userId/profile", handler.GetProfile)
		users.PUT("/:userId/profile", handler.UpdateProfile)
	}
}
