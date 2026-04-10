package middleware

import (
	"net/http"
	"portfolio/auth"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates the JWT token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")

		// Return an error if the header is missing
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required",
			})
			c.Abort() // Stop the request and do not continue to the next handlers
			return
		}

		// It comes in the format "Bearer token123456"
		// Remove the "Bearer " part and keep only the token
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		// Return an error if the token is empty
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token not found",
			})
			c.Abort()
			return
		}

		// Validate the token using the auth package function
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// If the token is valid, add user information to the context
		// This lets the next handlers access the current user
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)

		// Token is valid, continue to the next handler
		c.Next()
	}
}

// SuperAdminMiddleware - only for specific users
func SuperAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Normal auth must have already run
		username, exists := c.Get("username")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authentication required",
			})
			c.Abort()
			return
		}

		// Super admin check - only the "admin" user
		if username != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Super admin privileges required for user management",
			})
			c.Abort()
			return
		}

		// Continue if the user is a super admin
		c.Next()
	}
}

// GetCurrentUser gets user information from the context (helper function)
func GetCurrentUser(c *gin.Context) (userID int, username string, exists bool) {
	userIDInterface, exists1 := c.Get("user_id")
	usernameInterface, exists2 := c.Get("username")

	if !exists1 || !exists2 {
		return 0, "", false
	}

	userID, ok1 := userIDInterface.(int)
	username, ok2 := usernameInterface.(string)

	if !ok1 || !ok2 {
		return 0, "", false
	}

	return userID, username, true
}
