// File: cmd/server/main.go

package main

import (
	"AGRIJOBS/config"
	"AGRIJOBS/internal/application"
	"AGRIJOBS/internal/auth"
	"AGRIJOBS/internal/bookmark"
	"AGRIJOBS/internal/employerapplication"
	"AGRIJOBS/internal/employerprofile"
	"AGRIJOBS/internal/jobpost"
	"AGRIJOBS/internal/middleware"
	"AGRIJOBS/internal/notification"
	"AGRIJOBS/internal/storage"
	"AGRIJOBS/internal/userprofile"
	"AGRIJOBS/internal/worker"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize config and DB
	db := config.InitDB()
	defer config.CloseDB(db)

	router := gin.Default()
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.Logger())

	// Instantiate repositories, services, handlers for each module
	authRepo := auth.NewUserRepository(db)
	authService := auth.NewAuthService(authRepo)
	authHandler := auth.NewAuthHandler(authService)

	employerProfileRepo := employerprofile.NewEmployerProfileRepository(db)
	employerProfileService := employerprofile.NewEmployerProfileService(employerProfileRepo)
	employerProfileHandler := employerprofile.NewEmployerProfileHandler(employerProfileService)

	jobPostRepo := jobpost.NewJobPostRepository(db)
	jobPostService := jobpost.NewJobPostService(jobPostRepo)
	jobPostHandler := jobpost.NewJobPostHandler(jobPostService)

	applicationRepo := application.NewApplicationRepository(db)
	applicationService := application.NewApplicationService(applicationRepo)
	applicationHandler := application.NewApplicationHandler(applicationService)

	employerAppRepo := employerapplication.NewEmployerApplicationRepository(db)
	employerAppService := employerapplication.NewEmployerApplicationService(employerAppRepo)
	employerAppHandler := employerapplication.NewEmployerApplicationHandler(employerAppService)

	bookmarkRepo := bookmark.NewBookmarkRepository(db)
	bookmarkService := bookmark.NewBookmarkService(bookmarkRepo)
	bookmarkHandler := bookmark.NewBookmarkHandler(bookmarkService)

	userProfileRepo := userprofile.NewUserProfileRepository(db)
	userProfileService := userprofile.NewUserProfileService(userProfileRepo)
	userProfileHandler := userprofile.NewUserProfileHandler(userProfileService)

	storageService := storage.NewLocalStorageService("uploads")
	storageHandler := storage.NewStorageHandler(storageService)

	notificationService := notification.NewMailService()
	notificationHandler := notification.NewNotificationHandler(notificationService)

	jobService := worker.NewInMemoryJobService(100)
	workerHandler := worker.NewWorkerHandler(jobService)

	// API routes
	api := router.Group("/api")

	// Auth routes
	auth.RegisterRoutes(api, authHandler)

	// Middleware for authenticated routes
	authGroup := api.Group("/")
	authGroup.Use(middleware.AuthMiddleware())

	// Employer profile
	employerprofile.RegisterRoutes(authGroup, employerProfileHandler)

	// Job post
	jobpost.RegisterRoutes(authGroup, jobPostHandler)

	// Applications (student)
	application.RegisterRoutes(authGroup, applicationHandler)

	// Employer-side application management
	employerapplication.RegisterRoutes(authGroup, employerAppHandler)

	// Bookmarks
	bookmark.RegisterRoutes(authGroup, bookmarkHandler)

	// User profile
	userprofile.RegisterRoutes(authGroup, userProfileHandler)

	// File storage
	storage.RegisterRoutes(authGroup, storageHandler)

	// Notifications
	notification.RegisterRoutes(authGroup, notificationHandler)

	// Workers (background job queue)
	worker.RegisterRoutes(authGroup, workerHandler)

	router.Run(":3000")
}
