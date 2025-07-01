// File: internal/storage/service.go

package storage

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

type StorageService interface {
	SaveFile(fileHeader *multipart.FileHeader, folder string) (string, error)
}

type localStorageService struct {
	BasePath string
}

func NewLocalStorageService(basePath string) StorageService {
	return &localStorageService{BasePath: basePath}
}

func (s *localStorageService) SaveFile(fileHeader *multipart.FileHeader, folder string) (string, error) {
	if err := os.MkdirAll(filepath.Join(s.BasePath, folder), os.ModePerm); err != nil {
		return "", err
	}
	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), fileHeader.Filename)
	dst := filepath.Join(s.BasePath, folder, filename)
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	out, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = out.ReadFrom(file)
	if err != nil {
		return "", err
	}
	return filepath.Join(folder, filename), nil
}
