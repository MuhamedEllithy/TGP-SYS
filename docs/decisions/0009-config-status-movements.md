# ADR 0009: Versioned statuses and safe movement primitives

Date: 2026-10-31  
Status: Architecture recommendation; no implementation

## Context

Master plan §§0.1, 4, 8, 10 and 13 require evaluated choices, a small modular foundation and explicit unresolved decisions. No later-phase implementation is authorized.

## Decision / recommendation

Use immutable referenced configuration revisions; retire referenced statuses and movement types instead of deleting them. Configurable type names map only to implemented Go primitives (master plan §2.7); execution independently rechecks revision, scope, approvals, placement and child impacts inside one transaction. Verification is evidence-backed, never a free status toggle.

## Alternatives, consequences and risks

An arbitrary workflow/rule scripting language risks invariant bypass and is excluded. Hardcoding all status catalogs blocks configurability. Effective revision resolution and pending-proposal behavior after policy changes require design approval before phases 4/5/7. Snapshot referenced presentation and policy context; do not rewrite history.

## Verification / revisit gate

See Phase 0 report for executed checks. Revisit before the affected domain schema/API or production operation. No ADR overrides the master plan's company/project, unit precision, inventory separation or historical integrity invariants.
