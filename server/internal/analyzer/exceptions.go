package analyzer

import (
	neturl "net/url"
	"regexp"
	"strings"

	"github.com/urlvet/urlvet/internal/constants"
	"github.com/urlvet/urlvet/internal/service/checks"
	"github.com/urlvet/urlvet/internal/service/rank"
	"github.com/urlvet/urlvet/internal/service/threatfeeds"
	"github.com/urlvet/urlvet/internal/service/typosquat"
)

// lookupRank is rank.DomainRankLookup, swappable in tests.
var lookupRank = rank.DomainRankLookup

// isTopRankedDomain reports whether the domain is in the top-1M popularity list.
// Popular domains are trusted enough that name-based heuristics (typosquatting,
// risky TLD, sensitive keywords) produce far more false positives than real
// detections: paypal.com contains "paypal", login.gov "login", and pages like
// github.com/login or support.apple.com are where those words belong.
func isTopRankedDomain(rank int) bool {
	return rank > 0
}

// isBareURL reports whether the URL has no path, query, or fragment
// (e.g. "https://bit.ly" or "https://bit.ly/"). For a URL shortener this is
// the service's own homepage, not a short link hiding a destination.
func isBareURL(rawURL string) bool {
	u, err := neturl.Parse(rawURL)
	if err != nil {
		return false
	}
	return strings.Trim(u.Path, "/") == "" && u.RawQuery == "" && u.Fragment == ""
}

// applyExceptions clears signals that are known false positives, so the score,
// the flags shown in the UI, and the cached result all agree.
func applyExceptions(in *Input, out *Output) {
	out.mu.Lock()
	defer out.mu.Unlock()

	// The bare address of a provider's own endpoint (https://storage.googleapis.com/)
	// is the provider, so it takes the owning domain's rank. Anything under a
	// path is an upload and stays unranked.
	if _, ok := constants.ProviderServiceHosts[hostOf(in.URL)]; ok && isBareURL(in.URL) && out.Rank == 0 {
		out.Rank = lookupRank(out.TLD)
	}

	// A user's page on a link-in-bio service (linktr.ee/someone) is theirs,
	// not the service's: no rank, and it's judged like a hosted site.
	if host := strings.TrimPrefix(hostOf(in.URL), "www."); !isBareURL(in.URL) {
		if _, ok := constants.UserPageHosts[host]; ok {
			out.Rank = 0
			out.TLD = host
			out.TLDIsHostingPlatform = true
		}
	}

	// Restricted registries (.gov, .bank.in, .edu, …) verify who registers a
	// domain, so the same name-based heuristics don't apply there either.
	// Brand names in the URL. A well-known site (top 100,000) or a restricted
	// registry is trusted with them entirely. The rest of the top million also
	// holds abused and parked domains (binanceuz.co is ranked) next to honest
	// ones (killedbygoogle.com, lufthansa-industry-solutions.com), so there a
	// brand in the name only counts with something else pointing the same way:
	// a domain under two years old, a login or payment form, or words like
	// "login" in the link. A brand worked into a subdomain (paypal-login.…)
	// always counts; the plain name on another ending (roblox.dk) never does.
	ts := &out.TyposquatResult
	switch {
	case out.TLDTrusted && ts.PathBrand != "" && loginPath(in.URL):
		// A brand's login under a .gov or .edu address (/secure/netflix/login)
		// is a phishing kit on a hacked site; keep that, drop the rest.
		*ts = typosquat.TyposquatResult{PathBrand: ts.PathBrand}
	case (out.Rank > 0 && out.Rank <= wellKnownRank) || out.TLDTrusted:
		*ts = typosquat.TyposquatResult{}
	case out.Rank > 0:
		age, known := out.DomainInfo.Age()
		asks := out.ContentData != nil && (out.ContentData.HasLoginForm || out.ContentData.HasPaymentForm)
		corroborated := (known && age <= 2*oneYearDays) || asks || out.URLKeywordsPresent
		if ts.IsSuspicious && (!ts.IsComboSquat || ts.ComboExact || !corroborated) {
			ts.IsSuspicious, ts.IsComboSquat, ts.ComboExact = false, false, false
		}
		if !corroborated {
			ts.PathBrand = ""
			if ts.SubdomainBrandExact {
				ts.SubdomainBrand, ts.SubdomainBrandExact = "", false
			}
		}
	}

	if isTopRankedDomain(out.Rank) || out.TLDTrusted {
		out.TLDRisky = false
		out.URLKeywordsPresent = false
		out.URLKeywordMatches = nil
		out.URLKeywordCats = nil
	}

	// Long, deep URLs are normal on big sites (release downloads, docs, search
	// results); padding a URL only hides anything on a domain nobody knows.
	if isTopRankedDomain(out.Rank) {
		out.URLTooLong = false
		out.URLTooDeep = false
	}

	// Only a verified organisation can register on a restricted registry, so a
	// bank's page on hdfc.bank.in naming HDFC isn't impersonation, even before
	// the brand list learns the new address.
	if out.TLDTrusted && out.ContentData != nil && out.ContentData.BrandCheck.IsMismatch {
		brand := *out.ContentData
		brand.BrandCheck.IsMismatch = false
		brand.BrandCheck.BrandFound = ""
		out.ContentData = &brand
	}

	applyBrandContext(in, out)

	// A reported page elsewhere on the same host only implicates this page
	// when one owner controls the whole host: a customer's own subdomain on a
	// hosting platform, or a domain registered within the year (bought for the
	// campaign). Shared services (github.com, form builders like formbuddy.com)
	// and old, possibly hacked sites carry other people's pages too.
	if fm := out.ThreatFeeds; fm != nil && fm.Match == "host" {
		host := hostOf(in.URL)
		age, known := out.DomainInfo.Age()
		ownedHost := (out.TLDIsHostingPlatform && host != out.TLD && !isUserContentHost(host)) ||
			(known && age <= oneYearDays && !isTopRankedDomain(out.Rank))
		if !ownedHost {
			out.ThreatFeeds = &threatfeeds.FeedMatch{}
		}
	}

	// A well-known site moving to its own new address (zoom.us → zoom.com,
	// hdfcbank.com → hdfc.bank.in) isn't the kind of jump phishing makes: the
	// destination is either well-known too, or on a registry that only verified
	// organisations can use.
	if isTopRankedDomain(out.Rank) && out.RedirectionResult.HasDomainJump &&
		(isVouchedDestination(out.RedirectionResult.FinalURL) ||
			sameBrand(in.Domain, out.RedirectionResult.FinalURLHost)) {
		out.RedirectionResult.HasDomainJump = false
	}
	// Blogger moves a blog between its country addresses
	// (sbdrg.blogspot.fi → sbdrg.blogspot.com): the same blog.
	if out.RedirectionResult.HasDomainJump && sameBlog(hostOf(in.URL), out.RedirectionResult.FinalURLHost) {
		out.RedirectionResult.HasDomainJump = false
	}

	// A page on a brand's own hosting forwarding to the brand's own site
	// (googleresearch.blogspot.com → research.google): the brand moved it.
	// Its name isn't impersonation either: Google naming a Blogspot address
	// "googleresearch" is Google's own.
	if out.TLDIsHostingPlatform && out.RedirectionResult.HasDomainJump &&
		sameBrand(out.TLD, out.RedirectionResult.FinalURLHost) {
		out.RedirectionResult.HasDomainJump = false
		ts := &out.TyposquatResult
		if ts.IsComboSquat && brandRuns(ts.MatchedDomain, out.RedirectionResult.FinalURLHost) {
			*ts = typosquat.TyposquatResult{}
		}
	}

	if out.URLIsShortener && isBareURL(in.URL) {
		out.URLIsShortener = false
	}
}

// wellKnownRank is the rank within which a site is established enough that a
// brand in its title is news, a review or a reseller, not impersonation.
const wellKnownRank = 100000

// applyBrandContext adjusts the title brand check with what only the full scan
// knows: how established the site is and where it's hosted.
func applyBrandContext(in *Input, out *Output) {
	if out.ContentData == nil {
		return
	}
	content := *out.ContentData
	host := hostOf(content.URL)
	if host == "" {
		host = hostOf(in.URL)
	}
	// The title is the page the visitor ends up on, so judge it against that
	// page's host: a redirect to the brand's own site shows the brand's page.
	if final := out.RedirectionResult.FinalURLHost; final != "" {
		host = final
	}
	// Re-run the title check: page content is cached for a while, and brand
	// list changes should apply to the next scan, not after the cache expires.
	if content.Title != "" {
		content.BrandCheck = checks.CheckBrandMismatch(host, content.Title)
		out.ContentData = &content
	}
	// A link-in-bio page names the networks it links to ("Instagram &
	// Twitter Links | Linktree"); that's only telling where it asks for a login
	// or payment.
	if _, ok := constants.UserPageHosts[strings.TrimPrefix(host, "www.")]; ok && !content.HasLoginForm && !content.HasPaymentForm {
		content.BrandCheck.IsMismatch = false
		content.BrandCheck.BrandFound = ""
		content.BrandCheck.OfficialDomain = ""
		content.BrandCheck.DetectedNames = nil
		out.ContentData = &content
		return
	}

	switch {
	// "Netflix raises prices" on a major news site, or dhl.fr (which isn't in
	// DHL's list) naming DHL. Anyone can publish on a user-content host, so
	// those keep the check.
	case content.BrandCheck.IsMismatch && out.Rank > 0 && out.Rank <= wellKnownRank && !isUserContentHost(host):
		content.BrandCheck.IsMismatch = false
		content.BrandCheck.BrandFound = ""
		content.BrandCheck.OfficialDomain = ""
		// Naming the brand isn't being the brand: no "verified" match either.
		content.BrandCheck.DetectedNames = nil

	// A bare brand name ("Facebook", "Apple — Claim your gift card") is
	// telling on an unknown site that asks for a login or payment, or on a
	// free hosting subdomain where anyone can put up a page in seconds.
	case !content.BrandCheck.IsMismatch && out.Rank == 0 && !out.TLDTrusted &&
		(out.TLDIsHostingPlatform || content.HasLoginForm || content.HasPaymentForm):
		loose := checks.CheckBrandNames(host, content.Title)
		if !loose.IsMismatch {
			return
		}
		content.BrandCheck.IsMismatch = true
		content.BrandCheck.BrandFound = loose.BrandFound
		content.BrandCheck.OfficialDomain = loose.OfficialDomain

	default:
		return
	}
	out.ContentData = &content
}

// isUserContentHost reports whether anyone can publish under this host:
// file and repo hosts, Google Sites, provider storage endpoints.
func isUserContentHost(host string) bool {
	return userContentHost(host) != ""
}

// userContentHost returns the UserUploadHosts entry the host falls under, or "".
func userContentHost(host string) string {
	if _, ok := constants.ProviderServiceHosts[host]; ok {
		return host
	}
	for d := strings.ToLower(host); d != ""; {
		if _, ok := constants.UserUploadHosts[d]; ok {
			return d
		}
		i := strings.Index(d, ".")
		if i < 0 {
			break
		}
		d = d[i+1:]
	}
	return ""
}

// sameBrand reports whether both hosts belong to one brand's official
// domains, e.g. github.com handing a download to release-assets.githubusercontent.com.
func sameBrand(fromDomain, toHost string) bool {
	for _, entry := range constants.HighValueBrands {
		if hostInDomains(fromDomain, append(entry.Domains(), entry.Platforms...)) && hostInDomains(toHost, entry.Domains()) {
			return true
		}
	}
	return false
}

// loginPathWords are path words for signing in or verifying an account.
var loginPathWords = regexp.MustCompile(`(?i)(log[-_]?in|log[-_]?on|sign[-_]?in|signon|verif|auth|secure|account|webscr|wallet|update)`)

// loginPath reports whether a link's path is about signing in.
func loginPath(rawURL string) bool {
	u, err := neturl.Parse(rawURL)
	return err == nil && loginPathWords.MatchString(u.Path)
}

// sameBlog reports whether two hosts are one Blogger blog on two of its
// addresses (sbdrg.blogspot.fi and sbdrg.blogspot.com).
func sameBlog(a, b string) bool {
	blog := func(h string) string {
		h = strings.TrimPrefix(strings.ToLower(h), "www.")
		if name, ok := strings.CutSuffix(h, ".blogspot.com"); ok {
			return name
		}
		if p := checks.CustomerSitePlatform(h); strings.HasPrefix(p, "blogspot.") {
			return strings.TrimSuffix(h, "."+p)
		}
		return ""
	}
	return blog(a) != "" && blog(a) == blog(b)
}

// brandRuns reports whether the brand that owns officialDomain also runs host.
func brandRuns(officialDomain, host string) bool {
	for _, entry := range constants.HighValueBrands {
		if hostInDomains(officialDomain, entry.Domains()) && hostInDomains(host, entry.Domains()) {
			return true
		}
	}
	return false
}

func hostInDomains(host string, domains []string) bool {
	host = strings.ToLower(host)
	for _, d := range domains {
		if host == d || strings.HasSuffix(host, "."+d) {
			return true
		}
	}
	return false
}

// rankLookupForVouch is rank.DomainRankLookup, swappable in tests.
var rankLookupForVouch = rank.DomainRankLookup

// isVouchedDestination reports whether a redirect target is well-known (top-1M)
// or on a restricted registry such as .gov or .bank.in.
func isVouchedDestination(finalURL string) bool {
	domain, err := checks.GetDomain(finalURL)
	if err != nil || domain == "" {
		return false
	}
	if rankLookupForVouch(domain) > 0 {
		return true
	}
	trusted, _, _ := checks.IsTrustedTld(domain)
	return trusted
}
