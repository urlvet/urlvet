package checks

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/urlvet/urlvet/internal/constants"
)

type BrandResult struct {
	BrandFound    string   `json:"brand_found"`
	IsMismatch    bool     `json:"is_mismatch"`
	DetectedNames []string `json:"detected_names"`
}

func isOfficialDomain(domain string, officialDomains []string) bool {
	for _, official := range officialDomains {
		if domain == official || strings.HasSuffix(domain, "."+official) {
			return true
		}
	}
	return false
}

func CheckBrandMismatch(domain string, pageTitle string) BrandResult {
	domain = strings.ToLower(domain)
	pageTitle = strings.ToLower(pageTitle)

	res := BrandResult{
		DetectedNames: []string{},
	}

	for brand, entry := range constants.HighValueBrands {
		for _, kw := range entry.TitleKeywords {
			if containsWord(pageTitle, kw) {
				res.DetectedNames = append(res.DetectedNames, brand)
				if !isOfficialDomain(domain, entry.OfficialDomains) {
					res.BrandFound = brand
					res.IsMismatch = true
				}
				break
			}
		}
	}

	return res
}

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
