# Architecture

## Overview

url.vet is a real-time URL analysis engine. When a URL is submitted — via the web UI, Chrome extension, or REST API — the Go backend runs **21 concurrent analyzers** across 7 signal categories, aggregates a trust/risk score, assigns a verdict, and returns a fully explainable report. Results are cached in Valkey so repeat lookups are instant.

---

## System Architecture

![url.vet Architecture](../assets/architecture.png)

Four containerized services run on a shared Docker bridge network (`urlvet-net`). The Go backend is the **only** service that makes outbound calls — the frontend, Chrome, and cache are strictly internal.

| Service | Container | Role | Port |
|---|---|---|---|
| `urlvet-web` | SvelteKit UI | Renders the web interface; proxies API calls to the backend | `:3000` prod · `:5173` dev |
| `urlvet-backend` | Go REST API | Validates URLs, runs analyzers, aggregates scores, manages cache | `:8080` |
| `urlvet-chrome` | Headless Chrome | Takes page screenshots and serves content via WebSocket (chromedp) | `:9222` |
| `urlvet-valkey` | Valkey (Redis-compatible) | LRU result cache, volume-persisted across restarts | `:6379` |

**External services** (reached only by the backend over HTTPS/TCP):
- **Threat lists** — PhishTank's list (default), URLhaus and OpenPhish (opt-in) are downloaded on a schedule and matched in memory; no per-scan call. PhishTank's API is called only while its list isn't loaded
- **Google Safe Browsing / Web Risk** (opt-in) — Safe Browsing receives only 4-byte hash prefixes of the URL; Web Risk receives the URL
- **DNS resolvers** — NS, MX, IP resolution checks
- **RDAP / WHOIS servers** — domain age and registration data

Clients (browser, Chrome extension, API consumers) communicate directly with the Go backend. The SvelteKit frontend forwards all `/api/v1/` requests server-side to avoid CORS complexity.

---

## Request Lifecycle

![url.vet Analyzer Pipeline](../assets/pipeline.png)

```
Client
  │
  │  GET /api/v1/analyze?url=...
  ▼
Go Backend
  ├─ 1. Validate & normalize URL (add scheme if missing, keep the #fragment, reject private IPs)
  ├─ 2. Check Valkey cache
  │      └─ HIT  → return full cached result immediately (sub-millisecond)
  │      └─ MISS → continue
  ├─ 3. Short link? Follow it (HTTP redirects, meta refresh, script, data-url
  │      interstitials), through any further short links, up to 10 hops.
  │      └─ Resolved   → scan the destination instead; the result carries `short_link`
  │      └─ Unresolved → scan the short link itself, capped at Suspicious
  ├─ 4. Launch 21 goroutines via sync.WaitGroup
  │      ├─ Each task runs independently; panics are recovered per-task
  │      ├─ Tasks share a read-only Input struct and write to a mutex-guarded Output
  │      └─ All 21 complete (or timeout) before proceeding
  ├─ 5. Apply exceptions (known false positives: well-known sites, same-brand
  │      redirects, hosting platforms) → aggregate scores → assign verdict
  ├─ 5b. Leads to another site (HTTP redirect, or a page that forwards by
  │      script or meta refresh), from a ranked site or to an unknown one?
  │      Scan the destination (up to 2 jumps, 12 s each) and
  │      report that: trust from the destination only, plus the risks of the
  │      link as given and of any hop on the way (`redirected_from`, `origin`)
  ├─ 6. Store result in Valkey (24 h TTL) unless a check that could change it failed
  └─ 7. Return: trust score · verdict · per-signal reasons · redirect chain ·
              screenshot · per-task timings
```

---

## Detection Engine

21 goroutines run across **7 signal categories**, producing **39 individual signals**. Every check emits a labeled reason string — good, bad, or neutral — so the final score is always fully explainable.

### Scoring Formula

```
finalScore = clamp(50 + (trustScore − riskScore) × 0.5, 0, 100)
```

- **50** is the neutral baseline — an unknown URL with no signals scores exactly 50
- Trust signals pull the score up; risk signals pull it down, each weighted at 0.5× so neither dominates
- Both `trustScore` and `riskScore` are individually clamped to 0–100 before the formula runs

| Range | Verdict |
|---|---|
| ≥ 65 | Safe |
| 30 – 64 | Suspicious |
| < 30 | Risky |

Two signals cap the verdict at **Suspicious** (score 64) instead of adding risk, because a well-known host's reputation would otherwise carry them to Safe: a program download from a host anyone can upload to (GitHub releases, Google Drive, Dropbox…), and a short link whose destination couldn't be found.

### Signal Categories

**URL Signals** — 8 checks, purely structural, no network call

1. Raw IP address as hostname
2. Punycode / IDN encoding (lookalike domain spoofing)
3. URL shortener — followed to its destination, which is scanned instead (see step 3 above)
4. Excessive URL length (over 150 characters; not counted on well-known sites)
5. Excessive URL path depth (over 6 path segments; not counted on well-known sites)
6. Phishing keywords in URL path (`login`, `verify`, `secure`, `update` …)
7. Excessive subdomain count
8. Non-ASCII Unicode characters in hostname (IDN homograph attack)

**HTTP / Network** — 4 checks, single HTTP request via `httpCombinedTask`

9. Redirect chain hop count (only counted when the chain leaves the site)
10. Cross-domain redirect (not counted for a well-known site moving to a well-known or same-brand address, e.g. github.com → githubusercontent.com)
11. HSTS support, read from the scanned site's own HTTPS response rather than the last hop
12. HTTP status code

**DNS** — 3 checks

13. NS record validity
14. MX record validity
15. IP resolution

**TLS / SSL** — 2 checks, single TLS handshake to the link's own host

16. TLS presence, and a certificate issued for a different hostname (browsers show a full-page warning)
17. Certificate chain — validity, expiry, issuer, CT log status, known-bad fingerprints

**Domain Intelligence** — 6 checks

18. Domain rank (position in top-1M global popularity list). Unranked only counts against a site that's new or of unknown age
19. TLD trust / risk / ICANN status, and hosting platforms: a customer's subdomain on a hosting service or site builder (`*.github.io`, `*.vercel.app`, `*.weebly.com`, `*.godaddysites.com`…) is judged as its own site, without the provider's rank or age
20. Domain age via RDAP/WHOIS: under 30 days and under 90 days count as risk, 3–12 months is neutral, over 1, 3 and 5 years add increasing trust. Registries that don't publish a creation date (.de, .eu) give no age either way
21. DNSSEC (cryptographic DNS response integrity)
22. Shannon entropy score (flags algorithmically generated domains)
23. Typosquatting (edit distance to the 5,000 most-visited sites) and combo-squatting (514 curated brand names inside the domain; weighs more on free hosting subdomains)

**Content Analysis** — 12 checks, one HTTP GET to fetch the page (non-HTML responses such as downloads aren't parsed)

24. Login form (a password field, or an email/username field with a sign-in, next or verify button) on an unranked or new domain
25. Payment form (credit card, CVV fields)
26. Personal information form
27. Hidden `<iframe>` (credential theft / clickjacking vector)
28. Tracking pixels (1×1 hidden images)
29. Brand claimed in the page title vs. the hosting domain — 514 brands across banking, payments, crypto, shopping, delivery, telecoms, government and more. Lookalike letters (Ç, Cyrillic а) are folded first, and the page's own address is ignored. Bare brand names only count on a page asking for a login or payment, or on a hosting subdomain; well-known sites may mention brands
30. Form submitting to an external domain
31. Password field over unencrypted HTTP
32. Script or meta-refresh redirect, to a raw IP address or (on a near-empty page) another site
33. Host or CDN warning page in place of the site (Cloudflare's "Suspected Phishing" interstitial, a 451 takedown)
34. Page served from a WordPress system folder (`wp-content`, `wp-includes`, `wp-admin`), where phishing kits hide on hacked sites
35. Program download (`.exe`, `.msi`, `.apk`, `.dmg`…): caps the verdict at Suspicious on hosts anyone can upload to, adds risk on unknown sites

**Threat Intelligence** — 4 checks

36. On a locally held threat list (PhishTank by default; URLhaus, OpenPhish opt-in) — exact link
37. Other pages on the same host are listed — only where one owner controls the host (a hosting subdomain, or a domain under a year old)
38. Google Safe Browsing / Web Risk (opt-in)
39. PhishTank API, only while its list isn't loaded

---

## Code Layout

```
server/
├── cmd/urlvet/           entry point — init, router setup, graceful shutdown
├── internal/
│   ├── analyzer/
│   │   ├── analyze.go      task registration, cache integration
│   │   ├── runner.go       goroutine runner with panic recovery
│   │   ├── tasks.go        21 task implementations
│   │   ├── exceptions.go   known false positives cleared before scoring
│   │   └── result.go       score aggregation, verdict assignment
│   ├── handler/
│   │   ├── router.go       Gin router, middleware wiring
│   │   ├── analyze.go      /api/v1/analyze handler
│   │   └── middleware/     rate limiter, auth, Prometheus, request logger
│   ├── service/
│   │   ├── checks/         individual analyzer implementations, short-link resolver
│   │   ├── screenshot/     headless Chrome integration (chromedp)
│   │   ├── cache/          Valkey client wrapper
│   │   ├── threatfeeds/    local threat lists, PhishTank, Safe Browsing, Web Risk
│   │   └── typosquat/      brand similarity engine
│   ├── constants/          brands, URL shorteners, TLD lists, hosting and upload hosts, TTLs
│   ├── logger/             centralized slog-based logger (colors in DEV, JSON in prod)
│   └── admintoken/         admin JWT issuance and verification
web/
├── website/                SvelteKit UI
└── chrome-extension/       browser extension
docker/
├── dev/                    dev Compose (hot reload, exposed ports)
└── prod/                   prod Compose (optimized builds, restart policies)
docs/                       API reference, setup guide, architecture, security
```

---

## Deployment

Two fully separated Docker Compose stacks share the same image definitions but differ in configuration:

| | Dev | Prod |
|---|---|---|
| Backend | Air hot-reload, source mounted as volume | Compiled binary in distroless image |
| Frontend | Vite dev server `:5173` | Static build served by Nginx `:3000` |
| Chrome | Same `chromedp/headless-shell` image | Same |
| Valkey | Port exposed (`:6379`) for local inspection | Port not exposed; internal only |
| ENV | `ENV=DEV` — colored logs, debug endpoints | `ENV=PROD` — JSON logs, info level only |

Start everything with one command:

```bash
make start       # production stack
make dev         # development stack (hot reload)
```

See [docs/setup.md](setup.md) for full setup instructions including `.env` configuration.
