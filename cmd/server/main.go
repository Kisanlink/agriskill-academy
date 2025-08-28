// File: cmd/server/main.go

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"asa/config"
	"asa/internal/admin"
	"asa/internal/application"
	"asa/internal/auth"
	"asa/internal/bookmark"
	"asa/internal/contact"
	"asa/internal/employerapplication"
	"asa/internal/employerprofile"

	"asa/internal/jobpost"
	"asa/internal/middleware"
	"asa/internal/notification"
	"asa/internal/storage"
	"asa/internal/studentprofile"
	"asa/internal/worker"

	_ "asa/docs" // Import swagger docs

	kdb "asa/pkg/db"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
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
		&contact.ContactRequest{},
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

	// Storage Service setup
	var storageService storage.StorageService
	var s3Manager *kdb.S3Manager

	// Check if AWS S3 is configured
	s3Bucket := os.Getenv("AWS_S3_BUCKET")
	s3AccessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	s3SecretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	baseURL := os.Getenv("ASA_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:3000/api"
	}

	if s3Bucket != "" && s3AccessKeyID != "" && s3SecretAccessKey != "" {
		// Use S3 storage
		log.Printf("Using S3 storage with bucket: %s", s3Bucket)
		s3Region := os.Getenv("AWS_REGION")
		if s3Region == "" {
			s3Region = "us-east-1"
		}
		s3Endpoint := os.Getenv("AWS_S3_ENDPOINT")
		s3ForcePathStyle := os.Getenv("AWS_S3_FORCE_PATH_STYLE") == "true"
		s3DisableSSL := os.Getenv("AWS_S3_DISABLE_SSL") == "true"

		s3Config := &kdb.Config{
			S3Region:          s3Region,
			S3Bucket:          s3Bucket,
			S3Endpoint:        s3Endpoint,
			S3ForcePathStyle:  s3ForcePathStyle,
			S3DisableSSL:      s3DisableSSL,
			S3AccessKeyID:     s3AccessKeyID,
			S3SecretAccessKey: s3SecretAccessKey,
			LogLevel:          "info",
		}
		s3Logger := zap.NewNop()
		s3Manager = kdb.NewS3Manager(s3Config, s3Logger)
		if err := s3Manager.Connect(context.Background()); err != nil {
			log.Fatalf("Failed to connect to S3: %v", err)
		}
		storageService = storage.NewS3StorageService(s3Manager, s3Bucket, baseURL)
	} else {
		log.Fatalf("AWS S3 configuration is required. Please set AWS_S3_BUCKET, AWS_ACCESS_KEY_ID, and AWS_SECRET_ACCESS_KEY environment variables.")
	}

	// Create Gin router with middleware
	router := gin.Default()
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.Logger())

	// Repositories
	employerProfileRepo := employerprofile.NewEmployerProfileRepository(db)
	studentProfileRepo := studentprofile.NewStudentProfileRepository(db)
	authRepo := auth.NewUserRepository(db)

	// Using local authentication only - no AAA service dependency
	log.Printf("Using local authentication with kisanlink-db")

	// Services and handlers
	authService := auth.NewAuthService(authRepo, employerProfileRepo, studentProfileRepo)
	authHandler := auth.NewAuthHandler(authService)

	employerProfileService := employerprofile.NewEmployerProfileService(employerProfileRepo)
	employerProfileHandler := employerprofile.NewEmployerProfileHandler(employerProfileService, storageService)

	jobPostRepo := jobpost.NewJobPostRepository(db)
	jobPostService := jobpost.NewJobPostService(jobPostRepo, employerProfileRepo)
	jobPostHandler := jobpost.NewJobPostHandler(jobPostService)

	applicationRepo := application.NewApplicationRepository(db)
	applicationService := application.NewApplicationService(applicationRepo, jobPostRepo, s3Manager)
	applicationHandler := application.NewApplicationHandler(applicationService)

	employerAppRepo := employerapplication.NewEmployerApplicationRepository(db)
	employerAppService := employerapplication.NewEmployerApplicationService(employerAppRepo)
	employerAppHandler := employerapplication.NewEmployerApplicationHandler(employerAppService)

	bookmarkRepo := bookmark.NewBookmarkRepository(db)
	bookmarkService := bookmark.NewBookmarkService(bookmarkRepo, jobPostRepo)
	bookmarkHandler := bookmark.NewBookmarkHandler(bookmarkService)

	studentProfileService := studentprofile.NewStudentProfileService(studentProfileRepo)
	studentProfileHandler := studentprofile.NewStudentProfileHandler(studentProfileService, storageService)

	// File serving and storage handlers
	fileServeHandler := storage.NewFileServeHandler(s3Manager, db)
	storageHandler := storage.NewStorageHandler(storageService)

	notificationPrefsRepo := notification.NewNotificationPreferencesRepository(db)
	notificationService := notification.NewNotificationService(notificationPrefsRepo)
	notificationHandler := notification.NewNotificationHandler(notificationService)

	jobService := worker.NewInMemoryJobService(100)
	workerHandler := worker.NewWorkerHandler(jobService)

	adminRepo := admin.NewAdminRepository(db)
	adminService := admin.NewAdminService(adminRepo)
	adminHandler := admin.NewAdminHandler(adminService)

	contactRepo := contact.NewContactRepository(db)
	contactService := contact.NewContactService(contactRepo)
	contactHandler := contact.NewContactHandler(contactService)

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
	contact.RegisterPublicRoutes(api, contactHandler)

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
	contact.RegisterAdminRoutes(authGroup, contactHandler)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Starting ASA backend server on port %s", port)
	router.Run(":" + port)
}
