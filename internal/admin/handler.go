package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	service AdminService
}

func NewAdminHandler(s AdminService) *AdminHandler {
	return &AdminHandler{s}
}

// @Summary Get Job Analytics
// @Description Get job-related analytics for admin dashboard
// @Tags Admin Analytics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} AnalyticsResponse "Job analytics retrieved successfully"
// @Failure 500 {object} AnalyticsResponse "Internal server error"
// @Router /api/admin/analytics/jobs [get]
// GET /admin/analytics/jobs
func (h *AdminHandler) GetJobAnalytics(c *gin.Context) {
	analytics, err := h.service.GetJobAnalytics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, AnalyticsResponse{
			Success: false,
			Message: "Failed to fetch job analytics: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, AnalyticsResponse{
		Success: true,
		Message: "Job analytics retrieved successfully",
		Data:    analytics,
	})
}

// @Summary Get User Analytics
// @Description Get user-related analytics for admin dashboard
// @Tags Admin Analytics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} AnalyticsResponse "User analytics retrieved successfully"
// @Failure 500 {object} AnalyticsResponse "Internal server error"
// @Router /api/admin/analytics/users [get]
// GET /admin/analytics/users
func (h *AdminHandler) GetUserAnalytics(c *gin.Context) {
	analytics, err := h.service.GetUserAnalytics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, AnalyticsResponse{
			Success: false,
			Message: "Failed to fetch user analytics: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, AnalyticsResponse{
		Success: true,
		Message: "User analytics retrieved successfully",
		Data:    analytics,
	})
}

// @Summary Get Application Analytics
// @Description Get application-related analytics for admin dashboard
// @Tags Admin Analytics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} AnalyticsResponse "Application analytics retrieved successfully"
// @Failure 500 {object} AnalyticsResponse "Internal server error"
// @Router /api/admin/analytics/applications [get]
// GET /admin/analytics/applications
func (h *AdminHandler) GetApplicationAnalytics(c *gin.Context) {
	analytics, err := h.service.GetApplicationAnalytics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, AnalyticsResponse{
			Success: false,
			Message: "Failed to fetch application analytics: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, AnalyticsResponse{
		Success: true,
		Message: "Application analytics retrieved successfully",
		Data:    analytics,
	})
}

// @Summary Get Dashboard Analytics
// @Description Get comprehensive dashboard analytics for admin
// @Tags Admin Analytics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} AnalyticsResponse "Dashboard analytics retrieved successfully"
// @Failure 500 {object} AnalyticsResponse "Internal server error"
// @Router /api/admin/analytics/dashboard [get]
// GET /admin/analytics/dashboard
func (h *AdminHandler) GetDashboardAnalytics(c *gin.Context) {
	analytics, err := h.service.GetDashboardAnalytics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, AnalyticsResponse{
			Success: false,
			Message: "Failed to fetch dashboard analytics: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, AnalyticsResponse{
		Success: true,
		Message: "Dashboard analytics retrieved successfully",
		Data:    analytics,
	})
}

// User Management Handlers

// @Summary Get Users
// @Description Get list of users with pagination and filtering
// @Tags Admin User Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Param role query string false "Filter by role"
// @Param search query string false "Search term"
// @Success 200 {object} UserResponse "Users retrieved successfully"
// @Failure 400 {object} UserResponse "Invalid query parameters"
// @Failure 500 {object} UserResponse "Internal server error"
// @Router /api/admin/users [get]
// GET /admin/users
func (h *AdminHandler) GetUsers(c *gin.Context) {
	var req UserListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, UserResponse{
			Success: false,
			Message: "Invalid query parameters: " + err.Error(),
		})
		return
	}

	users, err := h.service.GetUsers(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, UserResponse{
			Success: false,
			Message: "Failed to fetch users: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		Success: true,
		Message: "Users retrieved successfully",
		Data:    users,
	})
}

// GET /admin/users/:id
func (h *AdminHandler) GetUserByID(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, UserResponse{
			Success: false,
			Message: "User ID is required",
		})
		return
	}

	user, err := h.service.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, UserResponse{
			Success: false,
			Message: "User not found: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		Success: true,
		Message: "User retrieved successfully",
		Data:    user,
	})
}

// PUT /admin/users/:id
// @Summary Update User
// @Description Update a user's information by user ID
// @Tags Admin User Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param body body UpdateUserRequest true "User update payload"
// @Success 200 {object} UserResponse "User updated successfully"
// @Failure 400 {object} UserResponse "Invalid request body or missing user ID"
// @Failure 500 {object} UserResponse "Failed to update user"
// @Router /api/admin/users/{id} [put]
func (h *AdminHandler) UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, UserResponse{
			Success: false,
			Message: "User ID is required",
		})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, UserResponse{
			Success: false,
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	err := h.service.UpdateUser(userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, UserResponse{
			Success: false,
			Message: "Failed to update user: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		Success: true,
		Message: "User updated successfully",
	})
}

// DELETE /admin/users/:id
// @Summary Delete User
// @Description Delete a user by their ID
// @Tags Admin User Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} UserResponse "User deleted successfully"
// @Failure 400 {object} UserResponse "User ID is required"
// @Failure 500 {object} UserResponse "Failed to delete user"
// @Router /api/admin/users/{id} [delete]
func (h *AdminHandler) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, UserResponse{
			Success: false,
			Message: "User ID is required",
		})
		return
	}

	err := h.service.DeleteUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, UserResponse{
			Success: false,
			Message: "Failed to delete user: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		Success: true,
		Message: "User deleted successfully",
	})
}

// Company Management Handlers

// GET /admin/companies
// @Summary Get Companies
// @Description Retrieve a list of companies with optional filters and pagination
// @Tags Admin Company Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param name query string false "Company name filter"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} CompanyResponse "Companies retrieved successfully"
// @Failure 400 {object} CompanyResponse "Invalid query parameters"
// @Failure 500 {object} CompanyResponse "Failed to fetch companies"
// @Router /api/admin/companies [get]
func (h *AdminHandler) GetCompanies(c *gin.Context) {
	var req CompanyListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, CompanyResponse{
			Success: false,
			Message: "Invalid query parameters: " + err.Error(),
		})
		return
	}

	companies, err := h.service.GetCompanies(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, CompanyResponse{
			Success: false,
			Message: "Failed to fetch companies: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, CompanyResponse{
		Success: true,
		Message: "Companies retrieved successfully",
		Data:    companies,
	})
}

// GET /admin/companies/:id
// @Summary Get Company by ID
// @Description Retrieve a company by its unique ID
// @Tags Admin Company Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Company ID"
// @Success 200 {object} CompanyResponse "Company retrieved successfully"
// @Failure 400 {object} CompanyResponse "Company ID is required"
// @Failure 404 {object} CompanyResponse "Company not found"
// @Router /api/admin/companies/{id} [get]
func (h *AdminHandler) GetCompanyByID(c *gin.Context) {
	companyID := c.Param("id")
	if companyID == "" {
		c.JSON(http.StatusBadRequest, CompanyResponse{
			Success: false,
			Message: "Company ID is required",
		})
		return
	}

	company, err := h.service.GetCompanyByID(companyID)
	if err != nil {
		c.JSON(http.StatusNotFound, CompanyResponse{
			Success: false,
			Message: "Company not found: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, CompanyResponse{
		Success: true,
		Message: "Company retrieved successfully",
		Data:    company,
	})
}

// PUT /admin/companies/:id
// @Summary Update Company
// @Description Update an existing company by its unique ID
// @Tags Admin Company Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Company ID"
// @Param body body UpdateCompanyRequest true "Company update data"
// @Success 200 {object} CompanyResponse "Company updated successfully"
// @Failure 400 {object} CompanyResponse "Invalid request or Company ID is required"
// @Failure 404 {object} CompanyResponse "Company not found"
// @Router /api/admin/companies/{id} [put]
func (h *AdminHandler) UpdateCompany(c *gin.Context) {
	companyID := c.Param("id")
	if companyID == "" {
		c.JSON(http.StatusBadRequest, CompanyResponse{
			Success: false,
			Message: "Company ID is required",
		})
		return
	}

	var req UpdateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, CompanyResponse{
			Success: false,
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	err := h.service.UpdateCompany(companyID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, CompanyResponse{
			Success: false,
			Message: "Failed to update company: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, CompanyResponse{
		Success: true,
		Message: "Company updated successfully",
	})
}

// DELETE /admin/companies/:id
// @Summary Delete Company
// @Description Delete a company by its unique ID
// @Tags Admin Company Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Company ID"
// @Success 200 {object} CompanyResponse "Company deleted successfully"
// @Failure 400 {object} CompanyResponse "Company ID is required"
// @Failure 404 {object} CompanyResponse "Company not found"
// @Failure 500 {object} CompanyResponse "Failed to delete company"
// @Router /api/admin/companies/{id} [delete]
func (h *AdminHandler) DeleteCompany(c *gin.Context) {
	companyID := c.Param("id")
	if companyID == "" {
		c.JSON(http.StatusBadRequest, CompanyResponse{
			Success: false,
			Message: "Company ID is required",
		})
		return
	}

	err := h.service.DeleteCompany(companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, CompanyResponse{
			Success: false,
			Message: "Failed to delete company: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, CompanyResponse{
		Success: true,
		Message: "Company deleted successfully",
	})
}

// GET /admin/analytics/companies
// @Summary Get Company Analytics
// @Description Retrieve analytics for all companies
// @Tags Admin Company Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} AnalyticsResponse "Company analytics retrieved successfully"
// @Failure 500 {object} AnalyticsResponse "Failed to fetch company analytics"
// @Router /api/admin/analytics/companies [get]
func (h *AdminHandler) GetCompanyAnalytics(c *gin.Context) {
	analytics, err := h.service.GetCompanyAnalytics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, AnalyticsResponse{
			Success: false,
			Message: "Failed to fetch company analytics: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, AnalyticsResponse{
		Success: true,
		Message: "Company analytics retrieved successfully",
		Data:    analytics,
	})
}
