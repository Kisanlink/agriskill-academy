package jobpost

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type JobPostHandler struct {
	service JobPostService
}

func NewJobPostHandler(s JobPostService) *JobPostHandler {
	return &JobPostHandler{s}
}

// POST /jobs
func (h *JobPostHandler) Create(c *gin.Context) {
	var req JobPost
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request"})
		return
	}
	employerID := c.GetString("user_id")
	if employerID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Missing employer ID"})
		return
	}
	req.EmployerID = employerID
	req.Status = "published" // default to published on Create

	if err := h.service.Create(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to create job"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Job created", "jobPost": req})
}

// POST /jobs/draft
func (h *JobPostHandler) CreateDraft(c *gin.Context) {
	var req JobPost
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request"})
		return
	}
	employerID := c.GetString("user_id")
	if employerID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Missing employer ID"})
		return
	}
	req.EmployerID = employerID
	req.Status = "draft"

	if err := h.service.Create(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to save draft"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Draft saved", "jobPost": req})
}

// POST /jobs/publish
func (h *JobPostHandler) Publish(c *gin.Context) {
	var req JobPost
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request"})
		return
	}
	employerID := c.GetString("user_id")
	if employerID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Missing employer ID"})
		return
	}
	req.EmployerID = employerID
	req.Status = "published"

	if err := h.service.Create(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to publish job"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Job published", "jobPost": req})
}

// PUT /jobs/:id
func (h *JobPostHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req JobPost
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request"})
		return
	}
	req.ID = id
	if err := h.service.Update(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to update job"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Job updated", "jobPost": req})
}

// DELETE /jobs/:id
func (h *JobPostHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to delete job"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Job deleted"})
}

// GET /jobs/:id
func (h *JobPostHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	job, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Job not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "jobPost": job})
}

// GET /jobs/my-posts
func (h *JobPostHandler) GetByEmployer(c *gin.Context) {
	employerID := c.GetString("user_id")
	jobs, err := h.service.GetByEmployer(employerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to fetch jobs"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "jobPosts": jobs})
}

// POST /jobs/search
func (h *JobPostHandler) Search(c *gin.Context) {
	var filter JobPostFilter
	if err := c.ShouldBindJSON(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request"})
		return
	}
	jobs, err := h.service.Search(&filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Search failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "jobs": jobs})
}
