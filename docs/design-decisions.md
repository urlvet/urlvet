# Design Decisions

Key architectural choices in url.vet and the reasoning behind them.

---

## Go for the backend

**Decision:** The analysis server is written in Go.

**Why:** The core bottleneck is I/O — DNS lookups, HTTP fetches, TLS handshakes, WHOIS queries. Go's goroutines make it trivial to run all 21 checks concurrently with a `sync.WaitGroup`. The result is 300–700ms median scan time with no async/await complexity and clean per-goroutine panic recovery. A Python or Node implementation would need explicit async machinery for the same concurrency and would carry more runtime overhead per request.

---

## No machine learning (intentional)

**Decision:** All detection is rule-based heuristics. No ML model.

**Why:** Explainability is a first-class goal. Every verdict can be traced to specific, named signals with human-readable reasons. An ML model would give accurate verdicts but opaque ones — users and integrators couldn't understand or trust why a URL was flagged. The tradeoff is accepted: url.vet will miss sophisticated phishing that mimics legitimate site structure, but it will never produce a verdict nobody can explain.

This decision may be revisited. If added, ML would be a separate signal alongside existing heuristics, not a replacement.

---

## 21 goroutines, all concurrent

**Decision:** All checks launch simultaneously via `sync.WaitGroup`. No prioritization or staged execution.

**Why:** Simplicity. The longest checks (WHOIS ~1–2s, screenshot ~10–30s) dominate latency regardless of ordering. Running fast structural checks first and gating slow network checks on their results would add coordination complexity for minimal gain — the total wall time is still bounded by the slowest check.

Known tradeoff: under high load, all goroutines (including slow screenshot/WHOIS) consume resources for every request. Making screenshot opt-in is on the roadmap.

---

## Valkey (not PostgreSQL) for caching

**Decision:** Full analysis results are cached in Valkey (Redis-compatible), not a relational database.

**Why:** The primary access pattern is `url → result` — a pure key-value lookup. Valkey handles this in sub-millisecond with built-in LRU eviction and TTL. PostgreSQL would add connection overhead, schema complexity, and index maintenance for a use case that maps directly to a key-value store's strengths.

Tradeoff: no persistent scan history, no queries across results. Scan history (for analytics, training data, abuse detection) requires PostgreSQL and is on the roadmap.

---

## Single HTTP request for redirect + HSTS + status code

**Decision:** `httpCombinedTask` makes one HTTP request shared across three checks (redirects, HSTS header, status code) instead of three separate requests.

**Why:** These three checks all need the same HTTP response. Before the optimization, `CheckRedirects` and `CheckHSTS` each made independent requests to the same URL. The combined task issues a single HEAD request (falling back to GET if needed), extracts all three results, and eliminates two redundant network round-trips.

HSTS is read from the scanned site's own HTTPS response, not from the last hop: a github.com download ends on a storage host that doesn't send the header.

---

## Threat lists downloaded, not queried per scan

**Decision:** PhishTank's list (and URLhaus and OpenPhish when enabled) is downloaded on a schedule and matched in memory. PhishTank's live API is only a fallback while its list isn't loaded.

**Why:** Without an API key, PhishTank rate-limited about three out of four scans, so the check mostly didn't happen and every result carried a "couldn't check PhishTank" note. A local copy answers every scan instantly, sends no scanned URL anywhere, and costs about 10 MB. The cost is freshness (up to 6 hours old) and missing PhishTank's unverified reports.

Only feeds whose terms allow url.vet's use are on by default: PhishTank allows commercial use; OpenPhish's free feed allows personal research only. See [configuration.md](configuration.md#threat-feeds-also-serverenv).

---

## Google Safe Browsing by hash prefix

**Decision:** Safe Browsing uses v5 `hashes:search`, which sends Google 4-byte hash prefixes, rather than `urls:search`, which sends the URL.

**Why:** The privacy page promises the scanned link only goes where it says. With hash prefixes, Google can't tell which link was checked; the full-hash match is made on our server. It needs more code (URL canonicalization and protobuf decoding, since v5 only answers in protobuf), but no third party learns what users scan.

---

## Short links: scan the destination

**Decision:** A short link is followed, through any further short links, and the destination is scanned in its place. If the destination can't be found, the short link itself is scanned and capped at Suspicious.

**Why:** A short link says nothing about itself, and the shortener's own reputation (bit.ly is a top-ranked site) would otherwise make any short link look safe. Some shorteners show an interstitial page instead of redirecting, so the resolver also reads meta refresh, script redirects and `data-url`-style attributes.

---

## Redirects: judge where the link lands

**Decision:** When a link sends visitors on to another site, by HTTP redirect or a page that forwards straight away, the destination is scanned and the result describes it. Trust (rank, age, HSTS, a known brand) comes only from the destination. Risk comes from wherever it's found: the destination, the link as given (a lookalike or brand name, login words, a victim's address, an IP, a fresh or high-risk domain, a threat-list entry) and the hops in between (IP addresses, listed hosts), which are checked from local data only. A link a list names is Risky wherever it lands. Moves within one site or brand (http → https, zoom.us → zoom.com, a Blogspot page to research.google) aren't jumps. One exception: an unknown site (unranked, not on a restricted registry) landing on a well-known one isn't followed. That's how cloaking works: kits send scanners to google.com or the brand's real login and victims to the phishing page. The link is judged itself, and the jump counts against it. An unknown site sending visitors on to another unknown one is followed, and the bounce adds risk.

**Why:** The visitor types their password on the page they land on, so that's whose reputation matters; google.com's rank says nothing about where `google.com/url?q=…` sends them. But a link that lies about itself is still a lie, and a clean destination shouldn't excuse it.

---

## Caps instead of risk points

**Decision:** Two findings cap the verdict at Suspicious instead of adding risk: a program downloaded from a host anyone can upload to, and a short link that couldn't be followed.

**Why:** In both, the host's reputation is high but says nothing about the content: an installer under `github.com/<anyone>/releases` is whoever uploaded it. github.com's trust alone saturates the score, so no reasonable risk weight would move it off Safe, while a cap states what's known: it can't be called safe.

---

## Hosting customers judged on their own

**Decision:** A customer's subdomain on a hosting service or site builder (`*.github.io`, `*.vercel.app`, `*.weebly.com`, `*.godaddysites.com`) is treated as its own site: no rank or age from the provider, no penalty for missing DNS records, and its own name is what the typosquatting check compares.

**Why:** Most of these providers aren't on the Public Suffix List, so by default `xfinitylogin.weebly.com` inherited Weebly's rank (#345) and 20-year age and scored Safe. The reverse holds too: a provider's reputation can't be lent to its customers, so `storage.googleapis.com/<bucket>` doesn't get Google's rank, while the bare endpoint does.

---

## Brand names: strict keywords and loose names

**Decision:** Each brand has strict title keywords that count anywhere, and bare names that only count on a page asking for a login or payment, or on a hosting subdomain. Well-known sites (top 100,000) may mention brands, except on hosts anyone can publish on.

**Why:** A bare name is both the strongest phishing tell ("Facebook" over a login form) and the commonest false positive ("Santander" is also a city, "Amazon" a rainforest, "Vodafone Shop Berlin" a real shop). Splitting them keeps the check sharp on pages that ask for something and quiet on everything else.

---

## SvelteKit for the frontend

**Decision:** The web UI is a SvelteKit application served as a static build in production.

**Why:** SvelteKit's server-side rendering lets API calls to the Go backend happen server-side, avoiding CORS issues in the browser. The compiled output is a static Nginx container with no Node.js runtime in production. Svelte's compiler-first approach produces minimal bundle sizes compared to React or Vue equivalents.

---

## Gin as the HTTP framework

**Decision:** The REST API uses the Gin framework.

**Why:** Gin provides routing, middleware chaining, and request binding with minimal overhead. It is the most widely used Go HTTP framework with good community support. The alternative (stdlib `net/http`) would require writing routing and middleware manually; heavier frameworks (Echo, Fiber) offer no meaningful advantage for url.vet's API surface.

---

## Structured logging via slog

**Decision:** All logging uses Go's stdlib `log/slog` wrapped in a centralized logger package.

**Why:** `slog` (Go 1.21+) provides structured, leveled logging with no external dependency. The centralized wrapper (`internal/logger`) lets the format be changed in one place — colored text in DEV, JSON in production — without touching any call site. Full URLs are kept out of logs (scheme+host only) to prevent token leakage in log files.

---

## AGPL-3.0 license

**Decision:** url.vet is licensed under AGPL-3.0 with a separate commercial license option.

**Why:** AGPL requires any modified version run over a network to make its source code available. This prevents organizations from running a proprietary fork as a SaaS without contributing back. The commercial license option lets organizations that cannot comply with AGPL (closed-source SaaS) pay for an exception.

---

## chromedp over a REST screenshot API

**Decision:** Screenshots are taken using chromedp (a Go-native Chrome DevTools Protocol client) against a shared headless Chrome container, not an external screenshot API.

**Why:** No external dependency, no API key, no per-screenshot cost, no data leaving the deployment. The shared browser allocator (`screenshot.GetService()`) reuses a single Chrome instance across requests, avoiding the overhead of launching a new browser per scan. Tradeoff: the Chrome container adds ~300MB memory baseline.
