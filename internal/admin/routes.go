package admin

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *AdminHandler) {
	admin := rg.Group("/admin")
	{
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
	}
}
