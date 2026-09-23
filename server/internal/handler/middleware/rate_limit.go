package middleware

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/urlvet/urlvet/internal/logger"
	"github.com/urlvet/urlvet/internal/service/cache"
)

// saltRotation is how often the secret behind hashed IPs is replaced.
const saltRotation = 24 * time.Hour

// ipHasher turns client IPs into opaque rate-limit keys, so no IP address is
// ever written to the cache (which persists to disk). The secret lives only in
// memory and is replaced daily, so a key can't be traced back to an IP.
type ipHasher struct {
	mu      sync.Mutex
	salt    []byte
	expires time.Time
	now     func() time.Time
}

func newIPHasher() *ipHasher { return &ipHasher{now: time.Now} }

func (h *ipHasher) key(ip string) string {
	h.mu.Lock()
	if h.salt == nil || !h.now().Before(h.expires) {
		h.salt = make([]byte, 32)
		if _, err := rand.Read(h.salt); err != nil {
			panic(err) // crypto/rand never fails on supported platforms
		}
		h.expires = h.now().Add(saltRotation)
	}
	mac := hmac.New(sha256.New, h.salt)
	h.mu.Unlock()

	mac.Write([]byte(ip))
	return "ratelimit:" + hex.EncodeToString(mac.Sum(nil)[:16])
}

// RateLimiter returns a Gin middleware that limits requests per client IP.
// IPs are hashed before use; see ipHasher.
func RateLimiter(limit int64, window time.Duration) gin.HandlerFunc {
	c, err := cache.New()
	if err != nil {
		logger.Warn("rate limiter cache unavailable, rate limiting disabled", "err", err)
		return func(ctx *gin.Context) { ctx.Next() }
	}

	hasher := newIPHasher()

	return func(ctx *gin.Context) {
		key := hasher.key(ctx.ClientIP())

		count, err := c.Increment(ctx, key)
		if err != nil {
			// On cache error, fail open to avoid blocking users
			ctx.Next()
			return
		}

		// If this is the first request in the window, set the expiration
		if count == 1 {
			_ = c.Expire(ctx, key, window)
		}

		// Set rate limit headers
		ctx.Header("X-RateLimit-Limit", strconv.FormatInt(limit, 10))
		ctx.Header("X-RateLimit-Remaining", strconv.FormatInt(max(0, limit-count), 10))
		ctx.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(window).Unix(), 10))

		if count > limit {
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":       "too many requests",
				"retry_after": window.Seconds(),
			})
			return
		}

		ctx.Next()
	}
}

// URLLengthValidator returns a middleware that rejects requests where the
// "url" query parameter exceeds maxLen bytes, preventing DoS via oversized inputs.
func URLLengthValidator(maxLen int) gin.HandlerFunc {
	return func(c *gin.Context) {
		if u := c.Query("url"); len(u) > maxLen {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("url exceeds maximum allowed length of %d characters", maxLen),
			})
			return
		}
		c.Next()
	}
}
