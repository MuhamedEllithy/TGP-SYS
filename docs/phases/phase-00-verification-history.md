# Historical Phase 0 verification snapshots

Archived before incorporating the owner's successful host verification. Earlier BLOCKED/FAIL statements describe the restricted execution environment at those times; they are not the current acceptance assessment. Old configuration/commands may be superseded. Use the [current report](phase-00-report.md), README and database runbook.

---

# Phase 0 — Delivery report

Date: 2026-10-31  
Status: **PARTIALLY COMPLETE — ACCEPTANCE BLOCKED**  
Commit/branch: unavailable; the supplied `.git` directory is not a valid repository. No Git initialization, commit, push or deployment was performed.  
Next phase: **NOT STARTED; not authorized.**

**Current state: see §15, Dedicated PostgreSQL follow-up.** Sections 1–14 retain historical evidence from earlier attempts; their old PostgreSQL version, shared-role setup, Git availability and missing `.env` statements are superseded by §15 and ADR 0011. Use the current README/database runbook for commands.


## 1. Executive summary

Read the master plan and inspected the initially empty workspace before creating the backend foundation. No existing application code, dependencies, conventions or AGENTS instructions were present. The master plan is unchanged. No functional contradictions require a business-scope change. All consequential future choices are ADR recommendations, not implicit approvals.

Created a small modular Go backend scaffold, environment validation, JSON structured logging, bounded HTTP transport, genuine pgx readiness adapter, graceful shutdown, local PostgreSQL Compose setup, paired SQL probe migration, OpenAPI 3.1.1 contract and local Scalar shell. Created onboarding, agent rules, requirement traceability, a real-DB integration test, Bruno probe requests and proposed CI. No business tables, authentication, asset modules, UI or production infrastructure were implemented.

**Verified:** formatting, config and HTTP unit/negative tests with race detector, targeted vet, OpenAPI lint, npm clean install from lock, Compose configuration parsing and source secret-pattern checks. **Not verified:** full build/test/vet, PostgreSQL connectivity, migration round trip, actual process startup/shutdown, real network smoke, browser Scalar Try-it or hosted CI. Missing pgx resolution, Docker denial and blocked network sockets prevent those gates. No infrastructure was mocked to declare acceptance.

## 2. Architecture decisions

See [ADR index](../decisions/README.md):

| ADR | Scope / current decision |
|---|---|
| 0001 | Go 1.26.10; standard net/http and slog; modular monolith |
| 0002 | PostgreSQL 18.6 local container; no extensions or production assumptions |
| 0003 | pgx 5.11.0; golang-migrate 4.20.1; defer sqlc until actual SQL exists |
| 0004 | Recommend UUIDv7, separate identities, home/current project and historical code claims; confirm before DDL |
| 0005 | Recommend revocable opaque sessions, CSRF protection and scoped default deny; Phase 2 gate |
| 0006 | Contract-first OpenAPI 3.1.1; pinned Scalar local reference; Redocly validation |
| 0007 | Future private storage port and quarantined scanning; provider undecided |
| 0008 | Initial test/CI gates; missing Go lock deliberately blocks CI |
| 0009 | Immutable referenced status/movement revisions and safe implemented primitives |
| 0010 | Transaction/idempotency boundaries, evidence-only observations and unresolved operating limits |

[Architecture](../architecture/overview.md) documents inward dependencies and future module ownership. [Traceability](../requirements/traceability.md) assigns IDs to all baseline requirement groups, links proposed tables/services/tests and marks future work deferred.

## 3. Repository structure

`cmd/api` composes config, DB and HTTP adapters. `internal/platform` contains only config/db/httpx infrastructure. `api` embeds the contract and docs shell. `db/migrations` has only an infrastructure probe. `deployments/local` runs PostgreSQL and explicit migration tooling. `tests/integration`, `scripts` and `bruno` provide testing foundations. `tools` holds pinned npm dependencies and lockfile. `docs` holds architecture, ADRs, runbooks, traceability and this report. `.github/workflows` contains an unexecuted CI proposal.

## 4. Files created / modified

All following files were newly created; the master plan was preserved. Ignored `tools/node_modules` is local installed tooling, not deliverable source. No `.env` with credentials was created. Temporary authoring helpers are outside the workspace in `/tmp`.

- `.env.example`
- `.github/workflows/ci.yml`
- `.gitignore`
- `AGENTS.md`
- `Makefile`
- `README.md`
- `api/docs.html`
- `api/embed.go`
- `api/openapi.yaml`
- `bruno/bruno.json`
- `bruno/health-live.bru`
- `bruno/health-ready.bru`
- `cmd/api/main.go`
- `db/migrations/000001_platform_probe.down.sql`
- `db/migrations/000001_platform_probe.up.sql`
- `deployments/local/compose.yaml`
- `docs/architecture/overview.md`
- `docs/database/local.md`
- `docs/decisions/0001-go-http-modular-monolith.md`
- `docs/decisions/0002-postgresql.md`
- `docs/decisions/0003-db-migrations.md`
- `docs/decisions/0004-identity-codes.md`
- `docs/decisions/0005-authentication-authorization.md`
- `docs/decisions/0006-openapi-docs.md`
- `docs/decisions/0007-private-storage.md`
- `docs/decisions/0008-ci-testing.md`
- `docs/decisions/0009-config-status-movements.md`
- `docs/decisions/0010-consistency-open-decisions.md`
- `docs/decisions/README.md`
- `docs/phases/phase-00-report.md`
- `docs/requirements/traceability.md`
- `go.mod`
- `internal/platform/config/config.go`
- `internal/platform/config/config_test.go`
- `internal/platform/db/db.go`
- `internal/platform/httpx/server.go`
- `internal/platform/httpx/server_test.go`
- `redocly.yaml`
- `scripts/api_smoke.py`
- `tests/integration/db_test.go`
- `tools/package-lock.json`
- `tools/package.json`

## 5. Dependencies and versions

| Dependency | Version / verification |
|---|---|
| Go | 1.26.10 cached toolchain; module language 1.26.0; machine default 1.24.0 is not used for project checks |
| pgx / pgxpool | v5.11.0 pinned, MIT; unavailable locally, download blocked, **no go.sum generated** |
| PostgreSQL | 18.6-alpine exact patch tag proposed; image pull/runtime blocked |
| golang-migrate | v4.20.1 image tag; migration execution blocked |
| Scalar API Reference | 1.76.4 installed and locked in tools; shell references exact CDN version; browser untested |
| Redocly CLI | 2.72.1 installed and locked; real lint passes |
| Node / npm | 24.21.0 / 11.19.0 observed locally |
| Docker / Compose | 29.8.2 / v5.6.0 observed; daemon socket denied |
| Python | 3.14 observed for smoke execution; no third-party runtime requirement for smoke |

The initial candidate Redocly 2.14.4 install reported dependency vulnerabilities; replaced with 2.72.1 and regenerated npm lock. The final Scalar/Redocly installation reported zero vulnerabilities, but an independent `npm audit` attempt failed on advisory-endpoint DNS access. This is not a completed dependency-security review. pgx upstream master mentions a SendBatch security fix after the selected release: current code only opens a pool/Pings, never uses SendBatch; review fixed releases and advisories before domain SQL. Direct licenses reviewed from primary upstream sources; transitive inventory is incomplete until Go resolution. Image/action tags are exact versions but not immutable digest/commit pins.

## 6. Database and migration setup

Compose binds disposable PostgreSQL to loopback port 55432 and persists its named volume. Password is required externally; example values are explicitly placeholders. Use `make db-up`, `make migrate-up`, `make test-integration`, `make migrate-down`, then up/test again per [database runbook](../database/local.md). Migration 000001 creates only `platform_probe` with a single true-key row and timestamp. Both directions use explicit transactions. Migration tooling owns `schema_migrations`; no hand-built engine or business schema introduced.

API never auto-migrates. Readiness uses a real bounded pgx Ping but does not yet assert domain schema compatibility. Local DB bootstrap user is intentionally an owner; production requires separate least-privileged application and migrator identities. Down removes only the disposable probe, never a domain/history table. No reset, volume deletion, production migration or database mutation actually ran here.

## 7. API and documentation setup

Implemented source routes: GET/HEAD `/health/live`, `/health/ready`, local-only `/docs` and `/openapi.yaml`. Unknown routes return JSON 404; known routes reject other methods with JSON 405. Readiness is redacted 503 on DB failure/timeout. JSON envelopes include a server-generated request ID; request-provided IDs are ignored. OpenAPI documents only these operations, examples, errors, headers and permissions exceptions. Contract is embedded to avoid runtime working-directory drift.

Scalar is permitted only in `APP_ENV=local`, with loopback binding; production docs enablement is rejected. Renderer is an exact-version external CDN dependency and requires network; shell/CSP use same-origin request permissions and disable auth persistence. UI rendering and Try-it remain untested. Bruno probe requests have no credentials; later authentication/authorization flows belong to Phase 2. Optional Hurl was not installed; a Python real-network smoke script supplies the foundation command without another HTTP-test dependency.

## 8. Exact verification commands and results

Commands ran from workspace root unless specified. No failing full-project command is treated as a pass because targeted packages passed.

| Command | Result |
|---|---|
| `gofmt -w cmd internal api tests` | PASS; formatting applied |
| `test -z "$(gofmt -l cmd internal api tests)"` | PASS; no unformatted Go files |
| `GOTOOLCHAIN=go1.26.10 GOCACHE=/tmp/facility-go-cache go test -race -count=1 ./internal/platform/config ./internal/platform/httpx ./api` | PASS; config/httpx pass, api embeds have no separate test files |
| `GOTOOLCHAIN=go1.26.10 GOCACHE=/tmp/facility-go-cache go vet ./internal/platform/config ./internal/platform/httpx ./api` | PASS |
| `make openapi-lint` | PASS; Redocly: API description valid, no warnings/errors |
| `npm ci --ignore-scripts --fetch-retries=0 --fetch-timeout=10000 --cache /tmp/facility-npm-cache` in `tools` | PASS; 284 installed packages from actual lock |
| `POSTGRES_PASSWORD=placeholder docker compose -f deployments/local/compose.yaml config --quiet` | PASS; configuration parses |
| `make lint` | BLOCKED/FAIL; full go vet requires missing pgx checksum/module |
| `make test` | BLOCKED/FAIL; full race suite cannot build cmd/api and db without pgx/go.sum |
| `GOTOOLCHAIN=go1.26.10 GOCACHE=/tmp/facility-go-cache go build ./cmd/api` | BLOCKED/FAIL; missing pgx go.sum |
| `GOTOOLCHAIN=local GOMODCACHE=/tmp/facility-mod-cache GOCACHE=/tmp/facility-go-cache /home/muhamed/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.10.linux-amd64/bin/go mod tidy` | BLOCKED/FAIL; proxy.golang.org DNS/socket operation not permitted. Writable cache avoids the initial read-only default-cache failure |
| `POSTGRES_PASSWORD=placeholder docker compose -f deployments/local/compose.yaml up -d --wait postgres` | BLOCKED/FAIL; Docker daemon socket permission denied |
| `POSTGRES_PASSWORD=placeholder docker compose -f deployments/local/compose.yaml run --rm migrate` | BLOCKED/FAIL; same Docker denial; up never ran |
| Direct Compose `migrate ... down 1` invocation | BLOCKED/FAIL before execution; Docker denial; round trip remains untested |
| `TEST_DATABASE_URL='postgres://facility:placeholder@127.0.0.1:55432/facility?sslmode=disable' make test-integration` | BLOCKED/FAIL; pgx checksum unavailable, DB also unavailable |
| `make api-smoke` | BLOCKED/FAIL; creating network socket returns operation not permitted; no running full server claimed |
| `npm audit --audit-level=low --cache /tmp/facility-npm-cache` in `tools` | BLOCKED/FAIL; EAI_AGAIN registry advisory endpoint; install-time report alone is insufficient |
| `git status --short` | Unavailable: workspace is not a valid Git repository |

Additional Python inspection parsed YAML, checked component schemas with JSON Schema Draft 2020-12 tooling already on the machine, confirmed the four-route inventory and absence of `.env`: PASS. This supplements, not replaces, Redocly validation.

Negative coverage in passing tests: invalid/missing/redacted DB URLs; invalid environment, log level, boolean, bind address/port and duration; public local binding rejection; production TLS/docs restrictions; database error redaction; timeout-driven 503; unknown/future route 404; unsupported method 405; HEAD bodies; generated IDs; no query strings/driver credentials in logs; server timeout/header limits. Test Pinger doubles are unit boundaries only; real connectivity acceptance is explicitly blocked.

## 9. Security and configuration checks

- Source scan for private-key headers and common AWS/GitHub credential token patterns returned no matches. Scanned newly created source excluding installed tooling/lock content and the untouched planning document. Review also found only explicit disposable CI/example password placeholders; no actual credential or local `.env` was added.
- `.gitignore` excludes `.env`, env variants, node_modules, build outputs and logs, preserving `.env.example`. Because Git metadata is invalid, no statement about an existing committed history can be verified.
- Config/DB errors never log URLs or driver error text; no configuration dumps, payload/query/header logging or embedded bootstrap credentials.
- Timeout/header bounds, local-only docs, no CORS and production DB TLS verification are present. Auth, authorization, audit events, rate limiting and business-body/upload limits remain later-phase work; no business route is exposed.
- No production infrastructure, credentials, migrations or deployments were accessed.

## 10. Limitations and technical debt

1. **Blocking:** missing Go modules/checksums means no full compilation/startup proof. Generate and commit genuine `go.sum` on an authorized networked host; rerun build, vet, all tests and module verification.
2. **Blocking:** Docker/network restrictions prevent real DB, up/down, constraint tests and network startup/shutdown checks. No fallback/fake DB substitutes were used.
3. **Untested:** clean source-copy/clone bootstrap, hosted CI, Scalar browser rendering, real Try-it, Bruno execution and graceful process shutdown. Handler recorder tests are not network/process evidence.
4. CDN renderer needs network; self-host/checksum verification is preferable before offline use. Installed Scalar dependency records the matching distribution; the API currently references CDN rather than serving installed node_modules.
5. CI is a foundation proposal: missing go.sum fails its explicit lock gate. Actions/images need immutable SHA/digest reconciliation; staticcheck, govulncheck and a dedicated secret scanner remain hardening items. gofmt/vet are the present lint baseline.
6. Readiness checks connectivity only. Add schema compatibility once the Phase 1 migration contract is approved.
7. No application or repository license/remote/module-path ownership has been selected. No SLA/volume, hosting, private-storage provider, backups/PITR/restore validation or monitoring provider selected.
8. Full dependency security/license review remains outstanding. No claim of production readiness or completed Phase 0 acceptance is made.

## 11. Unresolved decisions / blockers

Environmental blockers require an authorized development host with PostgreSQL/Docker, network sockets and Go module registry access; no extra production permission is needed or requested. Review ADR 0004 before schema work (code normalization/home versus current scope); ADR 0005 before auth (session/CSRF/expiry/deny precedence/MFA); ADRs 0007/0009/0010 before their later phases (storage, policy revisions, idempotency, operational limits). These design recommendations are documented without silently locking unsafe business assumptions.

## 12. Phase 0 acceptance checklist

- [x] Master plan read; original functional baseline preserved; repository inspected.
- [x] Concise phase plan, architecture boundaries, ADRs, onboarding and agent instructions.
- [x] Traceability foundation with future requirements explicitly deferred.
- [x] Minimal server/config/logging/health source and negative unit tests.
- [x] Versioned local DB/migration setup and real integration-test source.
- [x] OpenAPI 3.1.1 skeleton lint passes; local API documentation tooling foundation.
- [x] Go formatting, targeted vet and targeted race tests pass.
- [x] npm dependency lock and clean tool installation verified.
- [x] Source secret review; no real credentials added.
- [ ] Go dependency lock resolution and module verification.
- [ ] Clean clone/source-copy bootstrap and full `go test ./...` / `go vet ./...` / build.
- [ ] Actual DB connectivity, migrate up/down/up and constraint verification.
- [ ] Actual server startup, probes over HTTP and graceful shutdown.
- [ ] Scalar browser Try-it and Bruno requests against actual server.
- [ ] Hosted CI execution and full dependency-security review.

## 13. Recommended next steps and stop

Keep work in Phase 0. On an unrestricted local development host, generate/review `go.sum`, run `make setup`, choose a disposable local password, then execute all README/database verification commands, including rollback table-absence assertion, repeated up, full race suite, real API smoke and manual docs requests. Update this report with actual results; resolve security advisories and immutable pins before trusted CI/release. Request owner review only after the remaining Phase 0 evidence is concrete.

**Decision for owner:** review/revise Phase 0 artifacts and documented recommendations. Phase 0 acceptance is not claimed while mandatory gates remain blocked. No approval to Phase 1 is inferred or requested as a substitute for completing these gates. **Execution stops here; Phase 1 is not started.**


## 14. WSL Docker/database inspection follow-up — 2026-10-31

**Result: infrastructure access blocked; no database connection configured or verified.** Docker being installed on the WSL host does not establish that this restricted agent can access its daemon. This follow-up performed read-only inspection only; no containers, networks, volumes, databases, credentials or production infrastructure were changed. No replacement database was launched. Phase 1 remains unstarted.

### Docker and project association

- WSL kernel: `6.6.87.2-microsoft-standard-WSL2`.
- Docker client: `29.8.2`; daemon version unavailable.
- Active context: `default`, endpoint `unix:///var/run/docker.sock`.
- Socket exists with mode `srw-rw----`, owner `nobody`, group `nogroup`. Although the process reports supplementary `nogroup` membership, daemon requests still return permission denied; these facts alone do not prove a host group-permission fix is needed. No socket chmod/chown, group change, sudo, daemon change or alternative privileged access was attempted.
- `docker version --format '{{.Client.Version}} {{if .Server}}{{.Server.Version}}{{end}}'`: client version printed; daemon access denied.
- `docker ps -a --format '{{.ID}}\t{{.Names}}\t{{.Image}}\t{{.Ports}}\t{{.Status}}'`: denied.
- `docker --host unix:///var/run/docker.sock ps -a --format '{{.ID}}\t{{.Names}}\t{{.Image}}\t{{.Ports}}\t{{.Status}}'`: denied, including the explicit host path.
- `docker compose ls --all --format json`, `docker network ls --format '{{.Name}}\t{{.Driver}}'` and `docker volume ls --format '{{.Name}}'`: all denied.

**PostgreSQL container found/not found: UNKNOWN.** Container existence, running/stopped state, names, versions, published ports, project labels, network membership and mounted volumes cannot be established. Do not interpret this as an empty Docker host or permission to create a new database. No unrelated container was attributed to this project.

The only repository Compose file found is `deployments/local/compose.yaml`. It declares project `facility-local`, PostgreSQL image `postgres:18.6-alpine`, a loopback mapping `127.0.0.1:55432:5432`, service hostname `postgres`, configured database/user `facility`, and persistent volume `postgres_data`. These are **source declarations only**, not evidence of an existing container, installed PostgreSQL version, initialized database or valid credentials. Existing container association must be confirmed from Compose labels/working-directory/config paths and mounts when daemon access is available, without printing container environment values.

### Secure environment inspection

No `.env` exists. The repository contains `.env.example` with explicit placeholder credentials only. The backend's actual required connection variable is `DATABASE_URL`; optional runtime variables are documented in README and validated by `internal/platform/config`. Local Compose separately expects `POSTGRES_PASSWORD`. The existing local migration service uses its internal Compose path `postgres:5432`, whereas the declared host application path is `127.0.0.1:55432`; both currently declare user/database `facility`. No actual endpoint or credential pair was verified.

`.gitignore` excludes `.env` and `.env.*` while retaining `.env.example`. `git ls-files --error-unmatch .env` cannot determine tracking status because `.git` is not a valid repository. No Git history or index was altered. No passwords, tokens, full container environment dumps, or connection URLs with real credentials were printed.

**Environment configuration: intentionally unchanged and blocked.** Creating `.env` with guessed or newly generated credentials would not configure a verified existing database. No new password was generated or applied, and no unverified connection was represented as working. Existing example/Compose files were preserved. When the real project-owned connection is verified, populate ignored `.env` with the requested explicit host/port/name/user/password fields plus the required `DATABASE_URL`, preserving other values. Construct URL userinfo with proper percent encoding rather than interpolating a raw password; validate with pgx over the exact host/port used by the application. Container `POSTGRES_PASSWORD` and `pg_isready` alone are insufficient authentication evidence for an existing volume.

### Connectivity, migrations and startup

- Actual password-authenticated PostgreSQL connection: **BLOCKED / NOT ATTEMPTED**, because neither project association nor credentials/endpoint can be established.
- Actual PostgreSQL readiness and server version: **UNKNOWN / NOT VERIFIED**.
- `psql` executable: not found in the agent PATH.
- Disposable database identity: **NOT ESTABLISHED**. No database was created or selected for destructive migration tests.
- Migration up/down/reapplication: **NOT RUN in this follow-up**; existing-data safety requires a confirmed disposable target first. Earlier failed attempts remain recorded above.
- Full Go application startup, network health and readiness: **NOT RUN in this follow-up**, conditional on a verified database connection. The previous unresolved pgx/go.sum blocker also remains; no `go.sum` was found.
- `GOTOOLCHAIN=go1.26.10 GOCACHE=/tmp/facility-go-cache go test -race -count=1 ./internal/platform/config ./internal/platform/httpx ./api`: **PASS** (config and httpx; api has no separate test files). Unit boundary doubles are not evidence of actual PostgreSQL authentication.

### Owner checkpoint and safe continuation

Stop at this unresolved infrastructure blocker as requested. Resume Phase 0 in an execution environment authorized to access the existing WSL Docker daemon. First repeat the read-only inventory commands, confirm project association, then inspect only necessary credential sources without outputting secret values and prove password authentication through the application's host path. If that inventory establishes that no suitable project container exists, prepare the requested isolated, environment-driven proposal with a dedicated non-superuser application role, strong randomly generated local credentials stored only in an ignored mode-0600 file, persistent named volume, loopback port and health check; do not launch it automatically. The existing scaffold remains a proposal and must not be mistaken for a verified running database or the requested least-privileged role setup.

Only after a verified connection and explicitly disposable test database are available should the migration round trip and real startup/probe tests run. No Phase 0 acceptance, Phase 1 authorization, deployment or remote push is claimed by this follow-up.


## 15. Dedicated TGP-Sys PostgreSQL and Phase 0 verification follow-up

Date: 2026-10-31  
**Final recommendation: keep Phase 0 open — configuration prepared, mandatory live acceptance checks BLOCKED.** No Phase 1, deployment, push, container recreation, volume deletion or unrelated project mutation occurred.

### Execution environment and owner-supplied host facts

The owner confirms Docker works without sudo in the WSL host terminal under `muhamed`, context `default`, and that the host user belongs to the Docker group. The owner identifies healthy `supabase_db_thebookseller`, image `public.ecr.aws/supabase/postgres:17.11.0.006`, port 54322 as an unrelated application. Those are **owner-supplied host facts**, not fresh agent runtime results. It was not used, stopped, restarted or modified.

Current sandbox diagnostics: Docker client 29.8.2 and context `default` are available, but Docker version's daemon query and container listing return permission denied on `/var/run/docker.sock`. `make db-preflight` cannot obtain daemon information, and `make db-up` stops before provisioning. Socket binds also return operation not permitted. No sudo, socket-permission changes, alternate host socket or other sandbox bypass was attempted. A working host daemon and denied agent access are compatible facts.

### Dedicated PostgreSQL configuration and decisions

[ADR 0011](../decisions/0011-dedicated-postgresql17-local.md) supersedes the initial PostgreSQL 18 local selection with owner-authorized PostgreSQL **17.11**, official image `postgres:17.11-alpine`. PostgreSQL 17 remains supported according to the reviewed primary version policy. Existing probe SQL uses compatible PostgreSQL features; pgx and golang-migrate choices remain unchanged. Image pull, first-boot SQL and actual runtime compatibility remain unverified.

| Setting | Prepared value / status |
|---|---|
| Compose project | `tgp-sys-local-dev` |
| Service / container | `postgres` / `tgp-sys-postgres-dev` — **not launched here** |
| Network | `tgp-sys-dev-network`, dedicated, no external project dependency |
| Persistent volume | `tgp-sys-postgres17-dev-data`, new name, PostgreSQL 17 data-directory layout |
| Host binding | `127.0.0.1:55432` candidate; availability **not verified in sandbox** |
| Internal endpoint | `postgres:5432`, dedicated Compose network |
| Development database | `tgp_facility_dev` |
| Runtime login | `tgp_app_dev`, no superuser/database/role/schema creation |
| Migration owner | `tgp_migrator_dev`, non-superuser, owns only dedicated databases |
| Bootstrap administrator | `tgp_bootstrap_dev`, first-boot setup only; never runtime |
| Disposable test DB | `tgp_facility_phase0_test`; no user data permitted |
| Authentication | Host and Unix-socket SCRAM configured; live enforcement **not verified** |

First-boot hook creates roles/databases only on an uninitialized volume. It does not rotate/reset existing credentials or downgrade/reuse PostgreSQL 18 storage. Runtime receives CONNECT/schema USAGE only; future table grants await actual domain schema review. Health check is readiness only, not password verification. Preflight verifies dedicated resource ownership, image, volume, network and loopback port before Compose starts only `postgres` with `--no-recreate`. Any mismatch stops; unrelated infrastructure is not adopted.

### Environment configuration and security results

Created previously missing `.env`, with three independent 256-bit cryptographically random passwords written directly without output. File mode is **0600**. Explicit host/port/database/user/password fields and the application's required `DATABASE_URL` are present; internal migration and external test URLs use the corresponding migration credential. Password/userinfo encoding uses `urllib.parse.quote(..., safe='')`, with special-character round-trip coverage. `.env.example` has matching names and placeholders only. No real credential value is included in this report.

Git metadata is now available in this execution environment, unlike the earlier snapshot. `git check-ignore -q .env` succeeds and `git ls-files -- .env` is empty: `.env` is ignored and untracked. No history or index was changed. A private comparison of generated secret values against deliverable source found **no occurrences outside `.env`**; values were not printed. Existing unrelated environment values/comments and generated credentials are preserved on rerun; tracked/symlink/inconsistent configurations fail safely.

Credentials are **prepared, not live-verified**. The candidate host port is 55432, distinct from the excluded 54322, but no availability claim is made. Before first provisioning, the explicit host-side `--select-port` command selects a free port in 55432–55442 and updates only port-derived host settings, refusing to run if a TGP-Sys container already exists. There is no known/weak default secret.

### Migration and application verification foundation

- `cmd/dbcheck`: prepared genuine pgx host TCP check using the application's URL, database/user/version and role verification, plus mandatory wrong-password rejection with PostgreSQL authentication code 28P01. No URL/driver secret output. **Build and execution are blocked**, so this is source, not authenticated-connectivity evidence.
- Migrations use the separate owner. Forward development migration is explicit; the API never migrates automatically. All rollback convenience commands now target **only** `tgp_facility_phase0_test`. Guards check actual DB identity and reject unexpected user objects before test migration operations. `make migrate-down` aliases test rollback; no development rollback convenience is exposed.
- Host acceptance script checks up/up/down/absence/up on the disposable DB, real integration checks, actual API process startup, probe/docs/contract HTTP availability and graceful SIGTERM shutdown. It leaves PostgreSQL containers and volumes intact. It stops at the first failure and does not mark the report accepted automatically.
- `tools/go.mod` prevents npm's third-party Go fixtures from being treated as backend packages by `go test ./...`. Root package discovery confirms only backend packages. This is a tooling boundary, not removal of application tests.

### Exact sandbox checks and results

| Command/check | Result |
|---|---|
| `python3 scripts/local_env.py` / `make env` | **PASS** — generated/preserved private settings, no secret output; no live authentication claimed |
| `.env` mode check, `git check-ignore -q .env`, `git ls-files -- .env`, private source comparison against generated secret values | **PASS** — 0600, ignored/untracked, no generated secrets in deliverable source |
| `python3 -m unittest discover -s scripts -p 'test_*.py'` | **PASS: 7 tests**, including secret non-disclosure, random/independent credentials, idempotence, unrelated-value preservation, tracked-file refusal, inconsistent URL refusal, encoded special characters and port conflict selection. Unit socket/Git doubles do not prove host availability |
| `python3 -m py_compile scripts/local_env.py scripts/db_preflight.py scripts/migrate.py scripts/test_local_env.py` | **PASS** |
| `bash -n scripts/verify_phase0.sh` and `sh -n deployments/local/init-db.sh deployments/local/migrate.sh` | **PASS** syntax; not PostgreSQL execution |
| `docker compose --env-file .env -f deployments/local/compose.yaml config --quiet` | **PASS** — actual Compose parsing without printing interpolated secrets |
| Source Compose shape checks | **PASS** — dedicated names, PostgreSQL 17 image, loopback bind, new volume and no external network |
| `test -z "$(gofmt -l cmd internal api tests)"` | **PASS** |
| `GOTOOLCHAIN=go1.26.10 GOCACHE=/tmp/facility-go-cache go test -race -count=1 ./internal/platform/config ./internal/platform/httpx ./api` | **PASS**, config/httpx; api has no separate test files |
| `GOTOOLCHAIN=go1.26.10 GOCACHE=/tmp/facility-go-cache go vet ./internal/platform/config ./internal/platform/httpx ./api` | **PASS** |
| `make openapi-lint` | **PASS**, Redocly 2.72.1 validates OpenAPI 3.1.1 without warnings/errors |
| `python3 scripts/migrate.py down` | **PASS negative guard** — expected nonzero rejection before any rollback; `--test` mandatory |
| `go mod download` and `go mod verify` under pinned toolchain/cache | **BLOCKED/FAIL** — proxy.golang.org DNS/socket operation not permitted; neither is claimed passed |
| `go mod tidy` with writable `/tmp/tgp-sys-mod-cache` and cached Go 1.26.10 binary | **BLOCKED/FAIL**, same registry socket restriction; no fabricated `go.sum` |
| `go test ./...`, `go test -race ./...`, `make lint` | **BLOCKED/FAIL**, missing genuine pgx checksums/modules; targeted passes do not satisfy full-project gate |
| `make db-preflight`, `make db-up` | **BLOCKED/FAIL**, Docker access denied; no provisioning command proceeded |
| `make db-verify`, `make test-integration`, `go build ./cmd/api` | **BLOCKED/FAIL**, missing pgx dependencies; no connection/startup claim |
| `make migrate-test-up`, `make migrate-test-down` | **BLOCKED/FAIL** at Docker preflight; no database was modified |
| `make verify-phase0` | **BLOCKED/FAIL** at Docker preflight; deliberately stops rather than skipping required checks |

Actual PostgreSQL version, password-based connectivity, initialized role privileges, health/readiness, image availability, migrations and API startup remain **NOT VERIFIED**. Scalar shell/contract source passes unit/lint checks; actual HTTP docs availability and browser Try-it remain **NOT VERIFIED**. No disposable database was available for a real migration round trip. Prior npm lock/install evidence remains valid; no new full dependency security audit was claimed.

### Files created or updated in this follow-up

- Private local `.env` (ignored/untracked, no credentials listed here).
- `.env.example`, `.gitignore`, `Makefile`, `.github/workflows/ci.yml`.
- `deployments/local/compose.yaml`, new `init-db.sh` and `migrate.sh`.
- New `scripts/local_env.py`, `db_preflight.py`, `migrate.py`, `test_local_env.py`, `verify_phase0.sh`.
- New `cmd/dbcheck/main.go`; strengthened disposable DB guard in `tests/integration/db_test.go`.
- New tooling boundary `tools/go.mod`.
- `README.md`, `docs/database/local.md`, architecture overview, ADR index/0002/0003, new ADR 0011, traceability and this report.

The master plan, business schema and production infrastructure were not modified. No Phase 1 domain table/API was introduced.

### Exact host continuation and acceptance recommendation

Run these in the **owner's WSL terminal**, where Docker has been confirmed to work:

```sh
cd ~/Projects/TGP-Sys/backend
# First creation only; omit this line if the dedicated TGP-Sys container already exists:
python3 scripts/local_env.py --select-port
make db-up
make verify-phase0
```

The [database runbook](../database/local.md) provides read-only health/port inspection and individual commands. Review output without copying `.env`, full Docker environment/config dumps or raw credential-bearing errors. The host runner resolves genuine Go dependencies, verifies modules/tests/lint, performs actual authenticated checks and disposable migrations, then actual HTTP startup/docs availability. Update this report with **host-generated results** before acceptance. Browser Scalar Try-it is a separate manual check requiring the pinned CDN.

**Remaining risks/blockers:** denied agent Docker/sockets and Go registry access; candidate port/image unverified; new init hook/dbcheck/host runner require actual host execution; Go checksums absent; transitive security review and immutable image/action pins remain outstanding. No password reset or destructive workaround is authorized or supplied.

**Stop condition:** configuration preparation is complete; mandatory Phase 0 acceptance remains blocked. Stop for owner review and host-side execution. Phase 1 is **NOT STARTED**.


---

## Snapshot before final committed-checkout and Bruno acceptance review

# Phase 0 — Final verification and acceptance report

Review date: 2026-10-31 (host execution timestamp not supplied)  
Status: **READY FOR OWNER REVIEW**  
Overall acceptance: **PENDING** — fresh-clone bootstrap evidence and owner sign-off remain outstanding; manual Scalar Browser Try-it is unverified.  
Scope: **Phase 0 only. Phase 1 NOT STARTED.**  
Evidence source: owner's authorized WSL host verification summary in this conversation, supplemented by direct sandbox checks below.  
Reviewed workspace: branch `main`, HEAD `0de108d` plus uncommitted/untracked foundation files. The host-tested commit/run identifier was not supplied; do not treat HEAD alone as the tested deliverable revision.

## 1. Executive summary

The owner reports successful authorized host execution of PostgreSQL authentication, dependency verification, Go tests/vet, environment safety tests, OpenAPI lint, disposable migration round trip, integration tests, real HTTP endpoints and graceful shutdown. These results resolve the earlier host-execution blockers. No current implementation failure is reported. Sandbox restrictions previously prevented Docker/network execution; those historical failures are not product failures and are not relabeled as sandbox passes.

Direct follow-up inspection confirms required Phase 0 artifacts, all eleven ADRs, genuine `go.sum` and npm lockfiles, correct secret-file exclusion and local documentation link targets. With host-populated dependency caches, the full Go race suite, full vet and module verification now also pass in the sandbox. No Docker/database/process operation was repeated during this report-finalization task.

**Do not declare unconditional Phase 0 completion yet.** Master-plan §10 explicitly requires clean-clone bootstrap, which is not established by the supplied existing-workspace host results. Manual Scalar rendering/Browser Try-it, Bruno execution and hosted CI also lack evidence. Browser Try-it remains a documentation acceptance item; hosted CI end-to-end T23 belongs to the later release gate, while its Phase 0 foundation is present. No future business requirement is claimed implemented.

Previous chronological evidence is retained in [verification history](phase-00-verification-history.md). This report supersedes its current-status claims without changing what ran or failed historically.

## 2. Evidence provenance and boundaries

| Evidence class | What is established | What is not established |
|---|---|---|
| Owner-reported host execution | The concrete successful results enumerated in §3, on the authorized WSL host | No raw transcript, execution timestamp, commit SHA or hosted run URL was supplied; none is invented |
| Direct sandbox checks in this finalization | Files/ADRs/locks present; document links; local secret-file mode/Git exclusion; race tests, vet, module verification, seven Python tests, formatting and OpenAPI lint | No independent replay of live PostgreSQL, migrations, Docker health, HTTP startup or shutdown |
| Earlier sandbox evidence | Configuration/HTTP negative tests, Compose syntax, credential generation/preservation/redaction and migration safety guards | Unit doubles and configuration parsing do not prove real authentication or database behavior |

Host PASS means an owner-attested execution result, not a claim the agent ran it in the sandbox. The repository commands below are reproducible equivalents mapped to the supplied results, **not a fabricated exact host transcript**. The existing `make verify-phase0` runner includes normal/race tests, migration checks and HTTP checks, but no unreported step is inferred merely because that runner exists.

## 3. Final PASS / FAIL / PENDING checklist

| Check | Status | Evidence / applicable gate |
|---|---|---|
| Dedicated PostgreSQL 17.11 healthy and accessible | **PASS — host** | Owner reports healthy dedicated database; master-plan Phase 0 Docker/DB foundation |
| Password-authenticated host TCP connection | **PASS — host** | Owner reports success; repository reproducer: `make db-verify` |
| Incorrect password rejected | **PASS — host** | Explicit owner result; not inferred from pg_isready |
| Runtime DB login is non-superuser | **PASS — host** | Explicit owner result; privilege separation also present in bootstrap source |
| Dependencies downloaded and verified | **PASS — host + sandbox** | Owner reports download/verification; `go.sum` now present; direct `go mod verify` passes |
| `go test ./...` | **PASS — host** | Owner reports Go tests passed |
| `go test -race ./...` | **PASS — sandbox** | Direct full-module command passed using cached dependencies/results; not asserted as a separately attested host command |
| `go vet ./...` | **PASS — host + sandbox** | Owner result and direct full-module vet pass |
| gofmt check | **PASS — sandbox** | No unformatted files in cmd/internal/api/tests |
| Seven environment safety tests | **PASS — host + sandbox** | Explicit owner result; direct unittest rerun: 7 passed |
| OpenAPI 3.1.1 skeleton lint | **PASS — host + sandbox** | Owner result and direct Redocly lint pass |
| Disposable migration up/down/reapply | **PASS — host** | Explicit owner result; reproducer: `make migrate-test-up`, `make migrate-test-down`, up again; no existing-data rollback performed here |
| Real integration tests | **PASS — host** | Explicit owner result; reproducer: `make test-integration` |
| Real Go server startup | **PASS — host** | Real probe/HTTP and graceful shutdown results establish process startup |
| `/health/live` and `/health/ready` HTTP 200 | **PASS — host** | Explicit owner result; reproducer: `make api-smoke` |
| Unimplemented assets route HTTP 404 | **PASS — host** | Explicit owner result; no asset feature is represented as implemented |
| `/docs` and `/openapi.yaml` accessible | **PASS — host** | Explicit owner result; this proves HTTP availability only |
| Graceful shutdown | **PASS — host** | Explicit owner result; automated runner contains SIGTERM/wait check |
| Unrelated Docker infrastructure unchanged | **PASS — host report** | Explicit owner confirmation; no infrastructure operations in this finalization |
| Required Phase 0 foundation artifacts and ADR coverage | **PASS — sandbox review** | §4 inventory and §5 ADR mapping; 33 key artifacts checked, 11 ADRs present |
| Private `.env` protection | **PASS — sandbox review** | Mode 0600; `git check-ignore` matches; `git ls-files -- .env` empty; earlier generated-secret source comparison passed |
| Fresh-clone/bootstrap reproducibility | **PENDING — explicit Phase 0 exit gate** | No fresh-clone result supplied; existing-workspace validation is insufficient; foundation/locks still need to be included in a reviewed versioned source snapshot |
| Manual Scalar Browser rendering and Try-it | **PENDING** | No browser evidence supplied; HTTP 200 for HTML is not an interactive request demonstration |
| Bruno collection execution | **PENDING** | Collection exists; no execution evidence supplied; no auth flow is expected before Phase 2 |
| Hosted CI / complete T23 | **PENDING** | Workflow exists; no hosted run supplied. Phase 0 CI foundation is delivered; complete release acceptance deferred to Phase 11 |
| Dedicated vulnerability/transitive-license review; immutable image/action pins | **PENDING — hardening risk** | Module checksums/test success are not an advisory or license scan |
| Owner approval / Phase 0 sign-off | **PENDING** | This report requests review, not automatic advancement |
| Current demonstrated implementation failure | **None reported** | Earlier sandbox access/download failures remain historical; no current FAIL is invented |

## 4. Required deliverables and repository structure

| Phase 0 requirement | Present artifacts | Assessment |
|---|---|---|
| Authoritative scope and requirement IDs | `ASSET_MANAGEMENT_BACKEND_MASTER_PLAN.md`; `docs/requirements/traceability.md` | Baseline preserved; future module tables/services/tests explicitly deferred |
| Modular monolith foundation | `cmd/api`; `internal/platform/config`, `db`, `httpx`; `docs/architecture/overview.md` | System-only composition, no business modules or speculative endpoints |
| Configuration and structured logging | config implementation/tests; API composition; HTTP transport/tests | Validated environment, JSON slog, redacted failures, server timeouts and correlation IDs |
| Local PostgreSQL / migration tooling | `deployments/local/compose.yaml`, `init-db.sh`, `migrate.sh`; paired probe SQL; `scripts/migrate.py` | Dedicated PostgreSQL 17, separated roles, test-only rollback guard |
| Health server | `cmd/api/main.go`; HTTP handler and tests | Host live/ready/startup/shutdown results supplied |
| Database verification | `cmd/dbcheck/main.go`; `tests/integration/db_test.go` | Real host-path authentication and guarded disposable DB checks |
| API contract/docs foundation | `api/openapi.yaml`, `embed.go`, `docs.html`; `redocly.yaml`; locked tooling; Bruno probes | Lint and HTTP availability pass; interactive browser use pending |
| Development and test commands | `Makefile`; `scripts/local_env.py`, `db_preflight.py`, `test_local_env.py`, `api_smoke.py`, `verify_phase0.sh` | Reproducible equivalents documented in README/runbook |
| Onboarding and repository rules | `README.md`, `AGENTS.md`, `docs/database/local.md` | Present; no deployment/push/Phase 1 permission implied |
| CI foundation and locks | `.github/workflows/ci.yml`; `go.mod`, `go.sum`; `tools/package.json`, `package-lock.json`, `go.mod` boundary | Present; source snapshot/clean-clone and hosted execution remain unproven |
| Phase report | This report; historical evidence archive | Current status/provenance and unresolved gates explicit |

No required Phase 0 deliverable file was found missing. Presence does not mean every future policy recommendation or release test is implemented.

## 5. Architecture decisions and open design gates

| Required ADR topic | ADR | State |
|---|---|---|
| Go version/router/modular monolith | 0001 | Adopted Go/net/http/slog foundation |
| PostgreSQL version | 0002, superseded locally by 0011 | PostgreSQL 17.11 dedicated setup; owner reports host validation |
| DB access and migrations | 0003, 0011 | pgx + golang-migrate; host dependency/migration results supplied; sqlc deferred until real queries |
| IDs and codes | 0004 | Recommendation; confirm UUID choice, normalization and home/current project semantics before domain DDL |
| Authentication transport/scoped access | 0005 | Opaque-session recommendation; auth implementation and precise deny/expiry/MFA policy remain Phase 2 decisions |
| Contract-first OpenAPI/docs | 0006 | Contract/lint/HTTP foundation verified; Browser Try-it pending |
| Object storage/scanning | 0007 | Port/quarantine recommendation; provider, retention and scanning remain later decisions |
| CI/testing | 0008 | Foundation present; hosted execution not supplied |
| Status revisions/safe movement primitives | 0009 | Architecture documented; no status or movement implementation |
| Consistency/idempotency/operating limits | 0010 | Recommendations recorded; phase-specific policies and benchmarks remain open |

All required decision topics are covered; recommendations do not authorize future schema or business-policy changes. Pending deployment/storage/auth/identity choices are documented future gates, not missing Phase 0 business implementations.

## 6. Database, environment and migration evidence

Configuration inspected: official `postgres:17.11-alpine`, project `tgp-sys-local-dev`, container `tgp-sys-postgres-dev`, network `tgp-sys-dev-network`, persistent volume `tgp-sys-postgres17-dev-data`. Current non-secret local settings identify host `127.0.0.1:55432`, DB `tgp_facility_dev`, runtime user `tgp_app_dev`, migration owner `tgp_migrator_dev`, test DB `tgp_facility_phase0_test`. The owner reports the dedicated PostgreSQL 17.11 connection healthy/authenticated; exact runtime inspection output was not supplied, so configuration facts and reported results remain distinct.

Three independent local credentials remain only in ignored mode-0600 `.env`; no secret value is printed here. URL userinfo encoding and preservation/refusal behavior are covered by safety tests. Application runtime is separate from bootstrap/migration credentials. Runtime non-superuser status and incorrect-password rejection have explicit host evidence; additional privilege checks exist in the verifier but are not individually inferred from an unspecified transcript.

Migration 000001 is infrastructure-only (`platform_probe`); no domain schema exists. Up/down/reapplication and integration tests passed on the disposable test DB per the owner. No destructive command was run during this finalization. Refer to the [database runbook](../database/local.md) for exact reproduction and rollback safeguards; no production migration or volume removal is authorized.

## 7. API documentation and process evidence

Implemented system GET/HEAD routes remain `/health/live`, `/health/ready`, local `/docs`, and local `/openapi.yaml`. Unknown assets route returns 404. Host reports real probe status 200, documentation/contract availability and graceful process shutdown. OpenAPI remains the authoritative 3.1.1 skeleton and lint passes.

Scalar 1.76.4 is an exact-version external CDN dependency. **No manual browser rendering, Try-it request/response, browser console/CSP check, or auth interaction is claimed.** To close the manual item, open local `/docs`, render both probe operations, send each request and observe actual HTTP 200 response bodies and request IDs. Note browser, environment and outcome without credentials. Bruno execution likewise remains unreported. Authentication/authorization tests await actual Phase 2 APIs; full T22 is not claimed by Phase 0.

## 8. Direct checks executed during this finalization

These commands ran in the sandbox; no host result is silently substituted for a sandbox command:

| Exact command/check | Result |
|---|---|
| `GOTOOLCHAIN=go1.26.10 GOCACHE=/tmp/facility-go-cache GOPROXY=off go test -race ./...` | PASS; packages compile, config/httpx report cached successes; infrastructure/composition packages have no standalone tests. No DB/integration tag or live startup inferred |
| `GOTOOLCHAIN=go1.26.10 GOCACHE=/tmp/facility-go-cache GOPROXY=off go mod verify` | PASS: all modules verified |
| `GOTOOLCHAIN=go1.26.10 GOCACHE=/tmp/facility-go-cache GOPROXY=off go vet ./...` | PASS |
| `python3 -m unittest discover -s scripts -p 'test_*.py'` | PASS: 7 tests |
| `test -z "$(gofmt -l cmd internal api tests)"` | PASS |
| `make openapi-lint` | PASS: Redocly validates OpenAPI without warnings/errors |
| Nonempty required-file/ADR inventory and relative documentation target checks | PASS: 33 key artifacts, 11 ADRs, all checked local targets exist |
| `stat -c '%a %n' .env`; `git check-ignore .env`; `git ls-files -- .env` | 600; ignored; no tracked `.env` |
| `git branch --show-current`; `git rev-parse --short HEAD`; `git status --short` | Inspected main/0de108d and uncommitted deliverables; no staging/commit/push performed |

The original source secret-pattern review and generated-secret comparison are recorded in the history; no broader secret-history audit is claimed. No code changed in this task, so no new application regression or infrastructure rerun was required for documentation edits.

## 9. Dependencies and versions inspected

| Component | Pin / evidence |
|---|---|
| Go language/toolchain | 1.26.0 / go1.26.10; direct Go checks pass |
| PostgreSQL | 17.11 configured; owner reports actual 17.11 healthy |
| pgx | v5.11.0; real go.sum now present; direct module verification passes |
| golang-migrate | v4.20.1 image tag; disposable round trip host-reported PASS |
| Scalar | 1.76.4 pinned; HTTP shell available, browser use pending |
| Redocly | 2.72.1 pinned/locked; direct lint passes |
| Node/npm | Documented 24.21.0/11.19.0; no new host version transcript supplied |
| Transitive Go dependencies | Recorded in actual `go.mod`/`go.sum`; no checksum fabrication or dependency changes in this finalization |

Exact tags/versions are not immutable image/action digests. Module verification does not prove absence of vulnerabilities or complete license compliance.

## 10. Files and change safety

This task changes documentation only: the phase report, a history archive, related README/runbook status, traceability and ADR verification metadata. No Go/Python/shell/SQL/Compose/workflow implementation, dependency pin/lock, secret file, Git index or infrastructure is changed. The owner-created Go lockfile is inspected, not rewritten.

The full Phase 0 deliverable file inventory follows (excluding the unchanged master plan, private `.env`, installed dependencies, caches and Git metadata). These are working-tree deliverables, not a claim that they are all committed:

- `.env.example`
- `.github/workflows/ci.yml`
- `.gitignore`
- `AGENTS.md`
- `Makefile`
- `README.md`
- `api/docs.html`
- `api/embed.go`
- `api/openapi.yaml`
- `bruno/bruno.json`
- `bruno/health-live.bru`
- `bruno/health-ready.bru`
- `cmd/api/main.go`
- `cmd/dbcheck/main.go`
- `db/migrations/000001_platform_probe.down.sql`
- `db/migrations/000001_platform_probe.up.sql`
- `deployments/local/compose.yaml`
- `deployments/local/init-db.sh`
- `deployments/local/migrate.sh`
- `docs/architecture/overview.md`
- `docs/database/local.md`
- `docs/decisions/0001-go-http-modular-monolith.md`
- `docs/decisions/0002-postgresql.md`
- `docs/decisions/0003-db-migrations.md`
- `docs/decisions/0004-identity-codes.md`
- `docs/decisions/0005-authentication-authorization.md`
- `docs/decisions/0006-openapi-docs.md`
- `docs/decisions/0007-private-storage.md`
- `docs/decisions/0008-ci-testing.md`
- `docs/decisions/0009-config-status-movements.md`
- `docs/decisions/0010-consistency-open-decisions.md`
- `docs/decisions/0011-dedicated-postgresql17-local.md`
- `docs/decisions/README.md`
- `docs/phases/phase-00-report.md`
- `docs/phases/phase-00-verification-history.md`
- `docs/requirements/traceability.md`
- `go.mod`
- `go.sum`
- `internal/platform/config/config.go`
- `internal/platform/config/config_test.go`
- `internal/platform/db/db.go`
- `internal/platform/httpx/server.go`
- `internal/platform/httpx/server_test.go`
- `redocly.yaml`
- `scripts/api_smoke.py`
- `scripts/db_preflight.py`
- `scripts/local_env.py`
- `scripts/migrate.py`
- `scripts/test_local_env.py`
- `scripts/verify_phase0.sh`
- `tests/integration/db_test.go`
- `tools/go.mod`
- `tools/package-lock.json`
- `tools/package.json`

## 11. Genuine acceptance gaps, risks and outstanding decisions

1. **Explicit Phase 0 exit gate still PENDING: clean-clone bootstrap.** Successful existing-workspace host validation is not evidence of a fresh clone. Many foundation files/locks are currently untracked. Include the reviewed source and lockfiles in a versioned deliverable, then record clean-clone `make setup` and relevant acceptance results on an isolated local environment. Do not reuse secret files or alter unrelated infrastructure. No commit/push is performed or assumed here.
2. **Manual Scalar Browser Try-it PENDING.** HTTP availability does not establish renderer/CDN/CSP/request-control operation. Preserve this limitation even if the owner reviews the technical foundation favorably; do not claim complete interactive-documentation acceptance or T22.
3. **Owner review/sign-off PENDING.** No known runtime failure remains from supplied results. Missing acceptance evidence is not a reported implementation failure, but it prevents claiming all gates passed. No current Docker or dependency-access blocker remains for the successful host run; prior sandbox denials remain historical environment limitations.
4. **Hardening and release evidence outstanding:** hosted CI/T23, full dependency advisory and transitive-license review, immutable image/action pins, complete secret-history review and optionally self-hosted Scalar. The initial CI foundation is present; a green hosted release pipeline is not claimed.
5. **Recorded dependency risk:** ADR 0003 flags an upstream SendBatch concern for review before query expansion. The scaffold does not use SendBatch. No new advisory scan or security clearance is implied by passing tests.
6. **Future decisions remain deliberate:** UUID/code normalization/home-current project semantics before Phase 1 DDL; session/CSRF/authorization precedence and MFA/SSO before auth; storage/scanning/financial-field privacy; revision/idempotency rules and operating limits; real repository module path/license; production hosting, least privilege, backup/PITR, restore tests and monitoring before real data/deployment. These do not authorize implementation or silently change approved scope.

## 12. Master-plan Phase 0 gate reconciliation

- **PASS:** requirements traceability, required ADR coverage, modular skeleton, local PostgreSQL foundation, minimal healthy server, structured logging/config, README/AGENTS, OpenAPI skeleton/lint, test/lint/development commands, API-doc tooling foundation and initial CI workflow.
- **PASS — owner-reported host execution:** DB authentication, wrong-password rejection, non-superuser runtime, migration up/down/reapply on disposable DB, Go tests/vet, integration tests, network health/readiness, docs/contract availability and graceful shutdown.
- **PASS — direct sandbox corroboration:** Go module verification/full race suite/full vet, seven safety tests, formatting, OpenAPI lint, deliverable/ADR/link inventory and `.env` Git exclusion/mode.
- **PENDING — explicit exit gate:** clean-clone bootstrap reproducibility for the actual deliverable source/locks. Not waived or inferred from silence.
- **PENDING — documented usability/release evidence:** manual Scalar Browser Try-it, Bruno execution and hosted CI. Full future T22/T23 and all business tests remain deferred.
- **PENDING:** explicit owner Phase 0 acceptance. No current FAIL is demonstrated by the supplied host evidence.
- **Scope preserved:** no company/project/asset business schema or API, auth, inventory mutation logic, business UI, production operation or Phase 1 feature was implemented during this task.

## 13. Review decision and stop

**READY FOR OWNER REVIEW.** The previously blocked automated host verification is reported successful and has been reconciled with direct repository/sandbox evidence. **This status is not unconditional Phase 0 acceptance:** fresh-clone evidence is still required by the master plan, and manual Scalar Browser Try-it remains unverified. Review the documented gaps and recommendations, supply the remaining evidence, and explicitly approve or revise Phase 0 before any next-phase authorization.

STOP. Phase 1 is **NOT STARTED**. No deployment, push or infrastructure alteration was performed.


---

## Snapshot before authorized Git recovery attempt

# Phase 0 — Final acceptance verification report

Overall result: **BLOCKED**  
Review date: 2026-10-31  
Exact tested Git revision: `0de108d4f6279966c55f52e4a54b47b043f41fbb`  
Working branch: `main`; detached isolated checkout created outside the repository.  
Phase 1: not started. No deployment, push, staging, commit, reset, clean, database change or infrastructure mutation was performed.

Only result categories used here: **PASS / FAIL / BLOCKED / PENDING MANUAL VERIFICATION**. Provenance and scope are separate from result categories. Prior snapshots are preserved in [verification history](phase-00-verification-history.md); this report supersedes their current acceptance assessment.

## 1. Executive assessment and repository state

**FAIL — committed source completeness.** The exact commit contains only `README.md`. The current README is modified, while application source, instructions/master plan, all eleven ADRs, OpenAPI, migrations, Compose, tests, Go/npm locks and other Phase 0 deliverables are untracked. A detached clean checkout therefore does not contain the host-tested Phase 0 foundation. Successful earlier working-tree tests cannot be attributed to this commit.

**PASS — isolated checkout creation/preservation.** A real local Git clone was created under `/tmp/tgp-phase0-clean-f4s5a0r4/checkout`, without hardlinks or borrowing the working tree. It was detached at the exact SHA above and `git status --porcelain` was empty. No `.env`, node_modules, Go modules/build cache or uncommitted application source was copied into it. Original HEAD/index/uncommitted work were not reset or overwritten. Authorized changes in this task are limited to Bruno acceptance configuration and related documentation/evidence.

**BLOCKED — final Phase 0 acceptance.** The clean-bootstrap gate fails because the deliverable is absent from the committed revision. Live Bruno execution is additionally blocked by sandbox network policy. Scalar Browser Try-it/schema rendering remains **PENDING MANUAL VERIFICATION**. No Phase 1 work is permitted.

## 2. Fresh-checkout bootstrap: exact commands and results

Clone parent: `/tmp/tgp-phase0-clean-f4s5a0r4`. Installed Go 1.26.10 compiler used directly as a tool; build/module caches were new isolated paths, not reused working-repository artifacts:

```sh
GOTOOLCHAIN=local
GOCACHE=/tmp/tgp-phase0-clean-f4s5a0r4/go-build-cache
GOMODCACHE=/tmp/tgp-phase0-clean-f4s5a0r4/go-module-cache
# Go executable:
# /home/muhamed/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.10.linux-amd64/bin/go
```

These variables were supplied through the subprocess environment, not used to change global configuration. Full argv, cwd, exit codes and output are preserved in [clean-checkout command evidence](evidence/phase-00-clean-checkout.json).

| Command | Result | Actual observation |
|---|---|---|
| `git rev-parse HEAD` | PASS | Full SHA above |
| `git status --short` | PASS | Recorded modified README and untracked Phase 0 deliverables; `.env` ignored |
| `git ls-tree -r --name-only 0de108d4f6279966c55f52e4a54b47b043f41fbb` | FAIL | Completeness gate: only README.md exists in the tree; command itself exits 0 |
| `git clone --no-local --no-hardlinks --no-checkout /home/muhamed/Projects/TGP-Sys/backend /tmp/tgp-phase0-clean-f4s5a0r4/checkout` | PASS | Exit 0, isolated Git repository |
| `git checkout --detach 0de108d4f6279966c55f52e4a54b47b043f41fbb` in clone | PASS | Exit 0 |
| `git status --porcelain` in clone | PASS | Empty output, clean detached checkout |
| Required-file inventory in clone | FAIL | Makefile, go.mod, go.sum, AGENTS, master plan, API, main source, report and tooling lock absent; eleven ADRs absent too |
| `make setup` in clone | FAIL | Exit 2: no rule to make target setup; documented onboarding cannot run |
| Isolated Go `mod download` | BLOCKED | Exit 1: no modules specified; go.mod absent, not an application dependency installation result |
| Isolated Go `build ./...` | BLOCKED | Exit 1: directory does not contain main module |
| Isolated Go `vet ./...` | BLOCKED | Same missing-module error |
| Isolated Go `test ./...` | BLOCKED | Same missing-module error |
| `gofmt -l cmd internal api tests` in clone | BLOCKED | Exit 2: source directories absent; empty-source formatting is not counted as a pass |
| `make openapi-lint` in clone | BLOCKED | Exit 2: no lint target, contract/tooling absent |
| Isolated startup, DB connectivity and integration checks | BLOCKED | No server/Compose/migrations in tested revision; no safe isolated live environment available to this sandbox |

No database/container/port/volume was allocated for this failed-source checkout. Existing TGP-Sys PostgreSQL and unrelated infrastructure were neither reused for mutations nor changed. Do not substitute a dirty source export or copy the working `.env` into the checkout and call that a committed fresh-clone pass.

## 3. Bruno collection and CLI evidence

**PASS — static collection/configuration verification.** Original collection had two hardcoded probe URLs and no named local environment or negative requests. This task adds `bruno/environments/local.bru`, parameterizes both probe URLs with `{{base_url}}`, and adds 404/405 cases. The actual public API origin is **http://127.0.0.1:8080**, derived without exposing database credentials. All requests use auth:none; no secret/session or fabricated business API is introduced.

| Request | Expected actual HTTP response |
|---|---|
| GET `/health/live` | 200, `data.status=live` |
| GET `/health/ready` | 200, `data.status=ready`; DB must really be ready |
| GET `/api/v1/assets` | 404, `error.code=not_found`, expected message, empty details, no data |
| POST `/health/live` | 405, `error.code=method_not_allowed`, empty details, no data, `Allow: GET, HEAD` |

Each request tests JSON Content-Type and a generated 32-character lower-hex request ID matching the response header. Actual `@usebruno/lang` parser validated collection metadata, environment, four requests and ten enabled declarative assertions. Response scripts compile with Node's VM parser. This is syntax/structure evidence, not execution of response assertions or proof of safe-runtime behavior.

**PASS — temporary CLI installation/version.** No project dependency/lock was changed:

```sh
npm install --prefix /tmp/tgp-phase0-bruno-cli --no-save --ignore-scripts --fetch-retries=0 --fetch-timeout=10000 --cache /tmp/tgp-phase0-bruno-npm-cache @usebruno/cli@4.0.0
/tmp/tgp-phase0-bruno-cli/node_modules/.bin/bru --version
/tmp/tgp-phase0-bruno-cli/node_modules/.bin/bru run --help
```

Install exit 0; CLI reports 4.0.0. Deprecated transitive-package warnings were emitted; no advisory clearance is claimed. Temporary tooling is outside the repository.

**BLOCKED — live CLI execution**, actual command from the working `bruno/` directory:

```sh
/tmp/tgp-phase0-bruno-cli/node_modules/.bin/bru run --env local --sandbox safe --reporter-json /tmp/tgp-phase0-bruno-results.json
```

Exit 1. All four requests report `connect EPERM 127.0.0.1:8080`. CLI's raw summary says FAIL due to transport failure: four failed requests, **zero tests and zero assertions executed**. Acceptance result is BLOCKED because the sandbox cannot create the connection; no HTTP response/assertion failure is demonstrated. No host Bruno success was supplied. The [sanitized transport evidence](evidence/phase-00-bruno-transport.json) preserves counts/command/error; raw report remains at the temporary path above.

[Bruno README](../../bruno/README.md) contains exact WSL install/run commands, expected statuses, assertion details and a safe-mode invocation. Reuse an existing API or run only `make run` if needed; no DB startup, migration, reset or outage is required. Do not fabricate readiness or use developer sandbox mode to bypass execution policy.

## 4. Scalar/OpenAPI verification and manual browser task

**PASS — working-tree static contract checks.** OpenAPI 3.1.1 lint passes. All four source contract paths match implemented transport paths; GET/HEAD give eight documented read operations. No asset business API is advertised. Committed response examples validate against their JSON Schemas, and existing handler/config tests pass. Embedded `/docs` shell references the pinned Scalar 1.76.4 CDN script and `/openapi.yaml` contract.

**PASS — prior owner-reported host HTTP availability.** The owner explicitly reported `/docs`, `/openapi.yaml`, probes, DB connectivity, migrations and graceful shutdown successful in the authorized host environment. Those results apply to the uncommitted working-tree implementation, not this README-only commit. They are not silently relabeled as sandbox replay or Bruno/browser verification.

**BLOCKED — current live replay.** Read-only urllib attempts to `/health/live`, `/health/ready`, `/docs` and `/openapi.yaml` from this sandbox each raised URLError. Docker client exists, but daemon access is denied. No stopped process or unhealthy DB is inferred from these sandbox restrictions; no replacement service or database was launched.

**PENDING MANUAL VERIFICATION — actual schema rendering and browser Try-it.** Owner must use [the exact browser checklist](phase-00-browser-checklist.md):

- Open **http://127.0.0.1:8080/docs**; verify UI/CDN/contract loads and Health/Meta/Error schemas render without reference errors.
- Select server **http://127.0.0.1:8080**; send actual GET `/health/live` and GET `/health/ready`.
- Expect HTTP 200 and respectively `{"data":{"status":"live"},"meta":{"request_id":"<32 lowercase hex>"}}` and the equivalent `ready` body. Response `X-Request-ID` must equal the actual body value.
- Inspect actual Network requests; an example preview is not a sent request. Capture browser/version/date/revision and sanitized observations.
- If curl also fails, investigate API/bind/DB first. If curl works but browser fails, inspect CDN, CSP, browser origin, proxy/extensions and host loopback forwarding. `localhost` and `127.0.0.1` differ as origins. Do not enable wildcard CORS, loosen CSP or disable browser security as a workaround.

No manual browser/model-rendering or request-control result is claimed. A real readiness 503 is not a passing expected-ready test; do not alter existing DB infrastructure to manufacture an error test.

## 5. Architecture and security review

**PASS — all eleven ADRs reviewed for design consistency.** See the [ADR-by-ADR architecture/security review](../architecture/phase-00-acceptance-review.md) for findings on modular boundaries, pgx/migrations, project isolation, auth/authorization, secrets, contract ownership, dynamic configuration, approval stages, inventory integrity and pinning.

ADR 0002 is historical and superseded by 0011's dedicated PostgreSQL 17.11 setup. No business modules are implemented; future scoped denial/approval rules are documented gates, not silently treated as completed tests. No observation/sync/discrepancy approval may mutate official assets; separate authorized corrective commands remain required. No arbitrary movement/configuration scripting is introduced.

**FAIL — critical acceptance finding C1:** intended source/locks/docs are absent from the tested commit. Existing successful host results cannot establish clean-clone reproducibility of that revision. No additional confirmed critical application vulnerability or approved-scope contradiction was established by this bounded review; no full security certification/advisory scan is claimed.

## 6. Additional working-tree checks and evidence provenance

These do not turn the failed committed checkout into a pass:

| Exact command/check | Result | Scope |
|---|---|---|
| `GOTOOLCHAIN=go1.26.10 GOCACHE=/tmp/facility-go-cache GOPROXY=off go test -count=1 ./internal/platform/httpx ./internal/platform/config` | PASS | Working-tree fresh unit execution, no live DB |
| `make openapi-lint` | PASS | Working-tree Redocly 2.72.1, no lint errors/warnings |
| Python YAML + Draft202012Validator example validation | PASS | Four paths/eight read operations; examples satisfy component schemas; server origin matches configuration |
| Actual Bruno parser + Node VM script compilation | PASS | Four request URLs/methods/status assertions, ten declarative assertions, local environment; no live assertions executed |
| Private-settings inspection limited to public origin/docs flags | PASS | Actual HTTP_ADDR/local/docs configuration matches checklist; no passwords/DSNs displayed |
| `docker version --format '{{.Client.Version}} {{if .Server}}{{.Server.Version}}{{end}}'` | BLOCKED | Client 29.8.2; daemon query denied, no container mutation |
| urllib HTTP GET attempts for four system URLs | BLOCKED | URLError from restricted execution environment; no status code received |
| Prior owner-reported host DB/auth/migration/integration/Go/HTTP/shutdown checks | PASS | Owner attestation preserved; no raw host transcript/revision supplied; not a clean-checkout result |

Exact source-clone outputs are in the JSON evidence. CLI stdout counts/transport errors are recorded in §3 and evidence. Node/Python inline validators are described by operations and actual observations; they do not execute browser code or fabricate network responses. Earlier host/sandbox logs remain in verification history rather than being rewritten as new runs.

## 7. Critical blockers and non-blocking technical debt

| Finding | Result | Required action / boundary |
|---|---|---|
| C1: only README.md committed | FAIL | Review/version intended non-secret Phase 0 source/locks/docs, then repeat the isolated checkout using the resulting exact SHA; do not discard existing work or commit `.env` |
| Clean bootstrap/build/vet/test/lint/startup of actual deliverable revision | BLOCKED | Requires a committed deliverable; no working-tree file substitution |
| Live Bruno CLI assertions | BLOCKED | Run documented safe-mode CLI from authorized WSL host against real API; retain actual exit/assertion evidence |
| Scalar rendering/Browser Try-it | PENDING MANUAL VERIFICATION | Owner executes exact checklist and reports real outcomes |
| Owner final acceptance | PENDING MANUAL VERIFICATION | Review evidence after remaining gates; no silent next-phase approval |

Non-blocking improvements are separately documented: verified immutable container/action pins, advisory/transitive-license review, complete secret-history scan, optional self-hosted renderer and hosted CI/release matrix. A tool tag/checksum is not an advisory scan. Temporary Bruno emitted dependency deprecations and is not added to application dependencies. ADR 0003's recorded SendBatch concern requires review before query expansion; current scaffold does not use it.

Future decisions remain scoped to their phase gates: UUID/code normalization and home/current project semantics before domain DDL; session/deny/MFA policy before auth; configuration revision/idempotency/approval execution semantics; private storage/scanning/financial privacy; actual module-path ownership, licensing, production hosting/backups/restore/monitoring. No new architecture or implementation is introduced to address these now.

## 8. Final Phase 0 acceptance checklist

| Acceptance item | Result | Evidence boundary |
|---|---|---|
| Required foundation files and all eleven ADRs exist in working tree | PASS | Present/reviewed; not committed completeness |
| Required source/docs/locks committed at tested SHA | FAIL | README.md is the only committed file |
| Actual clean isolated checkout created, original work preserved | PASS | Detached clone, clean status, no copied secrets/uncommitted source |
| Fresh-clone onboarding succeeds | FAIL | make setup target absent |
| Fresh-clone dependencies/build/format/vet/tests/OpenAPI lint | BLOCKED | Source/manifests/targets absent |
| Fresh isolated API/DB startup | BLOCKED | No source in clone; no permitted isolated Docker/network environment |
| Prior core working-tree host DB/migration/integration/HTTP results | PASS | Owner-reported, not proof of tested commit |
| Working-tree contract/example/handler consistency | PASS | Direct static/lint/unit checks |
| Bruno structure/local origin/positive and 404/405 cases | PASS | Actual parser/script syntax checks; expected API behavior only |
| Bruno actual live CLI responses/assertions | BLOCKED | Four EPERM connections, zero executed assertions |
| Scalar schema rendering and Browser Try-it | PENDING MANUAL VERIFICATION | Owner evidence required |
| Architecture/master-plan consistency review | PASS | Eleven ADRs reviewed; future business rules not implemented |
| No secret output/copy, no unrelated infrastructure or DB mutation | PASS | Only public API origin in collection/evidence; no .env copied; no destructive operations |
| Owner acceptance | PENDING MANUAL VERIFICATION | No approval inferred |

## 9. Changes, retained evidence and mandatory stop

Changed only Bruno acceptance files and related documentation: existing two probe requests, new local environment and negative requests, Bruno README, browser checklist, architecture/security review, current report/history, README status/links and traceability. Added sanitized JSON acceptance evidence. No Go/Python/shell/SQL/Compose implementation, dependency pin/lock, private `.env`, original Git index/history or running infrastructure was changed. Temporary clone/CLI/raw evidence remain outside the working directory.

**BLOCKED.** Core working-tree host success is preserved, but final acceptance cannot pass a README-only committed checkout, blocked live Bruno assertions and unverified browser interaction. Present this report for owner review, preserve the source, and supply evidence for the actual reviewed revision and manual checks before acceptance.

STOP. Phase 1 is not started. No deployment or push was performed.
