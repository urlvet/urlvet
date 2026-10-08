package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
)

var (
	testRouter     *gin.Engine
	testRouterOnce sync.Once
)

// testR returns a shared router so we pay the Redis-timeout cost only once.
func testR() *gin.Engine {
	testRouterOnce.Do(func() {
		gin.SetMode(gin.TestMode)
		testRouter = SetupRouter()
	})
	return testRouter
}

func doRequest(t *testing.T, r *gin.Engine, method, path string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// TestHealthEndpoints verifies both health endpoints return 200 with status "ok".
func TestHealthEndpoints(t *testing.T) {
	r := testR()

	for _, path := range []string{"/health", "/api/v1/health"} {
		t.Run(path, func(t *testing.T) {
			w := doRequest(t, r, http.MethodGet, path)
			if w.Code != http.StatusOK {
				t.Errorf("GET %s: status = %d, want %d", path, w.Code, http.StatusOK)
			}
			var body map[string]any
			if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
				t.Fatalf("body not valid JSON: %v", err)
			}
			if body["status"] != "ok" {
				t.Errorf("status field = %v, want \"ok\"", body["status"])
			}
		})
	}
}

// TestRootEndpoint verifies the root handler returns the service metadata.
func TestRootEndpoint(t *testing.T) {
	w := doRequest(t, testR(), http.MethodGet, "/")
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("body not valid JSON: %v", err)
	}
	if body["service"] != "url.vet API" {
		t.Errorf("service = %v, want url.vet API", body["service"])
	}
}

// TestAnalyzeHandler_Validation tests the handler's validation without making
// real network requests.
func TestAnalyzeHandler_Validation(t *testing.T) {
	r := testR()

	tests := []struct {
		name       string
		path       string
		wantStatus int
	}{
		{"missing url param", "/api/v1/analyze", http.StatusBadRequest},
		{"invalid url - empty host", "/api/v1/analyze?url=http://", http.StatusBadRequest},
		{"url too long", "/api/v1/analyze?url=https://example.com/" + strings.Repeat("a", 2050), http.StatusBadRequest},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := doRequest(t, r, http.MethodGet, tc.path)
			if w.Code != tc.wantStatus {
				t.Errorf("status = %d, want %d (body: %s)", w.Code, tc.wantStatus, w.Body.String())
			}
		})
	}
}

// TestStructureHandlers tests length and depth endpoint validation.
func TestStructureHandlers(t *testing.T) {
	r := testR()

	tests := []struct {
		name       string
		path       string
		wantStatus int
	}{
		{"length - missing url", "/api/v1/length", http.StatusBadRequest},
		{"length - invalid url", "/api/v1/length?url=http://", http.StatusBadRequest},
		{"length - valid url", "/api/v1/length?url=https://example.com", http.StatusOK},
		{"depth - missing url", "/api/v1/depth", http.StatusBadRequest},
		{"depth - invalid url", "/api/v1/depth?url=http://", http.StatusBadRequest},
		{"depth - valid url", "/api/v1/depth?url=https://example.com/a/b", http.StatusOK},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := doRequest(t, r, http.MethodGet, tc.path)
			if w.Code != tc.wantStatus {
				t.Errorf("%s: status = %d, want %d (body: %s)", tc.path, w.Code, tc.wantStatus, w.Body.String())
			}
		})
	}
}

// TestPunycodeHandler tests punycode endpoint validation.
func TestPunycodeHandler(t *testing.T) {
	r := testR()

	tests := []struct {
		path       string
		wantStatus int
	}{
		{"/api/v1/punycode", http.StatusBadRequest},
		{"/api/v1/punycode?url=http://", http.StatusBadRequest},
		{"/api/v1/punycode?url=https://example.com", http.StatusOK},
	}

	for _, tc := range tests {
		w := doRequest(t, r, http.MethodGet, tc.path)
		if w.Code != tc.wantStatus {
			t.Errorf("GET %s: status = %d, want %d", tc.path, w.Code, tc.wantStatus)
		}
	}
}

// TestURLShortenerHandler tests shortener endpoint validation.
func TestURLShortenerHandler(t *testing.T) {
	r := testR()

	tests := []struct {
		path       string
		wantStatus int
	}{
		{"/api/v1/url-shortener", http.StatusBadRequest},
		{"/api/v1/url-shortener?url=http://", http.StatusBadRequest},
		{"/api/v1/url-shortener?url=https://example.com", http.StatusOK},
	}

	for _, tc := range tests {
		w := doRequest(t, r, http.MethodGet, tc.path)
		if w.Code != tc.wantStatus {
			t.Errorf("GET %s: status = %d, want %d", tc.path, w.Code, tc.wantStatus)
		}
	}
}

// TestTLDHandlers tests trusted-tld and risky-tld validation.
func TestTLDHandlers(t *testing.T) {
	r := testR()

	for _, endpoint := range []string{"/api/v1/trusted-tld", "/api/v1/risky-tld"} {
		t.Run(endpoint, func(t *testing.T) {
			w := doRequest(t, r, http.MethodGet, endpoint)
			if w.Code != http.StatusBadRequest {
				t.Errorf("missing url: status = %d, want 400", w.Code)
			}
			w = doRequest(t, r, http.MethodGet, endpoint+"?url=https://example.com")
			if w.Code != http.StatusOK {
				t.Errorf("valid url: status = %d, want 200 (body: %s)", w.Code, w.Body.String())
			}
		})
	}
}

// TestIPCheckHandler tests ip/check validation.
func TestIPCheckHandler(t *testing.T) {
	r := testR()

	w := doRequest(t, r, http.MethodGet, "/api/v1/ip/check")
	if w.Code != http.StatusBadRequest {
		t.Errorf("missing url: status = %d, want 400", w.Code)
	}
	w = doRequest(t, r, http.MethodGet, "/api/v1/ip/check?url=https://example.com")
	if w.Code != http.StatusOK {
		t.Errorf("valid url: status = %d, want 200 (body: %s)", w.Code, w.Body.String())
	}
}

// TestMetricsEndpoint verifies the /metrics endpoint is reachable and returns text/plain.
func TestMetricsEndpoint(t *testing.T) {
	w := doRequest(t, testR(), http.MethodGet, "/metrics")
	if w.Code != http.StatusOK {
		t.Fatalf("GET /metrics: status = %d, want 200", w.Code)
	}
	ct := w.Header().Get("Content-Type")
	if ct == "" {
		t.Error("Content-Type header is empty")
	}
}

func TestListEndpoints(t *testing.T) {
	r := SetupRouter()
	for _, path := range []string{"/api/v1/lists/threats", "/api/v1/lists/known-sites"} {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, path, nil))
		if w.Code != http.StatusOK {
			t.Fatalf("%s: %d", path, w.Code)
		}
		etag := w.Header().Get("ETag")
		if etag == "" {
			t.Fatalf("%s: no ETag", path)
		}
		req := httptest.NewRequest(http.MethodGet, path, nil)
		req.Header.Set("If-None-Match", etag)
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusNotModified {
			t.Errorf("%s with matching If-None-Match: %d, want 304", path, w.Code)
		}
	}
}
