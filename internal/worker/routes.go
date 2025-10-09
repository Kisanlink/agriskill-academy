// File: internal/worker/routes.go

package worker

import (
	"github.com/Kisanlink/agriskill-academy/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *WorkerHandler) {
	// Worker job management is admin-only functionality
	worker := rg.Group("/worker")
	worker.Use(middleware.RequireRole("asa_admin"))
	{
		worker.POST("/job", handler.EnqueueJob)
	}
}
