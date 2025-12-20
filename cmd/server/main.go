// File: cmd/server/main.go

// @title AgriJobs API
// @version 1.0
// @description AgriJobs backend API for agricultural job portal
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@agrijobs.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"reflect"
	"time"

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
	"github.com/Kisanlink/agriskill-academy/pkg/firebase"

	"github.com/gin-gonic/gin"
	scalar "github.com/MarceloPetrucio/go-scalar-api-reference"
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

// startWorker starts a background worker that processes jobs from the queue
func startWorker(jobService worker.JobService, notificationService notification.NotificationService, logger *zap.Logger) {
	go func() {
		ticker := time.NewTicker(2 * time.Second) // Check every 2 seconds
		defer ticker.Stop()

		logger.Info("Background worker started - processing jobs every 2 seconds")

		for range ticker.C {
			// Dequeue a job
			job, err := jobService.Dequeue()
			if err != nil {
				logger.Warn("Error dequeuing job", zap.Error(err))
				continue
			}

			if job == nil {
				continue // No jobs available
			}

			logger.Info("Processing job",
				zap.String("job_id", job.ID),
				zap.String("job_type", job.Type),
			)

			// Process the job based on type
			switch job.Type {
			case "send_email":
				err := notification.HandleSendEmail(job.Payload, notificationService)
				if err != nil {
					logger.Error("Failed to process email job",
						zap.String("job_id", job.ID),
						zap.Error(err),
					)
					if failErr := jobService.Fail(job.ID, err.Error()); failErr != nil {
						logger.Error("Failed to mark job as failed",
							zap.String("job_id", job.ID),
							zap.Error(failErr),
						)
					}
				} else {
					logger.Info("Successfully processed email job",
						zap.String("job_id", job.ID),
					)
					if completeErr := jobService.Complete(job.ID, "Email sent successfully"); completeErr != nil {
						logger.Error("Failed to mark job as completed",
							zap.String("job_id", job.ID),
							zap.Error(completeErr),
						)
					}
				}
			default:
				logger.Warn("Unknown job type",
					zap.String("job_id", job.ID),
					zap.String("job_type", job.Type),
				)
				if failErr := jobService.Fail(job.ID, fmt.Sprintf("Unknown job type: %s", job.Type)); failErr != nil {
					logger.Error("Failed to mark job as failed",
						zap.String("job_id", job.ID),
						zap.Error(failErr),
					)
				}
			}
		}
	}()
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

	// Initialize Firebase services (optional - for authentication and email)
	var firebaseEmail auth.FirebaseEmailService
	var firebaseAuth auth.FirebaseAuthClient

	if cfg.FirebaseProjectID != "" && (cfg.FirebaseCredentialsPath != "" || cfg.FirebaseCredentialsJSON != "") && cfg.FirebaseWebAPIKey != "" {
		log.Printf("Initializing Firebase services...")

		// Initialize Firebase Email Service (for sending verification/reset emails)
		fbEmail, err := firebase.NewEmailService(
			cfg.FirebaseCredentialsPath,
			cfg.FirebaseCredentialsJSON,
			cfg.FirebaseWebAPIKey,
		)
		if err != nil {
			logger.Warn("Failed to initialize Firebase email service",
				zap.Error(err))
		} else {
			firebaseEmail = fbEmail
			logger.Info("Firebase email service initialized successfully")
		}

		// Initialize Firebase Auth Client (for authentication via Firebase REST API)
		firebaseAuth = firebase.NewAuthAdapter(cfg.FirebaseWebAPIKey)
		logger.Info("Firebase authentication client initialized successfully")

		logger.Info("Using Firebase for authentication and email verification")
	} else {
		logger.Info("Firebase not configured - using local authentication only")
	}

	// Services and handlers
	authService := auth.NewAuthService(authRepo, employerProfileRepo, studentProfileRepo, firebaseEmail, firebaseAuth)
	authHandler := auth.NewAuthHandler(authService)

	employerProfileService := employerprofile.NewEmployerProfileService(employerProfileRepo)
	employerProfileHandler := employerprofile.NewEmployerProfileHandler(employerProfileService, storageService)

	// Initialize notification and worker services first (needed for job post and application handlers)
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

	// Initialize Email Sender Service for email notifications
	// Create an adapter to bridge JobService interface with JobEnqueuer interface
	enqueuer := notification.JobEnqueuerFunc(func(job interface{}) error {
		// Try to convert to BackgroundJob structure directly
		if bgJob, ok := job.(*worker.BackgroundJob); ok {
			return jobService.Enqueue(bgJob)
		}
		// Try map conversion
		jobMap, ok := job.(map[string]interface{})
		if ok {
			typeStr, _ := jobMap["type"].(string)
			payload, _ := jobMap["payload"].(map[string]interface{})
			bgJob := &worker.BackgroundJob{
				Type:    typeStr,
				Payload: payload,
			}
			return jobService.Enqueue(bgJob)
		}
		// Try struct with Type and Payload fields using reflection
		jobVal := reflect.ValueOf(job)
		if jobVal.Kind() == reflect.Ptr {
			jobVal = jobVal.Elem()
		}
		if jobVal.Kind() == reflect.Struct {
			typeField := jobVal.FieldByName("Type")
			payloadField := jobVal.FieldByName("Payload")
			if typeField.IsValid() && payloadField.IsValid() {
				bgJob := &worker.BackgroundJob{
					Type:    typeField.String(),
					Payload: payloadField.Interface().(map[string]interface{}),
				}
				return jobService.Enqueue(bgJob)
			}
		}
		return fmt.Errorf("invalid job type: %T", job)
	})
	emailSenderService := notification.NewEmailSenderService(notificationService, enqueuer)

	jobPostRepo := jobpost.NewJobPostRepository(db)
	jobPostService := jobpost.NewJobPostService(jobPostRepo, employerProfileRepo)
	jobPostHandler := jobpost.NewJobPostHandler(jobPostService, emailSenderService, db, notificationService)

	applicationRepo := application.NewApplicationRepository(db)
	applicationService := application.NewApplicationService(applicationRepo, jobPostRepo, s3Manager, emailSenderService, db, notificationService)
	applicationHandler := application.NewApplicationHandler(applicationService)

	employerAppRepo := employerapplication.NewEmployerApplicationRepository(db)
	employerAppService := employerapplication.NewEmployerApplicationService(employerAppRepo)
	employerAppHandler := employerapplication.NewEmployerApplicationHandler(
		employerAppService,
		emailSenderService,
		db,
	)

	bookmarkRepo := bookmark.NewBookmarkRepository(db)
	bookmarkService := bookmark.NewBookmarkService(bookmarkRepo, jobPostRepo)
	bookmarkHandler := bookmark.NewBookmarkHandler(bookmarkService)

	studentProfileService := studentprofile.NewStudentProfileService(studentProfileRepo)
	studentProfileHandler := studentprofile.NewStudentProfileHandler(studentProfileService, storageService)

	// File serving and storage handlers
	fileServeHandler := storage.NewFileServeHandler(s3Manager, db)
	storageHandler := storage.NewStorageHandler(storageService)

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

	// Scalar API documentation
	router.GET("/docs", func(c *gin.Context) {
		// Read the swagger.json file
		swaggerJSON, err := os.ReadFile("docs/swagger.json")
		if err != nil {
			c.String(500, "Error reading swagger spec: %v", err)
			return
		}

		htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
			SpecContent: string(swaggerJSON),
			CustomOptions: scalar.CustomOptions{
				PageTitle: "AgriJobs API Documentation",
			},
			DarkMode: true,
		})
		if err != nil {
			c.String(500, "Error generating Scalar documentation: %v", err)
			return
		}
		c.Data(200, "text/html; charset=utf-8", []byte(htmlContent))
	})

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

	// Start background worker to process jobs
	startWorker(jobService, notificationService, logger)
	logger.Info("Background worker initialized")

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
