package analyzer

import (
	"context"
	"strings"

	"github.com/redis/go-redis/v9"
	"github.com/urlvet/urlvet/internal/constants"
	"github.com/urlvet/urlvet/internal/service/checks"
	"github.com/urlvet/urlvet/internal/service/domaininfo"
	"github.com/urlvet/urlvet/internal/service/rank"
	"github.com/urlvet/urlvet/internal/service/threatfeeds"
	"github.com/urlvet/urlvet/internal/service/typosquat"
	"golang.org/x/net/publicsuffix"
)

// Rank
type rankTask struct{}

func (rankTask) Name() string { return "domain_rank" }
func (rankTask) Run(in *Input, out *Output) error {
	_, err := cachedTask(
		context.Background(),
		in.Cache,
		"domain_rank:"+in.Domain,
		constants.DomainRankTTL,
		func() (int, error) { return rank.DomainRankLookup(in.Domain), nil },
		func(o *Output, r int) { o.Rank = r },
		out,
	)
	return err
}

// HSTS
type hstsTask struct{}

func (hstsTask) Name() string { return "hsts_check" }
func (hstsTask) Run(in *Input, out *Output) error {
	h, err := checks.SupportsHSTS(in.URL)
	if err != nil {
		return err
	}
	updateOutput(func(o *Output) { o.SupportsHSTS = h })(out)
	return nil
}

// IP checks
type ipCheckTask struct{}

func (ipCheckTask) Name() string { return "ip_check" }
func (ipCheckTask) Run(in *Input, out *Output) error {
	b, err := checks.UsesIPInsteadOfDomain(in.URL)
	if err != nil {
		return err
	}
	updateOutput(func(o *Output) { o.URLUsesIP = b })(out)
	return nil
}

type ipResolveTask struct{}

func (ipResolveTask) Name() string { return "ip_resolution" }
func (ipResolveTask) Run(in *Input, out *Output) error {
	_, err := cachedTask(
		context.Background(),
		in.Cache,
		"ip_resolution:"+in.Domain,
		constants.IPResolutionTTL,
		func() ([]string, error) { return checks.GetIPAddress(in.Domain) },
		func(o *Output, ips []string) { o.IPs = ips },
		out,
	)
	return err
}

// Punycode
type punycodeTask struct{}

func (punycodeTask) Name() string { return "punycode_check" }
func (punycodeTask) Run(in *Input, out *Output) error {
	b, err := checks.ContainsPunycode(in.URL)
	if err != nil {
		return err
	}
	updateOutput(func(o *Output) { o.URLContainsPuny = b })(out)
	return nil
}

// Redirects
type redirectsTask struct{}

func (redirectsTask) Name() string { return "redirect_check" }
func (redirectsTask) Run(in *Input, out *Output) error {
	redir, err := checks.CheckRedirects(in.URL)
	if err != nil {
		return err
	}
	updateOutput(func(o *Output) { o.RedirectionResult = redir })(out)
	return nil
}

// TLD
type tldTask struct{}

func (tldTask) Name() string { return "tld_check" }
func (tldTask) Run(in *Input, out *Output) error {
	t, icann, tld := checks.IsTrustedTld(in.Domain)
	r, _, _ := checks.IsRiskyTld(in.Domain)
	isHosting := isHostingSuffix(tld)
	if platform := checks.CustomerSitePlatform(in.Domain); platform != "" {
		tld, isHosting = platform, true
	}
	updateOutput(func(o *Output) {
		o.TLDTrusted = t
		o.TLDICANN = icann
		o.TLDRisky = r
		o.TLD = tld
		o.TLDIsHostingPlatform = isHosting
	})(out)
	return nil
}

// isHostingSuffix reports whether a public suffix is, or sits under, one of
// the TrustedHostingPlatforms.
func isHostingSuffix(suffix string) bool {
	for d := suffix; d != ""; {
		if _, ok := constants.TrustedHostingPlatforms[d]; ok {
			return true
		}
		i := strings.Index(d, ".")
		if i < 0 {
			break
		}
		d = d[i+1:]
	}
	return false
}

// Shortener
type shortenerTask struct{}

func (shortenerTask) Name() string { return "shortener_check" }
func (shortenerTask) Run(in *Input, out *Output) error {
	updateOutput(func(o *Output) { o.URLIsShortener = checks.IsUrlShortener(in.Domain) })(out)
	return nil
}

// Status
type statusTask struct{}

func (statusTask) Name() string { return "status_code_check" }
func (statusTask) Run(in *Input, out *Output) error {
	code, text, success, redirect, err := checks.GetStatusCode(in.URL)
	if err != nil {
		return err
	}
	updateOutput(func(o *Output) {
		o.StatusCode = code
		o.StatusText = text
		o.StatusSuccess = success
		o.StatusIsRedirect = redirect
	})(out)
	return nil
}

// URL structure
type structureTask struct{}

func (structureTask) Name() string { return "url_structure_check" }
func (structureTask) Run(in *Input, out *Output) error {
	updateOutput(func(o *Output) {
		o.URLTooLong = checks.TooLongUrl(in.URL)
		o.URLTooDeep = checks.TooDeepUrl(in.URL)
	})(out)
	return nil
}

// Keywords
type keywordsTask struct{}

func (keywordsTask) Name() string { return "keywords_check" }
func (keywordsTask) Run(in *Input, out *Output) error {
	present, matches, cats := checks.CheckURLKeywords(in.URL)
	updateOutput(func(o *Output) {
		o.URLKeywordsPresent = present
		o.URLKeywordMatches = matches
		o.URLKeywordCats = cats
	})(out)
	return nil
}

// DNS validity (NS/MX)
type dnsValidityTask struct{}

type dnsValidityResult struct {
	NSValid bool     `json:"ns_valid"`
	NSHosts []string `json:"ns_hosts"`
	MXValid bool     `json:"mx_valid"`
	MXHosts []string `json:"mx_hosts"`
}

func (dnsValidityTask) Name() string { return "dns_validity_check" }
func (dnsValidityTask) Run(in *Input, out *Output) error {
	_, err := cachedTask(
		context.Background(),
		in.Cache,
		"dns_validity:"+in.Domain,
		constants.DNSValidityTTL,
		func() (dnsValidityResult, error) {
			ns, nsHosts, _ := checks.CheckNSValidity(in.Domain)
			mx, mxHosts, _ := checks.CheckMXValidity(in.Domain)
			return dnsValidityResult{
				NSValid: ns,
				NSHosts: nsHosts,
				MXValid: mx,
				MXHosts: mxHosts,
			}, nil
		},
		func(o *Output, cached dnsValidityResult) {
			o.NSValid = cached.NSValid
			o.NSHosts = cached.NSHosts
			o.MXValid = cached.MXValid
			o.MXHosts = cached.MXHosts
		},
		out,
	)
	return err
}

// Subdomains
type subdomainTask struct{}

func (subdomainTask) Name() string { return "subdomain_check" }
func (subdomainTask) Run(in *Input, out *Output) error {
	count, _ := checks.GetSubdomainCount(in.URL)
	updateOutput(func(o *Output) { o.URLSubdomainCount = count })(out)
	return nil
}

// Domain info
type whoisTask struct{}

func (whoisTask) Name() string { return "whois_lookup" }
func (whoisTask) Run(in *Input, out *Output) error {
	// A site under a private suffix (someone.github.io, a bucket on
	// storage.googleapis.com) or on a site builder (someone.weebly.com) isn't
	// registered on its own, so there's no
	// record to find; the lookup would only fail and leave the scan incomplete.
	if suffix, icann := publicsuffix.PublicSuffix(in.Domain); !icann && suffix != in.Domain {
		return nil
	}
	if checks.CustomerSitePlatform(in.Domain) != "" {
		return nil
	}
	_, err := cachedTask(
		context.Background(),
		in.Cache,
		"whois_lookup:"+in.Domain,
		constants.WHOISLookupTTL,
		func() (*domaininfo.RegistrationData, error) {
			return domaininfo.Lookup(in.Domain)
		},
		func(o *Output, di *domaininfo.RegistrationData) {
			o.DomainInfo = di
		},
		out,
	)
	if err != nil && strings.Contains(err.Error(), "no whois server found") {
		return nil
	}
	return err
}

// SSL info
type sslTask struct{}

func (sslTask) Name() string { return "ssl_check" }
func (sslTask) Run(in *Input, out *Output) error {
	updateOutput(func(o *Output) { o.SSLInfo = checks.AnalyzeSSLCert(in.Domain) })(out)
	return nil
}

// Entropy
type entropyTask struct{}

func (entropyTask) Name() string { return "entropy_check" }
func (entropyTask) Run(in *Input, out *Output) error {
	updateOutput(func(o *Output) { o.DomainRandomness = checks.AnalyzeDomainRandomness(in.Domain) })(out)
	return nil
}

// Page content
type contentTask struct{}

func (contentTask) Name() string { return "content_check" }
func (contentTask) Run(in *Input, out *Output) error {
	_, err := cachedTask(
		context.Background(),
		in.Cache,
		"content_check:"+in.URL,
		constants.ContentAnalysisTTL,
		func() (*checks.PageFormResult, error) {
			return checks.GetPageFormInfo(in.URL)
		},
		func(o *Output, c *checks.PageFormResult) { o.ContentData = c },
		out,
	)
	return err
}

// TLS
type tlsTask struct{}

func (tlsTask) Name() string { return "tls_check" }
func (tlsTask) Run(in *Input, out *Output) error {
	t, _ := checks.GetTLSInfo(in.Domain)
	updateOutput(func(o *Output) { o.TLSInfo = t })(out)
	return nil
}

// Homoglyph
type homoglyphTask struct{}

func (homoglyphTask) Name() string { return "homoglyph_check" }
func (homoglyphTask) Run(in *Input, out *Output) error {
	h, _ := checks.HasHomoglyphs(in.Domain)
	updateOutput(func(o *Output) { o.HomoglyphPresent = h })(out)
	return nil
}

// Typosquatting
type typosquatTask struct{}

func (typosquatTask) Name() string { return "typosquat_check" }
func (typosquatTask) Run(in *Input, out *Output) error {
	result := typosquat.CheckTyposquatting(in.Domain)
	typosquat.CheckURLParts(in.URL, in.Domain, &result)
	updateOutput(func(o *Output) { o.TyposquatResult = result })(out)
	return nil
}

// PhishTank
type phishtankTask struct{}

func (phishtankTask) Name() string { return "phishtank_check" }
func (phishtankTask) Run(in *Input, out *Output) error {
	// With PhishTank's dump loaded, threatFeedsTask already checks this URL
	// against every verified, online report. The live API adds only
	// unverified reports and is rate-limited on most scans, which left a
	// "couldn't check PhishTank" note on results the dump had checked.
	// FEED_PHISHTANK=0 or THREAT_LOOKUPS=0 means no PhishTank calls at all.
	if !threatfeeds.PhishTankAPIAllowed() {
		return nil
	}
	fromCache, err := cachedTask(
		context.Background(),
		in.Cache,
		"phishtank:"+in.URL,
		constants.PhishTankTTL,
		func() (*threatfeeds.PhishTankResult, error) {
			return threatfeeds.CheckPhishTank(in.URL)
		},
		func(o *Output, r *threatfeeds.PhishTankResult) { o.PhishTank = r },
		out,
	)
	if err == nil && out.PhishTank != nil {
		out.mu.Lock()
		out.PhishTank.FromCache = fromCache
		out.mu.Unlock()
	}
	return err
}

// Optimized HTTP check (combines redirects, HSTS, and status code)
type httpCombinedTask struct{}

type httpCombinedCacheResult struct {
	RedirectionResult checks.RedirectionResult `json:"redirection_result"`
	StatusCode        int                      `json:"status_code"`
	StatusText        string                   `json:"status_text"`
	StatusSuccess     bool                     `json:"status_success"`
	StatusIsRedirect  bool                     `json:"status_is_redirect"`
	SupportsHSTS      bool                     `json:"supports_hsts"`
}

func (httpCombinedTask) Name() string { return "http_combined_check" }
func (httpCombinedTask) Run(in *Input, out *Output) error {
	ctx := context.Background()
	cacheKey := "http_combined:" + in.URL

	// Try cache first
	if in.Cache != nil {
		var cached httpCombinedCacheResult
		if err := in.Cache.GetJSON(ctx, cacheKey, &cached); err == nil {
			out.mu.Lock()
			out.RedirectionResult = cached.RedirectionResult
			out.StatusCode = cached.StatusCode
			out.StatusText = cached.StatusText
			out.StatusSuccess = cached.StatusSuccess
			out.StatusIsRedirect = cached.StatusIsRedirect
			out.SupportsHSTS = cached.SupportsHSTS
			out.mu.Unlock()
			return nil
		} else if err != redis.Nil {
			// Cache error (not a miss) - log but continue
		} else {
		}
	}

	// Try combined HTTP check first
	combinedResult, err := checks.CheckHTTPCombined(in.URL)
	if err == nil {
		// Success - populate all fields from combined result
		// Store in cache
		if in.Cache != nil {
			cached := httpCombinedCacheResult{
				RedirectionResult: combinedResult.RedirectionResult,
				StatusCode:        combinedResult.StatusCode,
				StatusText:        combinedResult.StatusText,
				StatusSuccess:     combinedResult.StatusSuccess,
				StatusIsRedirect:  combinedResult.StatusIsRedirect,
				SupportsHSTS:      combinedResult.SupportsHSTS,
			}
			_ = in.Cache.SetJSON(ctx, cacheKey, cached, constants.HTTPCombinedTTL)
		}

		out.mu.Lock()
		out.RedirectionResult = combinedResult.RedirectionResult
		out.StatusCode = combinedResult.StatusCode
		out.StatusText = combinedResult.StatusText
		out.StatusSuccess = combinedResult.StatusSuccess
		out.StatusIsRedirect = combinedResult.StatusIsRedirect
		out.SupportsHSTS = combinedResult.SupportsHSTS
		out.mu.Unlock()
		return nil
	}

	// Fallback to individual checks if combined check fails
	var fallbackErrs []error

	// Try redirect check
	redir, err := checks.CheckRedirects(in.URL)
	if err != nil {
		fallbackErrs = append(fallbackErrs, err)
	} else {
		out.mu.Lock()
		out.RedirectionResult = redir
		out.mu.Unlock()
	}

	// Try status code check
	code, text, success, redirect, err := checks.GetStatusCode(in.URL)
	if err != nil {
		fallbackErrs = append(fallbackErrs, err)
	} else {
		out.mu.Lock()
		out.StatusCode = code
		out.StatusText = text
		out.StatusSuccess = success
		out.StatusIsRedirect = redirect
		out.mu.Unlock()
	}

	// Try HSTS check
	h, err := checks.SupportsHSTS(in.URL)
	if err != nil {
		fallbackErrs = append(fallbackErrs, err)
	} else {
		out.mu.Lock()
		out.SupportsHSTS = h
		out.mu.Unlock()
	}

	// If at least one fallback succeeded, cache the result
	if in.Cache != nil {
		cached := httpCombinedCacheResult{
			RedirectionResult: out.RedirectionResult,
			StatusCode:        out.StatusCode,
			StatusText:        out.StatusText,
			StatusSuccess:     out.StatusSuccess,
			StatusIsRedirect:  out.StatusIsRedirect,
			SupportsHSTS:      out.SupportsHSTS,
		}

		_ = in.Cache.SetJSON(ctx, cacheKey, cached, constants.HTTPCombinedTTL)
	}

	// Return error only if all fallbacks failed
	if len(fallbackErrs) == 3 {
		return err // Return the original combined check error
	}

	return nil
}

// Optimized TLS/SSL check (combines both checks into one connection)
type tlsCombinedTask struct{}

type tlsCombinedCacheResult struct {
	TLSInfo checks.TLSResult     `json:"tls_info"`
	SSLInfo checks.SSLCertResult `json:"ssl_info"`
}

func (tlsCombinedTask) Name() string { return "tls_combined_check" }
func (tlsCombinedTask) Run(in *Input, out *Output) error {
	ctx := context.Background()
	// Check the certificate of the host the link points at, which is the one a
	// browser checks: www.example.com can be served a certificate that only
	// covers example.com, or none of the two.
	host, _ := checks.GetHost(in.URL)
	if host == "" {
		host = in.Domain
	}
	cacheKey := "tls_combined:" + host

	// Try cache first
	if in.Cache != nil {
		var cached tlsCombinedCacheResult
		if err := in.Cache.GetJSON(ctx, cacheKey, &cached); err == nil {
			out.mu.Lock()
			out.TLSInfo = cached.TLSInfo
			out.SSLInfo = cached.SSLInfo
			out.mu.Unlock()
			return nil
		} else if err != redis.Nil {
			// Cache error (not a miss) - log but continue
		}
	}

	// Try combined TLS/SSL check first
	combinedResult, err := checks.CheckTLSCombined(host)
	if err == nil {
		// Success - populate both TLS and SSL info
		// Store in cache
		if in.Cache != nil {
			cached := tlsCombinedCacheResult{
				TLSInfo: combinedResult.TLSInfo,
				SSLInfo: combinedResult.SSLInfo,
			}
			_ = in.Cache.SetJSON(ctx, cacheKey, cached, constants.TLSCombinedTTL)
		}

		out.mu.Lock()
		out.TLSInfo = combinedResult.TLSInfo
		out.SSLInfo = combinedResult.SSLInfo
		out.mu.Unlock()
		return nil
	}

	// Fallback to individual checks if combined check fails
	var fallbackErrs []error

	// Try TLS check
	t, err := checks.GetTLSInfo(host)
	if err != nil {
		fallbackErrs = append(fallbackErrs, err)
	} else {
		out.mu.Lock()
		out.TLSInfo = t
		out.mu.Unlock()
	}

	// Try SSL check
	sslInfo := checks.AnalyzeSSLCert(host)
	out.mu.Lock()
	out.SSLInfo = sslInfo
	out.mu.Unlock()

	// Return error only if TLS check failed (SSL check doesn't return error)
	if len(fallbackErrs) > 0 {
		return err // Return the original combined check error
	}

	return nil
}

// Local threat feeds
type threatFeedsTask struct{}

func (threatFeedsTask) Name() string { return "threat_feeds_check" }
func (threatFeedsTask) Run(in *Input, out *Output) error {
	// Host matches are filtered later, once the site's rank is known.
	m := threatfeeds.LookupLocal(in.URL, false)
	updateOutput(func(o *Output) { o.ThreatFeeds = &m })(out)
	return nil
}

// Google Web Risk
type webRiskTask struct{}

func (webRiskTask) Name() string { return "webrisk_check" }
func (webRiskTask) Run(in *Input, out *Output) error {
	if !threatfeeds.WebRiskEnabled() {
		return nil
	}
	_, err := cachedTask(
		context.Background(),
		in.Cache,
		"webrisk:"+in.URL,
		constants.WebRiskTTL,
		func() (*threatfeeds.GoogleThreatResult, error) { return threatfeeds.CheckWebRisk(in.URL) },
		func(o *Output, r *threatfeeds.GoogleThreatResult) { o.WebRisk = r },
		out,
	)
	return err
}

// Google Safe Browsing (hash prefixes only; Google never sees the URL)
type safeBrowsingTask struct{}

func (safeBrowsingTask) Name() string { return "safe_browsing_check" }
func (safeBrowsingTask) Run(in *Input, out *Output) error {
	if !threatfeeds.SafeBrowsingEnabled() {
		return nil
	}
	_, err := cachedTask(
		context.Background(),
		in.Cache,
		"safe_browsing:"+in.URL,
		constants.SafeBrowsingTTL,
		func() (*threatfeeds.GoogleThreatResult, error) { return threatfeeds.CheckSafeBrowsing(in.URL) },
		func(o *Output, r *threatfeeds.GoogleThreatResult) { o.SafeBrowsing = r },
		out,
	)
	return err
}
