# ADR 0003: Database access and migration tooling

Date: 2026-10-31  
Status: Adopted; host dependency/migration checks owner-reported PASS

## Context

Master plan §§0.1, 4, 8, 10 and 13 require evaluated choices, a small modular foundation and explicit unresolved decisions. No later-phase implementation is authorized.

## Decision / recommendation

Use pgx v5.11.0 (MIT) with pgxpool capped at 10. Runtime uses Ping; the diagnostic command also performs bounded identity/privilege/catalog checks. Use golang-migrate v4.20.1 (MIT) with paired sequential SQL and explicit commands. Add sqlc only when Phase 1 query design exists; no ORM or custom migration engine. API never applies migrations.

## Alternatives, consequences and risks

database/sql plus pgx is useful for compatibility but adds no current benefit. GORM hides query/constraint behavior needed here. Goose is a strong alternative with Go migration support but not needed for SQL-only foundation. pgx upstream master documents a SendBatch comment-boundary security fix after 5.11.0; this scaffold never uses SendBatch. Reassess fixed release and run vulnerability scan before introducing domain queries. Earlier sandbox downloads were blocked; host resolution now supplied go.sum, and direct module verification passes.

## Verification / revisit gate

See Phase 0 report for executed checks. Revisit before the affected domain schema/API or production operation. No ADR overrides the master plan's company/project, unit precision, inventory separation or historical integrity invariants.

## Local credential update

ADR 0011 now separates bootstrap, non-superuser migration owner and restricted runtime credentials. pgx and golang-migrate choices remain unchanged; host verification is now owner-reported PASS; see the final Phase 0 report for provenance.

## Host evidence reconciliation

The owner reports dependency verification, disposable migration round trip and integration tests passed. Direct sandbox module verification, full race tests and vet also pass with the now-present dependencies. Earlier sandbox download failures are historical; no live database replay was performed during report finalization.
