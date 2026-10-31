#!/usr/bin/env python3
"""Run versioned migrations; down is restricted to the disposable Phase 0 DB."""
import argparse
import os
import subprocess
import sys
from local_env import ROOT, load, url


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("direction", choices=["up", "down"])
    parser.add_argument("--test", action="store_true")
    args = parser.parse_args()
    if args.direction == "down" and not args.test:
        parser.error("down is allowed only with --test against the dedicated disposable database")
    try:
        preflight = subprocess.run([sys.executable, "scripts/db_preflight.py"], cwd=ROOT)
        if preflight.returncode:
            return preflight.returncode
        values = load()
        if values.get("DB_NAME") != "tgp_facility_dev" or values.get("MIGRATION_USER") != "tgp_migrator_dev" or values.get("MIGRATION_DATABASE_URL") != url(values["MIGRATION_USER"], values["MIGRATION_PASSWORD"], "postgres", 5432, "tgp_facility_dev"):
            raise ValueError()
        child_env = dict(os.environ, **values)
        if args.test:
            if values.get("TEST_DB_NAME") != "tgp_facility_phase0_test" or values.get("TEST_MIGRATION_DATABASE_URL") != url(values["MIGRATION_USER"], values["MIGRATION_PASSWORD"], "postgres", 5432, "tgp_facility_phase0_test"):
                raise ValueError()
            # Verify DB identity and refuse foreign schemas/tables before any rollback.
            check = subprocess.run(["go", "run", "./cmd/dbcheck", "--test"], cwd=ROOT, env=child_env)
            if check.returncode:
                return check.returncode
            child_env["MIGRATION_DATABASE_URL"] = values["TEST_MIGRATION_DATABASE_URL"]
        cmd = ["docker", "compose", "--env-file", ".env", "-f", "deployments/local/compose.yaml",
               "run", "--rm", "--no-deps", "-e", "MIGRATION_DATABASE_URL", "migrate", args.direction]
        if args.direction == "down":
            cmd.append("1")
        result = subprocess.run(cmd, cwd=ROOT, env=child_env, capture_output=True)
        if result.returncode:
            print("Migration failed; raw tool output withheld because driver errors can contain credentials.", file=sys.stderr)
            return result.returncode
        print("PASS: migration", args.direction, "on", "disposable Phase 0 database" if args.test else "TGP-Sys development database")
        return 0
    except (KeyError, ValueError, OSError):
        print("Migration configuration unavailable or unsafe; no secret values printed.", file=sys.stderr)
        return 1


if __name__ == "__main__":
    sys.exit(main())
