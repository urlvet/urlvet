package analyzer

import (
	"testing"
	"time"

	"github.com/urlvet/urlvet/internal/service/checks"
	"github.com/urlvet/urlvet/internal/service/domaininfo"
	"github.com/urlvet/urlvet/internal/service/threatfeeds"
	"github.com/urlvet/urlvet/internal/service/typosquat"
)

func TestRedirectTarget(t *testing.T) {
	var resp Response
	if redirectTarget(resp) != "" {
		t.Error("no redirect, but a target")
	}
	resp.Analysis.RedirectionResult = checks.RedirectionResult{HasDomainJump: true, FinalURL: "https://dest.example/"}
	if got := redirectTarget(resp); got != "https://dest.example/" {
		t.Errorf("HTTP jump: %q", got)
	}
	resp.Analysis.RedirectionResult = checks.RedirectionResult{}
	resp.ContentData = &checks.PageFormResult{ScriptRedirect: &checks.ScriptRedirect{Target: "https://dest.example/x", CrossDomain: true}}
	if got := redirectTarget(resp); got != "https://dest.example/x" {
		t.Errorf("script redirect: %q", got)
	}
}

func TestOriginFindings_HonestRedirectAddsNothing(t *testing.T) {
	src := safeBase()
	src.URL = "https://oldname.example/"
	if o := originFindings(src); o.Risk != 0 || len(o.BadReasons) != 0 {
		t.Errorf("an honest old domain forwarding added %d: %v", o.Risk, o.BadReasons)
	}
}

func TestOriginFindings_LuresCount(t *testing.T) {
	src := Response{URL: "https://paypal-login.example/verify#alice@example.com"}
	src.TyposquatResult = typosquat.TyposquatResult{IsSuspicious: true, IsComboSquat: true, MatchedBrand: "paypal"}
	src.Features.URL.Keywords = Keywords{HasKeywords: true, Found: []string{"verify"}}
	days := 2
	src.DomainInfo = &domaininfo.RegistrationData{AgeDays: &days, AgeKnown: true, CreatedDate: time.Now().AddDate(0, 0, -days), AgeHuman: "2 days"}
	o := originFindings(src)
	for _, want := range []string{"brand name 'paypal'", "words like verify", "email address", "registered 2 days ago"} {
		if !hasReason(o.BadReasons, want) {
			t.Errorf("missing %q in %v", want, o.BadReasons)
		}
	}
	if o.Confirmed {
		t.Error("confirmed without a listing")
	}
}

func TestOriginFindings_ListedLinkIsConfirmed(t *testing.T) {
	src := Response{URL: "https://evil.example/x"}
	src.ThreatFeeds = &threatfeeds.FeedMatch{Listed: true, Match: "url", Sources: []string{"PhishTank"}}
	o := originFindings(src)
	if !o.Confirmed || o.Risk < 200 {
		t.Errorf("listed link: %+v", o)
	}
	// Even when it lands somewhere well known, it's Risky.
	dest := safeBase()
	dest.Origin = o
	if got := GenerateResult(dest); got.Verdict != "Risky" {
		t.Errorf("listed link landing on a good site = %s %d", got.Verdict, got.FinalScore)
	}
}

func TestOrigin_TrustComesFromDestination(t *testing.T) {
	// A trusted site forwarding to an unknown one: the unknown one's result.
	dest := Response{URL: "https://unknown.example/"}
	dest.Origin = originFindings(safeBase())
	if got := GenerateResult(dest); got.Verdict == "Safe" || got.TrustScore != 0 {
		t.Errorf("destination borrowed trust: %s trust %d", got.Verdict, got.TrustScore)
	}
	// An honest forward to a good site stays good.
	good := safeBase()
	good.Origin = &OriginFindings{URL: "https://oldname.example/"}
	if got := GenerateResult(good); got.Verdict != "Safe" {
		t.Errorf("honest forward = %s %d", got.Verdict, got.FinalScore)
	}
}

func TestHopFindings(t *testing.T) {
	o := &OriginFindings{}
	addHopFindings(o, []string{
		"https://start.example/",
		"http://203.0.113.9/r",
		"https://login.50-6-22-93.nip.io/x",
		"https://start.example/again", // same site as the start
		"https://end.example/",
	}, "https://start.example/", "https://end.example/")
	if !hasReason(o.BadReasons, "bare IP address (203.0.113.9)") || !hasReason(o.BadReasons, "50.6.22.93") {
		t.Errorf("hops: %v", o.BadReasons)
	}
	if len(o.BadReasons) != 2 {
		t.Errorf("same-site hop counted: %v", o.BadReasons)
	}
}

func TestWorthFollowing_Cloaking(t *testing.T) {
	saved := rankLookupForVouch
	rankLookupForVouch = func(d string) int { return map[string]int{"google.com": 1, "amazon.com": 10}[d] }
	defer func() { rankLookupForVouch = saved }()

	unknown := Response{URL: "https://oferta413.sbs/"}
	if worthFollowing(unknown, "https://www.google.com/") {
		t.Error("unknown site landing on google.com followed: that's cloaking")
	}
	if !worthFollowing(unknown, "https://ebeb.life/sec/") {
		t.Error("unknown site landing on another unknown one not followed")
	}
	ranked := Response{URL: "https://youtube.ng/"}
	ranked.Features.Rank = 400000
	if !worthFollowing(ranked, "https://www.google.com/") {
		t.Error("ranked site's redirect not followed")
	}
}

func TestSameBlog(t *testing.T) {
	if !sameBlog("sbdrg.blogspot.fi", "sbdrg.blogspot.com") || !sameBlog("www.x.blogspot.de", "x.blogspot.com") {
		t.Error("one blog on two Blogger addresses counted as two sites")
	}
	if sameBlog("a.blogspot.com", "b.blogspot.com") || sameBlog("a.blogspot.com", "a.example.com") {
		t.Error("different sites taken for one blog")
	}
}
