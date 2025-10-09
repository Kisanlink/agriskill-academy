package bookmark

import (
	"github.com/Kisanlink/agriskill-academy/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *BookmarkHandler) {
	// Bookmarks are student-only functionality
	bookmarks := rg.Group("/bookmarks")
	bookmarks.Use(middleware.RequireRole("student"))
	{
		bookmarks.GET("", handler.GetSaved)         // Get user's bookmarks
		bookmarks.POST("/:jobId", handler.Save)     // Save a job bookmark
		bookmarks.DELETE("/:jobId", handler.Remove) // Remove a job bookmark
	}
}
