# ADR 0011: Dedicated PostgreSQL 17 local environment

Date: 2026-10-31  
Status: Adopted; PostgreSQL/auth/migration host checks owner-reported PASS  
Supersedes: ADR 0002 local PostgreSQL 18 selection and ADR 0003 shared local owner credentials only.

## Context and decision

The owner explicitly authorized a dedicated PostgreSQL 17 development environment and identified `supabase_db_thebookseller` on port 54322 as unrelated and forbidden for application use or mutation. Use official `postgres:17.11-alpine`, project `tgp-sys-local-dev`, container `tgp-sys-postgres-dev`, network `tgp-sys-dev-network` and new persistent volume `tgp-sys-postgres17-dev-data`. PostgreSQL 17 is supported; the current minor listed by the reviewed primary support source is 17.11. No PostgreSQL 18 volume is reused or downgraded.

Three independent 256-bit random local passwords live only in ignored mode-0600 `.env`. A bootstrap administrator initializes a non-superuser migration owner (`tgp_migrator_dev`) and restricted application login (`tgp_app_dev`). The migrator owns `tgp_facility_dev` and explicitly disposable `tgp_facility_phase0_test`. The runtime login has CONNECT and schema USAGE only, no database/role/schema creation or superuser privileges. Future table-level runtime grants must be reviewed with actual Phase 1 SQL; no broad default DDL or migration-metadata write grants are added now.

Require SCRAM password authentication for host and Unix socket paths. Health check proves server availability only; `cmd/dbcheck` additionally uses the exact application host TCP URL with pgx, checks database/user/version/privileges, and requires an authentication-specific failure with a fresh wrong password. This avoids treating `pg_isready` or a container password variable as proof of authentication.

The candidate published port is 55432 on 127.0.0.1. Host preflight must verify it is free (or already correctly owned by this project); the sandbox cannot bind sockets. Before the first container exists, explicit `python3 scripts/local_env.py --select-port` selects a free port in 55432–55442 and updates only host-port-derived settings. Internal migrations use `postgres:5432`; host application/testing URLs use loopback and the chosen published port. URL userinfo is percent encoded.

## Compatibility and alternatives

The existing probe migration uses standard transactional PostgreSQL DDL, boolean checks and timestamptz available in PostgreSQL 17. pgx v5 and golang-migrate's PostgreSQL driver are the chosen adapters; host build/migration/database results are now owner-reported PASS. Direct sandbox Go checks pass; live Docker/database execution was not replayed during finalization. Image tag is pinned but a registry digest is not yet resolved. PostgreSQL 18 remains a future upgrade choice with an explicit data migration plan; shared Supabase infrastructure is rejected because it belongs to another project.

Separate runtime and migration privileges add two URLs and an initialization hook but prevent application DDL/superuser use. Hook runs only on a fresh volume. Existing credentials are never reset or silently rotated; a failed first initialization needs investigation, not volume deletion. Local bootstrap secrets are visible to users with authorized Docker inspection access; do not print full inspect/environment/config output. Production must use a dedicated secret manager and backup/restore plan.

## Safety and verification

Preflight checks container Compose labels, exact loopback port mapping and volume/network ownership; Compose starts only `postgres` with `--no-recreate`. No foreign container is adopted. Forward migration targets development DB; all rollback convenience commands target the explicitly disposable test DB and run a Go identity/object guard first. No DB, volume or network deletion command is provided.

Sandbox checks verify environment generation, no secret output, preservation/refusal cases, file permissions, Git exclusion, encoded URLs, shell/Python syntax and Compose parsing. Live checks must be run from the host WSL terminal using the database runbook. Neither configuration nor host information supplied by the owner is labeled as agent-verified runtime success.

## Primary references

- [PostgreSQL versioning/support](https://www.postgresql.org/support/versioning/).
- [Official PostgreSQL image](https://hub.docker.com/_/postgres).
- [Docker PostgreSQL setup and persistence](https://docs.docker.com/guides/postgresql/immediate-setup-and-data-persistence/).

- [Pinned golang-migrate image Dockerfile](https://github.com/golang-migrate/migrate/blob/v4.20.1/Dockerfile) confirms the Alpine shell and migrate binary used by the local wrapper. The host container is now reported healthy; no registry digest/pull transcript was supplied.

## Owner host verification update

The owner now reports healthy dedicated PostgreSQL 17.11, authenticated TCP access, incorrect-password rejection, non-superuser runtime and disposable migration/integration success. This supersedes prior blocked-status statements above as a host-reported result; the agent did not replay Docker/database operations during report finalization. See the final Phase 0 report for evidence boundaries, remaining clean-clone gate and owner review.
