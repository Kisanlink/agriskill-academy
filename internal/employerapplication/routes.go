package employerapplication

import (
	"github.com/Kisanlink/agriskill-academy/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *EmployerApplicationHandler) {
	middleware.DebugLog("DEBUG: Registering employer application routes\n")

	// Employer-only routes
	employerRoutes := rg.Group("/employer")
	employerRoutes.Use(middleware.RequireRole("employer"))
	{
		employerRoutes.GET("/jobs/:jobId/applications", handler.GetApplicationsForJob)
		employerRoutes.GET("/jobs/:jobId/applications/debug", handler.DebugApplications)
		employerRoutes.PUT("/applications/:applicationId/status", handler.UpdateStatus)
		employerRoutes.GET("/applicants/:studentId/profile", handler.GetApplicantProfile)
		employerRoutes.POST("/applications/:applicationId/message", handler.SendMessage)
		employerRoutes.GET("/applications/:applicationId/messages", handler.GetMessages)
	}

	// Student-only routes
	studentRoutes := rg.Group("/student")
	studentRoutes.Use(middleware.RequireRole("student"))
	{
		studentRoutes.GET("/applications", handler.GetApplicationsByStudent)
		studentRoutes.GET("/applications/:applicationId/messages", handler.GetMessages)
		studentRoutes.POST("/applications/:applicationId/message", handler.SendMessage)
	}

	middleware.DebugLog("DEBUG: Employer application routes registered successfully\n")
}
