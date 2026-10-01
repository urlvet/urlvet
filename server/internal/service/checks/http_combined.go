package checks

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

// CombinedHTTPResult contains all HTTP-related check results
type CombinedHTTPResult struct {
	RedirectionResult RedirectionResult
	StatusCode        int
	StatusText        string
	StatusSuccess     bool
	StatusIsRedirect  bool
	SupportsHSTS      bool
}

// CheckHTTPCombined performs a single HTTP request to gather redirect chain,
// status code, and HSTS header information. Uses HEAD first, falls back to GET.
func CheckHTTPCombined(rawURL string) (CombinedHTTPResult, error) {
	result := CombinedHTTPResult{}
	var redirects []string

	// Parse URL - use original URL for redirects
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return result, fmt.Errorf("invalid URL: %v", err)
	}

	// Ensure scheme exists
	checkURL := rawURL
	if parsedURL.Scheme == "" {
		checkURL = "https://" + rawURL
		parsedURL, _ = url.Parse(checkURL)
	}

	// Create HTTP client with aggressive timeouts and redirect tracking
	// Per-request timeouts are aggressive to avoid slow servers
	// Overall timeout is 5s as requested, but individual operations are faster
	transport := &http.Transport{
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		ResponseHeaderTimeout: 800 * time.Millisecond, // Max time to wait for response headers per request
		ExpectContinueTimeout: 500 * time.Millisecond,
		IdleConnTimeout:       90 * time.Second,
		DisableKeepAlives:     false, // Keep connections alive for redirects
		MaxIdleConns:          10,
		MaxIdleConnsPerHost:   2,
		DialContext: ssrfSafeDialContext(&net.Dialer{
			Timeout:   300 * time.Millisecond, // Connection timeout
			KeepAlive: 30 * time.Second,
		}),
		TLSHandshakeTimeout: 300 * time.Millisecond, // TLS handshake timeout
	}

	// Keep cookies across hops, like a browser: some sites set one and redirect,
	// and without it they bounce between the same pages. The jar lives only for
	// this scan.
	jar, _ := cookiejar.New(nil)

	client := &http.Client{
		Timeout:   5 * time.Second, // Overall timeout (shouldn't be reached with header timeout)
		Transport: transport,
		Jar:       jar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			redirects = append(redirects, req.URL.String())
			if len(via) >= 10 { // Max 10 redirects
				return errors.New("stopped after 10 redirects")
			}
			return nil
		},
	}

	// Create context with timeout for request cancellation
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Try HEAD first
	req, err := http.NewRequestWithContext(ctx, "HEAD", checkURL, nil)
	if err != nil {
		return result, fmt.Errorf("failed to create HEAD request: %v", err)
	}

	var usedGET bool
	resp, err := client.Do(req)
	if err != nil {
		// Fallback to GET if HEAD fails. Start the chain over, so hops from the
		// failed HEAD attempt aren't counted twice.
		usedGET = true
		redirects = nil
		req, err = http.NewRequestWithContext(ctx, "GET", checkURL, nil)
		if err != nil {
			return result, fmt.Errorf("failed to create GET request: %v", err)
		}
		resp, err = client.Do(req)
		if err != nil {
			return result, fmt.Errorf("HTTP request failed: %v", err)
		}
	}
	defer resp.Body.Close()

	// For GET requests, immediately discard body - we only need headers
	// This prevents waiting for slow servers to send body data
	// HEAD requests shouldn't have a body, but be safe
	if usedGET && resp.Body != nil {
		io.Copy(io.Discard, resp.Body)
	}

	// Build redirect chain
	chain := append([]string{rawURL}, redirects...)
	finalURL := chain[len(chain)-1]
	finalURLHost, _ := GetHost(finalURL)

	// Detect domain jumps
	origDomain, _ := GetDomain(rawURL)
	hasJump := false
	for _, u := range chain[1:] {
		urlDomain, _ := GetDomain(u)
		if !sameSite(urlDomain, origDomain) {
			hasJump = true
			break
		}
	}

	result.RedirectionResult = RedirectionResult{
		IsRedirected:  len(redirects) > 0,
		Chain:         chain,
		FinalURL:      finalURL,
		FinalURLHost:  finalURLHost,
		ChainLength:   len(chain),
		HasDomainJump: hasJump,
	}

	// Extract status code information
	result.StatusCode = resp.StatusCode
	result.StatusText = http.StatusText(resp.StatusCode)
	result.StatusSuccess = resp.StatusCode >= 200 && resp.StatusCode < 300
	result.StatusIsRedirect = resp.StatusCode >= 300 && resp.StatusCode < 400

	// Check for HSTS header (only available over HTTPS)
	finalParsedURL, _ := url.Parse(finalURL)
	if finalParsedURL != nil && finalParsedURL.Scheme == "https" {
		_, result.SupportsHSTS = resp.Header["Strict-Transport-Security"]
	} else {
		// If final URL is not HTTPS, check HSTS by making a separate HTTPS request
		// Use a separate client with very short timeout to avoid adding significant latency
		// (matching behavior of original SupportsHSTS function)
		domain, _ := GetDomain(rawURL)
		if domain != "" {
			httpsURL := "https://" + domain
			hstsCtx, hstsCancel := context.WithTimeout(context.Background(), 800*time.Millisecond)
			defer hstsCancel()

			// Use same transport but with context timeout
			hstsReq, err := http.NewRequestWithContext(hstsCtx, "HEAD", httpsURL, nil)
			if err == nil {
				hstsResp, err := client.Do(hstsReq)
				if err == nil {
					_, result.SupportsHSTS = hstsResp.Header["Strict-Transport-Security"]
					if hstsResp.Body != nil {
						io.Copy(io.Discard, hstsResp.Body) // Discard body
					}
					hstsResp.Body.Close()
				}
			}
		}
	}

	return result, nil
}
