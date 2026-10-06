package checks

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func TestDetectScriptRedirect(t *testing.T) {
	long := strings.Repeat("words ", 100)
	tests := []struct {
		name    string
		scripts []string
		meta    string
		text    int
		wantIP  bool
		wantAny bool
	}{
		{
			name: "redirector built from a raw IP",
			scripts: []string{`var p = window.location.href.split('#')[1];
var srv_ip = "185.80.128.4";
document.location.href = 'http://'+srv_ip+'/?'+p;`},
			wantIP: true, wantAny: true,
		},
		{
			name:    "thin page redirecting to another domain",
			scripts: []string{`window.location.replace("https://evil.example/login")`},
			wantAny: true,
		},
		{
			name:    "button handler on a full page",
			scripts: []string{`function go(){ location.href = "https://partner.example/"; }`},
			text:    len(long),
		},
		{
			name:    "same-site redirect",
			scripts: []string{`location.href = "https://www.site.example/home"`},
		},
		{
			name:    "IP in a script that never redirects",
			scripts: []string{`var ip = "8.8.8.8";`},
		},
		{
			name:    "private IP is not a raw-IP redirect",
			scripts: []string{`location.href = "http://192.168.1.1/"`},
			wantAny: true,
		},
		{
			name:    "meta refresh to another domain",
			meta:    "0; url=https://other.example/",
			text:    len(long),
			wantAny: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := detectScriptRedirect(tt.scripts, tt.meta, "site.example", tt.text)
			if (got != nil) != tt.wantAny {
				t.Fatalf("got %+v, want redirect=%v", got, tt.wantAny)
			}
			if got != nil && got.ToIP != tt.wantIP {
				t.Errorf("ToIP = %v, want %v", got.ToIP, tt.wantIP)
			}
		})
	}
}

func TestExecutableExt(t *testing.T) {
	tests := map[string]string{
		"https://github.com/a/b/releases/download/v1/VMware-Workstation.exe": ".exe",
		"VMware-Workstation-Full-26H1u1-25688693.exe":                        ".exe",
		"https://example.com/app.APK?ref=x":                                  ".apk",
		"https://example.com/docs/guide.pdf":                                 "",
		"https://cdn.example.com/app.js":                                     "",
		"https://example.com/go/example.com":                                 "",
		"":                                                                   "",
	}
	for in, want := range tests {
		if got := ExecutableExt(in); got != want {
			t.Errorf("ExecutableExt(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestDownloadFileName(t *testing.T) {
	u, _ := url.Parse("https://release-assets.example/asset/123?sig=x")
	resp := &http.Response{
		Header:  http.Header{"Content-Disposition": {"attachment; filename=VMware.exe"}},
		Request: &http.Request{URL: u},
	}
	if got := downloadFileName(resp); got != "VMware.exe" {
		t.Errorf("downloadFileName = %q, want VMware.exe", got)
	}
	resp.Header = http.Header{}
	if got := downloadFileName(resp); got != "123" {
		t.Errorf("downloadFileName without header = %q, want 123", got)
	}
}

func TestIsHTMLResponse(t *testing.T) {
	tests := []struct {
		ct   string
		body string
		want bool
	}{
		{"text/html; charset=utf-8", "", true},
		{"application/xhtml+xml", "", true},
		{"application/octet-stream", "", false},
		{"image/png", "", false},
		{"", "<!doctype html><html><body>hi</body></html>", true},
		{"", "MZ\x90\x00\x03\x00\x00\x00", false},
	}
	for _, tt := range tests {
		resp := &http.Response{Header: http.Header{}}
		if tt.ct != "" {
			resp.Header.Set("Content-Type", tt.ct)
		}
		if got := isHTMLResponse(resp, []byte(tt.body)); got != tt.want {
			t.Errorf("isHTMLResponse(%q, %q) = %v, want %v", tt.ct, tt.body, got, tt.want)
		}
	}
}

func formFrom(t *testing.T, src string) FormInfo {
	t.Helper()
	doc, err := html.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatal(err)
	}
	var form *html.Node
	var find func(*html.Node)
	find = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "form" && form == nil {
			form = n
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			find(c)
		}
	}
	find(doc)
	base, _ := url.Parse("https://site.example/")
	return extractFormInfo(form, base, "site.example")
}

func TestLoginIntent(t *testing.T) {
	tests := []struct {
		name string
		html string
		want bool // email field plus login intent
	}{
		{"newsletter / sample request", `<form><input type="email" placeholder="Work email"><button>Send the sample</button></form>`, false},
		{"subscribe", `<form><input type="email" name="email"><input type="submit" value="Subscribe"></form>`, false},
		{"email-first sign in", `<form><input type="email" name="loginfmt"><input type="submit" value="Next"></form>`, true},
		{"sign in button", `<form><input type="text" name="username"><button>Sign in</button></form>`, true},
		{"login action", `<form action="/auth/start"><input type="email"><button>Go</button></form>`, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := formFrom(t, tt.html)
			if got := f.ContainsUserLike && f.LoginIntent; got != tt.want {
				t.Errorf("userLike=%v loginIntent=%v, want login=%v", f.ContainsUserLike, f.LoginIntent, tt.want)
			}
		})
	}
}

func TestTooDeepAndLongUrl(t *testing.T) {
	if TooDeepUrl("https://github.com/201853910/VMwareWorkstation/releases/download/26H1/x.exe") {
		t.Error("six path segments flagged as too deep")
	}
	if !TooDeepUrl("https://e.example/a/b/c/d/e/f/g") {
		t.Error("seven path segments not flagged")
	}
	if TooDeepUrl("https://e.example/a?next=https://x/y/z/w/v/u/t") {
		t.Error("slashes in the query counted as depth")
	}
	if TooLongUrl("https://storage.googleapis.com/usales26/usales26.html%23?Z289MSZzMT0yMzc4MTUzJnMyPTgwODY0Nzk5OSZzMz1HTEI=") {
		t.Error("120-character URL flagged as too long")
	}
	if !TooLongUrl("https://e.example/" + strings.Repeat("a", 150)) {
		t.Error("170-character URL not flagged")
	}
}

func TestAsksPersonal(t *testing.T) {
	tests := map[string]bool{
		"emailAddress":            false, // github.blog's newsletter field
		"newsletter_emailAddress": false,
		"Your email address":      false,
		"e-mail address":          false,
		"ipAddress":               false,
		"statement":               false,
		"username":                false,
		"street_address":          true,
		"billingAddress":          true,
		"Home address":            true,
		"phoneNumber":             true,
		"dateOfBirth":             true,
		"ZIP":                     true,
		"postcode":                true,
		"state":                   true,
		"ssn":                     true,
	}
	for in, want := range tests {
		if got := asksPersonal(in); got != want {
			t.Errorf("asksPersonal(%q) = %v, want %v", in, got, want)
		}
	}
}

func TestIsMarketingFormService(t *testing.T) {
	tests := map[string]bool{
		"https://s88570519.t.eloqua.com/e/f2?elqFormName=newsletter": true,
		"https://github.us21.list-manage.com/subscribe/post":         true,
		"https://forms.hsforms.com/submissions/v3":                   true,
		"https://formspree.io/f/abc":                                 false, // phishing kits use these
		"https://getform.io/f/abc":                                   false,
		"https://evil.example/collect.php":                           false,
	}
	for in, want := range tests {
		if got := IsMarketingFormService(in); got != want {
			t.Errorf("IsMarketingFormService(%q) = %v, want %v", in, got, want)
		}
	}
}

func TestNewsletterFormIsNotPersonal(t *testing.T) {
	f := formFrom(t, `<form action="https://s88570519.t.eloqua.com/e/f2"><input type="email" name="emailAddress" placeholder="Your email address"><input type="checkbox" name="marketingEmailOptIn1"><button>Subscribe</button></form>`)
	if f.ContainsPersonal {
		t.Error("newsletter email field counted as personal details")
	}
}
