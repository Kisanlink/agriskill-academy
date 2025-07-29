package contact

import (
	"asa/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterPublicRoutes registers public contact endpoints (no auth required)
func RegisterPublicRoutes(rg *gin.RouterGroup, handler *ContactHandler) {
	// Public contact form submission endpoint
	rg.POST("/contact", handler.SubmitContactForm)
}

// RegisterAdminRoutes registers admin contact endpoints (admin auth required)
func RegisterAdminRoutes(rg *gin.RouterGroup, handler *ContactHandler) {
	// Admin contact management routes
	admin := rg.Group("/admin")
	admin.Use(middleware.RequireRole("asa_admin"))
	{
		contacts := admin.Group("/contacts")
		{
			contacts.GET("", handler.GetContactRequests)             // GET /admin/contacts
			contacts.GET("/analytics", handler.GetContactAnalytics)  // GET /admin/contacts/analytics
			contacts.GET("/:id", handler.GetContactByID)             // GET /admin/contacts/:id
			contacts.PUT("/:id/status", handler.UpdateContactStatus) // PUT /admin/contacts/:id/status
			contacts.DELETE("/:id", handler.DeleteContact)           // DELETE /admin/contacts/:id
		}
	}
}
