package bookmark

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *BookmarkHandler) {
	rg.POST("/bookmarks/:jobId", handler.Save)
	rg.DELETE("/bookmarks/:jobId", handler.Remove)
}
