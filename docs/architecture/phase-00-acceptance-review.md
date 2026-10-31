# Phase 0 architecture and security acceptance review

Design consistency result: **PASS** for the reviewed working-tree foundation, with the limits below. Source recovery is now **PASS** at actual host commit 46555ea; Phase 0 is now **ACCEPTED AND CLOSED** on later owner host/manual evidence; earlier blocked statements below are historical. The original C1 finding below is historical. Source/design review is not implementation of future business controls or a security certification.

Reviewed all eleven ADRs against master-plan §§0–4, 6–10 and 13, current system transport/configuration, Compose/bootstrap/migration guards, secret handling, API contract and collections. ADR 0002 is explicitly superseded by 0011; it does not require running PostgreSQL 18. No new architectural mechanism is proposed.

## Critical acceptance finding

**C1 — FAIL: committed source completeness.** Exact HEAD `0de108d4f6279966c55f52e4a54b47b043f41fbb` contains only README.md. The Phase 0 application, Go/npm locks, API contract, migrations, Compose, instructions and ADRs exist only as untracked working-tree files; README itself is modified. An actual detached clean checkout cannot bootstrap. Existing host results therefore cannot be represented as successful tests of that committed revision. Preserve all work, review/version the intended non-secret deliverable, and test that exact resulting revision. No commit, staging, reset or push was performed here.

This is a delivery/reproducibility defect, not evidence that the host-tested application failed at runtime. Never resolve it by copying uncommitted source or `.env` into the checkout and calling that a clean clone.

## Eleven-ADR consistency matrix

| ADR | Result | Review finding / remaining boundary |
|---|---|---|
| 0001 — Go/router/modular monolith | PASS | net/http, slog and composition/platform boundaries match Phase 0. No microservices, invented business services or repository coupling. Module import path remains provisional. |
| 0002 — original PostgreSQL choice | PASS | Historical PostgreSQL 18 recommendation is explicitly superseded by 0011; no downgrade/reuse of its volume is instructed. |
| 0003 — pgx/migrations | PASS | Versioned SQL and explicit migration commands; API never auto-migrates. Bootstrap, migration and runtime privileges are separated by 0011. SQL/query code for future domains is not implemented. |
| 0004 — IDs/codes | PASS | ID, code, serial and QR separation, historical claim retention, home/current project distinction and normalization approval gate preserve master-plan invariants. Recommendations are not finalized domain DDL. |
| 0005 — authentication/authorization | PASS | Revocable-session recommendation, CSRF, default deny, scoped resources and explicit precedence gate are consistent. Business routes are absent; no auth bypass or static admin token exists. This does not prove future project isolation. |
| 0006 — OpenAPI/docs | PASS | Contract-first 3.1.1, only implemented read operations, embedded source and local-only docs. Lint/examples/route checks are verified; owner-confirmed renderer/GET-live Try-it now supports Phase 0 PASS; no exhaustive browser-model coverage is claimed. |
| 0007 — private storage/scanning | PASS | Private provider boundary, quarantine and authorization are recommendations for the attachment phase. No public file endpoint or prematurely chosen provider is introduced. |
| 0008 — CI/tests | PASS | Commands/workflow/locks exist in working tree; no deploy job. C1 prevents a clean checkout/CI of those files. Hosted CI execution and stronger scan tooling remain unproven. |
| 0009 — statuses/movement primitives | PASS | Referenced immutable configuration revisions, safe implemented primitives and transactional execution preserve controlled dynamic configuration. No arbitrary workflow scripting or mutable historical references. |
| 0010 — consistency/inventory/future limits | PASS | Durable idempotency and revision policies remain design gates. Observations and individual discrepancy decisions are evidence; a separate authorized corrective command is mandatory. No automatic official-asset updates are introduced. |
| 0011 — isolated PostgreSQL 17 | PASS | Dedicated names/network/volume, loopback publish, non-superuser runtime, separate migration owner, first-boot-only role setup and guarded disposable-test rollback. Owner reports real authentication/migration success; no live replay was possible here. |

## Scope and security observations

- One company/multiple projects and cross-project identities remain approved constraints. Company/project/location/unit/asset/authz tables and their denial cases are future work, not claims covered by system probe tests.
- Unit independence, unspecified precise location, company-managed ownership, immutable identity/history and independent discrepancy approvals are retained in requirements/ADRs. No business records were created or changed in this task.
- `.env` remains ignored/untracked and mode 0600. Its passwords/DSNs are not copied into the clean checkout or Bruno. Only public API origin is present in Bruno. No secret rotation or environment rewrite was performed.
- Runtime uses its application URL; migration/bootstrap credentials are separate. Local host-path auth results are owner-reported; pg_isready/Compose source alone are not treated as authentication proof.
- Docs are local-only, no permissive CORS enabled, and CSP restricts API connections to the same origin. Browser failure must be investigated rather than bypassed by weakening browser/API security.
- Negative Bruno requests cover actual 404/405 behavior only; no fake login, role denial, asset implementation or readiness success is supplied. Readiness 503 must remain a real failure of the expected-ready check.

## Non-blocking improvements and future design gates

1. Pin trusted container images/actions by verified immutable digest/commit, not just mutable exact-version tags. Review updates deliberately; no version bump or infrastructure change in this task.
2. Complete advisory/license scans for application/tooling dependencies. Existing ADR 0003 records a SendBatch concern; current scaffold does not use SendBatch. No new advisory clearance is asserted. Temporary Bruno CLI installation emitted dependency deprecation warnings; it is not added to the application's dependency tree.
3. Consider a self-hosted verified Scalar renderer for offline use after manual evidence identifies actual requirements. Current CDN/CSP/browser behavior is not certified by static validation.
4. Hosted CI, complete contract/security matrix and full Bruno release coverage are later release evidence. The Phase 0 collection is a probe/negative foundation.
5. Confirm UUID/code normalization/home-current semantics before Phase 1 DDL; session/deny/MFA policy before Phase 2; config revision/idempotency/approval details before affected commands; storage/financial-field privacy before attachments. Production deployment, backup/PITR/restore and monitoring remain separate approval gates.

No additional confirmed critical application vulnerability or approved-scope contradiction was established by this bounded review. That statement is not a full security audit. Live collection execution is BLOCKED; browser/schema rendering is PENDING MANUAL VERIFICATION. See the final Phase 0 report for exact evidence and overall BLOCKED acceptance status.

## Authorized Git recovery follow-up

The original source-completeness finding remains true of original HEAD because sandbox .git writes are blocked. A complete local candidate commit was safely created in temporary storage and its actual clean checkout passed offline bootstrap/build/tests/lint. This resolves candidate-source completeness, not original branch recovery or live/manual verification. See the current phase report for exact candidate SHA, inventory and host recovery procedure. No additional ADR or Phase 1 design was introduced.

## Actual original-repository recovery resolved

Host commit `46555ea1a4f3ff14cd24354e1ec49c14e3ba74ec` contains exactly 68 manifest files. All implementation/config/locks/collections are byte-identical to the reviewed candidate. Its actual detached checkout passed offline bootstrap/build/unit/race/vet/lint. Original C1 source completeness is resolved; live exact-commit checks and browser Try-it remain separately gated. No Phase 1 code was found despite the commit-message label.

## Final owner sign-off

All mandatory Phase 0 evidence is reconciled: actual Git recovery/clean checkout and code checks, owner PostgreSQL/migration/lifecycle results, inspected real Bruno host JSON, and explicit manual Scalar confirmation. PHASE 0 — ACCEPTED AND CLOSED. Earlier C1/transport/manual statuses are historical; no architecture or future business feature was changed. Phase 1A is not started and requires separate authorization.
