# Deployment

## Prerequisites

- A Linux VPS (Ubuntu 22.04+ or Debian 12+ recommended)
- Docker Engine 24+ and Docker Compose v2
- A domain name with DNS pointing to the server
- Ports 80 and 443 open in the firewall

---

## 1. Clone and configure

```bash
git clone https://github.com/urlvet/urlvet.git
cd urlvet
```

Create the backend env file:

```bash
cp server/.env.example server/.env
```

Edit `server/.env` — all values are required for production:

```env
ENV=PROD
PORT=8080
CACHE_ADDR=urlvet-valkey:6379
CACHE_PASSWORD=<generate: openssl rand -hex 32>
CACHE_DB=0
CACHE_POOL_SIZE=50
CACHE_MIN_IDLE_CONNS=10
CORS_ALLOWED_ORIGINS=https://yourdomain.com
ADMIN_PASSWORD_HASH='<argon2id hash — see docs/security.md>'
ADMIN_JWT_SECRET=<generate: openssl rand -hex 32>
LOG_TIMEZONE=UTC
LOG_DIR=logs
```

Threat feeds are optional and also set in `server/.env`. PhishTank's list is on with no key; `URLHAUS_AUTH_KEY` and `SAFE_BROWSING_API_KEY` are free and need no card. Each feed has its own licence, so check it fits your use: see [configuration.md](configuration.md#threat-feeds-also-serverenv). Downloaded lists are cached in `/app/data/feeds`, on the backend's data volume, so they survive restarts and rebuilds.

Create the Compose env file (used for Valkey password substitution):

```bash
echo "CACHE_PASSWORD=<same value as above>" > docker/prod/.env
```

---

## 2. Start the stack

```bash
make start
```

This builds and starts four containers: `urlvet-backend`, `urlvet-web`, `urlvet-chrome`, `urlvet-valkey`.

Verify all containers are healthy:

```bash
make ps
# or
docker compose -f docker/prod/docker-compose.prod.yml ps
```

Check the backend is responding:

```bash
curl http://localhost:8080/health
# → {"status":"ok"}
```

---

## 3. Reverse proxy (Caddy)

Caddy is the recommended reverse proxy — it handles HTTPS automatically.

Install Caddy on the host (not in Docker):

```bash
sudo apt install -y debian-keyring debian-archive-keyring apt-transport-https curl
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | sudo gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | sudo tee /etc/apt/sources.list.d/caddy-stable.list
sudo apt update && sudo apt install caddy
```

`/etc/caddy/Caddyfile`:

```
yourdomain.com {
    reverse_proxy localhost:3000
}

api.yourdomain.com {
    reverse_proxy localhost:8080
}
```

```bash
sudo systemctl reload caddy
```

Caddy provisions TLS certificates from Let's Encrypt automatically.

Caddy replaces any `X-Forwarded-For` the client sent with the client's real address, and the backend trusts it because Caddy connects from the Docker gateway, a private address. So rate limiting works per user with no extra settings. See [configuration.md](configuration.md#client-ip-and-rate-limiting).

### Behind Cloudflare

With Cloudflare's proxy on (orange cloud), Caddy sees Cloudflare's addresses, not your users'. Left as is, the backend rate-limits per Cloudflare server, so unrelated users share a limit. Tell Caddy (2.7 or later) to trust Cloudflare and read `CF-Connecting-IP`, then pass the result on:

```
{
    servers {
        # Cloudflare's ranges from https://www.cloudflare.com/ips/ (they rarely
        # change; check that page when upgrading).
        trusted_proxies static 173.245.48.0/20 103.21.244.0/22 103.22.200.0/22 103.31.4.0/22 141.101.64.0/18 108.162.192.0/18 190.93.240.0/20 188.114.96.0/20 197.234.240.0/22 198.41.128.0/17 162.158.0.0/15 104.16.0.0/13 104.24.0.0/14 172.64.0.0/13 131.0.72.0/22 2400:cb00::/32 2606:4700::/32 2803:f800::/32 2405:b500::/32 2405:8100::/32 2a06:98c0::/29 2c0f:f248::/32
        client_ip_headers CF-Connecting-IP
    }
}

yourdomain.com {
    reverse_proxy localhost:3000
}

api.yourdomain.com {
    reverse_proxy localhost:8080 {
        # {client_ip} is CF-Connecting-IP for requests from Cloudflare and
        # the connecting address for anything else, so a request sent
        # straight to the server can't claim another IP.
        header_up X-Forwarded-For {client_ip}
    }
}
```

The backend needs no change. Don't set `TRUSTED_PLATFORM=cloudflare` here: the backend would then believe `CF-Connecting-IP` from anyone who reaches Caddy directly, not just Cloudflare. If you build Caddy with the [caddy-cloudflare-ip](https://github.com/WeidiDeng/caddy-cloudflare-ip) module, `trusted_proxies cloudflare` replaces the static list and keeps it up to date.

---

## 4. Nginx alternative

If you prefer Nginx:

```nginx
server {
    listen 80;
    server_name yourdomain.com;
    return 301 https://$host$request_uri;
}

server {
    listen 443 ssl;
    server_name yourdomain.com;

    ssl_certificate     /etc/letsencrypt/live/yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/yourdomain.com/privkey.pem;

    location / {
        proxy_pass http://localhost:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}

server {
    listen 443 ssl;
    server_name api.yourdomain.com;

    ssl_certificate     /etc/letsencrypt/live/api.yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/api.yourdomain.com/privkey.pem;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        # Required: without it nginx passes on whatever X-Forwarded-For the
        # client sent, and the backend reads that first.
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

Use Certbot for certificates: `sudo certbot --nginx -d yourdomain.com -d api.yourdomain.com`

Behind Cloudflare, use nginx's [realip module](https://nginx.org/en/docs/http/ngx_http_realip_module.html) so `$remote_addr` is the user's address: one `set_real_ip_from` line per [Cloudflare range](https://www.cloudflare.com/ips/), then `real_ip_header CF-Connecting-IP;`.

---

## 5. Firewall

Only expose ports 80 and 443 to the internet. All other ports (8080, 3000, 6379, 9222) must remain internal.

Docker-published ports (`8080` and `3000` in `docker-compose.prod.yml`) skip UFW's rules, because Docker writes its own iptables rules. Block them with your provider's firewall (e.g. a DigitalOcean Cloud Firewall), or publish them on localhost only (`"127.0.0.1:8080:8080"`), which still works with Caddy or nginx on the host. Rate limiting keeps working if the port is reachable, since direct requests are keyed on their own address, but nothing should reach the backend without going through the proxy.

**If you use `TRUSTED_PLATFORM`, make sure only that platform can reach the backend, for example with a firewall that allows only its IP ranges.** This applies when a CDN such as Cloudflare connects straight to the backend with no local proxy. Otherwise anyone who finds the server's IP can send the platform's header themselves and pick their own rate-limit key. For Cloudflare, allow 80 and 443 only from [its ranges](https://www.cloudflare.com/ips/) in your provider's firewall (a DigitalOcean Cloud Firewall is free), keep SSH open for yourself, and leave outbound traffic allowed, since scans fetch other sites.

```bash
# UFW example
sudo ufw allow 22/tcp    # SSH
sudo ufw allow 80/tcp    # HTTP (for ACME challenge)
sudo ufw allow 443/tcp   # HTTPS
sudo ufw enable
```

---

## 6. Updates

```bash
git pull
make start   # rebuilds images and restarts containers
```

To rebuild only the backend without restarting other services:

```bash
docker compose -f docker/prod/docker-compose.prod.yml up -d --build backend
```

---

## 7. Monitoring

Prometheus metrics are exposed at `http://localhost:8080/metrics` (not proxied — internal only). Scrape from a Prometheus instance on the same host or VPN.

Tail live logs:

```bash
make logs
```

See [maintenance.md](maintenance.md) for cache management and log rotation.
