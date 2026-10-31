# Phase 0 — Formal acceptance record

**PHASE 0 — ACCEPTED AND CLOSED**

## 1. Phase identification

TGP-Sys Facility Management Platform, Phase 0: requirements/technical decisions and production-minded backend foundation. Review date: 2026-10-31. No later phase is authorized.

## 2. Scope delivered

Go modular monolith foundation; PostgreSQL 17.11 local configuration; migrations; private environment/config validation; structured logs; health/readiness HTTP server; OpenAPI/Scalar/Bruno foundation; automated tests and CI scaffold; eleven ADRs; onboarding/security/architecture/traceability documentation. No business modules or UI.

## 3. Verified Git revision

Implementation: **`46555ea1a4f3ff14cd24354e1ec49c14e3ba74ec`**. Actual committed inventory: 68 original manifest files, reconciled against the 65-file candidate. This documentation-only closure adds three documentary artifacts, not application behavior. The containing sign-off commit is resolved with `git log -1 --format=%H -- docs/phases/phase-00-acceptance.md` if committed; no SHA is invented before Git succeeds. Current Git administrative result is in the [final report](phase-00-report.md#9-files-and-local-documentation-commit).

## 4. Final acceptance matrix

| Gate / requirement | Final status | Evidence source | Tested revision | Method / command / reference |
|---|---|---|---|---|
| Complete Git recovery and committed source | PASS | Direct repository/Git audit | `46555ea1a4f3ff14cd24354e1ec49c14e3ba74ec` | 68 committed manifest files; [inventory](evidence/phase-00-recovery-inventory.json) |
| Reproducible clean checkout | PASS | Recorded direct isolated verification | `46555ea1a4f3ff14cd24354e1ec49c14e3ba74ec` | Actual detached clone; fresh extraction/build/node_modules; offline `make setup`; [commands](evidence/phase-00-recovery-bootstrap.json) |
| Go dependency integrity, compilation, formatting, vet | PASS | Recorded direct verification; owner confirms | `46555ea1a4f3ff14cd24354e1ec49c14e3ba74ec` | `go mod verify`, `go build`, `make lint`; no source/lock changes |
| Unit and race tests | PASS | Recorded direct verification; owner confirms | `46555ea1a4f3ff14cd24354e1ec49c14e3ba74ec` | `go test -count=1 ./...`; `go test -race -count=1 ./...` |
| Seven environment safety tests | PASS | Recorded direct verification; owner WSL evidence | `46555ea1a4f3ff14cd24354e1ec49c14e3ba74ec` | `python3 -m unittest discover -s scripts -p 'test_*.py'` |
| PostgreSQL 17.11 health, role separation, TCP auth, wrong-password rejection, non-superuser runtime | PASS — Owner-provided WSL evidence | Explicit owner acceptance evidence | `46555ea1a4f3ff14cd24354e1ec49c14e3ba74ec` (owner context) | `make db-verify`; [final host evidence](evidence/phase-00-final-host-evidence.json) |
| Disposable DB migration up/down/reapply and integration | PASS — Owner-provided WSL evidence | Explicit owner acceptance evidence | `46555ea1a4f3ff14cd24354e1ec49c14e3ba74ec` (owner context) | `make migrate-test-up`, `make migrate-test-down`, up again; `make test-integration` |
| HTTP startup, health/readiness, contract and graceful shutdown | PASS — Owner-provided WSL evidence | Owner confirmation plus inspected Bruno responses | `46555ea1a4f3ff14cd24354e1ec49c14e3ba74ec` (owner context) | GET live/ready/openapi 200; actual process lifecycle; `make api-smoke` / host runner |
| OpenAPI 3.1.1 skeleton validation | PASS | Recorded direct validation; owner confirms | `46555ea1a4f3ff14cd24354e1ec49c14e3ba74ec` | `make openapi-lint`; implementation/contract consistency checks |
| Bruno real requests and assertions | PASS — Owner-provided WSL evidence | Codex directly inspected owner-generated JSON, not a new CLI execution | `46555ea1a4f3ff14cd24354e1ec49c14e3ba74ec` (owner context; JSON has no SHA field) | CLI 4.0.0; 4/4 requests; 7/7 tests; 10/10 assertions; owner reports 459 ms/exit 0 |
| Scalar interactive docs | PASS — Owner-confirmed manual verification | Explicit owner browser confirmation | `46555ea1a4f3ff14cd24354e1ec49c14e3ba74ec` (owner context) | UI/operations visible, real Try-it GET `/health/live`, HTTP 200 and expected live body |
| No committed/local secret leakage; safe environment foundation | PASS | Direct source/manifest/Git exclusion review and safety tests | `46555ea1a4f3ff14cd24354e1ec49c14e3ba74ec` + documentation diff | `.env` ignored/untracked/0600; actual-private-value/pattern checks; no credentials in closure records |
| Required foundation deliverables / ADRs / onboarding / CI scaffold | PASS | Direct file/decision coverage review | `46555ea1a4f3ff14cd24354e1ec49c14e3ba74ec` | All required artifacts and 11 ADRs present; CI foundation only, not claimed hosted execution |
| Owner Phase 0 closure authorization | PASS — Owner authorization | Current explicit final sign-off instruction | `46555ea1a4f3ff14cd24354e1ec49c14e3ba74ec` source baseline | Scope limited to Phase 0 closure and one reviewed local documentation commit |


## 5. PostgreSQL validation

**PASS — Owner-provided WSL evidence.** Dedicated PostgreSQL 17.11 healthy; separate runtime/migration credentials; authenticated TCP succeeds; incorrect password rejected; runtime role non-superuser; disposable migration up/down/reapply and integration tests pass. No infrastructure or DB changes during sign-off.

## 6. Go and tooling results

**PASS.** Recorded exact-revision clean checkout/offline installation, compilation, unit/race tests, formatting/vet, seven environment tests and OpenAPI lint. Commands, environment and outputs are preserved in [bootstrap evidence](evidence/phase-00-recovery-bootstrap.json). Original host confirmation supplements, not replaces, direct recorded checks.

## 7. Bruno execution summary

**PASS — Owner-provided WSL evidence.** CLI 4.0.0; requests **4/4**, tests **7/7**, assertions **10/10**, **459 ms**, **exit 0**. Codex directly inspected `/tmp/tgp-bruno-results.xIMTEf/results.json` and verified its outcomes. SHA-256 `d6ca62b7b9992405bb8322d35af985433122514080b09c27e1eb0b0dd76350de`. Execution was by the owner, not Codex. Version/duration/exit/revision context are owner-provided metadata; counts and per-request status/test/assertion success are independently visible in the JSON. [Sanitized evidence](evidence/phase-00-final-host-evidence.json).

## 8. Scalar browser verification

**PASS — Owner-confirmed manual browser verification.** Renderer/operations visible; Try-it worked; real GET `/health/live` returned 200 and the expected live response. No automated browser test or screenshot supplied/claimed. The owner did not separately attest every optional model/error/readiness UI step.

## 9. Security verification

**PASS for Phase 0 scope.** Private `.env` excluded/untracked/0600; source/documentation credential checks; no runtime/data/cache/build artifacts in intended commit. Role separation, non-superuser auth and wrong-password rejection have owner host evidence. Future business authorization is not implemented or certified. Application/config/lock source remains unchanged.

## 10. Known limitations

- Clean-checkout verification used offline registry archives with fresh extraction/build/install directories and committed checksum/lock validation. Earlier online registry/Docker/socket denials remain historical sandbox limitations, not retrospectively successful sandbox runs. Online cold-registry bootstrap was not independently demonstrated by the agent; the verified offline reproduction satisfies the recorded Phase 0 reproducibility scope explicitly confirmed by the owner.
- Host DB, migration, integration, lifecycle and live HTTP results are owner-provided WSL evidence. Codex did not independently execute inaccessible host checks during this closure.
- Bruno host JSON was actually accessible and inspected. CLI version 4.0.0, process exit 0, duration 459 ms and revision context were supplied by the owner; the JSON independently establishes response/test/assertion counts and outcomes, not those process/revision fields.
- Scalar PASS covers only the owner's explicit observations: successful renderer, visible operations, working Try-it, real GET `/health/live`, HTTP 200 and expected live response. No automated browser test, screenshot, browser version/time, browser readiness Try-it or exhaustive error-model/UI coverage is invented.
- No business API, authentication/authorization service, domain schema, inventory protocol or frontend is implemented or accepted by Phase 0. No production readiness/release/deployment approval is implied.


## 11. Deferred technical debt

- Hosted CI execution/full T23, verified immutable image/action pins, broader dependency-advisory/transitive-license and secret-history scanning remain hardening/release work; CI scaffold and pinned versions are delivered. Module verification/offline npm audit output is not an advisory clearance.
- Scalar remains an exact-version CDN dependency; a verified self-hosted bundle can be considered later. Optional additional browser/Bruno regression coverage remains future quality work.
- Existing ADR 0003 records a SendBatch concern for review before query expansion; the foundation does not use SendBatch. No new vulnerability scan is claimed.
- Module-path/license ownership, IDs/code normalization/home-current semantics, auth/session/MFA/scoped-deny details, configuration revision/approval/idempotency rules, storage/scanning/privacy, operating limits and production backup/restore/monitoring remain their documented future gates.
- The existing source commit's misleading `first commit phase 1` message is retained; content inspection confirms Phase 0 foundation only. No prior commit is amended or rewritten.


## 12. Owner approval basis

The owner's current final sign-off task explicitly confirms completion of mandatory technical verification, provides real Bruno/manual Scalar evidence, and authorizes Phase 0 closure plus one local documentation commit. This record relies on that scoped authorization; it is not an agent-forged signature, production release approval or next-phase authorization. Previous sandbox failures remain in [history](phase-00-verification-history.md).

## 13. Final closure statement

**PHASE 0 — ACCEPTED AND CLOSED**

Mandatory technical gates are satisfied with explicit evidence provenance. The documentation staging attempt was blocked by the read-only `.git/index.lock`; no documentation commit was created and HEAD remains the verified implementation SHA. The final report supplies the precise safe host procedure. This is an administrative delivery limitation, not a fabricated commit or revocation of the owner's technical acceptance.

**Phase 1A — Database Architecture and ERD Design: NOT STARTED — AWAITING SEPARATE OWNER AUTHORIZATION.** No deployment, push or infrastructure/data mutation. STOP.
