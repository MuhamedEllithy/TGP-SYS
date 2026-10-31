# Scalar manual browser verification — Phase 0

Result: **PENDING MANUAL VERIFICATION**. No browser evidence has been supplied. The owner must execute this checklist; HTTP HTML availability, OpenAPI lint and unit tests do not establish rendering or Try-it.

## Configuration and safe setup

Current private configuration has `HTTP_ADDR=127.0.0.1:8080`, `APP_ENV=local` and docs enabled. Public URLs derived from that address:

- Documentation: **http://127.0.0.1:8080/docs**
- Contract: **http://127.0.0.1:8080/openapi.yaml**
- Probes: **http://127.0.0.1:8080/health/live** and **http://127.0.0.1:8080/health/ready**

Reuse the running API. If it is not running, run `make run` in the WSL backend directory; do not start/restart databases or execute migrations for this check. In a second WSL terminal confirm:

```sh
curl -i http://127.0.0.1:8080/health/live
curl -i http://127.0.0.1:8080/health/ready
curl -i http://127.0.0.1:8080/docs
curl -i http://127.0.0.1:8080/openapi.yaml
```

Do not print or share `.env` or DB credentials. No authentication is required for these read-only system routes.

## Browser checklist

1. Open the exact documentation URL above in the owner's normal browser. Record browser/version, date, and the source commit or explicitly identify the uncommitted working tree.
2. Confirm the Scalar UI renders, not just the HTML shell. Verify the contract loads from `/openapi.yaml`, and the pinned Scalar 1.76.4 CDN script loads. Inspect the browser Console/Network tab for failed script, schema, CSP or network errors.
3. Confirm exactly the four implemented paths are listed: `/health/live`, `/health/ready`, `/docs`, `/openapi.yaml`, with GET/HEAD operations. There must be no implemented/authenticated asset API claim. No request bodies are required for probes.
4. Inspect `Health`, `Meta` and `Error` response models. Confirm required status/request_id/error fields, the live/ready status values, request ID format, error code/message/details, and readiness's documented 503 response/example. Confirm examples render rather than showing a schema/ref-resolution error.
5. Use Try-it for **GET `/health/live`** with server **http://127.0.0.1:8080**. Expect HTTP **200**, JSON Content-Type, and an actual body of this form (ID changes per request):

   ```json
   {"data":{"status":"live"},"meta":{"request_id":"0123456789abcdef0123456789abcdef"}}
   ```

6. Confirm actual `X-Request-ID` equals `meta.request_id`, consisting of 32 lowercase hexadecimal characters. Confirm UI sent a real network request, not merely displaying the documented example.
7. Use Try-it for **GET `/health/ready`**. Expect HTTP **200** and `data.status="ready"` with its own matching header/body request ID. A 503 with `error.code="not_ready"` is real API/DB unavailability, not a passing readiness result. Do not stop/change the DB to test it.
8. Record the outcome and any sanitized browser errors. Supply explicit successful render/schema/Try-it observations before changing this status to PASS. Do not infer completion from a screenshot of the page title or from curl alone.

## Browser/CORS problem versus API problem

- If the same curl probe fails, investigate API process/bind address or readiness/DB connectivity first. HTTP 503 is an API readiness failure; a connection refusal is an address/process/network issue.
- If curl succeeds and the browser UI is blank, check CDN script loading, extensions, proxy policies and Console CSP errors. The external renderer needs network; the local API may still be healthy.
- If the UI renders but Try-it fails, inspect the actual request URL and Console/Network errors. Keep both docs and selected server on **127.0.0.1:8080**. `localhost` and `127.0.0.1` are different browser origins; using a different alias/port may trigger the server's intentional same-origin CSP or a CORS restriction.
- Browser origin is independent of the WSL-to-Windows host network path. If curl works in WSL but the browser cannot reach the API at all, first check host loopback forwarding rather than weakening API CORS.
- A CORS/CSP message with a successful equivalent same-origin curl is browser-policy evidence, not proof the handler returned an error. Do not add wildcard CORS, disable browser security, loosen CSP, or expose the API externally as a workaround during acceptance.

Static route/schema and handler tests are recorded separately in the final report. Actual model rendering and interactive behavior remain **PENDING MANUAL VERIFICATION**.
