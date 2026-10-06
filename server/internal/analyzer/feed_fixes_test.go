package analyzer

import (
	"testing"
	"time"

	"github.com/urlvet/urlvet/internal/service/checks"
	"github.com/urlvet/urlvet/internal/service/domaininfo"
	"github.com/urlvet/urlvet/internal/service/threatfeeds"
	"github.com/urlvet/urlvet/internal/service/typosquat"
)

// Cases from scanning 300 OpenPhish URLs.

func hostingBase(url, suffix string) Response {
	return Response{
		URL:      url,
		Features: Features{TLD: TLDInfo{TLD: suffix, IsICANN: true, IsHostingPlatform: true}},
	}
}

func TestFeed_CloudflareInterstitialIsRisky(t *testing.T) {
	resp := hostingBase("http://ledgr--live.pages.dev/", "pages.dev")
	resp.ContentData = &checks.PageFormResult{ProviderBlock: &checks.ProviderBlock{Provider: "Cloudflare", Reason: "phishing"}}
	if got := GenerateResult(resp); got.Verdict != "Risky" || !hasReason(got.Reasons.BadReasons, "Cloudflare has flagged") {
		t.Errorf("verdict %s, bad %v", got.Verdict, got.Reasons.BadReasons)
	}
}

func TestFeed_TakenDownDeploymentAddsRisk(t *testing.T) {
	resp := hostingBase("https://posfihom.vercel.app/", "vercel.app")
	base := GenerateResult(resp).RiskScore
	resp.ContentData = &checks.PageFormResult{ProviderBlock: &checks.ProviderBlock{Provider: "Vercel", Reason: "blocked"}}
	if got := GenerateResult(resp); got.RiskScore-base != 40 || !hasReason(got.Reasons.BadReasons, "Vercel has taken this site down") {
		t.Errorf("risk +%d, bad %v", got.RiskScore-base, got.Reasons.BadReasons)
	}
}

func TestFeed_CertificateMismatch(t *testing.T) {
	resp := safeBase()
	resp.URL = "https://www.emmanuelbarrault.com/x"
	resp.TLSInfo = checks.TLSResult{Present: true, HostnameMismatch: true}
	if got := GenerateResult(resp); !hasReason(got.Reasons.BadReasons, "certificate doesn't belong") || got.RiskScore < 40 {
		t.Errorf("risk %d, bad %v", got.RiskScore, got.Reasons.BadReasons)
	}
	// Over plain HTTP nobody sees the certificate.
	resp.URL = "http://www.emmanuelbarrault.com/x"
	if got := GenerateResult(resp); hasReason(got.Reasons.BadReasons, "certificate") {
		t.Errorf("certificate flagged for an http:// link: %v", got.Reasons.BadReasons)
	}
}

func TestFeed_WordPressFolders(t *testing.T) {
	tests := map[string]bool{
		"https://hensonracingengines.com/wp-content/uploads/2026/10/alldomain-2-2.html": true,
		"https://www.emmanuelbarrault.com/wordpress/wp-content/plugins/hello-dolly/bb":  true,
		"https://www.conewagoventures.com/wp-admin/css/colors/":                         true,
		"https://example.com/wp-content/uploads/2026/10/photo.jpg":                      false,
		"https://example.com/wp-includes/js/jquery/jquery.min.js":                       false,
		"https://example.com/wp-admin/":                                                 false, // the real admin
		"https://example.com/wp-login.php":                                              false,
		"https://example.com/blog/post":                                                 false,
	}
	for url, want := range tests {
		if got := isWordPressSystemPage(url); got != want {
			t.Errorf("isWordPressSystemPage(%q) = %v, want %v", url, got, want)
		}
	}
}

func TestFeed_ComboSquatOnHostingWeighsMore(t *testing.T) {
	resp := hostingBase("http://www.instagram-blush.vercel.app/", "vercel.app")
	resp.TyposquatResult = typosquat.TyposquatResult{IsSuspicious: true, IsComboSquat: true, MatchedBrand: "instagram"}
	resp.ContentData = &checks.PageFormResult{Title: "Instagram fan art"}
	if got := GenerateResult(resp); got.RiskScore != 40 {
		t.Errorf("combo-squat on a hosting platform added %d risk, want 40", got.RiskScore)
	}
}

func TestFeed_UnresolvedShortLinkCapped(t *testing.T) {
	resp := safeBase()
	resp.URL = "https://bom.so/s2zyuwwezkn8"
	resp.ShortLink = &ShortLinkInfo{URL: resp.URL, Chain: []string{resp.URL}}
	if got := GenerateResult(resp); got.Verdict != "Suspicious" {
		t.Errorf("unresolved short link scored %s %d", got.Verdict, got.FinalScore)
	}
}

func TestFeed_BrandContext(t *testing.T) {
	// A bare "Facebook" title counts on an unknown site with a login form.
	out := &Output{}
	out.ContentData = &checks.PageFormResult{URL: "https://facebookbonus.emmanueldegreatltd.com/", Title: "Facebook", HasLoginForm: true}
	applyExceptions(&Input{URL: "https://facebookbonus.emmanueldegreatltd.com/"}, out)
	if !out.ContentData.BrandCheck.IsMismatch {
		t.Error("bare Facebook title with a login form not flagged")
	}

	// ...and on a free hosting subdomain, form or not.
	out = &Output{}
	out.TLDIsHostingPlatform = true
	out.ContentData = &checks.PageFormResult{URL: "https://pt-shopee-29.blogspot.com/", Title: "SHOPEE"}
	applyExceptions(&Input{URL: "https://pt-shopee-29.blogspot.com/"}, out)
	if !out.ContentData.BrandCheck.IsMismatch {
		t.Error("SHOPEE on blogspot not flagged")
	}

	// Not on an unknown site without a form: a blog post about Facebook.
	out = &Output{}
	out.ContentData = &checks.PageFormResult{URL: "https://smallblog.example/", Title: "Why I quit Facebook"}
	applyExceptions(&Input{URL: "https://smallblog.example/"}, out)
	if out.ContentData.BrandCheck.IsMismatch {
		t.Error("blog post title flagged")
	}

	// A well-known news site naming a brand isn't impersonating it.
	out = &Output{}
	out.Rank = 300
	out.ContentData = &checks.PageFormResult{URL: "https://news.example/", Title: "Netflix raises prices",
		BrandCheck: checks.CheckBrandMismatch("news.example", "Netflix raises prices")}
	applyExceptions(&Input{URL: "https://news.example/"}, out)
	if out.ContentData.BrandCheck.IsMismatch {
		t.Error("news headline on a well-known site flagged")
	}

	// But a well-known user-content host keeps the check.
	out = &Output{}
	out.Rank = 29
	out.ContentData = &checks.PageFormResult{URL: "https://sites.google.com/view/x", Title: "Netflix login",
		BrandCheck: checks.CheckBrandMismatch("sites.google.com", "Netflix login")}
	applyExceptions(&Input{URL: "https://sites.google.com/view/x"}, out)
	if !out.ContentData.BrandCheck.IsMismatch {
		t.Error("Netflix login page on Google Sites cleared")
	}
}

func TestFeed_LocalFeedListing(t *testing.T) {
	resp := hostingBase("https://candymesa.pages.dev/", "pages.dev")
	resp.ThreatFeeds = &threatfeeds.FeedMatch{Listed: true, Match: "url", Sources: []string{"PhishTank"}}
	if got := GenerateResult(resp); got.Verdict != "Risky" || !hasReason(got.Reasons.BadReasons, "listed as phishing or malware by PhishTank") {
		t.Errorf("verdict %s, bad %v", got.Verdict, got.Reasons.BadReasons)
	}

	hostMatch := func() *threatfeeds.FeedMatch {
		return &threatfeeds.FeedMatch{Listed: true, Match: "host", Sources: []string{"PhishTank"}}
	}
	// A customer's own hosting subdomain: the whole host is theirs.
	out0 := &Output{}
	out0.TLD, out0.TLDIsHostingPlatform = "pages.dev", true
	out0.ThreatFeeds = hostMatch()
	applyExceptions(&Input{URL: "https://candymesa.pages.dev/other"}, out0)
	if !out0.ThreatFeeds.Listed {
		t.Error("host-level listing dropped on a customer's hosting subdomain")
	}
	// An old shared service (a form builder): other customers' pages.
	out1 := &Output{}
	out1.DomainInfo = &domaininfo.RegistrationData{AgeDays: days(6000), CreatedDate: since(6000)}
	out1.ThreatFeeds = hostMatch()
	applyExceptions(&Input{URL: "http://www.formbuddy.com/some-other-page"}, out1)
	if out1.ThreatFeeds.Listed {
		t.Error("host-level listing kept on an old shared service")
	}
	// A domain registered weeks ago for the campaign.
	out2 := &Output{}
	out2.DomainInfo = &domaininfo.RegistrationData{AgeDays: days(20), CreatedDate: since(20)}
	out2.ThreatFeeds = hostMatch()
	applyExceptions(&Input{URL: "https://paypa1-login.top/other"}, out2)
	if !out2.ThreatFeeds.Listed {
		t.Error("host-level listing dropped on a fresh domain")
	}

	// A reported page elsewhere on github.com says nothing about this one.
	out := &Output{}
	out.Rank = 29
	out.ThreatFeeds = &threatfeeds.FeedMatch{Listed: true, Match: "host", Sources: []string{"PhishTank"}}
	applyExceptions(&Input{URL: "https://github.com/golang/go"}, out)
	if out.ThreatFeeds.Listed {
		t.Error("host-level listing kept on github.com")
	}
}

func TestFeed_WebRisk(t *testing.T) {
	resp := safeBase()
	resp.WebRisk = &threatfeeds.GoogleThreatResult{Listed: true, ThreatTypes: []string{"SOCIAL_ENGINEERING"}}
	if got := GenerateResult(resp); !hasReason(got.Reasons.BadReasons, "Google lists this link as phishing") {
		t.Errorf("bad %v", got.Reasons.BadReasons)
	}
}

func TestFeed_SafeBrowsing(t *testing.T) {
	resp := safeBase()
	resp.SafeBrowsing = &threatfeeds.GoogleThreatResult{Listed: true, ThreatTypes: []string{"MALWARE"}}
	got := GenerateResult(resp)
	if !hasReason(got.Reasons.BadReasons, "Google lists this link as malware") {
		t.Errorf("bad %v", got.Reasons.BadReasons)
	}
	// Listed by both Google sources: one reason, not two.
	resp.WebRisk = &threatfeeds.GoogleThreatResult{Listed: true, ThreatTypes: []string{"MALWARE"}}
	if got2 := GenerateResult(resp); got2.RiskScore != got.RiskScore {
		t.Errorf("both Google sources counted twice: %d vs %d", got2.RiskScore, got.RiskScore)
	}
}

// Cases from the second 300-URL run.

func TestRun2_BrandsInURLParts(t *testing.T) {
	base := func() Response {
		return Response{URL: "https://x.example/", Domain: "x.example", Features: Features{TLD: TLDInfo{TLD: "org", IsICANN: true}}, ContentData: &checks.PageFormResult{Title: "X Example"}}
	}
	tests := []struct {
		name string
		ts   typosquat.TyposquatResult
		risk int
		want string
	}{
		{"embedded domain", typosquat.TyposquatResult{EmbeddedDomain: "google.com", SubdomainBrand: "google", SubdomainBrandExact: true}, 60, "starts like google.com"},
		{"brand in a subdomain", typosquat.TyposquatResult{SubdomainBrand: "paypal"}, 40, "uses the brand name 'paypal'"},
		{"subdomain that is just the name", typosquat.TyposquatResult{SubdomainBrand: "github", SubdomainBrandExact: true}, 20, "is named 'github'"},
		{"brand in the path", typosquat.TyposquatResult{PathBrand: "paypal"}, 20, "path names 'paypal'"},
	}
	for _, tt := range tests {
		plain := GenerateResult(base())
		resp := base()
		resp.TyposquatResult = tt.ts
		got := GenerateResult(resp)
		if got.RiskScore-plain.RiskScore != tt.risk || !hasReason(got.Reasons.BadReasons, tt.want) {
			t.Errorf("%s: +%d risk, bad %v; want +%d and %q", tt.name, got.RiskScore-plain.RiskScore, got.Reasons.BadReasons, tt.risk, tt.want)
		}
	}

	// A login form makes a brand in the path count for more.
	resp := base()
	resp.TyposquatResult = typosquat.TyposquatResult{PathBrand: "correos"}
	resp.ContentData = &checks.PageFormResult{HasLoginForm: true}
	if got := GenerateResult(resp); !hasReason(got.Reasons.BadReasons, "asks you to log in") {
		t.Errorf("path brand with a login form: %v", got.Reasons.BadReasons)
	}

	// Well-known sites keep their own URLs: the exceptions clear these.
	out := &Output{}
	out.Rank = 50
	out.TyposquatResult = typosquat.TyposquatResult{SubdomainBrand: "github", PathBrand: "paypal", EmbeddedDomain: "google.com"}
	applyExceptions(&Input{URL: "https://github.blog/paypal"}, out)
	if out.TyposquatResult.SubdomainBrand != "" || out.TyposquatResult.PathBrand != "" || out.TyposquatResult.EmbeddedDomain != "" {
		t.Errorf("not cleared on a well-known site: %+v", out.TyposquatResult)
	}
}

func TestRun2_MidRankedBrandNames(t *testing.T) {
	combo := func(brand string, exact bool) typosquat.TyposquatResult {
		return typosquat.TyposquatResult{IsSuspicious: true, IsComboSquat: true, MatchedBrand: brand, ComboExact: exact}
	}
	age := func(d int) *domaininfo.RegistrationData {
		return &domaininfo.RegistrationData{AgeDays: days(d), CreatedDate: since(d)}
	}
	tests := []struct {
		name      string
		rank      int
		ts        typosquat.TyposquatResult
		di        *domaininfo.RegistrationData
		keywords  bool
		keepCombo bool
		keepSub   bool
	}{
		{"binanceuz.co: ranked, 1.6 years, /login/password", 722932, combo("binance", false), age(586), true, true, false},
		{"killedbygoogle.com: ranked, old, nothing else", 118270, combo("google", false), age(3000), false, false, false},
		{"roblox.dk: the plain name on another ending", 566203, combo("roblox", true), age(200), true, false, false},
		{"well-known site", 5000, combo("paypal", false), age(100), true, false, false},
		{"paypal-login subdomain on a ranked old domain", 119590, typosquat.TyposquatResult{SubdomainBrand: "paypal"}, age(2359), false, false, true},
		{"github.acme subdomain on a ranked old domain", 300000, typosquat.TyposquatResult{SubdomainBrand: "github", SubdomainBrandExact: true}, age(4000), false, false, false},
	}
	for _, tt := range tests {
		out := &Output{}
		out.Rank = tt.rank
		out.TyposquatResult = tt.ts
		out.DomainInfo = tt.di
		out.URLKeywordsPresent = tt.keywords
		applyExceptions(&Input{URL: "https://x.example/"}, out)
		if got := out.TyposquatResult.IsComboSquat; got != tt.keepCombo {
			t.Errorf("%s: combo kept = %v, want %v", tt.name, got, tt.keepCombo)
		}
		if got := out.TyposquatResult.SubdomainBrand != ""; got != tt.keepSub {
			t.Errorf("%s: subdomain brand kept = %v, want %v", tt.name, got, tt.keepSub)
		}
	}
}

func TestNewsletterToMarketingPlatform(t *testing.T) {
	newsletter := checks.FormInfo{ExternalAction: true, Action: "https://s88570519.t.eloqua.com/e/f2", ContainsUserLike: true, ContainsPersonal: true}
	resp := safeBase()
	resp.ContentData = &checks.PageFormResult{HasForms: true, Forms: []checks.FormInfo{newsletter}}
	if got := GenerateResult(resp); hasReason(got.Reasons.BadReasons, "different domain") {
		t.Errorf("newsletter to Eloqua flagged: %v", got.Reasons.BadReasons)
	}
	// A password going to the same platform is still out of place.
	newsletter.ContainsPassword = true
	resp.ContentData = &checks.PageFormResult{HasForms: true, Forms: []checks.FormInfo{newsletter}}
	if got := GenerateResult(resp); !hasReason(got.Reasons.BadReasons, "CRITICAL") {
		t.Errorf("password to a marketing platform not flagged: %v", got.Reasons.BadReasons)
	}
	// A phishing kit's form backend stays flagged.
	kit := checks.FormInfo{ExternalAction: true, Action: "https://formspree.io/f/abc", ContainsPersonal: true}
	resp.ContentData = &checks.PageFormResult{HasForms: true, Forms: []checks.FormInfo{kit}}
	if got := GenerateResult(resp); !hasReason(got.Reasons.BadReasons, "CRITICAL") {
		t.Errorf("personal details to formspree not flagged: %v", got.Reasons.BadReasons)
	}
}

// Cases from the third, unseen 300.

func TestRun3_BrandOnBarePage(t *testing.T) {
	sub := Response{URL: "https://outlook.verifytoken.com/s/63bzg", TyposquatResult: typosquat.TyposquatResult{SubdomainBrand: "outlook", SubdomainBrandExact: true}}
	combo := Response{URL: "https://admin.santandercitas.com", TyposquatResult: typosquat.TyposquatResult{IsSuspicious: true, IsComboSquat: true, MatchedBrand: "santander"}}
	for _, resp := range []Response{sub, combo} {
		resp.ContentData = &checks.PageFormResult{Title: "SantanderCitas — book an appointment"}
		base := GenerateResult(resp).RiskScore
		for why, page := range map[string]*checks.PageFormResult{
			"is blank":                  {},
			"is a bare file listing":    {Title: "Index of /"},
			"asks you to log in or pay": {Title: "Sign in", HasForms: true, HasLoginForm: true},
		} {
			resp.ContentData = page
			// A login form adds its own few points on an unranked site.
			if got := GenerateResult(resp); got.RiskScore-base < 20 || got.RiskScore-base > 25 || !hasReason(got.Reasons.BadReasons, why) {
				t.Errorf("%s, page %s: risk +%d, bad %v", resp.URL, why, got.RiskScore-base, got.Reasons.BadReasons)
			}
		}
		// A page that didn't load says nothing: brands park spare domains.
		resp.ContentData = nil
		if got := GenerateResult(resp); hasReason(got.Reasons.BadReasons, "little-known") {
			t.Errorf("%s: unloaded page counted as bare: %v", resp.URL, got.Reasons.BadReasons)
		}
		// A ranked site has a page of its own to judge.
		resp.ContentData, resp.Features.Rank = &checks.PageFormResult{}, 400000
		if got := GenerateResult(resp); hasReason(got.Reasons.BadReasons, "little-known") {
			t.Errorf("%s: ranked site counted as bare: %v", resp.URL, got.Reasons.BadReasons)
		}
	}
}

func TestCloudStorageRedirectHop(t *testing.T) {
	for suffix, want := range map[string]bool{"blob.core.windows.net": true, "s3.us-east-1.amazonaws.com": true, "web.core.windows.net": true, "github.io": true, "co.uk": false} {
		if got := isHostingSuffix(suffix); got != want {
			t.Errorf("isHostingSuffix(%q) = %v, want %v", suffix, got, want)
		}
	}
	resp := hostingBase("https://plh1288990.blob.core.windows.net/jud1288990/jtbluIh8.html", "blob.core.windows.net")
	resp.ContentData = &checks.PageFormResult{ScriptRedirect: &checks.ScriptRedirect{Target: "https://vmi3614668.contaboserver.net/gateway.php", CrossDomain: true}}
	got := GenerateResult(resp)
	if got.Verdict != "Risky" || !hasReason(got.Reasons.BadReasons, "anyone could upload") || !hasReason(got.Reasons.NeutralReasons, "Stored on blob.core.windows.net") || hasReason(got.Reasons.BadReasons, "Unregulated") {
		t.Errorf("verdict %s, bad %v, neutral %v", got.Verdict, got.Reasons.BadReasons, got.Reasons.NeutralReasons)
	}
	// A project page moving to its own domain is an ordinary redirect.
	resp = hostingBase("https://octocat.github.io/", "github.io")
	resp.ContentData = &checks.PageFormResult{ScriptRedirect: &checks.ScriptRedirect{Target: "https://octocat.dev/", CrossDomain: true}}
	if got := GenerateResult(resp); hasReason(got.Reasons.BadReasons, "anyone could upload") || got.RiskScore != 30 {
		t.Errorf("github.io redirect: risk %d, bad %v", got.RiskScore, got.Reasons.BadReasons)
	}
}

// Cases from the fourth 300 (OpenPhish and PhishTank).

func TestRun4_RankedSiteNamingABrandIsNotVerified(t *testing.T) {
	in := &Input{URL: "https://news.example/microsoft-outage"}
	out := &Output{}
	out.Rank = 3946
	out.ContentData = &checks.PageFormResult{URL: in.URL, Title: "Microsoft account security verification"}
	applyBrandContext(in, out)
	if bc := out.ContentData.BrandCheck; bc.IsMismatch || len(bc.DetectedNames) > 0 {
		t.Errorf("brand check = %+v, want neither a mismatch nor a verified match", bc)
	}
}

func TestRun4_ObjectStorage(t *testing.T) {
	for in, want := range map[string]string{
		"https://f005.backblazeb2.com/file/reimmm/20000k.html":          "f005.backblazeb2.com",
		"https://schroederseptic.s3.ca-east-006.backblazeb2.com/s.html": "schroederseptic.s3.ca-east-006.backblazeb2.com",
	} {
		if got, _ := checks.GetDomain(in); got != want {
			t.Errorf("GetDomain(%q) = %q, want %q", in, got, want)
		}
		if !isUserContentHost(hostOf(in)) {
			t.Errorf("%s not a user-content host", in)
		}
	}
	if r := checks.CheckBrandMismatch("ilhamnet.com", "ConstructConnect Log In"); !r.IsMismatch {
		t.Error("ConstructConnect login clone not flagged")
	}
}

func TestRun4_ConfirmedLinkIsAlwaysRisky(t *testing.T) {
	resp := safeBase()
	resp.Features.Rank = 50642
	resp.ThreatFeeds = &threatfeeds.FeedMatch{Listed: true, Match: "url", Sources: []string{"PhishTank"}}
	if got := GenerateResult(resp); got.Verdict != "Risky" {
		t.Errorf("listed link on a trusted site = %s %d", got.Verdict, got.FinalScore)
	}
	// Other pages on the site being listed is a warning, not a verdict.
	resp.ThreatFeeds.Match = "host"
	if got := GenerateResult(resp); got.Verdict == "Risky" && got.FinalScore == riskyThreshold-1 {
		t.Errorf("host match forced to Risky")
	}
}

// Cases from the 1000 OpenPhish links.

func TestRun5_IPHostnameIsLikeRawIP(t *testing.T) {
	resp := safeBase()
	resp.URL = "https://webmail.50-6-22-93.nip.io/desktop/"
	if got := GenerateResult(resp); !hasReason(got.Reasons.BadReasons, "bare server (50.6.22.93)") || got.Verdict == "Safe" {
		t.Errorf("verdict %s, bad %v", got.Verdict, got.Reasons.BadReasons)
	}
}

func TestRun5_MailLoginOnFreeHosting(t *testing.T) {
	login := func(title string) Response {
		resp := hostingBase("https://milkyfloors.vercel.app/", "vercel.app")
		resp.ContentData = &checks.PageFormResult{Title: title, HasForms: true, HasLoginForm: true}
		return resp
	}
	for _, title := range []string{"Mail Authentication", "PORTAL - Mail Authentication", "Sign into Webmail", "Password Protection", "Verify your email"} {
		if got := GenerateResult(login(title)); !hasReason(got.Reasons.BadReasons, "email sign-in page") {
			t.Errorf("%q not flagged: %v", title, got.Reasons.BadReasons)
		}
	}
	if got := GenerateResult(login("Dashboard login")); hasReason(got.Reasons.BadReasons, "email sign-in page") {
		t.Errorf("ordinary app login flagged")
	}
}

func TestRun5_VictimAddressInLink(t *testing.T) {
	tests := map[string]bool{
		"https://yuuinnovations.com/wk/indexx.html#ario@b705f5a01194e17bb1b9d0190848e3ceab53.net": true,
		"https://x.example/#YWxpY2VAZXhhbXBsZS5jb20=":                                             true, // alice@example.com
		"https://x.example/login?u=YWxpY2VAZXhhbXBsZS5jb20":                                       true,
		"https://ubu-two.vercel.app/#[[-Email-]]":                                                 true,
		"https://news.example/unsubscribe?email=alice@example.com":                                false,
		"https://annapolischarterfishing.com/cpsess/x/?uid=corar7@0689512e875.net":                true,
		"https://x.example/sign.html?eta=galalz@0aeee4a201d7.org":                                 true,
		"https://shop.example/account/reset?user=alice@example.com":                               false,
		"https://news.example/manage?id=alice@example.com&list=3":                                 true,
		"https://x.example/#section-2":                                                            false,
		"https://x.example/?ref=abcdefgh12345678":                                                 false,
	}
	for in, want := range tests {
		if got := victimAddressInLink(in); got != want {
			t.Errorf("victimAddressInLink(%q) = %v, want %v", in, got, want)
		}
	}
}

// From matching the phishing feeds against the top million.

func TestMined_HostingAgeIsTheProviders(t *testing.T) {
	resp := hostingBase("https://kg5.04d.mytemp.website/", "mytemp.website")
	days := 1006
	resp.DomainInfo = &domaininfo.RegistrationData{AgeDays: &days, AgeKnown: true, CreatedDate: time.Now().AddDate(0, 0, -days)}
	if got := GenerateResult(resp); hasReason(got.Reasons.GoodReasons, "Operational") || !hasReason(got.Reasons.NeutralReasons, "mytemp.website's, not this site's") {
		t.Errorf("good %v, neutral %v", got.Reasons.GoodReasons, got.Reasons.NeutralReasons)
	}
}

func TestMined_UserPageLosesServiceRank(t *testing.T) {
	in := &Input{URL: "https://linktr.ee/att-login"}
	out := &Output{}
	out.Rank = 272
	applyExceptions(in, out)
	if out.Rank != 0 || !out.TLDIsHostingPlatform || out.TLD != "linktr.ee" {
		t.Errorf("rank %d, hosting %v, tld %q", out.Rank, out.TLDIsHostingPlatform, out.TLD)
	}
	in, out = &Input{URL: "https://linktr.ee/"}, &Output{}
	out.Rank = 272
	applyExceptions(in, out)
	if out.Rank != 272 {
		t.Error("Linktree's own homepage lost its rank")
	}
}

func TestMined_PanelServerNames(t *testing.T) {
	for host, ip := range map[string]string{
		"login.50-6-193-131.cprapid.com":                "50.6.193.131",
		"anmelden-dkb-auth-de.91-218-65-223.plesk.page": "91.218.65.223",
		"102.175.153.160.host.secureserver.net":         "102.175.153.160",
		"focused-ivory-owl.31-22-7-7.cpanel.site":       "31.22.7.7",
	} {
		if _, got, ok := checks.IPHostname(host); !ok || got != ip {
			t.Errorf("IPHostname(%q) = %q %v, want %s", host, got, ok, ip)
		}
	}
}

func TestMined_LinkInBioNamesNetworks(t *testing.T) {
	page := func(login bool) *Output {
		out := &Output{}
		out.ContentData = &checks.PageFormResult{URL: "https://linktr.ee/selenagomez", Title: "Selena Gomez Songs, Instagram & Twitter Links | Linktree", HasLoginForm: login}
		applyBrandContext(&Input{URL: "https://linktr.ee/selenagomez"}, out)
		return out
	}
	if page(false).ContentData.BrandCheck.IsMismatch {
		t.Error("Linktree page naming Instagram flagged")
	}
	if !page(true).ContentData.BrandCheck.IsMismatch {
		t.Error("Instagram login on a Linktree page not flagged")
	}
}

func TestBrandHostingRedirectsToBrand(t *testing.T) {
	jump := func(source, platform, final string) bool {
		in := &Input{URL: source}
		out := &Output{}
		out.TLD, out.TLDIsHostingPlatform = platform, true
		out.RedirectionResult.HasDomainJump = true
		out.RedirectionResult.FinalURL = "https://" + final + "/"
		out.RedirectionResult.FinalURLHost = final
		applyExceptions(in, out)
		return out.RedirectionResult.HasDomainJump
	}
	if jump("https://googleresearch.blogspot.com/", "blogspot.com", "research.google") {
		t.Error("Blogspot page moving to research.google counted as a jump")
	}
	if jump("https://someone.blogspot.de/", "blogspot.de", "blog.google") {
		t.Error("country Blogspot page moving to Google counted as a jump")
	}
	if !jump("https://someone.blogspot.com/", "blogspot.com", "evil.example") {
		t.Error("Blogspot page jumping elsewhere not counted")
	}
	if !jump("https://someone.vercel.app/", "vercel.app", "research.google") {
		t.Error("a host Google doesn't run vouched for a jump to Google")
	}
}

func TestAlsoOfficialDomains(t *testing.T) {
	// The brand's other domains are its own…
	if r := checks.CheckBrandMismatch("www.primevideo.com", "Amazon Prime Video sign in"); r.IsMismatch {
		t.Errorf("primevideo.com flagged as impersonating %s", r.BrandFound)
	}
	// …but their names aren't brand names: "research" is no Google.
	if r := typosquat.CheckTyposquatting("research-lab.example"); r.IsComboSquat {
		t.Errorf("research-lab.example taken for a combo of %s", r.MatchedBrand)
	}
	// A customer's Blogspot page isn't Google's own.
	if r := checks.CheckBrandMismatch("google-support-help.blogspot.com", "Google Account sign in"); !r.IsMismatch {
		t.Error("Google login on a Blogspot page not flagged")
	}
}

// Cases from run 6 (CERT Polska, Phishing Army, Phishing.Database, SinkingYachts).

func TestRun6_BlogspotToGoogleIsGoogles(t *testing.T) {
	in := &Input{URL: "https://googleresearch.blogspot.com/"}
	out := &Output{}
	out.TLD, out.TLDIsHostingPlatform = "blogspot.com", true
	out.RedirectionResult.HasDomainJump = true
	out.RedirectionResult.FinalURL = "https://research.google/blog/"
	out.RedirectionResult.FinalURLHost = "research.google"
	out.TyposquatResult = typosquat.TyposquatResult{IsSuspicious: true, IsComboSquat: true, MatchedBrand: "google", MatchedDomain: "google.com"}
	out.ContentData = &checks.PageFormResult{URL: in.URL, Title: "Google Research Blog"}
	applyExceptions(in, out)
	if out.RedirectionResult.HasDomainJump || out.TyposquatResult.IsComboSquat || out.ContentData.BrandCheck.IsMismatch {
		t.Errorf("jump %v, combo %v, mismatch %+v", out.RedirectionResult.HasDomainJump, out.TyposquatResult.IsComboSquat, out.ContentData.BrandCheck)
	}
}

func TestRun6_HackedGovLoginKit(t *testing.T) {
	scan := func(url string) Result {
		in := &Input{URL: url}
		out := &Output{}
		out.TLDTrusted = true
		out.TyposquatResult = typosquat.TyposquatResult{PathBrand: "netflix"}
		applyExceptions(in, out)
		resp := safeBase()
		resp.URL = url
		resp.Features.TLD.IsTrusted = true
		resp.TyposquatResult = out.TyposquatResult
		return GenerateResult(resp)
	}
	if got := scan("https://prisa.ins.gob.pe/secure/netflix/login"); got.Verdict == "Safe" || !hasReason(got.Reasons.BadReasons, "government or university") {
		t.Errorf("hacked .gob.pe kit: %s %d %v", got.Verdict, got.FinalScore, got.Reasons.BadReasons)
	}
	// A university page about Netflix isn't a login kit.
	if got := scan("https://cs.example.edu/research/netflix-prize"); hasReason(got.Reasons.BadReasons, "government or university") {
		t.Errorf("ordinary .edu page flagged: %v", got.Reasons.BadReasons)
	}
}

func TestRun6_TradingScamTitle(t *testing.T) {
	for title, want := range map[string]bool{
		"Polriksobotiks | AI Trading Platform | Official Website":  true,
		"Frame 2u Avapro | AI Trading Platform | Official Gateway": true,
		"Immediate Peak ™ | The Official & Updated Site":           true,
		"Interactive Brokers: Trading platform for professionals":  false,
		"Immediate dental care in Leeds":                           false,
	} {
		if got := tradingScamTitle.MatchString(title); got != want {
			t.Errorf("tradingScamTitle(%q) = %v, want %v", title, got, want)
		}
	}
}

func TestRun6_MarkedBase64Address(t *testing.T) {
	// "$" + base64("mandy_torres@sympatico.ca")
	if !victimAddressInLink("https://www.example.cz/wp-datas/index.php?emb=$bWFuZHlfdG9ycmVzQHN5bXBhdGljby5jYQ==") {
		t.Error("$-marked base64 address missed")
	}
}

func TestRun6_BlogspotToGoogleIsSafe(t *testing.T) {
	resp := hostingBase("https://googleresearch.blogspot.com/", "blogspot.com")
	resp.Analysis.RedirectionResult.FinalURL = "https://research.google/blog/"
	resp.Analysis.RedirectionResult.FinalURLHost = "research.google"
	saved := rankLookupForVouch
	rankLookupForVouch = func(string) int { return 4000 }
	defer func() { rankLookupForVouch = saved }()
	if got := GenerateResult(resp); got.Verdict != "Safe" || !hasReason(got.Reasons.GoodReasons, "run by the same company") {
		t.Errorf("%s %d %v", got.Verdict, got.FinalScore, got.Reasons.GoodReasons)
	}
	// Another host's page forwarding to Google isn't vouched for by Google.
	resp = hostingBase("https://someone.vercel.app/", "vercel.app")
	resp.Analysis.RedirectionResult.FinalURL = "https://research.google/blog/"
	resp.Analysis.RedirectionResult.FinalURLHost = "research.google"
	if got := GenerateResult(resp); hasReason(got.Reasons.GoodReasons, "run by the same company") {
		t.Error("vercel.app vouched for by Google")
	}
}

func TestRun6_NoHSTSTrustOnBrandNewDomain(t *testing.T) {
	resp := safeBase()
	resp.Analysis.SupportsHSTS = true
	days := 0
	resp.DomainInfo = &domaininfo.RegistrationData{AgeDays: &days, AgeKnown: true, CreatedDate: time.Now()}
	if got := GenerateResult(resp); hasReason(got.Reasons.GoodReasons, "HSTS") {
		t.Errorf("HSTS counted on a domain registered today: %v", got.Reasons.GoodReasons)
	}
	days = 400
	resp.DomainInfo = &domaininfo.RegistrationData{AgeDays: &days, AgeKnown: true, CreatedDate: time.Now().AddDate(0, 0, -days)}
	if got := GenerateResult(resp); !hasReason(got.Reasons.GoodReasons, "HSTS") {
		t.Error("HSTS not counted on an established domain")
	}
}

func TestRun6_KitLoginPath(t *testing.T) {
	for url, want := range map[string]bool{
		"https://nikkimajor.com/otpservice/sign_in.php?session=6431613939363965":            true,
		"https://econsave.com.my/Owa/logon.php":                                             true,
		"https://vincentjanse.nl/uae/uae/authentification/colis=116/d4308c18788c8ec31ec5a":  true,
		"https://shop.example/customer/account/login/referer/aHR0cHM6Ly9zaG9wLmV4YW1wbGUv/": false,
		"https://shop.example/account/reset?token=6431613939363965abcdef":                   false,
		"https://mail.example.com/owa/auth/logon.aspx":                                      false,
		"https://blog.example/2019/10/how-we-secured-our-login":                             false,
	} {
		if got := kitLoginPath(url); got != want {
			t.Errorf("kitLoginPath(%q) = %v, want %v", url, got, want)
		}
	}
	resp := safeBase()
	resp.URL = "https://nikkimajor.com/otpservice/sign_in.php?session=6431613939363965"
	resp.Features.Rank = 0
	days := 4000
	resp.DomainInfo = &domaininfo.RegistrationData{AgeDays: &days, AgeKnown: true, CreatedDate: time.Now().AddDate(0, 0, -days)}
	if got := GenerateResult(resp); hasReason(got.Reasons.GoodReasons, "Long-standing") {
		t.Errorf("age counted for a planted sign-in page: %v", got.Reasons.GoodReasons)
	}
}
