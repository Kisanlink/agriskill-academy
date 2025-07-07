package storage

import (
	"errors"
	"mime/multipart"
	"path/filepath"
	"strings"
)

// Custom errors
var (
	ErrFileTooLarge    = errors.New("file size exceeds maximum allowed size")
	ErrInvalidFileType = errors.New("file type not allowed")
	ErrFileNotFound    = errors.New("file not found")
	ErrInvalidFolder   = errors.New("invalid folder name")
)

// File upload constants
const (
	MaxFileSize  = 10 * 1024 * 1024 // 10MB
	MaxImageSize = 5 * 1024 * 1024  // 5MB for images
)

// Allowed file types
var AllowedImageTypes = []string{
	".jpg", ".jpeg", ".png", ".gif", ".webp",
}

var AllowedDocumentTypes = []string{
	".pdf", ".doc", ".docx", ".txt", ".rtf",
}

var AllowedResumeTypes = []string{
	".pdf", ".doc", ".docx",
}

// Request/Response Models
type UploadRequest struct {
	File   *multipart.FileHeader `form:"file" binding:"required"`
	Folder string                `form:"folder" binding:"required"`
}

type UploadResponse struct {
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	FilePath string `json:"filePath,omitempty"`
	FileName string `json:"fileName,omitempty"`
	FileSize int64  `json:"fileSize,omitempty"`
	FileType string `json:"fileType,omitempty"`
	FileURL  string `json:"fileUrl,omitempty"`
}

type DeleteFileRequest struct {
	FilePath string `json:"filePath" binding:"required"`
}

type FileInfo struct {
	Name     string `json:"name"`
	Size     int64  `json:"size"`
	Type     string `json:"type"`
	Path     string `json:"path"`
	URL      string `json:"url"`
	Uploaded string `json:"uploaded"`
}

type ListFilesResponse struct {
	Success bool       `json:"success"`
	Message string     `json:"message"`
	Files   []FileInfo `json:"files,omitempty"`
}

// File validation functions
func IsValidImageFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	for _, allowedExt := range AllowedImageTypes {
		if ext == allowedExt {
			return true
		}
	}
	return false
}

func IsValidDocumentFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	for _, allowedExt := range AllowedDocumentTypes {
		if ext == allowedExt {
			return true
		}
	}
	return false
}

func IsValidResumeFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	for _, allowedExt := range AllowedResumeTypes {
		if ext == allowedExt {
			return true
		}
	}
	return false
}

func GetFileType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))

	// Check if it's an image
	for _, imgExt := range AllowedImageTypes {
		if ext == imgExt {
			return "image"
		}
	}

	// Check if it's a document
	for _, docExt := range AllowedDocumentTypes {
		if ext == docExt {
			return "document"
		}
	}

	return "unknown"
}

func ValidateFileSize(file *multipart.FileHeader, maxSize int64) error {
	if file.Size > maxSize {
		return ErrFileTooLarge
	}
	return nil
}

func ValidateFileType(filename string, allowedTypes []string) error {
	ext := strings.ToLower(filepath.Ext(filename))
	for _, allowedExt := range allowedTypes {
		if ext == allowedExt {
			return nil
		}
	}
	return ErrInvalidFileType
}
