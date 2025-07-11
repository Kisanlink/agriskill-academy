// File: internal/storage/routes.go

package storage

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *StorageHandler) {
	// File upload endpoints
	rg.POST("/upload/:folder", handler.UploadFile)
	rg.POST("/upload/image/:folder", handler.UploadImage)
	rg.POST("/upload/document/:folder", handler.UploadDocument)
	rg.POST("/upload/resume/:folder", handler.UploadResume)

	// File management endpoints - Using different path structure to avoid conflicts
	rg.GET("/files/info/*filePath", handler.GetFileInfo)
	rg.GET("/files/serve/*filePath", handler.ServeFile) // Serve/Download files
	rg.DELETE("/files/*filePath", handler.DeleteFile)
	rg.GET("/files/:folder", handler.ListFiles)
}

// RegisterPublicRoutes - Public routes (no auth required)
func RegisterPublicRoutes(rg *gin.RouterGroup, handler *StorageHandler) {
	// Public file serving routes (no auth required)
	rg.GET("/files/serve/*filePath", handler.ServeFile) // Serve/Download files
}

// RegisterAuthenticatedRoutes - Authenticated routes (require auth)
func RegisterAuthenticatedRoutes(rg *gin.RouterGroup, handler *StorageHandler) {
	// File upload endpoints (require auth)
	rg.POST("/upload/:folder", handler.UploadFile)
	rg.POST("/upload/image/:folder", handler.UploadImage)
	rg.POST("/upload/document/:folder", handler.UploadDocument)
	rg.POST("/upload/resume/:folder", handler.UploadResume)
	rg.POST("/upload/student/resume", handler.UploadStudentResume)
	rg.POST("/upload/student/certificate", handler.UploadStudentCertificate)

	// File management endpoints (require auth)
	rg.GET("/files/info/*filePath", handler.GetFileInfo)
	rg.DELETE("/files/*filePath", handler.DeleteFile)
	rg.GET("/files/:folder", handler.ListFiles)
}
