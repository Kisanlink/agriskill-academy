// File: cmd/server/main.go

package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"asa/config"
	"asa/internal/admin"
	"asa/internal/application"
	"asa/internal/auth"
	"asa/internal/bookmark"
	"asa/internal/employerapplication"
	"asa/internal/employerprofile"
	"asa/internal/jobpost"
	"asa/internal/middleware"
	"asa/internal/notification"
	"asa/internal/storage"
	"asa/internal/studentprofile"
	"asa/internal/worker"

	"mime/multipart"

	_ "asa/docs" // Import swagger docs

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

// runAutoMigrate runs GORM AutoMigrate for all models
func runAutoMigrate(db *gorm.DB) error {
	log.Println("Running GORM AutoMigrate...")

	// Enable UUID extension for generating UUIDs if needed
	err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error
	if err != nil {
		log.Printf("Warning: Could not create uuid-ossp extension: %v", err)
	}

	// List of all models to migrate
	models := []interface{}{
		&auth.User{},
		&employerprofile.EmployerProfile{},
		&studentprofile.StudentProfile{},
		&studentprofile.Certificate{},
		&jobpost.JobPost{},
		&jobpost.JobAlert{},
		&application.Application{},
		&bookmark.Bookmark{},
		&notification.NotificationPreferences{},
		&employerapplication.Message{},
	}

	// Run AutoMigrate for each model
	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			return fmt.Errorf("failed to migrate model %T: %w", model, err)
		}
		log.Printf("Successfully migrated: %T", model)
	}

	log.Println("AutoMigrate completed successfully!")
	return nil
}

// Mock storage service for compatibility with existing storage handler
type mockStorageService struct{}

func (m *mockStorageService) SaveFile(fileHeader *multipart.FileHeader, folder string) (string, error) {
	return "", fmt.Errorf("file uploads are now handled directly in profile handlers")
}

func (m *mockStorageService) SaveImage(fileHeader *multipart.FileHeader, folder string) (string, error) {
	return "", fmt.Errorf("file uploads are now handled directly in profile handlers")
}

func (m *mockStorageService) SaveDocument(fileHeader *multipart.FileHeader, folder string) (string, error) {
	return "", fmt.Errorf("file uploads are now handled directly in profile handlers")
}

func (m *mockStorageService) SaveResume(fileHeader *multipart.FileHeader, folder string) (string, error) {
	return "", fmt.Errorf("file uploads are now handled directly in profile handlers")
}

func (m *mockStorageService) DeleteFile(filePath string) error {
	return fmt.Errorf("file operations are now handled directly in profile handlers")
}

func (m *mockStorageService) ListFiles(folder string) ([]storage.FileInfo, error) {
	return []storage.FileInfo{}, nil
}

func (m *mockStorageService) GetFileInfo(filePath string) (*storage.FileInfo, error) {
	return nil, fmt.Errorf("file operations are now handled directly in profile handlers")
}

func main() {
	// Initialize config and DB
	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize DB: %v", err)
	}
	defer config.CloseDB(db)

	// Run AutoMigrate to create/update all tables
	if err := runAutoMigrate(db); err != nil {
		log.Fatalf("Failed to run auto migration: %v", err)
	}

	router := gin.Default()
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.Logger())

	// Instantiate repositories, services, handlers for each module
	employerProfileRepo := employerprofile.NewEmployerProfileRepository(db)
	studentProfileRepo := studentprofile.NewStudentProfileRepository(db)

	// Create auth service and handler
	authRepo := auth.NewUserRepository(db)
	authService := auth.NewAuthService(authRepo, employerProfileRepo, studentProfileRepo)
	authHandler := auth.NewAuthHandler(authService)

	employerProfileService := employerprofile.NewEmployerProfileService(employerProfileRepo)
	employerProfileHandler := employerprofile.NewEmployerProfileHandler(employerProfileService)

	jobPostRepo := jobpost.NewJobPostRepository(db)
	jobPostService := jobpost.NewJobPostService(jobPostRepo, employerProfileRepo)
	jobPostHandler := jobpost.NewJobPostHandler(jobPostService)

	applicationRepo := application.NewApplicationRepository(db)
	applicationService := application.NewApplicationService(applicationRepo)
	applicationHandler := application.NewApplicationHandler(applicationService)

	employerAppRepo := employerapplication.NewEmployerApplicationRepository(db)
	employerAppService := employerapplication.NewEmployerApplicationService(employerAppRepo)
	employerAppHandler := employerapplication.NewEmployerApplicationHandler(employerAppService)

	bookmarkRepo := bookmark.NewBookmarkRepository(db)
	jobRepo := jobpost.NewJobPostRepository(db)
	bookmarkService := bookmark.NewBookmarkService(bookmarkRepo, jobRepo)
	bookmarkHandler := bookmark.NewBookmarkHandler(bookmarkService)

	studentProfileService := studentprofile.NewStudentProfileService(studentProfileRepo)
	studentProfileHandler := studentprofile.NewStudentProfileHandler(studentProfileService)

	// Use binary storage service for file serving from database
	fileServeHandler := storage.NewFileServeHandler(db)

	// Create a mock storage service since file uploads are now handled directly in profile handlers
	// The traditional file upload endpoints are not needed for binary storage
	mockStorageService := &mockStorageService{}
	storageHandler := storage.NewStorageHandler(mockStorageService)

	// Notification module with preferences
	notificationPrefsRepo := notification.NewNotificationPreferencesRepository(db)
	notificationService := notification.NewNotificationService(notificationPrefsRepo)
	notificationHandler := notification.NewNotificationHandler(notificationService)

	jobService := worker.NewInMemoryJobService(100)
	workerHandler := worker.NewWorkerHandler(jobService)

	// Admin module
	adminRepo := admin.NewAdminRepository(db)
	adminService := admin.NewAdminService(adminRepo)
	adminHandler := admin.NewAdminHandler(adminService)

	// Health check endpoint (no auth required)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "healthy",
			"service":   "asa-backend",
			"timestamp": time.Now().Unix(),
		})
	})

	// Swagger documentation endpoint (no auth required)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	api := router.Group("/api")

	// Auth routes (no auth middleware)
	auth.RegisterRoutes(api, authHandler)

	// Public job routes (no auth required)
	jobpost.RegisterPublicRoutes(api, jobPostHandler)

	// Public file serving routes (no auth required for file access)
	storage.RegisterPublicRoutes(api, storageHandler, fileServeHandler)

	// Middleware for authenticated routes
	authGroup := api.Group("/")
	authGroup.Use(middleware.AuthMiddleware())

	// Admin routes (require admin role)
	admin.RegisterRoutes(authGroup, adminHandler)

	// Employer profile
	employerprofile.RegisterRoutes(authGroup, employerProfileHandler)

	// Job post (authenticated routes)
	jobpost.RegisterAuthenticatedRoutes(authGroup, jobPostHandler)

	// Applications (student)
	application.RegisterRoutes(authGroup, applicationHandler)

	// Employer-side application management
	employerapplication.RegisterRoutes(authGroup, employerAppHandler)

	// Bookmarks
	bookmark.RegisterRoutes(authGroup, bookmarkHandler)

	// Student profile
	studentprofile.RegisterRoutes(authGroup, studentProfileHandler)

	// File upload routes (require auth)
	storage.RegisterAuthenticatedRoutes(authGroup, storageHandler, fileServeHandler)

	// Notifications
	notification.RegisterRoutes(authGroup, notificationHandler)

	// Workers (background job queue)
	worker.RegisterRoutes(authGroup, workerHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	router.Run(":" + port)
}
