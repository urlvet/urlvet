package analyzer

import (
	"testing"

	"github.com/urlvet/urlvet/internal/service/typosquat"
)

func TestApplyExceptions_TopRankedDomain(t *testing.T) {
	out := &Output{
		Rank:            42000,
		TLDRisky:        true,
		TyposquatResult: typosquat.TyposquatResult{IsSuspicious: true, IsComboSquat: true, MatchedBrand: "paypal"},
	}
	applyExceptions(&Input{URL: "https://example.xyz"}, out)

	if out.TLDRisky {
		t.Error("risky TLD should be ignored for top-ranked domains")
	}
	if out.TyposquatResult.IsSuspicious || out.TyposquatResult.IsComboSquat {
		t.Error("typosquat should be ignored for top-ranked domains")
	}
}

func TestApplyExceptions_UnrankedDomainKeepsSignals(t *testing.T) {
	out := &Output{
		Rank:            0,
		TLDRisky:        true,
		TyposquatResult: typosquat.TyposquatResult{IsSuspicious: true},
	}
	applyExceptions(&Input{URL: "https://paypa1.xyz"}, out)

	if !out.TLDRisky || !out.TyposquatResult.IsSuspicious {
		t.Error("unranked domains must keep risky TLD and typosquat signals")
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
		out := &Output{URLIsShortener: true}
		applyExceptions(&Input{URL: tt.url}, out)
		if out.URLIsShortener != tt.want {
			t.Errorf("%s: URLIsShortener = %v, want %v", tt.url, out.URLIsShortener, tt.want)
		}
	}
}
