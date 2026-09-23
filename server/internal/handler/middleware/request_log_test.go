package middleware

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRequestLoggerOmitsQueryAndIP(t *testing.T) {
	var buf bytes.Buffer
	prev := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, nil)))
	defer slog.SetDefault(prev)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(Recovery(), RequestLogger())
	r.GET("/api/v1/analyze", func(c *gin.Context) { c.Status(http.StatusOK) })
	r.GET("/boom", func(c *gin.Context) { panic("kaboom") })

	req := httptest.NewRequest(http.MethodGet, "/api/v1/analyze?url=https://bank.com/reset?token=SECRET", nil)
	req.RemoteAddr = "203.0.113.7:5555"
	r.ServeHTTP(httptest.NewRecorder(), req)

	boom := httptest.NewRequest(http.MethodGet, "/boom", nil)
	boom.Header.Set("Authorization", "Bearer ADMINTOKEN")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, boom)

	out := buf.String()
	for _, leak := range []string{"SECRET", "bank.com", "203.0.113.7", "ADMINTOKEN"} {
		if strings.Contains(out, leak) {
			t.Errorf("request log leaked %q:\n%s", leak, out)
		}
	}
	if !strings.Contains(out, "path=/api/v1/analyze") || !strings.Contains(out, "status=200") {
		t.Errorf("expected method/path/status in log:\n%s", out)
	}
	if w.Code != http.StatusInternalServerError || !strings.Contains(out, "kaboom") {
		t.Errorf("panic should become a logged 500, got %d:\n%s", w.Code, out)
	}
}
