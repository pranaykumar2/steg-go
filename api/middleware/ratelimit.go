package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// IPRateLimiter stores rate limiters for different IP addresses
type IPRateLimiter struct {
	ips    map[string]*rate.Limiter
	mu     *sync.RWMutex
	rate   rate.Limit
	burst  int
	expiry time.Duration
}

// NewIPRateLimiter creates a new rate limiter for IP addresses
func NewIPRateLimiter(r rate.Limit, b int, expiry time.Duration) *IPRateLimiter {
	return &IPRateLimiter{
		ips:    make(map[string]*rate.Limiter),
		mu:     &sync.RWMutex{},
		rate:   r,
		burst:  b,
		expiry: expiry,
	}
}

// AddIP creates a new rate limiter and adds it to the map
func (i *IPRateLimiter) AddIP(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter := rate.NewLimiter(i.rate, i.burst)
	i.ips[ip] = limiter

	// Schedule cleanup after expiry
	time.AfterFunc(i.expiry, func() {
		i.mu.Lock()
		delete(i.ips, ip)
		i.mu.Unlock()
	})

	return limiter
}

// GetLimiter returns the rate limiter for the specified IP address
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.RLock()
	limiter, exists := i.ips[ip]
	i.mu.RUnlock()

	if !exists {
		return i.AddIP(ip)
	}

	return limiter
}

// RateLimit middleware for limiting request rate per IP
func RateLimit(rps int, expiry time.Duration) gin.HandlerFunc {
	// Create a new rate limiter with the specified requests per second
	// and a burst of twice that amount
	limiter := NewIPRateLimiter(rate.Limit(rps)/60, rps, expiry)

	return func(c *gin.Context) {
		// Get client IP address
		ip := c.ClientIP()

		// Get the rate limiter for this IP
		ipLimiter := limiter.GetLimiter(ip)

		// Check if this request is allowed
		if !ipLimiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Please try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
