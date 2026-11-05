# Changelog

All notable changes to this project are listed here.

The format is based on [Keep a Changelog](https://keepachangelog.com/). Versions follow [Semantic Versioning](https://semver.org/); url.vet is in beta (0.x), see [Versioning](.github/CONTRIBUTING.md#versioning).

## [0.4.0] – 2026-11-05 – Beta

### Added
- Detection of fake Adobe PDF viewers, fake Chrome downloads, "AI platform" investment scams and fake Outlook login paths
- Short brand names next to security words in domains (`hsbc-sec`, `verify-dhl`)
- Cloudflare's "Suspected Misleading Website" warning
- More brands, including Parivahan, Government of India, Citi, Banco Galicia, Codashop, CouriersPlease, TokenPocket and SushiSwap
- More free-hosting and link-page hosts

### Changed
- Phishing kits on old hacked domains no longer get credit for the domain's age
- Better redirect handling for hosting platforms and brands with several domains
- Home page tagline restyled

### Fixed
- Google sign-in copies titled "Google Accounts" weren't matched
- Brand names in Chinese and Hindi titles weren't matched
- IP lookup now uses the exact host scanned, not just its main domain

### Security
- Rate-limit bypass with a fake `X-Forwarded-For`: the header is now trusted only from private networks. New `TRUSTED_PROXIES` and `TRUSTED_PLATFORM` settings (see docs/configuration.md)
- Admin login limited to 5 attempts a minute per IP
- Updated website dependencies with security fixes (undici, Vite, Vitest, Satori)

## [0.3.0] – 2026-10-08 – Beta

### Added
- Local threat lists (PhishTank, URLhaus, OpenPhish) plus Google Safe Browsing and Web Risk lookups
- Checks for phishing kits, short links, IP hostnames and brand mismatches; redirects are judged where they land
- Users can report wrong results; admin dashboard with a reports tab
- Vetty, a guide for first-time visitors; About, How it works and Privacy pages
- Share preview images; links found in pasted messages; installable web app with share target
- Rebuilt Chrome extension

### Changed
- Editorial redesign of the website; Tailwind 4
- Fewer false alarms on well-known sites and banks; scans finish on time

### Security
- Client IPs are hashed in the rate limiter, and IPs and scanned URLs are kept out of server logs
- Stored screenshots are deleted after 24 hours

## [0.2.0] – 2026-06-09 – Beta

### Added
- url.vet name and website
- Typosquatting check
- Prometheus metrics, Swagger docs and unit tests
- Admin page
- Light and dark mode, 404 page, About and How it works pages
- Structured logging with slog

### Fixed
- Cookie issue in the website

## [0.1.0] – 2026-01-26 – Beta

First release of the Go rewrite.

### Added
- Go API with URL checks: domain rank, WHOIS and domain age, punycode, SSL, redirects, screenshots
- Score and verdict from heuristics
- Svelte website with shareable results
- Docker setup, with Chrome in its own container

## Legacy versions

These came before the Go rewrite, are no longer maintained, and share no code with the versions above.

- **SafeSurf Flask app** (tag `safesurf-flask`, 2024-09-01): a Python/Flask website that checked links with rules, domain data, SSL details and screenshots. No machine learning.
- **SafeSurf ML** (tag `safesurf-ml`, 2021-11-24): the first version, which classified links with an XGBoost model trained on lists of phishing and legitimate URLs.

[0.4.0]: https://github.com/urlvet/urlvet/compare/v0.3.0...v0.4.0
[0.3.0]: https://github.com/urlvet/urlvet/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/urlvet/urlvet/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/urlvet/urlvet/releases/tag/v0.1.0
