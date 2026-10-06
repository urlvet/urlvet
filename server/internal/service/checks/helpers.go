package checks

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/urlvet/urlvet/internal/constants"
	"golang.org/x/net/publicsuffix"
)

func GetTld(domain string) (string, bool) {
	tld, icann := publicsuffix.PublicSuffix(domain)
	return tld, icann
}

func GetDomain(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	host := strings.ToLower(parsedURL.Hostname())
	if site := customerSite(host); site != "" {
		return site, nil
	}
	domain, err := publicsuffix.EffectiveTLDPlusOne(host)
	if err != nil {
		// A host that is itself a public suffix (gov.uk, github.io) has no
		// registrable domain above it; the host is the closest thing to one.
		if ps, _ := publicsuffix.PublicSuffix(host); host != "" && ps == host {
			return host, nil
		}
		return "", err
	}
	return domain, nil
}

// CustomerSitePlatform returns the provider when host is a customer's site on
// a builder or host that isn't on the Public Suffix List (someone.weebly.com
// → weebly.com), or "" otherwise.
func CustomerSitePlatform(host string) string {
	host = strings.ToLower(host)
	for suffix := range constants.CustomerSiteHosts {
		if !strings.HasSuffix(host, "."+suffix) {
			continue
		}
		rest := strings.TrimSuffix(host, "."+suffix)
		if _, own := constants.ProviderSubdomains[suffix][rest[strings.LastIndex(rest, ".")+1:]]; own {
			return ""
		}
		return suffix
	}
	return ""
}

// customerSite returns the customer's site for a host on a CustomerSiteHosts
// provider: the whole host without "www." (www.someone.weebly.com →
// someone.weebly.com). Providers name customers differently (GoDaddy puts an
// ID first and a region second: 1286524.us23.myftpupload.com), so no single
// label can stand for the customer. The provider's own www host stays the
// provider. "" when host isn't on such a provider.
func customerSite(host string) string {
	suffix := CustomerSitePlatform(host)
	if suffix == "" {
		return ""
	}
	rest := strings.TrimPrefix(strings.TrimSuffix(host, "."+suffix), "www.")
	if rest == "www" || rest == "" {
		return suffix
	}
	return rest + "." + suffix
}

func GetHost(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	host := parsedURL.Hostname()
	return host, nil
}

func IsValidURL(rawURL string) (*url.URL, bool, error) {
	if !strings.Contains(rawURL, "://") {
		rawURL = "https://" + rawURL
	}

	// url.Parse, not ParseRequestURI: a request URI has no fragment, so
	// ParseRequestURI folds "#…" into the path as "%23…" and the page 404s.
	// Redirector links carry their payload in the fragment.
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return nil, false, err
	}
	if parsed.Scheme == "" || parsed.Host == "" {
		return nil, false, nil
	}
	// Hostnames are case-insensitive; normalize so "Google.com" matches rank,
	// typosquat and cache lookups for "google.com". Path and query keep their case.
	parsed.Host = strings.ToLower(parsed.Host)

	host := parsed.Hostname()
	// Allow bare IPs (IPv4 or IPv6)
	if net.ParseIP(host) != nil {
		return parsed, true, nil
	}
	// Require at least one dot and a TLD of ≥2 chars (rejects "financialcomm", "localhost", etc.)
	dot := strings.LastIndex(host, ".")
	if dot < 0 || len(host)-dot-1 < 2 {
		return nil, false, nil
	}
	return parsed, true, nil
}

func GetDomainAge(created time.Time) (string, int, error) {
	now := time.Now()
	if created.After(now) {
		return "not yet registered", 0, nil
	}

	years := now.Year() - created.Year()
	months := int(now.Month()) - int(created.Month())
	days := int(now.Sub(created).Hours() / 24)

	if months < 0 {
		years--
		months += 12
	}

	if years <= 0 && months <= 0 {
		switch {
		case days == 0:
			return "today", days, nil
		case days == 1:
			return "1 day", days, nil
		case days < 30:
			return fmt.Sprintf("%d days", days), days, nil
		default:
			return "less than a month", days, nil
		}
	}

	parts := []string{}
	if years > 0 {
		if years == 1 {
			parts = append(parts, "1 year")
		} else {
			parts = append(parts, fmt.Sprintf("%d years", years))
		}
	}
	if months > 0 {
		if months == 1 {
			parts = append(parts, "1 month")
		} else {
			parts = append(parts, fmt.Sprintf("%d months", months))
		}
	}
	return strings.Join(parts, " "), days, nil
}

// --- SSRF-safe HTTP utilities ---

// privateIPNets lists IP ranges that must never be dialed from the server.
var privateIPNets = func() []*net.IPNet {
	cidrs := []string{
		"0.0.0.0/8",      // "this" network
		"10.0.0.0/8",     // RFC 1918 private
		"100.64.0.0/10",  // shared address space (RFC 6598 / CGN)
		"127.0.0.0/8",    // loopback
		"169.254.0.0/16", // link-local — includes AWS/GCP metadata endpoints
		"172.16.0.0/12",  // RFC 1918 private
		"192.0.0.0/24",   // IETF protocol assignments
		"192.168.0.0/16", // RFC 1918 private
		"198.18.0.0/15",  // benchmark testing
		"224.0.0.0/4",    // multicast
		"240.0.0.0/4",    // reserved
		"::1/128",        // IPv6 loopback
		"fc00::/7",       // IPv6 unique local
		"fe80::/10",      // IPv6 link-local
	}
	nets := make([]*net.IPNet, 0, len(cidrs))
	for _, c := range cidrs {
		_, ipnet, _ := net.ParseCIDR(c)
		nets = append(nets, ipnet)
	}
	return nets
}()

func isPrivateIP(ip net.IP) bool {
	for _, block := range privateIPNets {
		if block.Contains(ip) {
			return true
		}
	}
	return false
}

// ssrfSafeDialContext wraps a Dialer's DialContext to reject connections to
// private/reserved IP ranges, preventing SSRF attacks.
// DNS is resolved once here and the first public IP is dialed directly,
// which also prevents DNS rebinding attacks.
func ssrfSafeDialContext(d *net.Dialer) func(context.Context, string, string) (net.Conn, error) {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		host, port, err := net.SplitHostPort(addr)
		if err != nil {
			return nil, err
		}
		ips, err := net.DefaultResolver.LookupHost(ctx, host)
		if err != nil {
			return nil, err
		}
		if len(ips) == 0 {
			return nil, fmt.Errorf("ssrf: no addresses resolved for %s", host)
		}
		for _, raw := range ips {
			ip := net.ParseIP(raw)
			if ip != nil && isPrivateIP(ip) {
				return nil, fmt.Errorf("ssrf: blocked connection to private/reserved address %s", raw)
			}
		}
		// Dial the first resolved IP directly to prevent re-resolution.
		return d.DialContext(ctx, network, net.JoinHostPort(ips[0], port))
	}
}

// newSafeHTTPClient returns an *http.Client that blocks SSRF.
// Use for plain fetches that don't need a custom CheckRedirect or transport.
func newSafeHTTPClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout:   timeout,
		Transport: newSafeTransport(),
	}
}

// newSafeTransport returns an *http.Transport with SSRF-safe DialContext.
// Use when building a custom http.Client (e.g. one with a CheckRedirect hook).
func newSafeTransport() *http.Transport {
	return &http.Transport{
		DialContext:           ssrfSafeDialContext(&net.Dialer{Timeout: 5 * time.Second}),
		ResponseHeaderTimeout: 5 * time.Second,
		TLSHandshakeTimeout:   5 * time.Second,
	}
}

// sameSite reports whether two registrable domains (from GetDomain) belong to
// the same site. They usually match exactly; the exception is a host that is
// itself a public suffix, like gov.uk, whose pages live on www.gov.uk.
func sameSite(a, b string) bool {
	return a == b || strings.HasSuffix(a, "."+b) || strings.HasSuffix(b, "."+a)
}
