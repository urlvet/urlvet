# Configuration

All configuration is via environment variables. Copy `server/.env.example` → `server/.env` and fill in your values.

## Backend (`server/.env`)

| Variable | Required | Description |
|---|---|---|
| `CACHE_ADDR` | Yes | Valkey/Redis address — `urlvet-valkey:6379` |
| `CACHE_PASSWORD` | Yes (prod) | Cache auth password — must match `CACHE_PASSWORD` set in `docker-compose.prod.yml` |
| `CACHE_DB` | No | Redis DB index (default `0`) |
| `CACHE_POOL_SIZE` | No | Connection pool size (default `50`) |
| `CACHE_MIN_IDLE_CONNS` | No | Min idle connections (default `10`) |
| `ENV` | No | Set to `DEV` for colored logs, debug endpoint, and verbose output |
| `CORS_ALLOWED_ORIGINS` | No | Comma-separated allowed origins (default `*`) — e.g. `https://url.vet,https://url.vet` |
| `PORT` | No | HTTP port (default `8080`) |
| `LOG_TIMEZONE` | No | IANA timezone for log timestamps (default `UTC`) — e.g. `Asia/Kolkata`, `America/New_York` |
| `LOG_DIR` | No | Directory for rotating daily log files — e.g. `logs`. Leave empty to disable file logging |
| `ADMIN_PASSWORD_HASH` | Yes | Argon2id hash of the admin password — see [security.md](security.md#setup) |
| `ADMIN_JWT_SECRET` | Yes | Signing secret for session tokens — `openssl rand -hex 32` |

### Threat feeds (also `server/.env`)

Lists of reported phishing and malware links, downloaded on a schedule and checked on every scan, plus optional Google lookups. Each source has its own licence; only PhishTank is on by default.

| Variable | Cost | Use allowed | Description |
|---|---|---|---|
| `FEED_PHISHTANK` | Free, no account | Commercial OK | PhishTank's verified, online list, refreshed every 6 hours (hourly with a key). On unless set to `0`, which means no PhishTank at all. Until the list has loaded, scans ask the PhishTank API instead (with or without a key); once it's loaded, they don't. |
| `PHISHTANK_API_KEY` | Free | — | Optional. PhishTank registration is currently closed. |
| `PHISHTANK_USER_AGENT` | — | — | PhishTank asks for `phishtank/<username>`. |
| `URLHAUS_AUTH_KEY` | Free key, no card ([auth.abuse.ch](https://auth.abuse.ch)) | Not-for-profit; commercial needs Spamhaus | abuse.ch's malware URLs, refreshed hourly. Set the key to enable. |
| `SAFE_BROWSING_API_KEY` | Free, no billing account | Non-commercial only | Google Safe Browsing v5. Sends Google only 4-byte hash prefixes, never the link. |
| `WEBRISK_API_KEY` | Billing account; 100k lookups/month free, then $0.50 per 1,000 | Commercial OK | Google Web Risk, Safe Browsing's commercial counterpart. Sends Google the link. |
| `FEED_OPENPHISH` | Free | Personal/academic research only | `1` enables OpenPhish's community feed. Its terms forbid use for detection or customer protection without their written permission. |
| `FEEDS_DIR` | — | — | Where downloaded lists are cached (default `data/feeds`, on the prod data volume). |
| `THREAT_LOOKUPS` | — | — | `0` switches every lookup above off, whatever keys are set: no lists are loaded or downloaded, and PhishTank, Safe Browsing and Web Risk aren't asked. For testing the scanner's own checks, for example against a batch of known phishing links, without calling anyone or hitting their rate limits. |

If you run url.vet under a commercial licence, the feed keys are yours and each feed's terms apply to you.

## Docker Compose (dev: `docker/dev/.env`, prod: `docker/prod/.env`)

Create the appropriate `.env` file with secrets used by docker-compose variable substitution:

| Variable | Required | Description |
|---|---|---|
| `CACHE_PASSWORD` | Yes | Password set on the Valkey container (`--requirepass`) and passed to the backend as `CACHE_PASSWORD`. Generate with `openssl rand -hex 32`. |

## Frontend (`web/website/.env`)

| Variable | Required | Description |
|---|---|---|
| `PUBLIC_BASE_URL` | No | Go API base URL (default `http://localhost:8080/api/v1`) |
