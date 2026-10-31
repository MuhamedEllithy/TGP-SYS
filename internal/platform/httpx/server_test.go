package httpx

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type pingFunc func(context.Context) error

func (f pingFunc) Ping(ctx context.Context) error { return f(ctx) }
func TestRoutes(t *testing.T) {
	for _, tt := range []struct {
		name, path, method string
		docs, fail         bool
		code               int
	}{
		{"live", "/health/live", "GET", false, true, 200}, {"ready", "/health/ready", "GET", false, false, 200},
		{"db down", "/health/ready", "GET", false, true, 503}, {"docs off", "/docs", "GET", false, false, 404},
		{"contract off", "/openapi.yaml", "GET", false, false, 404}, {"docs", "/docs", "GET", true, false, 200},
		{"contract", "/openapi.yaml", "GET", true, false, 200}, {"method", "/health/live", "POST", false, false, 405},
		{"future not exposed", "/api/v1/assets", "GET", false, false, 404}, {"unknown", "/missing", "GET", false, false, 404},
		{"head", "/health/live", "HEAD", false, false, 200},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var logs bytes.Buffer
			db := pingFunc(func(ctx context.Context) error {
				if _, ok := ctx.Deadline(); !ok {
					t.Error("missing timeout")
				}
				if tt.fail {
					return errors.New("postgres://user:secret@host/db")
				}
				return nil
			})
			h := Handler(db, time.Second, tt.docs, slog.New(slog.NewJSONHandler(&logs, nil)))
			r := httptest.NewRequest(tt.method, tt.path+"?secret=private", nil)
			r.Header.Set("X-Request-ID", "untrusted")
			w := httptest.NewRecorder()
			h.ServeHTTP(w, r)
			if w.Code != tt.code {
				t.Fatalf("got %d: %s", w.Code, w.Body.String())
			}
			if len(w.Header().Get("X-Request-ID")) != 32 {
				t.Fatal("missing server request ID")
			}
			if strings.Contains(w.Body.String()+logs.String(), "secret") || strings.Contains(logs.String(), "private") {
				t.Fatal("secret leaked")
			}
			if tt.method == "HEAD" && w.Body.Len() != 0 {
				t.Fatal("HEAD body")
			}
			if strings.HasPrefix(w.Header().Get("Content-Type"), "application/json") && tt.method != "HEAD" {
				var e envelope
				if err := json.Unmarshal(w.Body.Bytes(), &e); err != nil {
					t.Fatal(err)
				}
				if e.Meta.RequestID != w.Header().Get("X-Request-ID") {
					t.Fatal("request ID mismatch")
				}
				if tt.code >= 400 && (e.Error == nil || e.Error.Details == nil) {
					t.Fatal("error envelope incomplete")
				}
			}
		})
	}
}
func TestReadinessTimeout(t *testing.T) {
	h := Handler(pingFunc(func(ctx context.Context) error { <-ctx.Done(); return ctx.Err() }), time.Millisecond, false, slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil)))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("GET", "/health/ready", nil))
	if w.Code != 503 {
		t.Fatal(w.Code)
	}
}
func TestServerLimits(t *testing.T) {
	s := Server("127.0.0.1:8080", http.NotFoundHandler())
	if s.ReadHeaderTimeout <= 0 || s.WriteTimeout <= 0 || s.IdleTimeout <= 0 || s.MaxHeaderBytes > 16384 {
		t.Fatal("unbounded server")
	}
}
