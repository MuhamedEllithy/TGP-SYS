#!/usr/bin/env python3
"""Read-only safety checks before starting only the dedicated project service."""
import json
import socket
import subprocess
import sys
from pathlib import Path
from local_env import ROOT, load, url

PROJECT = "tgp-sys-local-dev"
CONTAINER = "tgp-sys-postgres-dev"


def docker(*args):
    result = subprocess.run(["docker", *args], cwd=ROOT, capture_output=True, text=True)
    if result.returncode:
        raise ValueError("Docker daemon unavailable to this execution environment; no provisioning attempted")
    return result.stdout.strip()


def main():
    try:
        env = load()
        if env.get("DB_HOST") != "127.0.0.1" or env.get("DB_NAME") != "tgp_facility_dev" or env.get("DB_USER") != "tgp_app_dev":
            raise ValueError("Dedicated local project settings required")
        port = int(env["DB_PORT"])
        if port == 54322 or not 1024 <= port <= 65535:
            raise ValueError("Refusing invalid or unrelated-project port")
        expected = {
            "DATABASE_URL": url(env["DB_USER"], env["DB_PASSWORD"], "127.0.0.1", port, env["DB_NAME"]),
            "MIGRATION_DATABASE_URL": url(env["MIGRATION_USER"], env["MIGRATION_PASSWORD"], "postgres", 5432, env["DB_NAME"]),
            "TEST_DATABASE_URL": url(env["MIGRATION_USER"], env["MIGRATION_PASSWORD"], "127.0.0.1", port, "tgp_facility_phase0_test"),
            "TEST_MIGRATION_DATABASE_URL": url(env["MIGRATION_USER"], env["MIGRATION_PASSWORD"], "postgres", 5432, "tgp_facility_phase0_test"),
        }
        if env.get("TEST_DB_NAME") != "tgp_facility_phase0_test" or any(env.get(k) != v for k, v in expected.items()):
            raise ValueError("Connection settings disagree; no provisioning attempted")
        docker("info", "--format", "{{.ServerVersion}}")
        # Never inspect container environments, including the unrelated Supabase container.
        names = docker("ps", "-a", "--filter", f"name=^/{CONTAINER}$", "--format", "{{.Names}}")
        if names:
            labels = json.loads(docker("inspect", CONTAINER, "--format", "{{json .Config.Labels}}")) or {}
            if labels.get("com.docker.compose.project") != PROJECT or labels.get("com.docker.compose.service") != "postgres" or Path(labels.get("com.docker.compose.project.working_dir", "")).resolve() != ROOT / "deployments/local":
                raise ValueError("Container-name collision or different project directory; refusing to adopt/recreate")
            image = docker("inspect", CONTAINER, "--format", "{{.Config.Image}}")
            mounts = json.loads(docker("inspect", CONTAINER, "--format", "{{json .Mounts}}"))
            data = [m for m in mounts if m.get("Destination") == "/var/lib/postgresql/data"]
            networks = json.loads(docker("inspect", CONTAINER, "--format", "{{json .NetworkSettings.Networks}}"))
            if image != "postgres:17.11-alpine" or len(data) != 1 or data[0].get("Name") != "tgp-sys-postgres17-dev-data" or set(networks) != {"tgp-sys-dev-network"}:
                raise ValueError("Existing project image/data/network configuration differs; refusing to adopt")
            state = docker("inspect", CONTAINER, "--format", "{{.State.Status}}")
            if state == "running":
                bindings = json.loads(docker("inspect", CONTAINER, "--format", "{{json .HostConfig.PortBindings}}"))
                if bindings.get("5432/tcp") != [{"HostIp":"127.0.0.1", "HostPort":str(port)}]:
                    raise ValueError("Existing project container port differs; no recreation allowed")
            else:
                with socket.socket() as probe:
                    probe.bind(("127.0.0.1", port))
        else:
            with socket.socket() as probe:
                probe.bind(("127.0.0.1", port))
        for kind, name in [("volume", "tgp-sys-postgres17-dev-data"), ("network", "tgp-sys-dev-network")]:
            existing = name in docker(kind, "ls", "--format", "{{.Name}}").splitlines()
            if existing:
                labels = json.loads(docker(kind, "inspect", name, "--format", "{{json .Labels}}")) or {}
                if labels.get("com.docker.compose.project") != PROJECT:
                    raise ValueError("Resource-name collision; refusing unrelated volume/network")
        print("PASS: dedicated project association and host port safety; no infrastructure changed")
        return 0
    except (ValueError, KeyError, OSError, json.JSONDecodeError) as error:
        # Messages constructed here contain no env values; OS exception details may not.
        print("Preflight blocked: check Docker daemon access, local settings, port availability and resource ownership. No provisioning attempted; no secrets displayed.", file=sys.stderr)
        return 1


if __name__ == "__main__":
    sys.exit(main())
