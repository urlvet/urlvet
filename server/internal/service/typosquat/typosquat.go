// Package typosquat detects typosquatting by comparing a domain against the
// top-N entries from the Tranco popularity list, and combo-squatting by
// looking for curated brand names inside it.
package typosquat

import (
	"encoding/csv"
	"fmt"
	"net/url"
	"os"
	"sort"
	"strings"

	"github.com/urlvet/urlvet/internal/constants"
	"github.com/urlvet/urlvet/internal/service/checks"
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
	ComboExact    bool   `json:"combo_exact,omitempty"` // just the brand name, on another ending (roblox.dk)

	// Brand names in other parts of the URL (see CheckURLParts).
	SubdomainBrand      string `json:"subdomain_brand,omitempty"`       // "paypal" in paypal-login.example.org
	SubdomainBrandExact bool   `json:"subdomain_brand_exact,omitempty"` // the label is just the name (github.acme.com)
	EmbeddedDomain      string `json:"embedded_domain,omitempty"`       // "google.com" in accounts.google.com.example.fr
	PathBrand           string `json:"path_brand,omitempty"`            // "paypal" in example.org/paypal-de/login
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
	}

	// --- Combo-squat check ---
	// Only curated brands: the top list is full of plain words ("service.gov.uk",
	// "online", "global"), and any small business name contains one of those.
	for _, entry := range comboBrands {
		if strings.Contains(inputSLD, entry.sld) && !isOfficial(strings.ToLower(domain), entry.officialDomains) {
			return TyposquatResult{
				IsSuspicious:  true,
				MatchedDomain: entry.officialDomains[0],
				MatchedBrand:  entry.sld,
				IsComboSquat:  true,
				ComboExact:    inputSLD == entry.sld,
			}
		}
	}

	// Short brands (hsbc, dhl, ups) are too common inside other words to look
	// for anywhere in a name, but as a whole word beside a security word
	// (hsbc-sec, dhl-verify) they're impersonation.
	if b, ok := shortBrandBesideSecurityWord(inputSLD); ok && !isOfficial(strings.ToLower(domain), b.officialDomains) {
		return TyposquatResult{
			IsSuspicious:  true,
			MatchedDomain: b.officialDomains[0],
			MatchedBrand:  b.sld,
			IsComboSquat:  true,
		}
	}

	return TyposquatResult{}
}

// comboBrand is a brand name that phishing domains embed ("paypal-secure-login").
type comboBrand struct {
	sld             string
	officialDomains []string
}

// minComboLen keeps short names out: "apple" is in "pineapple", "chase" in "purchase".
const minComboLen = 6

// genericBrandNames are brand names that are also everyday words, so finding
// one inside a domain says nothing about impersonation.
var genericBrandNames = map[string]struct{}{
	"battle": {}, "gemini": {}, "indeed": {}, "kraken": {}, "ledger": {},
	"notion": {}, "office": {}, "origin": {}, "phantom": {}, "signal": {},
	"stripe": {}, "netbank": {}, "medicare": {}, "proton": {}, "messenger": {},
	"freelancer": {}, "incometax": {},
	// Places and plain words: fragrancecanada.ca, fashionunited.nl, orangebikes.com.
	"canada": {}, "united": {}, "orange": {}, "target": {}, "booking": {},
	"spectrum": {}, "sunrise": {}, "fidelity": {}, "emirates": {}, "regions": {},
	"discover": {}, "square": {}, "citizens": {}, "nationwide": {}, "vanguard": {},
	// Inside other words: interac(tive), revolut(ion).
	"interac": {}, "revolut": {},
	// Hundreds of real regional banks and fan sites carry these names.
	"sparkasse": {}, "volksbank": {}, "raiffeisen": {}, "minecraft": {},
	// Words and parts of words: crypto(lake), blockchain(week), (l)ankama(trimony).
	"crypto": {}, "blockchain": {}, "turkiye": {}, "charter": {}, "secureserver": {}, "ankama": {},
}

var comboBrands = buildComboBrands()

// shortBrands are curated brand names too short for comboBrands.
var shortBrands = buildShortBrands()

// buildShortBrands collects short brand names from every brand's official
// domains. A name can belong to several entries (dhl to DHL and Deutsche
// Post, xbox to Xbox and Minecraft), so their domains are merged rather than
// letting map order pick one. The domains of the entry named after the brand
// (whose first official domain is dhl.com) come first, in that entry's order,
// since the first one is reported as the match.
func buildShortBrands() map[string]comboBrand {
	names := make([]string, 0, len(constants.HighValueBrands))
	for name := range constants.HighValueBrands {
		names = append(names, name)
	}
	sort.Strings(names)

	own := map[string][]string{}   // from the entry named after the brand
	other := map[string][]string{} // from entries that also own the name
	for _, name := range names {
		entry := constants.HighValueBrands[name]
		for _, d := range entry.OfficialDomains {
			sld := extractSLD(d)
			_, generic := genericBrandNames[sld]
			_, word := shortEverydayWords[sld]
			if generic || word || len(sld) < 3 || len(sld) >= minComboLen {
				continue
			}
			if extractSLD(entry.OfficialDomains[0]) == sld {
				own[sld] = append(own[sld], entry.Domains()...)
			} else {
				other[sld] = append(other[sld], entry.Domains()...)
			}
		}
	}

	out := map[string]comboBrand{}
	for _, m := range []map[string][]string{own, other} {
		for sld, list := range m {
			b := out[sld]
			b.sld = sld
			b.officialDomains = append(b.officialDomains, list...)
			out[sld] = b
		}
	}
	for sld, b := range out {
		b.officialDomains = dedupe(b.officialDomains)
		out[sld] = b
	}
	return out
}

// dedupe drops repeats from list, keeping the first of each.
func dedupe(list []string) []string {
	seen := make(map[string]struct{}, len(list))
	out := list[:0]
	for _, s := range list {
		if _, ok := seen[s]; !ok {
			seen[s] = struct{}{}
			out = append(out, s)
		}
	}
	return out
}

// shortEverydayWords are short brand names that are also ordinary words in
// domain names (web-security, cloud-update, free-online). visa and axis also
// name real businesses beside these words: visa agencies (uk-visa-services)
// and Axis Security.
var shortEverydayWords = map[string]struct{}{
	"ants": {}, "axis": {}, "bell": {}, "box": {}, "cash": {}, "cloud": {}, "cms": {}, "com": {},
	"delta": {}, "force": {}, "free": {}, "gate": {}, "globe": {}, "gov": {}, "gouv": {},
	"grab": {}, "jet": {}, "key": {}, "line": {}, "lido": {}, "live": {}, "mail": {},
	"match": {}, "max": {}, "nsw": {}, "post": {}, "sat": {}, "sky": {}, "trip": {},
	"visa": {}, "web": {}, "wise": {},
}

// securityWords are what phishing pairs a brand with in a domain.
var securityWords = map[string]struct{}{
	"sec": {}, "secure": {}, "security": {}, "login": {}, "logon": {}, "signin": {},
	"verify": {}, "verification": {}, "account": {}, "accounts": {}, "auth": {},
	"update": {}, "support": {}, "help": {}, "service": {}, "services": {},
	"online": {}, "web": {}, "portal": {}, "wallet": {}, "billing": {}, "pay": {},
	"payment": {}, "unlock": {}, "recovery": {}, "alert": {}, "confirm": {},
	"tracking": {}, "track": {}, "delivery": {}, "parcel": {}, "customs": {},
}

// shortBrandBesideSecurityWord finds a short brand as a whole word of a
// hyphenated name with a security word beside it ("hsbc-sec", "verify-dhl").
func shortBrandBesideSecurityWord(sld string) (comboBrand, bool) {
	words := strings.Split(sld, "-")
	if len(words) < 2 {
		return comboBrand{}, false
	}
	for i, w := range words {
		b, ok := shortBrands[w]
		if !ok {
			continue
		}
		for _, j := range []int{i - 1, i + 1} {
			if j < 0 || j >= len(words) {
				continue
			}
			if _, sec := securityWords[words[j]]; sec {
				return b, true
			}
		}
	}
	return comboBrand{}, false
}

func buildComboBrands() []comboBrand {
	seen := map[string]bool{}
	var out []comboBrand
	for _, entry := range constants.HighValueBrands {
		for _, d := range entry.OfficialDomains {
			sld := extractSLD(d)
			if _, generic := genericBrandNames[sld]; generic || len(sld) < minComboLen || seen[sld] {
				continue
			}
			seen[sld] = true
			out = append(out, comboBrand{sld: sld, officialDomains: entry.Domains()})
		}
	}
	// Longest first, so "disneyplus-free.com" reports "disneyplus", not "disney".
	sort.Slice(out, func(i, j int) bool {
		if len(out[i].sld) != len(out[j].sld) {
			return len(out[i].sld) > len(out[j].sld)
		}
		return out[i].sld < out[j].sld
	})
	return out
}

func isOfficial(domain string, officialDomains []string) bool {
	for _, d := range officialDomains {
		if domain == d || strings.HasSuffix(domain, "."+d) {
			return true
		}
	}
	return false
}

// extractSLD returns just the second-level domain (SLD) of a hostname,
// stripping both subdomains and the public suffix.
// e.g. "www.paypal.com" → "paypal", "example.co.uk" → "example"
func extractSLD(domain string) string {
	// A customer site on a builder: the customer's own label is the name
	// someone chose, so that's what to compare ("xfinitylogin.weebly.com").
	if platform := checks.CustomerSitePlatform(domain); platform != "" {
		// The whole customer part: 1286524.us23.myftpupload.com → "1286524.us23".
		return strings.TrimSuffix(domain, "."+platform)
	}
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

// lookupRank is rank.DomainRankLookup, swappable in tests.
var lookupRank = rank.DomainRankLookup

// embeddedRank is how popular a domain spelled out inside a hostname must be
// to count as impersonation: accounts.google.com.example.fr, not www.co.example.
const embeddedRank = 100000

// brandIn returns the first curated brand whose name s contains, longest first.
func brandIn(s string) (comboBrand, bool) {
	for _, b := range comboBrands {
		if strings.Contains(s, b.sld) {
			return b, true
		}
	}
	return comboBrand{}, false
}

// CheckURLParts looks for brands outside the registered name, where
// CheckTyposquatting doesn't: in the subdomains (horyzonix.paypal-login.
// antimoney-laundering.org), a well-known domain spelled out at the front
// (accounts.google.com.asso-entrautres.fr), and the path (/Paypal-de/login).
// domain is the registered name (or a hosting customer's site) from GetDomain.
func CheckURLParts(rawURL, domain string, res *TyposquatResult) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return
	}
	host := strings.ToLower(u.Hostname())
	domain = strings.ToLower(domain)

	bare := strings.TrimPrefix(host, "www.")
	if prefix := strings.TrimSuffix(bare, "."+domain); bare != domain && prefix != bare && prefix != "" && prefix != "www" {
		labels := strings.Split(prefix, ".")
		// A popular domain spelled out in front: google.com in accounts.google.com.example.fr.
		for i := 0; i < len(labels)-1 && res.EmbeddedDomain == ""; i++ {
			candidate := strings.Join(labels[i:], ".")
			if r := lookupRank(candidate); r > 0 && r <= embeddedRank {
				res.EmbeddedDomain = candidate
			}
		}
		for _, label := range labels {
			if b, ok := brandIn(label); ok && !isOfficial(host, b.officialDomains) {
				res.SubdomainBrand = b.sld
				res.SubdomainBrandExact = label == b.sld
				break
			}
		}
	}

	if path := strings.ToLower(u.EscapedPath() + "?" + u.RawQuery); len(path) > 2 {
		if b, ok := brandIn(path); ok && !isOfficial(host, b.officialDomains) {
			res.PathBrand = b.sld
		}
	}
}
