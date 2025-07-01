// File: internal/storage/routes.go

package storage

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *StorageHandler) {
	rg.POST("/upload/:folder", handler.UploadFile)
}
