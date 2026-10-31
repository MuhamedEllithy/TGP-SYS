# Architecture decisions

Technical adoption does not change approved functional scope. Recommendations require review before their implementation phase.

- [0001: Go, HTTP routing and modular monolith](0001-go-http-modular-monolith.md) — Adopted for Phase 0.
- [0002: PostgreSQL version and local environment](0002-postgresql.md) — Superseded for local configuration by ADR 0011.
- [0003: Database access and migration tooling](0003-db-migrations.md) — Adopted; dependency/migration host checks owner-reported PASS.
- [0004: Stable identity and administrative code scope](0004-identity-codes.md) — Recommendation; confirm before Phase 1 DDL.
- [0005: Authentication transport and scoped permissions](0005-authentication-authorization.md) — Recommendation; Phase 2 implementation gate.
- [0006: Contract-first OpenAPI and interactive reference](0006-openapi-docs.md) — Adopted; lint and host HTTP availability verified; Browser Try-it pending.
- [0007: Private object storage and scanning](0007-private-storage.md) — Recommendation; no provider selected.
- [0008: Testing and CI foundation](0008-ci-testing.md) — Adopted foundation; workflow unexecuted.
- [0009: Versioned statuses and safe movement primitives](0009-config-status-movements.md) — Architecture recommendation; no implementation.
- [0010: Consistency, observations and unresolved deployment choices](0010-consistency-open-decisions.md) — Recommendation; future phase gates.

- [0011: Dedicated PostgreSQL 17 and privilege separation](0011-dedicated-postgresql17-local.md) — adopted local configuration; supersedes 0002 version and earlier shared owner setup; host DB verification owner-reported PASS.

## Primary technical references

- [Go release downloads](https://go.dev/dl/) and [server module layout](https://go.dev/doc/modules/layout) (Go BSD-3-Clause).
- [PostgreSQL version policy](https://www.postgresql.org/support/versioning/) (PostgreSQL license).
- [pgx v5.11.0 release](https://github.com/jackc/pgx/releases/tag/v5.11.0) and [upstream changes](https://github.com/jackc/pgx/blob/master/CHANGELOG.md) (MIT).
- [golang-migrate v4.20.1](https://github.com/golang-migrate/migrate/releases/tag/v4.20.1) (MIT).
- [OpenAPI 3.1.1](https://spec.openapis.org/oas/v3.1.1.html), [Scalar configuration](https://github.com/scalar/scalar/blob/main/documentation/configuration.md) (MIT), [Redocly CLI](https://github.com/Redocly/redocly-cli) (MIT).

Sources reviewed during Phase 0; registry downloads, transitive licenses and security advisories require verification after network access is restored. Tool versions are explicit candidates until install and compatibility checks succeed; no claim that all are the latest release.

Current verification status/provenance is recorded in the [final Phase 0 report](../phases/phase-00-report.md). Host automation results do not establish manual Browser Try-it, clean-clone bootstrap or hosted CI.
