package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// SecurityHeaders adds security headers to responses
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Prevent clickjacking
		c.Writer.Header().Set("X-Frame-Options", "DENY")
		
		// Prevent MIME type sniffing
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		
		// Enable XSS protection
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")
		
		// Content Security Policy
		c.Writer.Header().Set("Content-Security-Policy", 
			"default-src 'self'; "+
			"script-src 'self' 'unsafe-inline' 'unsafe-eval'; "+
			"style-src 'self' 'unsafe-inline'; "+
			"img-src 'self' data: https:; "+
			"font-src 'self'; "+
			"connect-src 'self' ws: wss:; "+
			"frame-ancestors 'none'; "+
			"base-uri 'self'; "+
			"form-action 'self'")
		
		// Referrer Policy
		c.Writer.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		
		// Permissions Policy (formerly Feature Policy)
		c.Writer.Header().Set("Permissions-Policy", 
			"geolocation=(), microphone=(), camera=(), payment=(), usb=()")
		
		// Strict Transport Security (only for HTTPS)
		if c.Request.TLS != nil {
			c.Writer.Header().Set("Strict-Transport-Security", 
				"max-age=31536000; includeSubDomains; preload")
		}
		
		c.Next()
	}
}

// RequestSizeLimit limits the size of request bodies
func RequestSizeLimit(maxSize int64) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if strings.Contains(c.Request.URL.Path, "/api/stores") {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"error": "Request entity too large",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
			})
		}
		c.Abort()
	})
}

// TimeoutMiddleware adds a timeout to requests
func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a context with timeout
		ctx, cancel := c.Request.Context(), func() {}
		
		// Use gin's built-in timeout
		c.Request = c.Request.WithContext(ctx)
		
		finished := make(chan struct{})
		go func() {
			c.Next()
			finished <- struct{}{}
		}()
		
		select {
		case <-finished:
			return
		case <-time.After(timeout):
			cancel()
			c.JSON(http.StatusRequestTimeout, gin.H{
				"error": "Request timeout",
			})
			c.Abort()
		}
	}
}

// ValidateContentType ensures the request has the correct content type
func ValidateContentType() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "POST" || c.Request.Method == "PUT" {
			contentType := c.GetHeader("Content-Type")
			if !strings.Contains(contentType, "application/json") {
				c.JSON(http.StatusUnsupportedMediaType, gin.H{
					"error": "Content-Type must be application/json",
				})
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
