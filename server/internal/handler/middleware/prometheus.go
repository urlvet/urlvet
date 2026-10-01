package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/urlvet/urlvet/internal/metrics"
)

// PrometheusMiddleware records per-request latency and request counts.
func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		elapsed := time.Since(start).Seconds()

		path := c.FullPath() // matched route pattern, e.g. /api/v1/analyze
		if path == "" {
			path = c.Request.URL.Path // fallback for unmatched routes
		}
		status := strconv.Itoa(c.Writer.Status())
		method := c.Request.Method

		metrics.HTTPRequests.WithLabelValues(method, path, status).Inc()
		metrics.HTTPDuration.WithLabelValues(method, path).Observe(elapsed)
	}
}
