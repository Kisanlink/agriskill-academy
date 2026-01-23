package storage

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	db "github.com/Kisanlink/agriskill-academy/pkg/db"
)

// StorageService defines the interface for file storage operations
// (re-add here for handler/service compatibility)
type StorageService interface {
	SaveFile(fileHeader *multipart.FileHeader, folder string) (string, error)
	SaveImage(fileHeader *multipart.FileHeader, folder string) (string, error)
	SaveDocument(fileHeader *multipart.FileHeader, folder string) (string, error)
	SaveResume(fileHeader *multipart.FileHeader, folder string) (string, error)
	DeleteFile(filePath string) error
	ListFiles(folder string) ([]FileInfo, error)
	GetFileInfo(filePath string) (*FileInfo, error)
	GetPresignedURL(filePath string, expiration time.Duration) (string, error)
}

// getMimeTypeFromExtension returns a best-guess MIME type for a file extension
func getMimeTypeFromExtension(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".pdf":
		return "application/pdf"
	case ".doc":
		return "application/msword"
	case ".docx":
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	case ".txt":
		return "text/plain"
	case ".rtf":
		return "application/rtf"
	default:
		return "application/octet-stream"
	}
}

type s3StorageService struct {
	s3      *db.S3Manager
	bucket  string
	baseURL string
}

func NewS3StorageService(s3 *db.S3Manager, bucket, baseURL string) StorageService {
	return &s3StorageService{s3: s3, bucket: bucket, baseURL: baseURL}
}

func (s *s3StorageService) SaveFile(fileHeader *multipart.FileHeader, folder string) (string, error) {
	return s.saveFileInternal(fileHeader, folder)
}

func (s *s3StorageService) SaveImage(fileHeader *multipart.FileHeader, folder string) (string, error) {
	return s.saveFileInternal(fileHeader, folder)
}

func (s *s3StorageService) SaveDocument(fileHeader *multipart.FileHeader, folder string) (string, error) {
	return s.saveFileInternal(fileHeader, folder)
}

func (s *s3StorageService) SaveResume(fileHeader *multipart.FileHeader, folder string) (string, error) {
	return s.saveFileInternal(fileHeader, folder)
}

func (s *s3StorageService) saveFileInternal(fileHeader *multipart.FileHeader, folder string) (string, error) {
	ctx := context.Background()
	// Validate folder name
	if strings.Contains(folder, "..") || strings.Contains(folder, "/") || strings.Contains(folder, "\\") || strings.TrimSpace(folder) == "" {
		return "", ErrInvalidFolder
	}

	src, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open source file: %w", err)
	}
	defer src.Close()

	timestamp := time.Now().UnixNano()
	ext := filepath.Ext(fileHeader.Filename)
	baseName := strings.TrimSuffix(fileHeader.Filename, ext)
	safeBaseName := strings.ReplaceAll(baseName, " ", "_")
	filename := fmt.Sprintf("%d_%s%s", timestamp, safeBaseName, ext)
	// Use forward slashes for S3 keys (not filepath.Join which uses OS-specific separators)
	key := folder + "/" + filename

	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		contentType = getMimeTypeFromExtension(fileHeader.Filename)
	}

	err = s.s3.UploadFile(ctx, key, src, contentType, nil)
	if err != nil {
		return "", err
	}

	return key, nil
}

func (s *s3StorageService) DeleteFile(filePath string) error {
	if strings.Contains(filePath, "..") {
		return ErrInvalidFolder
	}
	return s.s3.Delete(context.Background(), filePath)
}

func (s *s3StorageService) ListFiles(folder string) ([]FileInfo, error) {
	ctx := context.Background()
	var files []db.S3File
	filters := []db.Filter{
		s.s3.BuildFilter("prefix", db.FilterOpEqual, folder+"/"),
	}
	err := s.s3.List(ctx, filters, &files)
	if err != nil {
		return nil, err
	}
	var result []FileInfo
	for _, f := range files {
		result = append(result, FileInfo{
			Name:     filepath.Base(f.Key),
			Size:     f.Size,
			Type:     f.ContentType,
			Path:     f.Key,
			URL:      s.buildFileURL(f.Key),
			Uploaded: f.CreatedAt.Format(time.RFC3339),
		})
	}
	return result, nil
}

func (s *s3StorageService) GetFileInfo(filePath string) (*FileInfo, error) {
	ctx := context.Background()
	var f db.S3File
	err := s.s3.GetByKey(ctx, filePath, &f)
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			return nil, ErrFileNotFound
		}
		return nil, err
	}
	return &FileInfo{
		Name:     filepath.Base(f.Key),
		Size:     f.Size,
		Type:     f.ContentType,
		Path:     f.Key,
		URL:      s.buildFileURL(f.Key),
		Uploaded: f.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *s3StorageService) buildFileURL(filePath string) string {
	if s.baseURL == "" {
		return filePath
	}
	return strings.TrimSuffix(s.baseURL, "/") + "/serve/" + strings.TrimPrefix(filePath, "/")
}

// GetPresignedURL generates a presigned S3 URL for a file
// This is useful for email clients that cannot access API endpoints
// expiration: how long the URL should be valid (e.g., 7*24*time.Hour for 7 days)
func (s *s3StorageService) GetPresignedURL(filePath string, expiration time.Duration) (string, error) {
	ctx := context.Background()
	return s.s3.GetPresignedURL(ctx, filePath, expiration)
}

// --- Begin model types and constants (from model.go) ---
// Only add model types and helpers if not already present. Remove duplicate constants and vars.
type UploadRequest struct {
	File   *multipart.FileHeader `form:"file" binding:"required"`
	Folder string                `form:"folder" binding:"required"`
}

type UploadResponse struct {
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	FilePath string `json:"file_path,omitempty"`
	FileName string `json:"file_name,omitempty"`
	FileSize int64  `json:"file_size,omitempty"`
	FileType string `json:"file_type,omitempty"`
	FileURL  string `json:"file_url,omitempty"`
}

type DeleteFileRequest struct {
	FilePath string `json:"file_path" binding:"required"`
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

func IsValidImageFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	for _, allowedExt := range AllowedImageTypes {
		if ext == allowedExt {
			return true
		}
	}
	return false
}

var (
	ErrInvalidFolder = fmt.Errorf("invalid folder name")
	ErrFileNotFound  = fmt.Errorf("file not found")
)

var AllowedImageTypes = []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
