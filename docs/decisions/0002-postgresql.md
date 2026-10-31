# ADR 0002: PostgreSQL version and local environment

Date: 2026-10-31  
Status: Superseded for local configuration by ADR 0011; prior decision retained for history

## Context

Master plan §§0.1, 4, 8, 10 and 13 require evaluated choices, a small modular foundation and explicit unresolved decisions. No later-phase implementation is authorized.

## Decision / recommendation

Pin PostgreSQL 18.6-alpine for disposable local development. Use relational constraints and selective JSONB in future phases, UTC timestamptz, recursive tree queries initially. No extensions now. Loopback port and persistent named volume; no production configuration.

## Alternatives, consequences and risks

PostgreSQL 17 remains supported and is a reasonable deployment fallback after compatibility tests. SQLite cannot validate PostgreSQL locking/constraint semantics. Patch tags are pinned but not immutable digests; resolve and record digests before trusted CI/release. Latest-version claims are avoided because registry access cannot verify availability here.

## Verification / revisit gate

See Phase 0 report for executed checks. Revisit before the affected domain schema/API or production operation. No ADR overrides the master plan's company/project, unit precision, inventory separation or historical integrity invariants.
