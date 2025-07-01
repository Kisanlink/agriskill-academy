// File: internal/bookmark/routes.go

package bookmark

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *BookmarkHandler) {
	// These will NOT conflict with /jobs/:id or other job routes
	rg.POST("/bookmarks/:jobId", handler.Save)
	rg.DELETE("/bookmarks/:jobId", handler.Remove)
	rg.GET("/jobs/saved", handler.GetSaved)
}
