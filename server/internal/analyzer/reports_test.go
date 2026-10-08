package analyzer

import (
	"strings"
	"testing"
	"time"

	"github.com/urlvet/urlvet/internal/service/checks"
	"github.com/urlvet/urlvet/internal/service/domaininfo"
)

// Cases from user reports, built from what the scans actually returned.

func hasReason(reasons []string, sub string) bool {
	for _, r := range reasons {
		if strings.Contains(r, sub) {
			return true
		}
	}
	return false
}

// url.vet: five months old, unranked, HSTS and DNSSEC, nothing bad found.
func TestReport_NewCleanSiteIsSafe(t *testing.T) {
	resp := Response{
		URL:            "https://url.vet/",
		Features:       Features{TLD: TLDInfo{TLD: "vet", IsICANN: true}},
		Infrastructure: Infrastructure{NameserversValid: true, MXRecordsValid: true},
		DomainInfo:     &domaininfo.RegistrationData{AgeDays: days(151), CreatedDate: since(151), AgeHuman: "5 months", DNSSEC: true},
		Analysis:       Analysis{SupportsHSTS: true},
		ContentData:    &checks.PageFormResult{},
	}
	got := GenerateResult(resp)
	if got.Verdict != "Safe" || got.RiskScore != 0 {
		t.Errorf("verdict = %s (risk %d, final %d), want Safe with no risk; bad: %v",
			got.Verdict, got.RiskScore, got.FinalScore, got.Reasons.BadReasons)
	}
}

// The same site a few weeks after registration still reads as new.
func TestReport_UnrankedFreshDomainStaysPenalised(t *testing.T) {
	resp := Response{
		Features:       Features{TLD: TLDInfo{TLD: "com", IsICANN: true}},
		Infrastructure: Infrastructure{NameserversValid: true, MXRecordsValid: true},
		DomainInfo:     &domaininfo.RegistrationData{AgeDays: days(20), CreatedDate: since(20), AgeHuman: "20 days"},
	}
	got := GenerateResult(resp)
	if !hasReason(got.Reasons.BadReasons, "Very low traffic") || !hasReason(got.Reasons.BadReasons, "Registered 20 days ago") {
		t.Errorf("bad reasons = %v, want low traffic and a new-registration warning", got.Reasons.BadReasons)
	}
}

// github.com/<anyone>/releases/download/…/VMware-Workstation….exe
func TestReport_ProgramOnUploadHostCappedAtSuspicious(t *testing.T) {
	resp := safeBase()
	resp.URL = "https://github.com/201853910/VMwareWorkstation/releases/download/26H1/VMware-Workstation-Full-26H1u1-25688693.exe"
	resp.Features.Rank = 29
	resp.Analysis.SupportsHSTS = true
	resp.Analysis.RedirectionResult = checks.RedirectionResult{
		IsRedirected: true,
		FinalURL:     "https://release-assets.githubusercontent.com/github-production-release-asset/225538126/e6eb9975",
		FinalURLHost: "release-assets.githubusercontent.com",
	}
	resp.ContentData = &checks.PageFormResult{ContentType: "application/octet-stream", FileName: "VMware-Workstation-Full-26H1u1-25688693.exe"}
	resp.DomainInfo = &domaininfo.RegistrationData{AgeDays: days(6937), CreatedDate: since(6937), AgeHuman: "19 years"}

	got := GenerateResult(resp)
	if got.Verdict != "Suspicious" || got.FinalScore != safeThreshold-1 {
		t.Errorf("verdict = %s (final %d), want Suspicious at %d", got.Verdict, got.FinalScore, safeThreshold-1)
	}
	if !hasReason(got.Reasons.BadReasons, "uploaded to github.com") {
		t.Errorf("bad reasons = %v, want the upload-host reason", got.Reasons.BadReasons)
	}
}

func TestReport_ProgramFromPublisherIsNeutral(t *testing.T) {
	resp := safeBase()
	resp.URL = "https://download.example-vendor.com/setup.msi"
	got := GenerateResult(resp)
	if got.Verdict != "Safe" || !hasReason(got.Reasons.NeutralReasons, "downloads a program (.msi)") {
		t.Errorf("verdict = %s, neutral = %v; want Safe with a download note", got.Verdict, got.Reasons.NeutralReasons)
	}
}

func TestReport_ProgramFromUnknownSiteAddsRisk(t *testing.T) {
	resp := safeBase()
	resp.Features.Rank = 0
	resp.URL = "https://free-vmware-keys.example/VMware.exe"
	base := safeBase()
	base.Features.Rank = 0
	if got, was := GenerateResult(resp).RiskScore, GenerateResult(base).RiskScore; got-was != 40 {
		t.Errorf("program download added %d risk, want 40", got-was)
	}
}

// storage.googleapis.com page whose script sends visitors to http://<ip>/?<fragment>.
func TestReport_ScriptRedirectToIPOnCloudStorage(t *testing.T) {
	resp := Response{
		URL: "https://storage.googleapis.com/usales26/usales26.html",
		Features: Features{
			TLD: TLDInfo{TLD: "googleapis.com", IsHostingPlatform: true},
		},
		ContentData: &checks.PageFormResult{
			ScriptRedirect: &checks.ScriptRedirect{Target: "185.80.128.4", ToIP: true, CrossDomain: true},
		},
	}
	got := GenerateResult(resp)
	if got.Verdict != "Risky" {
		t.Errorf("verdict = %s (final %d), want Risky; bad: %v", got.Verdict, got.FinalScore, got.Reasons.BadReasons)
	}
	for _, unwanted := range []string{"Unregulated", "low traffic", "DNS configuration"} {
		if hasReason(got.Reasons.BadReasons, unwanted) {
			t.Errorf("hosting platform got penalised for %q", unwanted)
		}
	}
}

func TestReport_TopRankedDomainIgnoresLongDeepURL(t *testing.T) {
	out := &Output{}
	out.Rank = 29
	out.URLTooLong = true
	out.URLTooDeep = true
	applyExceptions(&Input{URL: "https://github.com/a/b/releases/download/v1/x.exe", Domain: "github.com"}, out)
	if out.URLTooLong || out.URLTooDeep {
		t.Errorf("long/deep flags kept on a top-ranked domain")
	}
}

func TestReport_SameBrandRedirectIsNotAJump(t *testing.T) {
	out := &Output{}
	out.Rank = 29
	out.RedirectionResult = checks.RedirectionResult{
		IsRedirected:  true,
		HasDomainJump: true,
		FinalURL:      "https://release-assets.githubusercontent.com/x",
		FinalURLHost:  "release-assets.githubusercontent.com",
	}
	applyExceptions(&Input{URL: "https://github.com/a/b", Domain: "github.com"}, out)
	if out.RedirectionResult.HasDomainJump {
		t.Errorf("github.com → githubusercontent.com still counted as a domain jump")
	}
}

func TestReport_ProviderEndpointTakesOwnerRank(t *testing.T) {
	saved := lookupRank
	defer func() { lookupRank = saved }()
	lookupRank = func(d string) int {
		if d == "googleapis.com" || d == "github.io" {
			return 6
		}
		return 0
	}

	tests := []struct {
		url, suffix string
		wantRank    int
	}{
		{"https://storage.googleapis.com/", "googleapis.com", 6},
		{"https://storage.googleapis.com/usales26/usales26.html", "googleapis.com", 0}, // an upload
		{"https://someone.github.io/", "github.io", 0},                                 // a user's site
	}
	for _, tt := range tests {
		out := &Output{}
		out.TLD = tt.suffix
		out.TLDIsHostingPlatform = true
		applyExceptions(&Input{URL: tt.url}, out)
		if out.Rank != tt.wantRank {
			t.Errorf("%s: rank = %d, want %d", tt.url, out.Rank, tt.wantRank)
		}
	}
}

func TestReport_UploadOnProviderHostSaysSo(t *testing.T) {
	resp := Response{
		URL:      "https://storage.googleapis.com/usales26/usales26.html",
		Features: Features{TLD: TLDInfo{TLD: "googleapis.com", IsHostingPlatform: true}},
	}
	got := GenerateResult(resp)
	if !hasReason(got.Reasons.NeutralReasons, "File hosted on Google Cloud Storage") {
		t.Errorf("neutral reasons = %v, want the upload-host wording", got.Reasons.NeutralReasons)
	}
}

// ssiyad.com: five and a half years old, DNSSEC, no HSTS, nothing bad found.
// History should outweigh a missing header.
func TestReport_OldCleanSiteWithoutHSTSIsSafe(t *testing.T) {
	resp := Response{
		URL:            "https://ssiyad.com/",
		Features:       Features{TLD: TLDInfo{TLD: "com", IsICANN: true}},
		Infrastructure: Infrastructure{NameserversValid: true, MXRecordsValid: true},
		DomainInfo:     &domaininfo.RegistrationData{AgeDays: days(1994), CreatedDate: since(1994), AgeHuman: "5 years 6 months", DNSSEC: true},
		ContentData:    &checks.PageFormResult{},
	}
	if got := GenerateResult(resp); got.Verdict != "Safe" {
		t.Errorf("verdict = %s (final %d, trust %d), want Safe", got.Verdict, got.FinalScore, got.TrustScore)
	}
}

func days(n int) *int { return &n }

func since(n int) time.Time { return time.Now().AddDate(0, 0, -n) }

// .de and .eu registries publish no creation date. That used to read as
// "2025 years" and earn the full age bonus.
func TestReport_UnpublishedCreationDateIsUnknown(t *testing.T) {
	resp := safeBase()
	resp.Features.Rank = 0
	resp.DomainInfo = &domaininfo.RegistrationData{}
	got := GenerateResult(resp)
	if hasReason(got.Reasons.GoodReasons, "history") || hasReason(got.Reasons.BadReasons, "Registered") {
		t.Errorf("unknown age scored as known: good %v, bad %v", got.Reasons.GoodReasons, got.Reasons.BadReasons)
	}
	if !hasReason(got.Reasons.BadReasons, "Very low traffic") {
		t.Errorf("unranked site with unknown age lost the low-traffic penalty: %v", got.Reasons.BadReasons)
	}

	// A result cached before the fix still carries the bogus day count.
	resp.DomainInfo = &domaininfo.RegistrationData{AgeDays: days(106751), AgeHuman: "2025 years 9 months"}
	if got := GenerateResult(resp); hasReason(got.Reasons.GoodReasons, "history") {
		t.Errorf("stale cached age counted: %v", got.Reasons.GoodReasons)
	}
}

func TestReport_AgeBands(t *testing.T) {
	tests := []struct {
		days      int
		wantRisk  int
		wantTrust int
	}{
		{0, 25, 0}, {30, 25, 0}, {31, 15, 0}, {90, 15, 0}, {91, 0, 0}, {365, 0, 0},
		{366, 0, 10}, {3 * 365, 0, 10}, {3*365 + 1, 0, 15}, {5 * 365, 0, 15}, {5*365 + 1, 0, 25},
	}
	for _, tt := range tests {
		resp := Response{
			Features:       Features{Rank: 500000, TLD: TLDInfo{IsICANN: true}},
			Infrastructure: Infrastructure{NameserversValid: true, MXRecordsValid: true},
			DomainInfo:     &domaininfo.RegistrationData{AgeDays: days(tt.days), CreatedDate: since(tt.days), AgeHuman: "x"},
		}
		got := GenerateResult(resp)
		// rank 500000 adds 20 trust and no risk on its own.
		if got.RiskScore != tt.wantRisk || got.TrustScore-20 != tt.wantTrust {
			t.Errorf("%d days: risk %d trust %d, want risk %d trust %d", tt.days, got.RiskScore, got.TrustScore-20, tt.wantRisk, tt.wantTrust)
		}
	}
}
