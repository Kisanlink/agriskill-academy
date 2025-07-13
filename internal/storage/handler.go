// File: internal/storage/handler.go

package storage

import (
	"asa/pkg/authz"
	"fmt"
	"io"
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

// @Summary Upload File
// @Description Upload a file to a specific folder
// @Tags File Storage
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param folder path string true "Folder name"
// @Param file formData file true "File to upload"
// @Success 200 {object} UploadResponse "File uploaded successfully"
// @Failure 400 {object} map[string]interface{} "Invalid file or folder"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Permission denied"
// @Router /api/upload/{folder} [post]
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

// @Summary Upload Image
// @Description Upload an image file to a specific folder
// @Tags File Storage
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param folder path string true "Folder name"
// @Param file formData file true "Image file (JPG, PNG, GIF, WebP, max 5MB)"
// @Success 200 {object} UploadResponse "Image uploaded successfully"
// @Failure 400 {object} map[string]interface{} "Invalid image format or size"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Permission denied"
// @Router /api/upload/image/{folder} [post]
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

// @Summary Upload Profile Photo
// @Description Upload a profile photo for the current user
// @Tags File Storage
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "Profile photo (JPG, PNG, GIF, WebP, max 5MB)"
// @Success 200 {object} map[string]interface{} "Profile photo processed successfully"
// @Failure 400 {object} map[string]interface{} "Invalid image format or size"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Permission denied"
// @Router /api/upload/image/profile-photo [post]
// POST /upload/image/profile-photo - Special endpoint for profile photo upload
func (h *StorageHandler) UploadProfilePhoto(c *gin.Context) {
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
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Profile photo file is required"})
		return
	}

	// Validate file type
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	allowedTypes := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	isValid := false
	for _, allowedExt := range allowedTypes {
		if ext == allowedExt {
			isValid = true
			break
		}
	}
	if !isValid {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid image format. Allowed: JPG, PNG, GIF, WebP"})
		return
	}

	// Validate file size (5MB max)
	if fileHeader.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Image size exceeds maximum allowed size (5MB)"})
		return
	}

	// Read file into byte array
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to read file"})
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to read file"})
		return
	}

	// Get file metadata
	fileName := fileHeader.Filename
	fileType := fileHeader.Header.Get("Content-Type")
	if fileType == "" {
		fileType = getMimeTypeFromExtension(fileName)
	}
	fileSize := fileHeader.Size

	// Return the binary data for the frontend to use in profile update
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Profile photo processed successfully",
		"data": gin.H{
			"file_data": fileBytes,
			"file_name": fileName,
			"file_type": fileType,
			"file_size": fileSize,
		},
	})
}

// POST /upload/document/:folder
// @Summary Upload a document
// @Description Upload a document file (PDF, DOC, DOCX, TXT, RTF) to a specified folder
// @Tags Storage
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param folder path string true "Folder name"
// @Param file formData file true "Document file"
// @Success 200 {object} UploadResponse "Document uploaded successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request or file"
// @Failure 403 {object} map[string]interface{} "Permission denied"
// @Failure 500 {object} map[string]interface{} "Document upload failed"
// @Router /api/upload/document/{folder} [post]
// @x-swagger-ui true
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
// @Summary Upload a resume
// @Description Upload a resume file (PDF, DOC, DOCX) to a specified folder
// @Tags Storage
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @x-swagger-ui true
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
// @Summary Upload Student Resume
// @Description Upload a resume file for the current student user
// @Tags Storage
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "Resume file (PDF, DOC, DOCX)"
// @Success 200 {object} UploadResponse "Resume uploaded successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request or file"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Permission denied"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/upload/student/resume [post]
// @x-swagger-ui true
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
// @Summary Upload Student Certificate
// @Description Upload a certificate file for a student (PDF, DOC, DOCX, JPG, JPEG, PNG, max 10MB)
// @Tags Storage
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "Certificate file"
// @Param name formData string true "Certificate name"
// @Param issue_date formData string true "Issue date (YYYY-MM-DD)"
// @Success 200 {object} UploadResponse "Certificate uploaded successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request or file"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Permission denied"
// @Failure 500 {object} map[string]interface{} "Certificate upload failed"
// @Router /api/upload/student/certificate [post]
// @x-swagger-ui true
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
// @Summary Delete File
// @Description Delete a file by its path
// @Tags Files
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param filePath path string true "File path"
// @Success 200 {object} map[string]interface{} "File deleted successfully"
// @Failure 403 {object} map[string]interface{} "Permission denied"
// @Failure 404 {object} map[string]interface{} "File not found"
// @Failure 400 {object} map[string]interface{} "Invalid file path"
// @Failure 500 {object} map[string]interface{} "Failed to delete file"
// @Router /api/files/{filePath} [delete]
// @x-swagger-ui true
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
// @Summary List Files in Folder
// @Description List all files in a specific folder
// @Tags Files
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param folder path string true "Folder name"
// @Success 200 {object} ListFilesResponse "Files retrieved successfully"
// @Failure 403 {object} map[string]interface{} "Permission denied"
// @Failure 400 {object} map[string]interface{} "Invalid folder name"
// @Failure 500 {object} map[string]interface{} "Failed to list files"
// @Router /api/files/{folder} [get]
// @x-swagger-ui true
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
// @Summary Get File Info
// @Description Get information about a specific file
// @Tags Files
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param filePath path string true "File path"
// @Success 200 {object} map[string]interface{} "File info retrieved successfully"
// @Failure 403 {object} map[string]interface{} "Permission denied"
// @Failure 404 {object} map[string]interface{} "File not found"
// @Failure 400 {object} map[string]interface{} "Invalid file path"
// @Failure 500 {object} map[string]interface{} "Failed to get file info"
// @Router /api/files/info/{filePath} [get]
// @x-swagger-ui true
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
// @Summary Serve/Download File
// @Description Serve or download a file from the server
// @Tags Files
// @Accept json
// @Produce octet-stream
// @Security BearerAuth
// @Param filePath path string true "File path"
// @Success 200 {file} file "File served successfully"
// @Failure 403 {object} map[string]interface{} "Permission denied"
// @Failure 404 {object} map[string]interface{} "File not found"
// @Failure 400 {object} map[string]interface{} "Invalid file path"
// @Failure 500 {object} map[string]interface{} "Failed to serve file"
// @Router /api/files/{filePath} [get]
// @x-swagger-ui true
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
