# ADR 0007: Private object storage and scanning

Date: 2026-10-31  
Status: Recommendation; no provider selected

## Context

Master plan §§0.1, 4, 8, 10 and 13 require evaluated choices, a small modular foundation and explicit unresolved decisions. No later-phase implementation is authorized.

## Decision / recommendation

Define a future private storage port with server-generated keys, DB metadata/checksum, authorized streaming or short-lived links and quarantined scan-before-release. Recommend S3-compatible provider plus a disposable local emulator only in attachment phase. No storage service or upload endpoint now.

## Alternatives, consequences and risks

Local filesystem is simple but complicates replicated deployment and backups. Managed object storage reduces operations but adds provider constraints and egress cost. Provider, encryption, retention, MIME/size policy, malware scanner and sensitive financial-field permissions remain owner/security decisions before Phase 5; recommend explicit financial-field grants. Public buckets are incompatible with the baseline.

## Verification / revisit gate

See Phase 0 report for executed checks. Revisit before the affected domain schema/API or production operation. No ADR overrides the master plan's company/project, unit precision, inventory separation or historical integrity invariants.
