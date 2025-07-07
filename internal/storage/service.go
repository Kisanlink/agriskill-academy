// File: internal/storage/service.go

package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type StorageService interface {
	SaveFile(fileHeader *multipart.FileHeader, folder string) (string, error)
	SaveImage(fileHeader *multipart.FileHeader, folder string) (string, error)
	SaveDocument(fileHeader *multipart.FileHeader, folder string) (string, error)
	SaveResume(fileHeader *multipart.FileHeader, folder string) (string, error)
	DeleteFile(filePath string) error
	ListFiles(folder string) ([]FileInfo, error)
	GetFileInfo(filePath string) (*FileInfo, error)
}

type localStorageService struct {
	BasePath string
	BaseURL  string
}

func NewLocalStorageService(basePath, baseURL string) StorageService {
	return &localStorageService{
		BasePath: basePath,
		BaseURL:  baseURL,
	}
}

func (s *localStorageService) SaveFile(fileHeader *multipart.FileHeader, folder string) (string, error) {
	// Validate folder name
	if err := s.validateFolder(folder); err != nil {
		return "", err
	}

	// Validate file size
	if err := ValidateFileSize(fileHeader, MaxFileSize); err != nil {
		return "", err
	}

	// Validate file type
	if err := ValidateFileType(fileHeader.Filename, append(AllowedImageTypes, AllowedDocumentTypes...)); err != nil {
		return "", err
	}

	return s.saveFileInternal(fileHeader, folder)
}

func (s *localStorageService) SaveImage(fileHeader *multipart.FileHeader, folder string) (string, error) {
	// Validate folder name
	if err := s.validateFolder(folder); err != nil {
		return "", err
	}

	// Validate file size
	if err := ValidateFileSize(fileHeader, MaxImageSize); err != nil {
		return "", err
	}

	// Validate file type
	if err := ValidateFileType(fileHeader.Filename, AllowedImageTypes); err != nil {
		return "", err
	}

	return s.saveFileInternal(fileHeader, folder)
}

func (s *localStorageService) SaveDocument(fileHeader *multipart.FileHeader, folder string) (string, error) {
	// Validate folder name
	if err := s.validateFolder(folder); err != nil {
		return "", err
	}

	// Validate file size
	if err := ValidateFileSize(fileHeader, MaxFileSize); err != nil {
		return "", err
	}

	// Validate file type
	if err := ValidateFileType(fileHeader.Filename, AllowedDocumentTypes); err != nil {
		return "", err
	}

	return s.saveFileInternal(fileHeader, folder)
}

func (s *localStorageService) SaveResume(fileHeader *multipart.FileHeader, folder string) (string, error) {
	// Validate folder name
	if err := s.validateFolder(folder); err != nil {
		return "", err
	}

	// Validate file size
	if err := ValidateFileSize(fileHeader, MaxFileSize); err != nil {
		return "", err
	}

	// Validate file type
	if err := ValidateFileType(fileHeader.Filename, AllowedResumeTypes); err != nil {
		return "", err
	}

	return s.saveFileInternal(fileHeader, folder)
}

func (s *localStorageService) saveFileInternal(fileHeader *multipart.FileHeader, folder string) (string, error) {
	// Create folder if it doesn't exist
	folderPath := filepath.Join(s.BasePath, folder)
	if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create folder: %w", err)
	}

	// Generate unique filename
	timestamp := time.Now().UnixNano()
	ext := filepath.Ext(fileHeader.Filename)
	baseName := strings.TrimSuffix(fileHeader.Filename, ext)
	safeBaseName := strings.ReplaceAll(baseName, " ", "_")
	filename := fmt.Sprintf("%d_%s%s", timestamp, safeBaseName, ext)

	dst := filepath.Join(folderPath, filename)

	// Open source file
	src, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open source file: %w", err)
	}
	defer src.Close()

	// Create destination file
	out, err := os.Create(dst)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer out.Close()

	// Copy file content
	_, err = io.Copy(out, src)
	if err != nil {
		return "", fmt.Errorf("failed to copy file: %w", err)
	}

	// Return relative path
	return filepath.Join(folder, filename), nil
}

func (s *localStorageService) DeleteFile(filePath string) error {
	// Validate file path
	if strings.Contains(filePath, "..") {
		return ErrInvalidFolder
	}

	fullPath := filepath.Join(s.BasePath, filePath)

	// Check if file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return ErrFileNotFound
	}

	// Delete file
	return os.Remove(fullPath)
}

func (s *localStorageService) ListFiles(folder string) ([]FileInfo, error) {
	// Validate folder name
	if err := s.validateFolder(folder); err != nil {
		return nil, err
	}

	folderPath := filepath.Join(s.BasePath, folder)

	// Check if folder exists
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		return []FileInfo{}, nil // Return empty list if folder doesn't exist
	}

	// Read directory
	entries, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var files []FileInfo
	for _, entry := range entries {
		if !entry.IsDir() {
			filePath := filepath.Join(folder, entry.Name())
			fileInfo, err := s.GetFileInfo(filePath)
			if err == nil {
				files = append(files, *fileInfo)
			}
		}
	}

	return files, nil
}

func (s *localStorageService) GetFileInfo(filePath string) (*FileInfo, error) {
	// Validate file path
	if strings.Contains(filePath, "..") {
		return nil, ErrInvalidFolder
	}

	fullPath := filepath.Join(s.BasePath, filePath)

	// Get file info
	info, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrFileNotFound
		}
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	// Build file info
	fileInfo := &FileInfo{
		Name:     filepath.Base(filePath),
		Size:     info.Size(),
		Type:     GetFileType(filePath),
		Path:     filePath,
		URL:      s.buildFileURL(filePath),
		Uploaded: info.ModTime().Format(time.RFC3339),
	}

	return fileInfo, nil
}

func (s *localStorageService) validateFolder(folder string) error {
	// Check for path traversal attempts
	if strings.Contains(folder, "..") || strings.Contains(folder, "/") || strings.Contains(folder, "\\") {
		return ErrInvalidFolder
	}

	// Check for empty folder name
	if strings.TrimSpace(folder) == "" {
		return ErrInvalidFolder
	}

	return nil
}

func (s *localStorageService) buildFileURL(filePath string) string {
	if s.BaseURL == "" {
		return filePath
	}
	// Use the serve endpoint for file URLs
	return strings.TrimSuffix(s.BaseURL, "/") + "/serve/" + strings.TrimPrefix(filePath, "/")
}
