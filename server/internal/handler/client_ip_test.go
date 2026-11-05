package handler

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

const (
	cfEdge    = "172.68.1.1"  // a Cloudflare edge address: public
	dockerGW  = "172.18.0.1"  // Docker bridge gateway: how a host proxy reaches us
	realIP    = "203.0.113.7" // the actual client
	spoofedIP = "9.9.9.9"     // what the client claims to be
)

// clientIP sends one request from peer with the given headers through an engine
// set up by configureClientIP, and returns what ctx.ClientIP() saw.
func clientIP(t *testing.T, peer string, headers map[string]string) string {
	t.Helper()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	if warn := configureClientIP(r); warn != nil {
		r.Use(warn)
	}
	var got string
	r.GET("/", func(c *gin.Context) { got = c.ClientIP() })

	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = peer + ":12345"
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	r.ServeHTTP(httptest.NewRecorder(), req)
	return got
}

func TestConfigureClientIP(t *testing.T) {
	xff := func(v string) map[string]string { return map[string]string{"X-Forwarded-For": v} }

	tests := []struct {
		name, platform, proxies string
		peer                    string
		headers                 map[string]string
		want                    string
	}{
		// Default: private peers are trusted, public ones aren't.
		{"default, exposed directly, spoof ignored", "", "", realIP, xff(spoofedIP), realIP},
		{"default, Caddy on host", "", "", dockerGW, xff(realIP), realIP},
		// nginx's $proxy_add_x_forwarded_for appends the peer to whatever the
		// client sent; Gin reads right to left and stops at the first
		// untrusted address, so the spoofed entry is never reached.
		{"default, nginx appends to spoofed header", "", "", dockerGW, xff(spoofedIP + ", " + realIP), realIP},
		{"default, public proxy not trusted", "", "", cfEdge, xff(realIP), cfEdge},

		{"none ignores headers", "", "none", dockerGW, xff(realIP), dockerGW},
		{"all trusts anyone (old behaviour)", "", "all", realIP, xff(spoofedIP), spoofedIP},
		{"custom range", "", "172.64.0.0/13", cfEdge, xff(realIP), realIP},
		{"custom single IP", "", "172.68.1.1", cfEdge, xff(realIP), realIP},
		{"invalid list fails closed", "", "not-an-ip", dockerGW, xff(realIP), dockerGW},

		{"cloudflare platform", "cloudflare", "", cfEdge,
			map[string]string{"CF-Connecting-IP": realIP, "X-Forwarded-For": spoofedIP + ", " + realIP}, realIP},
		{"custom platform header", "X-Client-IP", "", cfEdge,
			map[string]string{"X-Client-IP": realIP}, realIP},
		// A request that skips the platform falls back to the proxy rules, so
		// a public peer's X-Forwarded-For is still ignored.
		{"platform missing header falls back", "cloudflare", "", realIP, xff(spoofedIP), realIP},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("TRUSTED_PLATFORM", tt.platform)
			t.Setenv("TRUSTED_PROXIES", tt.proxies)
			if got := clientIP(t, tt.peer, tt.headers); got != tt.want {
				t.Errorf("ClientIP() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestConfigureClientIP_Warning(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for _, tt := range []struct {
		platform, proxies string
		want              bool
	}{
		{"", "", true},
		{"", "10.0.0.0/8", true},
		{"cloudflare", "", false},
		{"", "none", false},
		{"", "all", false},
		{"", "bogus", false},
	} {
		t.Setenv("TRUSTED_PLATFORM", tt.platform)
		t.Setenv("TRUSTED_PROXIES", tt.proxies)
		if got := configureClientIP(gin.New()) != nil; got != tt.want {
			t.Errorf("platform=%q proxies=%q: warning middleware = %v, want %v", tt.platform, tt.proxies, got, tt.want)
		}
	}
}
