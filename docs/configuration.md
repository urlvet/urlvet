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
| `CORS_ALLOWED_ORIGINS` | No | Comma-separated allowed origins (default `*`) — e.g. `https://url.vet,https://www.url.vet` |
| `PORT` | No | HTTP port (default `8080`) |
| `TRUSTED_PROXIES` | No | Addresses whose `X-Forwarded-For` / `X-Real-IP` the backend believes: comma-separated IPs or CIDRs. Default: loopback and private networks, which fits a reverse proxy on the same host or Docker network. `none` ignores the headers; `all` trusts anyone (the old, spoofable behaviour). See [Client IP and rate limiting](#client-ip-and-rate-limiting) |
| `TRUSTED_PLATFORM` | No | Header a CDN or platform sets to the client IP, read before anything else: `cloudflare`, `flyio`, `google-app-engine`, or a header name. Only for a platform that connects to the backend directly, and only if nothing else can reach the backend (firewall it to the platform's IP ranges). See [Client IP and rate limiting](#client-ip-and-rate-limiting) |
| `LOG_TIMEZONE` | No | IANA timezone for log timestamps (default `UTC`) — e.g. `Asia/Kolkata`, `America/New_York` |
| `LOG_DIR` | No | Directory for rotating daily log files — e.g. `logs`. Leave empty to disable file logging |
| `ADMIN_PASSWORD_HASH` | Yes | Argon2id hash of the admin password — see [security.md](security.md#setup) |
| `ADMIN_JWT_SECRET` | Yes | Signing secret for session tokens — `openssl rand -hex 32` |

### Client IP and rate limiting

The rate limiter allows 20 requests a minute per client IP. Behind a reverse proxy, the backend only sees the proxy's address, so it takes the client IP from the `X-Forwarded-For` (or `X-Real-IP`) header the proxy adds. Clients can send that header too, so the backend only believes it when the request comes from a trusted proxy. Otherwise anyone could set a new IP on each request and never be limited.

By default the trusted proxies are loopback and private addresses (`127.0.0.0/8`, `::1`, `10.0.0.0/8`, `172.16.0.0/12`, `192.168.0.0/16`, `fc00::/7`). A proxy on the Docker host, as in [deployment.md](deployment.md#3-reverse-proxy-caddy), connects from the Docker bridge gateway (`172.17.0.1` or similar), which is covered, so most setups need no change. `X-Forwarded-For` is read right to left: the backend skips trusted addresses and uses the first one it doesn't trust, so anything a client puts in front is never reached.

| Setup | What to set |
| --- | --- |
| Caddy, nginx or Traefik on the same host or Docker network | Nothing. For nginx, see the headers in [deployment.md](deployment.md#4-nginx-alternative) |
| Cloudflare in front of Caddy or nginx | Nothing on the backend. Make the proxy trust Cloudflare: [deployment.md](deployment.md#behind-cloudflare) |
| Proxy on another machine with a private IP | Nothing if its IP is in a private range, otherwise `TRUSTED_PROXIES=<its IP or CIDR>` |
| Cloudflare connecting straight to the backend, no local proxy | `TRUSTED_PLATFORM=cloudflare`, and firewall the backend port to [Cloudflare's IPs](https://www.cloudflare.com/ips/) |
| Fly.io or Google App Engine | `TRUSTED_PLATFORM=flyio` or `google-app-engine` |
| Exposed directly, no proxy | Nothing, or `TRUSTED_PROXIES=none` to ignore the headers entirely |

Setting `TRUSTED_PROXIES` replaces the default list, so include every proxy in the chain. On startup the backend logs which addresses it trusts.

**If every user hits the limit together**, the backend is probably keying on your proxy's address rather than your users'. That happens when the proxy (or a CDN) connects from a public IP the backend doesn't trust, and the backend logs a warning the first time it sees such a request. Set `TRUSTED_PROXIES` or `TRUSTED_PLATFORM` for it as above. `TRUSTED_PROXIES=all` brings back the behaviour from before these settings existed, where any client can choose its own IP.

**`TRUSTED_PLATFORM` is believed from any address**, so it is only safe when nothing but that platform can reach the backend. If clients can connect directly (for example, to the server's IP), they can set the header themselves. When a local proxy sits between the platform and the backend, leave `TRUSTED_PLATFORM` unset and configure the proxy instead.

If you use `TRUSTED_PLATFORM`, make sure only that platform can reach the backend, for example with a firewall that allows only its IP ranges. On a VPS, use your provider's firewall (see [deployment.md](deployment.md#5-firewall)): Docker-published ports skip UFW.

**Private ranges include your LAN.** A machine on the same private network can still set its own IP. If that matters, set `TRUSTED_PROXIES` to just your proxy's address.

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
