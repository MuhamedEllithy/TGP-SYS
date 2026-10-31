SHELL := /bin/bash
export GOTOOLCHAIN := go1.26.10
export GOCACHE ?= /tmp/facility-go-cache
COMPOSE = docker compose --env-file .env -f deployments/local/compose.yaml
.PHONY: env setup db-preflight db-up db-stop db-verify migrate-up migrate-down migrate-test-up migrate-test-down run test test-integration lint fmt openapi-lint api-smoke docs check verify-phase0
env:
	@python3 scripts/local_env.py
setup: env
	go mod tidy
	go mod download
	go mod verify
	cd tools && npm ci --ignore-scripts
db-preflight:
	@python3 scripts/db_preflight.py
db-up: db-preflight
	@$(COMPOSE) up -d --wait --no-recreate postgres
db-stop:
	@$(COMPOSE) stop postgres
db-verify:
	@python3 scripts/local_env.py --run go run ./cmd/dbcheck
migrate-up:
	@python3 scripts/migrate.py up
# Never roll back the persistent development database using a convenience target.
migrate-down: migrate-test-down
migrate-test-up:
	@python3 scripts/migrate.py up --test
migrate-test-down:
	@python3 scripts/migrate.py down --test
run:
	@python3 scripts/local_env.py --run go run ./cmd/api
test:
	go test -race ./...
	python3 -m unittest discover -s scripts -p 'test_*.py'
test-integration:
	@python3 scripts/local_env.py --run go test -tags=integration -count=1 ./tests/integration
fmt:
	gofmt -w cmd internal api tests
lint:
	@test -z "$$(gofmt -l cmd internal api tests)" || (echo 'Run make fmt'; exit 1)
	go vet ./...
openapi-lint:
	cd tools && npm run lint
api-smoke:
	python3 scripts/api_smoke.py
docs:
	@echo 'Run make run, then open http://127.0.0.1:8080/docs (DOCS_ENABLED=true, APP_ENV=local).'
check: lint test openapi-lint
verify-phase0:
	@bash scripts/verify_phase0.sh
