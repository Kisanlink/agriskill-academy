package bookmark

import (
	"asa/pkg/authz"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type BookmarkHandler struct {
	service BookmarkService
}

func NewBookmarkHandler(s BookmarkService) *BookmarkHandler {
	return &BookmarkHandler{s}
}

func getJWT(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		return authHeader[7:]
	}
	return ""
}

// POST /jobs/:jobId/save
func (h *BookmarkHandler) Save(c *gin.Context) {
	username := c.GetString("email")
	jobID := c.Param("jobId")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckAAAPermission(username, "db_asa_bookmarks", "create", jobID, jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	userID := c.GetString("user_id")
	if err := h.service.Save(userID, jobID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to bookmark job"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Job saved"})
}

// DELETE /jobs/:jobId/unsave
func (h *BookmarkHandler) Remove(c *gin.Context) {
	username := c.GetString("email")
	jobID := c.Param("jobId")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckAAAPermission(username, "db_asa_bookmarks", "delete", jobID, jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	userID := c.GetString("user_id")
	if err := h.service.Remove(userID, jobID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to remove bookmark"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Removed from bookmarks"})
}

// GET /jobs/saved
func (h *BookmarkHandler) GetSaved(c *gin.Context) {
	username := c.GetString("email")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckAAAPermission(username, "db_asa_bookmarks", "read", "", jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	userID := c.GetString("user_id")
	jobs, err := h.service.GetByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Could not fetch saved jobs"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "jobs": jobs})
}
