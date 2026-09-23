package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/urlvet/urlvet/internal/logger"
)

// RequestLogger logs one line per request: method, path, status and duration.
// It deliberately leaves out the query string (it holds the scanned URL) and the
// client IP. Health checks and metrics scrapes are logged at debug level only.
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		path := c.Request.URL.Path
		attrs := []any{
			"method", c.Request.Method,
			"path", path,
			"status", c.Writer.Status(),
			"duration", time.Since(start).Round(time.Millisecond).String(),
		}
		switch {
		case path == "/health" || path == "/api/v1/health" || path == "/metrics":
			logger.Debug("request", attrs...)
		case c.Writer.Status() >= 500:
			logger.Error("request", attrs...)
		default:
			logger.Info("request", attrs...)
		}
	}
}

// Recovery turns a panic into a 500 and logs it without dumping the request.
// (gin.Recovery prints the full request, including headers like Authorization.)
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				logger.Error("panic", "path", c.Request.URL.Path, "err", fmt.Sprint(rec))
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
