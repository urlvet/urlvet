package analyzer

import (
	"encoding/base64"
	"fmt"
	"math"
	neturl "net/url"
	"path"
	"regexp"
	"slices"
	"strings"

	"github.com/urlvet/urlvet/internal/constants"
	"github.com/urlvet/urlvet/internal/service/checks"
	"github.com/urlvet/urlvet/internal/service/threatfeeds"
)

func GenerateResult(resp Response) Result {
	var neutralReasons []string
	var goodReasons []string
	var badReasons []string
	trustScore := 0
	riskScore := 0

	// --- 1. Popularity & Rank ---
	if resp.Features.Rank == 0 {
		if service, ok := constants.ProviderServiceHosts[hostOf(resp.URL)]; ok {
			neutralReasons = append(neutralReasons, fmt.Sprintf("File hosted on %s (anyone can upload here).", service))
		} else if service, ok := constants.UserPageHosts[strings.TrimPrefix(hostOf(resp.URL), "www.")]; ok {
			neutralReasons = append(neutralReasons, fmt.Sprintf("A page on %s, where anyone can make one.", service))
		} else if store := userContentHost(hostOf(resp.URL)); store != "" && resp.Features.TLD.IsHostingPlatform {
			neutralReasons = append(neutralReasons, fmt.Sprintf("Stored on %s, where anyone can upload files.", store))
		} else if resp.Features.TLD.IsHostingPlatform {
			// Personal/project pages on known hosting platforms are inherently unranked — not a risk signal.
			neutralReasons = append(neutralReasons, fmt.Sprintf("Unranked subdomain on %s (normal for personal or project pages).", resp.Features.TLD.TLD))
		} else if age, ok := resp.DomainInfo.Age(); ok && age > establishedAfterDays {
			// Most real sites never make the top million. Being unranked only
			// points at phishing alongside a domain registered weeks ago.
			neutralReasons = append(neutralReasons, "Not among the top 1M sites (normal for small sites).")
		} else {
			badReasons = append(badReasons, "Very low traffic volume.")
			riskScore += 10
		}
	} else if resp.Features.Rank > 0 && resp.Features.Rank <= 10000 {
		goodReasons = append(goodReasons, fmt.Sprintf("Global Giant: Ranked #%d worldwide.", resp.Features.Rank))
		trustScore += 90
	} else if resp.Features.Rank <= 50000 {
		// 10,001–50,000: well-known sites
		goodReasons = append(goodReasons, fmt.Sprintf("Well-known website (#%d worldwide).", resp.Features.Rank))
		trustScore += 50
	} else {
		// 50,001+: ranked but lower traffic
		goodReasons = append(goodReasons, fmt.Sprintf("Low-traffic but indexed website (#%d worldwide).", resp.Features.Rank))
		trustScore += 20
	}

	// --- 2. TLD (Top Level Domain) ---
	if resp.Features.TLD.IsRisky {
		badReasons = append(badReasons, "High-risk domain extension detected (often associated with spam).")
		riskScore += 20
	}

	// Exceptions keep a brand in the path on .gov or .edu only for a login
	// page (/secure/netflix/login): a hacked site, which the registry can't
	// vouch for.
	if resp.Features.TLD.IsTrusted && resp.TyposquatResult.PathBrand != "" {
		badReasons = append(badReasons, fmt.Sprintf("A '%s' sign-in page under a government or university address, which suggests the site was hacked.", resp.TyposquatResult.PathBrand))
		riskScore += 60
	} else if resp.Features.TLD.IsTrusted {
		goodReasons = append(goodReasons, "High-trust official domain extension (Gov/Edu).")
		trustScore += 100
	} else if resp.Features.TLD.IsICANN && !resp.Features.TLD.IsRisky {
		// Keep it simple
		neutralReasons = append(neutralReasons, "Standard, officially recognized domain extension.")
	}

	if !resp.Features.TLD.IsICANN {
		if resp.Features.TLD.IsHostingPlatform {
			// PSL private entry but operated by a reputable hosting company — small trust signal.
		} else {
			badReasons = append(badReasons, "Unregulated or non-standard domain extension.")
			riskScore += 30
		}
	}

	// --- 3. Security Protocols ---
	// HSTS is one switch at a CDN; on a domain registered this month it says
	// nothing about who runs it.
	if age, known := resp.DomainInfo.Age(); resp.Analysis.SupportsHSTS && known && age <= newDomainDays {
		neutralReasons = append(neutralReasons, "Sends an HSTS header, which on a domain this new says little.")
	} else if resp.Analysis.SupportsHSTS {
		goodReasons = append(goodReasons, "Enforces strict HTTPS security (HSTS Enabled).")
		trustScore += 20
	}

	// --- 4. URL Structure / Obfuscation ---
	if resp.Features.URL.IsURLShortener {
		badReasons = append(badReasons, "URL Shortener detected (hides the true destination).")
		riskScore += 25
	}

	// Uses IP
	if resp.Features.URL.UsesIP {
		badReasons = append(badReasons, "Raw IP address usage detected (common evasion tactic).")
		riskScore += 100
	}
	// The same, dressed up as a name: 1-2-3-4.sslip.io is 1.2.3.4.
	if service, ip, ok := checks.IPHostname(hostOf(resp.URL)); ok {
		if ip != "" {
			badReasons = append(badReasons, fmt.Sprintf("The address is a bare server (%s) written as a name on %s, which turns any IP address into a name.", ip, service))
		} else {
			badReasons = append(badReasons, fmt.Sprintf("The address is on %s, which turns any IP address into a name.", service))
		}
		riskScore += 100
	}

	// Punycode
	if resp.Features.URL.ContainsPunycode {
		badReasons = append(badReasons, "Punycode characters detected (potential phishing spoof).")
		riskScore += 100
	}

	// Too deep
	if resp.Features.URL.TooDeep {
		badReasons = append(badReasons, "Excessively deep URL path (potential request hiding).")
		riskScore += 30
	}

	// Too long
	if resp.Features.URL.TooLong {
		badReasons = append(badReasons, "URL length exceeds standard limits (potential buffer overflow/hiding).")
		riskScore += 20
	}

	// Subdomain Count
	if resp.Features.URL.SubdomainCount > 2 {
		badReasons = append(badReasons, "Suspicious number of subdomains detected.")
		riskScore += 15
	}

	// Keywords
	if resp.Features.URL.Keywords.HasKeywords {
		badReasons = append(badReasons, fmt.Sprintf("Sensitive security keywords found in URL: %s", strings.Join(resp.Features.URL.Keywords.Found, ", ")))
		riskScore += 10
	}

	// --- 5. Infrastructure Forensics ---
	if !resp.Infrastructure.NameserversValid {
		if resp.Features.TLD.IsHostingPlatform {
			// Hosting platforms manage the entire DNS zone; individual subdomains have no own NS records.
		} else {
			badReasons = append(badReasons, "Incomplete or missing DNS configuration.")
			riskScore += 10
		}
	}

	// MX records
	if !resp.Infrastructure.MXRecordsValid && !resp.Features.TLD.IsHostingPlatform {
		neutralReasons = append(neutralReasons, "No email server configured for this domain.")
		riskScore += 5
	}

	// A page on a brand's own hosting that forwards to the brand's own site
	// (googleresearch.blogspot.com → research.google) ends somewhere the
	// brand runs, which vouches for it.
	if final := resp.Analysis.RedirectionResult.FinalURLHost; resp.Features.TLD.IsHostingPlatform &&
		final != "" && final != hostOf(resp.URL) && sameBrand(resp.Features.TLD.TLD, final) && isVouchedDestination(resp.Analysis.RedirectionResult.FinalURL) {
		goodReasons = append(goodReasons, fmt.Sprintf("Leads to %s, run by the same company as %s.", final, resp.Features.TLD.TLD))
		trustScore += 50
	}

	// --- 6. Domain History ---
	// On a hosting platform the registration is the provider's, so its age
	// says nothing about the site.
	if resp.DomainInfo != nil && resp.Features.TLD.IsHostingPlatform && hostingRegistration(resp) {
		if _, known := resp.DomainInfo.Age(); known {
			neutralReasons = append(neutralReasons, fmt.Sprintf("The domain's age is %s's, not this site's.", resp.Features.TLD.TLD))
		}
	} else if resp.DomainInfo != nil {
		// Years of clean operation is real evidence, and should count for more
		// than a single header like HSTS. A domain whose registry publishes no
		// creation date gets neither the penalty nor the bonus.
		age, known := resp.DomainInfo.Age()
		ageHuman := resp.DomainInfo.AgeHuman
		switch {
		case !known:
			neutralReasons = append(neutralReasons, "Registration date not published by this domain's registry.")
		case age <= newDomainDays:
			badReasons = append(badReasons, fmt.Sprintf("Registered %s. Phishing sites are usually brand new.", ago(ageHuman)))
			riskScore += 25
		case age <= establishedAfterDays:
			badReasons = append(badReasons, fmt.Sprintf("Young domain, registered %s. Use caution.", ago(ageHuman)))
			riskScore += 15
		case age <= oneYearDays:
			// Phishing domains are used within weeks of registration and burn
			// out fast; one that has lasted a few months isn't that profile.
			neutralReasons = append(neutralReasons, fmt.Sprintf("Fairly new domain (%s old).", ageHuman))
		case (resp.Features.Rank == 0 || resp.Features.Rank > wellKnownRank) && kitLoginPath(resp.URL):
			// Hacked sites keep their age: a phishing kit dropped on an old
			// domain borrows years of history it had no part in.
			neutralReasons = append(neutralReasons, fmt.Sprintf("The domain is %s old, but this link is a sign-in page that looks planted in it, so its age isn't counted.", ageHuman))
		case age <= threeYearsDays:
			goodReasons = append(goodReasons, fmt.Sprintf("Operational for %s.", ageHuman))
			trustScore += 10
		case age <= fiveYearsDays:
			goodReasons = append(goodReasons, fmt.Sprintf("Operational for %s.", ageHuman))
			trustScore += 15
		default:
			goodReasons = append(goodReasons, fmt.Sprintf("Long-standing domain history (%s).", ageHuman))
			trustScore += 25
		}

		// DNSSEC Logic Updated
		if resp.DomainInfo.DNSSEC {
			goodReasons = append(goodReasons, "Advanced DNS security enabled (DNSSEC).")
			trustScore += 10
		} else {
			// Moved to Neutral. Not having DNSSEC is NOT a sign of phishing for .coms
			neutralReasons = append(neutralReasons, "Standard DNS security (DNSSEC not enabled).")
			// Removed riskScore penalty
		}
	}

	// --- 7. Redirection Analysis ---
	if resp.Analysis.RedirectionResult.IsRedirected {
		// ChainLength counts URLs, including the one scanned. Long chains only
		// matter when they leave the site: hopping www → locale → home on one
		// domain is normal, while phishing bounces through other domains.
		hops := resp.Analysis.RedirectionResult.ChainLength - 1
		if hops > 3 && resp.Analysis.RedirectionResult.HasDomainJump {
			badReasons = append(badReasons, fmt.Sprintf("Excessive redirection chain detected (%d hops).", hops))
			riskScore += 40
		}

		if resp.Analysis.RedirectionResult.HasDomainJump {
			badReasons = append(badReasons, "Cross-domain redirection detected (destination differs from source).")
			// Add the destination as a neutral fact so they can see where they are going
			badReasons = append(badReasons, fmt.Sprintf("Final Destination: %s. Check Report for more info.", resp.Analysis.RedirectionResult.FinalURLHost))
			riskScore += 50
		}
	}

	// --- 8. Homoglyphs ---
	if resp.Features.URL.HasHomoglyph {
		badReasons = append(badReasons, "Homoglyph attack detected (deceptive visual characters).")
		riskScore += 60
	}

	// --- 9. Typosquatting / Combo-squatting ---
	if resp.TyposquatResult.IsSuspicious {
		if resp.TyposquatResult.IsComboSquat {
			badReasons = append(badReasons, fmt.Sprintf("Combo-squatting detected: domain contains brand name '%s' but is not the official site.", resp.TyposquatResult.MatchedBrand))
			// On a hosting platform anyone can claim "instagram-login" for free
			// in seconds; a registered domain at least costs something.
			if resp.Features.TLD.IsHostingPlatform {
				riskScore += 40
			} else {
				riskScore += 20
			}
			if bare := brandPageIsBare(resp); bare != "" {
				badReasons = append(badReasons, fmt.Sprintf("Uses the name '%s' on a little-known site, and the page %s.", resp.TyposquatResult.MatchedBrand, bare))
				riskScore += 20
			}
		} else {
			badReasons = append(badReasons, fmt.Sprintf("Typosquatting detected: domain closely resembles '%s' (%d character difference).", resp.TyposquatResult.MatchedDomain, resp.TyposquatResult.Distance))
			riskScore += 40
		}
	}

	// Brands in the rest of the URL. Only reached on sites that aren't well
	// known: applyExceptions clears these for those.
	ts := resp.TyposquatResult
	if ts.EmbeddedDomain != "" {
		badReasons = append(badReasons, fmt.Sprintf("The address starts like %s but belongs to %s.", ts.EmbeddedDomain, resp.Domain))
		riskScore += 60
	}
	if ts.SubdomainBrand != "" && ts.EmbeddedDomain == "" {
		if bare := brandPageIsBare(resp); bare != "" {
			badReasons = append(badReasons, fmt.Sprintf("A subdomain is named after '%s' on a little-known site, and the page %s.", ts.SubdomainBrand, bare))
			riskScore += 20
		}
		if ts.SubdomainBrandExact {
			// github.acme.com is a common, honest pattern; it counts for less.
			badReasons = append(badReasons, fmt.Sprintf("A subdomain is named '%s', on a site that isn't theirs.", ts.SubdomainBrand))
			riskScore += 20
		} else {
			badReasons = append(badReasons, fmt.Sprintf("A subdomain uses the brand name '%s', on a site that isn't theirs.", ts.SubdomainBrand))
			riskScore += 40
		}
	}
	if ts.PathBrand != "" {
		if resp.ContentData != nil && resp.ContentData.HasLoginForm {
			badReasons = append(badReasons, fmt.Sprintf("The link's path names '%s' and the page asks you to log in, on a site that isn't theirs.", ts.PathBrand))
			riskScore += 40
		} else {
			badReasons = append(badReasons, fmt.Sprintf("The link's path names '%s', on a site that isn't theirs.", ts.PathBrand))
			riskScore += 20
		}
	}

	// --- 10. Threat Intelligence (PhishTank) ---
	// PhishTank semantics:
	//   valid=true  → URL IS phishing (the actual threat signal)
	//   valid=false → URL is NOT phishing
	//   verified    → whether the report has been reviewed at all (not a phishing indicator)
	if resp.Phishing != nil && resp.Phishing.InDatabase {
		if resp.Phishing.Valid && resp.Phishing.Verified {
			// Reviewed and confirmed as phishing — highest risk.
			badReasons = append(badReasons, "CONFIRMED PHISHING: This is a verified phishing URL.")
			riskScore += 200
			if resp.Phishing.Target != "" {
				badReasons = append(badReasons, fmt.Sprintf("Reported Target: %s", resp.Phishing.Target))
			}
		} else if resp.Phishing.Valid && !resp.Phishing.Verified {
			// Reported as phishing but not yet reviewed.
			badReasons = append(badReasons, "This URL has been reported as phishing, awaiting community verification.")
			riskScore += 70
		}
		// valid=false (regardless of verified) → confirmed NOT phishing → no penalty.
		// A verified=true, valid=false entry means PhishTank reviewed it and cleared it.
	}

	// A list naming this exact link outweighs anything good about its site:
	// the site's age and rank say nothing about one page on it.
	confirmed := false

	// Locally held feeds (PhishTank's dump, OpenPhish, URLhaus).
	if fm := resp.ThreatFeeds; fm != nil && fm.Listed {
		sources := strings.Join(fm.Sources, " and ")
		if fm.Match == "url" {
			badReasons = append(badReasons, fmt.Sprintf("CONFIRMED THREAT: This link is listed as phishing or malware by %s.", sources))
			riskScore += 200
			confirmed = true
		} else {
			badReasons = append(badReasons, fmt.Sprintf("Other pages on this site are listed as phishing or malware by %s.", sources))
			riskScore += 80
		}
	}

	// Google Web Risk or Safe Browsing (the same lists; one reason is enough).
	for _, g := range []*threatfeeds.GoogleThreatResult{resp.WebRisk, resp.SafeBrowsing} {
		if g != nil && g.Listed {
			// Google's terms require this attribution wherever its verdict is shown.
			badReasons = append(badReasons, fmt.Sprintf("CONFIRMED THREAT: Google lists this link as %s (advisory provided by Google).", webRiskThreatNames(g.ThreatTypes)))
			riskScore += 200
			confirmed = true
			break
		}
	}

	// A target's address carried in the link, for the page to pre-fill.
	if victimAddressInLink(resp.URL) {
		badReasons = append(badReasons, "The link carries an email address for the page to fill in, a common way phishing pages tailor a login to each target.")
		riskScore += 20
	}

	// --- 11. Page Content & Phishing Signals ---
	if resp.ContentData != nil && resp.Features.Rank == 0 && tradingScamTitle.MatchString(resp.ContentData.Title) {
		badReasons = append(badReasons, "The page uses the title template of fake crypto trading platforms (\"AI Trading Platform | Official Website\"), a common investment scam.")
		riskScore += 50
	}
	if resp.ContentData != nil {
		if resp.ContentData.HasLoginForm {
			// Check if domain is established
			isEstablished := resp.Features.Rank > 0 && resp.Features.Rank <= 100000
			age, known := resp.DomainInfo.Age()
			isOld := known && age > oneYearDays

			if !isEstablished && !isOld {
				badReasons = append(badReasons, "SUSPICIOUS: Login form detected on a new or unranked domain.")
				riskScore += 50
			} else {
				neutralReasons = append(neutralReasons, "Page contains a login form.")
			}
			// Mail providers don't host their sign-in pages on free hosting.
			if resp.Features.TLD.IsHostingPlatform && mailLoginTitle.MatchString(resp.ContentData.Title) {
				badReasons = append(badReasons, fmt.Sprintf("An email sign-in page on %s, free hosting where real mail providers don't put their logins.", resp.Features.TLD.TLD))
				riskScore += 30
			}
		}

		if resp.ContentData.HasPaymentForm {
			badReasons = append(badReasons, "WARNING: Payment-related fields detected (credit card, CVV, etc.).")
			riskScore += 30
		}

		if resp.ContentData.HasPersonalForm {
			neutralReasons = append(neutralReasons, "Page requests personal information (address, phone, etc.).")
		}

		if resp.ContentData.HasHiddenIframe {
			badReasons = append(badReasons, "WARNING: Hidden iframe detected (often used for background credential theft or clickjacking).")
			riskScore += 40
		}

		if resp.ContentData.HasTracking {
			neutralReasons = append(neutralReasons, "Background tracking elements (1x1 pixels) detected.")
		}

		if pb := resp.ContentData.ProviderBlock; pb != nil {
			switch pb.Reason {
			case "phishing", "malware":
				badReasons = append(badReasons, fmt.Sprintf("%s has flagged this page as suspected %s and shows a warning instead of it.", pb.Provider, pb.Reason))
				riskScore += 100
			default:
				badReasons = append(badReasons, fmt.Sprintf("%s has taken this site down.", pb.Provider))
				riskScore += 40
			}
		}

		if sr := resp.ContentData.ScriptRedirect; sr != nil {
			if sr.ToIP {
				badReasons = append(badReasons, "Page uses a script to send visitors on to a raw IP address (common in spam and phishing redirectors).")
				riskScore += 60
			} else if sr.CrossDomain && isUserContentHost(hostOf(resp.URL)) {
				// A file in cloud storage that only forwards visitors is a
				// throwaway hop: it borrows the provider's name to get past
				// filters and points at the real page somewhere else.
				badReasons = append(badReasons, fmt.Sprintf("A file anyone could upload here immediately sends visitors on to another site: %s.", hostOf(sr.Target)))
				riskScore += 60
			} else if sr.CrossDomain {
				badReasons = append(badReasons, fmt.Sprintf("Page immediately sends visitors on to another site: %s.", hostOf(sr.Target)))
				riskScore += 30
			}
		}

		if resp.ContentData.BrandCheck.IsMismatch {
			badReasons = append(badReasons, fmt.Sprintf("BRAND MISMATCH: Page mentions '%s' but is hosted on an unofficial domain.", resp.ContentData.BrandCheck.BrandFound))
			riskScore += 100
		} else if len(resp.ContentData.BrandCheck.DetectedNames) > 0 {
			goodReasons = append(goodReasons, fmt.Sprintf("Verified brand matching: %s", strings.Join(resp.ContentData.BrandCheck.DetectedNames, ", ")))
			trustScore += 20
		}

		if resp.ContentData.HasForms {
			for _, form := range resp.ContentData.Forms {
				// Where a form sends its data only matters for what it collects:
				// site search often posts to a search service on another domain.
				// Newsletter sign-ups post to email-marketing platforms
				// (Mailchimp, HubSpot, Eloqua…): only a password or card
				// going there would be out of place.
				marketing := checks.IsMarketingFormService(form.Action) && !form.ContainsPassword && !form.ContainsPayment
				if form.ExternalAction && !marketing {
					switch {
					case form.ContainsPassword || form.ContainsPayment || form.ContainsPersonal:
						badReasons = append(badReasons, "CRITICAL: Form submits data to a different domain (common phishing tactic).")
						riskScore += 80
					case form.ContainsUserLike:
						badReasons = append(badReasons, "Form sends an email address or username to a different domain.")
						riskScore += 20
					}
				}
				if form.ContainsPassword && !resp.SSLInfo.HasTLS {
					badReasons = append(badReasons, "DANGEROUS: Password form detected over insecure connection!")
					riskScore += 200
				}
			}
		}
	}

	// --- 12. Certificate ---
	// Browsers refuse a certificate issued for another name with a full-page
	// warning; a real site doesn't ship one. Phishing kits on hacked hosts and
	// shared servers often do.
	if resp.TLSInfo.Present && resp.TLSInfo.HostnameMismatch && servedOverHTTPS(resp) {
		badReasons = append(badReasons, "Security certificate doesn't belong to this site; browsers show a warning before opening it.")
		riskScore += 40
	}

	// --- 13. Hacked-site Paths ---
	// Phishing kits are dropped into WordPress system folders on compromised
	// sites. Images, styles and scripts live there legitimately; pages don't.
	if isWordPressSystemPage(resp.URL) {
		badReasons = append(badReasons, "Page sits inside a WordPress system folder, where phishing kits are hidden on hacked sites.")
		riskScore += 30
		if resp.ContentData != nil && resp.ContentData.HasLoginForm {
			badReasons = append(badReasons, "That page asks visitors to log in.")
			riskScore += 30
		}
	}

	// --- 14. Program Downloads ---
	// A well-known host can't vouch for a program any account uploaded to it.
	// Its reputation would carry the score to Safe, so cap it instead.
	capAtSuspicious := false
	if ext := downloadedProgram(resp); ext != "" {
		switch host := uploadHost(resp); {
		case host != "":
			badReasons = append(badReasons, fmt.Sprintf("Downloads a program (%s) that a user uploaded to %s. Nothing confirms it comes from the software's publisher.", ext, host))
			capAtSuspicious = true
		case resp.Features.Rank == 0:
			badReasons = append(badReasons, fmt.Sprintf("Downloads a program (%s) from a little-known site.", ext))
			riskScore += 40
		default:
			neutralReasons = append(neutralReasons, fmt.Sprintf("Link downloads a program (%s).", ext))
		}
	}

	// A short link we couldn't follow could lead anywhere; the shortener's
	// own reputation says nothing about it.
	if resp.ShortLink != nil && !resp.ShortLink.Resolved {
		badReasons = append(badReasons, "Couldn't see where this short link leads, so it can't be called safe.")
		capAtSuspicious = true
	}

	// --- 15. The link as given, when it led elsewhere ---
	// Trust above came from the page the link lands on; what the link itself
	// and the hops on the way get wrong still counts.
	if o := resp.Origin; o != nil {
		badReasons = append(badReasons, o.BadReasons...)
		riskScore += o.Risk
		if o.Confirmed {
			confirmed = true
		}
	}

	// --- Normalize / cap scores ---
	riskScore = clamp(riskScore)
	trustScore = clamp(trustScore)

	// Balanced formula: neutral point is 50.
	// trustScore pulls score up, riskScore pulls it down.
	// A site with no signals either way scores 50 → "Suspicious".
	finalScore := clamp(50 + int(math.Round(float64(trustScore-riskScore)*0.5)))
	if capAtSuspicious && finalScore >= safeThreshold {
		finalScore = safeThreshold - 1
	}
	if confirmed && finalScore >= riskyThreshold {
		finalScore = riskyThreshold - 1
	}

	var verdict string
	switch {
	// High risk, low trust
	case finalScore < riskyThreshold:
		verdict = "Risky"
	// Moderate risk OR conflicting/insufficient signals
	case finalScore < safeThreshold:
		verdict = "Suspicious"
	// Low risk, sufficient trust
	default:
		verdict = "Safe"
	}

	res := Result{
		RiskScore:  riskScore,
		TrustScore: trustScore,
		FinalScore: finalScore,
		Verdict:    verdict,
		Reasons: Reasons{
			NeutralReasons: neutralReasons,
			GoodReasons:    goodReasons,
			BadReasons:     badReasons,
		},
	}

	return res
}

// safeThreshold is the lowest final score that reads as Safe.
const safeThreshold = 65

// riskyThreshold: scores below it are Risky.
const riskyThreshold = 30

// Domain age bands, in days since registration.
const (
	// newDomainDays covers the weeks a phishing domain is typically used.
	newDomainDays = 30
	// establishedAfterDays is the age past which a domain no longer reads as
	// freshly registered for phishing.
	establishedAfterDays = 90
	oneYearDays          = 365
	threeYearsDays       = 3 * 365
	fiveYearsDays        = 5 * 365
)

// ago phrases an age for "registered …": "today", or "12 days ago".
func ago(ageHuman string) string {
	if ageHuman == "today" {
		return ageHuman
	}
	return ageHuman + " ago"
}

// downloadedProgram returns the executable extension of the file the link
// downloads, from the scanned URL, where it redirects, or the served file name.
func downloadedProgram(resp Response) string {
	candidates := []string{resp.URL, resp.Analysis.RedirectionResult.FinalURL}
	if resp.ContentData != nil {
		candidates = append(candidates, resp.ContentData.FileName)
	}
	for _, c := range candidates {
		if ext := checks.ExecutableExt(c); ext != "" {
			return ext
		}
	}
	return ""
}

// uploadHost returns the user-upload host the link or its download sits on,
// or "" when neither does.
func uploadHost(resp Response) string {
	if resp.Features.TLD.IsHostingPlatform {
		return resp.Features.TLD.TLD
	}
	for _, h := range []string{hostOf(resp.URL), resp.Analysis.RedirectionResult.FinalURLHost} {
		if d := userContentHost(h); d != "" {
			return d
		}
	}
	return ""
}

var wpSystemPathRe = regexp.MustCompile(`(?i)/wp-(?:content|includes|admin)/`)

// staticAssetExts are file types WordPress legitimately serves from its folders.
var staticAssetExts = map[string]struct{}{
	".css": {}, ".js": {}, ".mjs": {}, ".map": {}, ".png": {}, ".jpg": {}, ".jpeg": {},
	".gif": {}, ".webp": {}, ".avif": {}, ".svg": {}, ".ico": {}, ".bmp": {},
	".woff": {}, ".woff2": {}, ".ttf": {}, ".otf": {}, ".eot": {},
	".mp4": {}, ".webm": {}, ".mp3": {}, ".pdf": {}, ".zip": {}, ".json": {}, ".xml": {}, ".txt": {},
}

// isWordPressSystemPage reports whether a URL is a page (not a static file)
// under wp-content, wp-includes or wp-admin. WordPress's own login and admin
// entry points (/wp-login.php, /wp-admin/) aren't matched.
func isWordPressSystemPage(rawURL string) bool {
	u, err := neturl.Parse(rawURL)
	if err != nil || !wpSystemPathRe.MatchString(u.Path) {
		return false
	}
	p := strings.ToLower(strings.TrimSuffix(u.Path, "/"))
	for _, entry := range []string{"/wp-admin", "/wp-admin/index.php", "/wp-admin/admin-ajax.php", "/wp-admin/admin-post.php"} {
		if strings.HasSuffix(p, entry) {
			return false
		}
	}
	_, static := staticAssetExts[strings.ToLower(path.Ext(u.Path))]
	return !static
}

// servedOverHTTPS reports whether visitors reach the page over HTTPS, so a
// certificate problem is something they'd actually run into.
func servedOverHTTPS(resp Response) bool {
	return strings.HasPrefix(resp.URL, "https://") ||
		(strings.HasPrefix(resp.Analysis.RedirectionResult.FinalURL, "https://") &&
			resp.Analysis.RedirectionResult.FinalURLHost == hostOf(resp.URL))
}

// webRiskThreatNames turns Web Risk's threat types into plain words.
func webRiskThreatNames(types []string) string {
	names := map[string]string{
		"SOCIAL_ENGINEERING": "phishing",
		"MALWARE":            "malware",
		"UNWANTED_SOFTWARE":  "unwanted software",
		// Safe Browsing v5 adds this one.
		"POTENTIALLY_HARMFUL_APPLICATION": "a harmful app",
	}
	// Most serious first, whatever order Google lists them in.
	order := []string{"SOCIAL_ENGINEERING", "MALWARE", "UNWANTED_SOFTWARE", "POTENTIALLY_HARMFUL_APPLICATION"}
	var out []string
	for _, t := range order {
		if slices.Contains(types, t) {
			out = append(out, names[t])
		}
	}
	if len(out) == 0 {
		return "dangerous"
	}
	return strings.Join(out, " and ")
}

// hostOf returns the host of a URL, or the input itself when it has none (a bare IP).
func hostOf(raw string) string {
	if u, err := neturl.Parse(raw); err == nil && u.Hostname() != "" {
		return u.Hostname()
	}
	return raw
}

func clamp(score int) int {
	return int(math.Max(0, math.Min(100, float64(score))))
}

// brandPageIsBare says why a page on a little-known site using a brand's name
// looks like a phishing setup: it asks for a login or card, or has nothing to
// show (kits often serve a blank page or a directory listing to anyone but
// their targets). A company's own brand-named site or "github.acme.com" has a
// page of its own. "" when none of these hold or the site is ranked. A page
// that couldn't be loaded doesn't count: brands park spare domains
// (airfrance.es, emiratesnbd.ae) that time out or don't resolve.
func brandPageIsBare(resp Response) string {
	c := resp.ContentData
	if resp.Features.Rank > 0 || c == nil {
		return ""
	}
	switch {
	case c.FileName != "":
		return ""
	case c.HasLoginForm || c.HasPaymentForm:
		return "asks you to log in or pay"
	case strings.HasPrefix(c.Title, "Index of /"):
		return "is a bare file listing"
	case strings.TrimSpace(c.Title) == "" && !c.HasForms:
		return "is blank"
	}
	return ""
}

// tradingScamTitle matches the title template of fake "AI trading platform"
// investment sites ("Polriksobotiks | AI Trading Platform | Official
// Website", "Immediate Peak ™ | The Official …"). Kits churn out hundreds of
// names on one template.
var tradingScamTitle = regexp.MustCompile(`(?i)\btrading (platform|app|bot|software)\b.{0,20}\bofficial (website|site|gateway|page)\b|^immediate [a-z0-9 ]{2,20}(™|®)?\s*[|:-]\s*(the )?official\b`)

// mailLoginTitle matches page titles of email sign-in pages.
var mailLoginTitle = regexp.MustCompile(`(?i)\b(web ?mail|mail ?(authentication|verification|login|portal)|e-?mail (verification|authentication|login|sign[ -]?in|portal)|mailbox|password protect(ion|ed)|verify your (e-?mail|mailbox|account)|sign ?in ?to (your )?(web)?mail)\b`)

// emailToken is a whole email address; mailMergeField is a mass mailer's
// placeholder left unfilled ("#[[-Email-]]").
var (
	emailToken     = regexp.MustCompile(`(?i)^[a-z0-9._%+-]{1,64}@[a-z0-9-]+(\.[a-z0-9-]+)*\.[a-z]{2,}$`)
	mailMergeField = regexp.MustCompile(`(?i)\[\[-?e-?mail-?\]\]|\{\{\s*e-?mail\s*\}\}|##e-?mail##|%e-?mail%`)
)

// victimAddressInLink reports whether a link carries an email address in its
// fragment or, base64-encoded, in its query. A plain address in the query
// counts only under an odd name (?eta=, ?uid=): newsletter and account links
// carry one honestly as ?email= or ?e=, and on unsubscribe pages.
func victimAddressInLink(rawURL string) bool {
	u, err := neturl.Parse(rawURL)
	if err != nil {
		return false
	}
	if mailMergeField.MatchString(u.Fragment) {
		return true
	}
	split := func(s string) []string {
		return strings.FieldsFunc(s, func(r rune) bool { return strings.ContainsRune("#&=/?!;,", r) })
	}
	for _, tok := range split(u.Fragment) {
		if emailToken.MatchString(tok) || emailToken.MatchString(fromBase64(tok)) {
			return true
		}
	}
	mailingPage := mailingPath.MatchString(u.Path)
	for name, vals := range u.Query() {
		_, honest := honestAddressParams[strings.ToLower(name)]
		for _, v := range vals {
			if emailToken.MatchString(fromBase64(v)) || (!honest && !mailingPage && emailToken.MatchString(v)) {
				return true
			}
		}
	}
	return false
}

// honestAddressParams are query names under which sites pass an address
// openly; mailingPath matches pages that take one by design.
var (
	honestAddressParams = map[string]struct{}{
		"email": {}, "e-mail": {}, "e_mail": {}, "mail": {}, "e": {}, "em": {},
		"emailaddress": {}, "email_address": {}, "address": {}, "recipient": {},
		"subscriber": {}, "to": {}, "from": {}, "user": {}, "username": {},
		"login": {}, "contact": {}, "invite": {}, "invitee": {}, "q": {},
	}
	mailingPath = regexp.MustCompile(`(?i)unsubscribe|opt-?out|preferences|subscri|newsletter|invite|verify|confirm|reset`)
)

// fromBase64 decodes s as standard or URL-safe base64, padded or not, or
// returns "".
func fromBase64(s string) string {
	// Kits sometimes mark the encoded part ("$bWFu…").
	s = strings.TrimLeft(s, "$~!*")
	if len(s) < 8 {
		return ""
	}
	for _, enc := range []*base64.Encoding{base64.StdEncoding, base64.URLEncoding, base64.RawStdEncoding, base64.RawURLEncoding} {
		if b, err := enc.DecodeString(s); err == nil {
			return string(b)
		}
	}
	return ""
}

// hostingRegistration reports whether the domain info looked up for a site on
// a hosting platform is really the platform's: a customer's subdomain, or a
// user's page on a link-in-bio service. The platform's own site keeps its age.
func hostingRegistration(resp Response) bool {
	host := strings.TrimPrefix(hostOf(resp.URL), "www.")
	if _, ok := constants.UserPageHosts[host]; ok {
		return !isBareURL(resp.URL)
	}
	return host != resp.Features.TLD.TLD && strings.HasSuffix(host, "."+resp.Features.TLD.TLD)
}

// kitHexToken is a long run of hex, as kits put in paths and session
// parameters; owaOnPHP is an Outlook Web Access path served by PHP, which
// real OWA (ASP.NET) never is.
var (
	kitHexToken = regexp.MustCompile(`(?i)[0-9a-f]{16,}`)
	owaOnPHP    = regexp.MustCompile(`(?i)/owa/.*\.php\b`)
)

// kitLoginPath reports whether a link is a sign-in page that looks planted:
// a sign-in path carrying a long hex token in the path or a session-style
// parameter, or OWA on PHP. Reset and confirmation links carry tokens too,
// under names like token= or code=, so those don't count.
func kitLoginPath(rawURL string) bool {
	u, err := neturl.Parse(rawURL)
	if err != nil {
		return false
	}
	if owaOnPHP.MatchString(u.Path) {
		return true
	}
	if !loginPathWords.MatchString(u.Path) {
		return false
	}
	if kitHexToken.MatchString(u.Path) {
		return true
	}
	for name, vals := range u.Query() {
		switch strings.ToLower(name) {
		case "token", "code", "key", "reset", "confirm", "verify", "signature", "sig", "state", "nonce":
			continue
		}
		for _, v := range vals {
			if kitHexToken.MatchString(v) {
				return true
			}
		}
	}
	return false
}
