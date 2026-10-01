// Package typosquat detects typosquatting and combo-squatting by comparing
// a domain against the top-N entries from the Tranco popularity list.
package typosquat

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/urlvet/urlvet/internal/constants"
	"github.com/urlvet/urlvet/internal/service/rank"
	"golang.org/x/net/publicsuffix"
)

// topN controls how many Tranco entries are loaded and compared against.
// 500 covers essentially all consumer-facing brands worth impersonating while
// keeping false-positive rates low.
const topN = 5000

// minSLDLen is the minimum SLD character length required before running any
// string-distance comparison.  Short SLDs (e.g. "io", "ai") produce too many
// accidental near-matches.
const minSLDLen = 4

// TyposquatResult is the output of CheckTyposquatting.
type TyposquatResult struct {
	IsSuspicious  bool   `json:"is_suspicious"`
	MatchedDomain string `json:"matched_domain,omitempty"`
	MatchedBrand  string `json:"matched_brand,omitempty"` // SLD only
	Distance      int    `json:"distance,omitempty"`      // 0 for combo-squat
	IsComboSquat  bool   `json:"is_combo_squat,omitempty"`
}

// topEntry holds the full domain and its extracted SLD for a Tranco entry.
type topEntry struct {
	domain string
	sld    string
}

var topEntries []topEntry
var topSLDSet map[string]struct{}

// LoadTopDomains reads the Tranco CSV (rank,domain — no header) and populates
// the in-memory list used by CheckTyposquatting.  Call once at startup.
func LoadTopDomains() error {
	csvPath := constants.DOMAIN_RANK_FILE_PATH
	f, err := os.Open(csvPath)
	if err != nil {
		return fmt.Errorf("typosquat: open %s: %w", csvPath, err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		return fmt.Errorf("typosquat: read csv: %w", err)
	}

	entries := make([]topEntry, 0, topN)
	for _, rec := range records {
		if len(entries) >= topN {
			break
		}
		if len(rec) < 2 {
			continue
		}
		domain := strings.ToLower(strings.TrimSpace(rec[1]))
		sld := extractSLD(domain)
		if len(sld) < minSLDLen {
			continue // skip short SLDs
		}
		entries = append(entries, topEntry{domain: domain, sld: sld})
	}

	topEntries = entries

	sldSet := make(map[string]struct{}, len(entries))
	for _, e := range entries {
		sldSet[e.sld] = struct{}{}
	}
	topSLDSet = sldSet

	return nil
}

// GetTopEntries returns the loaded top-domain slice (read-only).
func GetTopEntries() []topEntry {
	return topEntries
}

// CheckTyposquatting returns a TyposquatResult for the given domain.
// It detects:
//   - Levenshtein typosquatting (edit distance 1-2 from a brand SLD)
//   - Combo-squatting (brand SLD embedded inside a longer domain SLD)
func CheckTyposquatting(domain string) TyposquatResult {
	inputSLD := extractSLD(strings.ToLower(domain))
	if len(inputSLD) < minSLDLen {
		return TyposquatResult{}
	}

	// If the input is itself a well-known domain, it cannot be a typosquat of another.
	if domainRank := rank.DomainRankLookup(domain); domainRank > 0 {
		return TyposquatResult{}
	}

	for _, entry := range topEntries {
		// --- Levenshtein check ---
		// Two edits is a lot for a short name ("hdfc" is two from "htsc"), so
		// allow it only when both names are long enough for it to look alike.
		maxDist := 1
		if len(inputSLD) >= 6 && len(entry.sld) >= 6 {
			maxDist = 2
		}
		dist := levenshtein(inputSLD, entry.sld)
		if dist >= 1 && dist <= maxDist {
			return TyposquatResult{
				IsSuspicious:  true,
				MatchedDomain: entry.domain,
				MatchedBrand:  entry.sld,
				Distance:      dist,
			}
		}

		// --- Combo-squat check ---
		// Require brand SLD ≥ 6 chars to avoid matching generic words like "mail"
		if len(entry.sld) >= 6 && strings.Contains(inputSLD, entry.sld) {
			return TyposquatResult{
				IsSuspicious:  true,
				MatchedDomain: entry.domain,
				MatchedBrand:  entry.sld,
				IsComboSquat:  true,
			}
		}
	}

	return TyposquatResult{}
}

// extractSLD returns just the second-level domain (SLD) of a hostname,
// stripping both subdomains and the public suffix.
// e.g. "www.paypal.com" → "paypal", "example.co.uk" → "example"
func extractSLD(domain string) string {
	// EffectiveTLDPlusOne strips subdomains: www.paypal.com → paypal.com
	reg, err := publicsuffix.EffectiveTLDPlusOne(domain)
	if err != nil {
		// Fallback: strip everything after first dot
		if i := strings.Index(domain, "."); i > 0 {
			return domain[:i]
		}
		return domain
	}
	// PublicSuffix returns the eTLD: paypal.com → "com", example.co.uk → "co.uk"
	etld, _ := publicsuffix.PublicSuffix(reg)
	sld := strings.TrimSuffix(reg, "."+etld)
	return sld
}

// levenshtein computes the edit distance between two strings.
func levenshtein(a, b string) int {
	ra, rb := []rune(a), []rune(b)
	la, lb := len(ra), len(rb)

	// dp[i][j] = edit distance between ra[:i] and rb[:j]
	dp := make([][]int, la+1)
	for i := range dp {
		dp[i] = make([]int, lb+1)
	}
	for i := 0; i <= la; i++ {
		dp[i][0] = i
	}
	for j := 0; j <= lb; j++ {
		dp[0][j] = j
	}
	for i := 1; i <= la; i++ {
		for j := 1; j <= lb; j++ {
			if ra[i-1] == rb[j-1] {
				dp[i][j] = dp[i-1][j-1]
			} else {
				dp[i][j] = 1 + min(dp[i-1][j], dp[i][j-1], dp[i-1][j-1])
			}
		}
	}
	return dp[la][lb]
}
