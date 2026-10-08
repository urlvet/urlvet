# Testing

## Running Tests

All tests live under `server/`. Run them from that directory.

```bash
cd server

# All tests
go test ./...

# Specific package
go test ./internal/analyzer/
go test ./internal/handler/
go test ./internal/service/checks/

# Verbose output
go test -v ./internal/analyzer/ -run TestGenerateResult

# With race detector
go test -race ./...
```

## What's Covered

### `internal/analyzer` — Scorer (`result_test.go`)

Table-driven tests for `GenerateResult()`. Every scoring branch is exercised:

| Test group | What it checks |
|---|---|
| `TestGenerateResult_Verdict` | 28 cases — rank, TLD, HSTS, IP, punycode, homoglyph, phishtank, redirects, domain age, DNSSEC, brand mismatch, login/payment/hidden-iframe signals |
| `TestGenerateResult_ScoreClamping` | Risk + final scores stay within 0–100 under max-risk pile-on |
| `TestGenerateResult_VerdictBoundaries` | Exact formula boundary: `finalScore = trust - risk×0.2` produces correct verdict bucket |

### `internal/analyzer` — Reported URLs (`reports_test.go`, `feed_fixes_test.go`)

Regression tests built from real scans: user reports and a run of 300 phishing URLs from a public feed. Each names the URL it came from.

| Area | What it checks |
|---|---|
| Domain age | Every band boundary (30, 90 days; 1, 3, 5 years); unpublished creation dates and stale cached ages count neither way |
| Downloads | Program on an upload host capped at Suspicious; from an unknown site adds risk; from the publisher is neutral |
| Short links | An unresolved short link is capped at Suspicious |
| Hosting | Provider endpoint takes the owner's rank, uploads under it don't; combo-squats on hosting weigh more |
| Content | Script redirect to a raw IP, Cloudflare interstitial, 451 takedown, certificate for another host, WordPress system folders |
| Brands | Bare names on login pages and hosting subdomains, news headlines on well-known sites, user-content hosts keep the check |
| Threat lists | Exact and host-level listings, host-level ignored on shared services, Safe Browsing and Web Risk counted once |

### `internal/handler` — HTTP smoke tests (`handler_test.go`)

Uses `net/http/httptest` against the full router (no mocks). Covers:

| Handler | Cases |
|---|---|
| `GET /health`, `GET /api/v1/health` | 200 + `status: ok` |
| `GET /` | 200 + service metadata |
| `GET /api/v1/analyze` | Missing URL → 400; empty-host URL → 400; URL > 2048 chars → 400 |
| `GET /api/v1/length`, `/depth` | Missing → 400; invalid → 400; valid → 200 |
| `GET /api/v1/punycode` | Missing → 400; invalid → 400; valid → 200 |
| `GET /api/v1/url-shortener` | Missing → 400; invalid → 400; valid → 200 |
| `GET /api/v1/trusted-tld`, `/risky-tld` | Missing → 400; valid → 200 |
| `GET /api/v1/ip/check` | Missing → 400; valid → 200 |
| `GET /metrics` | 200 + non-empty `Content-Type` |

### `internal/service/checks` — URL utilities and page signals (`checks_test.go`, `content_signals_test.go`, `feed_fixes_test.go`)

URL normalization (including `#fragment`), length and depth limits, customer-site domains (`x.weebly.com`), script-redirect detection, provider block pages, short-link resolution from interstitial pages, executable file types, login intent in forms, and brand matching: new sectors, bare-word false positives ("Portal Ayuntamiento Santander", "Amazon Rainforest Tours"), lookalike letters, and a site's own address in its title.

### `internal/service/threatfeeds` — Threat lists and Safe Browsing (`localfeeds_test.go`, `safebrowsing_test.go`)

Feed URL normalization and matching, PhishTank dump and URLhaus CSV parsing, Safe Browsing URL canonicalization against Google's published examples, host/path expressions, and protobuf response decoding.

### `internal/service/typosquat` (`typosquat_test.go`)

Edit-distance limits for short names, and combo-squatting against curated brands only (`takeoffservicesllc.com` doesn't match "service").

### `internal/service/domaininfo` (`lookup_test.go`)

Unknown creation dates serialize as `age_known: false` and `age_days: null`.

### `internal/service/rank` — Rank loader (`load_test.go`, `lookup_test.go`)

Rank file loading and lookup logic.

## URL Validation Behaviour

`IsValidURL` is permissive by design: bare hostnames (`example.com`) are accepted by prepending `https://`. Only truly malformed inputs (empty host, `http://`) are rejected. Handler tests reflect this. The `#fragment` is kept: redirector links carry their payload there, and folding it into the path breaks the link.

## Coverage Report

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Adding Tests

- Scorer tests: `server/internal/analyzer/result_test.go` — extend `TestGenerateResult_Verdict` table.
- Handler tests: `server/internal/handler/handler_test.go` — add rows to existing test tables or new `Test*` functions.
- Unit tests: place `_test.go` alongside the file under test.
