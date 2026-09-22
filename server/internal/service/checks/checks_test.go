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
