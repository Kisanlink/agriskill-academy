package storage

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"

	db "kisanlink-db/pkg/db"

	"github.com/gin-gonic/gin"
)

type FileServeHandler struct {
	s3 *db.S3Manager
}

func NewFileServeHandler(s3 *db.S3Manager) *FileServeHandler {
	return &FileServeHandler{s3: s3}
}

// GET /files/serve/resume/:user_id
func (h *FileServeHandler) ServeResume(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}
	// Assume resume key is resumes/{user_id}_resume.pdf (or similar convention)
	keyPrefix := "resumes/" + userID
	// List objects with this prefix and pick the first one
	var files []db.S3File
	filters := []db.Filter{
		h.s3.BuildFilter("prefix", db.FilterOpEqual, keyPrefix),
	}
	err := h.s3.List(c, filters, &files)
	if err != nil || len(files) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resume not found"})
		return
	}
	file := files[0]
	reader, err := h.s3.DownloadFile(c, file.Key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to download resume"})
		return
	}
	defer reader.Close()
	c.Header("Content-Type", file.ContentType)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(file.Key)))
	c.Header("Content-Length", fmt.Sprintf("%d", file.Size))
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
	// Assume certificate key is certificates/{certificate_id}.pdf (or similar convention)
	keyPrefix := "certificates/" + certificateID
	var files []db.S3File
	filters := []db.Filter{
		h.s3.BuildFilter("prefix", db.FilterOpEqual, keyPrefix),
	}
	err := h.s3.List(c, filters, &files)
	if err != nil || len(files) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Certificate not found"})
		return
	}
	file := files[0]
	reader, err := h.s3.DownloadFile(c, file.Key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to download certificate"})
		return
	}
	defer reader.Close()
	c.Header("Content-Type", file.ContentType)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(file.Key)))
	c.Header("Content-Length", fmt.Sprintf("%d", file.Size))
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
	// Assume profile photo key is profile_photos/{user_id}.jpg (or similar convention)
	keyPrefix := "profile_photos/" + userID
	var files []db.S3File
	filters := []db.Filter{
		h.s3.BuildFilter("prefix", db.FilterOpEqual, keyPrefix),
	}
	err := h.s3.List(c, filters, &files)
	if err != nil || len(files) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile photo not found"})
		return
	}
	file := files[0]
	reader, err := h.s3.DownloadFile(c, file.Key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to download profile photo"})
		return
	}
	defer reader.Close()
	c.Header("Content-Type", file.ContentType)
	c.Header("Content-Disposition", fmt.Sprintf("inline; filename=%s", filepath.Base(file.Key)))
	c.Header("Content-Length", fmt.Sprintf("%d", file.Size))
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
	// Assume logo key is logos/{employer_id}.png (or similar convention)
	keyPrefix := "logos/" + employerID
	var files []db.S3File
	filters := []db.Filter{
		h.s3.BuildFilter("prefix", db.FilterOpEqual, keyPrefix),
	}
	err := h.s3.List(c, filters, &files)
	if err != nil || len(files) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Logo not found"})
		return
	}
	file := files[0]
	reader, err := h.s3.DownloadFile(c, file.Key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to download logo"})
		return
	}
	defer reader.Close()
	c.Header("Content-Type", file.ContentType)
	c.Header("Content-Disposition", fmt.Sprintf("inline; filename=%s", filepath.Base(file.Key)))
	c.Header("Content-Length", fmt.Sprintf("%d", file.Size))
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
	// Assume application resume key is application_resumes/{application_id}.pdf (or similar convention)
	keyPrefix := "application_resumes/" + applicationID
	var files []db.S3File
	filters := []db.Filter{
		h.s3.BuildFilter("prefix", db.FilterOpEqual, keyPrefix),
	}
	err := h.s3.List(c, filters, &files)
	if err != nil || len(files) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application resume not found"})
		return
	}
	file := files[0]
	reader, err := h.s3.DownloadFile(c, file.Key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to download application resume"})
		return
	}
	defer reader.Close()
	c.Header("Content-Type", file.ContentType)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(file.Key)))
	c.Header("Content-Length", fmt.Sprintf("%d", file.Size))
	c.Status(http.StatusOK)
	_, _ = io.Copy(c.Writer, reader)
}
