package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/urlvet/urlvet/internal/logger"
	"github.com/urlvet/urlvet/internal/service/cache"
)

// RateLimiter returns a Gin middleware that limits requests per IP
func RateLimiter(limit int64, window time.Duration) gin.HandlerFunc {
	c, err := cache.New()
	if err != nil {
		logger.Warn("rate limiter cache unavailable, rate limiting disabled", "err", err)
		return func(ctx *gin.Context) { ctx.Next() }
	}

	return func(ctx *gin.Context) {
		ip := ctx.ClientIP()
		key := fmt.Sprintf("ratelimit:%s", ip)

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
