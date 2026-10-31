# Repository instructions

Read `ASSET_MANAGEMENT_BACKEND_MASTER_PLAN.md` completely before changes. It is the approved functional baseline; ADRs cannot silently change it. Current authorization: Phase 0 only. Phase 1 has NOT been authorized.

- Backend only. No business UI, deployment, production migration, remote push, reset, or deletion of user data.
- Modular monolith: domain owns invariants; application services own authorization and transaction boundaries; adapters own HTTP/SQL. Never import another module's private repository or mutate its tables.
- One company, many projects; shared employee identities; membership is not authorization. Units are independent from locations. Unknown exact unit location stays unknown.
- Company-managed assets only. Stable ID, code, serial and QR remain separate. Historical code claims are never reused. Referenced statuses/movement configuration/history cannot be deleted.
- Inventory observations, sync and discrepancy approvals never mutate official assets. Each discrepancy requires an independent decision and a separate authorized correction command.
- Future business routes require default-deny function + scope + resource checks, negative tests, and audit facts. Eligibility is not authorization or a maintenance task.
- Contract-first OpenAPI 3.1; update contract and tests together. Do not add speculative stub endpoints. Domain schemas begin only after Phase 1 approval.
- Run `make fmt`, `make lint`, `make test`, `make openapi-lint`, migration round trip and integration tests relevant to changes. Report blocked checks; never simulate real infrastructure to claim a gate passed.
- No secrets in source, logs, fixtures or Bruno. Do not log configuration structs, URLs, query strings, request headers, driver errors or credentials.
- Keep numbered ADRs, traceability and the phase report current. Stop after the authorized phase and request review.
