package middleware

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"boilerplate/app/infrastructure/config"
)

// TimeoutMiddleware creates a gin middleware for setting request timeouts
func TimeoutMiddleware(cfg *config.AppConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a context with timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), cfg.HandlerTimeout)
		defer cancel()

		// Replace the request's context with the new timeout context
		c.Request = c.Request.WithContext(ctx)

		// Use a channel to manage the flow of the request
		finish := make(chan struct{}, 1)
		panicChan := make(chan interface{}, 1)

		// Start a goroutine to handle the request
		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()
			c.Next()
			finish <- struct{}{}
		}()

		// Wait for either the request to finish or the context to timeout
		select {
		case <-ctx.Done():
			// If context timeout, abort with an error
			c.AbortWithStatusJSON(http.StatusRequestTimeout, gin.H{"error": "Request timeout"})
		case p := <-panicChan:
			// If a panic occurred, handle it
			panic(p)
		case <-finish:
			// If the request finished before timeout, do nothing
		}
	}
}
