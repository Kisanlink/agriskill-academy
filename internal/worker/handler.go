// File: internal/worker/handler.go

package worker

import (
	"asa/pkg/authz"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type WorkerHandler struct {
	service JobService
}

func NewWorkerHandler(s JobService) *WorkerHandler {
	return &WorkerHandler{s}
}

func getJWT(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		return authHeader[7:]
	}
	return ""
}

// @Summary Enqueue Background Job
// @Description Enqueue a background job for processing
// @Tags Background Jobs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body map[string]interface{} true "Job data"
// @Success 200 {object} map[string]interface{} "Job enqueued successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Permission denied"
// @Failure 503 {object} map[string]interface{} "Queue is full"
// @Router /api/worker/job [post]
// POST /worker/job
func (h *WorkerHandler) EnqueueJob(c *gin.Context) {
	username := c.GetString("username")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckLocalPermission(username, "db_asa_jobs", "create", "", jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	var req struct {
		Type    string                 `json:"type" binding:"required"`
		Payload map[string]interface{} `json:"payload"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request"})
		return
	}
	job := &Job{
		ID:        "", // generate UUID in production
		Type:      req.Type,
		Payload:   req.Payload,
		CreatedAt: time.Now(),
	}
	if err := h.service.Enqueue(job); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"success": false, "message": "Queue is full"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Job enqueued"})
}
