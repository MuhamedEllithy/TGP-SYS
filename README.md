# Facility Management Backend

Phase 0 Go modular-monolith foundation for TGP-Sys, initially The Gate Plaza. Backend only; no business modules or authentication are implemented. The [master plan](ASSET_MANAGEMENT_BACKEND_MASTER_PLAN.md) is authoritative. [Phase 0 report](docs/phases/phase-00-report.md) records actual passes and blockers; **PHASE 0 — ACCEPTED AND CLOSED** for implementation `46555ea1a4f3ff14cd24354e1ec49c14e3ba74ec`. Actual-commit clean-checkout verification passed; owner WSL evidence confirms database/live Bruno success, and the owner confirmed manual Scalar Try-it. See the [acceptance record](docs/phases/phase-00-acceptance.md) for provenance and non-blocking debt. Phase 1A is not started and requires separate owner authorization.

## Local setup (owner's WSL terminal)

Requires Go 1.26.10 (automatic toolchain selection), Docker/Compose, Node 24.21.0/npm 11.19.0 and Python 3. No business UI or deployment service is built.

```sh
cd ~/Projects/TGP-Sys/backend
make env
# Before the FIRST container creation only:
python3 scripts/local_env.py --select-port
make db-up
make verify-phase0
```

The [database runbook](docs/database/local.md) gives exact commands, ownership checks, migration safety and manual Scalar verification. The owner reports successful host PostgreSQL/API verification. Earlier sandbox Docker restrictions are historical execution limits; see the report for evidence provenance and pending gates.

Dedicated official PostgreSQL 17.11: container `tgp-sys-postgres-dev`, DB `tgp_facility_dev`, non-superuser runtime `tgp_app_dev`, separate migrator and bootstrap roles, named volume and private project network. Host binding is loopback only (candidate port 55432, selected/validated on host). The unrelated `supabase_db_thebookseller` and port 54322 are never used.

`make env` creates ignored mode-0600 `.env` with three strong random passwords and encoded URLs, without displaying values. It preserves existing credentials/unrelated settings and refuses inconsistent or tracked files. `.env.example` contains placeholders only; do not copy it over real credentials. No credentials are committed. `go.sum` is now present after host dependency resolution; `make setup` runs tidy/download/verify and installs locked npm dependencies. Do not fabricate checksums or treat missing dependencies as passed checks.

## Implemented HTTP surface

| GET/HEAD | Meaning |
|---|---|
| `/health/live` | HTTP process live, independent of DB |
| `/health/ready` | Actual bounded pgx Ping, 200 or redacted 503; no domain schema-version assertion yet |
| `/docs` | Local-only Scalar reference shell, pinned CDN renderer |
| `/openapi.yaml` | Local-only embedded authoritative OpenAPI 3.1.1 |

Unknown routes return JSON 404; unsupported methods return 405. JSON responses include a generated `meta.request_id` matching `X-Request-ID`. Business routes, scoped authorization and user authentication are deferred. No CORS is enabled. API never auto-migrates; SIGTERM/SIGINT drain requests before closing the pool.

## Configuration

The API consumes environment variables. `make run` loads `.env` without executing shell code. Core names: required `DATABASE_URL`; `APP_ENV` (local/test/production, default local); `HTTP_ADDR` (default 127.0.0.1:8080, local requires loopback); `DOCS_ENABLED` (false by default, true only locally); `LOG_LEVEL` (info); `DB_TIMEOUT` (2s); `SHUTDOWN_TIMEOUT` (10s). Durations must be positive and at most 1m. Production URL requires `sslmode=verify-full`, but no production infrastructure/deployment is delivered.

Compose/bootstrap connection fields and all migration/test URL names are listed in [.env.example](.env.example) and the runbook. Application, migration and bootstrap secrets differ. Runtime has CONNECT/schema USAGE only; future table grants require actual schema review.

## Development commands

```sh
make setup
make fmt
make lint                  # gofmt check + go vet
make test                  # Go race suite + Python environment safety tests
make openapi-lint
make db-verify             # host-path password authentication, identity/version/privileges
make migrate-up            # development DB forward only
make migrate-test-up
make test-integration
make migrate-test-down     # disposable test DB ONLY; never development rollback
make migrate-test-up
make run
make api-smoke
make docs
make verify-phase0         # all host acceptance checks; fails immediately on blockers
make db-stop               # only dedicated service, no data/volume deletion
```

The owner completed real Bruno execution (4/4 requests, 7/7 tests, 10/10 assertions) and manually confirmed Scalar Try-it GET `/health/live`. These are owner-executed results, not claims derived from handler tests or HTML availability. Bruno probe files in `bruno/` contain no secrets. Local evidence and blockers are in the phase report; no Phase 1 work may begin without explicit approval.

[Architecture](docs/architecture/overview.md), [ADRs](docs/decisions/README.md), [traceability](docs/requirements/traceability.md) and [repository instructions](AGENTS.md) describe module boundaries and future gates. Go module path `gateplaza.local/facility` remains provisional; production hosting, backups, private file storage and repository license are undecided.

Final acceptance guides: [Bruno CLI verification](bruno/README.md), [manual Scalar browser checklist](docs/phases/phase-00-browser-checklist.md), and [architecture/security review](docs/architecture/phase-00-acceptance-review.md). No Phase 1 authorization is implied.
