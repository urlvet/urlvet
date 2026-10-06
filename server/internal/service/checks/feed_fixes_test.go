package checks

import (
	"net/http"
	"net/url"
	"testing"
)

func TestGetDomain_CustomerSites(t *testing.T) {
	tests := map[string]string{
		"https://xfinityteamslogin.weebly.com/":       "xfinityteamslogin.weebly.com",
		"http://www.site-91q65jzg7.godaddysites.com/": "site-91q65jzg7.godaddysites.com",
		"https://hxfvuiyoiug.b-cdn.net/all1.html":     "hxfvuiyoiug.b-cdn.net",
		"https://www.weebly.com/login":                "weebly.com", // the provider itself
		"https://weebly.com/":                         "weebly.com",
		"https://www.paypal.com/":                     "paypal.com",
		"https://a.b.cf234276.tw1.ru/x":               "a.b.cf234276.tw1.ru",
		"https://1286524.us23.myftpupload.com/x":      "1286524.us23.myftpupload.com", // ID first, region second
		"https://vmi3580221.contaboserver.net/rd/":    "vmi3580221.contaboserver.net",
		"https://startio-trexor.squarespace.com":      "startio-trexor.squarespace.com",
		"https://account.squarespace.com/login":       "squarespace.com", // Squarespace's own
		"https://static1.squarespace.com/x.png":       "squarespace.com",
		"http://wellsfargosecurelogin.ct.ws":          "wellsfargosecurelogin.ct.ws",
	}
	for in, want := range tests {
		if got, _ := GetDomain(in); got != want {
			t.Errorf("GetDomain(%q) = %q, want %q", in, got, want)
		}
	}
	if CustomerSitePlatform("x.weebly.com") != "weebly.com" || CustomerSitePlatform("weebly.com") != "" {
		t.Error("CustomerSitePlatform misidentified the provider")
	}
}

func TestDetectProviderBlock(t *testing.T) {
	resp := &http.Response{StatusCode: 403, Header: http.Header{}}
	if b := detectProviderBlock(resp, "Suspected Phishing | Cloudflare"); b == nil || b.Provider != "Cloudflare" || b.Reason != "phishing" {
		t.Errorf("Cloudflare phishing interstitial = %+v", b)
	}
	if b := detectProviderBlock(resp, "Suspected Malware | Cloudflare"); b == nil || b.Reason != "malware" {
		t.Errorf("Cloudflare malware interstitial = %+v", b)
	}
	if b := detectProviderBlock(resp, "Attention Required! | Cloudflare"); b != nil {
		t.Errorf("ordinary Cloudflare challenge counted as a block: %+v", b)
	}
	resp = &http.Response{StatusCode: 451, Header: http.Header{"Server": {"Vercel"}}}
	if b := detectProviderBlock(resp, "Deployment Unavailable"); b == nil || b.Provider != "Vercel" || b.Reason != "blocked" {
		t.Errorf("Vercel 451 = %+v", b)
	}
}

func TestPageRedirectTarget(t *testing.T) {
	base, _ := url.Parse("https://goo.su/TIcsAp")
	tests := map[string]string{
		// goo.su keeps the destination in a data attribute
		`<div id="delay-page" data-delay="5" data-url="https://www.roblox.com.bn/users/1/profile">Redirecting...</div>`: "https://www.roblox.com.bn/users/1/profile",
		`<meta http-equiv="refresh" content="0; url=https://evil.example/x">`:                                           "https://evil.example/x",
		`<script>window.location.href = "https://evil.example/y";</script>`:                                             "https://evil.example/y",
		`<a href="https://other.example/">a link is not a redirect</a>`:                                                 "",
		`<div data-url="/same-site/path">same host</div>`:                                                               "",
	}
	for body, want := range tests {
		if got := pageRedirectTarget(body, base); got != want {
			t.Errorf("pageRedirectTarget(%q) = %q, want %q", body, got, want)
		}
	}
}

func TestIsShortLink(t *testing.T) {
	tests := map[string]bool{
		"https://bit.ly/abc":                true,
		"https://bit.ly/":                   false, // the shortener's homepage
		"https://goo.su/TIcsAp":             true,
		"https://u.gy/q2DGHw":               true,
		"https://example.com/abc":           false,
		"https://flowcode.com/p/etKFshjys9": true,  // a Flowcode QR link
		"https://www.flowcode.com/pricing":  false, // Flowcode's own page
	}
	for in, want := range tests {
		if got := IsShortLink(in); got != want {
			t.Errorf("IsShortLink(%q) = %v, want %v", in, got, want)
		}
	}
}

func TestCheckBrandNames(t *testing.T) {
	if r := CheckBrandNames("facebookbonus.emmanueldegreatltd.com", "Facebook"); !r.IsMismatch || r.BrandFound != "Facebook" {
		t.Errorf("bare Facebook title = %+v", r)
	}
	if r := CheckBrandNames("www.facebook.com", "Facebook"); r.IsMismatch {
		t.Errorf("Facebook's own site flagged: %+v", r)
	}
	if r := CheckBrandNames("applegiftreward.vercel.app", "Apple — Claim Your $1,000 Apple Gift Card"); !r.IsMismatch || r.BrandFound != "Apple" {
		t.Errorf("Apple gift card title = %+v", r)
	}
}

func TestBrandMismatch_NewSectors(t *testing.T) {
	tests := []struct{ host, title, brand string }{
		{"vavgxahwqq.vercel.app", "Crédito en Línea - Banco Pichincha", "Banco Pichincha"},
		{"moltravi-zekun39471628.vercel.app", "Meta for Business - Page Appeal", "Meta"},
		{"proxy-test-001.pages.dev", "Official Apple Support", "Apple"},
		{"app-exodusweb.pages.dev", "Exodus Web3 Wallet - Your Gateway", "Exodus"},
		{"tiktokshop-byte.github.io", "TikTok Shop", "TikTok Shop"},
	}
	for _, tt := range tests {
		if r := CheckBrandMismatch(tt.host, tt.title); !r.IsMismatch || r.BrandFound != tt.brand {
			t.Errorf("CheckBrandMismatch(%q, %q) = %+v, want %s", tt.host, tt.title, r, tt.brand)
		}
	}
	// Bare words small sites use stay out of the strict check.
	for _, title := range []string{"Portal Ayuntamiento Santander", "Amazon Rainforest Tours", "Medicare Plans in Florida", "Social Security Disability Lawyer", "Leather Messenger Bags", "Vodafone Shop Berlin", "Minecraft Server List", "Wise words for 2026", "Steam cleaning services"} {
		if r := CheckBrandMismatch("small-site.example", title); r.IsMismatch {
			t.Errorf("strict check flagged %q as %s", title, r.BrandFound)
		}
	}
}

func TestBrandMismatch_LookalikeLetters(t *testing.T) {
	for _, title := range []string{"Çoinbase Pro: Login", "Pаypal Login" /* Cyrillic а */, "Nеtflix" /* Cyrillic е */} {
		if r := CheckBrandMismatch("evil.example", title); !r.IsMismatch {
			t.Errorf("%q not matched to a brand", title)
		}
	}
	// Non-Latin keywords still match their own script.
	if r := CheckBrandMismatch("evil.example", "日本郵便 再配達"); !r.IsMismatch || r.BrandFound != "Japan Post" {
		t.Errorf("Japanese title = %+v", r)
	}
}

func TestBrandMismatch_OwnAddressInTitle(t *testing.T) {
	if r := CheckBrandMismatch("octocat.github.io", "Octocat.github.io"); r.IsMismatch {
		t.Errorf("site titled with its own address flagged as %s", r.BrandFound)
	}
	if r := CheckBrandNames("octocat.github.io", "Octocat.github.io"); r.IsMismatch {
		t.Errorf("loose check flagged own address as %s", r.BrandFound)
	}
	// The page can still claim the brand in the rest of its title.
	if r := CheckBrandMismatch("evil.github.io", "Sign in to GitHub · evil.github.io"); !r.IsMismatch {
		t.Error("GitHub login page on github.io not flagged")
	}
}

func TestBrandMismatch_FancyLetters(t *testing.T) {
	// Mathematical bold letters have no lowercase of their own.
	for title, brand := range map[string]string{
		"𝗠𝗲𝘁å𝗺å𝘀𝗸 𝗟𝗼𝗴𝗶𝗻":                         "MetaMask",
		"G𝓮miñi : Login | Sign In":               "Gemini",
		"SEUR: Envío y transporte de paquetería": "SEUR",
	} {
		if r := CheckBrandMismatch("evil.gitbook.io", title); !r.IsMismatch || r.BrandFound != brand {
			t.Errorf("%q = %+v, want %s", title, r, brand)
		}
	}
}

func TestBrandMismatch_MemberBanks(t *testing.T) {
	for host, title := range map[string]string{
		"www.sparkasse-hannover.de": "Sparkasse Hannover - Die Bank in Hannover und der Region",
		"www.ksk-koeln.de":          "Kreissparkasse Köln",
		"www.vrbank-sw.de":          "VR Bank Südwestpfalz – Volksbank",
	} {
		if r := CheckBrandMismatch(host, title); r.IsMismatch {
			t.Errorf("%s flagged as impersonating %s", host, r.BrandFound)
		}
	}
	if r := CheckBrandMismatch("evil.example", "Sparkasse Online-Banking"); !r.IsMismatch {
		t.Error("Sparkasse title on an unrelated domain not flagged")
	}
}

func TestIPHostname(t *testing.T) {
	tests := map[string][2]string{
		"vxv-securestanadradbankcoza-z7ma1m-10a930-139-59-0-44.sslip.io": {"sslip.io", "139.59.0.44"},
		"webmail.50-6-22-93.nip.io":                                      {"nip.io", "50.6.22.93"},
		"1.2.3.4.nip.io":                                                 {"nip.io", "1.2.3.4"},
		"myapp.traefik.me":                                               {"traefik.me", ""},
	}
	for host, want := range tests {
		if s, ip, ok := IPHostname(host); !ok || s != want[0] || ip != want[1] {
			t.Errorf("IPHostname(%q) = %q %q %v, want %v", host, s, ip, ok, want)
		}
	}
	// Dates in names aren't addresses, and other hosts aren't these services.
	if _, _, ok := IPHostname("porgu-8y3-87x4-6osm4-29-09-2026-hh.pages.dev"); ok {
		t.Error("pages.dev host taken for an IP hostname service")
	}
}

func TestBrandMismatch_ExactTitle(t *testing.T) {
	if r := CheckBrandMismatch("www.365securitydevices.com", "Sign in to your account"); !r.IsMismatch || r.BrandFound != "Microsoft" {
		t.Errorf("Microsoft's login title elsewhere = %+v", r)
	}
	if r := CheckBrandMismatch("login.microsoftonline.com", "Sign in to your account"); r.IsMismatch {
		t.Error("Microsoft's own login flagged")
	}
	if r := CheckBrandMismatch("app.acme.example", "Sign in to your account | Acme"); r.IsMismatch {
		t.Errorf("longer title flagged as %s", r.BrandFound)
	}
}
