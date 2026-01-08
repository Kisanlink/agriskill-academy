package admin

import (
	"github.com/Kisanlink/agriskill-academy/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *AdminHandler) {
	admin := rg.Group("/admin")
	admin.Use(middleware.RequireRole("asa_admin"))
	{
		// Admin Management
		admin.POST("/create-admin", handler.CreateAdmin)

		// Analytics endpoints
		analytics := admin.Group("/analytics")
		{
			analytics.GET("/jobs", handler.GetJobAnalytics)
			analytics.GET("/users", handler.GetUserAnalytics)
			analytics.GET("/applications", handler.GetApplicationAnalytics)
			analytics.GET("/companies", handler.GetCompanyAnalytics)
			analytics.GET("/dashboard", handler.GetDashboardAnalytics)
		}

		// User Management
		users := admin.Group("/users")
		{
			users.GET("", handler.GetUsers)
			users.GET("/:id", handler.GetUserByID)
			users.PUT("/:id", handler.UpdateUser)
			users.DELETE("/:id", handler.DeleteUser)
		}

		// Company Management
		companies := admin.Group("/companies")
		{
			companies.GET("", handler.GetCompanies)
			companies.GET("/:id", handler.GetCompanyByID)
			companies.PUT("/:id", handler.UpdateCompany)
			companies.DELETE("/:id", handler.DeleteCompany)
		}

		// Student List
		admin.GET("/students", handler.GetStudents)

		// Employer List
		admin.GET("/employers", handler.GetEmployers)

		// Job Management (View-Only)
		jobs := admin.Group("/jobs")
		{
			jobs.GET("", handler.GetJobs)                    // List all jobs
			jobs.GET("/statistics", handler.GetJobStatistics) // Job statistics (must come before /:id)
			jobs.GET("/:id", handler.GetJobByID)             // Get job details
		}
	}
}
