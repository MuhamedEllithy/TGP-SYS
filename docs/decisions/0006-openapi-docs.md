# ADR 0006: Contract-first OpenAPI and interactive reference

Date: 2026-10-31  
Status: Adopted; lint verified; HTTP docs availability host-reported PASS; Browser Try-it PENDING

## Context

Master plan §§0.1, 4, 8, 10 and 13 require evaluated choices, a small modular foundation and explicit unresolved decisions. No later-phase implementation is authorized.

## Decision / recommendation

Use OpenAPI 3.1.1 YAML as authoritative contract, embedded in the API binary. Document only implemented probes and gated docs. Scalar 1.76.4 CDN shell is local-only, auth persistence disabled and same-origin requests restricted by CSP. Bruno probe requests supplement Go tests. Redocly CLI 2.72.1 is installed and locked; OpenAPI lint passes. Security audit must be rerun when the registry advisory endpoint is reachable.

## Alternatives, consequences and risks

Code-first annotations risk hidden contract drift. Generated server interfaces can be evaluated when business schemas exist; generation is unnecessary for four system routes. Swagger UI is a renderer alternative. Self-hosting a checksum-verified Scalar bundle is preferable for offline reproducibility and removes CDN trust; current shell needs network and manual Try-it validation. Production docs are disabled, not merely hidden by a URL.

## Verification / revisit gate

See Phase 0 report for executed checks. Revisit before the affected domain schema/API or production operation. No ADR overrides the master plan's company/project, unit precision, inventory separation or historical integrity invariants.
