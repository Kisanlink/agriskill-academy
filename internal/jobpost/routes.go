// File: internal/jobpost/routes.go

package jobpost

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *JobPostHandler) {
	jobs := rg.Group("/jobs")
	{
		jobs.POST("", handler.Create)
		jobs.PUT("/:id", handler.Update)
		jobs.DELETE("/:id", handler.Delete)
		jobs.GET("/:id", handler.GetByID)
		jobs.GET("/my-posts", handler.GetByEmployer)
		jobs.POST("/search", handler.Search)
		// Add routes for draft, publish, save/unsave, etc as needed
	}
}
