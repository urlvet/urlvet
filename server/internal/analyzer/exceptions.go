package analyzer

import (
	neturl "net/url"
	"strings"

	"github.com/urlvet/urlvet/internal/service/checks"
	"github.com/urlvet/urlvet/internal/service/rank"
)

// isTopRankedDomain reports whether the domain is in the top-1M popularity list.
// Popular domains are trusted enough that name-based heuristics (typosquatting,
// risky TLD, sensitive keywords) produce far more false positives than real
// detections: paypal.com contains "paypal", login.gov "login", and pages like
// github.com/login or support.apple.com are where those words belong.
func isTopRankedDomain(rank int) bool {
	return rank > 0
}

// isBareURL reports whether the URL has no path, query, or fragment
// (e.g. "https://bit.ly" or "https://bit.ly/"). For a URL shortener this is
// the service's own homepage, not a short link hiding a destination.
func isBareURL(rawURL string) bool {
	u, err := neturl.Parse(rawURL)
	if err != nil {
		return false
	}
	return strings.Trim(u.Path, "/") == "" && u.RawQuery == "" && u.Fragment == ""
}

// applyExceptions clears signals that are known false positives, so the score,
// the flags shown in the UI, and the cached result all agree.
func applyExceptions(in *Input, out *Output) {
	out.mu.Lock()
	defer out.mu.Unlock()

	// Restricted registries (.gov, .bank.in, .edu, …) verify who registers a
	// domain, so the same name-based heuristics don't apply there either.
	if isTopRankedDomain(out.Rank) || out.TLDTrusted {
		out.TyposquatResult.IsSuspicious = false
		out.TyposquatResult.IsComboSquat = false
		out.TLDRisky = false
		out.URLKeywordsPresent = false
		out.URLKeywordMatches = nil
		out.URLKeywordCats = nil
	}

	// Only a verified organisation can register on a restricted registry, so a
	// bank's page on hdfc.bank.in naming HDFC isn't impersonation, even before
	// the brand list learns the new address.
	if out.TLDTrusted && out.ContentData != nil && out.ContentData.BrandCheck.IsMismatch {
		brand := *out.ContentData
		brand.BrandCheck.IsMismatch = false
		brand.BrandCheck.BrandFound = ""
		out.ContentData = &brand
	}

	// A well-known site moving to its own new address (zoom.us → zoom.com,
	// hdfcbank.com → hdfc.bank.in) isn't the kind of jump phishing makes: the
	// destination is either well-known too, or on a registry that only verified
	// organisations can use.
	if isTopRankedDomain(out.Rank) && out.RedirectionResult.HasDomainJump &&
		isVouchedDestination(out.RedirectionResult.FinalURL) {
		out.RedirectionResult.HasDomainJump = false
	}

	if out.URLIsShortener && isBareURL(in.URL) {
		out.URLIsShortener = false
	}
}

// isVouchedDestination reports whether a redirect target is well-known (top-1M)
// or on a restricted registry such as .gov or .bank.in.
func isVouchedDestination(finalURL string) bool {
	domain, err := checks.GetDomain(finalURL)
	if err != nil || domain == "" {
		return false
	}
	if rank.DomainRankLookup(domain) > 0 {
		return true
	}
	trusted, _, _ := checks.IsTrustedTld(domain)
	return trusted
}
