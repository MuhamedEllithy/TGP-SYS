# Requirements traceability

The master plan remains authoritative. IDs below index its sections rather than replace or narrow them. Each group covers all detailed rules in the cited section. Planned services/tables/endpoints are **not implemented** and must be refined in the relevant phase ERD/contract. No T01–T21 business acceptance scenario is claimed by Phase 0. T22/T23 have only foundation coverage.

| ID | Master plan | Requirement | Endpoint/service/artifact | DB | Verification | Phase | Status |
|---|---|---|---|---|---|---|---|
| P0-001 | §10 Phase 0 | Go/config/logging/modular skeleton | cmd/api; platform/config, db, httpx | platform_probe (infrastructure only) | config/httpx unit tests; go vet; build | 0 | Foundation built and host-verified; see final report |
| P0-002 | §10 Phase 0; §7 | Local DB, migrations, reproducible setup | Makefile; deployments/local | platform_probe; schema_migrations | integration; migration up/down; bootstrap | 0 | Host DB/migration PASS; committed clean-checkout gate FAIL |
| P0-003 | §10; §13 | ADRs and gated documentation | docs/decisions; README; AGENTS | — | review of baseline and ADRs | 0 | Complete as documents; recommendations remain open |
| CORE-001 | §2.1 | One company, many projects, shared identities/departments/teams | Future organization/projects application services | Phase 1 ERD: organization/project tables | T01 | 1–3 | Deferred |
| CORE-002 | §2.1 | Timed memberships, primary team, independent scoped grants | Future organization/authz services | Memberships, role bindings (proposed) | T01; revocation and effective-date tests | 2 | Deferred |
| LOC-001 | §2.2 | Project/location tree; no cycles or cross-project parents; historical safety | Future /api/v1/projects and locations | Phase 1 ERD locations/projects | T02 | 1,3 | Deferred |
| LOC-002 | §2.3 | Independent multi-location units with exactly one primary when active | Future units application service | units, unit_locations (proposed) | T03 | 1,3 | Deferred |
| LOC-003 | §2.3 | Optional precise location only from unit links; impact preview on unlink | Future units/assets services | Unit links + placement (proposed) | T04 | 3,5 | Deferred |
| CAT-001 | §2.4 | Category tree and typed reusable inherited dynamic schemas | Future catalog service | Category/attribute definitions + revisions (proposed) | T05; cycle/type/retirement negatives | 4 | Deferred |
| CAT-002 | §2.4 | Manufacturer/model and versioned templates; no retroactive rewrite | Future catalog/templates services | Model/template versions (proposed) | T06 | 4 | Deferred |
| AST-001 | §2.5 | Company-managed assets; ID/code/serial/QR separation; historical code reservation | Future assets service; ADR 0004 | Assets/code claims/history (proposed) | T07–T08; contention | 1,5 | Deferred |
| AST-002 | §2.3,2.5 | Unit or direct placement; inheritance/override; parent-child cycles | Future asset placement/assembly commands | Placement + parent link history (proposed) | T04,T09 | 5,7 | Deferred |
| AST-003 | §2.5 | Independent configurable status dimensions and evidence-backed verification | Future assets status service; ADR 0009 | Status revisions/assignments/events (proposed) | T10 | 5 | Deferred |
| AST-004 | §2.9 | Private files, checksum/MIME/scan, warranty/purchase metadata and privacy | Future file adapter; ADR 0007 | Attachment metadata/links (proposed) | T20; upload security negatives | 5 | Deferred |
| ELI-001 | §2.1,2.6 | Multiple eligible teams, inheritance/add/exclude/replace; no task creation | Future eligibility/responsibility service | Eligibility and responsibility rules (proposed) | T11 | 6 | Deferred |
| MOV-001 | §2.7 | Versioned configurable movement types mapping only to safe primitives | Future movement type catalog; ADR 0009 | Movement types/versions (proposed) | T12 | 7 | Deferred |
| MOV-002 | §2.7 | Atomic authorized execution, both project scopes, child impact and history | Future movement executor | Movement events/executions (proposed) | T09,T13; stale/replay/rollback tests | 7 | Deferred |
| APR-001 | §3.1 | Fine-grained default-deny scoped RBAC; delegated grants and deterministic precedence | Future authz application policy; ADR 0005 | Role/permission/scope bindings (proposed) | T01,T13,T19; enumeration denial | 2 | Deferred: precedence design open |
| APR-002 | §3.2 | Maker-checker stages, policy versions, expiry and independent execution | Future approval service | Approval versions/requests/decisions (proposed) | T13; stale/self-approval negatives | 1,7 | Deferred |
| IMP-001 | §2.8 | Manual/template/Excel preview and approved idempotent intake; no preview mutation | Future import/intake service | Import batch/row/issue/result (proposed) | T06,T14 | 8 | Deferred |
| INV-001 | §2.8 | Scoped frozen baselines/manifests, assignments, offline bounded sync and replay | Future inventory protocol service | Campaign/baseline/manifest/batch/observation (proposed) | T15–T16 | 9 | Deferred |
| INV-002 | §2.8 | Observation and approval never mutate official assets; independent discrepancy decisions | Future discrepancy review service | Discrepancies/individual decisions (proposed) | T16–T17; DB invariance assertions | 9 | Deferred |
| INV-003 | §2.8 | Separate authorized correction with linked approval, revision and deduplication | Future corrective command orchestration | Correction link to official command (proposed) | T18; replay/closed/revoked/stale tests | 9 | Deferred |
| RPT-001 | §2.10 | Scoped reports/exports; unit counts once; privacy/formula injection protection | Future reporting/export service | Read models/export jobs (proposed) | T19,T21 | 10 | Deferred |
| RPT-002 | §2.5,2.10 | QR resolves stable locator; no authorization grant; safe labels/reprints | Future QR service | QR identity/print jobs (proposed) | T08; QR access/label tests | 5,10 | Deferred |
| API-001 | §6,7 | Authoritative OpenAPI 3.1 and implemented routes only | /health/live; /health/ready; local /docs; /openapi.yaml | No domain tables | httpx tests; Redocly lint; network smoke | 0,11 | Phase 0 skeleton/lint and host HTTP PASS; full business contract deferred |
| API-002 | §6,7 | Interactive real requests, examples/errors; Bruno/Hurl convenience | Scalar shell; Bruno probes; Python smoke in place of optional Hurl | — | T22; manual Try-it | 0,2,11 | Host docs HTTP PASS; Bruno live BLOCKED; Browser Try-it PENDING MANUAL VERIFICATION; auth deferred |
| SEC-001 | §3.3,8.2 | Revocable sessions, KDF, CSRF/CORS, TLS, rate/size limits and scoped negatives | Future authn/authz; current HTTP timeout/config redaction | Session/audit schema later | config/httpx negative tests; later T01,T19,T20 | 0,2,5,11 | Partial: system surface controls only |
| SEC-002 | §8.1,8.2 | Immutable audit/domain histories, redaction, private storage and restore | Future audit/files; current structured logs | Audit/domain events (proposed) | Mutation/audit and restore tests later | 1–11 | Deferred except redacted system logs |
| SEC-003 | §8.3 | Bounded queries/payloads, metrics and benchmark assumptions | Future module handlers and observability | Indexes justified by query plans later | Volume/latency tests after real sample data | 10 | Deferred: no SLA invented |
| SEC-004 | §8.4,9 | Pinned CI, locks, real DB tests, no deployments or secrets | Proposed .github/workflows/ci.yml | Ephemeral PostgreSQL only | T23; secret review; clean bootstrap | 0,11 | Locks present in working tree; committed checkout FAIL; hosted CI BLOCKED |

## Change control and open items

No approved functional requirement changed. ADR recommendations do not approve schemas or business policies. Exact IDs should be split further when contracts are reviewed, retaining these parent IDs and master-plan links. Master plan §13 items map to ADR 0005 (1,3), 0003/0004 (2,10), 0009 (4), 0007 (5,9), 0010 (6–8). Owner review of code-scope/normalization and authorization precedence is required before affected implementation. All future modules and production deployment remain outside authorization.

## Dedicated PostgreSQL follow-up (Phase 0 only)

P0-001/P0-002/SEC-001/SEC-004 now also trace to ADR 0011, `scripts/local_env.py`, `scripts/db_preflight.py`, `cmd/dbcheck`, `scripts/migrate.py` and `scripts/verify_phase0.sh`. Environment safety, Git exclusion and Compose parsing are locally verified; host authentication, migration round trip and startup are now owner-reported PASS; direct Go race/vet/module checks pass. Clean-clone and Browser Try-it evidence remain pending. No domain requirement implementation was added.

## Final acceptance verification

An actual detached checkout of `0de108d4f6279966c55f52e4a54b47b043f41fbb` contains only README.md: source completeness FAIL, downstream clean-bootstrap checks BLOCKED. The working-tree Bruno collection now has a local environment and 404/405 cases; parser/script syntax PASS, live CLI BLOCKED by EPERM. Scalar browser/schema rendering and Try-it are PENDING MANUAL VERIFICATION. See the [current report](../phases/phase-00-report.md), [architecture review](../architecture/phase-00-acceptance-review.md), and [browser checklist](../phases/phase-00-browser-checklist.md). No future requirement is implemented by this acceptance work.

## Git recovery verification follow-up

Original repository commit is BLOCKED by read-only .git. A complete isolated candidate commit and clean-checkout offline bootstrap/tests/lint are PASS, with exact SHA/inventory in the current report. Live candidate DB/API and Bruno remain BLOCKED by sandbox sockets; Scalar remains PENDING MANUAL VERIFICATION. This is evidence for Phase 0 only; no future requirement implementation is implied.
