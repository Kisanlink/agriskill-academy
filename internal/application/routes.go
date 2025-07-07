// File: internal/application/routes.go

package application

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *ApplicationHandler) {
	// Student actions
	rg.POST("/jobs/:id/apply", handler.Apply)
	rg.GET("/applications/my", handler.GetMyApplications)
	rg.GET("/applications/:applicationId", handler.GetApplicationByID)
	rg.DELETE("/applications/:applicationId", handler.Remove)
	rg.PUT("/applications/:applicationId/status", handler.UpdateStatus)

	// Employer actions
	rg.GET("/jobs/:id/applications", handler.GetApplicationsByJob)
	rg.PUT("/jobs/:id/applications/:applicationId/status", handler.UpdateStatusByEmployer)
}
