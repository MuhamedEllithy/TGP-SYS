# Dedicated TGP-Sys local PostgreSQL

Current configuration: [ADR 0011](../decisions/0011-dedicated-postgresql17-local.md). PostgreSQL 17.11 is isolated from TheBookseller; its Supabase container/port 54322 must never be used, restarted, stopped or modified by these commands.

| Resource | Dedicated value |
|---|---|
| Compose project | `tgp-sys-local-dev` |
| Service / container | `postgres` / `tgp-sys-postgres-dev` |
| Official image | `postgres:17.11-alpine` |
| Development DB | `tgp_facility_dev` |
| Disposable test DB | `tgp_facility_phase0_test` — never store user data here |
| Runtime role | `tgp_app_dev` — non-superuser, CONNECT + public schema USAGE only |
| Migration owner | `tgp_migrator_dev` — non-superuser, owns only dedicated databases |
| Bootstrap role | `tgp_bootstrap_dev` — local first-boot administrator, never runtime |
| Volume / network | `tgp-sys-postgres17-dev-data` / `tgp-sys-dev-network` |
| Host endpoint | `127.0.0.1:55432` candidate; host preflight required |
| Internal migration endpoint | `postgres:5432` on the dedicated network |

## Exact commands for the owner's WSL terminal

The agent sandbox cannot access Docker or open network sockets. The following commands are **host procedures, not completed host verification**. They are authorized to provision only TGP-Sys resources; no sudo or sandbox bypass is used.

```sh
cd ~/Projects/TGP-Sys/backend
# Read-only daemon check (never print full container environments):
docker context show
docker version --format '{{.Client.Version}} {{if .Server}}{{.Server.Version}}{{end}}'
make env
# First provisioning only: require no existing TGP-Sys container, choose a free host port.
python3 scripts/local_env.py --select-port
make db-preflight
# Starts only the dedicated service, preserving existing containers and volumes:
make db-up
# Safe health/ownership inspection, no credential output:
docker inspect tgp-sys-postgres-dev --format '{{.State.Health.Status}}'
docker port tgp-sys-postgres-dev 5432/tcp
# Full reproducible Phase 0 checks, stopping at the first failure:
make verify-phase0
```

For repeat verification when the TGP-Sys container already exists, omit `--select-port` and run `make verify-phase0` directly. Preflight confirms ownership and port mapping before Compose starts with `--no-recreate`. It refuses resource-name collisions and unavailable Docker access. If an existing project container uses a different port/configuration, investigate explicitly; do not recreate it automatically.

`make verify-phase0` resolves Go modules with `go mod tidy`, runs `go mod download`/`verify`, installs locked npm tooling, runs normal/race Go tests, Python credential-handling tests, gofmt/vet and OpenAPI lint. It verifies real host TCP application authentication, migrates development DB forward, then performs all rollback checks only on the separate disposable DB. Finally it starts the actual Go binary on the default HTTP address, exercises probes/docs/contract availability and verifies graceful SIGTERM shutdown. It leaves database containers and volumes intact. The first unresolved check exits nonzero. No report is automatically marked accepted from a partial run.

## Environment handling

`make env` generates three independent 32-byte cryptographic passwords directly into `.env` (0600), with no secret output. It refuses tracked `.env`, unconfirmed Git exclusion, symlink files, unrelated preexisting DB settings, placeholder secrets or inconsistent existing URLs. It preserves unrelated assignments/comments and existing credentials. Python tooling parses simple KEY=value syntax without executing shell commands. Do not run `cat .env`, shell xtrace, full `docker inspect`, or plain `docker compose config` when real secrets are configured; use `config --quiet` only.

Explicit fields: `DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_USER`, `DB_PASSWORD`, `POSTGRES_USER`, `POSTGRES_PASSWORD`, `MIGRATION_USER`, `MIGRATION_PASSWORD`, `TEST_DB_NAME`. The Go API consumes **DATABASE_URL** only; all four derived URL fields use encoded userinfo and consistent settings. Migrations use `MIGRATION_DATABASE_URL`; test checks use host `TEST_DATABASE_URL`; test migrations use internal `TEST_MIGRATION_DATABASE_URL`. `.env.example` has matching names with placeholders only. Application runtime never uses bootstrap credentials.

First-boot initialization creates roles/databases and SCRAM authentication only on the newly named volume. It does not reset existing passwords or run against a previously initialized volume. If initialization was interrupted, do not delete the volume or assume changing `.env` changes stored passwords. Investigate with the owner, preserving data.

## Individual checks and safe migrations

```sh
make setup                 # generate/preserve .env, resolve real dependency locks
make db-verify             # pgx: exact host path, DB/user/version/privileges + wrong-password rejection
make migrate-up            # forward-only development DB
make migrate-test-up       # guarded disposable database only
make test-integration
make migrate-test-down     # down 1 on disposable test DB only
python3 scripts/local_env.py --run go run ./cmd/dbcheck --test --probe-absent
make migrate-test-up
make test-integration
make run                   # separate terminal for interactive docs
make api-smoke
```

`make migrate-down` is an alias for **test rollback only**, not development rollback. The test guard checks actual DB identity and rejects unexpected user objects before migrations. Rollback removes only migration 000001's infrastructure probe. Application schema/data is not created in Phase 0. Domain migrations, runtime table grants and finer role policies require Phase 1 design approval. Never run down/reset/force or remove volumes on a database containing user data.

Browse `http://127.0.0.1:8080/docs` and use Scalar Try-it on both probes. HTTP docs availability is automatically checked; browser renderer functionality requires its pinned CDN/network and remains a separate manual check. `pg_isready` and Docker healthy status do not prove the password; accept only the actual host-path pgx check. Do not paste credentials or full driver errors into chat.

## Current verification status

The owner reports healthy PostgreSQL 17.11, successful authenticated host TCP access, incorrect-password rejection, non-superuser runtime, disposable migration round trip, integration tests, HTTP endpoints and graceful shutdown. These are host-executed results, not an independent sandbox infrastructure replay. The full Go race suite/vet and module verification now also pass directly in the sandbox using available dependencies.

Clean-clone bootstrap evidence, manual Scalar Browser Try-it and owner acceptance remain pending. HTTP `/docs` availability is not browser request testing. See the [final Phase 0 report](../phases/phase-00-report.md) for PASS/PENDING provenance and risks. Phase 1 is not started.
