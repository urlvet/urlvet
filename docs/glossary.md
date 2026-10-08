# Glossary

Domain-specific terms, acronyms, and url.vet-internal concepts used throughout the codebase and documentation.

---

## A

**Analyzer**
The url.vet component that orchestrates all checks for a single URL. It launches 21 goroutines via `sync.WaitGroup`, collects their outputs, and feeds results to the scorer. See `server/internal/analyzer/`.

**AGPL-3.0**
GNU Affero General Public License v3. The open-source license url.vet is released under. Key clause: any modified version run over a network must make its source code available to users.

---

## C

**Cache hit / Cache miss**
When a URL is submitted, Valkey is checked first. A *hit* means a full result is already stored and is returned immediately (sub-millisecond). A *miss* means no cached result exists and the full analysis pipeline runs.

**chromedp**
A Go library that controls a headless Chrome browser over the Chrome DevTools Protocol (CDP) WebSocket. url.vet uses it to take page screenshots and for content fetching. Container: `urlvet-chrome` on port `:9222`.

**Combo-squatting**
A form of domain abuse where a brand name is combined with extra words to make the domain look legitimate, e.g. `paypal-login.com` or `apple-support.net`. url.vet looks for 514 curated brand names inside the domain (everyday words like "office" or "signal" are left out). On a hosting platform, where anyone can claim a name for free, a match weighs more.

**Content Analysis**
A category of 12 checks that fetch and parse the target page's HTML to look for phishing indicators — login forms, payment fields, hidden iframes, external form targets, brand-domain mismatches, script redirects, host or CDN warning pages, pages hidden in WordPress system folders, and program downloads.

**CT log (Certificate Transparency log)**
A public, append-only ledger of TLS certificates. Certificates without an embedded SCT (Signed Certificate Timestamp) may have been issued without public logging, which is a weak indicator of suspicious infrastructure.

---

## D

**DNS (Domain Name System)**
The internet's naming system. url.vet checks NS records (authoritative name servers), MX records (mail exchange), and IP resolution as part of domain infrastructure validation.

**DNSSEC (DNS Security Extensions)**
Cryptographic signatures on DNS responses that allow resolvers to verify authenticity. Absence of DNSSEC is a mild risk signal; many legitimate domains also lack it.

**Domain age**
How long ago a domain was first registered, retrieved via RDAP/WHOIS. Domains under 30 and under 90 days old count as risk, because attackers spin up fresh domains for each campaign; 3–12 months is neutral; over 1, 3 and 5 years adds increasing trust. Some registries (.de, .eu) don't publish a creation date, and then age counts neither way.

**Domain Intelligence**
A category of 6 checks covering domain rank, TLD classification, domain age, DNSSEC, Shannon entropy, and typosquatting detection.

---

## E

**eTLD (effective Top-Level Domain)**
The public suffix of a domain as defined by the [Public Suffix List](https://publicsuffix.org/). For `mail.google.co.uk`, the eTLD is `co.uk` and the eTLD+1 (registrable domain) is `google.co.uk`. url.vet uses `golang.org/x/net/publicsuffix` for this.

**eTLD+1**
The registrable portion of a domain — the eTLD plus one additional label. Used to determine the true owner of a domain and to compare redirect chains for cross-domain hops.

**Entropy (Shannon entropy)**
A measure of randomness in a string, calculated as `−Σ p(x) · log₂ p(x)`. High entropy in a domain name (e.g. `xq2k9f.net`) suggests it was algorithmically generated (DGA — Domain Generation Algorithm), which is common in malware C2 infrastructure.

---

## F

**False positive**
A legitimate URL that url.vet incorrectly scores as Suspicious or Risky. Common sources: newly registered legitimate domains, domains using IP-based hosting, or sites with unusual TLDs.

**Final URL**
The URL the browser lands on after following all redirects. url.vet records the full redirect chain and compares the final URL's domain to the original to detect cross-domain redirects.

---

## G

**Google Safe Browsing**
Google's lists of phishing, malware and unwanted-software sites. url.vet uses its v5 `hashes:search` method when `SAFE_BROWSING_API_KEY` is set: it sends only 4-byte prefixes of SHA-256 hashes of the URL's host/path combinations, never the URL, and matches the full hashes Google returns locally. Free, but for non-commercial use only; *Web Risk* is the commercial equivalent. Google's terms require "Advisory provided by Google" wherever its verdict is shown.

**Goroutine**
A lightweight concurrent execution unit in Go, cheaper than OS threads. url.vet launches one goroutine per analyzer task (18 total) via `sync.WaitGroup`, letting all checks run in parallel.

---

## H

**Homograph attack**
A phishing technique that uses visually similar Unicode characters to spoof domain names. Example: `аpple.com` where `а` is Cyrillic U+0430, not Latin `a`. Also called an IDN homograph attack.

**Hosting platform**
A service that gives each customer their own subdomain: GitHub Pages (`*.github.io`), Vercel, Netlify, Cloudflare Pages, and site builders like Weebly or GoDaddy Sites. url.vet judges such a subdomain as its own site: it doesn't inherit the provider's rank or age, isn't penalized for having no DNS records of its own, and its customer-chosen name is what the typosquatting check compares. See `TrustedHostingPlatforms` and `CustomerSiteHosts` in `server/internal/constants/`.

**HSTS (HTTP Strict Transport Security)**
An HTTP response header (`Strict-Transport-Security`) that instructs browsers to always use HTTPS for a domain. Its absence on an HTTPS site is a weak risk signal.

**HTTP Combined Task** (`httpCombinedTask`)
An internal url.vet optimization. Instead of making separate HTTP requests for redirects, HSTS, and status code checks, a single HTTP request is shared across all three. Falls back to individual checks if the combined request fails.

---

## I

**IDN (Internationalized Domain Name)**
A domain name containing non-ASCII Unicode characters, encoded using Punycode for DNS transport (e.g. `münchen.de` → `xn--mnchen-3ya.de`). url.vet flags IDN domains as they are frequently used in homograph attacks.

**iframe (hidden)**
An HTML `<iframe>` element with zero dimensions or `display:none`. Hidden iframes are used for clickjacking, credential theft, or silently loading malicious content.

**IP address as hostname**
Using a raw IPv4 or IPv6 address directly as the URL host (e.g. `http://192.168.1.1/login`). Legitimate services almost never use bare IP addresses; phishers use them to avoid domain registration and takedowns.

---

## L

**LRU (Least Recently Used)**
The cache eviction policy used by Valkey in url.vet. When the cache reaches its memory limit (`maxmemory`), the least recently accessed keys are evicted first. Configured as `allkeys-lru`.

---

## M

**MX record**
A DNS record that specifies the mail server for a domain. Missing or invalid MX records combined with other signals can indicate a disposable or fraudulent domain.

**Middleware**
In the Gin HTTP framework context, middleware is a function that runs before or after route handlers. url.vet middleware: request logger, rate limiter, CORS, URL length validator, Bearer auth, and Prometheus metrics.

---

## N

**NS record**
A DNS record that delegates a domain to its authoritative name servers. Missing or non-resolving NS records indicate a domain that is parked, expired, or misconfigured.

**Normalize (URL normalization)**
The process of bringing a URL to a canonical form before analysis: inferring a scheme (`https://` if missing), lowercasing the hostname, and percent-encoding special characters. Ensures consistent cache keys and check behavior.

---

## P

**PhishTank**
A community-driven database of verified phishing URLs maintained by Cisco Talos. url.vet downloads its list of verified, still-online phishing URLs every 6 hours (hourly with an API key) and checks every scan against it in memory. Its live API is only called while that list isn't loaded. Free, including for commercial use.

**Pipeline**
The end-to-end flow of a URL through url.vet: normalize → cache check → parallel analyzers → score aggregation → cache store → response. Visualized in `assets/pipeline.png`.

**Punycode**
An ASCII-compatible encoding for Unicode domain labels, defined in RFC 3492. `xn--` is the ACE prefix that identifies a Punycode-encoded label. url.vet detects domains with Punycode labels as a potential homograph risk.

---

## R

**RDAP (Registration Data Access Protocol)**
The modern replacement for WHOIS (RFC 7483). Returns structured JSON for domain registration data including creation date, registrar, and name servers. url.vet uses RDAP when available, falls back to WHOIS.

**Reason**
A labeled explanation string emitted by each check — classified as `good`, `bad`, or `neutral`. Every signal that fires produces a reason, making the final verdict fully explainable. Stored in `GoodReasons`, `BadReasons`, `NeutralReasons` in the response.

**Redirect chain**
The ordered list of URLs traversed when following HTTP redirects. url.vet records each hop and flags chains that cross domain boundaries or are excessively long.

**Risk score**
A 0–100 accumulator for negative signals. Each bad check adds weighted points. Combined with the trust score in the final formula. Clamped to 100 before scoring.

---

## S

**SCT (Signed Certificate Timestamp)**
A cryptographic token embedded in a TLS certificate (or delivered via TLS extension/OCSP stapling) proving it was logged to a CT log before issuance. Required by major browsers for public certificates issued after 2018.

**Shannon entropy** → see *Entropy*

**Short link**
A link on a URL shortener, as opposed to the shortener's own homepage. url.vet follows it (HTTP redirects, meta refresh, script redirects, and interstitial pages that name the destination), through any further short links, up to 10 hops, and scans the page it leads to. The response's `short_link` field records the chain. If the destination can't be found, the short link itself is scanned and the verdict is capped at Suspicious.

**SLD (Second-Level Domain)**
The label directly to the left of the TLD. In `mail.google.com`, the SLD is `google`. url.vet's typosquatting check compares the SLD against known brand names.

**SSRF (Server-Side Request Forgery)**
An attack where a server is tricked into making HTTP requests to internal or private network addresses on behalf of an attacker. url.vet blocks SSRF by rejecting URLs that resolve to private IP ranges (RFC 1918, loopback, link-local) before fetching.

**Signal**
A single, binary or graded indicator produced by one check — e.g. "domain registered 3 days ago" or "HSTS present". url.vet produces 39 signals across 7 categories. Each signal maps to one or more reasons.

**Suspicious** → see *Verdict*

---

## T

**Task**
The internal unit of work in the analyzer. Each task implements the `Task` interface (`Name() string`, `Run(*Input, *Output) error`) and corresponds to one analyzer goroutine. 18 tasks are registered in `analyze.go`.

**Threat Intelligence**
A category of 4 checks against lists of known-bad URLs: the locally held threat lists (exact link, or another page on the same single-owner host), Google Safe Browsing or Web Risk, and PhishTank's live API as a fallback.

**Threat list**
A list of reported phishing or malware URLs that url.vet downloads on a schedule, caches under `FEEDS_DIR`, and matches in memory: PhishTank (on by default), URLhaus and OpenPhish (opt-in). Each has its own licence; see [configuration.md](configuration.md#threat-feeds-also-serverenv).

**TLD (Top-Level Domain)**
The rightmost label of a domain (`.com`, `.net`, `.gov`). url.vet classifies TLDs as trusted (`.gov`, `.edu`, `.mil`), risky (commonly abused ccTLDs and gTLDs), or neutral.

**TLS (Transport Layer Security)**
The cryptographic protocol that secures HTTPS connections. url.vet performs a full TLS handshake to inspect the certificate chain — validity, expiry, issuer, and CT log status.

**Trust score**
A 0–100 accumulator for positive signals. Each good check adds weighted points. Combined with the risk score in the final formula. Clamped to 100 before scoring.

**Typosquatting**
Registering domains that are intentional misspellings of well-known brands to capture mistyped traffic or trick users — e.g. `googie.com`, `paypa1.com`. url.vet compares the domain's name against the 5,000 most-visited sites by edit distance (one edit for short names, up to two for longer ones).

---

## U

**URL shortener**
A service that maps a short URL (e.g. `bit.ly/abc`) to a longer destination URL. Shorteners are used by phishers to obscure the true target. url.vet keeps a list of 2,671 shortener domains and follows short links to their destination; see *Short link*.

**URLhaus**
abuse.ch's list of URLs distributing malware. url.vet downloads the recent export hourly when `URLHAUS_AUTH_KEY` is set and keeps the URLs still online. The key is free; use is free for not-for-profit purposes, while commercial use needs a Spamhaus subscription.

---

## V

**Valkey**
A Redis-compatible, open-source key-value store. url.vet uses it as an LRU cache for full analysis results (24 h TTL), content analysis (3 h TTL), and Google and PhishTank API lookups (1 h and 3 h TTL). Downloaded threat lists live in memory and under `FEEDS_DIR`, not in Valkey. Data is persisted to a Docker volume.

**Verdict**
The human-readable classification assigned to a URL based on its final score:

| Score | Verdict |
|---|---|
| ≥ 65 | **Safe** |
| 30 – 64 | **Suspicious** |
| < 30 | **Risky** |

---

## W

**Web Risk**
Google's commercial-use counterpart to Safe Browsing (`WEBRISK_API_KEY`). It needs a billing account (100,000 lookups a month free) and, unlike url.vet's Safe Browsing use, receives the URL itself.

**WHOIS**
A text-based protocol (RFC 3912) for querying domain registration data. Returns domain age, registrar, name servers, and contact information. Being replaced by RDAP. url.vet uses WHOIS as a fallback when RDAP data is unavailable.

---

## Scoring Reference

```
finalScore = clamp(50 + (trustScore − riskScore) × 0.5, 0, 100)
```

| Symbol | Meaning |
|---|---|
| `50` | Neutral baseline — default for a URL with no signals |
| `trustScore` | Sum of positive signal weights, clamped 0–100 |
| `riskScore` | Sum of negative signal weights, clamped 0–100 |
| `× 0.5` | Dampening factor so neither side dominates alone |
| `clamp(…, 0, 100)` | Final score is always in [0, 100] |

A program download from a host anyone can upload to, and a short link whose destination couldn't be found, cap the final score at 64 (Suspicious) instead of adding risk.
