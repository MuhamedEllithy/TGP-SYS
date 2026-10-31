#!/usr/bin/env python3
"""Network smoke test of a running local server; readiness must be real."""
import json
import urllib.request
import urllib.error
base = "http://127.0.0.1:8080"
for path, expected in [("/health/live", 200), ("/health/ready", 200), ("/api/v1/assets", 404)]:
    try:
        response = urllib.request.urlopen(base + path, timeout=5)
    except urllib.error.HTTPError as error:
        response = error
    assert response.code == expected, (path, response.code)
    body = json.load(response)
    assert body["meta"]["request_id"] == response.headers["X-Request-ID"]
    print("PASS", path, response.code)
