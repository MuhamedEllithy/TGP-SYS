# Enterprise Facility Management Platform — Asset Management Backend Master Plan

**Project:** Multi-project Facility & Mall Management Platform  
**Initial deployment:** The Gate Plaza (Egypt)  
**Delivery owner:** Codex, working in gated implementation phases  
**Document type:** Product SRS + Backend Architecture Blueprint + API Contract Plan + Implementation Runbook  
**Version:** 1.0 — Baseline scope agreed 2026-10-31  
**Status:** Approved *functional direction*; detailed technical decisions and schemas below are **proposals to validate during design**, not proof of implementation.  
**Language:** English is intentional to make this document directly usable by Codex.  

> **Instructions to Codex — READ FIRST:** This file is the principal delivery contract. Do not attempt the full system in a single pass. Execute **one phase only**, report evidence and blockers, and STOP for user review before advancing. Do not quietly invent product requirements, create a frontend, execute a production deployment, or claim success when any gate fails. Keep `docs/decisions/` and this plan updated with approved refinements. The user wants design-first, modular, testable, production-minded Go code.

---

## 0. Executive summary and non-negotiable rules

We are building a long-lived, enterprise-oriented system to operate multiple malls/properties for **one managing company**, beginning with The Gate Plaza. Later business modules will include Inventory/Warehousing, Maintenance, Operations, HR, Customer Service, Security, CCTV Incident/Violation Reporting, and cross-module reporting. **V1 delivers the Core Platform and Asset Management backend only**, documented and tested through an interactive API workspace. No business application frontend is to be built now.

**Non-negotiable product rules:**

1. One company, **multiple projects**; a single employee/account/team may serve multiple projects and assets with different permissions per scope.
2. **Locations are hierarchical**; **Units are independent records, NOT nodes** of the location tree. A unit may reference multiple location nodes, exactly one designated primary location, within the same project.
3. Asset placement is either **directly at a mall location** or **inside a unit**. For a unit spanning multiple locations, an optional exact location may be selected *only from that unit's linked locations*; otherwise location is deliberately unspecified at unit level and MUST NOT be fabricated as the primary node.
4. Track **mall/company-managed assets only**, including mall-owned equipment physically inside a leased/owned unit. **Exclude owners' and tenants' private assets** from the asset register.
5. Separate stable internal asset ID, administrative asset code, manufacturer serial number, and QR identifier. Codes can be generated or manually entered under project policy. QR resolves an opaque stable record locator and authorizes access at request time.
6. Support flexible hierarchical categories, reusable dynamic attributes, inheritance, templates, and controlled bulk operations.
7. Support **parent–child assets**, audited component re-parenting, location inheritance/override, and a future boundary with inventory items/consumables.
8. Multiple **independent configurable status dimensions** (single- or multi-select) with history. Critical lifecycle rules remain invariant. Verification is supported by inventory evidence and does not imply auto-correction.
9. **Movements have configurable types** created, edited, retired or deleted if truly unused. User configuration may choose only **safe, implemented movement primitives**; it cannot inject arbitrary logic or bypass invariant validation. Executed movements remain historically meaningful.
10. **Offline inventory observation NEVER changes official asset data**, including on sync or even on discrepancy approval. Each discrepancy has its **own individual decision/approval**. Any accepted correction requires a **separate authorized, auditable command**.
11. Super Admin controls fine-grained module/function permissions with scoped roles/groups, and configurable maker–checker workflows. Maintenance *eligibility* belongs in Assets; task assignment and repair details belong to the future Maintenance module.
12. **API-first, backend-first**: provide full OpenAPI contract, interactive request/response reference and runnable automated tests. **No React dashboard or offline PWA client** in V1. Implement offline inventory **server protocol/API** for future client.

### 0.1 Approval semantics

- **Approved functional decisions:** all items in the previous list and the detailed functional scope in this document, based on the planning dialogue.
- **Recommended technical choices:** Go, modular monolith, PostgreSQL, OpenAPI 3.1, Scalar, Bruno, Go tests; verify versions, licenses, library maturity and compatibility at Phase 0. Do not treat a library recommendation as a permanently locked choice.
- **Open design decisions:** exact authentication transport, migrations/ORM/query strategy, deployment target, file/object storage provider, and precise authorization and approval implementation must be recorded in Architecture Decision Records (ADRs) before irreversible work. Recommendations are given below.
- **Not approved:** production deployment; tenant private asset tracking; accounting depreciation; building the business frontend; implementing future modules.

## 1. Product boundary and vocabulary

| Term | Definition and boundary |
|---|---|
| Company | Managing organization; central employee/team/department structure. Not multi-tenant SaaS in V1. |
| Project | Managed property/mall (The Gate Plaza initially), containing location tree, units and assets. |
| Location | Physical node in one project's tree: building, floor, room, corridor, roof, outdoor area, etc. Unlimited practical depth subject to protective limits. |
| Unit | Independently identified commercial/medical/administrative/etc. space. Linked to one or more location nodes; may have contracts in a later dedicated module. |
| Asset | Identifiable company/mall-managed equipment or tracked component with its own immutable record ID. |
| Parent asset | An asset containing independently tracked child assets; not a quantity or stock bin. |
| Inventory item | Stock/consumable/spare part; may give rise to a uniquely tracked asset upon commissioning later. |
| Asset category | Configurable hierarchical classification controlling allowed attributes and default maintenance team eligibility. |
| Dynamic attribute | Typed, validated, reusable metadata property of an asset/category, with effective configuration revisions. |
| Asset template | Reusable set of defaults for a category/model; does not silently rewrite existing assets. |
| Movement | Business transaction (transfer/hand-over/external repair/disposal etc.) with workflow and immutable history. |
| Inventory campaign | A defined physical audit/census scope and snapshot, not a mutation of official asset data. |
| Observation | What an inspector actually saw, including offline observations and evidence. |
| Discrepancy | A specific mismatch between the campaign baseline and an observation; individually adjudicated. |
| Responsibility | Accountability for assets; distinct from permission to view/edit and distinct from maintenance eligibility. |
| Maintenance eligibility | Which teams may work on an asset/category/project; actual work orders belong in the future maintenance module. |
| Approval | Controlled decision(s) for a proposed operation; business action is executed separately/atomically only when policy allows. |

### 1.1 Scope: INCLUDED (V1 backend)

- Core: company seed, multi-project management, location tree, essential unit registry & location linking, staff identities, departments, teams, memberships, roles, permissions and approval policies; audit events.
- Asset classification: category tree; dynamic attribute library/config; field inheritance and validation; manufacturers, models; asset templates.
- Asset records: administrative code policy, optional manufacturer serial, stable QR, company ownership boundary; category-specific metadata, parent/child relationships; placements; flexible statuses; asset profiles, photos/attachments, warranty/purchase metadata; responsibility and maintenance-team eligibility.
- Operations: asset creation/edits with policy; configurable movement types and secure execution; asset lifecycle and disposition; events/history and explicit approvals.
- Intake: manual API registration, Excel import with preview and row-level errors, batch/template creation; identify verification state distinctly from registration.
- Physical inventory: campaign/scope snapshot, assignments, offline-capable sync API, idempotency, per-observation reconciliation, individually approved discrepancies, reports; NEVER implicit asset mutations.
- Reporting: basic asset inventory/coverage, counts by project/location/category/status/department/team, movements, import progress, discrepancies, attachments, audit, CSV/XLSX exports, QR label output; PDF report/label only where acceptance confirms a required format.
- API: well-defined REST API, OpenAPI 3.1 source of truth, interactive documentation/test UI, Bruno collection, automated Go tests and Hurl contract smoke tests, CI.

### 1.2 Scope: EXCLUDED / future integrations

- React/Vite admin frontend; mobile/offline PWA user interface; barcode/QR camera UI (protocol will be ready).
- Detailed Units/Property/Lease/Owner/Tenant contract module, sublease, financial lease obligations, customer service workflows. Keep `units` stub and extension boundaries.
- Maintenance orders, SLA, preventive maintenance, technician job/task dispatch, spare-part usage, vendor repair workflow. Define *eligibility metadata* only in V1.
- Warehouse stock ledgers, purchasing, receiving/issuing, valuation; advanced asset-to-inventory conversion transaction.
- HR payroll, attendance, shifts, recruitment; only users/employees/departments/teams necessary for access and responsibility.
- Security patrols/CCTV incident workflows/violations, operations scheduling, enterprise BI, full accounting/depreciation, SaaS tenant isolation.

## 2. Functional requirements and invariants

### 2.1 Multi-project, organization and teams

- `Company -> Departments/Employees/Teams`, and independently `Company -> Projects`. No repeated employee accounts per project.
- Employees can be assigned to multiple projects and teams; require at least one **primary team** when active in operations. Multiple membership allowed to support cross-skilled staff.
- A team may be centrally available to several projects or scoped to one; project membership and team membership have effective validity intervals and audited changes.
- Project assignment **does not automatically grant** asset read/write permissions; a role+scope evaluation is required on every request.
- Responsibility levels: project, location subtree, category subtree, or specific assets; many-to-many where multiple teams/people share responsibility. Provide optional primary responsible department and primary coordinating team; additional responsible parties allowed.
- Don't equate `responsible_team`, `allowed_maintenance_teams`, `assigned_work_order_team` (last is out of scope).

### 2.2 Project/location tree

- Project identity: immutable UUID, code, name, active flag, optional settings; code unique among all projects.
- Location: immutable UUID, `project_id`, `parent_id` nullable, code, name, type, status, metadata, revision; multiple roots permissible when appropriate.
- CRUD and relocate subtree with cycle detection, same-project parent validation, and explicit permission. Location delete only if unreferenced; archive otherwise.
- Query ancestors, descendants, subtree asset totals, paths. Use recursive queries/indexes initially; consider PostgreSQL `ltree` only if necessary and justified by benchmarks.
- Re-parenting location must not corrupt history; path can be derived while historic movements/observations retain sufficient snapshot labels/identifiers.
- Enforce referential integrity and write authorization for *all* project-scoped changes.

### 2.3 Essential Unit Registry (NOT a location node)

- Unit: immutable UUID, project ID, unique `unit_code` inside project, optional display name, active/status, basic configurable type(s), optional area/description.
- `unit_locations`: unit ID, location ID, `is_primary`, optional section/context. At least one linked location for a **fully active** unit; exactly one primary when active, constrained in transactional service logic with DB protections. All locations must match unit project.
- Unit can span floors/locations; no direct `parent_id` tree reference from unit to other units in V1. Unit contracts/parties are future modules, not required for asset creation.
- Asset placement `UNIT`: choose unit; optional location only among linked `unit_locations`. If no exact location selected, report the unit and its location set; **never infer false precision** from the primary location.
- Asset placement `LOCATION`: select a project location; no unit reference. Unit and direct mall placement mutually exclusive.
- A unit's location-link changes require impact preview of affected assets and historical/active references; do not silently invalidate a selected precise location.

### 2.4 Categories, attributes, models and templates

- Category tree is configurable: code, parent, name, description, active state, archived flag, applicable asset type, form/report configuration. Prevent cycles; changes must be audited.
- Common fields are strongly typed in `assets`. Dynamic category-specific fields are configured, NOT new physical database columns per attribute.
- `attribute_definitions`: unique stable key, display names/description, type (text, long-text, integer, decimal, boolean, date, datetime, enum, multiselect, IP, URL, etc.), validation hints, enum options, unit, privacy classification. Never evaluate user-entered executable code.
- `category_attribute_bindings`: category, attribute, required flag, default, searchable/reportable flag, inheritance/override rules, conditional visibility, effective revision. Child categories inherit until explicitly overridden; detect conflicts deterministically.
- Validate all dynamic values server-side; schema-driven field permissions apply separately to read and write. Store value data in JSONB where appropriate with relational definitions and selective indexes; do not indiscriminately GIN-index the entire payload.
- Changing a field from optional to required MUST NOT retroactively invalidate persisted records. Use config version, migration/quality gap reporting, and validation on future writes; offer a controlled backfill, no silent bulk rewrites.
- Manufacturer + Model catalog; models reference manufacturers; cross-category applicability is configurable. Serial uniqueness policy can be scoped per manufacturer+model or project and flagged for review, not assumed universally globally unique.
- Templates: versioned defaults for category/model attributes; individually edited asset values remain sovereign; new versions never silently update previous assets.
- Bulk template creation and Excel import must validate codes, serials, locations, units, custom attributes and permission scopes before commit.

### 2.5 Asset identity, ownership, status and relations

- Asset ID: immutable opaque UUID (recommend UUIDv7 if available/approved), primary key for all links.
- Administrative code: required, case-normalized according to a well-defined policy, generated or manual; configurable project prefix/pattern/sequence. **Never reuse a historically claimed code within its code scope**. Track previous codes and change reason.
- For moves across projects, preserve stable ID and original administrative code; decide separately if a re-code is necessary. Suggest `home_project_id` as administrative-code scope and `current_project_id` as current operational/placement scope; an ownership/code-scope transfer is a distinct approved action. This separation must be validated in an ADR before migration.
- Manufacturer serial: optional at platform level, may be mandatory per category; normalized duplicate checking with reviewed exceptions, not globally unique by default.
- QR: opaque unguessable token/locator or stable asset URL referencing immutable ID; the scanner always reauthenticates and re-authorizes the caller. Do not encode serial, confidential asset properties or personal data into QR itself. Log QR replacements/print jobs as needed.
- Core fields (suggested): ID, home/current project, administrative code, name, category, model/manufacturer refs, serial, ownership (always managing company domain for V1), type/criticality, registration source and status, revision, timestamps, creator/updater; explicit placement handled separately.
- Dynamic fields: typed attribute map checked against effective category schema; history on changes.
- Parent–child tree: arbitrary sane depth cap configurable for performance; one active parent per child; no cycles, no parent in a different project without explicit cross-project policy. Reparent/unlink only via audited operation, preserve assembly history.
- Location policy for child assets: `INHERIT_PARENT` or explicit placement. Parent relocation propagates **derived** location to inheriting children without fabricating independent movement history; override children receive conflict/impact review.
- Lifecycle status, operational status and additional configurable status groups coexist. Each group declares `single` or `multiple`; hard constraints protect mutually exclusive states and required dimensions.
- Verification status: `NOT_VERIFIED`, `VERIFIED`, `DISCREPANCY_PENDING` (or derived equivalent) based on immutable verification evidence; do not let an arbitrary user toggle `VERIFIED` without a qualifying verification process.
- State catalogs and workflow policy are configurable; avoid allowing deletion of referenced statuses. In-use state => deactivate and retain historical version.

### 2.6 Maintenance teams eligibility, not execution

- Map category -> allowed skill/team roles; project -> concrete eligible teams; asset -> optional overrides. Teams can overlap and multiple teams may be eligible for the same asset/category.
- Examples: CCTV => low-current + electrical; fire systems => fire safety + electrical + mechanical; HVAC => HVAC + electrical (examples only, user-configurable).
- Define override semantics explicitly: `INHERIT`, `ADD`, `EXCLUDE`, `REPLACE`; show effective eligibility, source/inheritance, and dates. In-use historic assignments retain context.
- Eligibility alone grants **no API authorization** and creates **no work order**. Future maintenance module will specify tasks, required roles and assignments per incident or category.

### 2.7 Movements and transfers — configurable types with safe primitives

Movement-type catalog must be manageable through protected APIs (future dashboard will use those APIs): `add`, `edit`, `archive/deactivate`, `delete only if unused/unreferenced`, `list`, `preview effective version`.

- Type fields: stable code, name, description, underlying action primitive, source/destination constraints, required fields/documents, eligible asset categories, role permissions, approval policy, status effects, parent/child behavior, effective revision.
- Safe built-in primitives: `RELOCATE`, `TRANSFER_PROJECT`, `HANDOVER_CUSTODY`, `SEND_EXTERNAL`, `RETURN_EXTERNAL`, `PLACE_IN_STORAGE`, `REACTIVATE`, `RETIRE`, `DISPOSE`, `ADMINISTRATIVE_RECODE` (subject to design review). Names are user-editable; primitive behavior is implemented in Go and cannot be dynamically authored by end users.
- New user-defined movement type MUST map to an implemented primitive. Unsupported new primitive => separate engineering change; never bypass hard domain rules using config.
- Requests have draft/pending/approved/rejected/cancelled/executed/failed states; approvals may have multiple stages. Approval does NOT automatically imply unvalidated state mutation: execution re-checks expected revision, eligibility, dependencies and scope in an atomic transaction.
- Once a movement type is referenced by history, it cannot be hard-deleted; archive and version. Historical movement records capture stable type ID + version + selected presentation snapshot.
- Sensitive movement execution is maker–checker; actor cannot approve own request when policy disallows self-approval. On cross-project moves verify both source and target permissions and active project memberships.
- Every successful movement writes immutable movement history/event metadata including from/to, reason, actor, approvals, time, revision; failure must not partially relocate child assets.

### 2.8 Physical inventory: manual + offline API, zero automatic corrections

**Input modes:** legacy partial Excel imports, ongoing manual registration, template-based bulk registration. These create asset records only through explicitly approved intake flows, not through audit observation processing.

**Campaign lifecycle:** draft -> scoped and assigned -> baseline snapshot frozen -> active -> reconcile -> finalized/closed; reopening requires recorded privileged action. Campaign scope may be project, location subtree, unit set, category subtree, or combinations; assign inspectors, deadlines and allowed download scope.

**Offline API contract principles:**

- Future mobile PWA will download an **authorized, bounded campaign manifest** (asset keys and minimal allowed fields/QR identifiers; avoid sensitive property data and bulk excess). API supports paginated/signed incremental snapshots and manifest revision.
- Client queues observations locally (for future frontend to implement), assigning UUID `client_observation_id`, `client_session_id`, observed timestamp, asset/QR if known, observed unit/location, optional dynamic notes/evidence identifiers.
- When online, `POST /inventory/sync-batches` accepts batched observations with idempotency keys and per-item results; dedupe by unique `(campaign_id, device/client_session_id, client_observation_id)` plus request idempotency key. Support safe retries, resume, duplicates, unordered arrival, missing attachments and partial errors without phantom duplicates.
- Server retains authoritative `received_at`; inspector-provided `observed_at` is untrusted/subject to skew rules. Revoked access/closed campaign/expired assignment yields explicit per-record error or quarantined review, **not silent acceptance**.
- An observation is an *evidence statement*, not a command to change asset.location, asset.status, asset.code or create an official asset.
- Comparing observation with immutable campaign baseline creates one or more **discrepancy cases** (missing, mismatch of unit/location/status/serial/label, unregistered-found, duplicate scan etc.); meaningful grouping/de-duplication without silently merging different issues.
- **Each discrepancy independently**: investigate -> review -> approve/reject/request-reinspect -> decision recorded (maker-checker). No bulk acceptance of distinct discrepancies.
- Even APPROVED discrepancy does **not** directly mutate assets. A reviewer with the right permission must initiate a separate corresponding asset command (e.g., move/recode/new registration), subject to its own policy, and link the resulting transaction to the discrepancy. If duplicate approvals accidentally arrive, ensure only one correction command can execute.
- Baseline comparison must not be retroactively rewritten by later asset movements; show both campaign snapshot and current official state; if stale, return `409` conflict and require explicit re-evaluation.
- Every finding, evidence attachment, inspector, decision and corrective transaction is auditable. Offline storage encryption/OS controls are the responsibility of the future client, specified as a client security requirement.

### 2.9 Attachments/documents

- Asset images, invoices, warranty certificates, technical manuals, drawings, QR label files, serial photo, campaign evidence. Extensible attachment types; category-specific required document rules.
- Private object storage abstraction (local dev provider is acceptable; production private buckets/provider chosen later), DB metadata, checksum, MIME sniffing, max size, virus scanning strategy, per-file authorization on upload/download, controlled short-lived links or streamed access, audit downloads of sensitive documents.
- Never trust client filename/path, never execute uploaded content, enforce allowed MIME and safe object keys. Attachment can be linked to movement, discrepancy or asset via explicit references and access rules.
- Purchase cost, vendor, purchase date, warranty start/end and documents can be stored with visibility permissions; **accounting depreciation is out of V1**.

### 2.10 Search, reports, exports, QR printing

- Filtering and pagination by project, location tree, unit, category tree, manufacturer/model, team, responsible department, dynamic attribute (per indexed/searchable fields), status group, verification, date ranges, parent/child.
- Reports: counts by project/building/location/category/status; assets inside units vs mall common areas; unverified count; discrepancies by campaign/status; age/warranty approaching expiry; maintenance eligibility; code/serial duplicates; movement history; import errors.
- CSV/XLSX exports with deterministic filters, audit and column-level permission redaction; protect against formula injection in exported spreadsheets. For large sets use background export jobs with ownership checks, expiry and streaming.
- QR label generation: code + human-readable administrative code, optional asset name and project; safe reprint; generated QR must not grant public asset access.
- Basic response/reporting performance budgets to be decided with baseline data; avoid prematurely promising specific performance or scale numbers without benchmarks.

## 3. Core security, permission and approval contract

### 3.1 Authorization model

Combine RBAC + explicit permissions + data scopes (project, location subtree, category subtree, asset set) + optional direct grants/denials with deterministic precedence. **Default deny** and evaluate in Go on every API resource, including linked objects in nested requests and export endpoints.

Identity tables separate `employees` (staff profile), `users` (authentication), `departments`, `teams`, `team_memberships`, `project_memberships`, `roles`, `permissions`, `role_permission_grants`, `user_role_assignments`, `user_permission_overrides`, `resource_scopes`.

Representative permission names (catalog is extensible, never free-form privilege escalation):

```text
projects.view/create/update/archive
locations.view/create/update/move/archive
units.view/create/update/link_locations/archive
categories.view/manage
attributes.view/manage
manufacturers.view/manage
asset_templates.view/manage
assets.view/create/update/recode/transfer/retire/dispose
assets.attachments.read/upload/delete
assets.history.view
asset_movements.create/review/approve/execute/configure
maintenance_eligibility.view/configure
inventory_campaigns.create/assign/activate/close
inventory_observations.capture/sync/view
inventory_discrepancies.review/approve/reject/resolve
imports.create/preview/commit
reports.view/export
qr.print
users.view/manage
teams.view/manage
roles.view/manage
approval_policies.view/manage
security.audit.view
```

- Permission scope can narrow by project, subtree, category, or explicit resources; support cross-project employees without duplicated accounts.
- Unit/location membership does not imply access to private tenant records. All sensitive financial documents require separate explicit visibility grants.
- `Super Admin` centrally manages roles/workflows, but all privileged activity is audited; protected bootstrap and at least one controlled break-glass procedure.
- Users must not grant permissions broader than their own delegated authority, except audited Super Admin bootstrap.

### 3.2 Approval engine

- Policy key: `module + operation + project + optional category/criticality/conditions` with priority and conflict-resolution semantics.
- Policies may specify one or more approval stages, roles/teams/users, escalation, separation of duties, effective revisions and expiry.
- Proposal and decisions should be append-only facts; only transition the command after optimistic concurrency validation.
- Record: requester, subject, policy version, mandatory evidence, stage assignments, who decided, timestamp, reason and outcome.
- **Per-discrepancy individual decision is mandatory** irrespective of configurable general workflow. No catch-all auto-approve for physical inventory discrepancies.
- Avoid a homegrown arbitrary workflow scripting language in V1; provide declarative, whitelisted configurations and explicit transition guards.

### 3.3 Authentication/security design decision gate

Recommend **revocable server-side opaque sessions** for future browser access (secure HttpOnly SameSite cookies + CSRF defenses), and separately scoped, short-lived **service/API tokens** only if needed for CLI/integration testing. Do not ship insecure long-lived static tokens. Alternatives (OIDC/SSO, bearer tokens) require Phase 0 ADR. Bruno and Scalar must be usable in a safe development authentication flow.

Minimum controls: TLS, secure password hashing if local accounts (Argon2id), MFA for privileged accounts where viable, secret rotation, auth rate limiting, audit access, request/body/file size limits, input validation, least privilege, CSRF/CORS policy, no secrets in logs, redaction of sensitive values, strict server-side object authorization, sanitization of generated docs, security headers, trace IDs.

## 4. Reference architecture and repository layout

**Proposed:** Go single deployable **Modular Monolith** with separated domains, ports/adapters where justified, synchronous application service calls and transaction boundaries; do not start with microservices, event bus or distributed consistency unless evidenced by actual needs.

**Proposed infrastructure:** PostgreSQL, structured migrations, pgx + sqlc (candidate, evaluate against alternatives), standard net/http or proven lightweight router (evaluate), private object storage adapter, optional durable jobs/outbox for exports/virus scans. No frontend service.

```text
repo/
├── README.md
├── AGENTS.md                         # Codex rules; generated from this contract
├── go.mod
├── cmd/
│   ├── api/main.go
│   └── worker/main.go                 # only if required
├── internal/
│   ├── platform/
│   │   ├── config/   db/   httpx/   logging/   errors/
│   │   ├── authn/    authz/ approvals/ audit/
│   │   ├── files/    jobs/  clock/   idempotency/
│   ├── organization/                 # people, teams, departments
│   ├── projects/
│   ├── locations/
│   ├── units/                        # basic registry only
│   ├── catalog/                      # categories, attributes, manufacturer/model
│   ├── assets/                       # identity, status, tree, placement
│   ├── assetmovements/
│   ├── imports/
│   ├── inventory/                    # campaigns, scans, discrepancies
│   └── reports/
├── api/
│   └── openapi.yaml                  # contract is source of truth
├── docs/
│   ├── architecture/
│   ├── decisions/                    # numbered ADRs
│   ├── security/
│   ├── database/
│   ├── runbooks/
│   └── phases/                       # one signed-off report per phase
├── db/migrations/
├── db/queries/                       # if sqlc chosen
├── tests/
│   ├── integration/
│   ├── contract/
│   ├── hurl/
│   └── fixtures/
├── bruno/                             # versioned collection, no secrets
├── deployments/
│   └── local/                         # local docker compose only at first
├── Makefile
├── .env.example
└── .github/workflows/ci.yml
```

**Module boundary rules:** domain services own business rules; HTTP handlers are adapters only; repositories may not leak SQL-specific types into domain; dependencies point inward; modules must not import other modules' private repositories or bypass application services. Standard audit & authorization hooks are shared, but business invariants stay owned by their domain.

**Concurrency & consistency:** DB transactions for all mutable commands, foreign keys, unique/partial unique indexes, version/ETag or expected revision for updates, no check-then-write races, idempotency for multi-step import/sync/execute endpoints; row locking or serialization where necessary.

**Database portability is not a goal**; use PostgreSQL well, but avoid unsanctioned exotic extensions. Prefer soft archive over hard delete where historical data is referenced. Avoid generic EAV everywhere; use normalized core + typed dynamic schema configuration + JSONB values.

## 5. Initial logical data model (ERD blueprint; finalize in Phase 1)

The following is an **entity checklist**, not migration-ready DDL. Codex must produce a real ERD, cardinalities, indexes, unique constraints, DELETE/ON UPDATE behavior, sample data and concurrency rationale before writing irreversible migrations.

| Domain | Suggested entities and key relationships |
|---|---|
| Organization | `companies`, `departments`, `employees`, `users`, `teams`, `team_memberships`, `project_memberships` |
| Access | `roles`, `permissions`, `role_permissions`, `user_role_bindings`, `user_permission_overrides`, `resource_scopes`, `auth_sessions`, `api_tokens` (only if adopted) |
| Approvals | `approval_policy_versions`, `approval_policy_stages`, `approval_requests`, `approval_decisions`, `approval_events` |
| Projects | `projects`, `project_code_policies`, `project_code_sequences`, `asset_code_claims` |
| Physical spaces | `location_types`, `locations`, `unit_types`, `units`, `unit_locations` |
| Catalog | `asset_categories`, `attribute_definitions`, `attribute_options`, `category_attribute_rules`, `category_schema_versions`, `manufacturers`, `models`, `asset_templates`, `asset_template_versions` |
| Assets | `assets`, `asset_dynamic_values` (JSONB or JSONB within asset, ADR), `asset_code_history`, `asset_parent_links_history`, `asset_placements`, `asset_placement_events`, `asset_status_groups`, `asset_status_values`, `asset_status_assignments`, `asset_status_events`, `asset_verifications` |
| Team responsibilities | `asset_responsibility_rules`, `asset_responsibility_assignments`, `category_team_eligibility_rules`, `project_team_eligibility`, `asset_team_eligibility_overrides` |
| Movements | `movement_types`, `movement_type_versions`, `asset_movement_requests`, `asset_movement_events`, `asset_movement_executions` |
| Imports | `import_batches`, `import_column_mappings`, `import_rows`, `import_issues`, `import_commit_results` |
| Physical inventories | `inventory_campaigns`, `inventory_campaign_scopes`, `inventory_assignments`, `inventory_baseline_entries`, `inventory_manifests`, `inventory_sync_batches`, `inventory_observations`, `inventory_discrepancies`, `discrepancy_decisions`, `discrepancy_corrections` |
| Files | `attachments`, `attachment_links`, `upload_sessions`, `file_scan_results` |
| Audit & jobs | `audit_events`, `outbox_events` (only if justified), `export_jobs`, `qr_print_jobs` |

### 5.1 Required relation/constraint examples

- `locations.parent_id -> locations.id`; `locations.project_id -> projects.id`; service checks same-project ancestor and prevents cycles; DB-level protections where feasible.
- `units.project_id -> projects.id`; `unit_locations(unit_id, location_id)` unique; enforce linked location belongs to same project. At most one primary by partial unique index; **exactly one** when active enforced by controlled transactional activation rule.
- `assets.category_id -> asset_categories.id`; project relationships explicit. `assets.internal_id` immutable; administrative code unique **within code-claim scope**, historical claims never recycled. `assets` must not reference private tenant assets.
- Asset placement uses a single active placement discriminated as `UNIT` or `LOCATION`; `UNIT` requires unit ID, optional constrained unit-linked location; `LOCATION` requires location ID and no unit. Enforce a single active placement and project consistency.
- `asset_parent_links_history` effective dates: exactly zero or one active parent for a child; no cycles; prohibit self-parent and unexpected cross-project parent.
- `asset_status_assignments`: for single-select group max one active value; multiple values for multi; deactivated values cannot be newly assigned; preserve historical names/config revisions.
- `movement_types` in use may be archived but never physically deleted; `movement_type_versions` immutable when referenced.
- `inventory_baseline_entries` frozen on activation; `inventory_observations` append only; `inventory_sync_batches` idempotent; `inventory_discrepancies` have one independent authoritative final decision record per decision version; any correction linked to a real asset command.
- Use UTC timestamps in storage; render user-local times in future clients (Cairo by default for project display), include `created_by`, `updated_by`, revision, and consistent lifecycle semantics. Never treat wall-clock time as a unique ID.

### 5.2 Suggested indexes/constraints to justify in ERD

- `projects(code_normalized)` unique.
- `locations(project_id,parent_id)`, `(project_id,code_normalized)` unique, subtree query optimized.
- `units(project_id,unit_code_normalized)` unique; `unit_locations(unit_id,is_primary)` partial unique for primary; FK project consistency.
- `asset_code_claims(code_scope_id, code_normalized)` unique forever; **reserve code transactionally**.
- `assets(current_project_id, category_id)`, `(home_project_id, created_at)`, asset serial normalized lookup without false universal uniqueness.
- `asset_placements(asset_id) WHERE active=true` unique.
- `asset_parent_links_history(child_asset_id) WHERE valid_to IS NULL` unique.
- `asset_status_assignments(asset_id,group_id) WHERE active=true` uniqueness appropriate to group's single-select mode (enforce via service + DB strategy).
- `inventory_observations(campaign_id,client_session_id,client_observation_id)` unique.
- `approval_requests(subject_type,subject_id,policy_version_id)` reference appropriately and avoid duplicate pending stage resolutions.
- Add selective JSONB indexes *only* for declared searchable attributes actually queried.

## 6. API contract design (API-first; examples are illustrative until OpenAPI is committed)

### 6.1 Standards

- Base path: `/api/v1`; JSON UTF-8; resource UUID IDs; consistent `snake_case` fields (confirm ADR); documented error codes and validation details.
- OpenAPI **3.1.x** is canonical (latest stable patch in supported family during implementation); use compatibility tooling and lint it. API server must implement contract or mark route intentionally future/not implemented.
- No fake 'working' responses for unimplemented business endpoints; return explicit `501` only in isolated experimental branches and never claim pass.
- Pagination: cursor preferred for changing datasets, or stable page pagination where safe; constrain page size, ordering, filtering. For bulk/import return per-item status.
- Auth required for all business routes; `/health/live`, `/health/ready`, and optionally protected docs are exceptions with documented semantics. No public production mutation sandbox.
- Standard success envelope recommended: `{ "data": { ... }, "meta": { "request_id": "..." } }` (list adds cursor); standard error envelope: `{ "error": { "code": "...", "message": "...", "details": [] }, "meta": { "request_id": "..." } }`.
- Status codes: `200` read/updated, `201` created, `202` accepted async, `204` delete/no content, `400` syntactic, `401` auth, `403` forbidden, `404` not found including concealed resource where appropriate, `409` conflicts/revision, `422` business validation, `429` throttled.
- For critical write/execute/sync/import calls: `Idempotency-Key` and conditional version handling (`If-Match` or `expected_version`) documented and tested.
- OpenAPI includes operationId, tags, permission/approval notes, errors, JSON examples, request body validation, pagination, auth, rate-limit constraints and versioning policy.

### 6.2 Endpoint inventory (proposed; refine and approve during Phase 1)

| Module | Key API endpoints/methods |
|---|---|
| System | `GET /health/live`, `GET /health/ready`, `GET /api/v1/me`, auth/session endpoints as chosen |
| Projects | `GET/POST /projects`, `GET/PATCH /projects/{projectId}`, `GET/PATCH /projects/{projectId}/code-policy` |
| Locations | `GET/POST /projects/{projectId}/locations`, `GET /projects/{projectId}/locations/tree`, `GET/PATCH /locations/{id}`, `POST /locations/{id}/move`, `POST /locations/{id}/archive` |
| Units stub | `GET/POST /projects/{projectId}/units`, `GET/PATCH /units/{id}`, `PUT /units/{id}/locations`, `GET /units/{id}/assets` |
| Organization | `GET/POST /departments`, `/employees`, `/teams`, team/project membership and responsibility endpoints |
| Authz | `GET /permissions`, CRUD `/roles`, role bindings, scoped assignments, read effective permissions, approval policy versions & decisions |
| Catalog | CRUD `/asset-categories`, `/attributes`, `/manufacturers`, `/models`, `/asset-templates`, versions/preview effective schema |
| Assets | `GET/POST /assets`, `GET/PATCH /assets/{id}`, `GET /assets/{id}/history`, `GET /assets/{id}/children`, `POST /assets/{id}/reparent`, `GET/PUT /assets/{id}/maintenance-eligibility`, status transitions, code-recode command |
| Movement config | `GET/POST /movement-types`, `GET/PATCH/DELETE /movement-types/{id}` (DELETE only if unused), `POST /movement-types/{id}/archive` |
| Movements | `GET/POST /asset-movements`, `GET /asset-movements/{id}`, `POST /asset-movements/{id}/submit`, `POST /asset-movements/{id}/execute`, `POST /asset-movements/{id}/cancel` |
| Approval | `GET /approval-requests`, `GET /approval-requests/{id}`, `POST /approval-requests/{id}/decisions` |
| Imports | `POST /imports/assets`, `GET /imports/assets/{id}`, `PUT /imports/assets/{id}/mapping`, `POST /imports/assets/{id}/validate`, `POST /imports/assets/{id}/commit`, `GET /imports/assets/{id}/issues` |
| Inventory | `GET/POST /inventory/campaigns`, `GET/PATCH /inventory/campaigns/{id}`, activation/assignment/finalize commands, `GET /inventory/campaigns/{id}/manifest`, `POST /inventory/sync-batches`, `GET /inventory/discrepancies`, `POST /inventory/discrepancies/{id}/decisions`, `POST /inventory/discrepancies/{id}/propose-correction` |
| Files | `POST /attachments/upload-sessions`, authorized completion, `GET /attachments/{id}`, `GET /attachments/{id}/download`, linking/unlinking endpoints |
| Reporting | `GET /reports/assets/summary`, other read-only reports, `POST /exports`, `GET /exports/{id}`, `POST /qr-label-jobs` |
| Audit | `GET /audit/events` restricted to audit permission with redacted fields |

**Conventions note:** Endpoint and resource names above are **planning inventory**, not signed OpenAPI schemas. Avoid prefix ambiguity; generate one full request/response schema for every implemented route. Atomic business commands (`move`, `recode`, `approve`, etc.) should not be hidden as unrestricted generic `PATCH` fields.

### 6.3 Example: create asset in a unit (illustrative)

`POST /api/v1/assets` authenticated; requires `assets.create` in the chosen project and allowed category and unit.

```json
{
  "name": "Hikvision IP Camera",
  "home_project_id": "00000000-0000-0000-0000-000000000101",
  "category_id": "00000000-0000-0000-0000-000000000102",
  "administrative_code": null,
  "serial_number": "SERIAL-EXAMPLE-123",
  "template_id": "00000000-0000-0000-0000-000000000103",
  "placement": {
    "type": "UNIT",
    "unit_id": "00000000-0000-0000-0000-000000000104",
    "exact_location_id": null,
    "position_description": "Above false ceiling"
  },
  "attributes": { "ip_address": "192.0.2.10", "resolution": "4MP" }
}
```

Illustrative response (IDs/dates are placeholders; no running API is implied):

```json
{
  "data": {
    "id": "00000000-0000-0000-0000-000000000105",
    "administrative_code": "TGP-AST-000125",
    "name": "Hikvision IP Camera",
    "placement": { "type": "UNIT", "unit_id": "00000000-0000-0000-0000-000000000104", "exact_location_id": null },
    "verification_status": "NOT_VERIFIED",
    "version": 1
  },
  "meta": { "request_id": "illustrative-request-id" }
}
```

**Examples policy:** All sample addresses, IDs and codes are test-only. Never treat them as real Gate Plaza records. A direct location asset uses `placement.type = LOCATION` plus `location_id` and no `unit_id`.

### 6.4 Example: sync offline observations (NO official asset updates)

```json
{
  "campaign_id": "00000000-0000-0000-0000-000000000201",
  "client_session_id": "session-2026-001",
  "manifest_revision": 3,
  "observations": [
    {
      "client_observation_id": "00000000-0000-0000-0000-000000000202",
      "asset_id": "00000000-0000-0000-0000-000000000105",
      "observed_at": "2026-10-31T10:30:00Z",
      "observed_unit_id": "00000000-0000-0000-0000-000000000104",
      "observed_location_id": null,
      "notes": "Visible and accessible"
    }
  ]
}
```

Return a `batch_id`, per-item accepted/duplicate/rejected results, server received timestamp and **observation or discrepancy IDs only**; never claim a placement changed. A separate reviewed command handles later corrections. Document batch limits, retry behavior, assignment expiry, revocation, evidence upload and conflict codes.

## 7. Interactive documentation and API testing strategy

**Single source of truth:** `api/openapi.yaml` (OpenAPI 3.1 family). Rendering must be replaceable: Scalar is **initial recommended renderer**, Swagger UI remains an alternative; do not put business correctness in UI-specific decorators.

- Serve local, access-controlled interactive reference at `/docs` backed by `/openapi.yaml` or `/openapi.json` with safe environment restrictions; demonstrate `Try it` for local development; protect production docs or disable interactive mutation per deployment policy.
- API reference per operation: purpose, access scope, request/response example(s), schema, headers, version conflicts, approval requirements, error examples, curl snippet, link to corresponding integration test.
- Commit **Bruno** collections in Git: environments `local` (and template staging with no secrets); every business area has success and denial/failure flows. Secrets belong in ignored local env/secret manager, never repository.
- Build **Go unit/integration tests** as primary quality evidence; `Hurl` optional/recommended for network-level deterministic smoke and contract tests; avoid maintaining two contradictory hand-written contracts.
- CI: format/vet/static checks, go tests/race where feasible, migration up/down dry run on empty ephemeral PostgreSQL, OpenAPI lint validation, API handler contract checks, privilege/negative tests, import/offline replay tests and security smoke.
- Generate docs and collections only from committed source/contracts where supported; compare generated output in CI or document intentional generation step.
- Provide `make setup`, `make migrate-up`, `make test`, `make test-integration`, `make lint`, `make openapi-lint`, `make api-smoke`, `make docs` after implementation. If a command cannot run, report reason and whether gate is blocked.

## 8. Nonfunctional requirements

### 8.1 Reliability and correctness

- Atomic domain mutations; serializable/locked operations for code allocation, transfers, nested asset relations and discrepancy resolution; test contention.
- Every imported or synced batch idempotent and resumable; no duplicate asset creation on retries.
- No data loss on job retries; transactional outbox considered only for external side effects; never send approval/notification before durable commit.
- Make normal deletes soft archives when referenced; events/approved facts are append-only. Audit log is not a replacement for domain history.
- Define recovery/backup strategy for PostgreSQL and private files before real data. Migrations must be reversible where possible and never remove historical data silently.

### 8.2 Security and privacy

- All project/asset/attachment endpoints authorize identity + function + scope + resource at the application and query level. Include negative tests for ID enumeration and cross-project references.
- Input size, rate limits, upload scan, secret handling, secure errors and least-privileged DB identities.
- Audit mutations with actor, subject, before/after relevant deltas (redacted), justification, source/device, request ID, time.
- Redact personal records and documents; log no passwords/tokens/secret attachment URLs. Store credentials using modern recommended KDF and user session rotation.
- Protection against spreadsheet formula injection, CSV/JSON injection, path traversal, malicious uploads, XSS in docs, mass assignment and SQL injection via parameterized queries.

### 8.3 Performance and observability

- Use bounded pagination, query plan tests, reasonable defaults/limits for tree depth, imports, request bodies and sync payloads.
- Structured logs, health endpoints, metrics for request latency/errors, job failures, import processing and sync conflicts, and trace correlation.
- Establish baseline load/volume assumptions from sample legacy data before committing SLA/throughput numbers.

### 8.4 Operational/deployment boundaries

- Local reproducible environment with disposable PostgreSQL and optional private object-store emulator; env variables documented.
- Production hosting and networking are **undecided**; no Hostinger/shared-host assumptions from other unrelated projects. No production migration/deploy without explicit approval and backup plan.
- API docs/test tooling must run locally without frontend app. Keep dependency lockfiles/Go module verification and pinned CI images.

## 9. Test matrix and acceptance scenarios

All key business invariants have both positive and failure tests, with table-driven Go tests and DB integration tests where required. Minimal scenarios to implement across phases:

| ID | Required scenario / observable acceptance |
|---|---|
| T01 | Employee serves two projects via one identity; rights differ by project; zero access to third project. |
| T02 | Create arbitrarily nested locations; forbid cycle, cross-project parent and deleting referenced location. |
| T03 | Unit links to 2+ locations with one primary; impossible to activate without exactly one primary or link across projects. |
| T04 | Create asset in unit without exact location: unit known; no falsely inferred floor; exact linked location accepted; unrelated location rejected. |
| T05 | Core fields always validated; dynamic schema inherited; invalid option/type rejected; changes to schema do not corrupt existing assets. |
| T06 | Template versions do not alter previously created assets; bulk import detects duplicates and malformed location/category links. |
| T07 | Generate/manual code collision handled transactionally; code never reused, even if original asset retired/deleted. |
| T08 | Optional serial per category, configurable duplicate policy; QR grants no access to unauthorized user. |
| T09 | Parent asset move evaluates inherited vs overridden children, prevents cycles and partial writes. |
| T10 | Multiple status dimensions coexist; single-select group forbids two active states; history retained on status retirement. |
| T11 | Maintenance allowed teams inherit category and project eligibility; per-asset add/exclude/replace works; no task creation. |
| T12 | Create/configure custom movement type safely; in-use type cannot hard-delete; wrong primitive or bypass rejected. |
| T13 | Maker cannot approve their own restricted action; stages/policy versions honored; stale execute returns conflict; no premature asset mutation. |
| T14 | Asset imports via Excel preview then approved commit; no mutation on preview; repeated commit idempotent; row errors visible. |
| T15 | Offline manifest limits to inspector's scope; repeat sync batch idempotent; out-of-order/revoked/closed campaign behavior explicit. |
| T16 | Sync matching or mismatching observations never changes official asset row or location. |
| T17 | Two independent discrepancies require two individual decisions, not one bulk approval; approved mismatch alone leaves official asset unchanged. |
| T18 | Authorized correction creates a separate audited command linking discrepancy and approval; repeat/replay cannot duplicate correction. |
| T19 | Import/discrepancy/file downloads enforce permissions, prevent cross-project disclosure and sanitize file names. |
| T20 | Private attachments enforce MIME/size/checksum and read authorization; downloads cannot be enumerated publicly. |
| T21 | Reports count unit-bound assets once (even if unit has multiple linked locations), prevent double counting across location rollups. |
| T22 | OpenAPI routes/errors match implementation, Scalar local try-it works, Bruno/Hurl smoke tests pass with non-admin failure cases. |
| T23 | CI runs DB migrations and tests on clean checkout and reports exact failure rather than hiding flaky errors. |

For each phase, include migration tests, project-scope permission denial, audit trail and concurrency tests relevant to new work. **A feature is not Done until its security/negative cases and API contract are demonstrated.**

## 10. Implementation roadmap — strict gated phases for Codex

**General phase rule:** Each phase begins from approved prior code. Codex must not proceed to the next phase automatically. Every phase ends with a *phase report*, exact changed-file list, tests run with pass/fail, unresolved decisions, risks and request for approval. Feature flags may guard incomplete API surfaces, never simulate functioning features.

### Phase 0 — Requirements lock, technical ADRs and project skeleton

**Deliverables:** translate this file into traceable requirement IDs; create ADRs for Go version/router, PostgreSQL version, DB access/migrations, IDs/codes, auth transport, OpenAPI code-first/contract-first workflow, object storage, CI, statuses and movement primitives; Docker local environment; minimal healthy Go server, lint/test commands, structured logs, README/AGENTS.md. No full domain implementation.

**Exit tests:** clean clone bootstrap; DB connectivity and migrate up/down; `go test ./...`; `go vet ./...`; OpenAPI skeleton lints; server health endpoints; no committed secrets. **STOP.**

### Phase 1 — Relational schema, ERD and core domain invariants

**Deliverables:** formal ERD, migration numbering, relationships and uniqueness/foreign keys, seed strategy, data ownership rules and transactional plans. Implement stable foundation tables + repository interfaces with tests; define approval and audit event shape. No unreviewed mass schema creation.

**Exit tests:** fresh DB round-trip migrations, FK violation cases, unique code reservation contention, project link integrity, documented ERD and ADR reconciliation. **STOP.**

### Phase 2 — Authentication, Super Admin, permissions, organization and teams

**Deliverables:** protected bootstrap, auth/session (per ADR), departments/users/employees/teams/memberships, project memberships, scoped RBAC and explicit permissions, group roles, effective permission explain API, audit logs, seed test users with zero production secrets.

**Exit tests:** unauthorized, role scope denial, cross-project access, least privilege, privilege-escalation denial, revocation, token/session protections; Bruno login flow. **STOP.**

### Phase 3 — Projects, locations and minimal unit registry

**Deliverables:** CRUD/projects; tree CRUD/reparent/cycle prevention; units independent from tree; multi-location link and one primary; unit and location read filters; safe archiving.

**Exit tests:** T01-T04, cross-project refs, ancestry queries, moving referenced subtree, unit primary validation, API samples. **STOP.**

### Phase 4 — Categories, dynamic attributes, manufacturer/models, templates

**Deliverables:** category tree, inherited schemas, typed attributes, validation, enum sets, manufacturer/model, template versioning, search configuration, schema-change safety report.

**Exit tests:** T05-T06; old assets/schema compatibility via fixtures; invalid conditional settings, cycles, option deletion and category archive protections. **STOP.**

### Phase 5 — Asset registration, identity, QR, placements, hierarchy and files

**Deliverables:** core+JSONB attribute storage, code generators/manual code, serial policy, QR stable lookup, placement in unit/direct location, parent–child, status dimensions, core history, attachments, warranty/purchase metadata.

**Exit tests:** T04, T07-T10, T20, code generation concurrency, location precision, parent asset relocation rules, gated file access. **STOP.**

### Phase 6 — Maintenance eligibility and responsibility mapping

**Deliverables:** department responsibility and rules; team skills/project eligibility; category defaults and per-asset overrides; effective team explanation endpoint. **No work orders**.

**Exit tests:** T11, project scoping, override precedence, inheritance and multi-team scenarios. **STOP.**

### Phase 7 — Movement type configuration, approvals and command execution

**Deliverables:** movement catalog CRUD + immutable versions, fixed safe primitives, generic approval stages/policies, maker–checker, safe executing business commands with audit, cross-project movement policy, child impact previews.

**Exit tests:** T09, T12-T13; unapproved state unchanged, stale/version conflict, rollback on any failure, in-use type archive vs deletion. **STOP.**

### Phase 8 — Excel migration/import, manual/bulk intake

**Deliverables:** Excel mapping, dry-run/preview/validation, duplicate matching, import-batch state machine, commit approvals, template bulk-create, external legacy references, no duplicate on retries.

**Exit tests:** T06/T14, malformed files, large bounded imports, partial error reports, audit, idempotent commit, rollback rules. **STOP.**

### Phase 9 — Physical inventory campaign and offline sync server

**Deliverables:** baseline snapshot/frozen manifest and scope assignments; observation APIs and batch sync; duplicates/retry/resume; discrepancy detection & individual reviewer decisions; separate corrective transaction handoff; reporting; no mobile PWA.

**Exit tests:** T15-T18; prove via DB assertions that even approved discrepancies do not mutate assets; device replay, lost connectivity simulation, stale manifest and multi-observation conflicts. **STOP.**

### Phase 10 — Reports, exports, QR labels, observability and hardening

**Deliverables:** filtered reporting, dynamic attribute query safeguards, PDF/CSV/XLSX when confirmed, safe exports, QR label jobs, metrics/logging, backups/restore instructions, performance baselines and security scan plan.

**Exit tests:** T19-T21, authorization/PII, code label validity, query indexes, high-cardinality list paging, spreadsheet injection. **STOP.**

### Phase 11 — End-to-end API quality, documentation and V1 release candidate

**Deliverables:** reviewed OpenAPI 3.1 full contract; Scalar interactive local docs; Bruno collections; Hurl/Go integration+contract suite; CI green; reproducible local demo dataset; signed phase reports; release notes; rollback/migration & disaster recovery runbooks, remaining limitations and future frontend/API integration notes.

**Exit tests:** T01-T23 green; OpenAPI valid and synchronized, local real request/response demos, no endpoints falsely represented as finished, security matrix pass, build reproducibility. **STOP for explicit release/deployment authorization.**

### 10.1 Definition of Done (every phase)

- Requirements implemented exactly or deviation documented in ADR and approved.
- Migrations isolated, forward rollback plan documented, tests use ephemeral PostgreSQL.
- New APIs have OpenAPI request/response/error examples and permissions documented.
- Security and cross-project denial tested; no unsafe mass assignment/secret leaks.
- `gofmt`, `go vet`, tests, chosen lint, DB tests and OpenAPI lint pass; tool versions pinned in repository.
- Full phase report provided with **what changed**, **test evidence**, **what is not done**, **risks** and **user approval requested**.
- No upcoming phase work initiated on the assumption that silence equals approval.

## 11. Codex operating instructions / reusable prompt

Copy the following into a Codex task, changing only the phase number when that phase has been approved:

```text
You are implementing the backend of a multi-project Facility Management Platform in Go.
Read ASSET_MANAGEMENT_BACKEND_MASTER_PLAN.md completely before changing files.
Follow its invariants, scope boundaries, architecture guidance, acceptance matrix,
security requirements, and phase gates. The business user has approved the
functional scope, not every unspecified technical design choice.

CURRENT PHASE: PHASE 0 ONLY.

1. Inspect the repository and identify pre-existing files and conventions.
2. State a concise phase plan with acceptance criteria.
3. If a technical decision is undecided, write an ADR with a recommended option
   and alternatives; never silently make a consequential decision.
4. Implement only the current phase in small, reviewable changes.
5. Update documentation/OpenAPI/tests as the code evolves.
6. Run every relevant test or report exactly why it cannot run.
7. Produce docs/phases/phase-00-report.md with changed files, commands/results,
   limitations, risks, and questions requiring explicit approval.
8. STOP. Do not start Phase 1 or build a frontend or deploy to production.

Never let inventory observations or approvals directly edit official assets.
Do not infer unit location precision. Enforce cross-project scoped authorization.
Never hard-delete referenced movement types, statuses, or history.
```

**Codex phase output template:**

```markdown
# Phase N — Delivery Report
- Commit/branch: ...
- Requirements addressed: ...
- ADRs created/changed: ...
- Files added/modified: ...
- Migrations: ...
- API/OpenAPI operations: ...
- Tests executed (exact commands and PASS/FAIL): ...
- Negative/security tests: ...
- Manual demo steps: ...
- Open blockers and deviations: ...
- Data safety / rollback notes: ...
- Decision requested from owner: APPROVE / REVISE PHASE N
- Next phase: NOT STARTED
```

## 12. Requirements traceability and change control

Create `docs/requirements/traceability.md` mapping each feature ID to endpoint, domain service, DB table(s), tests and phase. Start with the following IDs:

| Requirement IDs | Group |
|---|---|
| CORE-001.. | Company, projects, employees/teams, cross-project roles |
| LOC-001.. | Location tree, unit independence and multi-link |
| CAT-001.. | Category schema, attributes and template versions |
| AST-001.. | Asset identity, placements, parent–child, statuses, files |
| ELI-001.. | Maintenance eligible team rules |
| MOV-001.. | Configurable movement types and safe movement execution |
| APR-001.. | Permission policy and individual approvals |
| IMP-001.. | Excel/bulk/manual registration and verification boundary |
| INV-001.. | Offline campaign manifests, observations, discrepancy review, separate correction |
| RPT-001.. | Reports, exports and QR labels |
| API-001.. | Contract, doc interactive test workspace, error conventions |
| SEC-001.. | Scoped authorization, audit, security and recovery |

**Change control:** user-approved scope is a baseline. If discovery reveals contradictions, document alternatives and request decision; do not silently reinterpret. Future modules integrate through explicit APIs/domain events, not direct table edits. Any new feature outside V1 is added to a future backlog unless the owner explicitly changes V1 scope.

## 13. Open design items for Phase 0 ADRs (not blocking SRS baseline)

1. API authentication mechanism for browser-bound future and API test clients; MFA/SSO roadmap.
2. Database migrations and query tooling (pgx/sqlc recommended, verify fit), project identifiers, code normalization/scoping and project moves.
3. Project-to-company relationships and exact staff eligibility/role precedence, direct permission-deny semantics.
4. Effective config revision rules for attributes/statuses/movements and schema migrations, including approve stage changes.
5. Private file store and virus scan integration for local and future deployment.
6. Full idempotency/concurrency policy for imports/offline campaigns and how observation conflict resolution is surfaced.
7. Practical max tree depth, payload sizes, error/retry budgets, report/export volume and benchmark target.
8. Production hosting, backups, monitoring provider and deployment approvals; no deployment assumptions until selected.
9. Whether financial purchase cost is restricted by special per-field permission; recommended YES.
10. Representation of administrative `home_project_id` versus physical `current_project_id` and re-code policy on cross-project moves.

## 14. Reference materials (engineering guidance, not endorsements)

- Go module/server organization: https://go.dev/doc/modules/layout
- OpenAPI Specification 3.1 family: https://spec.openapis.org/oas/v3.1.1.html
- OpenAPI official releases/index: https://spec.openapis.org/oas/
- Scalar interactive API reference: https://scalar.com/
- Bruno API client and versioned collections: https://www.usebruno.com/
- Hurl HTTP test runner: https://hurl.dev/docs/manual.html
- OWASP API Security Top 10: https://api-security.owasp.org/editions/2023/en/0x11-t10/
- PostgreSQL reference: https://www.postgresql.org/docs/current/

---

## 15. Owner-facing Arabic checkpoint / ملخص الاعتماد

**المعتمد:** المشروع يبدأ بمول The Gate Plaza لكنه مصمم للمشروعات المتعددة. الوحدات منفصلة عن شجرة المواقع ومرتبطة بها. كل أصل تابع للمول فقط، ويمكن أن يكون داخل وحدة أو في منطقة بالمول. تصنيفات وخصائص وقوالب وحالات وحركات قابلة للإدارة. يوجد كود إداري وسيريال وQR، وأصول رئيسية وفرعية، وأكثر من فريق صيانة مؤهل حسب النوع أو الأصل. الموظف الواحد يعمل على أكثر من مشروع بصلاحيات محددة. الجرد يدعم الـOffline من ناحية الـAPI، والمزامنة لا تعدّل الأصول؛ كل اختلاف له مراجعة وموافقة فردية ثم إجراء تصحيحي مستقل إذا لزم.

**قرار التنفيذ:** Go Backend + PostgreSQL + API-first مع توثيق OpenAPI تفاعلي واختبارات فعلية؛ **تأجيل الـFrontend بالكامل**. يبدأ Codex بـ **Phase 0 فقط** ويقف بعد تقديم تقرير الاختبارات والمخرجات إلى أن تعتمد الانتقال للمرحلة التالية.

**Document integrity:** This file is a *plan/requirements artifact*, not proof that repository files, database migrations, API endpoints or tests already exist.
