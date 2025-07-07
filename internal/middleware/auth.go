// File: internal/middleware/auth.go

package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Missing or invalid token"})
			c.Abort()
			return
		}

		// Strip the 'Bearer ' prefix to get the token
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			secret = "secret" // Default secret for dev/test environments
		}

		// Parse the token
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			// Ensure the token is signed using the correct method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		})

		// Check for errors in token parsing
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Retrieve the claims and set the user context
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// Debug logging
			fmt.Printf("Token claims: %+v\n", claims)

			// Ensure correct type casting
			userID, ok := claims["user_id"].(string)
			if !ok {
				fmt.Printf("Failed to extract user_id from claims: %+v\n", claims)
				c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid user ID"})
				c.Abort()
				return
			}

			fmt.Printf("Extracted user_id: %s\n", userID)

			// Set values in the context
			c.Set("user_id", userID)
			c.Set("user_email", claims["email"])
			c.Set("user_name", claims["name"])
			c.Set("role", claims["role"])

			// Optionally, log the claims (useful for debugging)
			// fmt.Println("Token claims:", claims)
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid token claims"})
			c.Abort()
			return
		}

		// Proceed to the next handler
		c.Next()
	}
}
