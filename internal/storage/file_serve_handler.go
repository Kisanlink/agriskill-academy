package storage

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/Kisanlink/agriskill-academy/internal/middleware"
	db "github.com/Kisanlink/agriskill-academy/pkg/db"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FileServeHandler struct {
	s3 *db.S3Manager
	db *gorm.DB
}

func NewFileServeHandler(s3 *db.S3Manager, database *gorm.DB) *FileServeHandler {
	return &FileServeHandler{s3: s3, db: database}
}

// StudentProfile represents the student profile model for database lookup
type StudentProfile struct {
	ID              string `gorm:"primaryKey;type:varchar(255);default:gen_random_uuid()" json:"id"`
	UserID          string `gorm:"type:varchar(255);not null" json:"user_id"`
	ProfilePhotoKey string `json:"profile_photo_key,omitempty"`
	ResumeKey       string `json:"resume_key,omitempty"`
}

func (StudentProfile) TableName() string {
	return "student_profiles"
}

// Application represents the application model for database lookup
type Application struct {
	ID        string `gorm:"primaryKey;type:varchar(255);default:gen_random_uuid()" json:"id"`
	ResumeKey string `json:"resume_key,omitempty"`
}

func (Application) TableName() string {
	return "applications"
}

// Certificate represents the certificate model for database lookup
type Certificate struct {
	ID      string `gorm:"primaryKey;type:varchar(255);default:gen_random_uuid()" json:"id"`
	FileKey string `json:"file_key,omitempty"`
}

func (Certificate) TableName() string {
	return "certificates"
}

// EmployerProfile represents the employer profile model for database lookup
type EmployerProfile struct {
	ID      string `gorm:"primaryKey;type:varchar(255);default:gen_random_uuid()" json:"id"`
	UserID  string `gorm:"type:varchar(255);not null" json:"user_id"`
	LogoKey string `json:"logo_key,omitempty"`
}

func (EmployerProfile) TableName() string {
	return "employer_profiles"
}

// GET /files/serve/resume/:user_id
func (h *FileServeHandler) ServeResume(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	middleware.DebugLog("DEBUG: ServeResume - UserID: %s\n", userID)

	// First, get the resume key from the student profile
	var profile StudentProfile
	err := h.db.Where("user_id = ?", userID).First(&profile).Error
	if err != nil {
		middleware.DebugLog("DEBUG: Profile not found for user %s: %v\n", userID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		return
	}

	if profile.ResumeKey == "" {
		middleware.DebugLog("DEBUG: No resume key for user %s\n", userID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Resume not set"})
		return
	}

	middleware.DebugLog("DEBUG: Found resume key: %s\n", profile.ResumeKey)

	// Normalize the key to use forward slashes for S3
	normalizedKey := strings.ReplaceAll(profile.ResumeKey, "\\", "/")
	middleware.DebugLog("DEBUG: Normalized resume key: %s\n", normalizedKey)

	// Download the file from S3
	reader, err := h.s3.DownloadFile(c, normalizedKey)
	if err != nil {
		middleware.DebugLog("DEBUG: Failed to download resume from S3: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to download resume"})
		return
	}
	defer reader.Close()

	// Determine content type based on file extension
	ext := strings.ToLower(filepath.Ext(normalizedKey))
	contentType := "application/pdf" // default
	switch ext {
	case ".doc":
		contentType = "application/msword"
	case ".docx":
		contentType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	case ".pdf":
		contentType = "application/pdf"
	}

	middleware.DebugLog("DEBUG: Serving resume with content type: %s\n", contentType)

	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(normalizedKey)))
	c.Status(http.StatusOK)
	_, _ = io.Copy(c.Writer, reader)
}

// GET /files/serve/certificate/:certificate_id
func (h *FileServeHandler) ServeCertificate(c *gin.Context) {
	certificateID := c.Param("certificate_id")
	if certificateID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Certificate ID is required"})
		return
	}

	middleware.DebugLog("DEBUG: ServeCertificate - CertificateID: %s\n", certificateID)

	// First, get the file key from the certificate
	var certificate Certificate
	err := h.db.Where("id = ?", certificateID).First(&certificate).Error
	if err != nil {
		middleware.DebugLog("DEBUG: Certificate not found for ID %s: %v\n", certificateID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Certificate not found"})
		return
	}

	if certificate.FileKey == "" {
		middleware.DebugLog("DEBUG: No file key for certificate %s\n", certificateID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Certificate file not set"})
		return
	}

	middleware.DebugLog("DEBUG: Found certificate file key: %s\n", certificate.FileKey)

	// Normalize the key to use forward slashes for S3
	normalizedKey := strings.ReplaceAll(certificate.FileKey, "\\", "/")
	middleware.DebugLog("DEBUG: Normalized certificate key: %s\n", normalizedKey)

	// Download the file from S3
	reader, err := h.s3.DownloadFile(c, normalizedKey)
	if err != nil {
		middleware.DebugLog("DEBUG: Failed to download certificate from S3: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to download certificate"})
		return
	}
	defer reader.Close()

	// Determine content type based on file extension
	ext := strings.ToLower(filepath.Ext(normalizedKey))
	contentType := "application/pdf" // default
	switch ext {
	case ".doc":
		contentType = "application/msword"
	case ".docx":
		contentType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	case ".pdf":
		contentType = "application/pdf"
	case ".jpg", ".jpeg":
		contentType = "image/jpeg"
	case ".png":
		contentType = "image/png"
	}

	middleware.DebugLog("DEBUG: Serving certificate with content type: %s\n", contentType)

	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(normalizedKey)))
	c.Status(http.StatusOK)
	_, _ = io.Copy(c.Writer, reader)
}

// GET /files/serve/profile-photo/:user_id
func (h *FileServeHandler) ServeProfilePhoto(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	middleware.DebugLog("DEBUG: ServeProfilePhoto - UserID: %s\n", userID)

	// First, get the profile photo key from the student profile
	var profile StudentProfile
	err := h.db.Where("user_id = ?", userID).First(&profile).Error
	if err != nil {
		middleware.DebugLog("DEBUG: Profile not found for user %s: %v\n", userID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		return
	}

	if profile.ProfilePhotoKey == "" {
		middleware.DebugLog("DEBUG: No profile photo key for user %s\n", userID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile photo not set"})
		return
	}

	middleware.DebugLog("DEBUG: Found profile photo key: %s\n", profile.ProfilePhotoKey)

	// Normalize the key to use forward slashes for S3
	normalizedKey := strings.ReplaceAll(profile.ProfilePhotoKey, "\\", "/")
	middleware.DebugLog("DEBUG: Normalized key: %s\n", normalizedKey)

	// Download the file from S3
	reader, err := h.s3.DownloadFile(c, normalizedKey)
	if err != nil {
		middleware.DebugLog("DEBUG: Failed to download from S3: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to download profile photo"})
		return
	}
	defer reader.Close()

	// Determine content type based on file extension
	ext := strings.ToLower(filepath.Ext(normalizedKey))
	contentType := "image/jpeg" // default
	switch ext {
	case ".png":
		contentType = "image/png"
	case ".gif":
		contentType = "image/gif"
	case ".webp":
		contentType = "image/webp"
	case ".jpg", ".jpeg":
		contentType = "image/jpeg"
	}

	middleware.DebugLog("DEBUG: Serving profile photo with content type: %s\n", contentType)

	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", fmt.Sprintf("inline; filename=%s", filepath.Base(normalizedKey)))
	c.Status(http.StatusOK)
	_, _ = io.Copy(c.Writer, reader)
}

// GET /files/serve/logo/:employer_id
func (h *FileServeHandler) ServeLogo(c *gin.Context) {
	employerID := c.Param("employer_id")
	if employerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Employer ID is required"})
		return
	}

	middleware.DebugLog("DEBUG: ServeLogo - EmployerID: %s\n", employerID)

	// First, get the logo key from the employer profile
	var profile EmployerProfile
	err := h.db.Where("user_id = ?", employerID).First(&profile).Error
	if err != nil {
		middleware.DebugLog("DEBUG: Employer profile not found for user %s: %v\n", employerID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Employer profile not found"})
		return
	}

	if profile.LogoKey == "" {
		middleware.DebugLog("DEBUG: No logo key for employer %s\n", employerID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Logo not set"})
		return
	}

	middleware.DebugLog("DEBUG: Found logo key: %s\n", profile.LogoKey)

	// Normalize the key to use forward slashes for S3
	normalizedKey := strings.ReplaceAll(profile.LogoKey, "\\", "/")
	middleware.DebugLog("DEBUG: Normalized logo key: %s\n", normalizedKey)

	// Download the file from S3
	reader, err := h.s3.DownloadFile(c, normalizedKey)
	if err != nil {
		middleware.DebugLog("DEBUG: Failed to download logo from S3: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to download logo"})
		return
	}
	defer reader.Close()

	// Determine content type based on file extension
	ext := strings.ToLower(filepath.Ext(normalizedKey))
	contentType := "image/png" // default
	switch ext {
	case ".jpg", ".jpeg":
		contentType = "image/jpeg"
	case ".png":
		contentType = "image/png"
	case ".gif":
		contentType = "image/gif"
	case ".webp":
		contentType = "image/webp"
	}

	middleware.DebugLog("DEBUG: Serving logo with content type: %s\n", contentType)

	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", fmt.Sprintf("inline; filename=%s", filepath.Base(normalizedKey)))
	c.Status(http.StatusOK)
	_, _ = io.Copy(c.Writer, reader)
}

// GET /files/serve/avatar/:user_id
func (h *FileServeHandler) ServeAvatar(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}
	// Assume avatar key is avatars/{user_id}.jpg (or similar convention)
	keyPrefix := "avatars/" + userID
	var files []db.S3File
	filters := []db.Filter{
		h.s3.BuildFilter("prefix", db.FilterOpEqual, keyPrefix),
	}
	err := h.s3.List(c, filters, &files)
	if err != nil || len(files) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Avatar not found"})
		return
	}
	file := files[0]
	reader, err := h.s3.DownloadFile(c, file.Key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to download avatar"})
		return
	}
	defer reader.Close()
	c.Header("Content-Type", file.ContentType)
	c.Header("Content-Disposition", fmt.Sprintf("inline; filename=%s", filepath.Base(file.Key)))
	c.Header("Content-Length", fmt.Sprintf("%d", file.Size))
	c.Status(http.StatusOK)
	_, _ = io.Copy(c.Writer, reader)
}

// GET /files/serve/application-resume/:application_id
func (h *FileServeHandler) ServeApplicationResume(c *gin.Context) {
	applicationID := c.Param("application_id")
	if applicationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Application ID is required"})
		return
	}

	middleware.DebugLog("DEBUG: ServeApplicationResume - ApplicationID: %s\n", applicationID)

	// First, get the resume key from the application
	var application Application
	err := h.db.Where("id = ?", applicationID).First(&application).Error
	if err != nil {
		middleware.DebugLog("DEBUG: Application not found for ID %s: %v\n", applicationID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}

	if application.ResumeKey == "" {
		middleware.DebugLog("DEBUG: No resume key for application %s\n", applicationID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Application resume not set"})
		return
	}

	middleware.DebugLog("DEBUG: Found application resume key: %s\n", application.ResumeKey)

	// Normalize the key to use forward slashes for S3
	normalizedKey := strings.ReplaceAll(application.ResumeKey, "\\", "/")
	middleware.DebugLog("DEBUG: Normalized application resume key: %s\n", normalizedKey)

	// Download the file from S3
	reader, err := h.s3.DownloadFile(c, normalizedKey)
	if err != nil {
		middleware.DebugLog("DEBUG: Failed to download application resume from S3: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to download application resume"})
		return
	}
	defer reader.Close()

	// Determine content type based on file extension
	ext := strings.ToLower(filepath.Ext(normalizedKey))
	contentType := "application/pdf" // default
	switch ext {
	case ".doc":
		contentType = "application/msword"
	case ".docx":
		contentType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	case ".pdf":
		contentType = "application/pdf"
	}

	middleware.DebugLog("DEBUG: Serving application resume with content type: %s\n", contentType)

	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(normalizedKey)))
	c.Status(http.StatusOK)
	_, _ = io.Copy(c.Writer, reader)
}
