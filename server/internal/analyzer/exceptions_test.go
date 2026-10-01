package analyzer

import (
	"testing"

	"github.com/urlvet/urlvet/internal/service/checks"
	"github.com/urlvet/urlvet/internal/service/typosquat"
)

func TestApplyExceptions_TopRankedDomain(t *testing.T) {
	out := &Output{OutputData: OutputData{
		Rank:            42000,
		TLDRisky:        true,
		TyposquatResult: typosquat.TyposquatResult{IsSuspicious: true, IsComboSquat: true, MatchedBrand: "paypal"},

		URLKeywordsPresent: true,
		URLKeywordMatches:  []string{"paypal"},
		URLKeywordCats:     map[string][]string{"finance": {"paypal"}},
	}}
	applyExceptions(&Input{URL: "https://example.xyz"}, out)

	if out.TLDRisky {
		t.Error("risky TLD should be ignored for top-ranked domains")
	}
	if out.TyposquatResult.IsSuspicious || out.TyposquatResult.IsComboSquat {
		t.Error("typosquat should be ignored for top-ranked domains")
	}
	if out.URLKeywordsPresent || out.URLKeywordMatches != nil || out.URLKeywordCats != nil {
		t.Error("URL keywords should be ignored for top-ranked domains")
	}
}

func TestApplyExceptions_UnrankedDomainKeepsSignals(t *testing.T) {
	out := &Output{OutputData: OutputData{
		Rank:            0,
		TLDRisky:        true,
		TyposquatResult: typosquat.TyposquatResult{IsSuspicious: true},

		URLKeywordsPresent: true,
		URLKeywordMatches:  []string{"paypal", "login"},
	}}
	applyExceptions(&Input{URL: "https://paypal-login.xyz"}, out)

	if !out.TLDRisky || !out.TyposquatResult.IsSuspicious || !out.URLKeywordsPresent {
		t.Error("unranked domains must keep risky TLD, typosquat and keyword signals")
	}
}

func TestApplyExceptions_URLShortener(t *testing.T) {
	tests := []struct {
		url  string
		want bool
	}{
		{"https://bit.ly", false},
		{"https://bit.ly/", false},
		{"https://bit.ly/abc123", true},
		{"https://bit.ly/?x=1", true},
		{"https://bit.ly/#frag", true},
	}
	for _, tt := range tests {
		out := &Output{OutputData: OutputData{URLIsShortener: true}}
		applyExceptions(&Input{URL: tt.url}, out)
		if out.URLIsShortener != tt.want {
			t.Errorf("%s: URLIsShortener = %v, want %v", tt.url, out.URLIsShortener, tt.want)
		}
	}
}

func TestApplyExceptions_RedirectToVouchedDomain(t *testing.T) {
	tests := []struct {
		name  string
		rank  int
		final string
		want  bool // HasDomainJump after exceptions
	}{
		{"bank moving to .bank.in", 7237, "https://www.hdfc.bank.in/", false},
		{"well-known site to unknown domain", 7237, "https://hdfc-rewards.xyz/", true},
		{"unranked site to .bank.in", 0, "https://www.hdfc.bank.in/", true},
	}
	for _, tt := range tests {
		out := &Output{OutputData: OutputData{
			Rank:              tt.rank,
			RedirectionResult: checks.RedirectionResult{IsRedirected: true, HasDomainJump: true, FinalURL: tt.final},
		}}
		applyExceptions(&Input{URL: "https://hdfcbank.com"}, out)
		if out.RedirectionResult.HasDomainJump != tt.want {
			t.Errorf("%s: HasDomainJump = %v, want %v", tt.name, out.RedirectionResult.HasDomainJump, tt.want)
		}
	}
}

func TestApplyExceptions_TrustedTLD(t *testing.T) {
	out := &Output{OutputData: OutputData{
		TLDTrusted:         true,
		TyposquatResult:    typosquat.TyposquatResult{IsSuspicious: true, MatchedDomain: "htsc.com"},
		URLKeywordsPresent: true,
	}}
	applyExceptions(&Input{URL: "https://hdfc.bank.in"}, out)
	if out.TyposquatResult.IsSuspicious || out.URLKeywordsPresent {
		t.Error("domains on restricted registries should skip name-based heuristics")
	}
}

func TestApplyExceptions_TrustedTLDBrandMention(t *testing.T) {
	content := &checks.PageFormResult{BrandCheck: checks.BrandResult{BrandFound: "Kotak", IsMismatch: true, DetectedNames: []string{"Kotak"}}}
	out := &Output{OutputData: OutputData{TLDTrusted: true, ContentData: content}}
	applyExceptions(&Input{URL: "https://kotak.bank.in"}, out)
	if out.ContentData.BrandCheck.IsMismatch {
		t.Error("a brand named on a restricted registry should not count as impersonation")
	}
	if !content.BrandCheck.IsMismatch {
		t.Error("applyExceptions should not modify the cached content result in place")
	}
}
