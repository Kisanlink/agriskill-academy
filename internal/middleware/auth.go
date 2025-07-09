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
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Missing or invalid token"})
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate the token locally using the shared secret from SECRET_KEY
		claims, err := jwtutil.ParseToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid or expired token"})
			return
		}

		// Set all important JWT claims in context
		c.Set("user_id", claims["user_id"])
		c.Set("username", claims["username"])
		c.Set("email", claims["email"])
		c.Set("name", claims["name"])
		c.Set("roles", claims["roles"]) // <-- roles is now a []string
		c.Next()
	}
}
