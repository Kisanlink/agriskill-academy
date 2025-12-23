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
	}

	// Student-only routes
	studentRoutes := rg.Group("/student")
	studentRoutes.Use(middleware.RequireRole("student"))
	{
		studentRoutes.GET("/applications", handler.GetApplicationsByStudent)
	}

	middleware.DebugLog("DEBUG: Employer application routes registered successfully\n")
}
