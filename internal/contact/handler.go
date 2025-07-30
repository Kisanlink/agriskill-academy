package contact

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ContactHandler struct {
	service ContactService
}

func NewContactHandler(service ContactService) *ContactHandler {
	return &ContactHandler{service: service}
}

// @Summary Submit Contact Form
// @Description Submit a contact form (public endpoint, no authentication required)
// @Tags Contact
// @Accept json
// @Produce json
// @Param request body ContactSubmissionRequest true "Contact form data"
// @Success 201 {object} ContactSubmissionResponse "Contact form submitted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request data"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/contact [post]
// POST /contact
func (h *ContactHandler) SubmitContactForm(c *gin.Context) {
	var req ContactSubmissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data",
			"error":   err.Error(),
		})
		return
	}

	response, err := h.service.SubmitContactForm(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to submit contact form",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// @Summary Get Contact Requests (Admin)
// @Description Get paginated list of contact requests for admin panel
// @Tags Admin - Contact
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param status query string false "Filter by status"
// @Param search query string false "Search in name, email, or subject"
// @Param sort_by query string false "Sort by field" default(created_at)
// @Param sort_order query string false "Sort order" default(DESC)
// @Success 200 {object} ContactListResponse "Contact requests retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request parameters"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/admin/contacts [get]
// GET /admin/contacts
func (h *ContactHandler) GetContactRequests(c *gin.Context) {
	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	req := &ContactListRequest{
		Page:      page,
		Limit:     limit,
		Status:    c.Query("status"),
		Search:    c.Query("search"),
		SortBy:    c.DefaultQuery("sort_by", "created_at"),
		SortOrder: c.DefaultQuery("sort_order", "DESC"),
	}

	// Validate parameters
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 || req.Limit > 100 {
		req.Limit = 10
	}

	response, err := h.service.GetContactRequests(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get contact requests",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ContactResponse{
		Success: true,
		Message: "Contact requests retrieved successfully",
		Data:    response,
	})
}

// @Summary Get Contact Request by ID (Admin)
// @Description Get a specific contact request by ID for admin panel
// @Tags Admin - Contact
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Contact request ID"
// @Success 200 {object} ContactResponse "Contact request retrieved successfully"
// @Failure 404 {object} map[string]interface{} "Contact request not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/admin/contacts/{id} [get]
// GET /admin/contacts/:id
func (h *ContactHandler) GetContactByID(c *gin.Context) {
	id := c.Param("id")

	contact, err := h.service.GetContactByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Contact request not found",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ContactResponse{
		Success: true,
		Message: "Contact request retrieved successfully",
		Data:    contact,
	})
}

// @Summary Update Contact Status (Admin)
// @Description Update the status of a contact request
// @Tags Admin - Contact
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Contact request ID"
// @Param request body UpdateContactStatusRequest true "Status update data"
// @Success 200 {object} ContactResponse "Status updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request data"
// @Failure 404 {object} map[string]interface{} "Contact request not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/admin/contacts/{id}/status [put]
// PUT /admin/contacts/:id/status
func (h *ContactHandler) UpdateContactStatus(c *gin.Context) {
	id := c.Param("id")

	var req UpdateContactStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data",
			"error":   err.Error(),
		})
		return
	}

	err := h.service.UpdateContactStatus(id, req.Status)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "contact request not found" {
			status = http.StatusNotFound
		} else if err.Error() == "invalid status: "+req.Status {
			status = http.StatusBadRequest
		}

		c.JSON(status, gin.H{
			"success": false,
			"message": "Failed to update contact status",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ContactResponse{
		Success: true,
		Message: "Contact status updated successfully",
	})
}

// @Summary Delete Contact Request (Admin)
// @Description Delete a contact request
// @Tags Admin - Contact
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Contact request ID"
// @Success 200 {object} ContactResponse "Contact request deleted successfully"
// @Failure 404 {object} map[string]interface{} "Contact request not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/admin/contacts/{id} [delete]
// DELETE /admin/contacts/:id
func (h *ContactHandler) DeleteContact(c *gin.Context) {
	id := c.Param("id")

	err := h.service.DeleteContact(id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "contact request not found" {
			status = http.StatusNotFound
		}

		c.JSON(status, gin.H{
			"success": false,
			"message": "Failed to delete contact request",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ContactResponse{
		Success: true,
		Message: "Contact request deleted successfully",
	})
}

// @Summary Get Contact Analytics (Admin)
// @Description Get contact request analytics for admin dashboard
// @Tags Admin - Contact
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} ContactResponse "Contact analytics retrieved successfully"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/admin/contacts/analytics [get]
// GET /admin/contacts/analytics
func (h *ContactHandler) GetContactAnalytics(c *gin.Context) {
	analytics, err := h.service.GetContactAnalytics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get contact analytics",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ContactResponse{
		Success: true,
		Message: "Contact analytics retrieved successfully",
		Data:    analytics,
	})
}
