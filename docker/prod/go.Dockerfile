# ============================
# 1) BUILD STAGE
# ============================
FROM golang:1.27-alpine AS builder

WORKDIR /app

# Install git — required for go mod download (some deps)
RUN apk add --no-cache git

# Copy go.mod + go.sum first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy full source
COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o urlvet ./cmd/urlvet


# ============================
# 2) RUNTIME STAGE
# ============================
FROM alpine:latest

# Basic dependencies only (timezone, certificates)
RUN apk add --no-cache \
    ca-certificates \
    tzdata

WORKDIR /app

# Copy binary
COPY --from=builder /app/urlvet .

# Copy assets folder
COPY --from=builder /app/assets ./assets
COPY --from=builder /app/.env ./.env

ENV PORT=8080
EXPOSE 8080

CMD ["./urlvet"]
