# ============================
# 1) BUILD STAGE
# ============================
FROM node:24-alpine AS builder

WORKDIR /app

ENV DOCKER_BUILD=1

# Copy package files first for caching
COPY web/website/package*.json ./

# Install dependencies
RUN npm ci

# Copy full frontend source
COPY web/website .

# Build static site
RUN npm run build

# ============================
# 2) RUNTIME STAGE
# ============================
FROM caddy:2-alpine

# Set working dir for Caddy
WORKDIR /srv

# Copy built frontend from builder
COPY --from=builder /app/build /srv

# Copy Caddyfile
COPY docker/prod/Caddyfile /etc/caddy/Caddyfile

# Expose ports
EXPOSE 80

# Default Caddy command
CMD ["caddy", "run", "--config", "/etc/caddy/Caddyfile", "--adapter", "caddyfile"]
