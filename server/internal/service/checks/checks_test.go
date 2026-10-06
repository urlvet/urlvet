package checks

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func TestTooLongUrl(t *testing.T) {
	url := "https://google.com"
	result := TooLongUrl(url)
	if result == true {
		t.Errorf("Error while testing TooLongUrl(%v)", url)
	}
}

func TestIsValidURL_LowercasesHost(t *testing.T) {
	tests := map[string]string{
		"http://Google.com/aa":      "http://google.com/aa",
		"Google.com/Hh":             "https://google.com/Hh",
		"HTTPS://WWW.Example.COM":   "https://www.example.com",
		"https://Example.com/A?Q=B": "https://example.com/A?Q=B",
		"https://e.com/p.html#?x=1": "https://e.com/p.html#?x=1",
	}
	for in, want := range tests {
		u, ok, err := IsValidURL(in)
		if err != nil || !ok {
			t.Fatalf("IsValidURL(%q) = invalid, err=%v", in, err)
		}
		if got := u.String(); got != want {
			t.Errorf("IsValidURL(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestHasAncestor_SVGTitle(t *testing.T) {
	doc, err := html.Parse(strings.NewReader(`<html><head><title>GitHub</title></head><body><svg><title>Vodafone</title></svg></body></html>`))
	if err != nil {
		t.Fatal(err)
	}
	var titles []bool
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" {
			titles = append(titles, hasAncestor(n, "svg"))
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(doc)
	if len(titles) != 2 || titles[0] || !titles[1] {
		t.Errorf("hasAncestor(title, svg) = %v, want [false true]", titles)
	}
}

func TestGetDomain_PublicSuffixHost(t *testing.T) {
	tests := map[string]string{
		"https://gov.uk/":           "gov.uk",
		"https://www.gov.uk/":       "www.gov.uk",
		"https://www.hdfc.bank.in/": "hdfc.bank.in",
		"https://Example.COM/a":     "example.com",
	}
	for in, want := range tests {
		got, err := GetDomain(in)
		if err != nil || got != want {
			t.Errorf("GetDomain(%q) = %q, %v; want %q", in, got, err, want)
		}
	}
}

func TestCheckBrandMismatch_WholeWords(t *testing.T) {
	// HDFC's title mentions "NetBanking"; Commonwealth Bank's keyword is "netbank".
	if r := CheckBrandMismatch("hdfcbank.com", "HDFC Bank: Personal Banking Services | NetBanking"); r.IsMismatch {
		t.Errorf("NetBanking matched as %q", r.BrandFound)
	}
	if r := CheckBrandMismatch("commbank-verify.xyz", "CommBank NetBank - Log on"); !r.IsMismatch || r.BrandFound != "Commonwealth Bank" {
		t.Errorf("lookalike CommBank page = %+v, want a Commonwealth Bank mismatch", r)
	}
}

func TestCheckBrandMismatch_OfficialDomain(t *testing.T) {
	r := CheckBrandMismatch("paypal-login-verify.top", "PayPal Login")
	if r.OfficialDomain != "paypal.com" {
		t.Errorf("OfficialDomain = %q, want paypal.com", r.OfficialDomain)
	}
	if r := CheckBrandMismatch("paypal.com", "PayPal Login"); r.OfficialDomain != "" {
		t.Errorf("real site got OfficialDomain %q", r.OfficialDomain)
	}
}

func TestSameSite(t *testing.T) {
	tests := []struct {
		a, b string
		want bool
	}{
		{"paypal.com", "paypal.com", true},
		{"gov.uk", "www.gov.uk", true},
		{"zoom.us", "zoom.com", false},
		{"paypal.com", "evilpaypal.com", false},
	}
	for _, tt := range tests {
		if got := sameSite(tt.a, tt.b); got != tt.want {
			t.Errorf("sameSite(%q, %q) = %v, want %v", tt.a, tt.b, got, tt.want)
		}
	}
}
