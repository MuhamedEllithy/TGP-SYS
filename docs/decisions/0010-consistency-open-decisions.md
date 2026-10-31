# ADR 0010: Consistency, observations and unresolved deployment choices

Date: 2026-10-31  
Status: Recommendation; future phase gates

## Context

Master plan §§0.1, 4, 8, 10 and 13 require evaluated choices, a small modular foundation and explicit unresolved decisions. No later-phase implementation is authorized.

## Decision / recommendation

Plan optimistic expected revisions plus transactional locks/uniqueness where contention requires. Imports and offline sync need durable idempotency identity and payload fingerprints. Freeze inventory baselines; observations and individual discrepancy decisions are append-only evidence. Only a separate authorized correction command may change assets, even after approval.

## Alternatives, consequences and risks

Idempotency retention, conflicting replay responses, sync skew/quarantine policy, tree depth, batch limits and performance budgets require actual data and phase-specific contracts. Do not guess throughput SLAs. Production host, TLS/network, backup/PITR, restore tests, monitoring and operational ownership remain undecided; no deployment authorized. Every item in master plan §13 is either covered here or in ADRs 0003–0007/0009 and remains traceable before implementation.

## Verification / revisit gate

See Phase 0 report for executed checks. Revisit before the affected domain schema/API or production operation. No ADR overrides the master plan's company/project, unit precision, inventory separation or historical integrity invariants.
