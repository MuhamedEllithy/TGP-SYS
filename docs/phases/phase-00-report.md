# Phase 0 — Git recovery and final acceptance report

Overall result: **BLOCKED**

Original repository: `/home/muhamed/Projects/TGP-Sys/backend`, branch `main`.

Original and current original-repository HEAD: `0de108d4f6279966c55f52e4a54b47b043f41fbb`.

New local **candidate** commit: **`e7c15cac3a120c9977f1503226c8424bf9c06c80`**, created only in an isolated temporary repository, branch `phase0-recovery-candidate`.

Candidate tree: `268eb702c28523332d265920a50eaad7ac68f4f0`.

No new commit was created on the original branch: its Git metadata is read-only to this sandbox. Candidate verification is real committed-source verification, but is **not** a claim the working repository was recovered. Results use **PASS / FAIL / BLOCKED / PENDING MANUAL VERIFICATION**. Phase 1 is not started; no deployment or push occurred.

## 1. Git audit, staging review and permission boundary

**PASS — repository/security audit.** Root, branch and full original SHA were inspected. Initial HEAD tracks only README.md; README is modified and the foundation is untracked. There was no pre-existing staged work. All 63 pre-existing intended files were reviewed as Phase 0 source/docs/locks/CI and explicitly enumerated, without `git add .` or `git add -A`. Two recovery documents were added before candidate creation, giving 65 candidate files. Final evidence adds three documentary JSON files, giving 68 paths in the current [precise host staging manifest](phase-00-staging-manifest.txt).

Reviewed `.gitignore` excludes `.env`/variants, installed dependencies, caches, build outputs and logs. Actual `.env` is ignored, untracked and mode 0600. Known private password values and common key/token/private-key patterns were checked without displaying values. No real credentials, local databases, Docker volumes, runtime logs, node_modules or unrelated project files belong to the manifest. Historical sanitized JSON under docs is deliberate validation evidence, not application logs. Generated recovery bundle remains ignored under `.tools` and is never staged.

**BLOCKED — original-repository staging/local commit.** Actual explicit Python subprocess used `git add --` followed by the reviewed 63 individual paths. Exit 128:

```text
fatal: Unable to create '/home/muhamed/Projects/TGP-Sys/backend/.git/index.lock': Read-only file system
```

No escalation, chmod/chown, alternate index for the protected repository, reset, clean, force checkout or permission bypass was attempted. Original HEAD/index and unrelated work are preserved. This is an execution-environment blocker, not an unexplained failure to commit or an ambiguity about approved source.

**PASS — safe candidate commit.** A separate writable clone of the old revision was created under `/tmp/tgp-phase0-recovery-9vdwwwub/candidate`. Only the 65 reviewed non-secret manifest files were copied there and staged individually. Index content bytes were compared to reviewed SHA-256 fingerprints; tracked paths exactly matched the manifest. A local commit was created with the configured author identity and message:

```text
chore(phase0): establish Go backend foundation
```

It contains the complete Phase 0 foundation at candidate creation time. This writes only new temporary Git metadata, not the original repository's protected `.git`. Later report/evidence/host-guide edits are working-tree documentation changes, not part of the tested candidate; application source/config/locks/collections remain byte-identical to that candidate.

Default `git diff --cached --check` reported only intentional two-space Markdown hard breaks, including the unchanged master plan and ADR metadata. Every finding was inspected: Markdown files, exactly two trailing spaces, no tabs/other whitespace findings. These semantic Markdown line breaks were preserved, not mistaken for runtime/security errors. The host guide checks non-Markdown whitespace with Git and permits exactly two spaces in Markdown through a staged-content validator. No test failure was skipped.

## 2. Complete committed inventory and local handoff

The candidate's complete 65-file inventory, sizes and SHA-256 values are in [recovery inventory](evidence/phase-00-recovery-inventory.json). All eleven ADRs, the authoritative master plan, onboarding/instructions, Go module/checksums, npm lock, OpenAPI/embed/docs, scripts/tests, migrations, Compose and CI are represented. The original Git HEAD still contains only README.md.

Candidate files:

- `.env.example`
- `.github/workflows/ci.yml`
- `.gitignore`
- `AGENTS.md`
- `ASSET_MANAGEMENT_BACKEND_MASTER_PLAN.md`
- `Makefile`
- `README.md`
- `api/docs.html`
- `api/embed.go`
- `api/openapi.yaml`
- `bruno/README.md`
- `bruno/assets-not-found.bru`
- `bruno/bruno.json`
- `bruno/environments/local.bru`
- `bruno/health-live.bru`
- `bruno/health-method-not-allowed.bru`
- `bruno/health-ready.bru`
- `cmd/api/main.go`
- `cmd/dbcheck/main.go`
- `db/migrations/000001_platform_probe.down.sql`
- `db/migrations/000001_platform_probe.up.sql`
- `deployments/local/compose.yaml`
- `deployments/local/init-db.sh`
- `deployments/local/migrate.sh`
- `docs/architecture/overview.md`
- `docs/architecture/phase-00-acceptance-review.md`
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
- `docs/phases/evidence/phase-00-bruno-transport.json`
- `docs/phases/evidence/phase-00-clean-checkout.json`
- `docs/phases/phase-00-browser-checklist.md`
- `docs/phases/phase-00-git-recovery.md`
- `docs/phases/phase-00-report.md`
- `docs/phases/phase-00-staging-manifest.txt`
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

A verified, mode-0600, ignored recovery bundle retains the local candidate commit:

```text
.tools/phase0-recovery-e7c15cac.bundle
```

It contains no `.env` or installed/generated runtime artifacts. An actual bundle clone using `--branch phase0-recovery-candidate` verified the exact SHA and all 65 file hashes; that check clone was removed. The bundle advertises its named branch only: an initial clone without `--branch` could not select a default HEAD, so the explicit-branch form is required. It is not staged. Do not reset/cherry-pick over uncommitted work merely to import it. The [Git recovery guide](phase-00-git-recovery.md) gives explicit WSL manifest staging, staged-content secret/whitespace validation, local commit, SHA recording and inventory commands. The current 68-path manifest also includes the final documentation evidence. An owner host commit will have its own exact SHA and must be recorded; do not mislabel it as this candidate without comparison.

## 3. Clean committed checkout and bootstrap environment

**PASS — exact committed checkout.** A second clone was made from the candidate repository with `--no-local --no-hardlinks --no-checkout`, then detached at `e7c15cac3a120c9977f1503226c8424bf9c06c80` in `/tmp/tgp-phase0-recovery-9vdwwwub/checkout`. It had all 65 tracked files, clean status, no `.env`, no node_modules and no application source borrowed from the original working tree. No original source change or build artifact was copied into it after checkout.

Only task-local bootstrap credentials were generated there (mode 0600). DB target was explicitly set to **127.0.0.1:16543**, and API startup target to **127.0.0.1:16544**, distinct from existing development ports 55432/8080 and unrelated 54322. Those are isolated intended targets, not claims of available/running services. No Docker service, network, volume or database was created. Original private environment was not copied or changed.

**BLOCKED — online onboarding with empty caches.** `make setup` first ran with fresh module/build/npm cache paths. Registry DNS/socket creation for proxy.golang.org was denied; exit 2. This is recorded, not counted as a network bootstrap pass.

**PASS — offline dependency installation.** The same documented `make setup` command then succeeded using existing downloaded registry archives as an offline dependency source, with empty task-local extraction/build caches and no existing node_modules. Exact non-secret environment:

```text
PATH prefix=/home/muhamed/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.10.linux-amd64/bin
GOMODCACHE=/tmp/tgp-phase0-recovery-9vdwwwub/module-cache
GOCACHE=/tmp/tgp-phase0-recovery-9vdwwwub/build-cache
GOPROXY=file:///home/muhamed/go/pkg/mod/cache/download
GOSUMDB=off
npm_config_offline=true
npm_config_cache=/tmp/facility-npm-cache
```

Installed Go 1.26.10 is a compiler tool, not an original project artifact. The file proxy supplies Go registry archives; extraction and compilation occurred afresh. Committed go.sum and `go mod verify` validate archive/module content. No new online transparency/advisory lookup is claimed with GOSUMDB off. npm ci installed 284 packages from the committed lock using tarball cache, without reusing original node_modules. Its offline audit message is not an independent current vulnerability scan.

After setup/tests, both `git status --porcelain` and `git diff --exit-code` were clean. No committed manifest/lock/source changed during bootstrap.

## 4. Exact clean-checkout commands and results

Full argv, working directories, exit codes and outputs: [Git commands](evidence/phase-00-recovery-git.json) and [bootstrap commands](evidence/phase-00-recovery-bootstrap.json). These describe the candidate, not the old README-only checkout.

| Command / check | Result | Evidence |
|---|---|---|
| Original `git add -- <63 explicit paths>` | BLOCKED | Exit 128, original index.lock read-only; no original local commit |
| Candidate clone/copy/explicit stage/content comparison/local commit | PASS | Exact candidate SHA above, 65 files |
| Second clone and detached checkout of candidate | PASS | Clean committed source, required artifacts all present |
| `python3 scripts/local_env.py --port 16543` in clone | PASS | New private local-only configuration, not copied dev credentials |
| `make setup` with empty caches/default online registry | BLOCKED | Exit 2, proxy.golang.org socket/DNS restriction |
| `make setup` with offline dependency environment | PASS | Exit 0: tidy, download, verify, locked npm ci |
| `go mod verify` | PASS | All modules verified in fresh extraction cache |
| `go build -o /tmp/tgp-phase0-recovery-9vdwwwub/api ./cmd/api` | PASS | Actual compilation from candidate source |
| `go test -count=1 ./...` | PASS | All candidate packages compile; config/httpx pass, other packages have no standalone tests |
| `go test -race -count=1 ./...` | PASS | Actual uncached race test execution |
| `make lint` | PASS | gofmt check and full go vet |
| `make openapi-lint` | PASS | Redocly validates OpenAPI 3.1.1 |
| `python3 -m unittest discover -s scripts -p 'test_*.py'` | PASS | Seven environment-safety tests |
| `docker compose --env-file .env -f deployments/local/compose.yaml config --quiet` | PASS | Clone configuration resolves without printing secrets; not live Docker proof |
| `python3 scripts/local_env.py --run /tmp/tgp-phase0-recovery-9vdwwwub/api` | BLOCKED | Exit 1, redacted http_listen_failed; direct port-16544 socket probe confirms PermissionError errno 1 |
| `docker info --format '{{.ServerVersion}}'` | BLOCKED | Daemon socket access denied; no safe isolated live PostgreSQL available |
| Fresh isolated DB auth/migrations/integration/HTTP probes | BLOCKED | No live isolated DB/API; no development DB reused or mutated |
| `git status --porcelain`; `git diff --exit-code` after checks | PASS | Empty/exit 0: committed source and locks unchanged |

Prior owner-reported authorized host PostgreSQL 17.11/auth/wrong-password/non-superuser/migration/integration/health/docs/shutdown success remains **PASS** with owner-attestation provenance. It is not relabeled as a fresh live candidate-commit run or independent sandbox observation.

## 5. Bruno live results and exact WSL command

Collection structure/local environment, four implemented-system/negative behaviors, ten declarative assertions and response script syntax were previously checked. Current task reran the actual documented safe-mode CLI:

```sh
# cwd: bruno/
/tmp/tgp-phase0-bruno-cli/node_modules/.bin/bru run --env local --sandbox safe --reporter-json /tmp/tgp-phase0-bruno-recovery-results.json
```

**BLOCKED** — exit 1, four connection attempts, four transport errors (`connect EPERM 127.0.0.1:8080`), **zero HTTP responses, zero tests and zero assertions executed**. Passed assertions: 0. Failed response assertions: 0, because none ran. Raw CLI labels the run FAIL due to transport; this is not evidence that the API returned wrong statuses. No host Bruno PASS is fabricated.

Use the existing [Bruno guide](../../bruno/README.md). Exact owner WSL steps, reusing a running API; if absent run only `make run` separately, with no migration/container change:

```sh
cd ~/Projects/TGP-Sys/backend
curl --fail --silent --show-error http://127.0.0.1:8080/health/live
curl --fail --silent --show-error http://127.0.0.1:8080/health/ready
tgp_bruno_tools=$(mktemp -d /tmp/tgp-bruno-tools.XXXXXX)
tgp_bruno_results=$(mktemp -d /tmp/tgp-bruno-results.XXXXXX)
npm install --prefix "$tgp_bruno_tools" --no-save --ignore-scripts @usebruno/cli@4.0.0
cd bruno
"$tgp_bruno_tools/node_modules/.bin/bru" run --env local --sandbox safe --reporter-json "$tgp_bruno_results/results.json"
```

Expected four successful cases: GET live 200/live; GET ready 200/ready; GET unimplemented assets 404/not_found; POST live 405/method_not_allowed and Allow GET, HEAD. JSON/header/request-ID and error details assertions must actually execute. Do not pass database secrets to Bruno or disable failing readiness assertions.

## 6. Scalar and manual browser verification

OpenAPI lint is **PASS** from the candidate clone. Prior static route/schema/example/handler consistency checks are **PASS**. Prior owner host `/docs` and `/openapi.yaml` HTTP availability is **PASS** as owner-reported working-tree evidence, not a browser interaction claim. Fresh candidate live HTTP/docs verification is **BLOCKED** by socket permissions.

**PENDING MANUAL VERIFICATION** — exact original development configuration gives **http://127.0.0.1:8080/docs**, contract **http://127.0.0.1:8080/openapi.yaml**. Owner must follow [the browser checklist](phase-00-browser-checklist.md): render the pinned Scalar UI and schemas; send real GET `/health/live` and `/health/ready` from the same origin; observe HTTP 200, respectively data.status live/ready, and matching 32-hex meta.request_id/X-Request-ID. An example preview is not a request.

If equivalent curl fails, investigate API/DB/bind first. If curl works but browser fails, inspect CDN, CSP, browser origin/loopback forwarding and Network/Console errors. Do not substitute localhost for the contract's 127.0.0.1 origin, enable wildcard CORS, weaken CSP or disable browser security. No manual evidence was supplied in this task.

## 7. Security, preservation and cleanup

- **PASS:** explicit scoped manifests, no unrelated paths, generated-private-value and credential-pattern checks; no .env, tokens, local storage, node_modules, logs, binaries or ignored recovery bundle staged in candidate.
- **PASS:** original source/index/HEAD preserved, original .env remains ignored/untracked/0600 and unchanged. No production/unrelated infrastructure mutation or original development DB operation.
- **PASS:** candidate includes eleven ADRs and all required foundation code/config/docs/locks; no Phase 1 feature or new architecture introduced.
- **PASS:** exact non-secret command evidence and source inventory persisted in docs. Historical evidence is retained in [history](phase-00-verification-history.md), not rewritten as passing runs.
- Task-created clone/API/installed node_modules/extraction/build/npm caches and generated clone credentials contain no user data. After preserving evidence and the verified ignored bundle, only this task's temporary verification directory is eligible for removal. Shared/pre-existing caches and earlier-task temporary directories must not be removed. Cleanup result: **PASS** — the exact task-created temporary verification directory was removed after confirming both Git clones had no user edits and the ignored bundle verified. Generated clone credentials, node_modules, compiler output and task-local caches were removed with it. Shared/pre-existing caches and earlier-task directories were preserved. The bundle remains for owner review; no source/secret file in the original repository was removed.

## 8. Complete final acceptance checklist

| Gate | Result | Boundary |
|---|---|---|
| Repository/staged/untracked/ignored/security audit | PASS | Intended project-only files identified; no unrelated work |
| Original branch receives authorized reviewed local commit | BLOCKED | Original .git mounted read-only; current original SHA unchanged |
| New isolated local candidate commit and full inventory | PASS | e7c15cac3a120c9977f1503226c8424bf9c06c80, 65 reviewed files |
| Candidate fresh Git checkout and no uncommitted source reliance | PASS | Actual second clone, detached SHA, clean before/after |
| Online registry bootstrap from empty caches | BLOCKED | Sandbox DNS/socket policy |
| Offline dependency archive bootstrap, compilation, unit/race/vet/format/lint | PASS | Fresh extraction/build/installation, committed checksums/locks |
| New environment bootstrap | PASS | Task-generated isolated-port credentials, no existing environment reused |
| Isolated candidate live PostgreSQL/migration/API/docs checks | BLOCKED | No permitted Docker/listen/socket path; existing dev DB untouched |
| Earlier core authorized host execution | PASS | Owner-reported; not claimed as this candidate's live run |
| Bruno actual HTTP responses/assertions | BLOCKED | Four transport errors, zero executed assertions |
| Scalar browser rendering/Try-it | PENDING MANUAL VERIFICATION | Owner evidence required |
| Owner final acceptance | PENDING MANUAL VERIFICATION | No approval inferred |

## 9. Remaining risks and required owner action

The safe original-repository local commit is still required on the WSL host using the manifest/security guide; the sandbox cannot perform it. The bundle/candidate does not move original HEAD. Record the actual new original commit SHA, then verify that revision without overwriting uncommitted work. Current final report/evidence/guide edits are not falsely represented as part of the tested earlier candidate.

Non-blocking risks remain: online registry reproducibility is not independently verified here; offline audit output is not an advisory scan; container/action tags lack immutable digest pins; hosted CI/transitive license/advisory/secret-history review remains outstanding; Scalar depends on a CDN; module-path ownership, future identity/auth/config/approval/storage/operating policies and production backup/deployment choices remain gated. None authorizes Phase 1 or production work.

**BLOCKED — request owner review/acceptance only after original Git recovery and remaining host/manual evidence.** No reset, clean, push, deployment, database recreation or Phase 1 implementation was performed.
