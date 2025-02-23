package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
)

// CommonHeaders stores the headers we want to capture from the request
type CommonHeaders struct {
	RequestID string
	UserAgent string
	// Add more headers as needed
}

// commonHeadersKey is used as a unique key for storing CommonHeaders in the context
type commonHeadersKey struct{}

// CommonHeadersMiddleware creates a gin middleware for capturing common headers
func CommonHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		headers := &CommonHeaders{
			RequestID: c.GetHeader("X-Request-ID"), // Custom header for request tracking
			UserAgent: c.GetHeader("User-Agent"),   // Browser or client identifier
		}

		// Store the headers in the context
		ctx := context.WithValue(c.Request.Context(), commonHeadersKey{}, headers)
		c.Request = c.Request.WithContext(ctx)

		// Proceed to next handler
		c.Next()
	}
}

// GetCommonHeadersFromContext retrieves the CommonHeaders from the context
func GetCommonHeadersFromContext(ctx context.Context) (*CommonHeaders, bool) {
	headers, ok := ctx.Value(commonHeadersKey{}).(*CommonHeaders)
	return headers, ok
}
