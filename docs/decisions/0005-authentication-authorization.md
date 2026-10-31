# ADR 0005: Authentication transport and scoped permissions

Date: 2026-10-31  
Status: Recommendation; Phase 2 implementation gate

## Context

Master plan §§0.1, 4, 8, 10 and 13 require evaluated choices, a small modular foundation and explicit unresolved decisions. No later-phase implementation is authorized.

## Decision / recommendation

Recommend revocable opaque server sessions for future browser clients, Secure HttpOnly SameSite cookies, CSRF checks, rotation and bounded lifetime. CLI/Scalar/Bruno can use a documented development session flow. No static development backdoor. Use Argon2id for local password accounts; plan privileged MFA and optional OIDC integration. Every future resource requires default-deny function + scope checks; memberships/eligibility grant no permission.

## Alternatives, consequences and risks

JWT bearer tokens ease clients but complicate revocation; OIDC delegates credential security but needs a provider. Do not implement either in Phase 0. Direct-deny precedence, delegated authority, company relationship constraints, MFA/SSO provider and session expiry must be reviewed in Phase 2 design. Recommend explicit deny overrides within matching scope, with tested deterministic resolution; not approved product behavior yet.

## Verification / revisit gate

See Phase 0 report for executed checks. Revisit before the affected domain schema/API or production operation. No ADR overrides the master plan's company/project, unit precision, inventory separation or historical integrity invariants.
