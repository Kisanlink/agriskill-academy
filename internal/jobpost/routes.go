package jobpost

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *JobPostHandler) {
	jobs := rg.Group("/jobs")
	{
		jobs.POST("", handler.Create)
		jobs.POST("/draft", handler.CreateDraft) // new route for draft
		jobs.POST("/publish", handler.Publish)   // new route for publish
		jobs.PUT("/:id", handler.Update)
		jobs.DELETE("/:id", handler.Delete)
		jobs.GET("/:id", handler.GetByID)
		jobs.GET("/my-posts", handler.GetByEmployer)
		jobs.POST("/search", handler.Search)
	}
}
