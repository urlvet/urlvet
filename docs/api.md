# API Reference

## Interactive Docs (Swagger UI)

When the server is running, the full OpenAPI spec is browsable at:

```
http://localhost:8080/swagger/index.html
```

The raw spec files are also available:

| Format | URL |
|---|---|
| JSON | `http://localhost:8080/swagger/doc.json` |
| YAML | `server/internal/docs/swagger.yaml` (generated, committed) |

To regenerate the spec after editing handler annotations:

```bash
cd server
swag init -g cmd/urlvet/urlvet.go -o internal/docs
```

Install the CLI once with: `go install github.com/swaggo/swag/cmd/swag@v1.16.4`

---

## Endpoints

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

---

## Example

```bash
curl "http://localhost:8080/api/v1/analyze?url=https://example.com"
```
 <details>
<summary>Example API response</summary>
 <pre><code class="language-json">
{
  "url": "https://example.com",
  "domain": "example.com",
  "features": {
    "rank": 175,
    "tld": {
      "tld": "com",
      "is_trusted_tld": false,
      "is_risky_tld": false,
      "is_icann": true
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
        "found": [],
        "categories": {}
      }
    }
  },
  "infrastructure": {
    "ip_addresses": [
      "172.66.147.243",
      "104.20.23.154",
      "2606:4700:10::6814:179a",
      "2606:4700:10::ac42:93f3"
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
    "updated": "2026-01-16T18:26:50Z",
    "expiry": "2026-08-13T04:00:00Z",
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
    "age_human": "30 years 8 months",
    "age_days": 11202,
    "raw": "{\"ldhName\":\"EXAMPLE.COM\",\"nameservers\":[{\"ldhName\":\"ELLIOTT.NS.CLOUDFLARE.COM\"},{\"ldhName\":\"HERA.NS.CLOUDFLARE.COM\"}],\"events\":[{\"eventAction\":\"registration\",\"eventDate\":\"1995-08-14T04:00:00Z\"},{\"eventAction\":\"expiration\",\"eventDate\":\"2026-08-13T04:00:00Z\"},{\"eventAction\":\"last changed\",\"eventDate\":\"2026-01-16T18:26:50Z\"},{\"eventAction\":\"last update of RDAP database\",\"eventDate\":\"2026-04-15T19:04:14Z\"}],\"entities\":[{\"roles\":[\"registrar\"],\"vcardArray\":[\"vcard\",[[\"version\",{},\"text\",\"4.0\"],[\"fn\",{},\"text\",\"RESERVED-Internet Assigned Numbers Authority\"]]]}],\"status\":[\"client delete prohibited\",\"client transfer prohibited\",\"client update prohibited\"],\"secureDNS\":{\"delegationSigned\":true}}",
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
    "Issuer": "Cloudflare TLS Issuing ECC CA 1",
    "NotBefore": "2026-04-02T21:18:57Z",
    "NotAfter": "2026-07-01T21:24:46Z",
    "AgeDays": 12,
    "Fingerprint": "1AF627C6C2AC992E3C9102438F467C4C238D3112325AC7CF9003D77F75EFFFBA",
    "IsSuspicious": false,
    "Reasons": null,
    "CTLogged": true,
    "KnownBadChain": false
  },
  "tls_info": {
    "Present": true,
    "Issuer": "CLOUDFLARE, INC.",
    "AgeDays": 12,
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
    "fetch_duration": 137804093,
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
    "Entropy": 2.521640636343318,
    "EntropyPerChar": 0.36023437662047403,
    "NormalizedEntropy": 0.06050092369175979,
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
    "from_cache": false,
    "raw_response": {
      "meta": {
        "timestamp": "2026-04-15T19:04:30+00:00",
        "serverid": "e5f3084e",
        "status": "success",
        "requestid": "172.17.128.1.69dfe13e5ee121.10644345"
      },
      "results": {
        "url": "https://example.com",
        "in_database": true,
        "phish_id": 7366538,
        "phish_detail_page": "http://www.phishtank.com/phish_detail.php?phish_id=7366538",
        "verified": false,
        "verified_at": null,
        "valid": false
      }
    }
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
        "Long-standing domain history (30 years 8 months).",
        "Advanced DNS security enabled (DNSSEC)."
      ],
      "bad_reasons": null
    }
  },
  "incomplete": false,
  "errors": null
}
</code></pre>
</details>
