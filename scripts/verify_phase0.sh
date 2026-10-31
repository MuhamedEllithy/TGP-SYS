#!/usr/bin/env bash
# Run in the owner's WSL terminal. Never enable xtrace or print .env.
set -euo pipefail
cd "$(dirname "$0")/.."
make env
make db-up
make setup
# All commands must succeed; blocked checks never become passes.
go mod download
go mod verify
go test ./...
go test -race ./...
make lint
python3 -m unittest discover -s scripts -p 'test_*.py'
make openapi-lint
make db-verify
make migrate-up
make migrate-test-up
make migrate-test-up
make test-integration
make migrate-test-down
# Assert rollback actually removed the probe, with authenticated host-path pgx.
python3 scripts/local_env.py --run go run ./cmd/dbcheck --test --probe-absent
make migrate-test-up
make test-integration
# Check the configured local HTTP port before starting our own process.
python3 - <<'PY'
import sys
sys.path.insert(0, 'scripts')
from local_env import load
import socket
env = load()
assert env.get('APP_ENV') == 'local' and env.get('DOCS_ENABLED') == 'true'
assert env.get('HTTP_ADDR') == '127.0.0.1:8080', 'Smoke runner requires documented default HTTP address'
with socket.socket() as probe:
    probe.bind(('127.0.0.1', 8080))
PY
api_binary=$(mktemp /tmp/tgp-sys-api.XXXXXXXX)
api_log=$(mktemp /tmp/tgp-sys-api-log.XXXXXXXX)
api_pid=''
cleanup() {
    if [ -n "$api_pid" ]; then kill -TERM "$api_pid" 2>/dev/null || true; fi
    rm -f "$api_binary" "$api_log"
}
trap cleanup EXIT
go build -o "$api_binary" ./cmd/api
chmod 700 "$api_binary"
# Launch via Python exec: the recorded PID belongs to the actual API process.
python3 scripts/local_env.py --run "$api_binary" > "$api_log" 2>&1 &
api_pid=$!
for attempt in $(seq 1 30); do
    kill -0 "$api_pid"
    if curl --fail --silent http://127.0.0.1:8080/health/ready >/dev/null; then break; fi
    sleep 1
done
make api-smoke
python3 - <<'PY'
import urllib.request
for path, content in [('/docs', b'@scalar/api-reference@'), ('/openapi.yaml', b'openapi: 3.1.1')]:
    with urllib.request.urlopen('http://127.0.0.1:8080' + path, timeout=5) as response:
        assert response.status == 200 and content in response.read()
    print('PASS:', path, 'available')
PY
kill -TERM "$api_pid"
wait "$api_pid"
api_pid=''
echo 'PASS: host Phase 0 automated checks and graceful shutdown. Browser Scalar Try-it still requires manual verification.'
echo 'Database containers/volumes remain intact. Phase 1 has not started.'
