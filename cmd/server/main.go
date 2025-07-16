// File: cmd/server/main.go

package main

import (
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"time"

	"asa/config"
	"asa/internal/admin"
	"asa/internal/application"
	"asa/internal/auth"
	"asa/internal/bookmark"
	"asa/internal/employerapplication"
	"asa/internal/employerprofile"
	"asa/internal/grpc"
	"asa/internal/jobpost"
	"asa/internal/middleware"
	"asa/internal/notification"
	"asa/internal/storage"
	"asa/internal/studentprofile"
	"asa/internal/worker"

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
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		log.Printf("Warning: Could not create uuid-ossp extension: %v", err)
	}

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

	for _, m := range models {
		if err := db.AutoMigrate(m); err != nil {
			return fmt.Errorf("failed to migrate model %T: %w", m, err)
		}
		log.Printf("Successfully migrated: %T", m)
	}

	log.Println("AutoMigrate completed successfully!")
	return nil
}

// mockStorageService satisfies storage.Service but disables file uploads
type mockStorageService struct{}

func (m *mockStorageService) SaveFile(_ *multipart.FileHeader, _ string) (string, error) {
	return "", fmt.Errorf("file uploads are now handled in profile handlers")
}
func (m *mockStorageService) SaveImage(_ *multipart.FileHeader, _ string) (string, error) {
	return "", fmt.Errorf("file uploads are now handled in profile handlers")
}
func (m *mockStorageService) SaveDocument(_ *multipart.FileHeader, _ string) (string, error) {
	return "", fmt.Errorf("file uploads are now handled in profile handlers")
}
func (m *mockStorageService) SaveResume(_ *multipart.FileHeader, _ string) (string, error) {
	return "", fmt.Errorf("file uploads are now handled in profile handlers")
}
func (m *mockStorageService) DeleteFile(string) error {
	return fmt.Errorf("file operations are now handled in profile handlers")
}
func (m *mockStorageService) ListFiles(_ string) ([]storage.FileInfo, error) {
	return nil, nil
}
func (m *mockStorageService) GetFileInfo(string) (*storage.FileInfo, error) {
	return nil, fmt.Errorf("file operations are now handled in profile handlers")
}

func main() {
	// Initialize config and DB
	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize DB: %v", err)
	}
	defer config.CloseDB(db)

	// Run migrations
	if err := runAutoMigrate(db); err != nil {
		log.Fatalf("Failed to run auto migration: %v", err)
	}

	// Create Gin router with middleware
	router := gin.Default()
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.Logger())

	// Repositories
	employerProfileRepo := employerprofile.NewEmployerProfileRepository(db)
	studentProfileRepo := studentprofile.NewStudentProfileRepository(db)
	authRepo := auth.NewUserRepository(db)

	// Initialize AAA gRPC client (optional - will use local auth if not available)
	var aaaClient *grpc.AAAGrpcClient = nil
	aaaHost := os.Getenv("AAA_HOST")
	aaaPort := os.Getenv("AAA_GRPC_PORT")

	if aaaHost != "" && aaaPort != "" {
		aaaEndpoint := fmt.Sprintf("%s:%s", aaaHost, aaaPort)
		log.Printf("Attempting to connect to AAA gRPC service at: %s", aaaEndpoint)

		// Try to initialize gRPC client, but don't fail if it's not available
		// This allows the service to run with local authentication only
		grpcClient, err := initGrpcClient(aaaEndpoint)
		if err != nil {
			log.Printf("Warning: Failed to connect to AAA gRPC service: %v", err)
			log.Printf("Continuing with local authentication only")
		} else {
			aaaClient = grpcClient
			log.Printf("Successfully connected to AAA gRPC service")
		}
	} else {
		log.Printf("AAA_HOST or AAA_GRPC_PORT not set, using local authentication only")
	}

	// Services and handlers
	authService := auth.NewAuthService(authRepo, employerProfileRepo, studentProfileRepo, aaaClient)
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
	bookmarkService := bookmark.NewBookmarkService(bookmarkRepo, jobPostRepo)
	bookmarkHandler := bookmark.NewBookmarkHandler(bookmarkService)

	studentProfileService := studentprofile.NewStudentProfileService(studentProfileRepo)
	studentProfileHandler := studentprofile.NewStudentProfileHandler(studentProfileService)

	// File serving and storage handlers
	fileServeHandler := storage.NewFileServeHandler(db)
	storageHandler := storage.NewStorageHandler(&mockStorageService{})

	notificationPrefsRepo := notification.NewNotificationPreferencesRepository(db)
	notificationService := notification.NewNotificationService(notificationPrefsRepo)
	notificationHandler := notification.NewNotificationHandler(notificationService)

	jobService := worker.NewInMemoryJobService(100)
	workerHandler := worker.NewWorkerHandler(jobService)

	adminRepo := admin.NewAdminRepository(db)
	adminService := admin.NewAdminService(adminRepo)
	adminHandler := admin.NewAdminHandler(adminService)

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "healthy",
			"service":   "asa-backend",
			"timestamp": time.Now().Unix(),
		})
	})

	// Swagger docs
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Public API group
	api := router.Group("/api")
	auth.RegisterRoutes(api, authHandler)
	jobpost.RegisterPublicRoutes(api, jobPostHandler)
	storage.RegisterPublicRoutes(api, storageHandler, fileServeHandler)

	// Protected routes
	authGroup := api.Group("/")
	authGroup.Use(middleware.AuthMiddleware())

	admin.RegisterRoutes(authGroup, adminHandler)
	employerprofile.RegisterRoutes(authGroup, employerProfileHandler)
	jobpost.RegisterAuthenticatedRoutes(authGroup, jobPostHandler)
	application.RegisterRoutes(authGroup, applicationHandler)
	employerapplication.RegisterRoutes(authGroup, employerAppHandler)
	bookmark.RegisterRoutes(authGroup, bookmarkHandler)
	studentprofile.RegisterRoutes(authGroup, studentProfileHandler)
	storage.RegisterAuthenticatedRoutes(authGroup, storageHandler, fileServeHandler)
	notification.RegisterRoutes(authGroup, notificationHandler)
	worker.RegisterRoutes(authGroup, workerHandler)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Starting ASA backend server on port %s", port)
	router.Run(":" + port)
}

// initGrpcClient initializes the gRPC client with proper error handling
func initGrpcClient(endpoint string) (*grpc.AAAGrpcClient, error) {
	log.Printf("🔌 Attempting to connect to gRPC service at: %s", endpoint)

	// Create gRPC client using the correct function signature
	client, err := grpc.NewAAAGrpcClient(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server: %w", err)
	}

	log.Printf("✅ Successfully connected to gRPC service")
	return client, nil
}
