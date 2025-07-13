// File: internal/storage/routes.go

package storage

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *StorageHandler, fileServeHandler *FileServeHandler) {
	// File upload endpoints
	rg.POST("/upload/:folder", handler.UploadFile)
	rg.POST("/upload/image/:folder", handler.UploadImage)
	rg.POST("/upload/image/profile-photo", handler.UploadProfilePhoto)
	rg.POST("/upload/document/:folder", handler.UploadDocument)
	rg.POST("/upload/resume/:folder", handler.UploadResume)

	// File serving endpoints (serve files from database)
	rg.GET("/files/serve/resume/:user_id", fileServeHandler.ServeResume)
	rg.GET("/files/serve/certificate/:certificate_id", fileServeHandler.ServeCertificate)
	rg.GET("/files/serve/profile-photo/:user_id", fileServeHandler.ServeProfilePhoto)
	rg.GET("/files/serve/logo/:employer_id", fileServeHandler.ServeLogo)
	rg.GET("/files/serve/avatar/:user_id", fileServeHandler.ServeAvatar)
	rg.GET("/files/serve/application-resume/:application_id", fileServeHandler.ServeApplicationResume)

	// File management endpoints - Using different path structure to avoid conflicts
	rg.GET("/files/info/*filePath", handler.GetFileInfo)
	rg.DELETE("/files/*filePath", handler.DeleteFile)
	rg.GET("/files/:folder", handler.ListFiles)
}

// RegisterPublicRoutes - Public routes (no auth required)
func RegisterPublicRoutes(rg *gin.RouterGroup, handler *StorageHandler, fileServeHandler *FileServeHandler) {
	// Public file serving routes (no auth required)
	rg.GET("/files/serve/resume/:user_id", fileServeHandler.ServeResume)
	rg.GET("/files/serve/certificate/:certificate_id", fileServeHandler.ServeCertificate)
	rg.GET("/files/serve/profile-photo/:user_id", fileServeHandler.ServeProfilePhoto)
	rg.GET("/files/serve/logo/:employer_id", fileServeHandler.ServeLogo)
	rg.GET("/files/serve/avatar/:user_id", fileServeHandler.ServeAvatar)
	rg.GET("/files/serve/application-resume/:application_id", fileServeHandler.ServeApplicationResume)
}

// RegisterAuthenticatedRoutes - Authenticated routes (require auth)
func RegisterAuthenticatedRoutes(rg *gin.RouterGroup, handler *StorageHandler, fileServeHandler *FileServeHandler) {
	// File upload endpoints (require auth)
	rg.POST("/upload/:folder", handler.UploadFile)
	rg.POST("/upload/image/:folder", handler.UploadImage)
	rg.POST("/upload/image/profile-photo", handler.UploadProfilePhoto)
	rg.POST("/upload/document/:folder", handler.UploadDocument)
	rg.POST("/upload/resume/:folder", handler.UploadResume)
	rg.POST("/upload/student/resume", handler.UploadStudentResume)
	rg.POST("/upload/student/certificate", handler.UploadStudentCertificate)

	// File management endpoints (require auth)
	rg.GET("/files/info/*filePath", handler.GetFileInfo)
	rg.DELETE("/files/*filePath", handler.DeleteFile)
	rg.GET("/files/:folder", handler.ListFiles)
}
