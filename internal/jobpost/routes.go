package jobpost

import (
	"github.com/gin-gonic/gin"
)

// RegisterPublicRoutes registers public job endpoints (no auth required)
func RegisterPublicRoutes(rg *gin.RouterGroup, handler *JobPostHandler) {
	jobs := rg.Group("/jobs")
	{
		// Public endpoints (no auth required)
		jobs.GET("", handler.GetAllJobs) // Get all published jobs (for students)
		jobs.GET("/featured", handler.GetFeaturedJobs)
		jobs.GET("/recent", handler.GetRecentJobs)
		jobs.GET("/trending", handler.GetTrendingJobs)
		jobs.GET("/search-filters", handler.GetSearchFilters)
		jobs.GET("/:id", handler.GetByID)
		jobs.GET("/:id/similar", handler.GetSimilarJobs)
		jobs.POST("/search", handler.Search)
		jobs.POST("/advanced-search", handler.AdvancedSearch)
		jobs.POST("/recommendations", handler.GetRecommendedJobs)
	}
}

// RegisterAuthenticatedRoutes registers authenticated job endpoints (auth required)
func RegisterAuthenticatedRoutes(rg *gin.RouterGroup, handler *JobPostHandler) {
	jobs := rg.Group("/jobs")
	{
		// Specific routes (must come before parameterized routes)
		jobs.GET("/drafts", handler.GetDrafts)
		jobs.GET("/my-posts", handler.GetByEmployer)
		jobs.POST("/draft", handler.CreateDraft)
		jobs.POST("/publish", handler.Publish)

		// Job alerts endpoints (auth required)
		jobs.POST("/alerts", handler.CreateJobAlert)
		jobs.GET("/alerts", handler.GetJobAlertsByUser)
		jobs.GET("/alerts/:id", handler.GetJobAlertByID)
		jobs.PUT("/alerts/:id", handler.UpdateJobAlert)
		jobs.DELETE("/alerts/:id", handler.DeleteJobAlert)

		// Parameterized routes (must come after specific routes)
		jobs.POST("", handler.Create)
		jobs.PUT("/:id", handler.Update)
		jobs.DELETE("/:id", handler.Delete)
		jobs.POST("/:id/publish", handler.PublishDraft)
	}
}

// RegisterRoutes registers all job endpoints (legacy function for backward compatibility)
func RegisterRoutes(rg *gin.RouterGroup, handler *JobPostHandler) {
	jobs := rg.Group("/jobs")
	{
		// Public endpoints (no auth required) - specific routes first
		jobs.GET("", handler.GetAllJobs) // Get all published jobs (for students)
		jobs.GET("/featured", handler.GetFeaturedJobs)
		jobs.GET("/recent", handler.GetRecentJobs)
		jobs.GET("/trending", handler.GetTrendingJobs)
		jobs.GET("/search-filters", handler.GetSearchFilters)
		jobs.POST("/search", handler.Search)
		jobs.POST("/advanced-search", handler.AdvancedSearch)
		jobs.POST("/recommendations", handler.GetRecommendedJobs)

		// Specific authenticated routes (must come before parameterized routes)
		jobs.GET("/drafts", handler.GetDrafts)
		jobs.GET("/my-posts", handler.GetByEmployer)
		jobs.POST("/draft", handler.CreateDraft)
		jobs.POST("/publish", handler.Publish)

		// Job alerts endpoints (auth required)
		jobs.POST("/alerts", handler.CreateJobAlert)
		jobs.GET("/alerts", handler.GetJobAlertsByUser)
		jobs.GET("/alerts/:id", handler.GetJobAlertByID)
		jobs.PUT("/alerts/:id", handler.UpdateJobAlert)
		jobs.DELETE("/alerts/:id", handler.DeleteJobAlert)

		// Parameterized routes (must come after specific routes)
		jobs.POST("", handler.Create)
		jobs.GET("/:id", handler.GetByID)
		jobs.GET("/:id/similar", handler.GetSimilarJobs)
		jobs.PUT("/:id", handler.Update)
		jobs.DELETE("/:id", handler.Delete)
		jobs.POST("/:id/publish", handler.PublishDraft)
	}
}
