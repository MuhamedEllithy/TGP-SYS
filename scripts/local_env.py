#!/usr/bin/env python3
"""Generate or load local credentials without printing values or executing .env."""
import argparse
import os
from pathlib import Path
import re
import secrets
import shlex
import socket
import subprocess
import sys
from urllib.parse import quote

ROOT = Path(__file__).resolve().parents[1]
ENV_FILE = ROOT / ".env"


def load(path=None):
    if path is None:
        path = ENV_FILE
    values = {}
    if not path.exists():
        return values
    if path.is_symlink():
        raise ValueError("Refusing a symlink environment file")
    for line in path.read_text().splitlines():
        line = line.strip()
        if not line or line.startswith("#"):
            continue
        if line.startswith("export "):
            line = line[7:]
        key, separator, raw = line.partition("=")
        if not separator or not re.fullmatch(r"[A-Za-z_][A-Za-z0-9_]*", key):
            raise ValueError("Environment file must contain simple KEY=value assignments")
        tokens = shlex.split(raw, comments=True)
        if len(tokens) > 1:
            raise ValueError("Environment values containing spaces must be quoted")
        values[key] = tokens[0] if tokens else ""
    return values


def url(user, password, host, port, name):
    return f"postgres://{quote(user, safe='')}:{quote(password, safe='')}@{host}:{port}/{quote(name, safe='')}?sslmode=disable"


def configure(port=None, select_port=False):
    # Never silently rotate existing credentials or repoint a configured database.
    old = load()
    expected = {"DB_HOST": "127.0.0.1", "DB_NAME": "tgp_facility_dev", "DB_USER": "tgp_app_dev",
                "POSTGRES_USER": "tgp_bootstrap_dev", "MIGRATION_USER": "tgp_migrator_dev",
                "TEST_DB_NAME": "tgp_facility_phase0_test"}
    database_keys = set(expected) | {"DATABASE_URL", "DB_PORT", "DB_PASSWORD", "POSTGRES_PASSWORD", "MIGRATION_PASSWORD", "MIGRATION_DATABASE_URL", "TEST_DATABASE_URL", "TEST_MIGRATION_DATABASE_URL"}
    has_database_settings = any(key in old for key in database_keys)
    if has_database_settings and any(key not in old or old[key] != value for key, value in expected.items()):
        raise ValueError("Existing database settings need review; no values were changed")
    tracked = subprocess.run(["git", "ls-files", "--error-unmatch", ".env"], cwd=ROOT,
                             capture_output=True)
    if tracked.returncode == 0:
        raise ValueError("SECURITY: .env is tracked; no credentials were written")
    ignored = subprocess.run(["git", "check-ignore", "-q", ".env"], cwd=ROOT, capture_output=True)
    if ignored.returncode != 0:
        raise ValueError("Cannot verify .env Git exclusion; no credentials were written")
    values = dict(old)
    for key, value in expected.items():
        values.setdefault(key, value)
    values.setdefault("DB_PORT", "55432")
    for key in ("DB_PASSWORD", "MIGRATION_PASSWORD", "POSTGRES_PASSWORD"):
        if key in values and (not values[key] or values[key].startswith("replace-")):
            raise ValueError("Existing placeholder/empty secret needs review; no rotation performed")
        values.setdefault(key, secrets.token_hex(32))
    if port is not None:
        if not 1024 <= port <= 65535 or port == 54322:
            raise ValueError("Choose an unprivileged port other than the unrelated project's 54322")
        if old and str(port) != old["DB_PORT"]:
            raise ValueError("Existing port will not be changed automatically; review connections first")
        values["DB_PORT"] = str(port)
    # Explicit host-side port selection is allowed only before any project container exists.
    if select_port:
        listing = subprocess.run(["docker", "ps", "-a", "--filter", "name=^/tgp-sys-postgres-dev$", "--format", "{{.Names}}"], cwd=ROOT, capture_output=True)
        if listing.returncode or listing.stdout.strip():
            raise ValueError("Cannot select a new port while Docker is unavailable or a project container exists")
    port_verified = False
    # On a new configuration choose an actually free port when host sockets are available.
    candidates = [] if old and not select_port else ([port] if port is not None else range(55432, 55443))
    for candidate in candidates:
        try:
            with socket.socket() as sock:
                sock.bind(("127.0.0.1", candidate))
            values["DB_PORT"] = str(candidate)
            port_verified = True
            break
        except PermissionError:
            break  # Sandbox limitation: preserve candidate and require host preflight.
        except OSError:
            continue
    else:
        if candidates:
            raise ValueError("No candidate port is free; nothing was changed")
    derived = {
        "DATABASE_URL": url(values["DB_USER"], values["DB_PASSWORD"], values["DB_HOST"], values["DB_PORT"], values["DB_NAME"]),
        "MIGRATION_DATABASE_URL": url(values["MIGRATION_USER"], values["MIGRATION_PASSWORD"], "postgres", 5432, values["DB_NAME"]),
        "TEST_DATABASE_URL": url(values["MIGRATION_USER"], values["MIGRATION_PASSWORD"], values["DB_HOST"], values["DB_PORT"], values["TEST_DB_NAME"]),
        "TEST_MIGRATION_DATABASE_URL": url(values["MIGRATION_USER"], values["MIGRATION_PASSWORD"], "postgres", 5432, values["TEST_DB_NAME"]),
    }
    for key, value in derived.items():
        if key in old and old[key] != value and not (select_port and key in {"DATABASE_URL", "TEST_DATABASE_URL"}):
            raise ValueError("Existing connection URL disagrees with explicit fields; no changes made")
        values[key] = value
    for key, value in {"HTTP_ADDR":"127.0.0.1:8080", "APP_ENV":"local", "DOCS_ENABLED":"true",
                       "LOG_LEVEL":"info", "DB_TIMEOUT":"2s", "SHUTDOWN_TIMEOUT":"10s"}.items():
        values.setdefault(key, value)
    original = ENV_FILE.read_text() if ENV_FILE.exists() else "# Generated local development credentials. Never commit or print this file.\n"
    additions = [key for key in values if key not in old]
    if select_port and old:
        # Preserve comments/unrelated values; change only the port and its two host URLs.
        updated = []
        changed = {"DB_PORT", "DATABASE_URL", "TEST_DATABASE_URL"}
        for line in original.splitlines():
            key = line.removeprefix("export ").partition("=")[0].strip()
            updated.append(f"{key}={shlex.quote(values[key])}" if key in changed else line)
        original = "\n".join(updated)
    content = original.rstrip() + "\n" + "".join(f"{key}={shlex.quote(values[key])}\n" for key in additions)
    if not ENV_FILE.exists():
        fd = os.open(ENV_FILE, os.O_WRONLY | os.O_CREAT | os.O_EXCL, 0o600)
        with os.fdopen(fd, "w") as stream:
            stream.write(content)
    else:
        ENV_FILE.chmod(0o600)
        if additions or select_port:
            ENV_FILE.write_text(content)
    print("Local .env prepared (mode 0600); credentials were not displayed or verified against a server.")
    print("Host port availability checked." if port_verified else "Host port availability is not verified by this command; run make db-preflight on WSL.")


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--port", type=int)
    parser.add_argument("--select-port", action="store_true")
    parser.add_argument("--run", nargs=argparse.REMAINDER)
    args = parser.parse_args()
    try:
        if args.run:
            values = load()
            if not values:
                raise ValueError("Run make env first")
            env = dict(os.environ, **values)
            os.chdir(ROOT)
            os.execvpe(args.run[0], args.run, env)
        configure(args.port, args.select_port)
        return 0
    except (ValueError, OSError):
        print("Local environment operation failed safely; review Git exclusion, existing settings and permissions. No secret values are printed.", file=sys.stderr)
        return 1


if __name__ == "__main__":
    sys.exit(main())
