package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"

	"gorm.io/gorm"
)

type BinaryFileService interface {
	SaveFileToDB(fileHeader *multipart.FileHeader) ([]byte, string, string, int64, error)
	GetFileFromDB(db *gorm.DB, tableName, id, fileColumn string) ([]byte, string, string, int64, error)
	DeleteFileFromDB(db *gorm.DB, tableName, id, fileColumn string) error
}

type binaryFileService struct{}

func NewBinaryFileService() BinaryFileService {
	return &binaryFileService{}
}

func (s *binaryFileService) SaveFileToDB(fileHeader *multipart.FileHeader) ([]byte, string, string, int64, error) {
	// Read file into byte array
	file, err := fileHeader.Open()
	if err != nil {
		return nil, "", "", 0, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, "", "", 0, fmt.Errorf("failed to read file: %w", err)
	}

	// Get file metadata
	fileName := fileHeader.Filename
	fileType := fileHeader.Header.Get("Content-Type")
	if fileType == "" {
		fileType = getMimeTypeFromExtension(fileName)
	}
	fileSize := fileHeader.Size

	return fileBytes, fileName, fileType, fileSize, nil
}

func (s *binaryFileService) GetFileFromDB(db *gorm.DB, tableName, id, fileColumn string) ([]byte, string, string, int64, error) {
	var fileData []byte
	var fileName, fileType string
	var fileSize int64

	// Query the file data and metadata
	query := fmt.Sprintf(`
		SELECT %s, %s_name, %s_type, %s_size 
		FROM %s 
		WHERE id = $1
	`, fileColumn, fileColumn, fileColumn, fileColumn, tableName)

	err := db.Raw(query, id).Scan(&struct {
		FileData []byte `gorm:"column:file_data"`
		FileName string `gorm:"column:file_name"`
		FileType string `gorm:"column:file_type"`
		FileSize int64  `gorm:"column:file_size"`
	}{
		FileData: fileData,
		FileName: fileName,
		FileType: fileType,
		FileSize: fileSize,
	}).Error

	if err != nil {
		return nil, "", "", 0, fmt.Errorf("failed to get file from database: %w", err)
	}

	if len(fileData) == 0 {
		return nil, "", "", 0, fmt.Errorf("file not found")
	}

	return fileData, fileName, fileType, fileSize, nil
}

func (s *binaryFileService) DeleteFileFromDB(db *gorm.DB, tableName, id, fileColumn string) error {
	// Update the file columns to NULL
	query := fmt.Sprintf(`
		UPDATE %s 
		SET %s = NULL, %s_name = NULL, %s_type = NULL, %s_size = NULL 
		WHERE id = $1
	`, tableName, fileColumn, fileColumn, fileColumn, fileColumn)

	result := db.Exec(query, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete file from database: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no record found to delete")
	}

	return nil
}

// Helper function to get MIME type from file extension
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

// Helper function to validate file size
func ValidateBinaryFileSize(fileSize int64, maxSize int64) error {
	if fileSize > maxSize {
		return ErrFileTooLarge
	}
	return nil
}

// Helper function to validate file type
func ValidateBinaryFileType(filename string, allowedTypes []string) error {
	ext := strings.ToLower(filepath.Ext(filename))
	for _, allowedExt := range allowedTypes {
		if ext == allowedExt {
			return nil
		}
	}
	return ErrInvalidFileType
}
