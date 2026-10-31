// Package httpx contains system HTTP transport only; no business services.
package httpx

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"gateplaza.local/facility/api"
)

type Pinger interface{ Ping(context.Context) error }

type envelope struct {
	Data  any       `json:"data,omitempty"`
	Error *apiError `json:"error,omitempty"`
	Meta  metadata  `json:"meta"`
}
type metadata struct {
	RequestID string `json:"request_id"`
}
type apiError struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Details []string `json:"details"`
}
type health struct {
	Status string `json:"status"`
}

func Handler(db Pinger, timeout time.Duration, docs bool, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		id := make([]byte, 16)
		if _, err := rand.Read(id); err != nil {
			http.Error(w, "request initialization failed", http.StatusInternalServerError)
			return
		}
		requestID := hex.EncodeToString(id) // Never trust or echo an untrusted request ID.
		w.Header().Set("X-Request-ID", requestID)
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Referrer-Policy", "no-referrer")
		status := http.StatusOK
		defer func() {
			logger.Info("http_request", "request_id", requestID, "method", r.Method, "status", status, "duration_ms", time.Since(start).Milliseconds())
		}()
		send := func(code int, e envelope) {
			status = code
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(code)
			if r.Method != http.MethodHead {
				if err := json.NewEncoder(w).Encode(e); err != nil {
					logger.Warn("response_write_failed", "request_id", requestID)
				}
			}
		}
		fail := func(code int, key, message string) {
			send(code, envelope{Error: &apiError{key, message, []string{}}, Meta: metadata{requestID}})
		}
		known := r.URL.Path == "/health/live" || r.URL.Path == "/health/ready" || (docs && (r.URL.Path == "/docs" || r.URL.Path == "/openapi.yaml"))
		if !known {
			fail(404, "not_found", "Resource not found")
			return
		}
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			w.Header().Set("Allow", "GET, HEAD")
			fail(405, "method_not_allowed", "Method not allowed")
			return
		}
		switch r.URL.Path {
		case "/health/live":
			send(200, envelope{Data: health{"live"}, Meta: metadata{requestID}})
		case "/health/ready":
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()
			if db == nil || db.Ping(ctx) != nil {
				fail(503, "not_ready", "Service is not ready")
				return
			}
			send(200, envelope{Data: health{"ready"}, Meta: metadata{requestID}})
		case "/openapi.yaml":
			w.Header().Set("Content-Type", "application/yaml")
			if r.Method != http.MethodHead {
				_, _ = w.Write(api.Contract)
			}
		case "/docs":
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Header().Set("Content-Security-Policy", "default-src 'none'; script-src https://cdn.jsdelivr.net 'unsafe-inline'; style-src 'unsafe-inline'; img-src data:; font-src data:; connect-src 'self'; frame-ancestors 'none'; base-uri 'none'")
			if r.Method != http.MethodHead {
				_, _ = w.Write(api.Docs)
			}
		}
	})
}

func Server(address string, handler http.Handler) *http.Server {
	return &http.Server{Addr: address, Handler: handler, ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second, WriteTimeout: 15 * time.Second, IdleTimeout: 60 * time.Second, MaxHeaderBytes: 16 << 10}
}
