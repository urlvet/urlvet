package analyzer

import (
	"strings"
	"testing"

	"github.com/urlvet/urlvet/internal/service/checks"
	"github.com/urlvet/urlvet/internal/service/domaininfo"
)

// safeBase returns a minimal Response that should produce verdict "Safe".
// All boolean bad-signals are off; rank is top-10k; ICANN TLD; valid NS+MX.
func safeBase() Response {
	return Response{
		Features: Features{
			Rank: 5000,
			TLD:  TLDInfo{TLD: "com", IsTrusted: false, IsRisky: false, IsICANN: true},
			URL:  URLChecks{},
		},
		Infrastructure: Infrastructure{NameserversValid: true, MXRecordsValid: true},
		Analysis: Analysis{
			SupportsHSTS:      false,
			RedirectionResult: checks.RedirectionResult{},
		},
	}
}

func TestGenerateResult_Verdict(t *testing.T) {
	tests := []struct {
		name         string
		modify       func(*Response)
		wantVerdict  string
		wantRiskMin  int // risk score must be >= this
		wantTrustMin int // trust score must be >= this
	}{
		{
			name:         "top-ranked safe domain",
			modify:       func(r *Response) {},
			wantVerdict:  "Safe",
			wantTrustMin: 80,
		},
		{
			name: "unranked domain",
			modify: func(r *Response) {
				r.Features.Rank = 0
				// no rank trust (was +90), trust=0, risk=10
				// finalScore = 50 + (0-10)*0.5 = 45 → "Suspicious"
			},
			wantVerdict: "Suspicious",
			wantRiskMin: 10,
		},
		{
			name: "HSTS adds trust",
			modify: func(r *Response) {
				r.Analysis.SupportsHSTS = true
				// trust: 90(rank) + 10(NS) + 10(MX) + 20(HSTS) = 130 → 100
			},
			wantVerdict:  "Safe",
			wantTrustMin: 100,
		},
		{
			name: "URL shortener increases risk",
			modify: func(r *Response) {
				r.Features.URL.IsURLShortener = true
			},
			wantRiskMin: 25,
		},
		{
			// Top-10k rank trust (90) keeps final score high even with max risk.
			// Assert risk is maxed; verdict depends on combined formula.
			name: "raw IP usage - max risk signal",
			modify: func(r *Response) {
				r.Features.URL.UsesIP = true
			},
			wantRiskMin: 100,
		},
		{
			name: "punycode - max risk signal",
			modify: func(r *Response) {
				r.Features.URL.ContainsPunycode = true
			},
			wantRiskMin: 100,
		},
		{
			name: "homoglyph attack",
			modify: func(r *Response) {
				r.Features.URL.HasHomoglyph = true
			},
			wantRiskMin: 60,
		},
		{
			name: "risky TLD",
			modify: func(r *Response) {
				r.Features.TLD.IsRisky = true
			},
			wantRiskMin: 20,
		},
		{
			name: "trusted TLD adds trust",
			modify: func(r *Response) {
				r.Features.TLD.IsTrusted = true
			},
			wantTrustMin: 100,
		},
		{
			name: "non-ICANN TLD is bad",
			modify: func(r *Response) {
				r.Features.TLD.IsICANN = false
			},
			wantRiskMin: 30,
		},
		{
			name: "too-deep URL",
			modify: func(r *Response) {
				r.Features.URL.TooDeep = true
			},
			wantRiskMin: 30,
		},
		{
			name: "too-long URL",
			modify: func(r *Response) {
				r.Features.URL.TooLong = true
			},
			wantRiskMin: 20,
		},
		{
			name: "excessive subdomains",
			modify: func(r *Response) {
				r.Features.URL.SubdomainCount = 4
			},
			wantRiskMin: 15,
		},
		{
			name: "cross-domain redirect",
			modify: func(r *Response) {
				r.Analysis.RedirectionResult = checks.RedirectionResult{
					IsRedirected:  true,
					ChainLength:   2,
					HasDomainJump: true,
					FinalURLHost:  "evil.com",
				}
			},
			wantRiskMin: 50,
		},
		{
			// 4 hops across domains: excessive chain (+40) and domain jump (+50).
			name: "long cross-domain redirect chain",
			modify: func(r *Response) {
				r.Analysis.RedirectionResult = checks.RedirectionResult{
					IsRedirected:  true,
					ChainLength:   5,
					HasDomainJump: true,
					FinalURLHost:  "evil.com",
				}
			},
			wantRiskMin: 90,
		},
		{
			// valid=true + verified=true → confirmed phishing, max risk.
			name: "phishtank confirmed phishing (valid+verified)",
			modify: func(r *Response) {
				r.Phishing = &PhishingResult{
					InDatabase: true,
					Verified:   true,
					Valid:      true,
				}
			},
			wantRiskMin: 100,
		},
		{
			// valid=true + verified=false → reported but not yet reviewed.
			name: "phishtank reported phishing awaiting verification",
			modify: func(r *Response) {
				r.Phishing = &PhishingResult{
					InDatabase: true,
					Verified:   false,
					Valid:      true,
				}
			},
			wantRiskMin: 50,
		},
		{
			// valid=false → NOT phishing regardless of verified; no risk penalty.
			name: "phishtank in database but confirmed not phishing",
			modify: func(r *Response) {
				r.Phishing = &PhishingResult{
					InDatabase: true,
					Verified:   true,
					Valid:      false,
				}
			},
			wantRiskMin: 0,
		},
		{
			name: "new domain (<=30 days)",
			modify: func(r *Response) {
				r.DomainInfo = &domaininfo.RegistrationData{AgeDays: days(10), CreatedDate: since(10), AgeHuman: "10 days"}
			},
			wantRiskMin: 25,
		},
		{
			name: "young domain (<=90 days)",
			modify: func(r *Response) {
				r.DomainInfo = &domaininfo.RegistrationData{AgeDays: days(60), CreatedDate: since(60), AgeHuman: "2 months"}
			},
			wantRiskMin: 15,
		},

		{
			name: "old domain adds trust",
			modify: func(r *Response) {
				r.DomainInfo = &domaininfo.RegistrationData{AgeDays: days(3000), CreatedDate: since(3000), AgeHuman: "8 years"}
			},
			wantTrustMin: 90,
		},
		{
			name: "DNSSEC adds trust",
			modify: func(r *Response) {
				r.DomainInfo = &domaininfo.RegistrationData{AgeDays: days(3000), CreatedDate: since(3000), AgeHuman: "8 years", DNSSEC: true}
			},
			wantTrustMin: 100,
		},
		{
			name: "brand mismatch - critical",
			modify: func(r *Response) {
				r.ContentData = &checks.PageFormResult{
					BrandCheck: checks.BrandResult{IsMismatch: true, BrandFound: "PayPal"},
				}
			},
			wantRiskMin: 100,
		},
		{
			name: "login form on unranked domain",
			modify: func(r *Response) {
				r.Features.Rank = 0
				r.ContentData = &checks.PageFormResult{HasLoginForm: true}
			},
			wantRiskMin: 50,
		},
		{
			name: "login form on established domain is neutral",
			modify: func(r *Response) {
				r.ContentData = &checks.PageFormResult{HasLoginForm: true}
				// base has rank 5000 (established) → should not add risk for login form
			},
			wantVerdict: "Safe",
		},
		{
			name: "payment form warning",
			modify: func(r *Response) {
				r.ContentData = &checks.PageFormResult{HasPaymentForm: true}
			},
			wantRiskMin: 30,
		},
		{
			name: "hidden iframe",
			modify: func(r *Response) {
				r.ContentData = &checks.PageFormResult{HasHiddenIframe: true}
			},
			wantRiskMin: 40,
		},
		{
			name: "password form submits to external domain",
			modify: func(r *Response) {
				r.ContentData = &checks.PageFormResult{
					HasForms: true,
					Forms:    []checks.FormInfo{{ExternalAction: true, ContainsPassword: true}},
				}
			},
			wantRiskMin: 80,
		},
		{
			name: "email-only form submits to external domain",
			modify: func(r *Response) {
				r.ContentData = &checks.PageFormResult{
					HasForms: true,
					Forms:    []checks.FormInfo{{ExternalAction: true, ContainsUserLike: true}},
				}
			},
			wantRiskMin: 20,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp := safeBase()
			tc.modify(&resp)

			result := GenerateResult(resp)

			if tc.wantVerdict != "" && result.Verdict != tc.wantVerdict {
				t.Errorf("verdict = %q, want %q (risk=%d trust=%d final=%d)",
					result.Verdict, tc.wantVerdict, result.RiskScore, result.TrustScore, result.FinalScore)
			}
			if result.RiskScore < tc.wantRiskMin {
				t.Errorf("riskScore = %d, want >= %d", result.RiskScore, tc.wantRiskMin)
			}
			if result.TrustScore < tc.wantTrustMin {
				t.Errorf("trustScore = %d, want >= %d", result.TrustScore, tc.wantTrustMin)
			}
		})
	}
}

func TestGenerateResult_ScoreClamping(t *testing.T) {
	resp := safeBase()
	// pile on every risk signal
	resp.Features.URL.UsesIP = true                                                // +100
	resp.Features.URL.ContainsPunycode = true                                      // +100
	resp.Features.URL.HasHomoglyph = true                                          // +60
	resp.Phishing = &PhishingResult{InDatabase: true, Verified: true, Valid: true} // +100
	resp.ContentData = &checks.PageFormResult{
		BrandCheck: checks.BrandResult{IsMismatch: true, BrandFound: "PayPal"},
		HasForms:   true,
		Forms:      []checks.FormInfo{{ExternalAction: true}},
	}

	result := GenerateResult(resp)
	if result.RiskScore != 100 {
		t.Errorf("risk score should be clamped to 100, got %d", result.RiskScore)
	}
	if result.FinalScore < 0 || result.FinalScore > 100 {
		t.Errorf("final score out of range: %d", result.FinalScore)
	}
}

func TestGenerateResult_VerdictBoundaries(t *testing.T) {
	// finalScore = clamp(50 + (trust-risk)*0.5)
	// Build a response that lands exactly in each verdict bucket.

	// "Risky": finalScore < 30
	// rank=0 → risk=20 trust=0; NS invalid → risk+10; non-ICANN → risk+30; MX invalid → risk+5
	// total risk=65, trust=0; final = 50 + (0-65)*0.5 = 17.5 → 18 → Risky
	riskyResp := Response{
		Features: Features{
			Rank: 0,
			TLD:  TLDInfo{IsICANN: false},
		},
		Infrastructure: Infrastructure{},
	}
	if r := GenerateResult(riskyResp); r.Verdict != "Risky" {
		t.Errorf("expected Risky, got %s (final=%d)", r.Verdict, r.FinalScore)
	}

	// "Safe": finalScore >= 65
	// top-10k rank: trust=90; HSTS: +20 → trust=110→100; risk=0; final=50+50=100 → Safe
	safeResp := safeBase()
	safeResp.Analysis.SupportsHSTS = true
	if r := GenerateResult(safeResp); r.Verdict != "Safe" {
		t.Errorf("expected Safe, got %s (final=%d)", r.Verdict, r.FinalScore)
	}
}

func TestGenerateResult_SameSiteRedirectChainIsFine(t *testing.T) {
	// paypal.com → www.paypal.com → /in/home → … all on one domain.
	resp := safeBase()
	resp.Analysis.RedirectionResult = checks.RedirectionResult{
		IsRedirected: true,
		ChainLength:  5,
		FinalURLHost: "www.paypal.com",
	}
	base := GenerateResult(safeBase())
	got := GenerateResult(resp)
	if got.RiskScore != base.RiskScore {
		t.Errorf("same-site redirect chain added risk: %d, want %d", got.RiskScore, base.RiskScore)
	}
	for _, r := range got.Reasons.BadReasons {
		if strings.Contains(r, "redirection chain") {
			t.Errorf("unexpected bad reason %q", r)
		}
	}
}

func TestGenerateResult_ExternalSearchFormIsFine(t *testing.T) {
	// login.gov's site search sends the query to search.usa.gov.
	resp := safeBase()
	resp.ContentData = &checks.PageFormResult{
		HasForms: true,
		Forms:    []checks.FormInfo{{ExternalAction: true, Method: "GET", Action: "https://search.usa.gov/search"}},
	}
	base := GenerateResult(safeBase())
	if got := GenerateResult(resp); got.RiskScore != base.RiskScore {
		t.Errorf("external search form added risk: %d, want %d", got.RiskScore, base.RiskScore)
	}
}
