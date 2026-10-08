# Maintenance

Routine operational tasks for keeping a url.vet deployment healthy.

---

## Cache management

### Flush the entire cache

Requires a valid admin Bearer token. Obtain one via `POST /api/v1/admin/login` with your admin password.

```bash
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/admin/login \
  -H "Content-Type: application/json" \
  -d '{"password":"your-password"}' | jq -r '.token')

curl -X DELETE http://localhost:8080/api/v1/admin/cache \
  -H "Authorization: Bearer $TOKEN"
```

### Delete a specific cache entry

```bash
curl -X DELETE "http://localhost:8080/api/v1/admin/cache/analyze_result:https://example.com" \
  -H "Authorization: Bearer $TOKEN"
```

### Inspect cached keys directly in Valkey

```bash
docker exec -it urlvet-valkey redis-cli -a "$CACHE_PASSWORD"

# List all keys
KEYS *

# Check TTL of a key
TTL "analyze_result:https://example.com"

# Delete a key
DEL "analyze_result:https://example.com"

# Memory usage
INFO memory
```

### Cache TTLs

| Cache type | TTL | Key prefix |
|---|---|---|
| Full analysis result | 24 hours | `analyze_result:` |
| Content analysis | 3 hours | `content_check:` |
| HTTP combined result | 3 hours | `http_combined:` |
| TLS / certificate (per host) | 24 hours | `tls_combined:` |
| WHOIS / RDAP | 24 hours | `whois_lookup:` |
| Google Safe Browsing | 1 hour | `safe_browsing:` |
| Google Web Risk | 1 hour | `webrisk:` |
| PhishTank API (only while its list isn't loaded) | 3 hours | `phishtank:` |

All TTLs live in `server/internal/constants/ttls.go`. Brand checks are re-run on every scan from the cached page title, so brand-list changes apply without flushing `content_check:`. Other rule changes apply once the 24-hour results expire, or after flushing `analyze_result:`.

---

## Log management

Logs are written to `server/logs/YYYY-MM/DD.log` (when `LOG_DIR=logs`) and always to stderr.

### View today's log

```bash
cat server/logs/$(date +%Y-%m)/$(date +%d).log

# Follow live
tail -f server/logs/$(date +%Y-%m)/$(date +%d).log
```

### View container logs

```bash
make logs          # production
make dev-logs      # development
```

### Purge old logs

Log files rotate automatically at midnight — the logger opens a new file per calendar day. Clean up old files periodically:

```bash
# Delete logs older than 30 days
find server/logs -name "*.log" -mtime +30 -delete

# Remove empty monthly directories
find server/logs -type d -empty -delete
```

Add to a host cron job (`crontab -e`):

```cron
0 3 * * * find /path/to/urlvet/server/logs -name "*.log" -mtime +30 -delete
```

---

## Updating data files

### Domain rank list

The top-1M domain rank list is bundled at build time from `server/assets/`. To refresh it with the latest Tranco list:

```bash
cd server/assets
curl -L "https://tranco-list.eu/top-1m.csv.zip" -o top-1m.csv.zip
unzip -o top-1m.csv.zip
rm top-1m.csv.zip
```

Rebuild and restart the backend to pick up the new list.

### Brand list

The 514 brands used for brand-mismatch and combo-squatting checks live in `server/internal/constants/brands.go`, grouped by sector. Each entry has:

- `TitleKeywords` — phrases that mean a page claims to be the brand wherever they appear ("paypal", "banco santander", "vodafone login")
- `Names` — bare names ordinary sites also use in titles ("vodafone", "amazon"). These only count on a page asking for a login or payment, or on a hosting subdomain. Put a brand's name here, not in `TitleKeywords`, when it's also an everyday word, a place or a product category
- `OfficialDomains` — the brand's main domains, including country domains (a missing one makes the real site look like an impersonator). Their names are also what combo-squatting looks for inside other domains, so only list domains whose name *is* the brand
- `AlsoOfficial` — the brand's other domains (Google's `research.google`, Amazon's `primevideo.com`). They count as the brand's own everywhere, but their names aren't looked for in other domains, so "research" doesn't become a brand. Put any domain whose name is an ordinary word here
- `Platforms` — hosting the brand runs for anyone (`blogspot.com` for Google). A page there isn't the brand's own, but one that forwards to the brand's own site is the brand moving its content
- `ExactTitles` — whole page titles that are the brand's login but too generic to match inside a longer title ("Sign in to your account" is Microsoft's)
- `OwnNames` — for groups whose members each run their own domain (`sparkasse-hannover.de`): a domain containing one counts as the brand's own

Restart the backend after editing. Typosquatting by edit distance uses the top 5,000 entries of the domain rank list instead.

### Hosting, shortener and upload-host lists

Also in `server/internal/constants/`:

- `trusted_hosting_platforms.go` and `customer_site_hosts.go` — services whose subdomains belong to customers. A customer site doesn't inherit the provider's rank or age. Add a provider here when it isn't on the Public Suffix List; `ProviderSubdomains` keeps the provider's own subdomains (`account.squarespace.com`) as the provider's. `IPHostnameServices` (also in `customer_site_hosts.go`) are names that stand for an IP address (`1-2-3-4.sslip.io`, cPanel's `cprapid.com`), scored like a raw IP
- `provider_service_hosts.go` — a provider's own endpoints on those suffixes
- `user_upload_hosts.go` — where anyone can upload files
- `url_shorteners.go` — short-link services; `ShortLinkPaths` for sites where only some paths redirect (email click tracking, Flowcode QR codes); `UserPageHosts` for services that give each user a page under a path (`linktr.ee/someone`), which doesn't inherit the service's rank

To find providers missing from these lists, match a phishing feed against the top-1M rank list and look for ranked domains with many different subdomains, or many short random paths, in the feed.

### Threat lists

Downloaded lists refresh on their own and are cached under `FEEDS_DIR` (`data/feeds`, on the backend's data volume in production), so restarts don't download them again. The backend logs each load and refresh:

```bash
docker logs urlvet-backend 2>&1 | grep "threat feed"
```

A failed download is retried after 30 minutes, keeping the previous copy in use. To force a fresh download, delete the file (e.g. `data/feeds/urlhaus.txt`) and restart the backend. Without a PhishTank API key only a few downloads a day are allowed, so avoid restarting repeatedly with that file deleted.

---

## Container updates

```bash
# Pull latest code and rebuild all images
git pull
make start
```

Update only the backend:

```bash
docker compose -f docker/prod/docker-compose.prod.yml up -d --build backend
```

Update the headless Chrome image:

```bash
docker pull chromedp/headless-shell:latest
docker compose -f docker/prod/docker-compose.prod.yml up -d chrome
```

---

## Valkey backup and restore

Valkey snapshots to the `valkey_data` Docker volume automatically. To back up manually:

```bash
docker run --rm \
  -v urlvet_valkey_data:/data \
  -v $(pwd):/backup \
  alpine tar czf /backup/valkey-$(date +%Y%m%d).tar.gz /data
```

To restore:

```bash
docker run --rm \
  -v urlvet_valkey_data:/data \
  -v $(pwd):/backup \
  alpine tar xzf /backup/valkey-YYYYMMDD.tar.gz -C /
```

---

## Health checks

```bash
# Backend liveness
curl http://localhost:8080/health

# Prometheus metrics snapshot
curl -s http://localhost:8080/metrics | grep urlvet_

# Container status and restart counts
docker compose -f docker/prod/docker-compose.prod.yml ps
docker stats --no-stream
```
