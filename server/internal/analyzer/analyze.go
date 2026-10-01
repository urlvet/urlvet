package analyzer

import (
	"context"
	"errors"
	"slices"
	"time"

	"github.com/urlvet/urlvet/internal/constants"
	"github.com/urlvet/urlvet/internal/logger"
	"github.com/urlvet/urlvet/internal/metrics"
	"github.com/urlvet/urlvet/internal/service/cache"
	"github.com/urlvet/urlvet/internal/service/checks"
	"github.com/urlvet/urlvet/internal/service/threatfeeds"
	"github.com/urlvet/urlvet/internal/store"
)

// buildPhishingResult converts the internal PhishTankResult into the public PhishingResult.
// Returns nil when no check was performed (cache miss + fetch error).
func buildPhishingResult(r *threatfeeds.PhishTankResult) *PhishingResult {
	if r == nil {
		return nil
	}
	return &PhishingResult{
		InDatabase:      r.InDatabase,
		PhishID:         r.PhishID,
		PhishDetailPage: r.PhishDetailPage,
		Verified:        r.Verified,
		VerifiedAt:      r.VerifiedAt,
		Valid:           r.Valid,
		Target:          r.Target,
		Source:          "phishtank",
		FromCache:       r.FromCache,
		RawResponse:     r.RawResponse,
	}
}

// Analyze runs all tasks and builds the final response
func Analyze(ctx context.Context, rawURL string) (Response, []error) {
	// Validate and normalize the URL (adds https:// if scheme is missing).
	// Using the normalized form as the canonical URL ensures consistent cache
	// keys regardless of how the caller formatted the input.
	parsedURL, isValid, err := checks.IsValidURL(rawURL)
	if err != nil || !isValid {
		return Response{}, []error{ErrInvalidURL}
	}
	normalizedURL := parsedURL.String()

	domain, err := checks.GetDomain(normalizedURL)
	if err != nil {
		return Response{}, []error{err}
	}

	// Initialize cache (non-blocking - if cache fails, continue without it)
	var cacheInstance CacheInterface
	cacheConn, cacheErr := cache.New()
	if cacheErr != nil {
		logger.Warn("cache init failed, continuing without cache", "err", cacheErr)
	} else {
		cacheInstance = cacheConn
		defer cacheConn.Close()
	}

	in := &Input{URL: normalizedURL, Domain: domain, Cache: cacheInstance}

	// Full-result cache: if we have a recent scan for this URL, return it immediately
	// without re-running all tasks. TTL is 24h — same as the slowest individual task.
	start := time.Now()
	resultKey := "analyze_result:" + normalizedURL
	if cacheInstance != nil {
		var cached Response
		if err := cacheInstance.GetJSON(context.Background(), resultKey, &cached); err == nil {
			cached.Performance.TotalTime = time.Since(start).String()
			store.AddScan(store.ScanRecord{
				URL:      normalizedURL,
				Domain:   domain,
				Verdict:  cached.Result.Verdict,
				Score:    cached.Result.FinalScore,
				Duration: cached.Performance.TotalTime,
				Time:     time.Now(),
				Cached:   true,
			})
			return cached, nil
		}
	}

	tasks := []Task{
		rankTask{},
		httpCombinedTask{}, // Optimized: combines redirects, HSTS, and status code
		ipCheckTask{},
		ipResolveTask{},
		punycodeTask{},
		tldTask{},
		shortenerTask{},
		structureTask{},
		keywordsTask{},
		dnsValidityTask{},
		subdomainTask{},
		whoisTask{},
		tlsCombinedTask{}, // Optimized: combines TLS and SSL checks
		entropyTask{},
		contentTask{},
		homoglyphTask{},
		phishtankTask{},
		typosquatTask{},
	}

	// Reset timer to measure only the actual task execution time
	start = time.Now()
	out, errs := runTasks(ctx, in, tasks)
	applyExceptions(in, out)

	resp := Response{
		URL:    normalizedURL,
		Domain: domain,
		Features: Features{
			Rank: out.Rank,
			TLD: TLDInfo{
				TLD:               out.TLD,
				IsTrusted:         out.TLDTrusted,
				IsRisky:           out.TLDRisky,
				IsICANN:           out.TLDICANN,
				IsHostingPlatform: out.TLDIsHostingPlatform,
			},
			URL: URLChecks{
				IsURLShortener:   out.URLIsShortener,
				UsesIP:           out.URLUsesIP,
				ContainsPunycode: out.URLContainsPuny,
				TooLong:          out.URLTooLong,
				TooDeep:          out.URLTooDeep,
				SubdomainCount:   out.URLSubdomainCount,
				HasHomoglyph:     out.HomoglyphPresent,
				Keywords: Keywords{
					HasKeywords: out.URLKeywordsPresent,
					Found:       out.URLKeywordMatches,
					Categories:  out.URLKeywordCats,
				},
			},
		},
		Infrastructure: Infrastructure{
			IPAddresses:      out.IPs,
			NameserversValid: out.NSValid,
			NSHosts:          out.NSHosts,
			MXRecordsValid:   out.MXValid,
			MXHosts:          out.MXHosts,
		},
		DomainInfo: out.DomainInfo,
		Analysis: Analysis{
			RedirectionResult: out.RedirectionResult,
			SupportsHSTS:      out.SupportsHSTS,
			HTTPStatus: HTTPStatus{
				Code:                 out.StatusCode,
				Text:                 out.StatusText,
				Success:              out.StatusSuccess,
				IsRedirectStatusCode: out.StatusIsRedirect,
			},
		},
		SSLInfo:          out.SSLInfo,
		TLSInfo:          out.TLSInfo,
		ContentData:      out.ContentData,
		DomainRandomness: out.DomainRandomness,
		TyposquatResult:  out.TyposquatResult,
		Phishing:         buildPhishingResult(out.PhishTank),
		Performance: Performance{
			TotalTime: time.Since(start).String(),
			Timings:   ConvertTimings(out.Timings),
		},
	}

	result := GenerateResult(resp)
	resp.Result = result
	metrics.RiskScore.Observe(float64(result.RiskScore))
	metrics.TrustScore.Observe(float64(result.TrustScore))

	resp.Incomplete, resp.IncompleteChecks, resp.Errors = summarizeErrors(errs)

	// Only cache complete results — incomplete scans may be missing signals.
	if cacheInstance != nil && !resp.Incomplete {
		_ = cacheInstance.SetJSON(context.Background(), resultKey, resp, constants.AnalyzeResultTTL)
	}

	store.AddScan(store.ScanRecord{
		URL:      normalizedURL,
		Domain:   domain,
		Verdict:  result.Verdict,
		Score:    result.FinalScore,
		Duration: resp.Performance.TotalTime,
		Time:     time.Now(),
		Cached:   false,
	})

	for _, e := range errs {
		store.AddError(store.ErrorRecord{
			Task:  "analyze",
			Error: e.Error(),
			URL:   normalizedURL,
			Time:  time.Now(),
		})
	}

	return resp, errs
}

// summarizeErrors turns task errors into the response's incomplete flag, the
// names of the checks that didn't finish, and the raw messages.
func summarizeErrors(errs []error) (incomplete bool, checks []string, messages []string) {
	for _, e := range errs {
		messages = append(messages, e.Error())
		var taskErr *TaskError
		var timeout *TimeoutError
		switch {
		case errors.As(e, &timeout):
			checks = append(checks, timeout.Tasks...)
		case errors.As(e, &taskErr):
			checks = append(checks, taskErr.Task)
		}
		// PhishTank rate limits come and go with traffic; a scan missing only
		// that lookup is still worth caching rather than retrying all day.
		if !errors.Is(e, threatfeeds.ErrRateLimited) {
			incomplete = true
		}
	}
	slices.Sort(checks)
	return incomplete, slices.Compact(checks), messages
}
