# API Reference

url.vet is built to be self-hosted: run your own copy with Docker (see [Setup](setup.md) and [Deployment](deployment.md)) and you get the same API on your own server, with your own limits and no one else's traffic.

To try it before setting anything up, the public instance at `https://api.url.vet` serves the same endpoints. It's rate-limited (20 requests a minute per IP) and meant for trying things out, not for building on.

## Interactive Docs (Swagger UI)

Every running server serves its own OpenAPI spec and a "Try it out" form:

| | Your server | Public instance |
|---|---|---|
| Swagger UI | `http://localhost:8080/swagger/index.html` | [api.url.vet/swagger/index.html](https://api.url.vet/swagger/index.html) |
| Spec (JSON) | `http://localhost:8080/swagger/doc.json` | |

The spec is also committed as `server/internal/docs/swagger.yaml`. To regenerate it after editing handler annotations:

```bash
cd server
swag init -g cmd/urlvet/urlvet.go -o internal/docs
```

Install the CLI once with: `go install github.com/swaggo/swag/cmd/swag@v1.16.4`

---

## Analyze a URL

```
GET /api/v1/analyze?url=<url>
```

Runs every check in parallel and returns one scored report. This is the endpoint the website uses.

```bash
curl "http://localhost:8080/api/v1/analyze?url=https://example.com"
```

**Parameters**

| Name | In | Required | Notes |
|---|---|---|---|
| `url` | query | yes | Full URL or bare domain (`example.com`). Max 2048 characters. URL-encode it if it has its own query string. |

**Responses**

| Status | When | Body |
|---|---|---|
| `200` | The scan ran | The report below |
| `400` | `url` is missing, or isn't a valid URL | `{"error": "url query param is required"}` or `{"status": "ERROR", "error": "invalid url"}` |
| `422` | The URL is valid but has nothing to scan (no usable domain) | `{"status": "ERROR", "error": "could not analyze this URL"}` |
| `429` | Rate limit reached | `{"error": "too many requests", ...}` |

A scan takes 1 to 15 seconds. Checks still running at 15 seconds are dropped, and the report says which ones (see `incomplete` below).

### Reading the result

Most clients only need `result`:

| Field | Meaning |
|---|---|
| `result.verdict` | `Safe`, `Suspicious` or `Risky` |
| `result.final_score` | 0 to 100, higher is safer. Below 30 is `Risky`, below 65 is `Suspicious`, 65 and up is `Safe` |
| `result.risk_score` / `result.trust_score` | The two sides of the score, each 0 to 100. `final_score = 50 + (trust − risk) / 2`, rounded and clamped |
| `result.reasons.bad_reasons` | Red flags, in plain English (`null` when there are none) |
| `result.reasons.good_reasons` | Green flags |
| `result.reasons.neutral_reasons` | Facts worth knowing that don't move the score |

The rest of the report is the raw evidence behind those reasons:

| Field | What it holds |
|---|---|
| `url`, `domain` | The normalized URL that was scanned and its registrable domain |
| `features` | Popularity rank, TLD facts, and URL structure (shortener, raw IP, punycode, length, depth, subdomains, keywords, lookalike characters) |
| `infrastructure` | Resolved IPs, nameservers, mail servers |
| `domain_info` | Registration data from RDAP or WHOIS: registrar, created/expiry dates, age, DNSSEC. `null` if the lookup failed |
| `analysis` | Redirect chain (`chain`, `final_url`, `has_domain_jump`), HTTP status, HSTS |
| `ssl_info`, `tls_info` | Certificate issuer, validity, age, Certificate Transparency, hostname match |
| `content_data` | What the page contains: forms (login, payment, personal data), where they submit, hidden iframes, trackers, and whether the page claims to be a brand that doesn't own this domain (`brand_check`). `null` if the page couldn't be fetched |
| `domain_randomness` | How machine-generated the domain name looks |
| `typosquat_result` | Whether the domain imitates a well-known one (`matched_domain`, `distance`) |
| `phishing` | PhishTank lookup. `null` if the lookup didn't happen |
| `performance` | Total time and per-check timings |

### Incomplete scans

Some checks depend on other services (DNS, WHOIS/RDAP, the site itself, PhishTank) and can fail or time out.

| Field | Meaning |
|---|---|
| `incomplete` | `true` when a missing check could change the verdict. Treat the verdict with care and scan again later |
| `incomplete_checks` | Names of the checks that didn't finish, e.g. `["whois_lookup"]`. Omitted when everything ran. Can be non-empty while `incomplete` is `false`: a PhishTank rate limit is listed but doesn't make the result incomplete |
| `errors` | The underlying error messages, for debugging. `null` when there were none |

### Caching

Complete results are cached for 24 hours per normalized URL, so repeat scans are instant and identical. Incomplete results are never cached, so scanning again retries the failed checks.

### Example response

`https://example.com`, scanned on a local server. `domain_info.raw` and `performance.timings` are shortened here.

<details>
<summary>Show the full response</summary>

```json
{
  "url": "https://example.com",
  "domain": "example.com",
  "features": {
    "rank": 175,
    "tld": {
      "tld": "com",
      "is_trusted_tld": false,
      "is_risky_tld": false,
      "is_icann": true,
      "is_hosting_platform": false
    },
    "url": {
      "url_shortener": false,
      "uses_ip": false,
      "contains_punycode": false,
      "too_long": false,
      "too_deep": false,
      "has_homoglyph": false,
      "subdomain_count": 0,
      "keywords": {
        "has_keywords": false,
        "found": null,
        "categories": null
      }
    }
  },
  "infrastructure": {
    "ip_addresses": [
      "172.66.147.243",
      "104.20.23.154",
      "2606:4700:10::ac42:93f3",
      "2606:4700:10::6814:179a"
    ],
    "nameservers_valid": true,
    "ns_hosts": [
      "hera.ns.cloudflare.com."
    ],
    "mx_records_valid": false,
    "mx_hosts": [
      "."
    ]
  },
  "domain_info": {
    "domain": "EXAMPLE.COM",
    "registrar": "RESERVED-Internet Assigned Numbers Authority",
    "created": "1995-08-14T04:00:00Z",
    "updated": "2026-08-14T08:01:43Z",
    "expiry": "2027-08-13T04:00:00Z",
    "nameservers": [
      "ELLIOTT.NS.CLOUDFLARE.COM",
      "HERA.NS.CLOUDFLARE.COM"
    ],
    "status": [
      "client delete prohibited",
      "client transfer prohibited",
      "client update prohibited"
    ],
    "dnssec": true,
    "age_human": "31 years 2 months",
    "age_days": 11372,
    "raw": "…",
    "source": "RDAP"
  },
  "analysis": {
    "redirection_result": {
      "is_redirected": false,
      "chain_length": 1,
      "chain": [
        "https://example.com"
      ],
      "final_url": "https://example.com",
      "final_url_domain": "example.com",
      "has_domain_jump": false
    },
    "http_status": {
      "code": 200,
      "text": "OK",
      "success": true,
      "is_redirect": false
    },
    "is_hsts_supported": false
  },
  "ssl_info": {
    "Domain": "example.com",
    "HasTLS": true,
    "ChainValid": true,
    "Issuer": "Cloudflare TLS Issuing ECC CA 3",
    "NotBefore": "2026-09-26T22:49:11Z",
    "NotAfter": "2026-12-25T22:56:35Z",
    "AgeDays": 5,
    "Fingerprint": "85CA6AB068E9BCCE88B6C4AA3C47F7D17228134A457F870D3800E6223A0DF07A",
    "IsSuspicious": false,
    "Reasons": null,
    "CTLogged": true,
    "KnownBadChain": false
  },
  "tls_info": {
    "Present": true,
    "Issuer": "SSL Corporation",
    "AgeDays": 5,
    "HostnameMismatch": false
  },
  "content_data": {
    "url": "https://example.com",
    "title": "Example Domain",
    "has_forms": false,
    "has_login_form": false,
    "has_payment_form": false,
    "has_personal_form": false,
    "form_count": 0,
    "forms": null,
    "iframes": null,
    "has_hidden_iframe": false,
    "has_tracking": false,
    "fetch_duration": 233844276,
    "brand_check": {
      "brand_found": "",
      "is_mismatch": false,
      "detected_names": []
    }
  },
  "domain_randomness": {
    "Domain": "example.com",
    "Label": "example",
    "Length": 7,
    "Entropy": 2.5216406363433186,
    "EntropyPerChar": 0.36023437662047403,
    "NormalizedEntropy": 0.0605009236917598,
    "VowelRatio": 0.42857142857142855,
    "DigitRatio": 0,
    "UniqueCharRatio": 0.8571428571428571,
    "LongestConsonantRun": 3,
    "BigramEnglishiness": 0.16666666666666666,
    "RandomnessScore": 0.3567918975896066,
    "IsSuspicious": false,
    "Reasons": []
  },
  "typosquat_result": {
    "is_suspicious": false
  },
  "phishing": {
    "in_database": true,
    "phish_id": 7366538,
    "phish_detail_page": "http://www.phishtank.com/phish_detail.php?phish_id=7366538",
    "verified": false,
    "verified_at": "",
    "valid": false,
    "target": "",
    "source": "phishtank",
    "from_cache": false
  },
  "performance": {
    "total_time": "1.439368926s",
    "timings": [
      {
        "task": "phishtank_check",
        "time": "1.438956594s"
      },
      {
        "task": "content_check",
        "time": "234.59063ms"
      },
      {
        "task": "dns_validity_check",
        "time": "234.104966ms"
      },
      {
        "task": "…",
        "time": "…"
      }
    ]
  },
  "result": {
    "risk_score": 5,
    "trust_score": 100,
    "final_score": 98,
    "verdict": "Safe",
    "reasons": {
      "neutral_reasons": [
        "Standard, officially recognized domain extension.",
        "No email server configured for this domain."
      ],
      "good_reasons": [
        "Global Giant: Ranked #175 worldwide.",
        "Long-standing domain history (31 years 2 months).",
        "Advanced DNS security enabled (DNSSEC)."
      ],
      "bad_reasons": null
    }
  },
  "incomplete": false,
  "errors": null
}
```

</details>

---

## Other Endpoints

The single-check endpoints below are mostly for debugging and are not covered in detail here; Swagger UI lists their parameters and responses.

All endpoints are under `/api/v1/`. `GET` endpoints accept a `url` query parameter (max 2048 chars).

### Analysis

| Method | Path | Description |
|---|---|---|
| `GET` | `/api/v1/analyze` | Full URL analysis — runs all checks in parallel, returns scored report |
| `POST` | `/api/v1/report` | Report a wrong result. JSON body: `url` (required), `verdict`, `score`, `expected_verdict` (`Safe`/`Suspicious`/`Risky`), `comment` (max 1000 chars). Appended as one JSON line to `REPORTS_FILE` (default `data/reports.jsonl`) |

### URL Structure

| Method | Path | Description |
|---|---|---|
| `GET` | `/api/v1/length` | Check if URL exceeds safe length |
| `GET` | `/api/v1/depth` | Check if URL path is suspiciously deep |
| `GET` | `/api/v1/punycode` | Detect IDN/punycode characters |
| `GET` | `/api/v1/ip/check` | Detect raw IP address usage |
| `GET` | `/api/v1/url-shortener` | Detect known URL shortener services |
| `GET` | `/api/v1/trusted-tld` | Check for trusted TLD (gov/edu) |
| `GET` | `/api/v1/risky-tld` | Check for high-risk TLD |

### DNS / Infrastructure

| Method | Path | Description |
|---|---|---|
| `GET` | `/api/v1/ip/resolve` | Resolve domain to IP addresses |
| `GET` | `/api/v1/rank` | Global popularity rank of the domain |
| `GET` | `/api/v1/domain-info` | WHOIS / RDAP registration data |

### Security

| Method | Path | Description |
|---|---|---|
| `GET` | `/api/v1/hsts` | Check if host enforces HSTS |
| `GET` | `/api/v1/redirects` | Follow and report redirect chain |
| `GET` | `/api/v1/status-code` | Fetch HTTP status code |

### Utility

| Method | Path | Description |
|---|---|---|
| `GET` | `/health` | Service liveness check |
| `GET` | `/api/v1/health` | Same, versioned |
| `GET` | `/api/v1/screenshot` | Headless screenshot of the URL |
| `DELETE` | `/api/v1/cache` | Flush the Valkey cache |
| `GET` | `/metrics` | Prometheus metrics scrape endpoint |
| `GET` | `/swagger/*` | Swagger UI and spec |

---

## Rate Limiting

20 requests per minute per IP. Headers returned on every response:

```
X-RateLimit-Limit: 20
X-RateLimit-Remaining: 19
X-RateLimit-Reset: <unix timestamp>
```

Exceeding the limit returns `429 Too Many Requests`.

---

## Error Format

All 4xx/5xx responses return JSON:

```json
{ "error": "description of the problem" }
```
