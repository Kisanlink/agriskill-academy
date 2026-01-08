package jobpost

import (
	"github.com/Kisanlink/agriskill-academy/internal/middleware"

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
		// Employer-only routes (require employer role)
		employerJobs := jobs.Group("")
		employerJobs.Use(middleware.RequireRole("employer"))
		{
			employerJobs.GET("/drafts", handler.GetDrafts)
			employerJobs.GET("/my-posts", handler.GetByEmployer)
			employerJobs.POST("/draft", handler.CreateDraft)
			employerJobs.POST("/publish", handler.Publish)
			employerJobs.POST("", handler.Create)
			employerJobs.PUT("/:id", handler.Update)
			employerJobs.DELETE("/:id", handler.Delete)
			employerJobs.POST("/:id/publish", handler.PublishDraft)
			employerJobs.POST("/:id/close", handler.CloseJob)
			employerJobs.POST("/:id/reopen", handler.ReopenJob)
		}

		// Job alerts endpoints (auth required for any role)
		jobs.POST("/alerts", handler.CreateJobAlert)
		jobs.GET("/alerts", handler.GetJobAlertsByUser)
		jobs.GET("/alerts/:id", handler.GetJobAlertByID)
		jobs.PUT("/alerts/:id", handler.UpdateJobAlert)
		jobs.DELETE("/alerts/:id", handler.DeleteJobAlert)
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
		jobs.POST("/:id/close", handler.CloseJob)
		jobs.POST("/:id/reopen", handler.ReopenJob)
	}
}
