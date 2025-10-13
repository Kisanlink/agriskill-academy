// File: cmd/server/main.go

package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Kisanlink/agriskill-academy/config"
	"github.com/Kisanlink/agriskill-academy/internal/admin"
	"github.com/Kisanlink/agriskill-academy/internal/application"
	"github.com/Kisanlink/agriskill-academy/internal/auth"
	"github.com/Kisanlink/agriskill-academy/internal/bookmark"
	"github.com/Kisanlink/agriskill-academy/internal/contact"
	"github.com/Kisanlink/agriskill-academy/internal/employerapplication"
	"github.com/Kisanlink/agriskill-academy/internal/employerprofile"

	"github.com/Kisanlink/agriskill-academy/internal/jobpost"
	"github.com/Kisanlink/agriskill-academy/internal/middleware"
	"github.com/Kisanlink/agriskill-academy/internal/notification"
	"github.com/Kisanlink/agriskill-academy/internal/seeding"
	"github.com/Kisanlink/agriskill-academy/internal/storage"
	"github.com/Kisanlink/agriskill-academy/internal/studentprofile"
	"github.com/Kisanlink/agriskill-academy/internal/worker"

	_ "github.com/Kisanlink/agriskill-academy/docs" // Import swagger docs

	kdb "github.com/Kisanlink/agriskill-academy/pkg/db"

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

	application.InitializeCounterFromDatabase(db)
	auth.InitializeCounterFromDatabase(db)
	bookmark.InitializeCounterFromDatabase(db)
	contact.InitializeCounterFromDatabase(db)
	employerapplication.InitializeCounterFromDatabase(db)
	employerprofile.InitializeCounterFromDatabase(db)
	jobpost.InitializeCounterFromDatabase(db)
	notification.InitializeCounterFromDatabase(db)
	studentprofile.InitializeCounterFromDatabase(db)

	return nil
}

func main() {
	// Load complete configuration
	cfg := config.LoadConfig()

	// Initialize structured logging
	loggerConfig := middleware.LoggerConfig{
		Level:       cfg.LogLevel,
		OutputPath:  cfg.LogOutputPath,
		Format:      cfg.LogFormat,
		Development: cfg.LogDevelopment,
	}

	if err := middleware.InitLogger(loggerConfig); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	logger := middleware.GetLogger()
	logger.Info("Starting AgriJobs application",
		zap.String("version", cfg.AppVersion),
		zap.String("environment", cfg.AppEnv),
	)

	// Initialize config and DB
	db, err := config.InitDB()
	if err != nil {
		logger.Fatal("Failed to initialize DB", zap.Error(err))
	}
	defer config.CloseDB(db)

	// Run migrations
	if err := runAutoMigrate(db); err != nil {
		log.Fatalf("Failed to run auto migration: %v", err)
	}

	// Run seeding
	seedingService := seeding.NewSeedingService(db)
	if err := seedingService.RunSeeding(); err != nil {
		logger.Fatal("Failed to run seeding", zap.Error(err))
	}

	// Storage Service setup
	var storageService storage.StorageService
	var s3Manager *kdb.S3Manager

	// Check if AWS S3 is configured
	s3Bucket := cfg.AWSS3Bucket
	s3AccessKeyID := cfg.AWSAccessKeyID
	s3SecretAccessKey := cfg.AWSSecretKey

	baseURL := cfg.ASABaseURL
	if baseURL == "" {
		logger.Fatal("ASA_BASE_URL environment variable is required")
	}

	// Only S3 bucket is required. Credentials are optional:
	// - If provided (MinIO/local): uses explicit credentials
	// - If empty (AWS ECS): AWS SDK automatically uses IAM Task Role
	if s3Bucket != "" {
		s3Region := cfg.AWSRegion
		if s3Region == "" {
			logger.Fatal("AWS_REGION environment variable is required")
		}

		// Log storage configuration
		if s3AccessKeyID != "" && s3SecretAccessKey != "" {
			log.Printf("Using S3 storage with explicit credentials (MinIO/local): bucket=%s", s3Bucket)
		} else {
			log.Printf("Using S3 storage with IAM role (AWS ECS): bucket=%s", s3Bucket)
		}

		s3Endpoint := cfg.AWSS3Endpoint
		s3ForcePathStyle := cfg.AWSS3ForcePathStyle
		s3DisableSSL := cfg.AWSS3DisableSSL

		s3Config := &kdb.Config{
			S3Region:          s3Region,
			S3Bucket:          s3Bucket,
			S3Endpoint:        s3Endpoint,
			S3ForcePathStyle:  s3ForcePathStyle,
			S3DisableSSL:      s3DisableSSL,
			S3AccessKeyID:     s3AccessKeyID,     // Optional: empty for IAM role
			S3SecretAccessKey: s3SecretAccessKey, // Optional: empty for IAM role
			LogLevel:          "info",
		}
		s3Logger := zap.NewNop()
		s3Manager = kdb.NewS3Manager(s3Config, s3Logger)
		if err := s3Manager.Connect(context.Background()); err != nil {
			log.Fatalf("Failed to connect to S3: %v", err)
		}
		storageService = storage.NewS3StorageService(s3Manager, s3Bucket, baseURL)
	} else {
		log.Fatalf("AWS S3 configuration is required. Please set AWS_S3_BUCKET and AWS_REGION environment variables.")
	}

	// Create Gin router with production middleware
	router := gin.Default()

	// Production security and monitoring middleware
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.StructuredLoggingMiddleware())
	router.Use(middleware.PerformanceMonitoringMiddleware())
	router.Use(middleware.ErrorLoggingMiddleware())
	router.Use(middleware.SecurityHeadersMiddleware())
	router.Use(middleware.RateLimitMiddleware(cfg.RateLimitRequests, cfg.RateLimitWindow))
	router.Use(middleware.InputSanitizationMiddleware())
	router.Use(middleware.SQLInjectionProtectionMiddleware())
	router.Use(middleware.RequestSizeLimitMiddleware(cfg.MaxRequestSize))
	router.Use(middleware.ContextTimeoutMiddleware(cfg.HealthCheckTimeout))

	// CORS middleware
	if cfg.EnableCORS {
		router.Use(middleware.CORSMiddleware())
		router.Use(middleware.CORSValidationMiddleware(cfg.AllowedOrigins))
	}

	// Legacy middleware (keeping for compatibility)
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

	// Initialize Redis job service
	jobService, err := worker.NewRedisJobService(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	if err != nil {
		logger.Fatal("Failed to initialize Redis job service", zap.Error(err))
	}
	defer jobService.Close()

	workerHandler := worker.NewWorkerHandler(jobService)

	adminRepo := admin.NewAdminRepository(db)
	adminService := admin.NewAdminService(adminRepo)
	adminHandler := admin.NewAdminHandler(adminService)

	contactRepo := contact.NewContactRepository(db)
	contactService := contact.NewContactService(contactRepo)
	contactHandler := contact.NewContactHandler(contactService)

	// Health check with production monitoring
	router.GET("/health", middleware.HealthCheckMiddleware())
	router.GET("/health/db", middleware.DatabaseHealthCheck(db))

	// Swagger docs
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Public API group
	api := router.Group("/api")
	auth.RegisterPublicRoutes(api, authHandler)
	jobpost.RegisterPublicRoutes(api, jobPostHandler)
	storage.RegisterPublicRoutes(api, storageHandler, fileServeHandler)
	contact.RegisterPublicRoutes(api, contactHandler)

	// Protected routes
	authGroup := api.Group("/")
	authGroup.Use(middleware.AuthMiddleware())

	auth.RegisterProtectedRoutes(authGroup, authHandler)
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
	port := cfg.ServerPort
	if port == "" {
		logger.Fatal("SERVER_PORT environment variable is required")
	}

	logger.Info("Starting ASA backend server",
		zap.String("port", port),
		zap.String("environment", cfg.AppEnv),
	)

	if err := router.Run(":" + port); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
