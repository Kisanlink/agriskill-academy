// File: internal/worker/routes.go

package worker

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *WorkerHandler) {
	rg.POST("/worker/job", handler.EnqueueJob)
}
