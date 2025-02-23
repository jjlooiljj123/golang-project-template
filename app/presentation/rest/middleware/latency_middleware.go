package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// LatencyLogger logs the time taken for each request to complete
func LatencyLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Stop timer
		latency := time.Since(start)
		path := c.Request.URL.Path
		clientIP := c.ClientIP()

		// Log latency
		fmt.Printf("[%s] %s - Latency: %v\n", clientIP, path, latency)
	}
}
