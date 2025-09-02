package middleware

import (
	"asa/pkg/jwtutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// contains checks if a slice of strings contains a specific string
func contains(list []string, val string) bool {
	for _, item := range list {
		if item == val {
			return true
		}
	}
	return false
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		DebugLog("🔐 === AUTH MIDDLEWARE START ===")
		DebugLog("🔐 Request: %s %s", c.Request.Method, c.Request.URL.Path)

		authHeader := c.GetHeader("Authorization")
		DebugLog("🔐 Authorization header: %s", authHeader)

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			DebugLog("❌ Missing or invalid Authorization header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Missing or invalid token"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		DebugLog("🔐 Token string length: %d", len(tokenString))
		DebugLog("🔐 Token preview: %s...", tokenString[:min(50, len(tokenString))])

		// Validate the token locally using the shared secret from JWT_SECRET
		DebugLog("🔐 Parsing JWT token...")
		claims, err := jwtutil.ParseToken(tokenString)
		if err != nil {
			DebugLog("❌ Failed to parse JWT token: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid or expired token"})
			return
		}
		DebugLog("✅ JWT token parsed successfully")
		DebugLog("🔐 All JWT claims: %+v", claims)

		// Extract roles from JWT claims - handle both 'role' (legacy) and 'roles' (local)
		var roles []string

		DebugLog("🔐 Extracting roles from JWT claims...")

		// First try to get 'roles' (plural array)
		if rolesInterface, exists := claims["roles"]; exists {
			DebugLog("🔐 Found 'roles' in claims: %+v (type: %T)", rolesInterface, rolesInterface)
			switch v := rolesInterface.(type) {
			case []string:
				roles = v
				DebugLog("✅ Roles extracted as []string: %v", roles)
			case []interface{}:
				for _, r := range v {
					if s, ok := r.(string); ok {
						roles = append(roles, s)
					}
				}
				DebugLog("✅ Roles extracted from []interface{}: %v", roles)
			default:
				DebugLog("⚠️ Unknown roles type: %T", rolesInterface)
			}
		} else {
			DebugLog("🔐 No 'roles' found in claims")
		}

		// If no roles found, try 'role' (singular from legacy auth)
		if len(roles) == 0 {
			DebugLog("🔐 No roles found, trying 'role' (singular)...")
			if roleInterface, exists := claims["role"]; exists {
				DebugLog("🔐 Found 'role' in claims: %+v (type: %T)", roleInterface, roleInterface)
				if role, ok := roleInterface.(string); ok {
					roles = []string{role}
					DebugLog("✅ Role extracted as string: %v", roles)
				} else {
					DebugLog("❌ Role is not a string: %T", roleInterface)
				}
			} else {
				DebugLog("🔐 No 'role' found in claims either")
			}
		}

		DebugLog("🔐 Final roles array: %v", roles)

		// Set all important JWT claims in context
		userID := claims["user_id"]
		username := claims["username"]
		email := claims["email"]
		name := claims["name"]

		DebugLog("🔐 Setting context values:")
		DebugLog("   user_id: %v", userID)
		DebugLog("   username: %v", username)
		DebugLog("   email: %v", email)
		DebugLog("   name: %v", name)
		DebugLog("   roles: %v", roles)

		c.Set("user_id", userID)
		c.Set("username", username)
		c.Set("email", email)
		c.Set("name", name)
		c.Set("roles", roles) // <-- roles is now a []string

		DebugLog("✅ === AUTH MIDDLEWARE COMPLETE ===")
		c.Next()
	}
}

// RequireRole creates a middleware that requires specific roles
func RequireRole(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		DebugLog("🔒 === ROLE CHECK START ===")
		DebugLog("🔒 Request: %s %s", c.Request.Method, c.Request.URL.Path)
		DebugLog("🔒 Required roles: %v", requiredRoles)

		rolesInterface, exists := c.Get("roles")
		if !exists {
			DebugLog("❌ No roles found in context")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "message": "No roles found in token"})
			return
		}

		userRoles, ok := rolesInterface.([]string)
		if !ok {
			DebugLog("❌ Invalid roles format in context: %T", rolesInterface)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid roles format in token"})
			return
		}

		DebugLog("🔒 User roles: %v", userRoles)

		// Check if user has any of the required roles
		hasRequiredRole := false
		for _, requiredRole := range requiredRoles {
			for _, userRole := range userRoles {
				if userRole == requiredRole {
					hasRequiredRole = true
					DebugLog("✅ User has required role: %s", requiredRole)
					break
				}
			}
			if hasRequiredRole {
				break
			}
		}

		if !hasRequiredRole {
			DebugLog("❌ User lacks required roles. User roles: %v, Required: %v", userRoles, requiredRoles)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Insufficient permissions. Required roles: " + strings.Join(requiredRoles, ", "),
			})
			return
		}

		DebugLog("✅ === ROLE CHECK PASSED ===")
		c.Next()
	}
}

// RequireAnyRole creates a middleware that requires at least one of the specified roles
func RequireAnyRole(requiredRoles ...string) gin.HandlerFunc {
	return RequireRole(requiredRoles...)
}

// RequireAllRoles creates a middleware that requires all specified roles
// Note: Since we now use single role per user, this function is equivalent to RequireRole
func RequireAllRoles(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		DebugLog("🔒 === ALL ROLES CHECK START ===")
		DebugLog("🔒 Request: %s %s", c.Request.Method, c.Request.URL.Path)
		DebugLog("🔒 Required roles (ALL): %v", requiredRoles)
		DebugLog("⚠️ Note: Single role system - checking if user has any of the required roles")

		rolesInterface, exists := c.Get("roles")
		if !exists {
			DebugLog("❌ No roles found in context")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "message": "No roles found in token"})
			return
		}

		userRoles, ok := rolesInterface.([]string)
		if !ok {
			DebugLog("❌ Invalid roles format in context: %T", rolesInterface)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid roles format in token"})
			return
		}

		DebugLog("🔒 User roles: %v", userRoles)

		// Check if user has ALL of the required roles
		hasAllRequiredRoles := true
		for _, requiredRole := range requiredRoles {
			hasThisRole := false
			for _, userRole := range userRoles {
				if userRole == requiredRole {
					hasThisRole = true
					break
				}
			}
			if !hasThisRole {
				hasAllRequiredRoles = false
				DebugLog("❌ User missing required role: %s", requiredRole)
				break
			}
		}

		if !hasAllRequiredRoles {
			DebugLog("❌ User lacks required roles. User roles: %v, Required: %v", userRoles, requiredRoles)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Insufficient permissions. Required roles: " + strings.Join(requiredRoles, ", "),
			})
			return
		}

		DebugLog("✅ === ALL ROLES CHECK PASSED ===")
		c.Next()
	}
}

// Helper function for min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
