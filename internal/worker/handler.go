// File: internal/worker/handler.go

package worker

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type WorkerHandler struct {
	service JobService
}

func NewWorkerHandler(s JobService) *WorkerHandler {
	return &WorkerHandler{s}
}

// POST /worker/job
func (h *WorkerHandler) EnqueueJob(c *gin.Context) {
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
