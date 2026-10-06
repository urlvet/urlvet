package checks

import (
	"errors"
	"github.com/urlvet/urlvet/internal/constants"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// maxShortLinkHops bounds how many short links in a row are followed.
const maxShortLinkHops = 10

// ErrShortLinkUnresolved means a short link's destination couldn't be found:
// it needs a click or a captcha, the service refused us, or it's dead.
var ErrShortLinkUnresolved = errors.New("short link destination not found")

// IsShortLink reports whether a URL is a link on a URL shortener, rather than
// the shortener's own homepage.
func IsShortLink(rawURL string) bool {
	u, err := url.Parse(rawURL)
	if err != nil {
		return false
	}
	if strings.Trim(u.Path, "/") == "" && u.RawQuery == "" {
		return false
	}
	// The host or any parent: click trackers sit on per-sender subdomains
	// (us14.list-manage.com).
	for host := strings.ToLower(u.Hostname()); host != ""; {
		for _, prefix := range constants.ShortLinkPaths[host] {
			if strings.HasPrefix(u.Path, prefix) {
				return true
			}
		}
		i := strings.Index(host, ".")
		if i < 0 {
			break
		}
		host = host[i+1:]
	}
	domain, err := GetDomain(rawURL)
	if err != nil {
		return false
	}
	return IsUrlShortener(domain) || IsUrlShortener(strings.ToLower(u.Hostname()))
}

// ResolveShortLink follows a short link to where it leads, through any
// further short links, and returns every URL on the way (starting with the
// link itself) and the destination. A destination that's still a short link
// returns ErrShortLinkUnresolved along with the chain so far.
func ResolveShortLink(rawURL string) (chain []string, target string, err error) {
	chain = []string{rawURL}
	seen := map[string]bool{rawURL: true}
	cur := rawURL
	for i := 0; i < maxShortLinkHops && IsShortLink(cur); i++ {
		hops, next := resolveShortLinkOnce(cur)
		if next == "" || next == cur || seen[next] {
			// Nowhere to go, or back to a link already visited (a loop).
			return chain, cur, ErrShortLinkUnresolved
		}
		for _, h := range append(hops, next) {
			if !seen[h] {
				seen[h] = true
				chain = append(chain, h)
			}
		}
		cur = next
	}
	if IsShortLink(cur) {
		return chain, cur, ErrShortLinkUnresolved
	}
	return chain, cur, nil
}

// resolveShortLinkOnce opens one short link and returns the HTTP redirects it
// went through and where it leads: the end of those redirects, or failing
// that a destination named in the page (interstitial "Redirecting…" pages).
func resolveShortLinkOnce(rawURL string) (hops []string, next string) {
	client := &http.Client{
		Timeout:   8 * time.Second,
		Transport: newSafeTransport(),
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			hops = append(hops, req.URL.String())
			// Stop at the first non-shortener: that's the destination, and
			// it's for the full scan to open, not us.
			if !IsShortLink(req.URL.String()) || len(via) >= maxShortLinkHops {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return nil, ""
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,*/*;q=0.8")
	resp, err := client.Do(req)
	if err != nil {
		if len(hops) > 0 {
			return hops, hops[len(hops)-1]
		}
		return nil, ""
	}
	defer resp.Body.Close()

	if len(hops) > 0 {
		last := hops[len(hops)-1]
		if !IsShortLink(last) || resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return hops, last
		}
	}

	body, _ := io.ReadAll(io.LimitReader(resp.Body, 512*1024))
	base := resp.Request.URL
	if t := pageRedirectTarget(string(body), base); t != "" {
		return hops, t
	}
	if len(hops) > 0 {
		return hops, hops[len(hops)-1]
	}
	return hops, ""
}

// redirectAttrs are attributes interstitial pages keep the destination in
// (goo.su: <div data-url="https://…">).
var redirectAttrs = []string{"data-url", "data-href", "data-redirect", "data-redirect-url", "data-target-url", "data-link", "data-destination"}

// pageRedirectTarget finds the destination an interstitial page names: a meta
// refresh, a literal URL assigned to location, or a data-url-style attribute.
// Only absolute http(s) URLs on another host count.
func pageRedirectTarget(body string, base *url.URL) string {
	resolve := func(raw string) string {
		raw = strings.TrimSpace(html.UnescapeString(raw))
		u, err := url.Parse(raw)
		if err != nil {
			return ""
		}
		u = base.ResolveReference(u)
		if (u.Scheme != "http" && u.Scheme != "https") || strings.EqualFold(u.Hostname(), base.Hostname()) {
			return ""
		}
		return u.String()
	}

	doc, err := html.Parse(strings.NewReader(body))
	if err != nil {
		return ""
	}
	var found string
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if found != "" {
			return
		}
		if n.Type == html.ElementNode {
			switch {
			case n.Data == "meta" && strings.EqualFold(getAttr(n, "http-equiv"), "refresh"):
				if m := metaRefreshURLRe.FindStringSubmatch(getAttr(n, "content")); m != nil {
					found = resolve(m[1])
				}
			case n.Data == "script" && getAttr(n, "src") == "" && n.FirstChild != nil:
				for _, m := range jsRedirectTargetRe.FindAllStringSubmatch(n.FirstChild.Data, -1) {
					if found = resolve(m[1] + m[2]); found != "" {
						break
					}
				}
			default:
				for _, a := range redirectAttrs {
					if v := getAttr(n, a); v != "" {
						if found = resolve(v); found != "" {
							break
						}
					}
				}
			}
		}
		for c := n.FirstChild; c != nil && found == ""; c = c.NextSibling {
			walk(c)
		}
	}
	walk(doc)
	return found
}
