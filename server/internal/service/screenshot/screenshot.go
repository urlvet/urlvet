package screenshot

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/urlvet/urlvet/internal/logger"
)

// Service manages screenshot operations with a shared browser allocator
type Service struct {
	allocCtx    context.Context
	allocCancel context.CancelFunc
	mu          sync.RWMutex
	initialized bool
}

var (
	// Default service instance (singleton pattern)
	defaultService *Service
	serviceOnce    sync.Once
)

// NewService creates a new screenshot service with a shared browser allocator
func NewService(chromeURL string) (*Service, error) {
	var allocCtx context.Context
	var allocCancel context.CancelFunc

	if chromeURL != "" {
		// Use remote allocator for enhanced sandboxing (e.g., separate container)
		logger.Info("initializing remote screenshot service", "url", chromeURL)
		allocCtx, allocCancel = chromedp.NewRemoteAllocator(context.Background(), chromeURL)
	} else {
		// Fallback to local allocator
		logger.Info("initializing local screenshot service")
		opts := append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.NoSandbox,
			chromedp.Headless,
			chromedp.DisableGPU,
			chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
		)
		allocCtx, allocCancel = chromedp.NewExecAllocator(context.Background(), opts...)
	}

	service := &Service{
		allocCtx:    allocCtx,
		allocCancel: allocCancel,
		initialized: true,
	}

	return service, nil
}

// GetService returns the default service instance (singleton)
func GetService() (*Service, error) {
	var err error
	serviceOnce.Do(func() {
		chromeURL := os.Getenv("CHROME_URL")
		defaultService, err = NewService(chromeURL)
	})
	return defaultService, err
}

// Close shuts down the browser allocator
func (s *Service) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.allocCancel != nil {
		s.allocCancel()
	}
	return nil
}

// normalizeURL normalizes a URL for consistent cache keys
// - Lowercase host
// - Remove fragment
// - Sort query params
// - Trim trailing slash (except root)
func normalizeURL(rawURL string) (string, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	// Enforce scheme allowlist (security)
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return "", errors.New("invalid or unsafe URL scheme")
	}

	// Normalize host (lowercase)
	parsed.Host = strings.ToLower(parsed.Host)

	// Remove fragment
	parsed.Fragment = ""

	// Normalize path (trim trailing slash except root)
	if parsed.Path != "/" && strings.HasSuffix(parsed.Path, "/") {
		parsed.Path = strings.TrimSuffix(parsed.Path, "/")
	}

	// Sort query parameters for consistency
	query := parsed.Query()
	if len(query) > 0 {
		parsed.RawQuery = ""
		keys := make([]string, 0, len(query))
		for k := range query {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		values := make([]string, 0, len(keys))
		for _, k := range keys {
			vals := query[k]
			sort.Strings(vals)
			for _, v := range vals {
				values = append(values, fmt.Sprintf("%s=%s", k, v))
			}
		}
		parsed.RawQuery = strings.Join(values, "&")
	}

	return parsed.String(), nil
}

// sanitizeURL creates a filesystem-safe filename from a URL
func sanitizeURL(rawURL string) string {
	normalized, err := normalizeURL(rawURL)
	if err != nil {
		// Fallback to basic sanitization
		normalized = rawURL
	}

	re := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	safe := re.ReplaceAllString(normalized, "-")
	safe = strings.Trim(safe, "-")
	if len(safe) > 100 {
		safe = safe[:100]
	}
	return safe
}

// GetScreenshotPath returns the file path for a screenshot based on normalized URL
func GetScreenshotPath(url string) string {
	urlStr := sanitizeURL(url)
	filename := "screenshot-" + urlStr + ".png"
	dir := filepath.Join(".", "tmp", "screenshots")
	return filepath.Join(dir, filename)
}

// validateURL ensures the URL is safe to navigate to
func validateURL(rawURL string) (string, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", errors.New("invalid URL format")
	}

	// Enforce scheme allowlist (security critical)
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return "", errors.New("only http and https schemes are allowed")
	}

	// Block dangerous schemes explicitly
	dangerousSchemes := []string{"file", "chrome", "data", "javascript", "about"}
	for _, scheme := range dangerousSchemes {
		if parsed.Scheme == scheme {
			return "", fmt.Errorf("unsafe URL scheme: %s", scheme)
		}
	}

	return parsed.String(), nil
}

// TakeScreenshot takes a screenshot and returns the image bytes
// It first checks if a cached screenshot exists, otherwise takes a new one
func (s *Service) TakeScreenshot(rawURL string) ([]byte, error) {
	// Validate and normalize URL
	validatedURL, err := validateURL(rawURL)
	if err != nil {
		return nil, fmt.Errorf("URL validation failed: %w", err)
	}

	// Check if cached screenshot exists
	filePath := GetScreenshotPath(validatedURL)
	if fileInfo, err := os.Stat(filePath); err == nil && !fileInfo.IsDir() {
		// File exists, read and return it
		imageBytes, err := os.ReadFile(filePath)
		if err == nil {
			return imageBytes, nil
		}
		// If read fails, continue to take new screenshot
	}

	// Ensure service is initialized
	s.mu.RLock()
	if !s.initialized {
		s.mu.RUnlock()
		return nil, errors.New("screenshot service not initialized")
	}
	allocCtx := s.allocCtx
	s.mu.RUnlock()

	// Create a new context for this screenshot operation
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// Timeout to prevent hanging (increased for slow pages)
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var buf []byte

	// Configure viewport and anti-bot measures
	err = chromedp.Run(ctx,
		// Set desktop viewport for consistent screenshots
		chromedp.EmulateViewport(1366, 768),
		// Anti-bot hardening
		chromedp.ActionFunc(func(ctx context.Context) error {
			// Remove webdriver flag
			return chromedp.Evaluate(`Object.defineProperty(navigator, 'webdriver', {get: () => undefined})`, nil).Do(ctx)
		}),
		// Navigate to URL
		chromedp.Navigate(validatedURL),
		// Wait for page to be ready (better than Sleep)
		chromedp.WaitReady("body", chromedp.ByQuery),
		// Wait for network idle (more reliable for JS-heavy pages)
		chromedp.ActionFunc(func(ctx context.Context) error {
			// Wait a bit for dynamic content, but with timeout
			select {
			case <-time.After(2 * time.Second):
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		}),
		// Take screenshot with quality 90, but limit height to prevent memory issues
		// Using viewport screenshot instead of full page for memory safety
		chromedp.CaptureScreenshot(&buf),
	)
	if err != nil {
		return nil, fmt.Errorf("screenshot capture failed: %w", err)
	}

	// Save screenshot to cache for future use
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		// If we can't create directory, still return the screenshot
		return buf, nil
	}

	// Try to save, but don't fail if save fails
	if err := os.WriteFile(filePath, buf, 0644); err != nil {
		logger.Warn("failed to cache screenshot", "err", err)
	}

	return buf, nil
}

// TakeScreenshot is a convenience function that uses the default service
func TakeScreenshot(rawURL string) ([]byte, error) {
	service, err := GetService()
	if err != nil {
		return nil, fmt.Errorf("failed to get screenshot service: %w", err)
	}
	return service.TakeScreenshot(rawURL)
}

// TakeScreenshotAndSave takes a screenshot, saves it to disk, and returns the file path
// This is kept for backward compatibility if needed elsewhere
func TakeScreenshotAndSave(url string) string {
	buf, err := TakeScreenshot(url)
	if err != nil {
		logger.Error("screenshot failed", "err", err)
		return ""
	}

	timestamp := time.Now().Format("20060102-150405")
	urlStr := sanitizeURL(url)
	filename := "screenshot-" + timestamp + "-" + urlStr + ".png"
	dir := filepath.Join(".", "tmp", "screenshots") // ./server/tmp/screenshots

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		logger.Error("failed to create screenshots directory", "err", err)
		return ""
	}
	fullPath := filepath.Join(dir, filename)

	if err := os.WriteFile(fullPath, buf, 0644); err != nil {
		logger.Error("failed to write screenshot file", "err", err)
		return ""
	}

	fmt.Println("Screenshot saved")
	return fullPath
}
