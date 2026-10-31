# ADR 0008: Testing and CI foundation

Date: 2026-10-31  
Status: Adopted foundation; workflow unexecuted

## Context

Master plan §§0.1, 4, 8, 10 and 13 require evaluated choices, a small modular foundation and explicit unresolved decisions. No later-phase implementation is authorized.

## Decision / recommendation

Go test/race, gofmt/vet, real PostgreSQL integration and migration round-trip, OpenAPI lint and network smoke form the initial gate. Pin exact language/tool/service versions. CI fails on missing dependency locks; no silent skip. Commit locks only after genuine resolution. No deploy job or production credentials.

## Alternatives, consequences and risks

Hosted CI is proposed GitHub Actions because no repository host exists yet; translate if owner selects another service. No remote repository was created or pushed. Dedicated staticcheck, govulncheck and secret scanner require verified pins/downloads before release. GitHub action tags are versioned but mutable; immutable commit and image digest pins remain a tracked hardening item. A clean source-copy bootstrap is required where no Git metadata exists.

## Verification / revisit gate

See Phase 0 report for executed checks. Revisit before the affected domain schema/API or production operation. No ADR overrides the master plan's company/project, unit precision, inventory separation or historical integrity invariants.
