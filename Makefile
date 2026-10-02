# ============================================
#   Makefile
# ============================================

.DEFAULT_GOAL := help

# --- Project Paths ---
BACKEND_DIR := server
FRONTEND_DIR := web/website

# --- Docker Compose Files ---
DOCKER_DEV  := docker/dev/docker-compose.dev.yml
DOCKER_PROD := docker/prod/docker-compose.prod.yml

DC := docker compose
compose-dev  = $(DC) -f $(DOCKER_DEV)
compose-prod = $(DC) -f $(DOCKER_PROD)

# --- Helpers: Colors ---
RESET := \033[0m
BLUE  := \033[36m
GREEN := \033[32m

info = @echo "$(BLUE)==>$(RESET) $1"
success = @echo "$(GREEN)✔$(RESET) $1"

# --- Helpers: Timer ---
start-timer:
	@START_TIME=$$(date +%s); echo $$START_TIME > .make_timer

timer:
	@START_TIME=$$(cat .make_timer); \
	END_TIME=$$(date +%s); \
	DURATION=$$((END_TIME - START_TIME)); \
	echo ""; \
	echo "⏱  Done in $${DURATION}s"; \
	rm -f .make_timer




# ============================================
# Default Target
# ============================================
help: ## Show all commands
	@echo ""
	@echo "Available commands:"
	@grep -E '^[a-zA-Z0-9_-]+:.*?## ' $(MAKEFILE_LIST) \
		| sort \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "  $(BLUE)%-22s$(RESET) %s\n", $$1, $$2}'
	@echo ""


# ============================================
# Production (Docker)
# ============================================

build: start-timer ## Build production Docker images
	$(info "Building production images...")
	$(compose-prod) build
	$(timer)

up: start-timer ## Start production (detached)
	$(info "Starting production stack...")
	$(compose-prod) up -d
	$(success "Production started")
	$(timer)

start: start-timer ## Build and start production stack
	$(info "Building and starting production stack...")
	$(compose-prod) build
	$(compose-prod) up -d
	$(success "Production started")
	$(timer)

build-b: start-timer ## Build backend containers only (backend, chrome, valkey)
	$(info "Building backend containers...")
	$(compose-prod) build backend chrome valkey
	$(timer)

up-b: start-timer ## Start backend containers only (backend, chrome, valkey)
	$(info "Starting backend containers...")
	$(compose-prod) up -d backend chrome valkey
	$(success "Backend containers started")
	$(timer)

start-b: start-timer ## Build and start backend containers only
	$(info "Building and starting backend containers...")
	$(compose-prod) build backend chrome valkey
	$(compose-prod) up -d backend chrome valkey
	$(success "Backend containers started")
	$(timer)

down: ## Stop production
	$(info "Stopping production...")
	$(compose-prod) down
	$(success "Production stopped")

logs: ## Tail production logs
	$(compose-prod) logs -f


# ============================================
# Development (Docker)
# ============================================

dev-build: start-timer ## Build dev images only
	$(info "Building dev images...")
	$(compose-dev) build
	$(timer)

dev-start: start-timer ## Build and start dev stack
	$(info "Building and starting dev stack...")
	$(compose-dev) build
	$(compose-dev) up -d
	$(success "Dev environment up")
	$(timer)

dev-up: start-timer ## Start dev environment
	$(info "Starting dev stack...")
	$(compose-dev) up -d
	$(success "Dev environment up")
	$(timer)

dev-down: ## Stop dev environment
	$(info "Stopping dev...")
	$(compose-dev) down
	$(success "Dev environment stopped")

dev-logs: ## Tail dev logs
	$(compose-dev) logs -f

dev-restart: ## Restart dev environment
	$(info "Restarting dev environment...")
	$(compose-dev) down
	$(compose-dev) up -d
	$(success "Restart complete")


# Utility Targets
ps: ## Show running containers (dev)
	$(compose-dev) ps

sh-backend: ## Enter backend container shell (dev)
	$(compose-dev) exec backend sh


# ============================================
# Local (No Docker)
# ============================================

local-build-backend: start-timer ## Local Go build
	$(info "Building backend locally...")
	cd $(BACKEND_DIR) && go build -o urlvet ./cmd/urlvet
	$(timer)

local-build-frontend: start-timer ## Local Svelte build
	$(info "Building frontend locally...")
	cd $(FRONTEND_DIR) && npm install && npm run build
	$(timer)

local-run-backend: ## Run backend via Air
	cd $(BACKEND_DIR) && air

local-run-frontend: ## Run Svelte dev server
	cd $(FRONTEND_DIR) && npm run dev


# ============================================
# Testing / Lint / CI
# ============================================

doctor: ## Check dev environment
	$(info "Checking environment...")
	@which docker > /dev/null || (echo "Docker not installed"; exit 1)
	@which go > /dev/null || (echo "Go not installed"; exit 1)
	$(success "Environment looks good")

format: ## Format backend (go fmt)
	cd $(BACKEND_DIR) && go fmt ./...

format-frontend: ## Format frontend (prettier)
	cd $(FRONTEND_DIR) && npm run format

tidy: ## Tidy backend go modules
	cd $(BACKEND_DIR) && go mod tidy

update-ranks: ## Update the Tranco top-1M list (server/assets/top-1m.csv)
	./scripts/update-tranco.sh

test: ## Run backend tests
	cd $(BACKEND_DIR) && go test ./...

lint: ## Lint backend (go vet)
	cd $(BACKEND_DIR) && go vet ./...

check-frontend: ## Type-check frontend (svelte-check)
	cd $(FRONTEND_DIR) && npm run check

test-frontend: ## Run frontend tests (vitest)
	cd $(FRONTEND_DIR) && npm run test

ci: tidy format format-frontend lint check-frontend test test-frontend ## Run all CI checks locally (backend + frontend)


# ============================================
# Cleanup
# ============================================

clean: ## Stop all containers + prune
	$(info "Cleaning Docker resources...")
	$(compose-dev) down -v || true
	$(compose-prod) down -v || true
	docker system prune -f
	$(success "Cleanup complete")


# ============================================
# PHONY
# ============================================
.PHONY: $(shell grep -E '^[a-zA-Z_-]+:' Makefile | sed 's/:.*//')
