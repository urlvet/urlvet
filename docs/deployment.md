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
    }
}
```

Use Certbot for certificates: `sudo certbot --nginx -d yourdomain.com -d api.yourdomain.com`

---

## 5. Firewall

Only expose ports 80 and 443 to the internet. All other ports (8080, 3000, 6379, 9222) must remain internal.

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
