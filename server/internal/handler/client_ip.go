package handler

import (
	"net/netip"
	"os"
	"strings"
	"sync/atomic"

	"github.com/gin-gonic/gin"
	"github.com/urlvet/urlvet/internal/logger"
)

// Client IP resolution
//
// The rate limiter keys on ctx.ClientIP(). Behind a reverse proxy the TCP peer
// is the proxy, so the real client IP has to come from a header the proxy
// sets, usually X-Forwarded-For. A client can send that header too, so it may
// only be believed when the request arrives from a proxy we trust. Gin's
// default trusts every peer, which lets anyone pick their own rate-limit key.
//
// Two settings control this (see docs/configuration.md):
//
//   - TRUSTED_PROXIES: peers whose X-Forwarded-For / X-Real-IP is believed.
//     Unset means loopback and private networks, which covers a proxy on the
//     same host or Docker network (the setup in docs/deployment.md). It is the
//     same default Rails, Echo and authentik use. "none" ignores the headers,
//     "all" restores Gin's old trust-everyone behaviour.
//   - TRUSTED_PLATFORM: a header set by a CDN or platform that connects to us
//     directly, read before anything else. Gin trusts it from any peer, so the
//     origin must only accept traffic from that platform.

// defaultTrustedProxies are the loopback and private ranges a local reverse
// proxy connects from. A proxy on the Docker host reaches the backend through
// the bridge gateway (172.17.0.0/16 and up), which 172.16.0.0/12 covers.
var defaultTrustedProxies = []string{
	"127.0.0.0/8", "::1/128", // loopback
	"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16", // RFC 1918
	"fc00::/7", // IPv6 unique local
}

// trustedPlatforms maps TRUSTED_PLATFORM shorthands to the header that platform
// sets to the client IP. Any other value is used as a header name as-is.
var trustedPlatforms = map[string]string{
	"cloudflare":        gin.PlatformCloudflare,
	"flyio":             gin.PlatformFlyIO,
	"google-app-engine": gin.PlatformGoogleAppEngine,
}

// configureClientIP applies TRUSTED_PROXIES and TRUSTED_PLATFORM to r. It
// returns a middleware that warns once about a likely misconfiguration, or nil
// when that check doesn't apply.
func configureClientIP(r *gin.Engine) gin.HandlerFunc {
	platform := strings.TrimSpace(os.Getenv("TRUSTED_PLATFORM"))
	if platform != "" {
		if header, ok := trustedPlatforms[strings.ToLower(platform)]; ok {
			platform = header
		}
		r.TrustedPlatform = platform
		logger.Info("client IP: reading platform header; only the platform must be able to reach the backend (firewall it to the platform's IP ranges)", "header", platform)
	}

	raw := strings.TrimSpace(os.Getenv("TRUSTED_PROXIES"))
	var proxies []string
	switch {
	case raw == "":
		proxies = defaultTrustedProxies
	case strings.EqualFold(raw, "none"):
		proxies = nil
	case strings.EqualFold(raw, "all"):
		proxies = []string{"0.0.0.0/0", "::/0"}
		logger.Warn("TRUSTED_PROXIES=all: any client can set its own IP with X-Forwarded-For and bypass the rate limit")
	default:
		for _, p := range strings.Split(raw, ",") {
			if p = strings.TrimSpace(p); p != "" {
				proxies = append(proxies, p)
			}
		}
	}

	trusted, err := parsePrefixes(proxies)
	if err == nil {
		err = r.SetTrustedProxies(proxies)
	}
	if err != nil {
		// Fail closed: a typo shouldn't reopen the bypass.
		logger.Error("invalid TRUSTED_PROXIES, ignoring forwarding headers", "err", err)
		_ = r.SetTrustedProxies(nil)
		return nil
	}
	if len(proxies) == 0 {
		logger.Info("client IP: ignoring forwarding headers, using the connecting address")
	} else {
		logger.Info("client IP: trusting forwarding headers", "from", strings.Join(proxies, ","))
	}

	// With a platform header, or nothing left to distrust, the warning has
	// nothing to catch.
	if platform != "" || len(proxies) == 0 || strings.EqualFold(raw, "all") {
		return nil
	}
	return warnUntrustedForwarder(trusted)
}

// parsePrefixes parses IPs and CIDRs the way SetTrustedProxies accepts them, so
// the warning middleware sees the same set Gin does.
func parsePrefixes(list []string) ([]netip.Prefix, error) {
	out := make([]netip.Prefix, 0, len(list))
	for _, s := range list {
		if strings.Contains(s, "/") {
			p, err := netip.ParsePrefix(s)
			if err != nil {
				return nil, err
			}
			out = append(out, p.Masked())
			continue
		}
		a, err := netip.ParseAddr(s)
		if err != nil {
			return nil, err
		}
		out = append(out, netip.PrefixFrom(a, a.BitLen()))
	}
	return out, nil
}

// warnUntrustedForwarder logs once when a request carries forwarding headers
// from a peer outside the trusted set. That is either a client trying to spoof
// its IP (harmless, the header is ignored) or a proxy with a public address,
// such as a CDN connecting straight to the backend. In the second case every
// user shares the proxy's rate-limit bucket until TRUSTED_PROXIES or
// TRUSTED_PLATFORM is set, so the log says how to fix it. The peer IP itself is
// not logged, in line with the request logger.
func warnUntrustedForwarder(trusted []netip.Prefix) gin.HandlerFunc {
	var warned atomic.Bool
	return func(c *gin.Context) {
		if !warned.Load() && (c.GetHeader("X-Forwarded-For") != "" || c.GetHeader("X-Real-IP") != "") {
			if ip, err := netip.ParseAddr(c.RemoteIP()); err == nil && !inPrefixes(trusted, ip.Unmap()) && warned.CompareAndSwap(false, true) {
				logger.Warn("got X-Forwarded-For/X-Real-IP from an untrusted address and ignored it. " +
					"If a CDN or proxy with a public IP sits in front of the backend, set TRUSTED_PLATFORM (e.g. cloudflare) or TRUSTED_PROXIES, " +
					"or all users share one rate limit; see docs/configuration.md")
			}
		}
		c.Next()
	}
}

func inPrefixes(prefixes []netip.Prefix, ip netip.Addr) bool {
	for _, p := range prefixes {
		if p.Contains(ip) {
			return true
		}
	}
	return false
}
