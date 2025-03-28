package middleware

import (
	"github.com/gin-gonic/gin"
)

// SecurityHeaders adds security-related HTTP headers to all responses
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Prevent MIME type sniffing
		c.Header("X-Content-Type-Options", "nosniff")

		// Protection against clickjacking
		c.Header("X-Frame-Options", "DENY")

		// Enable XSS filter in browser
		c.Header("X-XSS-Protection", "1; mode=block")

		// Control what features and APIs can be used in the browser
		c.Header("Feature-Policy", "camera 'none'; microphone 'none'")

		// Content Security Policy - More permissive for web UI with blob support
		c.Header("Content-Security-Policy",
            "default-src 'self'; " +
            "script-src 'self' 'unsafe-inline' 'unsafe-eval'; " +
            "style-src 'self' 'unsafe-inline' https://fonts.googleapis.com https://cdnjs.cloudflare.com; " +
            "font-src 'self' https://fonts.gstatic.com https://cdnjs.cloudflare.com; " +
            "img-src 'self' data: blob:; " + // Added 'blob:' here
            "connect-src 'self'")

		c.Next()
	}
}
