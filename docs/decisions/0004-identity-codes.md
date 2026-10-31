# ADR 0004: Stable identity and administrative code scope

Date: 2026-10-31  
Status: Recommendation; confirm before Phase 1 DDL

## Context

Master plan §§0.1, 4, 8, 10 and 13 require evaluated choices, a small modular foundation and explicit unresolved decisions. No later-phase implementation is authorized.

## Decision / recommendation

Recommend UUIDv7 internal IDs; keep code, serial and opaque QR separate. Recommend home_project_id for immutable original code-claim scope and current_project_id for operational placement. Project transfer does not silently recode; code-claim reservations never recycle. UUID timestamp is not an authorization mechanism.

## Alternatives, consequences and risks

UUIDv4 avoids timestamp leakage with less index locality; either is viable before first schema. Globally unique administrative codes simplify lookup but alter approved configurable scope. Exact Unicode/case normalization, code-policy edits and ownership/code-scope transfer rules require owner review before constraints are committed; no identity library or domain schema introduced.

## Verification / revisit gate

See Phase 0 report for executed checks. Revisit before the affected domain schema/API or production operation. No ADR overrides the master plan's company/project, unit precision, inventory separation or historical integrity invariants.
