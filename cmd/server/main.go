// File: cmd/server/main.go

package main

import (
	"log"
	"os"

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

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize config and DB
	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize DB: %v", err)
	}
	defer config.CloseDB(db)

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

	storageService := storage.NewLocalStorageService("uploads", "http://localhost:3000/api/files")
	storageHandler := storage.NewStorageHandler(storageService)

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

	// API routes
	api := router.Group("/api")

	// Auth routes (no auth middleware)
	auth.RegisterRoutes(api, authHandler)

	// Public job routes (no auth required)
	jobpost.RegisterPublicRoutes(api, jobPostHandler)

	// Public file serving routes (no auth required for file access)
	storage.RegisterPublicRoutes(api, storageHandler)

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
	storage.RegisterAuthenticatedRoutes(authGroup, storageHandler)

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
