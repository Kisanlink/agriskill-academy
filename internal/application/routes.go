// File: internal/application/routes.go

package application

import (
	"github.com/Kisanlink/agriskill-academy/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *ApplicationHandler) {
	// Student applies to a job
	rg.POST("/jobs/:id/apply", middleware.RequireRole("student"), handler.Apply)

	// Student views their own applications
	rg.GET("/applications/my", middleware.RequireRole("student"), handler.GetMyApplications)

	// Employer views applications for a job
	rg.GET("/jobs/:id/applications", middleware.RequireRole("employer"), handler.GetApplicationsByJob)

	// Student or employer views a specific application
	rg.GET("/applications/:applicationId", handler.GetApplicationByID)

	// Student removes their application
	rg.DELETE("/applications/:applicationId", middleware.RequireRole("student"), handler.Remove)

	// Student updates status (withdraw)
	rg.PUT("/applications/:applicationId/status", middleware.RequireRole("student"), handler.UpdateStatus)

	// Employer updates application status
	rg.PUT("/jobs/:id/applications/:applicationId/status", middleware.RequireRole("employer"), handler.UpdateStatusByEmployer)
}
