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
