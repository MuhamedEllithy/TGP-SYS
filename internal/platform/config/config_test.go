package config

import (
	"strings"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	for _, tt := range []struct {
		name, key, value string
		bad              bool
	}{
		{"defaults", "", "", false}, {"bad env", "APP_ENV", "staging", true}, {"bad URL", "DATABASE_URL", "postgres://secret:password@%zz/db", true},
		{"missing DB", "DATABASE_URL", "", true}, {"bad duration", "DB_TIMEOUT", "0s", true}, {"large duration", "DB_TIMEOUT", "2m", true},
		{"bad bool", "DOCS_ENABLED", "yes", true}, {"bad level", "LOG_LEVEL", "secret", true}, {"bad port", "HTTP_ADDR", "127.0.0.1:0", true},
		{"public local", "HTTP_ADDR", "0.0.0.0:8080", true}, {"valid timeout", "DB_TIMEOUT", "3s", false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			env := map[string]string{"DATABASE_URL": "postgres://local:private@127.0.0.1/db?sslmode=disable"}
			if tt.key != "" {
				env[tt.key] = tt.value
			}
			c, e := Load(func(k string) string { return env[k] })
			if (e != nil) != tt.bad {
				t.Fatalf("unexpected error: %v", e)
			}
			if e != nil && (strings.Contains(e.Error(), "private") || strings.Contains(e.Error(), "password")) {
				t.Fatal("secret leaked")
			}
			if e == nil && c.DBTimeout <= 0 {
				t.Fatal("invalid timeout")
			}
		})
	}
}
func TestProduction(t *testing.T) {
	env := map[string]string{"APP_ENV": "production", "DATABASE_URL": "postgres://local:private@db/db?sslmode=verify-full"}
	get := func(k string) string { return env[k] }
	if c, e := Load(get); e != nil || c.Docs || c.ShutdownTimeout != 10*time.Second {
		t.Fatalf("invalid production defaults: %v", e)
	}
	env["DOCS_ENABLED"] = "true"
	if _, e := Load(get); e == nil {
		t.Fatal("public production docs allowed")
	}
	env["DOCS_ENABLED"] = "false"
	env["DATABASE_URL"] = "postgres://u:p@db/db?sslmode=require"
	if _, e := Load(get); e == nil {
		t.Fatal("unverified production TLS allowed")
	}
}
