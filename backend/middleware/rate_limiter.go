package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter implements a simple token bucket rate limiter
type RateLimiter struct {
	clients map[string]*ClientLimiter
	mutex   sync.RWMutex
	rate    int           // requests per minute
	burst   int           // maximum burst size
}

type ClientLimiter struct {
	tokens    int
	lastSeen  time.Time
	mutex     sync.Mutex
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(rate, burst int) *RateLimiter {
	return &RateLimiter{
		clients: make(map[string]*ClientLimiter),
		rate:    rate,
		burst:   burst,
	}
}

// Middleware returns a Gin middleware for rate limiting
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		
		// Get or create client limiter
		rl.mutex.RLock()
		limiter, exists := rl.clients[clientIP]
		rl.mutex.RUnlock()
		
		if !exists {
			rl.mutex.Lock()
			// Double-check after acquiring write lock
			limiter, exists = rl.clients[clientIP]
			if !exists {
				limiter = &ClientLimiter{
					tokens:   rl.burst,
					lastSeen: time.Now(),
				}
				rl.clients[clientIP] = limiter
			}
			rl.mutex.Unlock()
		}
		
		// Check if request is allowed
		limiter.mutex.Lock()
		
		// Add tokens based on time elapsed
		now := time.Now()
		elapsed := now.Sub(limiter.lastSeen)
		tokensToAdd := int(elapsed.Minutes() * float64(rl.rate))
		
		if tokensToAdd > 0 {
			limiter.tokens += tokensToAdd
			if limiter.tokens > rl.burst {
				limiter.tokens = rl.burst
			}
		}
		
		limiter.lastSeen = now
		
		if limiter.tokens <= 0 {
			limiter.mutex.Unlock()
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Please try again later.",
			})
			c.Abort()
			return
		}
		
		limiter.tokens--
		limiter.mutex.Unlock()
		
		c.Next()
	}
}

// CleanupExpiredClients removes clients that haven't been seen for a while
func (rl *RateLimiter) CleanupExpiredClients() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	
	for range ticker.C {
		rl.mutex.Lock()
		now := time.Now()
		for ip, limiter := range rl.clients {
			limiter.mutex.Lock()
			if now.Sub(limiter.lastSeen) > 10*time.Minute {
				delete(rl.clients, ip)
			}
			limiter.mutex.Unlock()
		}
		rl.mutex.Unlock()
	}
}
