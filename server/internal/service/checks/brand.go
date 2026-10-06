package checks

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/urlvet/urlvet/internal/constants"
	"golang.org/x/text/unicode/norm"
)

type BrandResult struct {
	BrandFound string `json:"brand_found"`
	// OfficialDomain is the brand's own site when the page claims to be a brand
	// this domain doesn't belong to, so the UI can point people to the real one.
	OfficialDomain string   `json:"official_domain,omitempty"`
	IsMismatch     bool     `json:"is_mismatch"`
	DetectedNames  []string `json:"detected_names"`
}

func isOfficialDomain(domain string, officialDomains []string) bool {
	for _, official := range officialDomains {
		if domain == official || strings.HasSuffix(domain, "."+official) {
			return true
		}
	}
	return false
}

// isBrandsOwn reports whether domain belongs to the brand: one of its official
// domains, or for member groups like the Sparkassen, any domain carrying the
// group's name.
func isBrandsOwn(domain string, entry constants.BrandEntry) bool {
	if isOfficialDomain(domain, entry.Domains()) {
		return true
	}
	for _, name := range entry.OwnNames {
		if strings.Contains(domain, name) {
			return true
		}
	}
	return false
}

func CheckBrandMismatch(domain string, pageTitle string) BrandResult {
	domain = strings.ToLower(domain)
	pageTitle = withoutOwnAddress(foldTitle(pageTitle), domain)

	res := BrandResult{
		DetectedNames: []string{},
	}

	for brand, entry := range constants.HighValueBrands {
		for _, kw := range append(entry.TitleKeywords, exactTitle(pageTitle, entry.ExactTitles)...) {
			if containsWord(pageTitle, foldTitle(kw)) {
				res.DetectedNames = append(res.DetectedNames, brand)
				if !isBrandsOwn(domain, entry) {
					res.BrandFound = brand
					res.IsMismatch = true
					if len(entry.OfficialDomains) > 0 {
						res.OfficialDomain = entry.OfficialDomains[0]
					}
				}
				break
			}
		}
	}

	return res
}

// exactTitle returns the one of titles the whole page title is, if any.
func exactTitle(pageTitle string, titles []string) []string {
	for _, t := range titles {
		if strings.TrimSpace(pageTitle) == foldTitle(t) {
			return []string{t}
		}
	}
	return nil
}

// CheckBrandNames is the loose check: it matches brands' bare names
// ("facebook", "vodafone"), which ordinary sites also put in titles. Callers
// use it only where that's telling: a page asking for a login or payment, or
// a free hosting subdomain.
func CheckBrandNames(domain string, pageTitle string) BrandResult {
	domain = strings.ToLower(domain)
	pageTitle = withoutOwnAddress(foldTitle(pageTitle), domain)
	res := BrandResult{DetectedNames: []string{}}
	for brand, entry := range constants.HighValueBrands {
		for _, name := range entry.Names {
			if !containsWord(pageTitle, foldTitle(name)) {
				continue
			}
			res.DetectedNames = append(res.DetectedNames, brand)
			if !isBrandsOwn(domain, entry) && !res.IsMismatch {
				res.BrandFound = brand
				res.IsMismatch = true
				if len(entry.OfficialDomains) > 0 {
					res.OfficialDomain = entry.OfficialDomains[0]
				}
			}
			break
		}
	}
	return res
}

// withoutOwnAddress removes the page's own hostname from its title. Sites
// titled with their address ("Octocat.github.io", the GitHub Pages default)
// would otherwise claim the hosting provider's brand.
func withoutOwnAddress(title, host string) string {
	host = strings.TrimPrefix(host, "www.")
	if host == "" {
		return title
	}
	title = strings.ReplaceAll(title, "www."+host, " ")
	return strings.ReplaceAll(title, host, " ")
}

// lookalikes maps letters phishing titles swap in for Latin ones
// ("Çoinbase", "Pаypal" with a Cyrillic а) that don't decompose to them.
var lookalikes = strings.NewReplacer(
	"ø", "o", "ł", "l", "đ", "d", "ħ", "h", "ı", "i", "ŧ", "t", "ß", "ss", "æ", "ae", "œ", "oe",
	// Cyrillic
	"а", "a", "е", "e", "о", "o", "р", "p", "с", "c", "у", "y", "х", "x", "і", "i", "ѕ", "s", "ј", "j", "ԁ", "d", "һ", "h", "ӏ", "l",
	// Greek
	"α", "a", "ο", "o", "ν", "v", "ρ", "p", "τ", "t", "ι", "i", "κ", "k", "υ", "u",
)

// foldTitle lowercases a title and folds accented and lookalike letters to
// plain Latin, so "Çoinbase Pro: Løgin" matches "coinbase".
func foldTitle(s string) string {
	var b strings.Builder
	for _, r := range norm.NFKD.String(s) {
		if unicode.Is(unicode.Mn, r) {
			continue // the accent split off a letter
		}
		b.WriteRune(r)
	}
	// Lowercase after decomposing: "𝗠" (mathematical bold) has no lowercase
	// form of its own, only the plain "M" it decomposes to does.
	folded := lookalikes.Replace(strings.ToLower(b.String()))
	// Titles separate words with ":", "|", "–" and the like ("Gemini : Login");
	// treat those as spaces. Characters brand names contain (crypto.com,
	// at&t, t-mobile, lowe's) stay.
	return strings.Join(strings.Fields(titleSeparators.ReplaceAllString(folded, " ")), " ")
}

var titleSeparators = regexp.MustCompile(`[^\p{L}\p{N}.&+'’\-\s]+`)

// containsWord reports whether kw appears in s as a whole word or phrase, so
// "netbank" (Commonwealth Bank) doesn't match "NetBanking" on another bank's page.
func containsWord(s, kw string) bool {
	isWordChar := func(r rune) bool { return unicode.IsLetter(r) || unicode.IsDigit(r) }
	for i := 0; ; {
		j := strings.Index(s[i:], kw)
		if j < 0 {
			return false
		}
		start, end := i+j, i+j+len(kw)
		before, _ := utf8.DecodeLastRuneInString(s[:start])
		after, _ := utf8.DecodeRuneInString(s[end:])
		if (start == 0 || !isWordChar(before)) && (end == len(s) || !isWordChar(after)) {
			return true
		}
		i = start + 1
	}
}
