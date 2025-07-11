// File: internal/storage/handler.go

package storage

import (
	"asa/pkg/authz"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type StorageHandler struct {
	service StorageService
}

func NewStorageHandler(s StorageService) *StorageHandler {
	return &StorageHandler{s}
}

func getJWT(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		return authHeader[7:]
	}
	return ""
}

// POST /upload/:folder
func (h *StorageHandler) UploadFile(c *gin.Context) {
	username := c.GetString("email")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckAAAPermission(username, "db_asa_files", "create", "", jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	folder := c.Param("folder")
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "File is required"})
		return
	}

	path, err := h.service.SaveFile(fileHeader, folder)
	if err != nil {
		switch err {
		case ErrFileTooLarge:
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "File size exceeds maximum allowed size"})
		case ErrInvalidFileType:
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "File type not allowed"})
		case ErrInvalidFolder:
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid folder name"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Upload failed: " + err.Error()})
		}
		return
	}

	// Get file info for response
	fileInfo, err := h.service.GetFileInfo(path)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": true, "filePath": path})
		return
	}

	c.JSON(http.StatusOK, UploadResponse{
		Success:  true,
		Message:  "File uploaded successfully",
		FilePath: path,
		FileName: fileInfo.Name,
		FileSize: fileInfo.Size,
		FileType: fileInfo.Type,
		FileURL:  fileInfo.URL,
	})
}

// POST /upload/image/:folder
func (h *StorageHandler) UploadImage(c *gin.Context) {
	username := c.GetString("email")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckAAAPermission(username, "db_asa_files", "create", "", jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	folder := c.Param("folder")
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Image file is required"})
		return
	}

	path, err := h.service.SaveImage(fileHeader, folder)
	if err != nil {
		switch err {
		case ErrFileTooLarge:
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Image size exceeds maximum allowed size (5MB)"})
		case ErrInvalidFileType:
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid image format. Allowed: JPG, PNG, GIF, WebP"})
		case ErrInvalidFolder:
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid folder name"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Image upload failed: " + err.Error()})
		}
		return
	}

	// Get file info for response
	fileInfo, err := h.service.GetFileInfo(path)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": true, "filePath": path})
		return
	}

	c.JSON(http.StatusOK, UploadResponse{
		Success:  true,
		Message:  "Image uploaded successfully",
		FilePath: path,
		FileName: fileInfo.Name,
		FileSize: fileInfo.Size,
		FileType: fileInfo.Type,
		FileURL:  fileInfo.URL,
	})
}

// POST /upload/document/:folder
func (h *StorageHandler) UploadDocument(c *gin.Context) {
	username := c.GetString("email")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckAAAPermission(username, "db_asa_files", "create", "", jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	folder := c.Param("folder")
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Document file is required"})
		return
	}

	path, err := h.service.SaveDocument(fileHeader, folder)
	if err != nil {
		switch err {
		case ErrFileTooLarge:
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Document size exceeds maximum allowed size (10MB)"})
		case ErrInvalidFileType:
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid document format. Allowed: PDF, DOC, DOCX, TXT, RTF"})
		case ErrInvalidFolder:
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid folder name"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Document upload failed: " + err.Error()})
		}
		return
	}

	// Get file info for response
	fileInfo, err := h.service.GetFileInfo(path)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": true, "filePath": path})
		return
	}

	c.JSON(http.StatusOK, UploadResponse{
		Success:  true,
		Message:  "Document uploaded successfully",
		FilePath: path,
		FileName: fileInfo.Name,
		FileSize: fileInfo.Size,
		FileType: fileInfo.Type,
		FileURL:  fileInfo.URL,
	})
}

// POST /upload/resume/:folder
func (h *StorageHandler) UploadResume(c *gin.Context) {
	username := c.GetString("email")
	userID := c.GetString("user_id")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckAAAPermission(username, "db_asa_files", "create", "", jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	folder := c.Param("folder")
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Resume file is required"})
		return
	}

	path, err := h.service.SaveResume(fileHeader, folder)
	if err != nil {
		switch err {
		case ErrFileTooLarge:
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Resume size exceeds maximum allowed size (10MB)"})
		case ErrInvalidFileType:
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid resume format. Allowed: PDF, DOC, DOCX"})
		case ErrInvalidFolder:
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid folder name"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Resume upload failed: " + err.Error()})
		}
		return
	}

	// If this is a student uploading a resume to the resumes folder, also update their profile
	if userID != "" && folder == "resumes" {
		fmt.Printf("DEBUG: Resume uploaded for user %s, path: %s - Profile update needed\n", userID, path)
		// Note: In a production environment, you would inject the student profile service here
		// For now, we'll just log that the profile needs to be updated
		// The frontend should make a separate call to update the profile with this file path
	}

	// Get file info for response
	fileInfo, err := h.service.GetFileInfo(path)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": true, "filePath": path})
		return
	}

	c.JSON(http.StatusOK, UploadResponse{
		Success:  true,
		Message:  "Resume uploaded successfully",
		FilePath: path,
		FileName: fileInfo.Name,
		FileSize: fileInfo.Size,
		FileType: fileInfo.Type,
		FileURL:  fileInfo.URL,
	})
}

// POST /upload/student/resume - Special endpoint for student resume upload
func (h *StorageHandler) UploadStudentResume(c *gin.Context) {
	username := c.GetString("email")
	userID := c.GetString("user_id")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckAAAPermission(username, "db_asa_files", "create", "", jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Resume file is required"})
		return
	}

	// Validate file type
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	allowedTypes := []string{".pdf", ".doc", ".docx"}
	isValid := false
	for _, allowedExt := range allowedTypes {
		if ext == allowedExt {
			isValid = true
			break
		}
	}
	if !isValid {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid file type. Allowed: PDF, DOC, DOCX"})
		return
	}

	// Validate file size (10MB max)
	if fileHeader.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "File size exceeds maximum allowed size (10MB)"})
		return
	}

	// Save file using the storage service
	path, err := h.service.SaveResume(fileHeader, "resumes")
	if err != nil {
		switch err {
		case ErrFileTooLarge:
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Resume size exceeds maximum allowed size (10MB)"})
		case ErrInvalidFileType:
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid resume format. Allowed: PDF, DOC, DOCX"})
		case ErrInvalidFolder:
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid folder name"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Resume upload failed: " + err.Error()})
		}
		return
	}

	// Get file info for response
	fileInfo, err := h.service.GetFileInfo(path)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": true, "filePath": path})
		return
	}

	fmt.Printf("DEBUG: Student resume uploaded successfully - UserID: %s, Path: %s\n", userID, path)

	c.JSON(http.StatusOK, UploadResponse{
		Success:  true,
		Message:  "Resume uploaded successfully",
		FilePath: path,
		FileName: fileInfo.Name,
		FileSize: fileInfo.Size,
		FileType: fileInfo.Type,
		FileURL:  fileInfo.URL,
	})
}

// POST /upload/student/certificate
func (h *StorageHandler) UploadStudentCertificate(c *gin.Context) {
	username := c.GetString("email")
	userID := c.GetString("user_id")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckAAAPermission(username, "db_asa_files", "create", "", jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Certificate file is required"})
		return
	}

	// Get certificate details from form
	certificateName := c.PostForm("name")
	issueDate := c.PostForm("issue_date")

	if certificateName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Certificate name is required"})
		return
	}

	if issueDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Issue date is required"})
		return
	}

	// Validate file type
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	allowedTypes := []string{".pdf", ".doc", ".docx", ".jpg", ".jpeg", ".png"}
	isValid := false
	for _, allowedExt := range allowedTypes {
		if ext == allowedExt {
			isValid = true
			break
		}
	}
	if !isValid {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid file type. Allowed: PDF, DOC, DOCX, JPG, JPEG, PNG"})
		return
	}

	// Validate file size (10MB max)
	if fileHeader.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "File size exceeds maximum allowed size (10MB)"})
		return
	}

	// Save file using the storage service
	path, err := h.service.SaveDocument(fileHeader, "certificates")
	if err != nil {
		switch err {
		case ErrFileTooLarge:
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Certificate size exceeds maximum allowed size (10MB)"})
		case ErrInvalidFileType:
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid certificate format. Allowed: PDF, DOC, DOCX, JPG, JPEG, PNG"})
		case ErrInvalidFolder:
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid folder name"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Certificate upload failed: " + err.Error()})
		}
		return
	}

	// Get file info for response
	fileInfo, err := h.service.GetFileInfo(path)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": true, "filePath": path})
		return
	}

	fmt.Printf("DEBUG: Student certificate uploaded successfully - UserID: %s, Path: %s, Name: %s\n", userID, path, certificateName)

	c.JSON(http.StatusOK, UploadResponse{
		Success:  true,
		Message:  "Certificate uploaded successfully",
		FilePath: path,
		FileName: fileInfo.Name,
		FileSize: fileInfo.Size,
		FileType: fileInfo.Type,
		FileURL:  fileInfo.URL,
	})
}

// DELETE /files/:filePath
func (h *StorageHandler) DeleteFile(c *gin.Context) {
	username := c.GetString("email")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckAAAPermission(username, "db_asa_files", "delete", "", jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	filePath := c.Param("filePath")

	// Decode URL-encoded file path
	filePath = filepath.Clean(filePath)

	err = h.service.DeleteFile(filePath)
	if err != nil {
		switch err {
		case ErrFileNotFound:
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "File not found"})
		case ErrInvalidFolder:
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid file path"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to delete file: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "File deleted successfully"})
}

// GET /files/:folder
func (h *StorageHandler) ListFiles(c *gin.Context) {
	username := c.GetString("email")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckAAAPermission(username, "db_asa_files", "read", "", jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	folder := c.Param("folder")

	files, err := h.service.ListFiles(folder)
	if err != nil {
		switch err {
		case ErrInvalidFolder:
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid folder name"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to list files: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, ListFilesResponse{
		Success: true,
		Message: "Files retrieved successfully",
		Files:   files,
	})
}

// GET /files/info/:filePath
func (h *StorageHandler) GetFileInfo(c *gin.Context) {
	username := c.GetString("email")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckAAAPermission(username, "db_asa_files", "read", "", jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	filePath := c.Param("filePath")

	// Decode URL-encoded file path
	filePath = filepath.Clean(filePath)

	fileInfo, err := h.service.GetFileInfo(filePath)
	if err != nil {
		switch err {
		case ErrFileNotFound:
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "File not found"})
		case ErrInvalidFolder:
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid file path"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to get file info: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "File info retrieved successfully",
		"file":    fileInfo,
	})
}

// GET /files/*filePath - Serve/Download file
func (h *StorageHandler) ServeFile(c *gin.Context) {
	username := c.GetString("email")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckAAAPermission(username, "db_asa_files", "read", "", jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	filePath := c.Param("filePath")

	// Decode URL-encoded file path
	filePath = filepath.Clean(filePath)

	// Validate file path
	if strings.Contains(filePath, "..") {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid file path"})
		return
	}

	// Get file info to check if it exists
	_, err = h.service.GetFileInfo(filePath)
	if err != nil {
		switch err {
		case ErrFileNotFound:
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "File not found"})
		case ErrInvalidFolder:
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid file path"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to serve file: " + err.Error()})
		}
		return
	}

	// Serve the file
	fullPath := filepath.Join("uploads", filePath)
	c.File(fullPath)
}
