// File: internal/bookmark/handler.go

package bookmark

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type BookmarkHandler struct {
	service BookmarkService
}

func NewBookmarkHandler(s BookmarkService) *BookmarkHandler {
	return &BookmarkHandler{s}
}

// POST /jobs/:jobId/save
func (h *BookmarkHandler) Save(c *gin.Context) {
	userID := c.GetString("user_id")
	jobID := c.Param("jobId")
	if err := h.service.Save(userID, jobID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to bookmark job"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Job saved"})
}

// DELETE /jobs/:jobId/unsave
func (h *BookmarkHandler) Remove(c *gin.Context) {
	userID := c.GetString("user_id")
	jobID := c.Param("jobId")
	if err := h.service.Remove(userID, jobID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to remove bookmark"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Removed from bookmarks"})
}

// GET /jobs/saved
func (h *BookmarkHandler) GetSaved(c *gin.Context) {
	userID := c.GetString("user_id")
	bookmarks, err := h.service.GetByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Could not fetch saved jobs"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "jobs": bookmarks})
}
