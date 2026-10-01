package analyzer

import (
	"fmt"
	"math"
	"strings"
)

func GenerateResult(resp Response) Result {
	var neutralReasons []string
	var goodReasons []string
	var badReasons []string
	trustScore := 0
	riskScore := 0

	// --- 1. Popularity & Rank ---
	if resp.Features.Rank == 0 {
		if resp.Features.TLD.IsHostingPlatform {
			// Personal/project pages on known hosting platforms are inherently unranked — not a risk signal.
			neutralReasons = append(neutralReasons, fmt.Sprintf("Unranked subdomain on %s (normal for personal or project pages).", resp.Features.TLD.TLD))
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

	if resp.Features.TLD.IsTrusted {
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
	if resp.Analysis.SupportsHSTS {
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

	// --- 6. Domain History ---
	if resp.DomainInfo != nil {
		if resp.DomainInfo.AgeDays <= 30 {
			badReasons = append(badReasons, fmt.Sprintf("Newly created domain (%s old). High Risk.", resp.DomainInfo.AgeHuman))
			riskScore += 25
		} else if resp.DomainInfo.AgeDays <= 365 {
			badReasons = append(badReasons, fmt.Sprintf("Young domain (%s old). Use caution.", resp.DomainInfo.AgeHuman))
			riskScore += 15
		} else if resp.DomainInfo.AgeDays <= 1825 {
			neutralReasons = append(neutralReasons, fmt.Sprintf("Operational for %s.", resp.DomainInfo.AgeHuman))
			trustScore += 5
		} else {
			goodReasons = append(goodReasons, fmt.Sprintf("Long-standing domain history (%s).", resp.DomainInfo.AgeHuman))
			trustScore += 15
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
			riskScore += 20
		} else {
			badReasons = append(badReasons, fmt.Sprintf("Typosquatting detected: domain closely resembles '%s' (%d character difference).", resp.TyposquatResult.MatchedDomain, resp.TyposquatResult.Distance))
			riskScore += 40
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

	// --- 11. Page Content & Phishing Signals ---
	if resp.ContentData != nil {
		if resp.ContentData.HasLoginForm {
			// Check if domain is established
			isEstablished := resp.Features.Rank > 0 && resp.Features.Rank <= 100000
			isOld := resp.DomainInfo != nil && resp.DomainInfo.AgeDays > 365

			if !isEstablished && !isOld {
				badReasons = append(badReasons, "SUSPICIOUS: Login form detected on a new or unranked domain.")
				riskScore += 50
			} else {
				neutralReasons = append(neutralReasons, "Page contains a login form.")
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
				if form.ExternalAction {
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

	// --- Normalize / cap scores ---
	riskScore = clamp(riskScore)
	trustScore = clamp(trustScore)

	// Balanced formula: neutral point is 50.
	// trustScore pulls score up, riskScore pulls it down.
	// A site with no signals either way scores 50 → "Suspicious".
	finalScore := clamp(50 + int(math.Round(float64(trustScore-riskScore)*0.5)))

	var verdict string
	switch {
	// High risk, low trust
	case finalScore < 30:
		verdict = "Risky"
	// Moderate risk OR conflicting/insufficient signals
	case finalScore < 65:
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

func clamp(score int) int {
	return int(math.Max(0, math.Min(100, float64(score))))
}
