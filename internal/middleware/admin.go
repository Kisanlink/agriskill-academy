package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdminAuthMiddleware ensures the user has admin role
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user roles from context (set by AuthMiddleware)
		rolesInterface, exists := c.Get("roles")
		if !exists {
			DebugLog("❌ AdminAuthMiddleware: No roles found in context")
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Unauthorized - No roles found",
			})
			c.Abort()
			return
		}

		// Type assert roles to []string
		roles, ok := rolesInterface.([]string)
		if !ok {
			DebugLog("❌ AdminAuthMiddleware: Invalid roles type in context")
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Unauthorized - Invalid roles type",
			})
			c.Abort()
			return
		}

		// Check if user has admin role
		hasAdminRole := false
		for _, role := range roles {
			if role == "asa_admin" {
				hasAdminRole = true
				break
			}
		}

		if !hasAdminRole {
			DebugLog("❌ AdminAuthMiddleware: User roles '%v' do not include admin", roles)
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Forbidden - Admin access required",
			})
			c.Abort()
			return
		}

		DebugLog("✅ AdminAuthMiddleware: Admin access granted for user roles '%v'", roles)
		c.Next()
	}
}
