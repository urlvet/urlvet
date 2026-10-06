package analyzer

import (
	"fmt"
	"net"
	"strings"

	"github.com/urlvet/urlvet/internal/service/checks"
	"github.com/urlvet/urlvet/internal/service/threatfeeds"
)

// When a link leads to another site, the scan judges the page it lands on:
// trust comes only from there, since the site that sent the visitor on can't
// vouch for where it sent them. What the link itself gets wrong still counts,
// though, and so does anything bad it passes through on the way. These are
// the findings carried over from the link as given (the origin) and from the
// hops between it and the destination: risk only, never trust.

// OriginFindings is what the submitted link and the hops after it add to the
// destination's result.
type OriginFindings struct {
	URL        string   `json:"url"`
	BadReasons []string `json:"bad_reasons,omitempty"`
	Risk       int      `json:"risk"`
	// Confirmed is set when a list names the link or a hop: the result is
	// Risky, whatever the destination looks like.
	Confirmed bool `json:"confirmed,omitempty"`
}

func (o *OriginFindings) add(reason string, risk int) {
	o.BadReasons = append(o.BadReasons, reason)
	o.Risk += risk
}

// redirectTarget returns where the scanned page sends visitors on to another
// site, by HTTP redirect or straight away by script or meta refresh, or "".
// Moves within one site or brand were already cleared by applyExceptions.
func redirectTarget(resp Response) string {
	if rr := resp.Analysis.RedirectionResult; rr.HasDomainJump && rr.FinalURL != "" {
		return rr.FinalURL
	}
	if c := resp.ContentData; c != nil && c.ScriptRedirect != nil && c.ScriptRedirect.CrossDomain && c.ScriptRedirect.Target != "" {
		return c.ScriptRedirect.Target
	}
	return ""
}

// worthFollowing reports whether a redirect's destination should be judged
// in the link's place. A ranked site or one on a restricted registry sending
// visitors on, or anything leading to a site nobody knows, is judged where it
// lands. An unknown site landing on a well-known one is the shape of
// cloaking: kits send scanners to google.com or the brand's real login and
// victims to the phishing page, so the destination's good name can't be
// borrowed; the link is judged itself, and the jump counts against it.
func worthFollowing(src Response, dest string) bool {
	return src.Features.Rank > 0 || src.Features.TLD.IsTrusted || !isVouchedDestination(dest)
}

// originFindings collects the risks of the submitted link itself: how it's
// written, what lists say about it, and how its own domain looks, but not
// anything that would vouch for it.
func originFindings(src Response) *OriginFindings {
	o := &OriginFindings{URL: src.URL}
	host := hostOf(src.URL)
	link := fmt.Sprintf("The link you were given (%s)", host)

	// How it's written.
	if src.Features.URL.UsesIP {
		o.add(link+" is a bare IP address.", 100)
	}
	if service, ip, ok := checks.IPHostname(host); ok {
		what := "a bare server"
		if ip != "" {
			what = fmt.Sprintf("a bare server (%s)", ip)
		}
		o.add(fmt.Sprintf("%s is %s written as a name on %s.", link, what, service), 100)
	}
	if src.Features.URL.ContainsPunycode {
		o.add(link+" uses punycode, which can make one address look like another.", 100)
	}
	if src.Features.URL.HasHomoglyph {
		o.add(link+" uses lookalike letters.", 60)
	}
	if src.Features.URL.Keywords.HasKeywords {
		o.add(fmt.Sprintf("%s contains words like %s.", link, strings.Join(src.Features.URL.Keywords.Found, ", ")), 10)
	}
	if victimAddressInLink(src.URL) {
		o.add(link+" carries an email address for a page to fill in, a common way phishing tailors a login to each target.", 20)
	}

	// Brands it borrows.
	ts := src.TyposquatResult
	switch {
	case ts.IsSuspicious && ts.IsComboSquat:
		risk := 20
		if src.Features.TLD.IsHostingPlatform {
			risk = 40
		}
		o.add(fmt.Sprintf("%s uses the brand name '%s' but isn't theirs.", link, ts.MatchedBrand), risk)
	case ts.IsSuspicious:
		o.add(fmt.Sprintf("%s closely resembles %s.", link, ts.MatchedDomain), 40)
	}
	switch {
	case ts.EmbeddedDomain != "":
		o.add(fmt.Sprintf("%s starts like %s but isn't.", link, ts.EmbeddedDomain), 60)
	case ts.SubdomainBrand != "" && ts.SubdomainBrandExact:
		o.add(fmt.Sprintf("%s has a subdomain named '%s', on a site that isn't theirs.", link, ts.SubdomainBrand), 20)
	case ts.SubdomainBrand != "":
		o.add(fmt.Sprintf("%s has a subdomain using the brand name '%s', on a site that isn't theirs.", link, ts.SubdomainBrand), 40)
	}
	if ts.PathBrand != "" {
		if src.Features.TLD.IsTrusted {
			o.add(fmt.Sprintf("%s is a '%s' sign-in path under a government or university address, which suggests the site was hacked.", link, ts.PathBrand), 60)
		} else {
			o.add(fmt.Sprintf("%s names '%s' in its path, on a site that isn't theirs.", link, ts.PathBrand), 20)
		}
	}

	// What lists say about it.
	if fm := src.ThreatFeeds; fm != nil && fm.Listed {
		sources := strings.Join(fm.Sources, " and ")
		if fm.Match == "url" {
			o.add(fmt.Sprintf("CONFIRMED THREAT: %s is listed as phishing or malware by %s.", link, sources), 200)
			o.Confirmed = true
		} else {
			o.add(fmt.Sprintf("Other pages on %s are listed as phishing or malware by %s.", host, sources), 80)
		}
	}
	for _, g := range []*threatfeeds.GoogleThreatResult{src.WebRisk, src.SafeBrowsing} {
		if g != nil && g.Listed {
			o.add(fmt.Sprintf("CONFIRMED THREAT: Google lists %s as %s (advisory provided by Google).", host, webRiskThreatNames(g.ThreatTypes)), 200)
			o.Confirmed = true
			break
		}
	}
	if p := src.Phishing; p != nil && p.InDatabase && p.Valid && p.Verified {
		o.add(fmt.Sprintf("CONFIRMED PHISHING: %s is a verified phishing link.", link), 200)
		o.Confirmed = true
	}

	// How its own domain looks.
	if age, known := src.DomainInfo.Age(); known && age <= newDomainDays &&
		!(src.Features.TLD.IsHostingPlatform && hostingRegistration(src)) {
		o.add(fmt.Sprintf("%s is on a domain registered %s.", link, ago(src.DomainInfo.AgeHuman)), 25)
	}
	if src.Features.TLD.IsRisky {
		o.add(link+" uses a high-risk domain ending.", 20)
	}
	if isWordPressSystemPage(src.URL) {
		o.add(link+" points into a WordPress system folder, where phishing kits hide on hacked sites.", 30)
	}
	if store := userContentHost(host); store != "" && src.ContentData != nil && src.ContentData.ScriptRedirect != nil {
		o.add(fmt.Sprintf("%s is a file anyone could upload to %s, and it only forwards visitors.", link, store), 60)
	}
	return o
}

// addHopFindings adds what the hops between the link and its destination
// pass through. Only local checks: the hops aren't fetched or looked up.
func addHopFindings(o *OriginFindings, chain []string, from, to string) {
	if len(chain) < 3 {
		return
	}
	fromDomain, _ := checks.GetDomain(from)
	toDomain, _ := checks.GetDomain(to)
	seen := map[string]bool{}
	for _, hop := range chain[1 : len(chain)-1] {
		host := hostOf(hop)
		domain, _ := checks.GetDomain(hop)
		if host == "" || seen[host] || sameSite(domain, fromDomain) || sameSite(domain, toDomain) {
			continue
		}
		seen[host] = true
		switch {
		case net.ParseIP(strings.Trim(host, "[]")) != nil:
			o.add(fmt.Sprintf("On the way, it passes through a bare IP address (%s).", host), 60)
		default:
			if service, ip, ok := checks.IPHostname(host); ok {
				o.add(fmt.Sprintf("On the way, it passes through a bare server (%s) written as a name on %s.", ip, service), 60)
			}
		}
		if risky, _, _ := checks.IsRiskyTld(domain); risky {
			o.add(fmt.Sprintf("On the way, it passes through %s, on a high-risk domain ending.", host), 10)
		}
		if fm := threatfeeds.LookupLocal(hop, false); fm.Listed {
			sources := strings.Join(fm.Sources, " and ")
			if fm.Match == "url" {
				o.add(fmt.Sprintf("CONFIRMED THREAT: on the way, it passes through %s, listed as phishing or malware by %s.", host, sources), 200)
				o.Confirmed = true
			} else {
				o.add(fmt.Sprintf("On the way, it passes through %s, where other pages are listed as phishing or malware by %s.", host, sources), 80)
			}
		}
	}
}

// sameSite reports whether two registrable domains are the same site.
func sameSite(a, b string) bool {
	return a != "" && b != "" && (a == b || strings.HasSuffix(a, "."+b) || strings.HasSuffix(b, "."+a))
}
