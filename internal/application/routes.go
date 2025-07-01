// File: internal/application/routes.go

package application

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *ApplicationHandler) {
	// Student actions
	rg.POST("/jobs/:jobId/apply", handler.Apply)
	rg.GET("/applications/my", handler.GetMyApplications)
	rg.DELETE("/applications/:applicationId", handler.Remove)
	rg.PUT("/applications/:applicationId/status", handler.UpdateStatus)
}
