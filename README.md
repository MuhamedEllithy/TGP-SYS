# Facility Management Backend

Phase 0 foundation for a multi-project facility management modular monolith, initially The Gate Plaza. **Backend only; no business modules or authentication are implemented.** The [master plan](ASSET_MANAGEMENT_BACKEND_MASTER_PLAN.md) defines the functional baseline. See the [Phase 0 report](docs/phases/phase-00-report.md) for actual verification status and blockers.

## Requirements and setup

Go **1.26.10**, Docker Engine + Compose v2+, Node 24.21.0/npm 11.19.0 for contract tooling, Python 3 for network smoke tests, and outbound dependency registry access. Go can automatically select the pinned toolchain. Module path `gateplaza.local/facility` is an internal placeholder; select the real repository import path before distributing Go packages.

```sh
make setup
# Edit .env: replace BOTH occurrences of the local password.
# Use URL-safe characters (e.g. openssl rand -hex 24); never commit .env.
make db-up
make migrate-up
make run
```

`make setup` copies `.env.example` only if `.env` does not exist, then resolves Go modules and installs documentation tools. **Bootstrap currently requires Go dependency-lock completion on an unrestricted host**: `go mod tidy` must generate `go.sum`. The npm lockfile is present and Redocly lint passes. Go dependencies could not be downloaded in the restricted implementation environment; Phase 0 is not accepted until `go.sum` is generated, reviewed and checked in and the full gate is rerun. Do not use this scaffold for production.

Compose binds PostgreSQL only to `127.0.0.1:55432`; data persists in a named volume. `make db-stop` stops it without removing data. Migration credentials are local-only; production must split owner/migrator/application roles. Password changes do not update an already initialized PostgreSQL volume: use an explicit SQL rotation procedure, never delete a volume containing data to fix credentials.

## Implemented HTTP surface

| GET or HEAD | Behavior |
|---|---|
| `/health/live` | 200 when HTTP is serving; independent of database availability |
| `/health/ready` | Real pgx pool Ping; 200 or redacted 503 within `DB_TIMEOUT`; does not check future domain migrations |
| `/openapi.yaml` | Embedded contract, local docs mode only |
| `/docs` | Scalar shell, local docs mode only; pinned external renderer requires network |

Unknown routes return JSON 404; known routes reject writes with JSON 405. All JSON responses carry a generated `meta.request_id` matching `X-Request-ID`. No CORS is enabled. Business APIs will live under `/api/v1` in later phases. The server stays live during DB outages and closes gracefully on SIGINT/SIGTERM. It never runs migrations automatically.

```sh
curl -i http://127.0.0.1:8080/health/live
curl -i http://127.0.0.1:8080/health/ready
make docs
make api-smoke
```

Open `/docs` and use Scalar's request control on liveness/readiness. This is an unverified manual procedure until dependencies and network are available. Bruno request files are in `bruno/`; open the folder in Bruno. They are manual conveniences, not another contract. Auth testing is deferred to Phase 2.

## Configuration

The process reads environment variables only; `make run` explicitly sources trusted `.env` shell syntax. Do not source an untrusted file. Never print `.env` or pass credentials into logs.

| Variable | Default / rule |
|---|---|
| `DATABASE_URL` | Required PostgreSQL URL including user, host, database; local SSL disabled only for local DB |
| `APP_ENV` | `local`, `test`, or `production`; defaults local |
| `HTTP_ADDR` | `127.0.0.1:8080`; local environment requires loopback |
| `DOCS_ENABLED` | false; true allowed only in local environment |
| `LOG_LEVEL` | info; debug/info/warn/error accepted by slog |
| `DB_TIMEOUT` | 2s; positive, maximum 1m |
| `SHUTDOWN_TIMEOUT` | 10s; positive, maximum 1m |
| `POSTGRES_PASSWORD` | Required by local Compose; not read by the API |
| `TEST_DATABASE_URL` | Explicit local disposable database for integration tests |

Production mode requires `sslmode=verify-full`, but production infrastructure, TLS termination, least privilege and operational hardening are not delivered or approved.

## Commands

```sh
make fmt
make lint                  # gofmt check + go vet (stdlib baseline static lint)
make test                  # all Go tests with race detector
make openapi-lint          # Redocly OpenAPI 3.1 validation
make migrate-up
export TEST_DATABASE_URL='postgres://facility:YOUR_LOCAL_PASSWORD@127.0.0.1:55432/facility?sslmode=disable'
make test-integration
make migrate-down          # one migration, ONLY on disposable Phase 0 DB
make migrate-up
make api-smoke             # running API + real PostgreSQL required
```

Migration round-trip details: [database runbook](docs/database/local.md). Architecture and directory boundaries: [architecture](docs/architecture/overview.md). Decisions and alternatives: [ADRs](docs/decisions/README.md). Requirements: [traceability](docs/requirements/traceability.md). No Git metadata is initialized in this environment; repository hosting/branch/remote and license are owner decisions.
# TGP-SYS
