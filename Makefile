# Arthur: everything you need is a make command. Run `make` to see them.

COMPOSE := docker compose

# Passed to docker-compose.yml so the worker container runs as you and can use Docker.
export LOCAL_UID := $(shell id -u)
export LOCAL_GID := $(shell id -g)
export DOCKER_GID := $(shell stat -c %g /var/run/docker.sock 2>/dev/null || stat -f %g /var/run/docker.sock 2>/dev/null || echo 999)

.DEFAULT_GOAL := help

.PHONY: help
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "} {printf "  \033[36m%-14s\033[0m %s\n", $$1, $$2}'

.PHONY: check-env
check-env:
	@test -f backend/.env || { echo "backend/.env is missing. Run: cp backend/.env.example backend/.env  (then fill in the values)"; exit 1; }

.PHONY: prepare
prepare: check-env
	@mkdir -p /tmp/project-arthur
	@chmod 1777 /tmp/project-arthur 2>/dev/null || true
	@test -w /tmp/project-arthur || { echo "/tmp/project-arthur is not writable by you (probably created earlier with sudo). Run: sudo chmod 1777 /tmp/project-arthur"; exit 1; }

.PHONY: up
up: prepare ## Build and start everything (db, migrations, api, worker, frontend)
	$(COMPOSE) up --build -d
	@echo ""
	@echo "  Dashboard:  http://localhost:5173"
	@echo "  API:        http://localhost:8080/v1/healthz"
	@echo "  Postgres:   localhost:5433  (user/password/db: arthur)"
	@echo "  Logs:       make logs"

.PHONY: up-lite
up-lite: prepare ## Start everything EXCEPT the worker (then run `make worker-local`)
	$(COMPOSE) up --build -d db migrate api frontend

.PHONY: worker-local
worker-local: check-env ## Run the worker on your machine instead of in Docker (use with up-lite)
	cd backend && DB_URL="postgres://arthur:arthur@localhost:5433/arthur?sslmode=disable" go run ./cmd/worker

.PHONY: down
down: ## Stop everything (keeps the database)
	$(COMPOSE) down

.PHONY: restart
restart: ## Restart all services without rebuilding
	$(COMPOSE) restart

.PHONY: rebuild
rebuild: prepare ## Rebuild images from scratch (use after changing package.json / go.mod)
	$(COMPOSE) build --no-cache
	$(COMPOSE) up -d -V

.PHONY: logs
logs: ## Follow logs. One service: make logs s=worker
	$(COMPOSE) logs -f --tail=100 $(s)

.PHONY: ps
ps: ## Show service status
	$(COMPOSE) ps

.PHONY: migrate
migrate: ## Re-run database migrations
	$(COMPOSE) run --rm migrate

.PHONY: psql
psql: ## Open a psql shell in the database
	$(COMPOSE) exec db psql -U arthur -d arthur

.PHONY: test
test: ## Run the backend tests
	cd backend && go test ./...

.PHONY: sqlc
sqlc: ## Regenerate database code after editing sql/
	cd backend && sqlc generate

.PHONY: smoke
smoke: ## End-to-end test against the running stack: make smoke REPO=https://github.com/you/repo
	@test -n "$(REPO)" || { echo "usage: make smoke REPO=https://github.com/<you>/<npm-buildable-repo>"; exit 1; }
	REPO=$(REPO) ./scripts/smoke_test.sh

.PHONY: clean
clean: ## Stop everything AND delete the database volume (asks first)
	@read -p "This deletes ALL local data (users, projects, deployments). Continue? [y/N] " a; \
	if [ "$$a" = "y" ] || [ "$$a" = "Y" ]; then $(COMPOSE) down -v; else echo "aborted"; fi
