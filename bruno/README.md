# Phase 0 Bruno verification

This collection uses implemented system behavior only. No credentials or database URLs are required. `environments/local.bru` defines `base_url=http://127.0.0.1:8080`, matching the current private `HTTP_ADDR` and OpenAPI server origin. Requests interpolate `{{base_url}}`.

| Request | Expected HTTP / assertions |
|---|---|
| GET `/health/live` | 200; `data.status=live` |
| GET `/health/ready` | 200; `data.status=ready`; actual DB must already be ready |
| GET `/api/v1/assets` | 404; `error.code=not_found`, expected message, empty details, no data |
| POST `/health/live` | 405; `error.code=method_not_allowed`, empty details, no data; `Allow: GET, HEAD` |

Each request also checks JSON Content-Type and a 32-character lowercase-hex `meta.request_id` equal to `X-Request-ID`. The asset request verifies absence of an unimplemented API; it does not implement or advertise an asset endpoint. The unsupported POST is rejected by the probe and never mutates data. A real readiness 503 must fail the expected-ready case; do not fake a DB response, disable that assertion, or stop the development DB to manufacture a negative case.

## Exact WSL execution

If the API is already serving, reuse it. Otherwise run `make run` from the backend directory in a separate terminal. Do not start/recreate containers, migrate or modify existing DBs for this collection.

```sh
cd ~/Projects/TGP-Sys/backend
curl --fail --silent --show-error http://127.0.0.1:8080/health/live
curl --fail --silent --show-error http://127.0.0.1:8080/health/ready
# Isolated CLI installation; no global or project dependency changes:
tgp_bruno_tools=$(mktemp -d /tmp/tgp-bruno-tools.XXXXXX)
tgp_bruno_results=$(mktemp -d /tmp/tgp-bruno-results.XXXXXX)
npm install --prefix "$tgp_bruno_tools" --no-save --ignore-scripts @usebruno/cli@4.0.0
cd bruno
"$tgp_bruno_tools/node_modules/.bin/bru" --version
"$tgp_bruno_tools/node_modules/.bin/bru" run --env local --sandbox safe --reporter-json "$tgp_bruno_results/results.json"
```

Record actual exit status and request/assertion counts. All four requests and their tests must complete successfully before marking live collection execution PASS. Do not use developer sandbox mode or `--insecure`. Do not copy backend `.env` into Bruno. If HTTP_ADDR is intentionally changed, update the public base URL or explicitly pass `--env-var 'base_url=http://127.0.0.1:PORT'` to match it; never include credentials.

## Actual verification status

PASS: CLI 4.0.0 installed in a temporary tools directory, version/help verified, collection/environment parsed and all four expected request URLs attempted. BLOCKED: live execution from this sandbox returned `connect EPERM` for all four; exit 1, zero tests/assertions executed. This is not a passed collection run or evidence of API response failure. No host Bruno result was supplied. See the Phase 0 report for recorded commands and artifacts.

Primary references: [Bruno CLI options](https://docs.usebruno.com/bru-cli/run/options), [response scripting API](https://docs.usebruno.com/testing/script/javascript-reference).
