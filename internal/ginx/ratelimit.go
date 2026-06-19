package ginx

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type rateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
}

type visitor struct {
	mu      sync.Mutex
	tokens  int
	last    time.Time
	burst   int
	refill  time.Duration
}

func newRateLimiter() *rateLimiter {
	return &rateLimiter{visitors: make(map[string]*visitor)}
}

func (rl *rateLimiter) getVisitor(key string, rps int, burst int) *visitor {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, ok := rl.visitors[key]
	if !ok {
		v = &visitor{tokens: burst, last: time.Now(), burst: burst, refill: time.Second / time.Duration(rps)}
		rl.visitors[key] = v
	}
	return v
}

func (rl *rateLimiter) allow(key string, rps int, burst int) bool {
	v := rl.getVisitor(key, rps, burst)
	v.mu.Lock()
	defer v.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(v.last)
	v.last = now

	// Refill tokens based on elapsed time
	v.tokens += int(elapsed / v.refill)
	if v.tokens > v.burst {
		v.tokens = v.burst
	}

	if v.tokens > 0 {
		v.tokens--
		return true
	}
	return false
}

// GitRateLimit limits Git HTTP operations per IP.
// rps = requests per second, burst = max burst size.
func GitRateLimit(rps int, burst int) gin.HandlerFunc {
	rl := newRateLimiter()
	return func(c *gin.Context) {
		key := c.ClientIP()
		if !rl.allow(key, rps, burst) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, map[string]string{
				"message": "Rate limit exceeded. Try again later.",
			})
			return
		}
		c.Next()
	}
}
