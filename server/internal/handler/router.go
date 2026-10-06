package handler

import (
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/urlvet/urlvet/internal/handler/middleware"
)

// allowedOrigins reads CORS_ALLOWED_ORIGINS (comma-separated) from the environment.
func allowedOrigins() []string {
	raw := os.Getenv("CORS_ALLOWED_ORIGINS")
	if raw == "" {
		return []string{"*"}
	}
	origins := strings.Split(raw, ",")
	for i, o := range origins {
		origins[i] = strings.TrimSpace(o)
	}
	return origins
}

func SetupRouter() *gin.Engine {
	// Not gin.Default(): its logger records client IPs and full query strings
	// (the scanned URL), and its recovery dumps request headers.
	r := gin.New()
	r.Use(middleware.Recovery(), middleware.RequestLogger())

	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins(),
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	// Global Rate Limiter: 20 requests per minute per IP
	r.Use(middleware.RateLimiter(20, time.Minute))

	// Prometheus metrics — record latency and request counts for all routes
	r.Use(middleware.PrometheusMiddleware())

	// Prometheus scrape endpoint (not rate-limited, not behind /api/v1)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Swagger UI — served at /swagger/index.html
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// RootHandler returns basic info about the url.vet API service
	r.GET("/", RootHandler)

	// Unversioned global health check
	r.GET("/health", HealthHandler)

	v1 := r.Group("/api/v1")
	v1.Use(middleware.URLLengthValidator(2048))
	{
		v1.GET("/health", HealthHandler)

		if os.Getenv("ENV") == "DEV" {
			v1.GET("/test", TestApiHandler)
		}

		v1.GET("/analyze", AnalyzeURLHandler)
		v1.POST("/report", ReportResultHandler)
		v1.GET("/rank", GetDomainRankHandler)
		v1.GET("/ip/check", CheckIfUsesIPHandler)
		v1.GET("/ip/resolve", ResolveIPHandler)
		v1.GET("/length", CheckUrlLengthHandler)
		v1.GET("/depth", CheckUrlDepthHandler)
		v1.GET("/hsts", CheckHSTSHandler)
		v1.GET("/redirects", CheckRedirectsHandler)
		v1.GET("/punycode", CheckPunycodeHandler)
		v1.GET("/trusted-tld", CheckTrustedTLDHandler)
		v1.GET("/risky-tld", CheckRiskyTLDHandler)
		v1.GET("/url-shortener", CheckUrlShortenerHandler)
		v1.GET("/status-code", CheckStatusCodeHandler)
		v1.GET("/domain-info", DomainInfoHandler)

		// Lists the browser extension keeps locally
		v1.GET("/lists/threats", ThreatListHandler)
		v1.GET("/lists/known-sites", KnownSitesHandler)

		// Admin login — issues a signed session token
		v1.POST("/admin/login", AdminLoginHandler)

		// Admin — all routes require a valid Bearer token
		admin := v1.Group("/admin")
		admin.Use(middleware.BearerAuth())
		{
			admin.GET("/stats", AdminStatsHandler)
			admin.GET("/recent", AdminRecentHandler)
			admin.GET("/errors", AdminErrorsHandler)
			admin.GET("/reports", AdminReportsHandler)
			admin.DELETE("/reports/:id", AdminDeleteReportHandler)
			admin.GET("/cache", ListCacheHandler)
			admin.DELETE("/cache", FlushCacheHandler)
			admin.DELETE("/cache/*key", DeleteCacheKeyHandler)
		}

		v1.GET("/screenshot", ScreenshotHandler)

	}

	return r
}
