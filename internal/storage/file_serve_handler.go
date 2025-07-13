package storage

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FileServeHandler struct {
	db *gorm.DB
}

func NewFileServeHandler(db *gorm.DB) *FileServeHandler {
	return &FileServeHandler{db: db}
}

// GET /files/serve/resume/:user_id
func (h *FileServeHandler) ServeResume(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// Get resume from student_profiles table
	var profile struct {
		Resume     []byte `gorm:"column:resume"`
		ResumeName string `gorm:"column:resume_name"`
		ResumeType string `gorm:"column:resume_type"`
		ResumeSize int64  `gorm:"column:resume_size"`
	}

	err := h.db.Table("student_profiles").
		Select("resume, resume_name, resume_type, resume_size").
		Where("user_id = ?", userID).
		First(&profile).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resume not found"})
		return
	}

	if len(profile.Resume) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resume file is empty"})
		return
	}

	// Set response headers
	c.Header("Content-Type", profile.ResumeType)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", profile.ResumeName))
	c.Header("Content-Length", fmt.Sprintf("%d", profile.ResumeSize))

	// Return binary data
	c.Data(http.StatusOK, profile.ResumeType, profile.Resume)
}

// GET /files/serve/certificate/:certificate_id
func (h *FileServeHandler) ServeCertificate(c *gin.Context) {
	certificateID := c.Param("certificate_id")
	if certificateID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Certificate ID is required"})
		return
	}

	// Get certificate from certificates table
	var certificate struct {
		File     []byte `gorm:"column:file"`
		FileName string `gorm:"column:file_name"`
		FileType string `gorm:"column:file_type"`
		FileSize int64  `gorm:"column:file_size"`
	}

	err := h.db.Table("certificates").
		Select("file, file_name, file_type, file_size").
		Where("id = ?", certificateID).
		First(&certificate).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Certificate not found"})
		return
	}

	if len(certificate.File) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Certificate file is empty"})
		return
	}

	// Set response headers
	c.Header("Content-Type", certificate.FileType)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", certificate.FileName))
	c.Header("Content-Length", fmt.Sprintf("%d", certificate.FileSize))

	// Return binary data
	c.Data(http.StatusOK, certificate.FileType, certificate.File)
}

// GET /files/serve/profile-photo/:user_id
func (h *FileServeHandler) ServeProfilePhoto(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// Get profile photo from student_profiles table
	var profile struct {
		ProfilePhoto     []byte `gorm:"column:profile_photo"`
		ProfilePhotoName string `gorm:"column:profile_photo_name"`
		ProfilePhotoType string `gorm:"column:profile_photo_type"`
		ProfilePhotoSize int64  `gorm:"column:profile_photo_size"`
	}

	err := h.db.Table("student_profiles").
		Select("profile_photo, profile_photo_name, profile_photo_type, profile_photo_size").
		Where("user_id = ?", userID).
		First(&profile).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile photo not found"})
		return
	}

	if len(profile.ProfilePhoto) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile photo is empty"})
		return
	}

	// Set response headers
	c.Header("Content-Type", profile.ProfilePhotoType)
	c.Header("Content-Disposition", fmt.Sprintf("inline; filename=%s", profile.ProfilePhotoName))
	c.Header("Content-Length", fmt.Sprintf("%d", profile.ProfilePhotoSize))

	// Return binary data
	c.Data(http.StatusOK, profile.ProfilePhotoType, profile.ProfilePhoto)
}

// GET /files/serve/logo/:employer_id
func (h *FileServeHandler) ServeLogo(c *gin.Context) {
	employerID := c.Param("employer_id")
	if employerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Employer ID is required"})
		return
	}

	// Get logo from employer_profiles table
	var profile struct {
		Logo     []byte `gorm:"column:logo"`
		LogoName string `gorm:"column:logo_name"`
		LogoType string `gorm:"column:logo_type"`
		LogoSize int64  `gorm:"column:logo_size"`
	}

	err := h.db.Table("employer_profiles").
		Select("logo, logo_name, logo_type, logo_size").
		Where("user_id = ?", employerID).
		First(&profile).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Logo not found"})
		return
	}

	if len(profile.Logo) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Logo is empty"})
		return
	}

	// Set response headers
	c.Header("Content-Type", profile.LogoType)
	c.Header("Content-Disposition", fmt.Sprintf("inline; filename=%s", profile.LogoName))
	c.Header("Content-Length", fmt.Sprintf("%d", profile.LogoSize))

	// Return binary data
	c.Data(http.StatusOK, profile.LogoType, profile.Logo)
}

// GET /files/serve/avatar/:user_id
func (h *FileServeHandler) ServeAvatar(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// Get avatar from users table
	var user struct {
		Avatar     []byte `gorm:"column:avatar"`
		AvatarName string `gorm:"column:avatar_name"`
		AvatarType string `gorm:"column:avatar_type"`
		AvatarSize int64  `gorm:"column:avatar_size"`
	}

	err := h.db.Table("users").
		Select("avatar, avatar_name, avatar_type, avatar_size").
		Where("id = ?", userID).
		First(&user).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Avatar not found"})
		return
	}

	if len(user.Avatar) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Avatar is empty"})
		return
	}

	// Set response headers
	c.Header("Content-Type", user.AvatarType)
	c.Header("Content-Disposition", fmt.Sprintf("inline; filename=%s", user.AvatarName))
	c.Header("Content-Length", fmt.Sprintf("%d", user.AvatarSize))

	// Return binary data
	c.Data(http.StatusOK, user.AvatarType, user.Avatar)
}

// GET /files/serve/application-resume/:application_id
func (h *FileServeHandler) ServeApplicationResume(c *gin.Context) {
	applicationID := c.Param("application_id")
	if applicationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Application ID is required"})
		return
	}

	// Get resume from applications table
	var application struct {
		ResumeFile     []byte `gorm:"column:resume_file"`
		ResumeFileName string `gorm:"column:resume_file_name"`
		ResumeFileType string `gorm:"column:resume_file_type"`
		ResumeFileSize int64  `gorm:"column:resume_file_size"`
	}

	err := h.db.Table("applications").
		Select("resume_file, resume_file_name, resume_file_type, resume_file_size").
		Where("id = ?", applicationID).
		First(&application).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application resume not found"})
		return
	}

	if len(application.ResumeFile) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application resume is empty"})
		return
	}

	// Set response headers
	c.Header("Content-Type", application.ResumeFileType)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", application.ResumeFileName))
	c.Header("Content-Length", fmt.Sprintf("%d", application.ResumeFileSize))

	// Return binary data
	c.Data(http.StatusOK, application.ResumeFileType, application.ResumeFile)
}
