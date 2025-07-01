// File: internal/storage/handler.go

package storage

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type StorageHandler struct {
	service StorageService
}

func NewStorageHandler(s StorageService) *StorageHandler {
	return &StorageHandler{s}
}

// POST /upload/:folder
func (h *StorageHandler) UploadFile(c *gin.Context) {
	folder := c.Param("folder")
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "File is required"})
		return
	}
	path, err := h.service.SaveFile(fileHeader, folder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Upload failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "filePath": path})
}
