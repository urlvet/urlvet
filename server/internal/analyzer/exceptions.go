package analyzer

import (
	neturl "net/url"
	"strings"
)

// isTopRankedDomain reports whether the domain is in the top-1M popularity list.
// Popular domains are trusted enough that name-based heuristics (typosquatting,
// risky TLD) produce far more false positives than real detections.
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

	if isTopRankedDomain(out.Rank) {
		out.TyposquatResult.IsSuspicious = false
		out.TyposquatResult.IsComboSquat = false
		out.TLDRisky = false
	}

	if out.URLIsShortener && isBareURL(in.URL) {
		out.URLIsShortener = false
	}
}
