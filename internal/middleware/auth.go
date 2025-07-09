package middleware

import (
	"asa/pkg/jwtutil"
	"log"
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
		log.Printf("🔐 === AUTH MIDDLEWARE START ===")
		log.Printf("🔐 Request: %s %s", c.Request.Method, c.Request.URL.Path)

		authHeader := c.GetHeader("Authorization")
		log.Printf("🔐 Authorization header: %s", authHeader)

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			log.Printf("❌ Missing or invalid Authorization header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Missing or invalid token"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		log.Printf("🔐 Token string length: %d", len(tokenString))
		log.Printf("🔐 Token preview: %s...", tokenString[:min(50, len(tokenString))])

		// Validate the token locally using the shared secret from SECRET_KEY
		log.Printf("🔐 Parsing JWT token...")
		claims, err := jwtutil.ParseToken(tokenString)
		if err != nil {
			log.Printf("❌ Failed to parse JWT token: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid or expired token"})
			return
		}
		log.Printf("✅ JWT token parsed successfully")
		log.Printf("🔐 All JWT claims: %+v", claims)

		// Extract roles from JWT claims - handle both 'role' (AAA) and 'roles' (local)
		var roles []string

		log.Printf("🔐 Extracting roles from JWT claims...")

		// First try to get 'roles' (plural array)
		if rolesInterface, exists := claims["roles"]; exists {
			log.Printf("🔐 Found 'roles' in claims: %+v (type: %T)", rolesInterface, rolesInterface)
			switch v := rolesInterface.(type) {
			case []string:
				roles = v
				log.Printf("✅ Roles extracted as []string: %v", roles)
			case []interface{}:
				for _, r := range v {
					if s, ok := r.(string); ok {
						roles = append(roles, s)
					}
				}
				log.Printf("✅ Roles extracted from []interface{}: %v", roles)
			default:
				log.Printf("⚠️ Unknown roles type: %T", rolesInterface)
			}
		} else {
			log.Printf("🔐 No 'roles' found in claims")
		}

		// If no roles found, try 'role' (singular from AAA service)
		if len(roles) == 0 {
			log.Printf("🔐 No roles found, trying 'role' (singular)...")
			if roleInterface, exists := claims["role"]; exists {
				log.Printf("🔐 Found 'role' in claims: %+v (type: %T)", roleInterface, roleInterface)
				if role, ok := roleInterface.(string); ok {
					roles = []string{role}
					log.Printf("✅ Role extracted as string: %v", roles)
				} else {
					log.Printf("❌ Role is not a string: %T", roleInterface)
				}
			} else {
				log.Printf("🔐 No 'role' found in claims either")
			}
		}

		log.Printf("🔐 Final roles array: %v", roles)

		// Set all important JWT claims in context
		userID := claims["user_id"]
		username := claims["username"]
		email := claims["email"]
		name := claims["name"]

		log.Printf("🔐 Setting context values:")
		log.Printf("   user_id: %v", userID)
		log.Printf("   username: %v", username)
		log.Printf("   email: %v", email)
		log.Printf("   name: %v", name)
		log.Printf("   roles: %v", roles)

		c.Set("user_id", userID)
		c.Set("username", username)
		c.Set("email", email)
		c.Set("name", name)
		c.Set("roles", roles) // <-- roles is now a []string

		log.Printf("✅ === AUTH MIDDLEWARE COMPLETE ===")
		c.Next()
	}
}

// RequireRole creates a middleware that requires specific roles
func RequireRole(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("🔒 === ROLE CHECK START ===")
		log.Printf("🔒 Request: %s %s", c.Request.Method, c.Request.URL.Path)
		log.Printf("🔒 Required roles: %v", requiredRoles)

		rolesInterface, exists := c.Get("roles")
		if !exists {
			log.Printf("❌ No roles found in context")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "message": "No roles found in token"})
			return
		}

		roles, ok := rolesInterface.([]string)
		if !ok {
			log.Printf("❌ Invalid roles format in context: %T", rolesInterface)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid roles format in token"})
			return
		}

		log.Printf("🔒 User roles: %v", roles)

		// Check if user has any of the required roles
		hasRequiredRole := false
		for _, requiredRole := range requiredRoles {
			if contains(roles, requiredRole) {
				hasRequiredRole = true
				log.Printf("✅ User has required role: %s", requiredRole)
				break
			}
		}

		if !hasRequiredRole {
			log.Printf("❌ User lacks required roles. User roles: %v, Required: %v", roles, requiredRoles)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Insufficient permissions. Required roles: " + strings.Join(requiredRoles, ", "),
			})
			return
		}

		log.Printf("✅ === ROLE CHECK PASSED ===")
		c.Next()
	}
}

// RequireAnyRole creates a middleware that requires at least one of the specified roles
func RequireAnyRole(requiredRoles ...string) gin.HandlerFunc {
	return RequireRole(requiredRoles...)
}

// RequireAllRoles creates a middleware that requires all specified roles
func RequireAllRoles(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("🔒 === ALL ROLES CHECK START ===")
		log.Printf("🔒 Request: %s %s", c.Request.Method, c.Request.URL.Path)
		log.Printf("🔒 Required roles (ALL): %v", requiredRoles)

		rolesInterface, exists := c.Get("roles")
		if !exists {
			log.Printf("❌ No roles found in context")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "message": "No roles found in token"})
			return
		}

		roles, ok := rolesInterface.([]string)
		if !ok {
			log.Printf("❌ Invalid roles format in context: %T", rolesInterface)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid roles format in token"})
			return
		}

		log.Printf("🔒 User roles: %v", roles)

		// Check if user has ALL required roles
		for _, requiredRole := range requiredRoles {
			if !contains(roles, requiredRole) {
				log.Printf("❌ User missing required role: %s", requiredRole)
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"success": false,
					"message": "Insufficient permissions. Required roles: " + strings.Join(requiredRoles, ", "),
				})
				return
			}
			log.Printf("✅ User has required role: %s", requiredRole)
		}

		log.Printf("✅ === ALL ROLES CHECK PASSED ===")
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
