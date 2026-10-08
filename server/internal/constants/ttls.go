package constants

import "time"

// Cache TTL constants for analyzer tasks
// These values balance freshness with performance for cached network operations

// DomainRankTTL - Domain rank doesn't change frequently
const DomainRankTTL = 24 * time.Hour

// IPResolutionTTL - DNS IP addresses can change but not frequently
const IPResolutionTTL = 3 * time.Hour

// DNSValidityTTL - DNS records (NS/MX) can change occasionally
const DNSValidityTTL = 3 * time.Hour

// WHOISLookupTTL - Domain registration info is relatively stable
const WHOISLookupTTL = 24 * time.Hour

// HTTPCombinedTTL - HTTP responses (redirects, status codes, HSTS) can change more frequently
const HTTPCombinedTTL = 3 * time.Hour

// TLSCombinedTTL - TLS/SSL certificates don't change frequently
const TLSCombinedTTL = 24 * time.Hour

// PhishTankTTL - Threat intel data should be reasonably fresh
const PhishTankTTL = 3 * time.Hour

// ContentAnalysisTTL - Page content can change, but usually stays consistent for a while
const ContentAnalysisTTL = 3 * time.Hour

// ScreenshotTTL - Page screenshots on disk are reused, then deleted, after this long
const ScreenshotTTL = 24 * time.Hour

// WebRiskTTL - Google Web Risk verdicts; Google's lists change by the hour
const WebRiskTTL = time.Hour

// SafeBrowsingTTL - Google Safe Browsing verdicts
const SafeBrowsingTTL = time.Hour

// AnalyzeResultTTL - Full scan result cache: avoids re-running all 17 tasks for the same URL
const AnalyzeResultTTL = 24 * time.Hour
