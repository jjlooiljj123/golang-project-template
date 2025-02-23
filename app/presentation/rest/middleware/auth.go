package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware creates a gin middleware for handling authentication
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header value
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Authorization header required"})
			return
		}

		// Split the header into parts, expecting "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid Authorization header format"})
			return
		}

		// Here, you would typically validate the token. For this example, we'll just check if the token is "valid"
		token := parts[1]
		if token != "valid" { // In real scenarios, you'd check against a token store or JWT validation
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token"})
			return
		}

		// If authentication is successful, proceed to the next handler
		c.Next()
	}
}
