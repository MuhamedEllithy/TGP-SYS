# ADR 0001: Go, HTTP routing and modular monolith

Date: 2026-10-31  
Status: Adopted for Phase 0

## Context

Master plan §§0.1, 4, 8, 10 and 13 require evaluated choices, a small modular foundation and explicit unresolved decisions. No later-phase implementation is authorized.

## Decision / recommendation

Use Go 1.26.10 (supported stable line, cached locally), standard net/http routing and log/slog JSON. A single binary composes platform packages; domain modules are added only in their approved phases. Module path is provisional gateplaza.local/facility.

## Alternatives, consequences and risks

Chi offers useful middleware but adds a dependency with no current routing need. Gin adds a distinct handler abstraction. Microservices introduce distributed transactions without evidence. Standard routing keeps the system small; add a router only with concrete requirements. Upgrade Go patch versions deliberately; the machine default 1.24.0 is not the project toolchain.

## Verification / revisit gate

See Phase 0 report for executed checks. Revisit before the affected domain schema/API or production operation. No ADR overrides the master plan's company/project, unit precision, inventory separation or historical integrity invariants.
