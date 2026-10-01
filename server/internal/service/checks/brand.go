package checks

import (
	"strings"

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
			if strings.Contains(pageTitle, kw) {
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
